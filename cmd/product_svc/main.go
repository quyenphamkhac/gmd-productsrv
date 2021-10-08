package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/labstack/echo/v4"
	"github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/quyenphamkhac/gmd-productsrv/config"
	"github.com/quyenphamkhac/gmd-productsrv/internal/interceptors"
	"github.com/quyenphamkhac/gmd-productsrv/internal/jaeger"
	"github.com/quyenphamkhac/gmd-productsrv/internal/logger"
	"github.com/quyenphamkhac/gmd-productsrv/internal/metrics"
	"github.com/quyenphamkhac/gmd-productsrv/internal/rabbitmq"
	"github.com/quyenphamkhac/gmd-productsrv/pkg/adapter"
	pb "github.com/quyenphamkhac/gmd-productsrv/pkg/api/v1"
	"github.com/quyenphamkhac/gmd-productsrv/pkg/handler"
	"github.com/quyenphamkhac/gmd-productsrv/pkg/usecase"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatalf("loading config: %v", err)
	}

	svcLogger, err := logger.NewServiceLogger(cfg)
	svcLogger.InitLogger()
	if err != nil {
		log.Fatalf("unable init logger: %v", err)
	}

	rabbitmqConn, err := rabbitmq.NewRabbitMQConn(cfg)
	if err != nil {
		svcLogger.Fatalf("connect rabbitmq: %v", err)
	}
	defer rabbitmqConn.Close()

	rabbitmqCh, err := rabbitmqConn.Channel()
	if err != nil {
		svcLogger.Fatalf("create rabbitmq channel: %v", err)
	}
	defer rabbitmqCh.Close()

	tracer, closer, err := jaeger.InitJaegerTracing(cfg)
	if err != nil {
		svcLogger.Fatalf("could not init jaeger tracer: %v", err)
	}
	svcLogger.Info("jaeger conntected")

	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()
	svcLogger.Info("opentracing conntected")

	metr, err := metrics.NewPrometheusMetrics(cfg.Metrics.Url, cfg.Metrics.ServiceName)
	if err != nil {
		svcLogger.Fatalf("could not init metrics connection: %v", err)
	}
	svcLogger.Info("prometheus metrics collector conntected")

	im := interceptors.NewInterceptorManager(svcLogger, cfg, metr)

	ctx, cancel := context.WithCancel(context.Background())

	router := echo.New()
	router.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	go func() {
		if err := router.Start(cfg.Metrics.Url); err != nil {
			svcLogger.Errorf("router.Start metrics: %v", err)
			cancel()
		}
	}()

	grpcServer := grpc.NewServer(grpc.KeepaliveParams(keepalive.ServerParameters{
		MaxConnectionIdle: cfg.Service.MaxConnectionIdle * time.Minute,
		Timeout:           cfg.Service.Timeout * time.Second,
		MaxConnectionAge:  cfg.Service.MaxConnectionAge * time.Minute,
		Time:              cfg.Service.Time * time.Minute,
	}),
		grpc.UnaryInterceptor(im.Logger),
		grpc.ChainUnaryInterceptor(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_prometheus.UnaryServerInterceptor,
			grpcrecovery.UnaryServerInterceptor(),
		),
	)

	mockDb := adapter.NewMockAdaper()
	productUsecase := usecase.NewProductUseCase(mockDb)
	productSvc := handler.NewProductService(productUsecase, svcLogger)
	pb.RegisterProductSrvServer(grpcServer, productSvc)
	grpc_prometheus.Register(grpcServer)
	svcLogger.Info("product service initialized")

	port := cfg.Service.Port
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		svcLogger.Fatalf("failed to listen: %v", err)
	}

	go func() {
		svcLogger.Info("start grpc server port: ", port)
		grpcServer.Serve(lis)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	select {
	case v := <-quit:
		svcLogger.Errorf("signal.Notify: %v", v)
	case done := <-ctx.Done():
		svcLogger.Errorf("ctx.Done: %v", done)
	}

	if err := router.Shutdown(ctx); err != nil {
		svcLogger.Errorf("Metrics router.Shutdown: %v", err)
	}
	grpcServer.GracefulStop()
	svcLogger.Info("product service exited properly")
}
