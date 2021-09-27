package kcp

import (
	"crypto/sha1"

	kcp "github.com/xtaci/kcp-go/v5"
	"github.com/yomorun/yomo/internal/core"
	"golang.org/x/crypto/pbkdf2"
)

type KcpDialer struct {
}

func NewDialer() *KcpDialer {
	return &KcpDialer{}
}

func (d *KcpDialer) Name() string {
	return "KCP-Client"
}

func (d *KcpDialer) Dial(addr string) (core.Session, error) {
	key := pbkdf2.Key([]byte(pass), []byte(salt), 4096, 32, sha1.New)
	block, _ := kcp.NewAESBlockCrypt(key)
	session, err := kcp.DialWithOptions(addr, block, dataShards, parityShards)
	if err != nil {
		return nil, err
	}

	return NewKcpSession(session), nil
}
