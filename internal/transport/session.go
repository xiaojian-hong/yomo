package transport

import "io"

type Session interface {
	io.ReadWriter
}
