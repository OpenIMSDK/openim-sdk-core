package interaction

import (
	"time"

	"github.com/openimsdk/protocol/constant"
)

type ConnContext struct {
	RemoteAddr string
}

func (c *ConnContext) Deadline() (deadline time.Time, ok bool) {
	return
}

func (c *ConnContext) Done() <-chan struct{} {
	return nil
}

func (c *ConnContext) Err() error {
	return nil
}

func (c *ConnContext) Value(key any) any {
	switch key {
	case constant.RemoteAddr:
		return c.RemoteAddr
	default:
		return ""
	}
}

func newContext(remoteAddr string) *ConnContext {
	return &ConnContext{
		RemoteAddr: remoteAddr,
	}
}
