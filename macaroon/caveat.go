package macaroon

import (
	"errors"
	"fmt"
	"strings"
)

type caveat struct {
	key   string
	value string
}

func (c *caveat) toString() string {
	return fmt.Sprintf("%s=%s", c.key, c.value)
}

func (c *caveat) encode() []byte {
	s := c.toString()
	return []byte(s)
}

func decodeCaveat(s string) (*caveat, error) {
	parts := strings.SplitN(s, "=", 2)

	if len(parts) != 2 {
		return nil, errors.New(fmt.Sprintf("invalid caveat string: %s", s))
	}

	return &caveat{
		key:   parts[0],
		value: parts[1],
	}, nil
}
