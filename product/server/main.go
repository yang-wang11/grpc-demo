package main

import (
	"context"
	"fmt"
	"io"
	"mygrpc/pkg/product/product"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
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

// 客户端发送stream
func (p *ProductServer) UpdateProductStockClient(stream product.ProdService_UpdateProductStockClientServer) error {
	count := 0
	for {
		msg, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		count++
		fmt.Printf("第%d次服务端接受到 %d\n", count, msg.ProdId)
		if count > 10 {
			time.Sleep(time.Second)
			err = stream.SendAndClose(&product.ProductResponse{
				ProdStock: int32(count),
			})
			if err != nil {
				return err
			}
			return nil
		}
	}
}

// 服务端发送stream
func (p *ProductServer) UpdateProductStockServer(request *product.ProductRequest, stream product.ProdService_UpdateProductStockServerServer) error {
	count := 0
	fmt.Println("UpdateProductStockServer called")
	for {
		if err := stream.Send(&product.ProductResponse{
			ProdStock: request.ProdId,
		}); err != nil {
			return err
		}
		count++
		if count > 100 {
			return nil
		}
	}
}

func (p *ProductServer) UpdateProductStockBidirect(stream product.ProdService_UpdateProductStockBidirectServer) error {
	for {
		request, err := stream.Recv()
		if err != nil {
			return err
		}
		fmt.Println("服务端获取的id", request.ProdId)
		if err := stream.Send(&product.ProductResponse{
			ProdStock: request.ProdId,
		}); err != nil {
			return err
		}

	}
}

func (p *ProductServer) mustEmbedUnimplementedProdServiceServer() {}

func main() {
	creds, err := credentials.NewServerTLSFromFile("../../certs/server.pem", "../../certs/server-key.pem")
	if err != nil {
		panic(err)
	}
	// server := grpc.NewServer()
	server := grpc.NewServer(grpc.Creds(creds))
	product.RegisterProdServiceServer(server, NewProductServer())
	listener, err := net.Listen("tcp", ":8000")
	if err != nil {
		panic(err)
	}
	server.Serve(listener)
}
