package open_im_sdk

import (
	"context"
	"errors"
	"github.com/OpenIMSDK/tools/log"
	"github.com/openimsdk/openim-sdk-core/v3/open_im_sdk_callback"
)

var ErrNotImplemented = errors.New("not set listener")

type emptyGroupListener struct {
	ctx context.Context
}

func newEmptyGroupListener(ctx context.Context) open_im_sdk_callback.OnGroupListener {
	return &emptyGroupListener{ctx}
}

func (e *emptyGroupListener) OnJoinedGroupAdded(groupInfo string) {
	log.ZWarn(e.ctx, "GroupListener is not implemented", nil, "groupInfo", groupInfo)
}

func (e *emptyGroupListener) OnJoinedGroupDeleted(groupInfo string) {
	log.ZWarn(e.ctx, "GroupListener is not implemented", nil, "groupInfo", groupInfo)

}

func (e *emptyGroupListener) OnGroupMemberAdded(groupMemberInfo string) {
	log.ZWarn(e.ctx, "GroupListener is not implemented", nil, "groupMemberInfo", groupMemberInfo)

}

func (e *emptyGroupListener) OnGroupMemberDeleted(groupMemberInfo string) {
	log.ZWarn(e.ctx, "GroupListener is not implemented", nil, "groupMemberInfo", groupMemberInfo)

}

func (e *emptyGroupListener) OnGroupApplicationAdded(groupApplication string) {
	log.ZWarn(e.ctx, "GroupListener is not implemented", nil, "groupApplication", groupApplication)

}

func (e *emptyGroupListener) OnGroupApplicationDeleted(groupApplication string) {
	log.ZWarn(e.ctx, "GroupListener is not implemented", nil, "groupApplication", groupApplication)

}

func (e *emptyGroupListener) OnGroupInfoChanged(groupInfo string) {
	log.ZWarn(e.ctx, "GroupListener is not implemented", nil, "groupInfo", groupInfo)

}

func (e *emptyGroupListener) OnGroupDismissed(groupInfo string) {
	log.ZWarn(e.ctx, "GroupListener is not implemented", nil, "groupInfo", groupInfo)
}

func (e *emptyGroupListener) OnGroupMemberInfoChanged(groupMemberInfo string) {
	log.ZWarn(e.ctx, "GroupListener is not implemented", nil, "groupMemberInfo", groupMemberInfo)
}

func (e *emptyGroupListener) OnGroupApplicationAccepted(groupApplication string) {
	log.ZWarn(e.ctx, "GroupListener is not implemented", nil, "groupApplication", groupApplication)
}

func (e *emptyGroupListener) OnGroupApplicationRejected(groupApplication string) {
	log.ZWarn(e.ctx, "GroupListener is not implemented", nil, "groupApplication", groupApplication)
}

type emptyFriendshipListener struct {
	ctx context.Context
}

func newEmptyFriendshipListener(ctx context.Context) open_im_sdk_callback.OnFriendshipListener {
	return &emptyFriendshipListener{ctx}
}

func (e *emptyFriendshipListener) OnFriendApplicationAdded(friendApplication string) {
	log.ZWarn(e.ctx, "FriendshipListener is not implemented", nil,
		"friendApplication", friendApplication)
}

func (e *emptyFriendshipListener) OnFriendApplicationDeleted(friendApplication string) {
	log.ZWarn(e.ctx, "FriendshipListener is not implemented", nil,
		"friendApplication", friendApplication)
}

func (e *emptyFriendshipListener) OnFriendApplicationAccepted(friendApplication string) {
	log.ZWarn(e.ctx, "FriendshipListener is not implemented", nil,
		"friendApplication", friendApplication)
}

func (e *emptyFriendshipListener) OnFriendApplicationRejected(friendApplication string) {
	log.ZWarn(e.ctx, "FriendshipListener is not implemented", nil,
		"friendApplication", friendApplication)
}

func (e *emptyFriendshipListener) OnFriendAdded(friendInfo string) {
	log.ZWarn(e.ctx, "FriendshipListener is not implemented", nil,
		"friendInfo", friendInfo)
}

func (e *emptyFriendshipListener) OnFriendDeleted(friendInfo string) {
	log.ZWarn(e.ctx, "FriendshipListener is not implemented", nil,
		"friendInfo", friendInfo)
}

func (e *emptyFriendshipListener) OnFriendInfoChanged(friendInfo string) {
	log.ZWarn(e.ctx, "FriendshipListener is not implemented", nil,
		"friendInfo", friendInfo)
}

func (e *emptyFriendshipListener) OnBlackAdded(blackInfo string) {
	log.ZWarn(e.ctx, "FriendshipListener is not implemented", nil,
		"blackInfo", blackInfo)
}

func (e *emptyFriendshipListener) OnBlackDeleted(blackInfo string) {
	log.ZWarn(e.ctx, "FriendshipListener is not implemented", nil,
		"blackInfo", blackInfo)
}

type emptyConversationListener struct {
	ctx context.Context
}

func newEmptyConversationListener(ctx context.Context) open_im_sdk_callback.OnConversationListener {
	return &emptyConversationListener{ctx: ctx}
}

func (e *emptyConversationListener) OnSyncServerStart() {

	log.ZWarn(e.ctx, "ConversationListener is not implemented", nil)
}

func (e *emptyConversationListener) OnSyncServerFinish() {
	log.ZWarn(e.ctx, "ConversationListener is not implemented", nil)

}

func (e *emptyConversationListener) OnSyncServerFailed() {

	log.ZWarn(e.ctx, "ConversationListener is not implemented", nil)
}

func (e *emptyConversationListener) OnNewConversation(conversationList string) {
	log.ZWarn(e.ctx, "ConversationListener is not implemented", nil,
		"conversationList", conversationList)

}

func (e *emptyConversationListener) OnConversationChanged(conversationList string) {

	log.ZWarn(e.ctx, "ConversationListener is not implemented", nil,
		"conversationList", conversationList)
}

func (e *emptyConversationListener) OnTotalUnreadMessageCountChanged(totalUnreadCount int32) {
	log.ZWarn(e.ctx, "ConversationListener is not implemented", nil,
		"totalUnreadCount", totalUnreadCount)
}

type emptyAdvancedMsgListener struct {
	ctx context.Context
}

func newEmptyAdvancedMsgListener(ctx context.Context) open_im_sdk_callback.OnAdvancedMsgListener {
	return &emptyAdvancedMsgListener{ctx}
}

func (e *emptyAdvancedMsgListener) OnRecvNewMessage(message string) {
	log.ZWarn(e.ctx, "AdvancedMsgListener is not implemented", nil, "message", message)
}

func (e *emptyAdvancedMsgListener) OnRecvC2CReadReceipt(msgReceiptList string) {

	log.ZWarn(e.ctx, "AdvancedMsgListener is not implemented", nil,
		"msgReceiptList", msgReceiptList)
}

func (e *emptyAdvancedMsgListener) OnRecvGroupReadReceipt(groupMsgReceiptList string) {
	log.ZWarn(e.ctx, "AdvancedMsgListener is not implemented", nil,
		"groupMsgReceiptList", groupMsgReceiptList)
}

func (e *emptyAdvancedMsgListener) OnNewRecvMessageRevoked(messageRevoked string) {
	log.ZWarn(e.ctx, "AdvancedMsgListener is not implemented", nil, "messageRevoked", messageRevoked)
}

func (e *emptyAdvancedMsgListener) OnRecvMessageExtensionsChanged(msgID string, reactionExtensionList string) {
	log.ZWarn(e.ctx, "AdvancedMsgListener is not implemented", nil, "msgID", msgID,
		"reactionExtensionList", reactionExtensionList)
}

func (e *emptyAdvancedMsgListener) OnRecvMessageExtensionsDeleted(msgID string, reactionExtensionKeyList string) {
	log.ZWarn(e.ctx, "AdvancedMsgListener is not implemented", nil, "msgID", msgID,
		"reactionExtensionKeyList", reactionExtensionKeyList)
}

func (e *emptyAdvancedMsgListener) OnRecvMessageExtensionsAdded(msgID string, reactionExtensionList string) {
	log.ZWarn(e.ctx, "AdvancedMsgListener is not implemented", nil, "msgID", msgID,
		"reactionExtensionList", reactionExtensionList)
}

func (e *emptyAdvancedMsgListener) OnRecvOfflineNewMessage(message string) {
	log.ZWarn(e.ctx, "AdvancedMsgListener is not implemented", nil, "message", message)
}

func (e *emptyAdvancedMsgListener) OnMsgDeleted(message string) {
	log.ZWarn(e.ctx, "AdvancedMsgListener is not implemented", nil, "message", message)
}

type emptyBatchMsgListener struct{}

func newEmptyBatchMsgListener() *emptyBatchMsgListener {
	return &emptyBatchMsgListener{}
}

func (e *emptyBatchMsgListener) OnRecvNewMessages(messageList string) {

}

func (e *emptyBatchMsgListener) OnRecvOfflineNewMessages(messageList string) {

}

type emptyUserListener struct {
	ctx context.Context
}

func newEmptyUserListener(ctx context.Context) open_im_sdk_callback.OnUserListener {
	return &emptyUserListener{ctx: ctx}
}

func (e *emptyUserListener) OnSelfInfoUpdated(userInfo string) {
	log.ZWarn(e.ctx, "UserListener is not implemented", nil, "userInfo", userInfo)
}

func (e *emptyUserListener) OnUserStatusChanged(statusMap string) {
	log.ZWarn(e.ctx, "UserListener is not implemented", nil, "statusMap", statusMap)
}

type emptyCustomBusinessListener struct {
	ctx context.Context
}

func newEmptyCustomBusinessListener(ctx context.Context) open_im_sdk_callback.OnCustomBusinessListener {
	return &emptyCustomBusinessListener{ctx: ctx}
}

func (e *emptyCustomBusinessListener) OnRecvCustomBusinessMessage(businessMessage string) {
	log.ZWarn(e.ctx, "CustomBusinessListener is not implemented", nil,
		"businessMessage", businessMessage)

}
