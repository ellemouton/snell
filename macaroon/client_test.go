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

func (mockEtcd) Get(ctx context.Context, key string, opts ...clientv3.OpOption) (*clientv3.GetResponse, error) {
	return nil, nil
}

func (mockEtcd) Delete(ctx context.Context, key string, opts ...clientv3.OpOption) (*clientv3.DeleteResponse, error) {
	return nil, nil
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

// etcd needs to be running for the following test to work
func TestVerifyMacaroon(t *testing.T) {
	c, err := New()
	require.NoError(t, err)

	var payHash lntypes.Hash
	payHashBytes, err := hex.DecodeString("35fc86eda88ea58ebcbc24a9a9a487588a6114dccaf9bc9bc4f5ef086d71a96e")
	require.NoError(t, err)
	copy(payHash[:], payHashBytes)

	preimage1, err := hex.DecodeString("cba21a8096d21ae456f5fba0676407bfe3fa30525704ad48a5e7c6b24f1b8b86")
	require.NoError(t, err)

	preimage2, err := hex.DecodeString("bca21a8096d21ae456f5fba0676407bfe3fa30525704ad48a5e7c6b24f1b8b78")
	require.NoError(t, err)

	mac, err := c.Create(payHash, "article", 1)
	require.NoError(t, err)

	valid, err := c.Verify(mac, preimage1, "article", 1)
	require.NoError(t, err)
	require.True(t, valid)

	valid, err = c.Verify(mac, preimage2, "article", 1)
	require.NoError(t, err)
	require.False(t, valid)

}
