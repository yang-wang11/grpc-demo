package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"mygrpc/pkg/product/product"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/credentials/insecure"

	"google.golang.org/grpc"
)

const etcdUrl = "http://localhost:2379"
const serviceName = "grpc_server"

func main() {

	// server address
	var ProdId int
	flag.IntVar(&ProdId, "stock_id", 1000, "stock_id")
	flag.Parse()

	etcdClient, err := clientv3.NewFromURL(etcdUrl)
	if err != nil {
		panic(err)
	}
	etcdResolver, err := resolver.NewBuilder(etcdClient)

	conn, err := grpc.Dial(fmt.Sprintf("etcd:///%s", serviceName),
		grpc.WithResolvers(etcdResolver),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)))
	if err != nil {
		fmt.Printf("err: %v", err)
		return
	}

	grpcClient := product.NewProdServiceClient(conn)

	for {
		if helloRespone, err := grpcClient.GetProductStock(context.Background(), &product.ProductRequest{ProdId: int32(ProdId)}); err != nil {
			fmt.Printf("err: %v", err)
			return
		} else {
			log.Println(helloRespone, err)
		}

		time.Sleep(500 * time.Millisecond)
	}

}
