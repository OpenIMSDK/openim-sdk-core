package server_api_params

type Conversation struct {
	OwnerUserID           string `json:"ownerUserID" binding:"required"`
	ConversationID        string `json:"conversationID" binding:"required"`
	ConversationType      int32  `json:"conversationType" binding:"required"`
	UserID                string `json:"userID"`
	GroupID               string `json:"groupID"`
	RecvMsgOpt            int32  `json:"recvMsgOpt"`
	UnreadCount           int32  `json:"unreadCount"`
	DraftTextTime         int64  `json:"draftTextTime"`
	IsPinned              bool   `json:"isPinned"`
	IsPrivateChat         bool   `json:"isPrivateChat"`
	BurnDuration          int32  `json:"burnDuration"`
	GroupAtType           int32  `json:"groupAtType"`
	IsNotInGroup          bool   `json:"isNotInGroup"`
	UpdateUnreadCountTime int64  ` json:"updateUnreadCountTime"`
	AttachedInfo          string `json:"attachedInfo"`
	Ex                    string `json:"ex"`
}
