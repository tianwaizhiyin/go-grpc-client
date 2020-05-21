package main

import (
	"context"
	"fmt"
	"github.com/tianwaizhiyin/go-grpc-client.git/services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
)

func main()  {
	creds, err := credentials.NewClientTLSFromFile("keys/server.crt", "dongfangfuli.com")
	if err != nil {
		log.Fatal(err)
	}

	conn, err := grpc.Dial(":8080", grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	prodClient := services.NewProdServiceClient(conn)
	prodRes, err := prodClient.GetProdStock(context.Background(), &services.ProdRequest{ProdId:12})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(prodRes.ProdStock)

}