package bridge

import (
	"net/http"
	"net/url"
	"os"
	"sync"

	"github.com/yomorun/yomo/core"
	"github.com/yomorun/yomo/core/frame"
	"github.com/yomorun/yomo/pkg/logger"
	"golang.org/x/net/websocket"
)

// WebSocketBridge implements the Bridge interface for WebSocket.
type WebSocketBridge struct {
	addr   string
	server *websocket.Server

	// Registered the connections in each room.
	// Key: connID
	// Value: *websocket.Conn
	conns sync.Map
}

// meshID is the access point for users. This services can be deployed at anywhere
// to decrease connection_time. Users can connect to the nearest access point.
var meshID = os.Getenv("MESH_ID")

// NewWebSocketBridge initializes an instance for WebSocketBridge.
func NewWebSocketBridge(addr string) *WebSocketBridge {
	// server config
	conf := websocket.Config{
		Origin: &url.URL{
			Host: addr,
		},
		Header: make(http.Header),
	}

	// add mesh ID to header.
	if meshID != "" {
		conf.Header.Add("YoMo-Mesh-ID", meshID)
	}

	return &WebSocketBridge{
		addr: addr,
		server: &websocket.Server{
			Config: conf,
			Handshake: func(c *websocket.Config, r *http.Request) error {
				// TODO: check Origin header for auth.
				return nil
			},
		},
		conns: sync.Map{},
	}
}

// Name returns the name of WebSocket bridge.
func (ws *WebSocketBridge) Name() string {
	return nameOfWebSocket
}

// Addr returns the address of bridge.
func (ws *WebSocketBridge) Addr() string {
	return ws.addr
}

// ListenAndServe starts a WebSocket server with a given handler.
func (ws *WebSocketBridge) ListenAndServe(handler func(ctx *core.Context)) error {
	// wrap the WebSocket handler.
	ws.server.Handler = func(c *websocket.Conn) {
		// set payload type
		c.PayloadType = websocket.BinaryFrame
		// TODO: support multi rooms.
		connID := c.Request().RemoteAddr
		ws.conns.Store(connID, c)

		// trigger the YoMo Server's Handler in bridge.
		handler(&core.Context{
			ConnID: connID,
			Stream: c,
			OnClose: func(code uint64, msg string) {
				// remove this connection in room.
				ws.conns.Delete(connID)
			},
		})
	}

	// serve
	return http.ListenAndServe(ws.addr, ws.server)
}

// Send the data to WebSocket clients.
func (ws *WebSocketBridge) Send(f frame.Frame) error {
	ws.conns.Range(func(key, value interface{}) bool {
		if c, ok := value.(*websocket.Conn); ok {
			_, err := c.Write(f.Encode())
			if err != nil {
				logger.Errorf("[WebSocketBridge] send data to conn failed, connID=%s", key)
			}
		}
		return true
	})
	return nil
}
