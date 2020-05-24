package macaroon

type Client interface {
	Close() error
}
