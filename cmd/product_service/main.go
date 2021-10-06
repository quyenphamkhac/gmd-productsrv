package main

import (
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/opentracing/opentracing-go"
	"github.com/quyenphamkhac/gmd-productsrv/config"
	"github.com/quyenphamkhac/gmd-productsrv/internal/jaeger"
	"github.com/quyenphamkhac/gmd-productsrv/internal/logger"
	"github.com/quyenphamkhac/gmd-productsrv/internal/rabbitmq"
	"github.com/quyenphamkhac/gmd-productsrv/pkg/adapter"
	pb "github.com/quyenphamkhac/gmd-productsrv/pkg/api/v1"
	"github.com/quyenphamkhac/gmd-productsrv/pkg/handler"
	"github.com/quyenphamkhac/gmd-productsrv/pkg/usecase"
	"google.golang.org/grpc"
)

func main() {

	log.Println("starting product service")

	config, err := config.GetConfig()
	if err != nil {
		log.Fatalf("loading config: %v", err)
	}
	rabbitmqConn, err := rabbitmq.NewRabbitMQConn(config)
	if err != nil {
		log.Fatalf("connect rabbitmq: %v", err)
	}
	defer rabbitmqConn.Close()

	rabbitmqCh, err := rabbitmqConn.Channel()
	if err != nil {
		log.Fatalf("create rabbitmq channel: %v", err)
	}
	defer rabbitmqCh.Close()

	srvLogger, err := logger.NewServiceLogger(config)
	srvLogger.InitLogger()
	if err != nil {
		log.Fatalf("unable init logger: %v", err)
	}

	tracer, closer, err := jaeger.InitJeagerTracing(config)
	if err != nil {
		log.Fatalf("could not init jaeger tracer: %v", err)
	}
	log.Println("jaeger conntected")

	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()
	log.Println("opentracing conntected")

	grpcServer := grpc.NewServer()
	mockAdapter := adapter.NewMockAdaper()
	productUsecase := usecase.NewProductUseCase(mockAdapter)
	productSrvHandler := handler.NewProductService(productUsecase, srvLogger)
	pb.RegisterProductSrvServer(grpcServer, productSrvHandler)

	port := config.Service.Port
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	go func() {
		log.Printf("start grpc server port: %s\n", port)
		grpcServer.Serve(lis)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("stopping grpc server...")
	grpcServer.GracefulStop()
}
