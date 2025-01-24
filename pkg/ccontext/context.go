package ccontext

import (
	"context"

	"github.com/openimsdk/openim-sdk-core/v3/open_im_sdk_callback"
	"github.com/openimsdk/openim-sdk-core/v3/sdk_struct"

	"github.com/openimsdk/tools/mcontext"
)

const (
	Callback = "callback"
)

type GlobalConfig struct {
	UserID string
	Token  string

	sdk_struct.IMConfig
}

type ContextInfo interface {
	UserID() string
	Token() string
	PlatformID() int32
	ApiAddr() string
	WsAddr() string
	DataDir() string
	LogLevel() uint32
	OperationID() string
	IsExternalExtensions() bool
}

func Info(ctx context.Context) ContextInfo {
	conf := ctx.Value(GlobalConfigKey{}).(*GlobalConfig)
	return &info{
		conf: conf,
		ctx:  ctx,
	}
}

func WithInfo(ctx context.Context, conf *GlobalConfig) context.Context {
	return context.WithValue(ctx, GlobalConfigKey{}, conf)
}

func WithOperationID(ctx context.Context, operationID string) context.Context {
	return mcontext.SetOperationID(ctx, operationID)
}
func WithSendMessageCallback(ctx context.Context, callback open_im_sdk_callback.SendMsgCallBack) context.Context {
	return context.WithValue(ctx, Callback, callback)
}

func WithApiErrCode(ctx context.Context, cb ApiErrCodeCallback) context.Context {
	return context.WithValue(ctx, apiErrCode{}, cb)
}

func GetApiErrCodeCallback(ctx context.Context) ApiErrCodeCallback {
	fn, _ := ctx.Value(apiErrCode{}).(ApiErrCodeCallback)
	if fn == nil {
		return &emptyApiErrCodeCallback{}
	}
	return fn
}

type GlobalConfigKey struct{}

type info struct {
	conf *GlobalConfig
	ctx  context.Context
}

func (i *info) UserID() string {
	return i.conf.UserID
}

func (i *info) Token() string {
	return i.conf.Token
}

func (i *info) PlatformID() int32 {
	return i.conf.PlatformID
}

func (i *info) ApiAddr() string {
	return i.conf.ApiAddr
}

func (i *info) WsAddr() string {
	return i.conf.WsAddr
}

func (i *info) DataDir() string {
	return i.conf.DataDir
}

func (i *info) LogLevel() uint32 {
	return i.conf.LogLevel
}

func (i *info) OperationID() string {
	return mcontext.GetOperationID(i.ctx)
}

func (i *info) IsExternalExtensions() bool {
	return i.conf.IsExternalExtensions
}

type apiErrCode struct{}

type ApiErrCodeCallback interface {
	OnError(ctx context.Context, err error)
}

type emptyApiErrCodeCallback struct{}

func (e *emptyApiErrCodeCallback) OnError(ctx context.Context, err error) {}
