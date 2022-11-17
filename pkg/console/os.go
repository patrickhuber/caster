package console

import (
	"io"
	"os"
)

type Console interface {
	Error() io.Writer
	Out() io.Writer
	In() io.Reader
}

func New() Console {
	return &console{}
}

type console struct {
}

func (c *console) Error() io.Writer {
	return os.Stderr
}

func (c *console) Out() io.Writer {
	return os.Stdout
}

func (c *console) In() io.Reader {
	return os.Stdin
}
