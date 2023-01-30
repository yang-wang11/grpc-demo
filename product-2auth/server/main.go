package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"mygrpc/pkg/product/product"
	"net"

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

func main() {
	// 证书认证-双向认证
	// 从证书相关文件中读取和解析信息，得到证书公钥、密钥对
	cert, err := tls.LoadX509KeyPair("../../certs/server.pem", "../../certs/server-key.pem")
	if err != nil {
		log.Fatal("证书读取错误", err)
	}
	// 创建一个新的、空的 CertPool 用于client CA存放
	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile("../../certs/ca.pem")
	if err != nil {
		log.Fatal("ca证书读取错误", err)
	}
	// 尝试解析所传入的 PEM 编码的证书。如果解析成功会将其加到 CertPool 中，便于后面的使用
	certPool.AppendCertsFromPEM(ca)
	// 构建基于 TLS 的 TransportCredentials 选项
	creds := credentials.NewTLS(&tls.Config{
		// 设置证书链，允许包含一个或多个
		Certificates: []tls.Certificate{cert},
		// 要求必须校验客户端的证书。可以根据实际情况选用以下参数
		ClientAuth: tls.RequireAndVerifyClientCert,
		// 设置根证书的集合，校验方式使用 ClientAuth 中设定的模式
		ClientCAs: certPool,
	})

	// server := grpc.NewServer()
	server := grpc.NewServer(grpc.Creds(creds))
	product.RegisterProdServiceServer(server, NewProductServer())
	listener, err := net.Listen("tcp", ":8000")
	if err != nil {
		panic(err)
	}
	server.Serve(listener)
}
