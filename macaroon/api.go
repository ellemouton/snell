package macaroon

import (
	"github.com/lightningnetwork/lnd/lntypes"
	"gopkg.in/macaroon.v2"
)

type Client interface {
	Close() error
	Create(paymentHash lntypes.Hash, resourceType string, resourceID int64) (*macaroon.Macaroon, error)
}
