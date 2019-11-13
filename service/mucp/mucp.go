// Package mucp initialises a mucp service
package mucp

import (
	// TODO: change to go-micro/service
	"github.com/arun-spire/go-micro"
	cmucp "github.com/arun-spire/go-micro/client/mucp"
	smucp "github.com/arun-spire/go-micro/server/mucp"
)

// NewService returns a new mucp service
func NewService(opts ...micro.Option) micro.Service {
	options := []micro.Option{
		micro.Client(cmucp.NewClient()),
		micro.Server(smucp.NewServer()),
	}

	options = append(options, opts...)

	return micro.NewService(options...)
}
