package kcp

import (
	"context"
	"net"
	"time"

	kcp "github.com/xtaci/kcp-go/v5"
	"github.com/yomorun/yomo/internal/core"

	"github.com/xtaci/smux"
)

const (
	smuxVer       = 2
	smuxKeepAlive = 10
	smuxFrameSize = 65535
)

type KcpSession struct {
	Session *smux.Session
}

func NewKcpSession(s *kcp.UDPSession) *KcpSession {
	// stream multiplex
	smuxConfig := smux.DefaultConfig()
	smuxConfig.Version = smuxVer
	smuxConfig.MaxReceiveBuffer = sockBuf
	smuxConfig.MaxStreamBuffer = streamBuf
	smuxConfig.KeepAliveInterval = time.Duration(smuxKeepAlive) * time.Second
	smuxConfig.MaxFrameSize = smuxFrameSize

	sesssion, err := smux.Server(s, smuxConfig)
	if err != nil {
		panic(err)
	}

	return &KcpSession{sesssion}
}

func (s *KcpSession) AcceptStream(ctx context.Context) (core.Stream, error) {
	stream, err := s.Session.AcceptStream()
	if err != nil {
		return nil, err
	}

	return NewKcpStream(stream), nil
}

func (s *KcpSession) OpenStream(ctx context.Context) (core.Stream, error) {
	stream, err := s.Session.OpenStream()
	return NewKcpStream(stream), err
}

func (s *KcpSession) Close() error {
	if s.Session != nil {
		return s.Session.Close()
	}
	return nil
}

func (s *KcpSession) CloseWithError(ecode uint64, msg string) error {
	// TODO:
	return s.Close()
}

func (s *KcpSession) LocalAddr() net.Addr {
	return s.Session.LocalAddr()
}

func (s *KcpSession) RemoteAddr() net.Addr {
	return s.Session.RemoteAddr()
}
