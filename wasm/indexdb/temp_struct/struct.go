//go:build js && wasm

package temp_struct

type LocalChatLog struct {
	ServerMsgID          string ` json:"serverMsgID,omitempty"`
	SendID               string ` json:"sendID,omitempty"`
	RecvID               string ` json:"recvID,omitempty"`
	SenderPlatformID     int32  ` json:"senderPlatformID,omitempty"`
	SenderNickname       string ` json:"senderNickname,omitempty"`
	SenderFaceURL        string ` json:"senderFaceURL,omitempty"`
	SessionType          int32  ` json:"sessionType,omitempty"`
	MsgFrom              int32  ` json:"msgFrom,omitempty"`
	ContentType          int32  ` json:"contentType,omitempty"`
	Content              string ` json:"content,omitempty"`
	IsRead               bool   ` json:"isRead,omitempty"`
	Status               int32  ` json:"status,omitempty"`
	Seq                  int64  ` json:"seq,omitempty"`
	SendTime             int64  ` json:"sendTime,omitempty"`
	CreateTime           int64  ` json:"createTime,omitempty"`
	AttachedInfo         string ` json:"attachedInfo,omitempty"`
	Ex                   string ` json:"ex,omitempty"`
	IsReact              bool   ` json:"isReact,omitempty"`
	IsExternalExtensions bool   ` json:"isExternalExtensions,omitempty"`
	MsgFirstModifyTime   int64  ` json:"msgFirstModifyTime,omitempty"`
}

type LocalConversation struct {
	ConversationID        string ` json:"conversationID,omitempty"`
	ConversationType      int32  ` json:"conversationType,omitempty"`
	UserID                string ` json:"userID,omitempty"`
	GroupID               string ` json:"groupID,omitempty"`
	ShowName              string ` json:"showName,omitempty"`
	FaceURL               string ` json:"faceURL,omitempty"`
	RecvMsgOpt            int32  ` json:"recvMsgOpt,omitempty"`
	UnreadCount           int32  ` json:"unreadCount,omitempty"`
	GroupAtType           int32  ` json:"groupAtType,omitempty"`
	LatestMsg             string ` json:"latestMsg,omitempty"`
	LatestMsgSendTime     int64  ` json:"latestMsgSendTime,omitempty"`
	DraftText             string ` json:"draftText,omitempty"`
	DraftTextTime         int64  ` json:"draftTextTime,omitempty"`
	IsPinned              bool   ` json:"isPinned,omitempty"`
	IsPrivateChat         bool   ` json:"isPrivateChat,omitempty"`
	BurnDuration          int32  ` json:"burnDuration,omitempty"`
	IsNotInGroup          bool   ` json:"isNotInGroup,omitempty"`
	UpdateUnreadCountTime int64  ` json:"updateUnreadCountTime,omitempty"`
	AttachedInfo          string ` json:"attachedInfo,omitempty"`
	Ex                    string ` json:"ex,omitempty"`
}

type LocalPartConversation struct {
	RecvMsgOpt            int32  ` json:"recvMsgOpt"`
	GroupAtType           int32  ` json:"groupAtType"`
	IsPinned              bool   ` json:"isPinned,"`
	IsPrivateChat         bool   ` json:"isPrivateChat"`
	IsNotInGroup          bool   ` json:"isNotInGroup"`
	UpdateUnreadCountTime int64  ` json:"updateUnreadCountTime"`
	BurnDuration          int32  ` json:"burnDuration,omitempty"`
	AttachedInfo          string ` json:"attachedInfo"`
	Ex                    string ` json:"ex"`
}

type LocalGroup struct {
	GroupID                string `json:"groupID,omitempty"`
	GroupName              string `json:"groupName,omitempty"`
	Notification           string `json:"notification,omitempty"`
	Introduction           string `json:"introduction,omitempty"`
	FaceURL                string `json:"faceURL,omitempty"`
	CreateTime             uint32 `json:"createTime,omitempty"`
	Status                 int32  `json:"status,omitempty"`
	CreatorUserID          string `json:"creatorUserID,omitempty"`
	GroupType              int32  `json:"groupType,omitempty"`
	OwnerUserID            string `json:"ownerUserID,omitempty"`
	MemberCount            int32  `json:"memberCount,omitempty"`
	Ex                     string `json:"ex,omitempty"`
	AttachedInfo           string `json:"attachedInfo,omitempty"`
	NeedVerification       int32  `json:"needVerification,omitempty"`
	LookMemberInfo         int32  `json:"lookMemberInfo,omitempty"`
	ApplyMemberFriend      int32  `json:"applyMemberFriend,omitempty"`
	NotificationUpdateTime uint32 `json:"notificationUpdateTime,omitempty"`
	NotificationUserID     string `json:"notificationUserID,omitempty"`
}

type LocalFriendRequest struct {
	FromUserID    string `json:"fromUserID,omitempty"`
	FromNickname  string `json:"fromNickname,omitempty"`
	FromFaceURL   string `json:"fromFaceURL,omitempty"`
	ToUserID      string `json:"toUserID,omitempty"`
	ToNickname    string `json:"toNickname,omitempty"`
	ToFaceURL     string `json:"toFaceURL,omitempty"`
	HandleResult  int32  `json:"handleResult,omitempty"`
	ReqMsg        string `json:"reqMsg,omitempty"`
	CreateTime    int64  `json:"createTime,omitempty"`
	HandlerUserID string `json:"handlerUserID,omitempty"`
	HandleMsg     string `json:"handleMsg,omitempty"`
	HandleTime    int64  `json:"handleTime,omitempty"`
	Ex            string `json:"ex,omitempty"`
	AttachedInfo  string `json:"attachedInfo,omitempty"`
}

type LocalBlack struct {
	OwnerUserID    string `json:"ownerUserID,omitempty"`
	BlockUserID    string `json:"blockUserID,omitempty"`
	Nickname       string `json:"nickname,omitempty"`
	FaceURL        string `json:"faceURL,omitempty"`
	CreateTime     int64  `json:"createTime,omitempty"`
	AddSource      int32  `json:"addSource,omitempty"`
	OperatorUserID string `json:"operatorUserID,omitempty"`
	Ex             string `json:"ex,omitempty"`
	AttachedInfo   string `json:"attachedInfo,omitempty"`
}

type LocalFriend struct {
	OwnerUserID    string `json:"ownerUserID,omitempty"`
	FriendUserID   string `json:"friendUserID,omitempty"`
	Remark         string `json:"remark,omitempty"`
	CreateTime     int64  `json:"createTime,omitempty"`
	AddSource      int32  `json:"addSource,omitempty"`
	OperatorUserID string `json:"operatorUserID,omitempty"`
	Nickname       string `json:"nickname,omitempty"`
	FaceURL        string `json:"faceURL,omitempty"`
	Ex             string `json:"ex,omitempty"`
	AttachedInfo   string `json:"attachedInfo,omitempty"`
	IsPinned       bool   `json:"isPinned,omitempty"`
}

type LocalUser struct {
	UserID           string `json:"userID,omitempty"`
	Nickname         string `json:"nickname,omitempty"`
	FaceURL          string `json:"faceURL,omitempty"`
	CreateTime       int64  `json:"createTime,omitempty"`
	AppMangerLevel   int32  `json:"-"`
	Ex               string `json:"ex,omitempty"`
	AttachedInfo     string `json:"attachedInfo,omitempty"`
	GlobalRecvMsgOpt int32  `json:"globalRecvMsgOpt,omitempty"`
}

type LocalUserCommand struct {
	UserID     string `json:"userID,omitempty"`
	Type       int32  `json:"type,omitempty"`
	CreateTime int64  `json:"createTime,omitempty"`
	Uuid       string `json:"uuid,omitempty"`
	Value      string `json:"value,omitempty"`
}
