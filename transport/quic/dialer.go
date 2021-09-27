package quic

import (
	"crypto/tls"
	"time"

	"github.com/lucas-clemente/quic-go"
	"github.com/yomorun/yomo/internal/core"
)

type QuicDialer struct {
	c  *quic.Config
	tc *tls.Config
}

func NewDialer() *QuicDialer {
	// tls config
	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"yomo"},
		ClientSessionCache: tls.NewLRUClientSessionCache(64),
	}
	// quic config
	quicConf := &quic.Config{
		Versions:                       []quic.VersionNumber{quic.Version1, quic.VersionDraft29},
		MaxIdleTimeout:                 time.Second * 3,
		KeepAlive:                      true,
		MaxIncomingStreams:             10000,
		MaxIncomingUniStreams:          10000,
		HandshakeIdleTimeout:           time.Second * 3,
		InitialStreamReceiveWindow:     1024 * 1024 * 2,
		InitialConnectionReceiveWindow: 1024 * 1024 * 2,
		TokenStore:                     quic.NewLRUTokenStore(1, 1),
		DisablePathMTUDiscovery:        true,
	}

	return &QuicDialer{
		tc: tlsConf,
		c:  quicConf,
	}
}

func (d *QuicDialer) Name() string {
	return "QUIC-Client"
}

func (d *QuicDialer) Dial(addr string) (core.Session, error) {
	session, err := quic.DialAddr(addr, d.tc, d.c)
	if err != nil {
		return nil, err
	}
	return NewQuicSession(session), nil
}
