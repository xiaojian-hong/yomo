package kcp

import (
	"context"
	"crypto/sha1"
	"net"

	kcp "github.com/xtaci/kcp-go/v5"
	"github.com/yomorun/yomo/internal/core"
	"golang.org/x/crypto/pbkdf2"
)

const (
	// listener config
	// salt is use for pbkdf2 key expansion
	salt = "yomo-kcp-go"
	// maximum supported smux version
	maxSmuxVer = 2
	// stream copy buffer size
	bufSize = 4096
	// 'dataShards', 'parityShards' specify how many parity packets will be generated following the data packets.
	dataShards   = 10
	parityShards = 3
	pass         = "yomo"

	// session config
	sessionNoDelay      = 1
	sessionInterval     = 40
	sessionResend       = 2
	sessionNoCongestion = 1
	sessionMTU          = 1400
	sessionSndWnd       = 2048
	sessionRcvWnd       = 2048
	sessionAckNodelay   = false
)

type KcpListener struct {
	*kcp.Listener
}

func NewListener() *KcpListener {
	return &KcpListener{}
}

func (l *KcpListener) Name() string {
	return "KCP-Server"
}

func (l *KcpListener) Listen(ctx context.Context, addr string) error {
	key := pbkdf2.Key([]byte(pass), []byte(salt), 4096, 32, sha1.New)
	block, _ := kcp.NewAESBlockCrypt(key)
	listener, err := kcp.ListenWithOptions(addr, block, dataShards, parityShards)
	if err != nil {
		return err
	}
	l.Listener = listener

	return nil
}

// Close the server. All active sessions will be closed.
func (l *KcpListener) Close() error {
	return l.Listener.Close()
}

// Addr returns the local network addr that the server is listening on.
func (l *KcpListener) Addr() net.Addr {
	return l.Listener.Addr()
}

// Accept returns new sessions. It should be called in a loop.
func (l *KcpListener) Accept(_ context.Context) (core.Session, error) {
	session, err := l.Listener.AcceptKCP()
	if err != nil {
		return nil, err
	}
	session.SetStreamMode(true)
	session.SetWriteDelay(false)
	session.SetNoDelay(sessionNoDelay, sessionInterval, sessionResend, sessionNoCongestion)
	session.SetMtu(sessionMTU)
	session.SetWindowSize(sessionSndWnd, sessionRcvWnd)
	session.SetACKNoDelay(sessionAckNodelay)

	return NewKcpSession(session), err
}

// Versions
func (l *KcpListener) Versions() []string {
	return []string{"v5"}
}
