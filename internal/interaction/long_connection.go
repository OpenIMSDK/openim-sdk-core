package interaction

import (
	"net/http"
	"time"
)

type PingPongHandler func(string) error

type LongConn interface {
	// Close closes this connection.
	Close() error

	// WriteMessage writes a message to the connection.
	// messageType indicates the type of data and can be set to binary (2) or text (1).
	WriteMessage(messageType int, message []byte) error

	// ReadMessage reads a message from the connection.
	ReadMessage() (int, []byte, error)

	// SetReadDeadline sets the deadline for reading from the underlying network connection.
	// After a timeout, there will be an error in the writing process.
	SetReadDeadline(timeout time.Duration) error

	// SetWriteDeadline sets the deadline for writing to the connection.
	// After a timeout, there will be an error in the writing process.
	SetWriteDeadline(timeout time.Duration) error

	// Dial tries to establish a connection.
	// urlStr must include authentication arguments; requestHeader can control data compression.
	Dial(urlStr string, requestHeader http.Header) (*http.Response, error)

	// IsNil checks whether the current long connection is nil.
	IsNil() bool

	// SetReadLimit sets the maximum size for a message read from the peer in bytes.
	SetReadLimit(limit int64)

	// SetPingHandler sets the handler for ping messages.
	SetPingHandler(handler PingPongHandler)

	// SetPongHandler sets the handler for pong messages.
	SetPongHandler(handler PingPongHandler)

	// LocalAddr returns the local network address.
	LocalAddr() string
}
