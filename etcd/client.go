package etcd

import (
	"context"
	"flag"
	"fmt"
	"time"

	"go.etcd.io/etcd/clientv3"
)

var etcdHost = flag.String("etcd_host", "localhost:2379", "etcd host")
var etcdUser = flag.String("etcd_user", "", "etcd user")
var etcdPassword = flag.String("etcd_password", "", "etcd password")

type client struct {
	*clientv3.Client
}

type Client interface {
	Close() error
	Put(ctx context.Context, key, val string, opts ...clientv3.OpOption) (*clientv3.PutResponse, error)
	Get(ctx context.Context, key string, opts ...clientv3.OpOption) (*clientv3.GetResponse, error)
	Delete(ctx context.Context, key string, opts ...clientv3.OpOption) (*clientv3.DeleteResponse, error)
}

func New() (Client, error) {
	ec, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{*etcdHost},
		DialTimeout: 5 * time.Second,
		Username:    *etcdUser,
		Password:    *etcdPassword,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to connect to etcd: %v", err)
	}

	return &client{
		ec,
	}, nil
}
