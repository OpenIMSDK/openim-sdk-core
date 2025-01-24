//go:build js && wasm

package main

import (
	"fmt"
	"runtime"
	"runtime/debug"

	"github.com/openimsdk/openim-sdk-core/v3/wasm/wasm_wrapper"

	"syscall/js"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("MAIN", "panic info is:", r, debug.Stack())
		}
	}()
	fmt.Println("runtime env", runtime.GOARCH, runtime.GOOS)
	registerFunc()
	<-make(chan bool)
}

func registerFunc() {
	//register global listener function
	globalFuc := wasm_wrapper.NewWrapperCommon()
	js.Global().Set(wasm_wrapper.COMMONEVENTFUNC, js.FuncOf(globalFuc.CommonEventFunc))
	//register init login function
	wrapperInitLogin := wasm_wrapper.NewWrapperInitLogin(globalFuc)
	js.Global().Set("initSDK", js.FuncOf(wrapperInitLogin.InitSDK))
	js.Global().Set("login", js.FuncOf(wrapperInitLogin.Login))
	js.Global().Set("logout", js.FuncOf(wrapperInitLogin.Logout))
	js.Global().Set("getLoginStatus", js.FuncOf(wrapperInitLogin.GetLoginStatus))
	js.Global().Set("setAppBackgroundStatus", js.FuncOf(wrapperInitLogin.SetAppBackgroundStatus))
	js.Global().Set("networkStatusChanged", js.FuncOf(wrapperInitLogin.NetworkStatusChanged))
	//register conversation and message function
	wrapperConMsg := wasm_wrapper.NewWrapperConMsg(globalFuc)
	js.Global().Set("createTextMessage", js.FuncOf(wrapperConMsg.CreateTextMessage))
	js.Global().Set("createImageMessage", js.FuncOf(wrapperConMsg.CreateImageMessage))
	js.Global().Set("createImageMessageByURL", js.FuncOf(wrapperConMsg.CreateImageMessageByURL))
	js.Global().Set("createSoundMessageByURL", js.FuncOf(wrapperConMsg.CreateSoundMessageByURL))
	js.Global().Set("createVideoMessageByURL", js.FuncOf(wrapperConMsg.CreateVideoMessageByURL))
	js.Global().Set("createFileMessageByURL", js.FuncOf(wrapperConMsg.CreateFileMessageByURL))
	js.Global().Set("createCustomMessage", js.FuncOf(wrapperConMsg.CreateCustomMessage))
	js.Global().Set("createQuoteMessage", js.FuncOf(wrapperConMsg.CreateQuoteMessage))
	js.Global().Set("createAdvancedQuoteMessage", js.FuncOf(wrapperConMsg.CreateAdvancedQuoteMessage))
	js.Global().Set("createAdvancedTextMessage", js.FuncOf(wrapperConMsg.CreateAdvancedTextMessage))
	js.Global().Set("createCardMessage", js.FuncOf(wrapperConMsg.CreateCardMessage))
	js.Global().Set("createTextAtMessage", js.FuncOf(wrapperConMsg.CreateTextAtMessage))
	js.Global().Set("createVideoMessage", js.FuncOf(wrapperConMsg.CreateVideoMessage))
	js.Global().Set("createFileMessage", js.FuncOf(wrapperConMsg.CreateFileMessage))
	js.Global().Set("createMergerMessage", js.FuncOf(wrapperConMsg.CreateMergerMessage))
	js.Global().Set("createFaceMessage", js.FuncOf(wrapperConMsg.CreateFaceMessage))
	js.Global().Set("createForwardMessage", js.FuncOf(wrapperConMsg.CreateForwardMessage))
	js.Global().Set("createLocationMessage", js.FuncOf(wrapperConMsg.CreateLocationMessage))
	js.Global().Set("createVideoMessageFromFullPath", js.FuncOf(wrapperConMsg.CreateVideoMessageFromFullPath))
	js.Global().Set("createImageMessageFromFullPath", js.FuncOf(wrapperConMsg.CreateImageMessageFromFullPath))

	js.Global().Set("createSoundMessageFromFullPath", js.FuncOf(wrapperConMsg.CreateSoundMessageFromFullPath))
	js.Global().Set("createFileMessageFromFullPath", js.FuncOf(wrapperConMsg.CreateFileMessageFromFullPath))
	js.Global().Set("createSoundMessage", js.FuncOf(wrapperConMsg.CreateSoundMessage))
	js.Global().Set("createForwardMessage", js.FuncOf(wrapperConMsg.CreateForwardMessage))
	js.Global().Set("createLocationMessage", js.FuncOf(wrapperConMsg.CreateLocationMessage))
	js.Global().Set("createVideoMessageFromFullPath", js.FuncOf(wrapperConMsg.CreateVideoMessageFromFullPath))
	js.Global().Set("createImageMessageFromFullPath", js.FuncOf(wrapperConMsg.CreateImageMessageFromFullPath))
	js.Global().Set("getAtAllTag", js.FuncOf(wrapperConMsg.GetAtAllTag))
	js.Global().Set("markConversationMessageAsRead", js.FuncOf(wrapperConMsg.MarkConversationMessageAsRead))
	js.Global().Set("markAllConversationMessageAsRead", js.FuncOf(wrapperConMsg.MarkAllConversationMessageAsRead))
	js.Global().Set("markMessagesAsReadByMsgID", js.FuncOf(wrapperConMsg.MarkMessagesAsReadByMsgID))
	js.Global().Set("sendMessage", js.FuncOf(wrapperConMsg.SendMessage))
	js.Global().Set("sendMessageNotOss", js.FuncOf(wrapperConMsg.SendMessageNotOss))
	js.Global().Set("getAllConversationList", js.FuncOf(wrapperConMsg.GetAllConversationList))
	js.Global().Set("getConversationListSplit", js.FuncOf(wrapperConMsg.GetConversationListSplit))
	js.Global().Set("getOneConversation", js.FuncOf(wrapperConMsg.GetOneConversation))
	js.Global().Set("deleteConversationAndDeleteAllMsg", js.FuncOf(wrapperConMsg.DeleteConversationAndDeleteAllMsg))
	js.Global().Set("getAdvancedHistoryMessageList", js.FuncOf(wrapperConMsg.GetAdvancedHistoryMessageList))
	js.Global().Set("getAdvancedHistoryMessageListReverse", js.FuncOf(wrapperConMsg.GetAdvancedHistoryMessageListReverse))
	js.Global().Set("getMultipleConversation", js.FuncOf(wrapperConMsg.GetMultipleConversation))
	js.Global().Set("hideConversation", js.FuncOf(wrapperConMsg.HideConversation))
	js.Global().Set("setConversationDraft", js.FuncOf(wrapperConMsg.SetConversationDraft))
	js.Global().Set("setConversation", js.FuncOf(wrapperConMsg.SetConversation))

	js.Global().Set("getTotalUnreadMsgCount", js.FuncOf(wrapperConMsg.GetTotalUnreadMsgCount))
	js.Global().Set("findMessageList", js.FuncOf(wrapperConMsg.FindMessageList))

	js.Global().Set("revokeMessage", js.FuncOf(wrapperConMsg.RevokeMessage))
	js.Global().Set("typingStatusUpdate", js.FuncOf(wrapperConMsg.TypingStatusUpdate))
	js.Global().Set("deleteMessageFromLocalStorage", js.FuncOf(wrapperConMsg.DeleteMessageFromLocalStorage))
	js.Global().Set("deleteMessage", js.FuncOf(wrapperConMsg.DeleteMessage))
	js.Global().Set("hideAllConversations", js.FuncOf(wrapperConMsg.HideAllConversations))
	js.Global().Set("deleteAllMsgFromLocalAndSvr", js.FuncOf(wrapperConMsg.DeleteAllMsgFromLocalAndSvr))
	js.Global().Set("deleteAllMsgFromLocal", js.FuncOf(wrapperConMsg.DeleteAllMsgFromLocal))
	js.Global().Set("clearConversationAndDeleteAllMsg", js.FuncOf(wrapperConMsg.ClearConversationAndDeleteAllMsg))
	js.Global().Set("insertSingleMessageToLocalStorage", js.FuncOf(wrapperConMsg.InsertSingleMessageToLocalStorage))
	js.Global().Set("insertGroupMessageToLocalStorage", js.FuncOf(wrapperConMsg.InsertGroupMessageToLocalStorage))
	js.Global().Set("searchLocalMessages", js.FuncOf(wrapperConMsg.SearchLocalMessages))
	js.Global().Set("setMessageLocalEx", js.FuncOf(wrapperConMsg.SetMessageLocalEx))
	js.Global().Set("searchConversation", js.FuncOf(wrapperConMsg.SearchConversation))

	js.Global().Set("changeInputStates", js.FuncOf(wrapperConMsg.ChangeInputStates))
	js.Global().Set("getInputStates", js.FuncOf(wrapperConMsg.GetInputStates))

	//register group func
	wrapperGroup := wasm_wrapper.NewWrapperGroup(globalFuc)
	js.Global().Set("createGroup", js.FuncOf(wrapperGroup.CreateGroup))
	js.Global().Set("getSpecifiedGroupsInfo", js.FuncOf(wrapperGroup.GetSpecifiedGroupsInfo))
	js.Global().Set("joinGroup", js.FuncOf(wrapperGroup.JoinGroup))
	js.Global().Set("quitGroup", js.FuncOf(wrapperGroup.QuitGroup))
	js.Global().Set("dismissGroup", js.FuncOf(wrapperGroup.DismissGroup))
	js.Global().Set("changeGroupMute", js.FuncOf(wrapperGroup.ChangeGroupMute))
	js.Global().Set("changeGroupMemberMute", js.FuncOf(wrapperGroup.ChangeGroupMemberMute))
	js.Global().Set("setGroupMemberInfo", js.FuncOf(wrapperGroup.SetGroupMemberInfo))
	js.Global().Set("getJoinedGroupList", js.FuncOf(wrapperGroup.GetJoinedGroupList))
	js.Global().Set("getJoinedGroupListPage", js.FuncOf(wrapperGroup.GetJoinedGroupListPage))
	js.Global().Set("searchGroups", js.FuncOf(wrapperGroup.SearchGroups))
	js.Global().Set("setGroupInfo", js.FuncOf(wrapperGroup.SetGroupInfo))
	js.Global().Set("getGroupMemberList", js.FuncOf(wrapperGroup.GetGroupMemberList))
	js.Global().Set("getGroupMemberOwnerAndAdmin", js.FuncOf(wrapperGroup.GetGroupMemberOwnerAndAdmin))
	js.Global().Set("getGroupMemberListByJoinTimeFilter", js.FuncOf(wrapperGroup.GetGroupMemberListByJoinTimeFilter))
	js.Global().Set("getSpecifiedGroupMembersInfo", js.FuncOf(wrapperGroup.GetSpecifiedGroupMembersInfo))
	js.Global().Set("kickGroupMember", js.FuncOf(wrapperGroup.KickGroupMember))
	js.Global().Set("transferGroupOwner", js.FuncOf(wrapperGroup.TransferGroupOwner))
	js.Global().Set("inviteUserToGroup", js.FuncOf(wrapperGroup.InviteUserToGroup))
	js.Global().Set("getGroupApplicationListAsRecipient", js.FuncOf(wrapperGroup.GetGroupApplicationListAsRecipient))
	js.Global().Set("getGroupApplicationListAsApplicant", js.FuncOf(wrapperGroup.GetGroupApplicationListAsApplicant))
	js.Global().Set("acceptGroupApplication", js.FuncOf(wrapperGroup.AcceptGroupApplication))
	js.Global().Set("refuseGroupApplication", js.FuncOf(wrapperGroup.RefuseGroupApplication))
	js.Global().Set("searchGroupMembers", js.FuncOf(wrapperGroup.SearchGroupMembers))
	js.Global().Set("isJoinGroup", js.FuncOf(wrapperGroup.IsJoinGroup))
	js.Global().Set("getUsersInGroup", js.FuncOf(wrapperGroup.GetUsersInGroup))

	wrapperUser := wasm_wrapper.NewWrapperUser(globalFuc)
	js.Global().Set("getSelfUserInfo", js.FuncOf(wrapperUser.GetSelfUserInfo))
	js.Global().Set("setSelfInfo", js.FuncOf(wrapperUser.SetSelfInfo))
	js.Global().Set("getUsersInfo", js.FuncOf(wrapperUser.GetUsersInfo))
	js.Global().Set("subscribeUsersStatus", js.FuncOf(wrapperUser.SubscribeUsersStatus))
	js.Global().Set("unsubscribeUsersStatus", js.FuncOf(wrapperUser.UnsubscribeUsersStatus))
	js.Global().Set("getSubscribeUsersStatus", js.FuncOf(wrapperUser.GetSubscribeUsersStatus))
	js.Global().Set("getUserStatus", js.FuncOf(wrapperUser.GetUserStatus))
	js.Global().Set("addUserCommand", js.FuncOf(wrapperUser.AddUserCommand))
	js.Global().Set("deleteUserCommand", js.FuncOf(wrapperUser.DeleteUserCommand))
	js.Global().Set("getAllUserCommands", js.FuncOf(wrapperUser.GetAllUserCommands))

	wrapperFriend := wasm_wrapper.NewWrapperFriend(globalFuc)
	js.Global().Set("getSpecifiedFriendsInfo", js.FuncOf(wrapperFriend.GetSpecifiedFriendsInfo))
	js.Global().Set("getFriendList", js.FuncOf(wrapperFriend.GetFriendList))
	js.Global().Set("getFriendListPage", js.FuncOf(wrapperFriend.GetFriendListPage))
	js.Global().Set("searchFriends", js.FuncOf(wrapperFriend.SearchFriends))
	js.Global().Set("checkFriend", js.FuncOf(wrapperFriend.CheckFriend))
	js.Global().Set("addFriend", js.FuncOf(wrapperFriend.AddFriend))
	js.Global().Set("updateFriends", js.FuncOf(wrapperFriend.UpdateFriends))
	js.Global().Set("deleteFriend", js.FuncOf(wrapperFriend.DeleteFriend))
	js.Global().Set("getFriendApplicationListAsRecipient", js.FuncOf(wrapperFriend.GetFriendApplicationListAsRecipient))
	js.Global().Set("getFriendApplicationListAsApplicant", js.FuncOf(wrapperFriend.GetFriendApplicationListAsApplicant))
	js.Global().Set("acceptFriendApplication", js.FuncOf(wrapperFriend.AcceptFriendApplication))
	js.Global().Set("refuseFriendApplication", js.FuncOf(wrapperFriend.RefuseFriendApplication))
	js.Global().Set("getBlackList", js.FuncOf(wrapperFriend.GetBlackList))
	js.Global().Set("removeBlack", js.FuncOf(wrapperFriend.RemoveBlack))
	js.Global().Set("addBlack", js.FuncOf(wrapperFriend.AddBlack))

	wrapperThird := wasm_wrapper.NewWrapperThird(globalFuc)
	js.Global().Set("updateFcmToken", js.FuncOf(wrapperThird.UpdateFcmToken))
	js.Global().Set("uploadFile", js.FuncOf(wrapperThird.UploadFile))
}
