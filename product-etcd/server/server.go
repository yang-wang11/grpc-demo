package main

import (
	"context"
	"mygrpc/pkg/product/product"
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
