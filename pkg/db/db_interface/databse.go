package db_interface

import (
	"context"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
)

type GroupModel interface {
	InsertGroup(ctx context.Context, groupInfo *model_struct.LocalGroup) error
	DeleteGroup(ctx context.Context, groupID string) error
	UpdateGroup(ctx context.Context, groupInfo *model_struct.LocalGroup) error
	BatchInsertGroup(ctx context.Context, groupList []*model_struct.LocalGroup) error
	DeleteAllGroup(ctx context.Context) error
	GetJoinedGroupListDB(ctx context.Context) ([]*model_struct.LocalGroup, error)
	GetGroups(ctx context.Context, groupIDs []string) ([]*model_struct.LocalGroup, error)
	GetGroupInfoByGroupID(ctx context.Context, groupID string) (*model_struct.LocalGroup, error)
	GetAllGroupInfoByGroupIDOrGroupName(ctx context.Context, keyword string, isSearchGroupID bool, isSearchGroupName bool) ([]*model_struct.LocalGroup, error)

	InsertAdminGroupRequest(ctx context.Context, groupRequest *model_struct.LocalAdminGroupRequest) error
	DeleteAdminGroupRequest(ctx context.Context, groupID, userID string) error
	UpdateAdminGroupRequest(ctx context.Context, groupRequest *model_struct.LocalAdminGroupRequest) error
	GetAdminGroupApplication(ctx context.Context) ([]*model_struct.LocalAdminGroupRequest, error)
	InsertGroupRequest(ctx context.Context, groupRequest *model_struct.LocalGroupRequest) error
	DeleteGroupRequest(ctx context.Context, groupID, userID string) error
	UpdateGroupRequest(ctx context.Context, groupRequest *model_struct.LocalGroupRequest) error
	GetSendGroupApplication(ctx context.Context) ([]*model_struct.LocalGroupRequest, error)
	GetGroupMemberInfoByGroupIDUserID(ctx context.Context, groupID, userID string) (*model_struct.LocalGroupMember, error)
	GetGroupSomeMemberInfo(ctx context.Context, groupID string, userIDList []string) ([]*model_struct.LocalGroupMember, error)
	GetGroupMemberListByGroupID(ctx context.Context, groupID string) ([]*model_struct.LocalGroupMember, error)
	GetGroupMemberListSplit(ctx context.Context, groupID string, filter int32, offset, count int) ([]*model_struct.LocalGroupMember, error)
	GetGroupMemberListByUserIDs(ctx context.Context, groupID string, filter int32, userIDs []string) ([]*model_struct.LocalGroupMember, error)
	GetGroupMemberOwnerAndAdminDB(ctx context.Context, groupID string) ([]*model_struct.LocalGroupMember, error)
	GetGroupMemberListSplitByJoinTimeFilter(ctx context.Context, groupID string, offset, count int, joinTimeBegin, joinTimeEnd int64, userIDList []string) ([]*model_struct.LocalGroupMember, error)
	InsertGroupMember(ctx context.Context, groupMember *model_struct.LocalGroupMember) error
	BatchInsertGroupMember(ctx context.Context, groupMemberList []*model_struct.LocalGroupMember) error
	DeleteGroupMember(ctx context.Context, groupID, userID string) error
	DeleteGroupAllMembers(ctx context.Context, groupID string) error
	UpdateGroupMember(ctx context.Context, groupMember *model_struct.LocalGroupMember) error
	SearchGroupMembersDB(ctx context.Context, keyword string, groupID string, isSearchMemberNickname, isSearchUserID bool, offset, count int) (result []*model_struct.LocalGroupMember, err error)
}

type MessageModel interface {
	BatchInsertMessageList(ctx context.Context, conversationID string, MessageList []*model_struct.LocalChatLog) error
	InsertMessage(ctx context.Context, conversationID string, Message *model_struct.LocalChatLog) error
	SearchMessageByKeyword(ctx context.Context, contentType []int, keywordList []string, keywordListMatchType int, conversationID string, startTime, endTime int64, offset, count int) (result []*model_struct.LocalChatLog, err error)
	SearchMessageByContentType(ctx context.Context, contentType []int, conversationID string, startTime, endTime int64, offset, count int) (result []*model_struct.LocalChatLog, err error)
	SearchMessageByContentTypeAndKeyword(ctx context.Context, contentType []int, conversationID string, keywordList []string, keywordListMatchType int, startTime, endTime int64) (result []*model_struct.LocalChatLog, err error)
	GetMessage(ctx context.Context, conversationID, clientMsgID string) (*model_struct.LocalChatLog, error)
	GetMessageBySeq(ctx context.Context, conversationID string, seq int64) (*model_struct.LocalChatLog, error)
	UpdateColumnsMessage(ctx context.Context, conversationID string, ClientMsgID string, args map[string]interface{}) error
	UpdateMessage(ctx context.Context, conversationID string, c *model_struct.LocalChatLog) error
	UpdateMessageBySeq(ctx context.Context, conversationID string, c *model_struct.LocalChatLog) error
	UpdateMessageTimeAndStatus(ctx context.Context, conversationID, clientMsgID string, serverMsgID string, sendTime int64, status int32) error
	GetMessageList(ctx context.Context, conversationID string, count int, startTime int64, isReverse bool) (result []*model_struct.LocalChatLog, err error)
	MarkConversationMessageAsReadDB(ctx context.Context, conversationID string, msgIDs []string) (rowsAffected int64, err error)
	MarkConversationMessageAsReadBySeqs(ctx context.Context, conversationID string, seqs []int64) (rowsAffected int64, err error)
	GetUnreadMessage(ctx context.Context, conversationID string) (result []*model_struct.LocalChatLog, err error)
	GetMessagesByClientMsgIDs(ctx context.Context, conversationID string, msgIDs []string) (result []*model_struct.LocalChatLog, err error)
	GetMessagesBySeqs(ctx context.Context, conversationID string, seqs []int64) (result []*model_struct.LocalChatLog, err error)
	GetConversationNormalMsgSeq(ctx context.Context, conversationID string) (int64, error)
	CheckConversationNormalMsgSeq(ctx context.Context, conversationID string) (int64, error)
	GetConversationPeerNormalMsgSeq(ctx context.Context, conversationID string) (int64, error)
	GetLatestActiveMessage(ctx context.Context, conversationID string, isReverse bool) (result []*model_struct.LocalChatLog, err error)
	GetLatestValidServerMessage(ctx context.Context, conversationID string, startTime int64, isReverse bool) (*model_struct.LocalChatLog, error)

	UpdateMsgSenderFaceURLAndSenderNickname(ctx context.Context, conversationID, sendID, faceURL, nickname string) error
	DeleteConversationAllMessages(ctx context.Context, conversationID string) error
	MarkDeleteConversationAllMessages(ctx context.Context, conversationID string) error

	SetNotificationSeq(ctx context.Context, conversationID string, seq int64) error
	BatchInsertNotificationSeq(ctx context.Context, notificationSeqs []*model_struct.NotificationSeqs) error
	GetNotificationAllSeqs(ctx context.Context) ([]*model_struct.NotificationSeqs, error)
}

type ConversationModel interface {
	GetConversationByUserID(ctx context.Context, userID string) (*model_struct.LocalConversation, error)
	GetAllConversationListDB(ctx context.Context) ([]*model_struct.LocalConversation, error)
	FindAllUnreadConversationConversationID(ctx context.Context) ([]string, error)
	GetHiddenConversationList(ctx context.Context) ([]*model_struct.LocalConversation, error)
	GetAllConversations(ctx context.Context) ([]*model_struct.LocalConversation, error)
	GetAllSingleConversationIDList(ctx context.Context) (result []string, err error)
	GetAllConversationIDList(ctx context.Context) (result []string, err error)
	GetConversationListSplitDB(ctx context.Context, offset, count int) ([]*model_struct.LocalConversation, error)
	BatchInsertConversationList(ctx context.Context, conversationList []*model_struct.LocalConversation) error
	InsertConversation(ctx context.Context, conversationList *model_struct.LocalConversation) error
	DeleteConversation(ctx context.Context, conversationID string) error
	DeleteAllConversation(ctx context.Context) error
	GetConversation(ctx context.Context, conversationID string) (*model_struct.LocalConversation, error)
	UpdateConversation(ctx context.Context, c *model_struct.LocalConversation) error
	BatchUpdateConversationList(ctx context.Context, conversationList []*model_struct.LocalConversation) error
	ResetConversation(ctx context.Context, conversationID string) error
	ResetAllConversation(ctx context.Context) error
	ClearConversation(ctx context.Context, conversationID string) error
	SetConversationDraftDB(ctx context.Context, conversationID, draftText string) error
	RemoveConversationDraft(ctx context.Context, conversationID, draftText string) error
	UpdateColumnsConversation(ctx context.Context, conversationID string, args map[string]interface{}) error
	DecrConversationUnreadCount(ctx context.Context, conversationID string, count int64) (err error)
	GetTotalUnreadMsgCountDB(ctx context.Context) (totalUnreadCount int32, err error)
	GetMultipleConversationDB(ctx context.Context, conversationIDList []string) (result []*model_struct.LocalConversation, err error)
	SearchAllMessageByContentType(ctx context.Context, conversationID string, contentType int) ([]*model_struct.LocalChatLog, error)
	SearchConversations(ctx context.Context, searchParam string) ([]*model_struct.LocalConversation, error)
}

type UserModel interface {
	GetLoginUser(ctx context.Context, userID string) (*model_struct.LocalUser, error)
	UpdateLoginUser(ctx context.Context, user *model_struct.LocalUser) error
	InsertLoginUser(ctx context.Context, user *model_struct.LocalUser) error
	ProcessUserCommandAdd(ctx context.Context, command *model_struct.LocalUserCommand) error
	ProcessUserCommandUpdate(ctx context.Context, command *model_struct.LocalUserCommand) error
	ProcessUserCommandDelete(ctx context.Context, command *model_struct.LocalUserCommand) error
	ProcessUserCommandGetAll(ctx context.Context) ([]*model_struct.LocalUserCommand, error)
}

type FriendModel interface {
	InsertFriend(ctx context.Context, friend *model_struct.LocalFriend) error
	DeleteFriendDB(ctx context.Context, friendUserID string) error
	UpdateFriend(ctx context.Context, friend *model_struct.LocalFriend) error
	GetAllFriendList(ctx context.Context) ([]*model_struct.LocalFriend, error)
	BatchInsertFriend(ctx context.Context, friendList []*model_struct.LocalFriend) error
	DeleteAllFriend(ctx context.Context) error

	SearchFriendList(ctx context.Context, keyword string, isSearchUserID, isSearchNickname, isSearchRemark bool) ([]*model_struct.LocalFriend, error)
	GetFriendInfoByFriendUserID(ctx context.Context, FriendUserID string) (*model_struct.LocalFriend, error)
	GetFriendInfoList(ctx context.Context, friendUserIDList []string) ([]*model_struct.LocalFriend, error)
	InsertFriendRequest(ctx context.Context, friendRequest *model_struct.LocalFriendRequest) error
	DeleteFriendRequestBothUserID(ctx context.Context, fromUserID, toUserID string) error
	UpdateFriendRequest(ctx context.Context, friendRequest *model_struct.LocalFriendRequest) error
	GetRecvFriendApplication(ctx context.Context) ([]*model_struct.LocalFriendRequest, error)
	GetSendFriendApplication(ctx context.Context) ([]*model_struct.LocalFriendRequest, error)
	GetBothFriendReq(ctx context.Context, fromUserID, toUserID string) ([]*model_struct.LocalFriendRequest, error)

	GetBlackListDB(ctx context.Context) ([]*model_struct.LocalBlack, error)
	GetBlackInfoList(ctx context.Context, blockUserIDList []string) ([]*model_struct.LocalBlack, error)
	InsertBlack(ctx context.Context, black *model_struct.LocalBlack) error
	UpdateBlack(ctx context.Context, black *model_struct.LocalBlack) error
	DeleteBlack(ctx context.Context, blockUserID string) error
}

type S3Model interface {
	GetUpload(ctx context.Context, partHash string) (*model_struct.LocalUpload, error)
	InsertUpload(ctx context.Context, upload *model_struct.LocalUpload) error
	DeleteUpload(ctx context.Context, partHash string) error
	UpdateUpload(ctx context.Context, upload *model_struct.LocalUpload) error
}

type SendingMessagesModel interface {
	InsertSendingMessage(ctx context.Context, message *model_struct.LocalSendingMessages) error
	DeleteSendingMessage(ctx context.Context, conversationID, clientMsgID string) error
	GetAllSendingMessages(ctx context.Context) (friendRequests []*model_struct.LocalSendingMessages, err error)
}

type VersionSyncModel interface {
	GetVersionSync(ctx context.Context, tableName, entityID string) (*model_struct.LocalVersionSync, error)
	SetVersionSync(ctx context.Context, version *model_struct.LocalVersionSync) error
	DeleteVersionSync(ctx context.Context, tableName, entityID string) error
}

type AppSDKVersion interface {
	GetAppSDKVersion(ctx context.Context) (*model_struct.LocalAppSDKVersion, error)
	SetAppSDKVersion(ctx context.Context, version *model_struct.LocalAppSDKVersion) error
}

type DataBase interface {
	Close(ctx context.Context) error
	GroupModel
	MessageModel
	ConversationModel
	UserModel
	FriendModel
	S3Model
	SendingMessagesModel
	VersionSyncModel
	AppSDKVersion
}
