// Package registry uses the go-micro registry for selection
package registry

import (
	"github.com/arun-spire/go-micro/client/selector"
)

// NewSelector returns a new registry selector
func NewSelector(opts ...selector.Option) selector.Selector {
	return selector.NewSelector(opts...)
}
