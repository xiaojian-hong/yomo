package core

import (
	"context"
	"net"
)

type Listener interface {
	Name() string
	Listen(ctx context.Context, addr string) error
	// Close the server. All active sessions will be closed.
	Close() error
	// Addr returns the local network addr that the server is listening on.
	Addr() net.Addr
	// Accept returns new sessions. It should be called in a loop.
	Accept(context.Context) (Session, error)
	// Versions
	Versions() []string
}
