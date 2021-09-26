package core

import "io"

type Stream interface {
	io.ReadWriteCloser
	StreamID() int64
}
