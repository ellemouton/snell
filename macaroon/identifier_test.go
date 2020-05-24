package macaroon

import (
	"encoding/hex"
	"testing"

	"github.com/lightningnetwork/lnd/lntypes"
	"github.com/stretchr/testify/require"
)

func TestIdentifier(t *testing.T) {
	var userID [32]byte
	copy(userID[:], generateRandomBytes(31))

	var payHash lntypes.Hash

	payHashBytes, err := hex.DecodeString("b47f88a477a018b4dbeb85753d3ad8926a861fc247390e33ca17d7c80609807d")
	require.NoError(t, err)

	copy(payHash[:], payHashBytes)

	id := &identifier{
		version:     0,
		paymentHash: payHash,
		userID:      userID,
	}

	b, err := id.encode()
	require.NoError(t, err)

	id2, err := decodeIdentifier(b)
	require.NoError(t, err)

	require.Equal(t, *id, *id2)
}
