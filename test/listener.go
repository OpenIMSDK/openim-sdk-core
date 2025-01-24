package test

import (
	"context"
	"fmt"

	"github.com/openimsdk/tools/log"
)

type OnConnListener struct{}

func (c *OnConnListener) OnUserTokenInvalid(errMsg string) {}

func (c *OnConnListener) OnConnecting() {
	// fmt.Println("OnConnecting")
}

func (c *OnConnListener) OnConnectSuccess() {
	// fmt.Println("OnConnectSuccess")
}

func (c *OnConnListener) OnConnectFailed(errCode int32, errMsg string) {
	// fmt.Println("OnConnectFailed")
}

func (c *OnConnListener) OnKickedOffline() {
	// fmt.Println("OnKickedOffline")
}

func (c *OnConnListener) OnUserTokenExpired() {
	// fmt.Println("OnUserTokenExpired")
}

type onConversationListener struct {
	ctx context.Context
	ch  chan error
}

func (o *onConversationListener) OnSyncServerStart(reinstalled bool) {
	log.ZInfo(o.ctx, "OnSyncServerStart")
}

func (o *onConversationListener) OnSyncServerFinish(reinstalled bool) {
	log.ZInfo(o.ctx, "OnSyncServerFinish")
	o.ch <- nil
}

func (o *onConversationListener) OnSyncServerFailed(reinstalled bool) {
	log.ZInfo(o.ctx, "OnSyncServerFailed")
	o.ch <- fmt.Errorf("OnSyncServerFailed")
}

func (o *onConversationListener) OnSyncServerProgress(progress int) {
	log.ZInfo(o.ctx, "OnSyncServerProgress", "progress", progress)
}

func (o *onConversationListener) OnNewConversation(conversationList string) {
	log.ZInfo(o.ctx, "OnNewConversation", "conversationList", conversationList)
}

func (o *onConversationListener) OnConversationChanged(conversationList string) {
	log.ZInfo(o.ctx, "OnConversationChanged", "conversationList", conversationList)
}

func (o *onConversationListener) OnTotalUnreadMessageCountChanged(totalUnreadCount int32) {
	log.ZInfo(o.ctx, "OnTotalUnreadMessageCountChanged", "totalUnreadCount", totalUnreadCount)
}

func (o *onConversationListener) OnConversationUserInputStatusChanged(change string) {
	log.ZInfo(o.ctx, "OnConversationUserInputStatusChanged", "change", change)
}

type onGroupListener struct {
	ctx context.Context
}

func (o *onGroupListener) OnGroupDismissed(groupInfo string) {
	log.ZInfo(o.ctx, "OnGroupDismissed", "groupInfo", groupInfo)
}

func (o *onGroupListener) OnJoinedGroupAdded(groupInfo string) {
	log.ZInfo(o.ctx, "OnJoinedGroupAdded", "groupInfo", groupInfo)
}

func (o *onGroupListener) OnJoinedGroupDeleted(groupInfo string) {
	log.ZInfo(o.ctx, "OnJoinedGroupDeleted", "groupInfo", groupInfo)
}

func (o *onGroupListener) OnGroupMemberAdded(groupMemberInfo string) {
	log.ZInfo(o.ctx, "OnGroupMemberAdded", "groupMemberInfo", groupMemberInfo)
}

func (o *onGroupListener) OnGroupMemberDeleted(groupMemberInfo string) {
	log.ZInfo(o.ctx, "OnGroupMemberDeleted", "groupMemberInfo", groupMemberInfo)
}

func (o *onGroupListener) OnGroupApplicationAdded(groupApplication string) {
	log.ZInfo(o.ctx, "OnGroupApplicationAdded", "groupApplication", groupApplication)
}

func (o *onGroupListener) OnGroupApplicationDeleted(groupApplication string) {
	log.ZInfo(o.ctx, "OnGroupApplicationDeleted", "groupApplication", groupApplication)
}

func (o *onGroupListener) OnGroupInfoChanged(groupInfo string) {
	log.ZInfo(o.ctx, "OnGroupInfoChanged", "groupInfo", groupInfo)
}

func (o *onGroupListener) OnGroupMemberInfoChanged(groupMemberInfo string) {
	log.ZInfo(o.ctx, "OnGroupMemberInfoChanged", "groupMemberInfo", groupMemberInfo)
}

func (o *onGroupListener) OnGroupApplicationAccepted(groupApplication string) {
	log.ZInfo(o.ctx, "OnGroupApplicationAccepted", "groupApplication", groupApplication)
}

func (o *onGroupListener) OnGroupApplicationRejected(groupApplication string) {
	log.ZInfo(o.ctx, "OnGroupApplicationRejected", "groupApplication", groupApplication)
}

type onAdvancedMsgListener struct {
	ctx context.Context
}

func (o *onAdvancedMsgListener) OnRecvOnlineOnlyMessage(message string) {
	log.ZDebug(o.ctx, "OnRecvOnlineOnlyMessage", "message", message)
}

func (o *onAdvancedMsgListener) OnRecvOfflineNewMessage(message string) {
	//TODO implement me
	panic("implement me")
}

func (o *onAdvancedMsgListener) OnMsgDeleted(message string) {
	log.ZInfo(o.ctx, "OnMsgDeleted", "message", message)
}

func (o *onAdvancedMsgListener) OnRecvOfflineNewMessages(messageList string) {
	log.ZInfo(o.ctx, "OnRecvOfflineNewMessages", "messageList", messageList)
}

func (o *onAdvancedMsgListener) OnRecvNewMessage(message string) {
	log.ZInfo(o.ctx, "OnRecvNewMessage", "message", message)
}

func (o *onAdvancedMsgListener) OnRecvC2CReadReceipt(msgReceiptList string) {
	log.ZInfo(o.ctx, "OnRecvC2CReadReceipt", "msgReceiptList", msgReceiptList)
}

func (o *onAdvancedMsgListener) OnRecvGroupReadReceipt(groupMsgReceiptList string) {
	log.ZInfo(o.ctx, "OnRecvGroupReadReceipt", "groupMsgReceiptList", groupMsgReceiptList)
}

func (o *onAdvancedMsgListener) OnRecvMessageRevoked(msgID string) {
	log.ZInfo(o.ctx, "OnRecvMessageRevoked", "msgID", msgID)
}

func (o *onAdvancedMsgListener) OnNewRecvMessageRevoked(messageRevoked string) {
	log.ZInfo(o.ctx, "OnNewRecvMessageRevoked", "messageRevoked", messageRevoked)
}

func (o *onAdvancedMsgListener) OnRecvMessageExtensionsChanged(msgID string, reactionExtensionList string) {
	log.ZInfo(o.ctx, "OnRecvMessageExtensionsChanged", "msgID", msgID, "reactionExtensionList", reactionExtensionList)
}

func (o *onAdvancedMsgListener) OnRecvMessageExtensionsDeleted(msgID string, reactionExtensionKeyList string) {
	log.ZInfo(o.ctx, "OnRecvMessageExtensionsDeleted", "msgID", msgID, "reactionExtensionKeyList", reactionExtensionKeyList)
}

func (o *onAdvancedMsgListener) OnRecvMessageExtensionsAdded(msgID string, reactionExtensionList string) {
	log.ZInfo(o.ctx, "OnRecvMessageExtensionsAdded", "msgID", msgID, "reactionExtensionList", reactionExtensionList)
}

type onFriendshipListener struct {
	ctx context.Context
}

func (o *onFriendshipListener) OnFriendApplicationAdded(friendApplication string) {
	log.ZDebug(context.Background(), "OnFriendApplicationAdded", "friendApplication", friendApplication)
}

func (o *onFriendshipListener) OnFriendApplicationDeleted(friendApplication string) {
	log.ZDebug(context.Background(), "OnFriendApplicationDeleted", "friendApplication", friendApplication)
}

func (o *onFriendshipListener) OnFriendApplicationAccepted(friendApplication string) {
	log.ZDebug(context.Background(), "OnFriendApplicationAccepted", "friendApplication", friendApplication)
}

func (o *onFriendshipListener) OnFriendApplicationRejected(friendApplication string) {
	log.ZDebug(context.Background(), "OnFriendApplicationRejected", "friendApplication", friendApplication)
}

func (o *onFriendshipListener) OnFriendAdded(friendInfo string) {
	log.ZDebug(context.Background(), "OnFriendAdded", "friendInfo", friendInfo)
}

func (o *onFriendshipListener) OnFriendDeleted(friendInfo string) {
	log.ZDebug(context.Background(), "OnFriendDeleted", "friendInfo", friendInfo)
}

func (o *onFriendshipListener) OnFriendInfoChanged(friendInfo string) {
	log.ZDebug(context.Background(), "OnFriendInfoChanged", "friendInfo", friendInfo)
}

func (o *onFriendshipListener) OnBlackAdded(blackInfo string) {
	log.ZDebug(context.Background(), "OnBlackAdded", "blackInfo", blackInfo)
}

func (o *onFriendshipListener) OnBlackDeleted(blackInfo string) {
	log.ZDebug(context.Background(), "OnBlackDeleted", "blackInfo", blackInfo)
}

type onUserListener struct {
	ctx context.Context
}

func (o *onUserListener) OnSelfInfoUpdated(userInfo string) {
	log.ZDebug(context.Background(), "OnSelfInfoUpdated", "userInfo", userInfo)
}
func (o *onUserListener) OnUserCommandAdd(userInfo string) {
	log.ZDebug(context.Background(), "OnUserCommandAdd", "blackInfo", userInfo)
}
func (o *onUserListener) OnUserCommandDelete(userInfo string) {
	log.ZDebug(context.Background(), "OnUserCommandDelete", "blackInfo", userInfo)
}
func (o *onUserListener) OnUserCommandUpdate(userInfo string) {
	log.ZDebug(context.Background(), "OnUserCommandUpdate", "blackInfo", userInfo)
}
func (o *onUserListener) OnUserStatusChanged(statusMap string) {
	log.ZDebug(context.Background(), "OnUserStatusChanged", "OnUserStatusChanged", statusMap)
}
