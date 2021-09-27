package quic

import "github.com/lucas-clemente/quic-go"

type QuicStream struct {
	quic.Stream
}

func NewQuicStream(s quic.Stream) *QuicStream {
	return &QuicStream{s}
}

func (s *QuicStream) StreamID() int64 {
	return int64(s.Stream.StreamID())
}
