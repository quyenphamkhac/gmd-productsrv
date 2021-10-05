package main

import (
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/quyenphamkhac/gmd-productsrv/config"
	"github.com/quyenphamkhac/gmd-productsrv/internal/rabbitmq"
	"github.com/quyenphamkhac/gmd-productsrv/pkg/adapter"
	pb "github.com/quyenphamkhac/gmd-productsrv/pkg/api/v1"
	"github.com/quyenphamkhac/gmd-productsrv/pkg/handler"
	"github.com/quyenphamkhac/gmd-productsrv/pkg/usecase"
	"google.golang.org/grpc"
)

func main() {

	log.Println("starting service")

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

	grpcServer := grpc.NewServer()
	mockAdapter := adapter.NewMockAdaper()
	productUsecase := usecase.NewProductUseCase(mockAdapter)
	productSrvHandler := handler.NewProductService(productUsecase)
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
