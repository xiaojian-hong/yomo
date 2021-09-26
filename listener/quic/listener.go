package quic

import (
	"context"
	"time"

	"github.com/lucas-clemente/quic-go"
	"github.com/yomorun/yomo/internal/core"
	"github.com/yomorun/yomo/pkg/tls"
)

type QuicListener struct {
	c *quic.Config
	quic.Listener
}

func NewListener() *QuicListener {
	qconf := &quic.Config{
		Versions:                       []quic.VersionNumber{quic.Version1, quic.VersionDraft29},
		MaxIdleTimeout:                 time.Second * 3,
		KeepAlive:                      true,
		MaxIncomingStreams:             10000,
		MaxIncomingUniStreams:          10000,
		HandshakeIdleTimeout:           time.Second * 3,
		InitialStreamReceiveWindow:     1024 * 1024 * 2,
		InitialConnectionReceiveWindow: 1024 * 1024 * 2,
		DisablePathMTUDiscovery:        true,
	}

	return &QuicListener{
		c: qconf,
	}
}
func (l *QuicListener) Listen(ctx context.Context, addr string) error {
	// listen the address
	listener, err := quic.ListenAddr(addr, tls.GenerateTLSConfig(addr), l.c)
	if err != nil {
		return err
	}
	l.Listener = listener
	return nil
}

func (l *QuicListener) Name() string {
	return "quic"
}

func (l *QuicListener) Versions() []string {
	vers := make([]string, 0)
	for _, v := range l.c.Versions {
		vers = append(vers, v.String())
	}
	return vers
}

func (l *QuicListener) Close() error {
	return l.Listener.Close()
}

func (l *QuicListener) Accept(ctx context.Context) (core.Session, error) {
	session, err := l.Listener.Accept(ctx)
	if err != nil {
		return nil, err
	}

	return NewQuicSession(session), err
}
