//go:build !js

package interaction

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type Default struct {
	ConnType  int
	conn      *websocket.Conn
	isSetConf bool
}

func (d *Default) SetReadDeadline(timeout time.Duration) error {
	return d.conn.SetReadDeadline(time.Now().Add(timeout))
}

func (d *Default) SetWriteDeadline(timeout time.Duration) error {
	return d.conn.SetWriteDeadline(time.Now().Add(timeout))
}

func (d *Default) SetReadLimit(limit int64) {
	d.conn.SetReadLimit(limit)

}

func (d *Default) SetPingHandler(handler PingPongHandler) {
	d.conn.SetPingHandler(handler)
}

func (d *Default) SetPongHandler(handler PingPongHandler) {
	d.conn.SetPongHandler(handler)
}

func (d *Default) LocalAddr() string {
	return d.conn.LocalAddr().String()
}

func NewWebSocket(connType int) *Default {
	return &Default{ConnType: connType}
}
func (d *Default) Close() error {
	return d.conn.Close()
}

func (d *Default) WriteMessage(messageType int, message []byte) error {
	return d.conn.WriteMessage(messageType, message)
}

func (d *Default) ReadMessage() (int, []byte, error) {
	return d.conn.ReadMessage()
}

func (d *Default) Dial(urlStr string, requestHeader http.Header) (*http.Response, error) {
	conn, httpResp, err := websocket.DefaultDialer.Dial(urlStr, requestHeader)
	if err == nil {
		d.conn = conn
	}
	return httpResp, err
}

func (d *Default) IsNil() bool {
	if d.conn != nil {
		return false
	}
	return true
}
