//go:build js && wasm

package event_listener

import (
	"syscall/js"

	"github.com/openimsdk/openim-sdk-core/v3/internal/third/file"
	"github.com/openimsdk/openim-sdk-core/v3/open_im_sdk_callback"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/utils"
	"github.com/openimsdk/openim-sdk-core/v3/sdk_struct"
)

type ConnCallback struct {
	uid string
	CallbackWriter
}

func NewConnCallback(funcName string, callback *js.Value) *ConnCallback {
	return &ConnCallback{CallbackWriter: NewEventData(callback).SetEvent(funcName)}
}

func (i *ConnCallback) OnConnecting() {
	i.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SendMessage()
}

func (i *ConnCallback) OnConnectSuccess() {
	i.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SendMessage()

}
func (i *ConnCallback) OnConnectFailed(errCode int32, errMsg string) {
	i.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SetErrCode(errCode).SetErrMsg(errMsg).SendMessage()
}

func (i *ConnCallback) OnKickedOffline() {
	i.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SendMessage()
}

func (i *ConnCallback) OnUserTokenExpired() {
	i.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SendMessage()
}

func (i *ConnCallback) OnUserTokenInvalid(errMsg string) {
	i.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SetData(errMsg).SendMessage()
}
func (i *ConnCallback) OnUserCommandAdd(userInfo string) {
	i.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SetData(userInfo).SendMessage()
}
func (i *ConnCallback) OnUserCommandDelete(userInfo string) {
	i.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SetData(userInfo).SendMessage()
}
func (i *ConnCallback) OnUserCommandUpdate(userInfo string) {
	i.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SetData(userInfo).SendMessage()
}

type ConversationCallback struct {
	uid string
	CallbackWriter
}

func NewConversationCallback(callback *js.Value) *ConversationCallback {
	return &ConversationCallback{CallbackWriter: NewEventData(callback)}
}
func (c ConversationCallback) OnSyncServerStart(reinstalled bool) {
	c.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SetData(reinstalled).SendMessage()
}

func (c ConversationCallback) OnSyncServerFinish(reinstalled bool) {
	c.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SetData(reinstalled).SendMessage()
}
func (c ConversationCallback) OnSyncServerProgress(progress int) {
	c.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SetData(progress).SendMessage()
}

func (c ConversationCallback) OnSyncServerFailed(reinstalled bool) {
	c.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SetData(reinstalled).SendMessage()

}

func (c ConversationCallback) OnNewConversation(conversationList string) {
	c.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SetData(conversationList).SendMessage()

}

func (c ConversationCallback) OnConversationChanged(conversationList string) {
	c.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SetData(conversationList).SendMessage()

}

func (c ConversationCallback) OnTotalUnreadMessageCountChanged(totalUnreadCount int32) {
	c.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SetData(totalUnreadCount).SendMessage()
}

func (c ConversationCallback) OnConversationUserInputStatusChanged(change string) {
	c.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SetData(change).SendMessage()
}

type AdvancedMsgCallback struct {
	CallbackWriter
}

func NewAdvancedMsgCallback(callback *js.Value) *AdvancedMsgCallback {
	return &AdvancedMsgCallback{CallbackWriter: NewEventData(callback)}
}
func (a AdvancedMsgCallback) OnRecvNewMessage(message string) {
	a.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SetData(message).SendMessage()
}

func (a AdvancedMsgCallback) OnRecvC2CReadReceipt(msgReceiptList string) {
	a.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SetData(msgReceiptList).SendMessage()
}

func (a AdvancedMsgCallback) OnRecvGroupReadReceipt(groupMsgReceiptList string) {
	a.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SetData(groupMsgReceiptList).SendMessage()
}

func (a AdvancedMsgCallback) OnRecvMessageRevoked(msgID string) {
	a.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SetData(msgID).SendMessage()
}

func (a AdvancedMsgCallback) OnNewRecvMessageRevoked(messageRevoked string) {
	a.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SetData(messageRevoked).SendMessage()
}
func (a AdvancedMsgCallback) OnRecvMessageModified(message string) {
	a.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SetData(message).SendMessage()
}
func (a AdvancedMsgCallback) OnRecvMessageExtensionsChanged(clientMsgID string, reactionExtensionList string) {
	m := make(map[string]interface{})
	m["clientMsgID"] = clientMsgID
	m["reactionExtensionList"] = reactionExtensionList
	a.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SetData(utils.StructToJsonString(m)).SendMessage()
}

func (a AdvancedMsgCallback) OnRecvMessageExtensionsDeleted(clientMsgID string, reactionExtensionKeyList string) {
	m := make(map[string]interface{})
	m["clientMsgID"] = clientMsgID
	m["reactionExtensionKeyList"] = reactionExtensionKeyList
	a.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SetData(utils.StructToJsonString(m)).SendMessage()
}
func (a AdvancedMsgCallback) OnRecvMessageExtensionsAdded(clientMsgID string, reactionExtensionList string) {
	m := make(map[string]interface{})
	m["clientMsgID"] = clientMsgID
	m["reactionExtensionList"] = reactionExtensionList
	a.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SetData(utils.StructToJsonString(m)).SendMessage()
}
func (a AdvancedMsgCallback) OnRecvOfflineNewMessage(message string) {
	a.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SetData(message).SendMessage()
}

func (a AdvancedMsgCallback) OnMsgDeleted(message string) {
	a.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SetData(message).SendMessage()
}

func (a AdvancedMsgCallback) OnRecvOnlineOnlyMessage(message string) {
	a.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SetData(message).SendMessage()
}

type BaseCallback struct {
	CallbackWriter
}

func (b *BaseCallback) EventData() CallbackWriter {
	return b.CallbackWriter
}

func NewBaseCallback(funcName string, _ *js.Value) *BaseCallback {
	return &BaseCallback{CallbackWriter: NewPromiseHandler().SetEvent(funcName)}
}

func (b *BaseCallback) OnError(errCode int32, errMsg string) {
	b.CallbackWriter.SetErrCode(errCode).SetErrMsg(errMsg).SendMessage()
}
func (b *BaseCallback) OnSuccess(data string) {
	b.CallbackWriter.SetData(data).SendMessage()
}

type SendMessageCallback struct {
	BaseCallback
	globalEvent CallbackWriter
	clientMsgID string
}

func (s *SendMessageCallback) SetClientMsgID(args *[]js.Value) *SendMessageCallback {
	m := sdk_struct.MsgStruct{}
	utils.JsonStringToStruct((*args)[1].String(), &m)
	s.clientMsgID = m.ClientMsgID
	return s
}
func NewSendMessageCallback(funcName string, callback *js.Value) *SendMessageCallback {
	return &SendMessageCallback{BaseCallback: BaseCallback{CallbackWriter: NewPromiseHandler().SetEvent(funcName)}, globalEvent: NewEventData(callback).SetEvent(funcName)}
}

func (s *SendMessageCallback) OnProgress(progress int) {
	mReply := make(map[string]interface{})
	mReply["progress"] = progress
	mReply["clientMsgID"] = s.clientMsgID
	s.globalEvent.SetEvent(utils.GetSelfFuncName()).SetData(utils.StructToJsonString(mReply)).SendMessage()
}

type UploadInterface interface {
	open_im_sdk_callback.Base
	open_im_sdk_callback.UploadFileCallback
}

var _ UploadInterface = (*UploadFileCallback)(nil)

type UploadFileCallback struct {
	BaseCallback
	globalEvent CallbackWriter
	Uuid        string
}

func NewUploadFileCallback(funcName string, callback *js.Value) *UploadFileCallback {
	return &UploadFileCallback{BaseCallback: BaseCallback{CallbackWriter: NewPromiseHandler().SetEvent(funcName)}, globalEvent: NewEventData(callback).SetEvent(funcName)}
}
func (u *UploadFileCallback) SetUuid(args *[]js.Value) *UploadFileCallback {
	f := file.UploadFileReq{}
	utils.JsonStringToStruct((*args)[1].String(), &f)
	u.Uuid = f.Uuid
	return u
}
func (u *UploadFileCallback) Open(size int64) {
	mReply := make(map[string]interface{})
	mReply["size"] = size
	mReply["uuid"] = u.Uuid
	u.globalEvent.SetEvent(utils.GetSelfFuncName()).SetData(utils.StructToJsonString(mReply)).SendMessage()
}

func (u *UploadFileCallback) PartSize(partSize int64, num int) {
	mReply := make(map[string]interface{})
	mReply["partSize"] = partSize
	mReply["num"] = num
	mReply["uuid"] = u.Uuid
	u.globalEvent.SetEvent(utils.GetSelfFuncName()).SetData(utils.StructToJsonString(mReply)).SendMessage()
}

func (u *UploadFileCallback) HashPartProgress(index int, size int64, partHash string) {
	mReply := make(map[string]interface{})
	mReply["index"] = index
	mReply["size"] = size
	mReply["partHash"] = partHash
	mReply["uuid"] = u.Uuid
	u.globalEvent.SetEvent(utils.GetSelfFuncName()).SetData(utils.StructToJsonString(mReply)).SendMessage()
}

func (u *UploadFileCallback) HashPartComplete(partsHash string, fileHash string) {
	mReply := make(map[string]interface{})
	mReply["partsHash"] = partsHash
	mReply["fileHash"] = fileHash
	mReply["uuid"] = u.Uuid
	u.globalEvent.SetEvent(utils.GetSelfFuncName()).SetData(utils.StructToJsonString(mReply)).SendMessage()
}

func (u *UploadFileCallback) UploadID(uploadID string) {
	mReply := make(map[string]interface{})
	mReply["uploadID"] = uploadID
	mReply["uuid"] = u.Uuid
	u.globalEvent.SetEvent(utils.GetSelfFuncName()).SetData(utils.StructToJsonString(mReply)).SendMessage()
}

func (u *UploadFileCallback) UploadPartComplete(index int, partSize int64, partHash string) {
	mReply := make(map[string]interface{})
	mReply["index"] = index
	mReply["partSize"] = partSize
	mReply["partHash"] = partHash
	mReply["uuid"] = u.Uuid
	u.globalEvent.SetEvent(utils.GetSelfFuncName()).SetData(utils.StructToJsonString(mReply)).SendMessage()
}

func (u *UploadFileCallback) UploadComplete(fileSize int64, streamSize int64, storageSize int64) {
	mReply := make(map[string]interface{})
	mReply["fileSize"] = fileSize
	mReply["streamSize"] = streamSize
	mReply["storageSize"] = storageSize
	mReply["uuid"] = u.Uuid
	u.globalEvent.SetEvent(utils.GetSelfFuncName()).SetData(utils.StructToJsonString(mReply)).SendMessage()
}

func (u *UploadFileCallback) Complete(size int64, url string, typ int) {
	mReply := make(map[string]interface{})
	mReply["size"] = size
	mReply["url"] = url
	mReply["typ"] = typ
	mReply["uuid"] = u.Uuid
	u.globalEvent.SetEvent(utils.GetSelfFuncName()).SetData(utils.StructToJsonString(mReply)).SendMessage()
}

type BatchMessageCallback struct {
	CallbackWriter
}

func NewBatchMessageCallback(callback *js.Value) *BatchMessageCallback {
	return &BatchMessageCallback{CallbackWriter: NewEventData(callback)}
}

func (b *BatchMessageCallback) OnRecvNewMessages(messageList string) {
	b.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SetData(messageList).SendMessage()
}
func (b *BatchMessageCallback) OnRecvOfflineNewMessages(messageList string) {
	b.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SetData(messageList).SendMessage()
}

type FriendCallback struct {
	CallbackWriter
}

func NewFriendCallback(callback *js.Value) *FriendCallback {
	return &FriendCallback{CallbackWriter: NewEventData(callback)}
}

func (f *FriendCallback) OnFriendApplicationAdded(friendApplication string) {
	f.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SetData(friendApplication).SendMessage()
}

func (f *FriendCallback) OnFriendApplicationDeleted(friendApplication string) {
	f.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SetData(friendApplication).SendMessage()
}

func (f *FriendCallback) OnFriendApplicationAccepted(groupApplication string) {
	f.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SetData(groupApplication).SendMessage()

}
func (f *FriendCallback) OnFriendApplicationRejected(friendApplication string) {
	f.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SetData(friendApplication).SendMessage()
}

func (f *FriendCallback) OnFriendAdded(friendInfo string) {
	f.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SetData(friendInfo).SendMessage()
}

func (f *FriendCallback) OnFriendDeleted(friendInfo string) {
	f.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SetData(friendInfo).SendMessage()
}
func (f *FriendCallback) OnFriendInfoChanged(friendInfo string) {
	f.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SetData(friendInfo).SendMessage()
}
func (f *FriendCallback) OnBlackAdded(blackInfo string) {
	f.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SetData(blackInfo).SendMessage()
}

func (f *FriendCallback) OnBlackDeleted(blackInfo string) {
	f.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SetData(blackInfo).SendMessage()
}

type GroupCallback struct {
	CallbackWriter
}

func NewGroupCallback(callback *js.Value) *GroupCallback {
	return &GroupCallback{CallbackWriter: NewEventData(callback)}
}

func (f *GroupCallback) OnJoinedGroupAdded(groupInfo string) {
	f.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SetData(groupInfo).SendMessage()
}
func (f *GroupCallback) OnJoinedGroupDeleted(groupInfo string) {
	f.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SetData(groupInfo).SendMessage()
}
func (f *GroupCallback) OnGroupMemberAdded(groupMemberInfo string) {
	f.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SetData(groupMemberInfo).SendMessage()
}
func (f *GroupCallback) OnGroupMemberDeleted(groupMemberInfo string) {
	f.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SetData(groupMemberInfo).SendMessage()
}
func (f *GroupCallback) OnGroupApplicationAdded(groupApplication string) {
	f.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SetData(groupApplication).SendMessage()
}
func (f *GroupCallback) OnGroupApplicationDeleted(groupApplication string) {
	f.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SetData(groupApplication).SendMessage()
}
func (f *GroupCallback) OnGroupInfoChanged(groupInfo string) {
	f.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SetData(groupInfo).SendMessage()
}
func (f *GroupCallback) OnGroupMemberInfoChanged(groupMemberInfo string) {
	f.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SetData(groupMemberInfo).SendMessage()
}
func (f *GroupCallback) OnGroupApplicationAccepted(groupApplication string) {
	f.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SetData(groupApplication).SendMessage()
}
func (f *GroupCallback) OnGroupApplicationRejected(groupApplication string) {
	f.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SetData(groupApplication).SendMessage()
}
func (f *GroupCallback) OnGroupDismissed(groupInfo string) {
	f.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SetData(groupInfo).SendMessage()
}

type UserCallback struct {
	CallbackWriter
}

func (u UserCallback) OnUserStatusChanged(statusMap string) {
	u.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SetData(statusMap).SendMessage()
}

func (u UserCallback) OnUserInputStatusChanged(change string) {
	u.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SetData(change).SendMessage()
}

func NewUserCallback(callback *js.Value) *UserCallback {
	return &UserCallback{CallbackWriter: NewEventData(callback)}
}
func (u UserCallback) OnSelfInfoUpdated(userInfo string) {
	u.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SetData(userInfo).SendMessage()
}
func (u UserCallback) OnUserCommandAdd(userInfo string) {
	u.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SetData(userInfo).SendMessage()
}
func (u UserCallback) OnUserCommandDelete(userInfo string) {
	u.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SetData(userInfo).SendMessage()
}
func (u UserCallback) OnUserCommandUpdate(userInfo string) {
	u.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SetData(userInfo).SendMessage()
}

type CustomBusinessCallback struct {
	CallbackWriter
}

func NewCustomBusinessCallback(callback *js.Value) *CustomBusinessCallback {
	return &CustomBusinessCallback{CallbackWriter: NewEventData(callback)}
}

func (c CustomBusinessCallback) OnRecvCustomBusinessMessage(businessMessage string) {
	c.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SetData(businessMessage).SendMessage()

}

type SignalingCallback struct {
	CallbackWriter
}

func (s SignalingCallback) OnRoomParticipantConnected(onRoomParticipantConnectedCallback string) {
	s.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SetData(onRoomParticipantConnectedCallback).SendMessage()
}

func (s SignalingCallback) OnRoomParticipantDisconnected(onRoomParticipantDisconnectedCallback string) {
	s.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SetData(onRoomParticipantDisconnectedCallback).SendMessage()
}

func (s SignalingCallback) OnReceiveNewInvitation(receiveNewInvitationCallback string) {
	s.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SetData(receiveNewInvitationCallback).SendMessage()
}

func (s SignalingCallback) OnInviteeAccepted(inviteeAcceptedCallback string) {
	s.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SetData(inviteeAcceptedCallback).SendMessage()

}
func (s SignalingCallback) OnInviteeAcceptedByOtherDevice(inviteeAcceptedCallback string) {
	s.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SetData(inviteeAcceptedCallback).SendMessage()
}

func (s SignalingCallback) OnInviteeRejected(inviteeRejectedCallback string) {
	s.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SetData(inviteeRejectedCallback).SendMessage()
}

func (s SignalingCallback) OnInviteeRejectedByOtherDevice(inviteeRejectedCallback string) {
	s.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SetData(inviteeRejectedCallback).SendMessage()
}

func (s SignalingCallback) OnInvitationCancelled(invitationCancelledCallback string) {
	s.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SetData(invitationCancelledCallback).SendMessage()
}

func (s SignalingCallback) OnInvitationTimeout(invitationTimeoutCallback string) {
	s.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SetData(invitationTimeoutCallback).SendMessage()
}

func (s SignalingCallback) OnHangUp(hangUpCallback string) {
	s.CallbackWriter.SetEvent(utils.GetSelfFuncName()).SetData(hangUpCallback).SendMessage()
}

func NewSignalingCallback(callback *js.Value) *SignalingCallback {
	return &SignalingCallback{CallbackWriter: NewEventData(callback)}
}
