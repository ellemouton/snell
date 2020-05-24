package macaroon

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/lightningnetwork/lnd/lntypes"
	"go.etcd.io/etcd/clientv3"
	"gopkg.in/macaroon.v2"
)

var etcdHost = flag.String("etcd_host", "localhost:2379", "etcd host")
var etcdUser = flag.String("etcd_user", "", "etcd user")
var etcdPassword = flag.String("etcd_password", "", "etcd password")

type client struct {
	etcdClient *clientv3.Client
}

func New() (Client, error) {
	mac := new(client)

	ec, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{*etcdHost},
		DialTimeout: 5 * time.Second,
		Username:    *etcdUser,
		Password:    *etcdPassword,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to connect to etcd: %v", err)
	}
	mac.etcdClient = ec

	return mac, nil
}

func (c *client) Close() error {
	return c.etcdClient.Close()
}

func (c *client) Create(paymentHash lntypes.Hash, resourceType string, resourceID int64) (*macaroon.Macaroon, error) {
	var rootKey [32]byte
	copy(rootKey[:], generateRandomBytes(32))

	var userID [32]byte
	copy(userID[:], generateRandomBytes(32))

	id := &identifier{
		version:     0,
		paymentHash: paymentHash,
		userID:      userID,
	}

	idBytes, err := id.encode()
	if err != nil {
		return nil, err
	}

	idHash := sha256.Sum256(idBytes)

	idKey := strings.Join([]string{"snell", "secrets", hex.EncodeToString(idHash[:])}, "/")

	// Store the key-value pair
	c.etcdClient.Put(context.TODO(), idKey, string(rootKey[:]))

	_, err = macaroon.New(rootKey[:], idBytes, "snell", macaroon.LatestVersion)
	if err != nil {
		return nil, err
	}

	//if err := mac.AddFirstPartyCaveat(); err != nil {
	//	return nil, err
	//}

	// add first party caveat to Macaroon. (condition = "resource Type", value = resourceID)
	// return Mac (can convert to bytes later with .MarshalBinary(), and then encode to base64 string.)
	return nil, nil
}

func generateRandomBytes(n int64) []byte {
	b := make([]byte, n)
	rand.Read(b)
	return b
}
