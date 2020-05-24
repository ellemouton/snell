package macaroon

import (
	"context"
	"encoding/hex"
	"testing"

	"github.com/lightningnetwork/lnd/lntypes"
	"github.com/stretchr/testify/require"
	"go.etcd.io/etcd/clientv3"
)

func TestGenerateRandomBytes(t *testing.T) {
	b := generateRandomBytes(2)
	require.Equal(t, len(b), 2)
}

type mockEtcd struct{}

func (mockEtcd) Put(ctx context.Context, key, val string, opts ...clientv3.OpOption) (*clientv3.PutResponse, error) {
	return nil, nil
}

func (mockEtcd) Close() error {
	return nil
}

func TestCreateMacaroon(t *testing.T) {
	c := &client{etcdClient: mockEtcd{}}

	var payHash lntypes.Hash
	payHashBytes, err := hex.DecodeString("b47f88a477a018b4dbeb85753d3ad8926a861fc247390e33ca17d7c80609807d")
	require.NoError(t, err)
	copy(payHash[:], payHashBytes)

	mac, err := c.Create(payHash, "article", 1)
	require.NoError(t, err)

	cavs := mac.Caveats()
	require.Equal(t, len(cavs), 1)

	snellCav, err := decodeCaveat(string(cavs[0].Id))
	require.NoError(t, err)
	require.Equal(t, "article", snellCav.key)
	require.Equal(t, "1", snellCav.value)
}
