//go:build !linux
// +build !linux

package kcp

import (
	kcp "github.com/xtaci/kcp-go/v5"
)

func setSessionSocketBuf(_ *kcp.UDPSession) error {
	return nil
}

func setListenerSocketBuf(_ *kcp.Listener) error {
	return nil
}
