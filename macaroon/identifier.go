package macaroon

import (
	"bytes"
	"encoding/binary"

	"github.com/lightningnetwork/lnd/lntypes"
)

type identifier struct {
	version     uint16
	paymentHash lntypes.Hash
	userID      [32]byte
}

func (i *identifier) encode() ([]byte, error) {
	var buf bytes.Buffer

	if err := binary.Write(&buf, binary.BigEndian, i.version); err != nil {
		return nil, err
	}

	if _, err := buf.Write(i.paymentHash[:]); err != nil {
		return nil, err
	}

	if _, err := buf.Write(i.userID[:]); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func decodeIdentifier(b []byte) (*identifier, error) {

	id := new(identifier)
	r := bytes.NewReader(b)

	if err := binary.Read(r, binary.BigEndian, &id.version); err != nil {
		return nil, err
	}

	if _, err := r.Read(id.paymentHash[:]); err != nil {
		return nil, err
	}

	if _, err := r.Read(id.userID[:]); err != nil {
		return nil, err
	}

	return id, nil
}
