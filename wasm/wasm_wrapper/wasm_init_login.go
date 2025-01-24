//go:build js && wasm

package wasm_wrapper

import (
	"errors"
	"syscall/js"

	"github.com/openimsdk/openim-sdk-core/v3/open_im_sdk"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/utils"
	"github.com/openimsdk/openim-sdk-core/v3/wasm/event_listener"
)

const COMMONEVENTFUNC = "commonEventFunc"

var ErrArgsLength = errors.New("from javascript args length err")
var ErrFunNameNotSet = errors.New("reflect func not to set")

type SetListener struct {
	*WrapperCommon
}

func NewSetListener(wrapperCommon *WrapperCommon) *SetListener {
	return &SetListener{WrapperCommon: wrapperCommon}
}

func (s *SetListener) setConversationListener() {
	callback := event_listener.NewConversationCallback(s.commonFunc)
	open_im_sdk.SetConversationListener(callback)
}
func (s *SetListener) setAdvancedMsgListener() {
	callback := event_listener.NewAdvancedMsgCallback(s.commonFunc)
	open_im_sdk.SetAdvancedMsgListener(callback)
}

func (s *SetListener) setBatchMessageListener() {
	callback := event_listener.NewBatchMessageCallback(s.commonFunc)
	open_im_sdk.SetBatchMsgListener(callback)
}

func (s *SetListener) setFriendListener() {
	callback := event_listener.NewFriendCallback(s.commonFunc)
	open_im_sdk.SetFriendListener(callback)
}

func (s *SetListener) setGroupListener() {
	callback := event_listener.NewGroupCallback(s.commonFunc)
	open_im_sdk.SetGroupListener(callback)
}

func (s *SetListener) setUserListener() {
	callback := event_listener.NewUserCallback(s.commonFunc)
	open_im_sdk.SetUserListener(callback)
}

func (s *SetListener) setSignalingListener() {
	//callback := event_listener.NewSignalingCallback(s.commonFunc)
	//open_im_sdk.SetSignalingListener(callback)
}
func (s *SetListener) setCustomBusinessListener() {
	callback := event_listener.NewCustomBusinessCallback(s.commonFunc)
	open_im_sdk.SetCustomBusinessListener(callback)
}

func (s *SetListener) SetAllListener() {
	s.setConversationListener()
	s.setAdvancedMsgListener()
	s.setBatchMessageListener()
	s.setFriendListener()
	s.setGroupListener()
	s.setUserListener()
	s.setSignalingListener()
	s.setCustomBusinessListener()
}

type WrapperCommon struct {
	commonFunc *js.Value
}

func NewWrapperCommon() *WrapperCommon {
	return &WrapperCommon{}
}
func (w *WrapperCommon) CommonEventFunc(_ js.Value, args []js.Value) interface{} {
	if len(args) >= 1 {
		w.commonFunc = &args[len(args)-1]
		return js.ValueOf(true)
	} else {
		return js.ValueOf(false)
	}
}

type WrapperInitLogin struct {
	*WrapperCommon
}

func NewWrapperInitLogin(wrapperCommon *WrapperCommon) *WrapperInitLogin {
	return &WrapperInitLogin{WrapperCommon: wrapperCommon}
}
func (w *WrapperInitLogin) InitSDK(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewConnCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return js.ValueOf(event_listener.NewCaller(open_im_sdk.InitSDK, callback, &args).SyncCall())
}
func (w *WrapperInitLogin) Login(_ js.Value, args []js.Value) interface{} {
	listener := NewSetListener(w.WrapperCommon)
	listener.SetAllListener()
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.Login, callback, &args).AsyncCallWithCallback()
}
func (w *WrapperInitLogin) Logout(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.Logout, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperInitLogin) NetworkStatusChanged(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.NetworkStatusChanged, callback, &args).AsyncCallWithCallback()
}
func (w *WrapperInitLogin) GetLoginStatus(_ js.Value, args []js.Value) interface{} {
	return event_listener.NewCaller(open_im_sdk.GetLoginStatus, nil, &args).AsyncCallWithOutCallback()
}
func (w *WrapperInitLogin) SetAppBackgroundStatus(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.SetAppBackgroundStatus, callback, &args).AsyncCallWithCallback()
}
