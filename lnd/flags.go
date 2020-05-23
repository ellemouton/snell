package lnd

import "flag"

var (
	lnd_address = flag.String("lnd_address", "127.0.0.1:10009", "host:port of lnd gRPC service")
	lnd_cert    = flag.String("lnd_cert", "/Users/ellemouton/.lnd/tls.cert", "lnd cert location")
	mac_path    = flag.String("mac_path", "/Users/ellemouton/.lnd/data/chain/bitcoin/testnet/admin.macaroon", "path to macaroon file")
)
