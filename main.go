package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/golang/protobuf/ptypes/timestamp"
	. "github.com/tianwaizhiyin/go-grpc-client.git/services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"log"
	"time"
)

func main()  {
	//creds, err := credentials.NewClientTLSFromFile("keys/server.crt", "dongfangfuli.com")
	//if err != nil {
	//	log.Fatal(err)
	//}

	//tls证书认证
	cert, _ := tls.LoadX509KeyPair("cert/client.pem", "cert/client.key")
	certPool := x509.NewCertPool()
	ca, _ := ioutil.ReadFile("cert/ca.pem")
	certPool.AppendCertsFromPEM(ca)

	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert}, //客户端证书
		ServerName:    "localhost",
		RootCAs:       certPool,
	})

	conn, err := grpc.Dial(":8081", grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	prodClient := NewProdServiceClient(conn)
	ctx := context.Background()
	t:=timestamp.Timestamp{Seconds:time.Now().Unix()}
	orderClient := NewOrderServiceClient(conn)
	res, _ := orderClient.NewOrder(ctx, &OrderMain{
		OrderId:1001,
		OrderNo:"20190809",
		OrderMoney:90,
		OrderTime:&t,
	})
	fmt.Println(res)
	//获取单个商品
	//prodRes, err := prodClient.GetProdStock(ctx,
	//	&ProdRequest{ProdId:12, ProdArea:ProdAreas_C}) //获取商品库存

	//获取商品模型
	prod, err := prodClient.GetProdSInfo(ctx,&ProdRequest{ProdId:12})

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(prod.ProdName)
	//fmt.Println(prodRes.ProdStock)

	//获取多个商品库存
	response, err := prodClient.GetProdStocks(ctx,
		&QuerySize{Size:10})
	if err != nil {
		log.Fatal(err)
	}
	//获取所有
	fmt.Println(response.Prodres)
	//获取单个
	fmt.Println(response.Prodres[2].ProdStock)



}