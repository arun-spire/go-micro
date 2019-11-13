package grpc

import (
	"crypto/tls"

	"github.com/arun-spire/go-micro"
	gc "github.com/arun-spire/go-micro/client/grpc"
	gs "github.com/arun-spire/go-micro/server/grpc"
)

// WithTLS sets the TLS config for the service
func WithTLS(t *tls.Config) micro.Option {
	return func(o *micro.Options) {
		o.Client.Init(
			gc.AuthTLS(t),
		)
		o.Server.Init(
			gs.AuthTLS(t),
		)
	}
}
