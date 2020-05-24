package macaroon

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCaveat(t *testing.T) {
	c := &caveat{
		key:   "article",
		value: "1",
	}

	s := c.toString()
	require.Equal(t, "article=1", s)

	b := c.encode()

	c2, err := decodeCaveat(string(b))
	require.NoError(t, err)

	require.Equal(t, *c, *c2)
}
