package main

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
	"log"
)

const etcdUrl = "http://localhost:2379"
const serviceName = "grpc_server"
const ttl = 10

func NewEtcdClient() (*clientv3.Client, error) {

	return clientv3.NewFromURL(etcdUrl)
}

func etcdRegister(client *clientv3.Client, addr string) error {

	log.Printf("Register %s on etcd.\n", addr)

	ctx := context.Background()

	em, err := endpoints.NewManager(client, serviceName)
	if err != nil {
		return err
	}

	lease, _ := client.Grant(ctx, ttl)

	if err = em.AddEndpoint(ctx, fmt.Sprintf("%s/%s", serviceName, addr), endpoints.Endpoint{Addr: addr}, clientv3.WithLease(lease.ID)); err != nil {
		return err
	}

	go func(c *clientv3.Client) {
		alive, err := c.KeepAlive(ctx, lease.ID)
		if err != nil {
			log.Fatalln(err)
		}
		for {
			<-alive
			fmt.Println("etcd server keep alive")
		}
	}(client)

	return nil
}

func etcdUnRegister(client *clientv3.Client, addr string) error {

	log.Printf("UnRegister %s on etcd.\n", addr)

	ctx := context.Background()

	em, err := endpoints.NewManager(client, serviceName)
	if err != nil {
		return err
	}
	if err = em.DeleteEndpoint(ctx, fmt.Sprintf("%s/%s", serviceName, addr)); err != nil {
		return err
	}

	return nil
}
