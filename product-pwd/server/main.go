package main

import (
	"context"
	"fmt"
	"mygrpc/pkg/product/product"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type ProductServer struct {
	product.UnimplementedProdServiceServer
}

func NewProductServer() *ProductServer {
	return &ProductServer{}
}

func (p *ProductServer) GetProductStock(ctx context.Context, request *product.ProductRequest) (*product.ProductResponse, error) {
	return &product.ProductResponse{
		ProdStock: request.ProdId,
	}, nil
}

func main() {
	creds, err := credentials.NewServerTLSFromFile("../../certs/server.pem", "../../certs/server-key.pem")
	if err != nil {
		panic(err)
	}

	authInterceptor := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, fmt.Errorf("missing credentials")
		}

		result := false
		if user, ok := md["user"]; ok {
			if passwd, ok := md["password"]; ok {
				if user[0] != "admin" || passwd[0] != "admin" {
					status.Errorf(codes.Unauthenticated, "token is illegal")
				} else {
					result = true
				}
			}
		}

		if !result {
			return nil, fmt.Errorf("token is illegal")
		}

		return handler(ctx, req)
	}

	server := grpc.NewServer(grpc.Creds(creds), grpc.UnaryInterceptor(authInterceptor))
	product.RegisterProdServiceServer(server, NewProductServer())
	listener, err := net.Listen("tcp", ":8000")
	if err != nil {
		panic(err)
	}
	server.Serve(listener)
}
