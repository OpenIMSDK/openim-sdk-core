package open_im_sdk_callback

type Base interface {
	OnError(errCode int32, errMsg string)
	OnSuccess(data string)
}
type SendMsgCallBack interface {
	Base
	OnProgress(progress int)
}

type OnConnListener interface {
	OnConnecting()
	OnConnectSuccess()
	OnConnectFailed(errCode int32, errMsg string)
	OnKickedOffline()
	OnUserTokenExpired()
	OnUserTokenInvalid(errMsg string)
}

type OnGroupListener interface {
	OnJoinedGroupAdded(groupInfo string)
	OnJoinedGroupDeleted(groupInfo string)
	OnGroupMemberAdded(groupMemberInfo string)
	OnGroupMemberDeleted(groupMemberInfo string)
	OnGroupApplicationAdded(groupApplication string)
	OnGroupApplicationDeleted(groupApplication string)
	OnGroupInfoChanged(groupInfo string)
	OnGroupDismissed(groupInfo string)
	OnGroupMemberInfoChanged(groupMemberInfo string)
	OnGroupApplicationAccepted(groupApplication string)
	OnGroupApplicationRejected(groupApplication string)
}
type OnFriendshipListener interface {
	OnFriendApplicationAdded(friendApplication string)
	OnFriendApplicationDeleted(friendApplication string)
	OnFriendApplicationAccepted(friendApplication string)
	OnFriendApplicationRejected(friendApplication string)
	OnFriendAdded(friendInfo string)
	OnFriendDeleted(friendInfo string)
	OnFriendInfoChanged(friendInfo string)
	OnBlackAdded(blackInfo string)
	OnBlackDeleted(blackInfo string)
}
type OnConversationListener interface {
	OnSyncServerStart(reinstalled bool)
	OnSyncServerFinish(reinstalled bool)
	OnSyncServerProgress(progress int)
	OnSyncServerFailed(reinstalled bool)
	OnNewConversation(conversationList string)
	OnConversationChanged(conversationList string)
	OnTotalUnreadMessageCountChanged(totalUnreadCount int32)
	OnConversationUserInputStatusChanged(change string)
}

type OnAdvancedMsgListener interface {
	OnRecvNewMessage(message string)
	OnRecvC2CReadReceipt(msgReceiptList string)
	OnNewRecvMessageRevoked(messageRevoked string)
	OnRecvOfflineNewMessage(message string)
	OnMsgDeleted(message string)
	OnRecvOnlineOnlyMessage(message string)
}

type OnBatchMsgListener interface {
	OnRecvNewMessages(messageList string)
	OnRecvOfflineNewMessages(messageList string)
}

type OnUserListener interface {
	OnSelfInfoUpdated(userInfo string)
	OnUserStatusChanged(userOnlineStatus string)
	OnUserCommandAdd(userCommand string)
	OnUserCommandDelete(userCommand string)
	OnUserCommandUpdate(userCommand string)
}

type OnCustomBusinessListener interface {
	OnRecvCustomBusinessMessage(businessMessage string)
}
type OnMessageKvInfoListener interface {
	OnMessageKvInfoChanged(messageChangedList string)
}

type OnListenerForService interface {
	// OnGroupApplicationAdded Someone applied to join a group
	OnGroupApplicationAdded(groupApplication string)
	// OnGroupApplicationAccepted Group join application has been accepted
	OnGroupApplicationAccepted(groupApplication string)
	// OnFriendApplicationAdded Someone applied to add you as a friend
	OnFriendApplicationAdded(friendApplication string)
	// OnFriendApplicationAccepted Friend request has been accepted
	OnFriendApplicationAccepted(friendApplication string)
	// OnRecvNewMessage Received a new message
	OnRecvNewMessage(message string)
}

type OnSignalingListener interface {
	OnReceiveNewInvitation(receiveNewInvitationCallback string)

	OnInviteeAccepted(inviteeAcceptedCallback string)

	OnInviteeAcceptedByOtherDevice(inviteeAcceptedCallback string)

	OnInviteeRejected(inviteeRejectedCallback string)

	OnInviteeRejectedByOtherDevice(inviteeRejectedCallback string)

	OnInvitationCancelled(invitationCancelledCallback string)

	OnInvitationTimeout(invitationTimeoutCallback string)

	OnHangUp(hangUpCallback string)

	OnRoomParticipantConnected(onRoomParticipantConnectedCallback string)

	OnRoomParticipantDisconnected(onRoomParticipantDisconnectedCallback string)
}

type UploadFileCallback interface {
	// Open a file with a given size
	Open(size int64)
	// PartSize Set the size of each part and the total number of parts
	PartSize(partSize int64, num int)
	// HashPartProgress Progress of hashing each part, including the part index, size, and hash value
	HashPartProgress(index int, size int64, partHash string)
	// HashPartComplete All parts have been hashed, providing the combined hash of all parts and the final file hash
	HashPartComplete(partsHash string, fileHash string)
	// UploadID Upload ID is generated and provided
	UploadID(uploadID string)
	// UploadPartComplete A specific part has completed uploading, providing the part index, size, and hash value
	UploadPartComplete(index int, partSize int64, partHash string)
	// UploadComplete The entire file upload progress, including the file size, stream size, and storage size
	UploadComplete(fileSize int64, streamSize int64, storageSize int64)
	// Complete The file upload is complete, providing the final size, URL, and type of the file
	Complete(size int64, url string, typ int)
}

type UploadLogProgress interface {
	OnProgress(current int64, size int64)
}
