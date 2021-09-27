package core

type ClientOptions struct {
	Dialer Dialer
}

func WithDialer(d Dialer) ClientOption {
	return func(o *ClientOptions) {
		o.Dialer = d
	}
}
