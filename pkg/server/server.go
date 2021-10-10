package server

import (
	"context"
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
	"github.com/quyenphamkhac/gmd-productsrv/pkg/adapter"
	pb "github.com/quyenphamkhac/gmd-productsrv/pkg/api/v1"
	"github.com/quyenphamkhac/gmd-productsrv/pkg/handler"
	"github.com/quyenphamkhac/gmd-productsrv/pkg/usecase"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

type Server struct {
	logger   logger.Logger
	cfg      *config.Config
	amqpConn *amqp.Connection
}

func NewProductServer(logger logger.Logger, cfg *config.Config, amqpConn *amqp.Connection) *Server {
	return &Server{
		amqpConn: amqpConn,
		logger:   logger,
		cfg:      cfg,
	}
}

func (s *Server) Run() error {
	rabbitmqCh, err := s.amqpConn.Channel()
	if err != nil {
		s.logger.Fatalf("create rabbitmq channel: %v", err)
	}
	defer rabbitmqCh.Close()

	tracer, closer, err := jaeger.InitJaegerTracing(s.cfg)
	if err != nil {
		s.logger.Fatalf("could not init jaeger tracer: %v", err)
	}
	s.logger.Info("jaeger conntected")

	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()
	s.logger.Info("opentracing conntected")

	metr, err := metrics.NewPrometheusMetrics(s.cfg.Metrics.Url, s.cfg.Metrics.ServiceName)
	if err != nil {
		s.logger.Fatalf("could not init metrics connection: %v", err)
	}
	s.logger.Info("prometheus metrics collector conntected")

	im := interceptors.NewInterceptorManager(s.logger, s.cfg, metr)

	ctx, cancel := context.WithCancel(context.Background())

	router := echo.New()
	router.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	go func() {
		if err := router.Start(s.cfg.Metrics.Url); err != nil {
			s.logger.Errorf("router.Start metrics: %v", err)
			cancel()
		}
	}()

	grpcServer := grpc.NewServer(grpc.KeepaliveParams(keepalive.ServerParameters{
		MaxConnectionIdle: s.cfg.Service.MaxConnectionIdle * time.Minute,
		Timeout:           s.cfg.Service.Timeout * time.Second,
		MaxConnectionAge:  s.cfg.Service.MaxConnectionAge * time.Minute,
		Time:              s.cfg.Service.Time * time.Minute,
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
	productSvc := handler.NewProductService(productUsecase, s.logger)
	pb.RegisterProductSrvServer(grpcServer, productSvc)
	grpc_prometheus.Register(grpcServer)
	s.logger.Info("product service initialized")

	servicePort := s.cfg.Service.Port
	lis, err := net.Listen("tcp", ":"+servicePort)
	if err != nil {
		s.logger.Fatalf("failed to listen: %v", err)
	}

	go func() {
		s.logger.Info("start grpc server port: ", servicePort)
		grpcServer.Serve(lis)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	select {
	case v := <-quit:
		s.logger.Errorf("signal.Notify: %v", v)
	case done := <-ctx.Done():
		s.logger.Errorf("ctx.Done: %v", done)
	}

	if err := router.Shutdown(ctx); err != nil {
		s.logger.Errorf("Metrics router.Shutdown: %v", err)
	}
	grpcServer.GracefulStop()
	s.logger.Info("product service exited properly")
	return err
}
