package lnd

import (
	"context"

	"github.com/lightningnetwork/lnd/lnrpc"
)

type Client interface {
	AddInvoice(ctx context.Context, numSats int64, expiry int64, memo string) (*lnrpc.AddInvoiceResponse, error)
	Close() error
}
