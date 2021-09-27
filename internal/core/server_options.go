package core

type ServerOptions struct {
	Listener Listener
	Addr     string
}

func WithListener(l Listener) ServerOption {
	return func(o *ServerOptions) {
		o.Listener = l
	}
}

func WithAddr(addr string) ServerOption {
	return func(o *ServerOptions) {
		o.Addr = addr
	}
}
