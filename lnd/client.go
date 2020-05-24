package lnd

import (
	"context"
	"io/ioutil"

	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/lightningnetwork/lnd/macaroons"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"gopkg.in/macaroon.v2"
)

type client struct {
	rpcConn   *grpc.ClientConn
	rpcClient lnrpc.LightningClient
}

func New() (Client, error) {
	creds, err := credentials.NewClientTLSFromFile(*lnd_cert, "")
	if err != nil {
		return nil, err
	}

	macBytes, err := ioutil.ReadFile(*mac_path)
	if err != nil {
		return nil, err
	}

	mac := &macaroon.Macaroon{}
	if err = mac.UnmarshalBinary(macBytes); err != nil {
		return nil, errors.Wrap(err, "unable to decode macaroon")
	}

	macCreds := macaroons.NewMacaroonCredential(mac)

	conn, err := grpc.Dial(*lnd_address, grpc.WithTransportCredentials(creds), grpc.WithPerRPCCredentials(macCreds))
	if err != nil {
		return nil, err
	}

	return &client{
		rpcConn:   conn,
		rpcClient: lnrpc.NewLightningClient(conn),
	}, nil
}

func (c *client) AddInvoice(ctx context.Context, numSats int64, expiry int64, memo string) (*lnrpc.AddInvoiceResponse, error) {
	req := &lnrpc.Invoice{
		Value:  numSats,
		Expiry: expiry,
		Memo:   memo,
	}
	return c.rpcClient.AddInvoice(ctx, req)
}

func (c *client) Close() error {
	return c.rpcConn.Close()
}
