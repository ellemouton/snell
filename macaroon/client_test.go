package macaroon

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerateRandomBytes(t *testing.T) {
	b := generateRandomBytes(2)
	require.Equal(t, len(b), 2)
}
