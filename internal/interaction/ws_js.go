// Copyright © 2023 OpenIM SDK. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build js && wasm
// +build js,wasm

package interaction

import (
	"bytes"
	"context"
	"fmt"
	"github.com/OpenIMSDK/tools/log"
	"io"
	"net/http"
	"net/url"
	"nhooyr.io/websocket"
	"time"
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

func (w *JSWebSocket) SetPongHandler(handler PongHandler) {

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
	log.ZDebug(ctx, "ws msg read resp", "data", string(data))
	httpResp.Body = io.NopCloser(bytes.NewReader(data))
	return conn, httpResp, nil
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
