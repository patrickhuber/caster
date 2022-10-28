package fs

import "io"

// File defines an abstraction for FileSystem file operations
type File interface {
	io.Closer
	io.Reader
	io.ReaderAt
	io.Seeker
	io.Writer
	io.WriterAt
	io.StringWriter
}
