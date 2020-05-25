package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	. "github.com/tianwaizhiyin/go-grpc-client.git/services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io"
	"io/ioutil"
	"log"
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
	//t:=timestamp.Timestamp{Seconds:time.Now().Unix()}
	//orderClient := NewOrderServiceClient(conn)
	//res, _ := orderClient.NewOrder(ctx, &OrderMain{
	//	OrderId:1001,
	//	OrderNo:"20190809",
	//	OrderMoney:90,
	//	OrderTime:&t,
	//})
	//fmt.Println(res)

	userClient:= NewUserServiceClient(conn)

	//服务端流模式，客户端代码
	//req := UserScoreRequest{}
	//req.Users=make([]*UserInfo,0)
	//for i:=1; i<=6; i++ {
	//	var i int32
	//	req.Users=append(req.Users, &UserInfo{UserId:i})
	//}
	//stream, err := userClient.GetUserScoreByServerStream(ctx, &req)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//for {
	//	res, err := stream.Recv()
	//	if err == io.EOF {
	//		break
	//	}
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	fmt.Println(res.Users)
	//}
	//fmt.Println(res.Users)

	//客户端流模式，客户端代码
	var i int32
	if err != nil {
		log.Fatal(err)
	}
	//stream, err := userClient.GetUserScoreByClientStream(ctx)
	stream, err := userClient.GetUserScoreByTWS(ctx)
	if err != nil {
		log.Fatal(err)
	}
	var uid int32=1
	for j:=1; j <=3; j++ {
		req:=UserScoreRequest{}
		req.Users=make([]*UserInfo,0)
		for i=1; i < 6; i++ { //加了5条用户数据
			req.Users=append(req.Users, &UserInfo{UserId:uid})
			uid++
		}
		err := stream.Send(&req)
		if err != nil {
			log.Println(err)
		}
		res,err := stream.Recv()
		if err == io.EOF{
			break
		}
		if err != nil {
			log.Println(err)
		}
		fmt.Println(res.Users)
	}
	//res, _ := stream.CloseAndRecv()
	//fmt.Println(res.Users)


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