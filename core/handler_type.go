package core

import (
	"io"

	"github.com/yomorun/yomo/core/frame"
)

// AsyncHandler is the request-response mode (asnyc)
type AsyncHandler func([]byte) (byte, []byte)

// PipeHandler is the bidirectional stream mode (blocking).
type PipeHandler func(in <-chan []byte, out chan<- *frame.PayloadFrame)

// StreamHandler is the multiple stream mode.
type StreamHandler func(in io.Reader) io.Reader
