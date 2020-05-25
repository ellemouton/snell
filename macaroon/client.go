package macaroon

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"strings"

	"github.com/lightningnetwork/lnd/lntypes"
	"gopkg.in/macaroon.v2"

	"github.com/ellemouton/snell/etcd"
)

type client struct {
	etcdClient etcd.Client
}

func New() (Client, error) {
	ec, err := etcd.New()
	if err != nil {
		return nil, err
	}

	return &client{
		etcdClient: ec,
	}, nil
}

func (c *client) Close() error {
	return c.etcdClient.Close()
}

func (c *client) Create(paymentHash lntypes.Hash, resourceType string, resourceID int64) (*macaroon.Macaroon, error) {
	var rootKey [32]byte
	copy(rootKey[:], generateRandomBytes(32))

	var userID [32]byte
	copy(userID[:], generateRandomBytes(32))

	id := &identifier{
		version:     0,
		paymentHash: paymentHash,
		userID:      userID,
	}

	idBytes, err := id.encode()
	if err != nil {
		return nil, err
	}

	idHash := sha256.Sum256(idBytes)

	idKey := strings.Join([]string{"snell", "secrets", hex.EncodeToString(idHash[:])}, "/")

	// Store the key-value pair
	c.etcdClient.Put(context.TODO(), idKey, string(rootKey[:]))

	mac, err := macaroon.New(rootKey[:], idBytes, "snell", macaroon.LatestVersion)
	if err != nil {
		return nil, err
	}

	cav := &caveat{
		key:   resourceType,
		value: strconv.FormatInt(resourceID, 10),
	}

	if err := mac.AddFirstPartyCaveat(cav.encode()); err != nil {
		return nil, err
	}

	return mac, nil
}

func (c *client) Verify(macBytes []byte, preimage []byte, resourceType string, resourceID int64) (bool, error) {

	mac := &macaroon.Macaroon{}
	err := mac.UnmarshalBinary(macBytes)
	if err != nil {
		return false, err
	}

	id, err := decodeIdentifier(mac.Id())
	if err != nil {
		return false, err
	}

	// 1. check preimage and paymentHash
	if id.paymentHash != sha256.Sum256(preimage) {
		return false, nil
	}

	// 2. check that this macaroon was made by this server
	idHash := sha256.Sum256(mac.Id())
	idKey := strings.Join([]string{"snell", "secrets", hex.EncodeToString(idHash[:])}, "/")

	resp, err := c.etcdClient.Get(context.TODO(), idKey)
	if err != nil {
		return false, err
	}

	if len(resp.Kvs) == 0 {
		return false, nil
	}

	rawCaveats, err := mac.VerifySignature(resp.Kvs[0].Value, nil)
	if err != nil {
		return false, err
	}

	// 3. check that correct caveats for given resource are present
	for _, rawCaveat := range rawCaveats {
		caveat, err := decodeCaveat(rawCaveat)
		if err != nil {
			continue
		}

		if caveat.key != resourceType {
			continue
		}

		if caveat.value != strconv.FormatInt(resourceID, 10) {
			continue
		}

		return true, nil
	}

	return false, nil
}

func generateRandomBytes(n int64) []byte {
	b := make([]byte, n)
	rand.Read(b)
	return b
}
