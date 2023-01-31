package main

import (
	"context"
	"flag"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"server/pb"
	"syscall"
)

func main() {

	// server address
	var port int
	flag.IntVar(&port, "port", 8001, "port")
	flag.Parse()
	addr := fmt.Sprintf("localhost:%d", port)

	// etcd client
	etcdClient, err := NewEtcdClient()
	if err != nil {
		panic(err)
	}

	// unregister server
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)

	go func(client *clientv3.Client) {
		<-ch
		if err := etcdUnRegister(client, addr); err != nil {
			log.Fatalf("failed to unregister %s\n", addr)
		}
		os.Exit(0)
	}(etcdClient)

	// register server
	err = etcdRegister(etcdClient, addr)
	if err != nil {
		log.Fatalf("failed to register %s\n", addr)
	}

	// create the listener
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	// start the grpc server
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(UnaryInterceptor()))
	pb.RegisterServerServer(grpcServer, Server{})
	log.Printf("service start port %d\n", port)
	if err = grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}

func UnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		log.Printf("call %s\n", info.FullMethod)
		resp, err = handler(ctx, req)
		return resp, err
	}
}
