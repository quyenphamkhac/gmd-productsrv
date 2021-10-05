package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/quyenphamkhac/gmd-productsrv/pkg/api/v1"
	"google.golang.org/grpc"
)

const (
	address = "localhost:8080"
)

func main() {
	fmt.Println("starting grpc client")

	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewProductSrvClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.GetById(ctx, &pb.GetByIdRequest{Id: 1})

	if err != nil {
		log.Fatalf("could not get product by id: %v", err)
	}
	log.Println("product: ", r.GetData())
}
