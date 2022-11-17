package console

import (
	"bytes"
	"io"
)

type memory struct {
	err *bytes.Buffer
	out *bytes.Buffer
	in  *bytes.Buffer
}

func NewMemory() Console {
	return &memory{
		err: &bytes.Buffer{},
		out: &bytes.Buffer{},
		in:  &bytes.Buffer{},
	}
}

func (m *memory) Error() io.Writer {
	return m.err
}

func (m *memory) Out() io.Writer {
	return m.out
}

func (m *memory) In() io.Reader {
	return m.in
}
