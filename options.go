package yomo

import (
	"os"

	"github.com/yomorun/yomo/internal/core"
	"github.com/yomorun/yomo/transport/kcp"
	"github.com/yomorun/yomo/transport/quic"
)

const (
	// DefaultZipperAddr is the default address of downstream zipper.
	DefaultZipperAddr = "localhost:9000"
	// DefaultZipperListenAddr set default listening port to 9000 and binding to all interfaces.
	DefaultZipperListenAddr = "0.0.0.0:9000"
)

var (
	DefaultListener = quic.NewListener()
	DefaultDialer   = quic.NewDialer()
)

// Option is a function that applies a YoMo-Client option.
type Option func(o *options)

// options are the options for YoMo-Client.
type options struct {
	ZipperAddr string // target Zipper endpoint address
	// ZipperListenAddr     string // Zipper endpoint address
	ZipperWorkflowConfig string // Zipper workflow file
	MeshConfigURL        string // meshConfigURL is the URL of edge-mesh config
	// Listener
	Listener core.Listener
	// Dialer
	Dialer core.Dialer
}

// WithZipperAddr return a new options with ZipperAddr set to addr.
func WithZipperAddr(addr string) Option {
	return func(o *options) {
		o.ZipperAddr = addr
	}
}

// WithMeshConfigURL sets the initial edge-mesh config URL for the YoMo-Zipper.
func WithMeshConfigURL(url string) Option {
	return func(o *options) {
		o.MeshConfigURL = url
	}
}

func WithListener(listener core.Listener) Option {
	return func(o *options) {
		o.Listener = listener
	}
}

func WithDialer(dialer core.Dialer) Option {
	return func(o *options) {
		o.Dialer = dialer
	}
}

func WithEnv() Option {
	return func(o *options) {
		transport := os.Getenv("YOMO_TRANSPORT")
		switch transport {
		case "kcp":
			o.Dialer = kcp.NewDialer()
			o.Listener = kcp.NewListener()
		default:
			o.Dialer = quic.NewDialer()
			o.Listener = quic.NewListener()
		}
	}
}

// newOptions creates a new options for YoMo-Client.
func newOptions(opts ...Option) *options {
	options := &options{}

	for _, o := range opts {
		o(options)
	}

	if options.ZipperAddr == "" {
		options.ZipperAddr = DefaultZipperAddr
	}

	if options.Listener == nil {
		options.Listener = DefaultListener
	}

	if options.Dialer == nil {
		options.Dialer = DefaultDialer
	}

	return options
}
