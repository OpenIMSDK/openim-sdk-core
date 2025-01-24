//go:build js && wasm

package event_listener

import (
	"github.com/openimsdk/openim-sdk-core/v3/pkg/utils"
	"syscall/js"
)

type CallbackWriter interface {
	SendMessage()
	SetEvent(event string) CallbackWriter
	SetData(data interface{}) CallbackWriter
	SetErrCode(errCode int32) CallbackWriter
	SetOperationID(operationID string) CallbackWriter
	SetErrMsg(errMsg string) CallbackWriter
	GetOperationID() string
	HandlerFunc(fn FuncLogic) interface{}
}

type EventData struct {
	Event       string      `json:"event"`
	ErrCode     int32       `json:"errCode"`
	ErrMsg      string      `json:"errMsg"`
	Data        interface{} `json:"data,omitempty"`
	OperationID string      `json:"operationID"`
	callback    *js.Value
}

func (e *EventData) HandlerFunc(fn FuncLogic) interface{} {
	panic("implement me")
}

func (e *EventData) GetOperationID() string {
	return e.OperationID
}

func NewEventData(callback *js.Value) *EventData {
	return &EventData{callback: callback}
}
func (e *EventData) SendMessage() {
	e.callback.Invoke(utils.StructToJsonString(e))
}
func (e *EventData) SetEvent(event string) CallbackWriter {
	e.Event = event
	return e
}

func (e *EventData) SetData(data interface{}) CallbackWriter {
	e.Data = data
	return e
}
func (e *EventData) SetErrCode(errCode int32) CallbackWriter {
	e.ErrCode = errCode
	return e
}
func (e *EventData) SetOperationID(operationID string) CallbackWriter {
	e.OperationID = operationID
	return e
}
func (e *EventData) SetErrMsg(errMsg string) CallbackWriter {
	e.ErrMsg = errMsg
	return e
}

var (
	jsErr     = js.Global().Get("Error")
	jsPromise = js.Global().Get("Promise")
)

type PromiseHandler struct {
	Event       string      `json:"event"`
	ErrCode     int32       `json:"errCode"`
	ErrMsg      string      `json:"errMsg"`
	Data        interface{} `json:"data,omitempty"`
	OperationID string      `json:"operationID"`
	resolve     *js.Value
	reject      *js.Value
}

func NewPromiseHandler() *PromiseHandler {
	return &PromiseHandler{}
}
func (p *PromiseHandler) HandlerFunc(fn FuncLogic) interface{} {
	handler := js.FuncOf(func(_ js.Value, promFn []js.Value) interface{} {
		p.resolve, p.reject = &promFn[0], &promFn[1]
		fn()
		return nil
	})
	return jsPromise.New(handler)
}

func (p *PromiseHandler) GetOperationID() string {
	return p.OperationID
}

func (p *PromiseHandler) SendMessage() {
	if p.Data != nil {
		p.resolve.Invoke(p.Data)
	} else {
		//p.reject.Invoke(jsErr.New(fmt.Sprintf("erCode:%d,errMsg:%s,operationID:%s", p.ErrCode, p.ErrMsg, p.OperationID)))
		errInfo := make(map[string]interface{})
		errInfo["errCode"] = p.ErrCode
		errInfo["errMsg"] = p.ErrMsg
		errInfo["operationID"] = p.OperationID
		p.reject.Invoke(errInfo)
	}
}
func (p *PromiseHandler) SetEvent(event string) CallbackWriter {
	p.Event = event
	return p
}

func (p *PromiseHandler) SetData(data interface{}) CallbackWriter {
	p.Data = data
	return p
}
func (p *PromiseHandler) SetErrCode(errCode int32) CallbackWriter {
	p.ErrCode = errCode
	return p
}
func (p *PromiseHandler) SetOperationID(operationID string) CallbackWriter {
	p.OperationID = operationID
	return p
}
func (p *PromiseHandler) SetErrMsg(errMsg string) CallbackWriter {
	p.ErrMsg = errMsg
	return p
}
