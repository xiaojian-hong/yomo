package kcp

import "github.com/xtaci/smux"

type KcpStream struct {
	*smux.Stream
}

func NewKcpStream(s *smux.Stream) *KcpStream {
	return &KcpStream{s}
}

func (s *KcpStream) StreamID() int64 {
	return int64(s.Stream.ID())
}
