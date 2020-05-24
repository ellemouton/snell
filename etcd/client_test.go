package etcd

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

// Note that etcd needs to be running for this test to work
func TestEtcd(t *testing.T) {
	ctx := context.Background()

	ec, err := New()
	require.NoError(t, err)

	_, err = ec.Put(ctx, "snell/test_key", "snell/test_value")
	require.NoError(t, err)

	_, err = ec.Get(ctx, "snell/test_key")
	require.NoError(t, err)

	_, err = ec.Delete(ctx, "snell/test_key")
	require.NoError(t, err)

	err = ec.Close()
	require.NoError(t, err)
}
