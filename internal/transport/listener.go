package transport

import "context"

type Listener interface {
	// Close the server. All active sessions will be closed.
	Close() error
	// Accept returns new sessions. It should be called in a loop.
	Accept(context.Context) (Session, error)
}
