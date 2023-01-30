package main

import (
	"context"
	"fmt"
	"log"
	"mygrpc/pkg/product/product"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {

	creds, err := credentials.NewClientTLSFromFile("../../certs/server.pem", "grpc.server")
	if err != nil {
		panic(err)
	}
	// conn, err := grpc.Dial(":8000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.Dial(":8000", grpc.WithTransportCredentials(creds))
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	client := product.NewProdServiceClient(conn)
	// -----------------------------无流
	// resp, err := client.GetProductStock(context.Background(), &product.ProductRequest{
	// 	ProdId: 20,
	// })
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println("调用gRPC方法成功，ProdStock = ", resp.ProdStock)

	// -----------------------------客户端流（客户端发送）
	// steamClient, err := client.UpdateProductStockClient(context.Background())
	// if err != nil {
	// 	panic(err)
	// }
	// stop := make(chan struct{})
	// go func(stream product.ProdService_UpdateProductStockClientClient, ch chan struct{}) {
	// 	count := 0
	// 	for {
	// 		if err := stream.Send(&product.ProductRequest{ProdId: 1}); err != nil {
	// 			if err == io.EOF {
	// 				stop <- struct{}{}
	// 			}
	// 			panic(err)
	// 		}
	// 		count++
	// 		if count > 10 {
	// 			stop <- struct{}{}
	// 		}
	// 	}
	// }(steamClient, stop)
	// select {
	// case <-stop:
	// 	fmt.Println("received the stop signal")
	// 	if resp, err := steamClient.CloseAndRecv(); err != nil {
	// 		panic(err)
	// 	} else {
	// 		fmt.Println(resp.ProdStock)
	// 	}
	// }

	// -----------------------------服务端流（服务端发送）
	// streamClient, err := client.UpdateProductStockServer(context.Background(), &product.ProductRequest{
	// 	ProdId: 123,
	// })
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// for {
	// 	resp, err := streamClient.Recv()
	// 	if err != nil {
	// 		if err == io.EOF {
	// 			streamClient.CloseSend()
	// 			fmt.Println("服务端数据接受完成")
	// 			break
	// 		}
	// 		log.Fatal(err)
	// 	}
	// 	fmt.Printf("获取到信息 %d\n", resp.ProdStock)
	// }

	// -----------------------------双向流（服务端&客户端相互发送）

	streamClient, err := client.UpdateProductStockBidirect(context.Background())
	if err != nil {
		panic(err)
	}

	for {
		streamClient.Send(&product.ProductRequest{
			ProdId: 123,
		})
		req, err := streamClient.Recv()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("客户端收到id: %d\n", req.ProdStock)
	}
}
