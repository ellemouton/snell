package main

import (
	"database/sql"
	"flag"
	"fmt"
	"time"

	"github.com/ellemouton/snell/db"
	"github.com/ellemouton/snell/lnd"
	"go.etcd.io/etcd/clientv3"
)

var etcdHost = flag.String("etcd_host", "localhost:2379", "etcd host")
var etcdUser = flag.String("etcd_user", "", "etcd user")
var etcdPassword = flag.String("etcd_password", "", "etcd password")

type State struct {
	db         *sql.DB
	etcdClient *clientv3.Client
	lndClient  lnd.Client
}

func NewState() (*State, error) {
	s := &State{}

	db, err := db.Connect()
	if err != nil {
		return nil, fmt.Errorf("problem connecting to db: %s", err)
	}
	s.db = db

	ec, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{*etcdHost},
		DialTimeout: 5 * time.Second,
		Username:    *etcdUser,
		Password:    *etcdPassword,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to connect to etcd: %v", err)
	}
	s.etcdClient = ec

	return s, nil
}

func (s *State) GetDB() *sql.DB {
	return s.db
}

func (s *State) GetEtcdClient() *clientv3.Client {
	return s.etcdClient
}
