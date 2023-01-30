package main

import (
	"context"
	"fmt"
	"mygrpc/pkg/product/product"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {

	creds, err := credentials.NewClientTLSFromFile("../../certs/server.pem", "grpc.server")
	if err != nil {
		panic(err)
	}
	conn, err := grpc.Dial(":8000", grpc.WithTransportCredentials(creds),
		grpc.WithPerRPCCredentials(NewPasswordRPCCredentials("admin", "admin")))
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	client := product.NewProdServiceClient(conn)
	resp, err := client.GetProductStock(context.Background(), &product.ProductRequest{
		ProdId: 20,
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("调用gRPC方法成功，ProdStock = ", resp.ProdStock)
}

type PasswordRPCCredentials struct {
	User     string
	Password string
}

func NewPasswordRPCCredentials(u string, p string) *PasswordRPCCredentials {
	return &PasswordRPCCredentials{
		User:     u,
		Password: p,
	}
}

func (p *PasswordRPCCredentials) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"user":     p.User,
		"password": p.Password,
	}, nil
}

func (p *PasswordRPCCredentials) RequireTransportSecurity() bool {
	return false
}
