package eventStore

import (
	"context"
	"crypto/tls"
)

type Options struct {
	Address   string
	Secure    bool
	TLSConfig *tls.Config
	Context   context.Context
}
