//go:build js && wasm

package wasm_wrapper

import (
	"syscall/js"

	"github.com/openimsdk/openim-sdk-core/v3/open_im_sdk"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/utils"
	"github.com/openimsdk/openim-sdk-core/v3/wasm/event_listener"
)

// ------------------------------------group---------------------------
type WrapperGroup struct {
	*WrapperCommon
}

func NewWrapperGroup(wrapperCommon *WrapperCommon) *WrapperGroup {
	return &WrapperGroup{WrapperCommon: wrapperCommon}
}

func (w *WrapperGroup) CreateGroup(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.CreateGroup, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperGroup) JoinGroup(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.JoinGroup, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperGroup) QuitGroup(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.QuitGroup, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperGroup) DismissGroup(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.DismissGroup, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperGroup) ChangeGroupMute(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.ChangeGroupMute, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperGroup) ChangeGroupMemberMute(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.ChangeGroupMemberMute, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperGroup) SetGroupMemberInfo(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.SetGroupMemberInfo, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperGroup) GetJoinedGroupList(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.GetJoinedGroupList, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperGroup) GetJoinedGroupListPage(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.GetJoinedGroupListPage, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperGroup) SearchGroups(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.SearchGroups, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperGroup) SetGroupInfo(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.SetGroupInfo, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperGroup) GetGroupMemberList(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.GetGroupMemberList, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperGroup) GetGroupMemberOwnerAndAdmin(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.GetGroupMemberOwnerAndAdmin, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperGroup) GetGroupMemberListByJoinTimeFilter(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.GetGroupMemberListByJoinTimeFilter, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperGroup) GetSpecifiedGroupMembersInfo(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.GetSpecifiedGroupMembersInfo, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperGroup) KickGroupMember(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.KickGroupMember, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperGroup) TransferGroupOwner(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.TransferGroupOwner, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperGroup) InviteUserToGroup(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.InviteUserToGroup, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperGroup) GetGroupApplicationListAsRecipient(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.GetGroupApplicationListAsRecipient, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperGroup) GetGroupApplicationListAsApplicant(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.GetGroupApplicationListAsApplicant, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperGroup) AcceptGroupApplication(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.AcceptGroupApplication, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperGroup) RefuseGroupApplication(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.RefuseGroupApplication, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperGroup) SearchGroupMembers(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.SearchGroupMembers, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperGroup) IsJoinGroup(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.IsJoinGroup, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperGroup) GetUsersInGroup(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.GetUsersInGroup, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperGroup) GetSpecifiedGroupsInfo(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.GetSpecifiedGroupsInfo, callback, &args).AsyncCallWithCallback()
}
