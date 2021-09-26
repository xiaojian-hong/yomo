package core

type ServerOptions struct {
	Listener Listener
}

func WithListener(l Listener) ServerOption {
	return func(o *ServerOptions) {
		o.Listener = l
	}
}
