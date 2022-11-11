package setup

import (
	"io"

	"github.com/patrickhuber/go-di"
)

type Setup interface {
	io.Closer
	Container() di.Container
}
