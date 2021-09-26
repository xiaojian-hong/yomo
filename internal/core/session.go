package core

import (
	"context"
	"net"
)

type Session interface {
	// LocalAddr returns the local address.
	LocalAddr() net.Addr
	// RemoteAddr returns the address of the peer.
	RemoteAddr() net.Addr
	AcceptStream(ctx context.Context) (Stream, error)
	Close() error
	CloseWithError(uint64, string) error
}
