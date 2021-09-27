package quic

import (
	"context"

	"github.com/lucas-clemente/quic-go"
	"github.com/yomorun/yomo/internal/core"
)

type QuicSession struct {
	quic.Session
}

func NewQuicSession(s quic.Session) *QuicSession {
	return &QuicSession{s}
}

func (s *QuicSession) AcceptStream(ctx context.Context) (core.Stream, error) {
	stream, err := s.Session.AcceptStream(ctx)
	return NewQuicStream(stream), err
}

func (s *QuicSession) OpenStream(ctx context.Context) (core.Stream, error) {
	stream, err := s.Session.OpenStreamSync(ctx)
	return NewQuicStream(stream), err
}

func (s *QuicSession) Close() error {
	return nil
}

func (s *QuicSession) CloseWithError(ecode uint64, msg string) error {
	return s.Session.CloseWithError(quic.ApplicationErrorCode(ecode), msg)
}
