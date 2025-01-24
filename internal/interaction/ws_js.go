//go:build js && wasm

package interaction

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"nhooyr.io/websocket"

	"github.com/openimsdk/tools/log"
)

type JSWebSocket struct {
	ConnType int
	conn     *websocket.Conn
	sendConn *websocket.Conn
}

func (w *JSWebSocket) SetReadDeadline(timeout time.Duration) error {
	return nil
}

func (w *JSWebSocket) SetWriteDeadline(timeout time.Duration) error {
	return nil
}

func (w *JSWebSocket) SetReadLimit(limit int64) {
	w.conn.SetReadLimit(limit)
}

func (w *JSWebSocket) SetPingHandler(handler PingPongHandler) {

}

func (w *JSWebSocket) SetPongHandler(handler PingPongHandler) {

}

func (w *JSWebSocket) LocalAddr() string {
	return ""
}

func NewWebSocket(connType int) *JSWebSocket {
	return &JSWebSocket{ConnType: connType}
}

func (w *JSWebSocket) Close() error {
	return w.conn.Close(websocket.StatusGoingAway, "Actively close the conn have old conn")
}

func (w *JSWebSocket) WriteMessage(messageType int, message []byte) error {
	if messageType == PingMessage || messageType == PongMessage {
		return nil
	}
	return w.conn.Write(context.Background(), websocket.MessageType(messageType), message)
}

func (w *JSWebSocket) ReadMessage() (int, []byte, error) {
	messageType, b, err := w.conn.Read(context.Background())
	return int(messageType), b, err
}

func (w *JSWebSocket) dial(ctx context.Context, urlStr string) (*websocket.Conn, *http.Response, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, nil, err
	}
	query := u.Query()
	query.Set("isMsgResp", "true")
	u.RawQuery = query.Encode()
	conn, httpResp, err := websocket.Dial(ctx, u.String(), nil)
	if err != nil {
		return nil, nil, err
	}
	if httpResp == nil {
		httpResp = &http.Response{
			StatusCode: http.StatusSwitchingProtocols,
		}
	}
	_, data, err := conn.Read(ctx)
	if err != nil {
		_ = conn.CloseNow()
		return nil, nil, fmt.Errorf("read response error %w", err)
	}
	var apiResp struct {
		ErrCode int    `json:"errCode"`
		ErrMsg  string `json:"errMsg"`
		ErrDlt  string `json:"errDlt"`
	}
	if err := json.Unmarshal(data, &apiResp); err != nil {
		return nil, nil, fmt.Errorf("unmarshal response error %w", err)
	}
	if apiResp.ErrCode == 0 {
		return conn, httpResp, nil
	}
	log.ZDebug(ctx, "ws msg read resp", "data", string(data))
	httpResp.Body = io.NopCloser(bytes.NewReader(data))
	return conn, httpResp, fmt.Errorf("read response error %d %s %s",
		apiResp.ErrCode, apiResp.ErrMsg, apiResp.ErrDlt)
}

func (w *JSWebSocket) Dial(urlStr string, _ http.Header) (*http.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	conn, httpResp, err := w.dial(ctx, urlStr)
	if err == nil {
		w.conn = conn
	}
	return httpResp, err
}

func (w *JSWebSocket) IsNil() bool {
	if w.conn != nil {
		return false
	}
	return true
}
