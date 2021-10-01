package main

import (
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/quyenphamkhac/gmd-productsrv/internal/driver"
	"github.com/quyenphamkhac/gmd-productsrv/pkg/adapter"
	pb "github.com/quyenphamkhac/gmd-productsrv/pkg/api/v1"
	"github.com/quyenphamkhac/gmd-productsrv/pkg/handler"
	"github.com/quyenphamkhac/gmd-productsrv/pkg/usecase"
	"google.golang.org/grpc"
)

func main() {

	rabbitmqConn := driver.NewRabbitMQConn()
	defer rabbitmqConn.Close()

	rabbitmqCh, err := rabbitmqConn.Channel()
	if err != nil {
		log.Fatal("Failed to create rabbitmq channel")
	}
	defer rabbitmqCh.Close()

	grpcServer := grpc.NewServer()
	mockAdapter := adapter.NewMockAdaper()
	productUsecase := usecase.NewProductUseCase(mockAdapter)
	productSrvHandler := handler.NewProductService(productUsecase)
	pb.RegisterProductSrvServer(grpcServer, productSrvHandler)

	port := os.Getenv("PORT")
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	go func(serverPort string) {
		log.Printf("Start grpc server port: %s\n", serverPort)
		grpcServer.Serve(lis)
	}(port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Stopping grpc server...")
	grpcServer.GracefulStop()
}
