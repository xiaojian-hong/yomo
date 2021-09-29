//go:build linux
// +build linux

package kcp

import (
	kcp "github.com/xtaci/kcp-go/v5"
)

func setSessionSocketBuf(session *kcp.UDPSession) error {
	if err := session.SetReadBuffer(sockBuf); err != nil {
		return err
	}
	if err := session.SetWriteBuffer(sockBuf); err != nil {
		return err
	}

	return nil
}

func setListenerSocketBuf(listener *kcp.Listener) error {
	if err := listener.SetReadBuffer(sockBuf); err != nil {
		return err
	}
	if err := listener.SetWriteBuffer(sockBuf); err != nil {
		return err
	}

	return nil
}
