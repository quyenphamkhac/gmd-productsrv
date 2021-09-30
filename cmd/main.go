package main

import (
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/quyenphamkhac/gmd-productsrv/pkg/adapter"
	"github.com/quyenphamkhac/gmd-productsrv/pkg/service"
	"github.com/quyenphamkhac/gmd-productsrv/pkg/usecase"
	pb "github.com/quyenphamkhac/gmd-productsrv/protos"
	"google.golang.org/grpc"
)

func main() {
	grpcServer := grpc.NewServer()
	mockAdapter := adapter.NewMockAdaper()
	productUsecase := usecase.NewProductUseCase(mockAdapter)
	productSrv := service.NewProductService(productUsecase)
	pb.RegisterProductSrvServer(grpcServer, productSrv)

	port := os.Getenv("PORT")
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer.Serve(lis)
	go func(serverPort string) {
		log.Printf("Start grpc server port: %s\n", serverPort)
	}(port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Stopping grpc server...")
	grpcServer.GracefulStop()
}
