//go:build js && wasm

package wasm_wrapper

import (
	"syscall/js"

	"github.com/openimsdk/openim-sdk-core/v3/open_im_sdk"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/utils"
	"github.com/openimsdk/openim-sdk-core/v3/wasm/event_listener"
)

// ------------------------------------group---------------------------
type WrapperFriend struct {
	*WrapperCommon
}

func NewWrapperFriend(wrapperCommon *WrapperCommon) *WrapperFriend {
	return &WrapperFriend{WrapperCommon: wrapperCommon}
}

func (w *WrapperFriend) GetSpecifiedFriendsInfo(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.GetSpecifiedFriendsInfo, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperFriend) GetFriendList(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.GetFriendList, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperFriend) GetFriendListPage(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.GetFriendListPage, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperFriend) SearchFriends(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.SearchFriends, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperFriend) CheckFriend(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.CheckFriend, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperFriend) AddFriend(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.AddFriend, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperFriend) UpdateFriends(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.UpdateFriends, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperFriend) DeleteFriend(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.DeleteFriend, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperFriend) GetFriendApplicationListAsRecipient(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.GetFriendApplicationListAsRecipient, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperFriend) GetFriendApplicationListAsApplicant(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.GetFriendApplicationListAsApplicant, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperFriend) AcceptFriendApplication(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.AcceptFriendApplication, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperFriend) RefuseFriendApplication(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.RefuseFriendApplication, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperFriend) GetBlackList(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.GetBlackList, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperFriend) RemoveBlack(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.RemoveBlack, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperFriend) AddBlack(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.AddBlack, callback, &args).AsyncCallWithCallback()
}
