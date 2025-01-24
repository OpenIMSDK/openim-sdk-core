//go:build js && wasm

package wasm_wrapper

import (
	"syscall/js"

	"github.com/openimsdk/openim-sdk-core/v3/open_im_sdk"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/utils"
	"github.com/openimsdk/openim-sdk-core/v3/wasm/event_listener"
)

// ------------------------------------group---------------------------
type WrapperUser struct {
	*WrapperCommon
}

func NewWrapperUser(wrapperCommon *WrapperCommon) *WrapperUser {
	return &WrapperUser{WrapperCommon: wrapperCommon}
}

func (w *WrapperUser) GetSelfUserInfo(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.GetSelfUserInfo, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperUser) SetSelfInfo(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.SetSelfInfo, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperUser) GetUsersInfo(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.GetUsersInfo, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperUser) SubscribeUsersStatus(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.SubscribeUsersStatus, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperUser) UnsubscribeUsersStatus(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.UnsubscribeUsersStatus, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperUser) GetSubscribeUsersStatus(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.GetSubscribeUsersStatus, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperUser) GetUserStatus(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.GetUserStatus, callback, &args).AsyncCallWithCallback()
}
func (w *WrapperUser) AddUserCommand(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.AddUserCommand, callback, &args).AsyncCallWithCallback()
}
func (w *WrapperUser) DeleteUserCommand(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.DeleteUserCommand, callback, &args).AsyncCallWithCallback()
}
func (w *WrapperUser) GetAllUserCommands(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.GetAllUserCommands, callback, &args).AsyncCallWithCallback()
}
