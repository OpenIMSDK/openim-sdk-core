package model_struct

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/openimsdk/tools/errs"
)

type LocalFriend struct {
	OwnerUserID    string `gorm:"column:owner_user_id;primary_key;type:varchar(64)" json:"ownerUserID"`
	FriendUserID   string `gorm:"column:friend_user_id;primary_key;type:varchar(64)" json:"userID"`
	Remark         string `gorm:"column:remark;type:varchar(255)" json:"remark"`
	CreateTime     int64  `gorm:"column:create_time" json:"createTime"`
	AddSource      int32  `gorm:"column:add_source" json:"addSource"`
	OperatorUserID string `gorm:"column:operator_user_id;type:varchar(64)" json:"operatorUserID"`
	Nickname       string `gorm:"column:name;type:varchar;type:varchar(255)" json:"nickname"`
	FaceURL        string `gorm:"column:face_url;type:varchar;type:varchar(255)" json:"faceURL"`
	Ex             string `gorm:"column:ex;type:varchar(1024)" json:"ex"`
	AttachedInfo   string `gorm:"column:attached_info;type:varchar(1024)" json:"attachedInfo"`
	IsPinned       bool   `gorm:"column:is_pinned;" json:"isPinned"`
}

func (LocalFriend) TableName() string {
	return "local_friends"
}

type LocalFriendRequest struct {
	FromUserID    string `gorm:"column:from_user_id;primary_key;type:varchar(64)" json:"fromUserID"`
	FromNickname  string `gorm:"column:from_nickname;type:varchar;type:varchar(255)" json:"fromNickname"`
	FromFaceURL   string `gorm:"column:from_face_url;type:varchar;type:varchar(255)" json:"fromFaceURL"`
	ToUserID      string `gorm:"column:to_user_id;primary_key;type:varchar(64)" json:"toUserID"`
	ToNickname    string `gorm:"column:to_nickname;type:varchar;type:varchar(255)" json:"toNickname"`
	ToFaceURL     string `gorm:"column:to_face_url;type:varchar;type:varchar(255)" json:"toFaceURL"`
	HandleResult  int32  `gorm:"column:handle_result" json:"handleResult"`
	ReqMsg        string `gorm:"column:req_msg;type:varchar(255)" json:"reqMsg"`
	CreateTime    int64  `gorm:"column:create_time" json:"createTime"`
	HandlerUserID string `gorm:"column:handler_user_id;type:varchar(64)" json:"handlerUserID"`
	HandleMsg     string `gorm:"column:handle_msg;type:varchar(255)" json:"handleMsg"`
	HandleTime    int64  `gorm:"column:handle_time" json:"handleTime"`
	Ex            string `gorm:"column:ex;type:varchar(1024)" json:"ex"`
	AttachedInfo  string `gorm:"column:attached_info;type:varchar(1024)" json:"attachedInfo"`
}

type LocalGroup struct {
	GroupID                string `gorm:"column:group_id;primary_key;type:varchar(64)" json:"groupID" binding:"required"`
	GroupName              string `gorm:"column:name;size:255" json:"groupName"`
	Notification           string `gorm:"column:notification;type:varchar(255)" json:"notification"`
	Introduction           string `gorm:"column:introduction;type:varchar(255)" json:"introduction"`
	FaceURL                string `gorm:"column:face_url;type:varchar(255)" json:"faceURL"`
	CreateTime             int64  `gorm:"column:create_time" json:"createTime"`
	Status                 int32  `gorm:"column:status" json:"status"`
	CreatorUserID          string `gorm:"column:creator_user_id;type:varchar(64)" json:"creatorUserID"`
	GroupType              int32  `gorm:"column:group_type" json:"groupType"`
	OwnerUserID            string `gorm:"column:owner_user_id;type:varchar(64)" json:"ownerUserID"`
	MemberCount            int32  `gorm:"column:member_count" json:"memberCount"`
	Ex                     string `gorm:"column:ex;type:varchar(1024)" json:"ex"`
	AttachedInfo           string `gorm:"column:attached_info;type:varchar(1024)" json:"attachedInfo"`
	NeedVerification       int32  `gorm:"column:need_verification"  json:"needVerification"`
	LookMemberInfo         int32  `gorm:"column:look_member_info" json:"lookMemberInfo"`
	ApplyMemberFriend      int32  `gorm:"column:apply_member_friend" json:"applyMemberFriend"`
	NotificationUpdateTime int64  `gorm:"column:notification_update_time" json:"notificationUpdateTime"`
	NotificationUserID     string `gorm:"column:notification_user_id;size:64" json:"notificationUserID"`
}

func (LocalGroup) TableName() string {
	return "local_groups"
}

type LocalGroupMember struct {
	GroupID        string `gorm:"column:group_id;primary_key;type:varchar(64)" json:"groupID"`
	UserID         string `gorm:"column:user_id;primary_key;type:varchar(64)" json:"userID"`
	Nickname       string `gorm:"column:nickname;type:varchar(255)" json:"nickname"`
	FaceURL        string `gorm:"column:user_group_face_url;type:varchar(255)" json:"faceURL"`
	RoleLevel      int32  `gorm:"column:role_level;index:index_role_level;" json:"roleLevel"`
	JoinTime       int64  `gorm:"column:join_time;index:index_join_time;" json:"joinTime"`
	JoinSource     int32  `gorm:"column:join_source" json:"joinSource"`
	InviterUserID  string `gorm:"column:inviter_user_id;size:64"  json:"inviterUserID"`
	MuteEndTime    int64  `gorm:"column:mute_end_time;default:0" json:"muteEndTime"`
	OperatorUserID string `gorm:"column:operator_user_id;type:varchar(64)" json:"operatorUserID"`
	Ex             string `gorm:"column:ex;type:varchar(1024)" json:"ex"`
	AttachedInfo   string `gorm:"column:attached_info;type:varchar(1024)" json:"attachedInfo"`
}

func (LocalGroupMember) TableName() string {
	return "local_group_members"
}

type LocalGroupRequest struct {
	GroupID       string `gorm:"column:group_id;primary_key;type:varchar(64)" json:"groupID"`
	GroupName     string `gorm:"column:group_name;size:255" json:"groupName"`
	Notification  string `gorm:"column:notification;type:varchar(255)" json:"notification"`
	Introduction  string `gorm:"column:introduction;type:varchar(255)" json:"introduction"`
	GroupFaceURL  string `gorm:"column:face_url;type:varchar(255)" json:"groupFaceURL"`
	CreateTime    int64  `gorm:"column:create_time" json:"createTime"`
	Status        int32  `gorm:"column:status" json:"status"`
	CreatorUserID string `gorm:"column:creator_user_id;type:varchar(64)" json:"creatorUserID"`
	GroupType     int32  `gorm:"column:group_type" json:"groupType"`
	OwnerUserID   string `gorm:"column:owner_user_id;type:varchar(64)" json:"ownerUserID"`
	MemberCount   int32  `gorm:"column:member_count" json:"memberCount"`
	UserID        string `gorm:"column:user_id;primary_key;type:varchar(64)" json:"userID"`
	Nickname      string `gorm:"column:nickname;type:varchar(255)" json:"nickname"`
	UserFaceURL   string `gorm:"column:user_face_url;type:varchar(255)" json:"userFaceURL"`
	HandleResult  int32  `gorm:"column:handle_result" json:"handleResult"`
	ReqMsg        string `gorm:"column:req_msg;type:varchar(255)" json:"reqMsg"`
	HandledMsg    string `gorm:"column:handle_msg;type:varchar(255)" json:"handledMsg"`
	ReqTime       int64  `gorm:"column:req_time" json:"reqTime"`
	HandleUserID  string `gorm:"column:handle_user_id;type:varchar(64)" json:"handleUserID"`
	HandledTime   int64  `gorm:"column:handle_time" json:"handledTime"`
	Ex            string `gorm:"column:ex;type:varchar(1024)" json:"ex"`
	AttachedInfo  string `gorm:"column:attached_info;type:varchar(1024)" json:"attachedInfo"`
	JoinSource    int32  `gorm:"column:join_source" json:"joinSource"`
	InviterUserID string `gorm:"column:inviter_user_id;size:64"  json:"inviterUserID"`
}

type LocalUser struct {
	UserID           string `gorm:"column:user_id;primary_key;type:varchar(64)" json:"userID"`
	Nickname         string `gorm:"column:name;type:varchar(255)" json:"nickname"`
	FaceURL          string `gorm:"column:face_url;type:varchar(255)" json:"faceURL"`
	CreateTime       int64  `gorm:"column:create_time" json:"createTime"`
	AppMangerLevel   int32  `gorm:"column:app_manger_level" json:"-"`
	Ex               string `gorm:"column:ex;type:varchar(1024)" json:"ex"`
	AttachedInfo     string `gorm:"column:attached_info;type:varchar(1024)" json:"attachedInfo"`
	GlobalRecvMsgOpt int32  `gorm:"column:global_recv_msg_opt" json:"globalRecvMsgOpt"`
}

type LocalBlack struct {
	OwnerUserID    string `gorm:"column:owner_user_id;primary_key;type:varchar(64)" json:"ownerUserID"`
	BlockUserID    string `gorm:"column:block_user_id;primary_key;type:varchar(64)" json:"userID"`
	Nickname       string `gorm:"column:nickname;type:varchar(255)" json:"nickname"`
	FaceURL        string `gorm:"column:face_url;type:varchar(255)" json:"faceURL"`
	CreateTime     int64  `gorm:"column:create_time" json:"createTime"`
	AddSource      int32  `gorm:"column:add_source" json:"addSource"`
	OperatorUserID string `gorm:"column:operator_user_id;type:varchar(64)" json:"operatorUserID"`
	Ex             string `gorm:"column:ex;type:varchar(1024)" json:"ex"`
	AttachedInfo   string `gorm:"column:attached_info;type:varchar(1024)" json:"attachedInfo"`
}

type LocalSeqData struct {
	UserID string `gorm:"column:user_id;primary_key;type:varchar(64)"`
	Seq    uint32 `gorm:"column:seq"`
}

type LocalSeq struct {
	ID     string `gorm:"column:id;primary_key;type:varchar(64)"`
	MinSeq uint32 `gorm:"column:min_seq"`
}

type LocalChatLog struct {
	ClientMsgID      string `gorm:"column:client_msg_id;primary_key;type:char(64)" json:"clientMsgID"`
	ServerMsgID      string `gorm:"column:server_msg_id;type:char(64)" json:"serverMsgID"`
	SendID           string `gorm:"column:send_id;type:char(64)" json:"sendID"`
	RecvID           string `gorm:"column:recv_id;index:index_recv_id;type:char(64)" json:"recvID"`
	SenderPlatformID int32  `gorm:"column:sender_platform_id" json:"senderPlatformID"`
	SenderNickname   string `gorm:"column:sender_nick_name;type:varchar(255)" json:"senderNickname"`
	SenderFaceURL    string `gorm:"column:sender_face_url;type:varchar(255)" json:"senderFaceURL"`
	SessionType      int32  `gorm:"column:session_type" json:"sessionType"`
	MsgFrom          int32  `gorm:"column:msg_from" json:"msgFrom"`
	ContentType      int32  `gorm:"column:content_type;index:content_type_alone" json:"contentType"`
	Content          string `gorm:"column:content;type:varchar(1000)" json:"content"`
	IsRead           bool   `gorm:"column:is_read" json:"isRead"`
	Status           int32  `gorm:"column:status" json:"status"`
	Seq              int64  `gorm:"column:seq;index:index_seq;default:0" json:"seq"`
	SendTime         int64  `gorm:"column:send_time;index:index_send_time;" json:"sendTime"`
	CreateTime       int64  `gorm:"column:create_time" json:"createTime"`
	AttachedInfo     string `gorm:"column:attached_info;type:varchar(1024)" json:"attachedInfo"`
	Ex               string `gorm:"column:ex;type:varchar(1024)" json:"ex"`
	LocalEx          string `gorm:"column:local_ex;type:varchar(1024)" json:"localEx"`
}

type LocalErrChatLog struct {
	Seq              int64  `gorm:"column:seq;primary_key" json:"seq"`
	ClientMsgID      string `gorm:"column:client_msg_id;type:char(64)" json:"clientMsgID"`
	ServerMsgID      string `gorm:"column:server_msg_id;type:char(64)" json:"serverMsgID"`
	SendID           string `gorm:"column:send_id;type:char(64)" json:"sendID"`
	RecvID           string `gorm:"column:recv_id;type:char(64)" json:"recvID"`
	SenderPlatformID int32  `gorm:"column:sender_platform_id" json:"senderPlatformID"`
	SenderNickname   string `gorm:"column:sender_nick_name;type:varchar(255)" json:"senderNickname"`
	SenderFaceURL    string `gorm:"column:sender_face_url;type:varchar(255)" json:"senderFaceURL"`
	SessionType      int32  `gorm:"column:session_type" json:"sessionType"`
	MsgFrom          int32  `gorm:"column:msg_from" json:"msgFrom"`
	ContentType      int32  `gorm:"column:content_type" json:"contentType"`
	Content          string `gorm:"column:content;type:varchar(1000)" json:"content"`
	IsRead           bool   `gorm:"column:is_read" json:"isRead"`
	Status           int32  `gorm:"column:status" json:"status"`
	SendTime         int64  `gorm:"column:send_time" json:"sendTime"`
	CreateTime       int64  `gorm:"column:create_time" json:"createTime"`
	AttachedInfo     string `gorm:"column:attached_info;type:varchar(1024)" json:"attachedInfo"`
	Ex               string `gorm:"column:ex;type:varchar(1024)" json:"ex"`
}
type TempCacheLocalChatLog struct {
	ClientMsgID      string `gorm:"column:client_msg_id;primary_key;type:char(64)" json:"clientMsgID"`
	ServerMsgID      string `gorm:"column:server_msg_id;type:char(64)" json:"serverMsgID"`
	SendID           string `gorm:"column:send_id;type:char(64)" json:"sendID"`
	RecvID           string `gorm:"column:recv_id;type:char(64)" json:"recvID"`
	SenderPlatformID int32  `gorm:"column:sender_platform_id" json:"senderPlatformID"`
	SenderNickname   string `gorm:"column:sender_nick_name;type:varchar(255)" json:"senderNickname"`
	SenderFaceURL    string `gorm:"column:sender_face_url;type:varchar(255)" json:"senderFaceURL"`
	SessionType      int32  `gorm:"column:session_type" json:"sessionType"`
	MsgFrom          int32  `gorm:"column:msg_from" json:"msgFrom"`
	ContentType      int32  `gorm:"column:content_type" json:"contentType"`
	Content          string `gorm:"column:content;type:varchar(1000)" json:"content"`
	IsRead           bool   `gorm:"column:is_read" json:"isRead"`
	Status           int32  `gorm:"column:status" json:"status"`
	Seq              int64  `gorm:"column:seq;default:0" json:"seq"`
	SendTime         int64  `gorm:"column:send_time;" json:"sendTime"`
	CreateTime       int64  `gorm:"column:create_time" json:"createTime"`
	AttachedInfo     string `gorm:"column:attached_info;type:varchar(1024)" json:"attachedInfo"`
	Ex               string `gorm:"column:ex;type:varchar(1024)" json:"ex"`
}

type LocalConversation struct {
	ConversationID        string `gorm:"column:conversation_id;primary_key;type:char(128)" json:"conversationID"`
	ConversationType      int32  `gorm:"column:conversation_type" json:"conversationType"`
	UserID                string `gorm:"column:user_id;type:char(64)" json:"userID"`
	GroupID               string `gorm:"column:group_id;type:char(128)" json:"groupID"`
	ShowName              string `gorm:"column:show_name;type:varchar(255)" json:"showName"`
	FaceURL               string `gorm:"column:face_url;type:varchar(255)" json:"faceURL"`
	RecvMsgOpt            int32  `gorm:"column:recv_msg_opt" json:"recvMsgOpt"`
	UnreadCount           int32  `gorm:"column:unread_count" json:"unreadCount"`
	GroupAtType           int32  `gorm:"column:group_at_type" json:"groupAtType"`
	LatestMsg             string `gorm:"column:latest_msg;type:varchar(1000)" json:"latestMsg"`
	LatestMsgSendTime     int64  `gorm:"column:latest_msg_send_time;index:index_latest_msg_send_time" json:"latestMsgSendTime"`
	DraftText             string `gorm:"column:draft_text" json:"draftText"`
	DraftTextTime         int64  `gorm:"column:draft_text_time" json:"draftTextTime"`
	IsPinned              bool   `gorm:"column:is_pinned" json:"isPinned"`
	IsPrivateChat         bool   `gorm:"column:is_private_chat" json:"isPrivateChat"`
	BurnDuration          int32  `gorm:"column:burn_duration;default:30" json:"burnDuration"`
	IsNotInGroup          bool   `gorm:"column:is_not_in_group" json:"isNotInGroup"`
	UpdateUnreadCountTime int64  `gorm:"column:update_unread_count_time" json:"updateUnreadCountTime"`
	AttachedInfo          string `gorm:"column:attached_info;type:varchar(1024)" json:"attachedInfo"`
	Ex                    string `gorm:"column:ex;type:varchar(1024)" json:"ex"`
	MaxSeq                int64  `gorm:"column:max_seq" json:"maxSeq"`
	MinSeq                int64  `gorm:"column:min_seq" json:"minSeq"`
	MsgDestructTime       int64  `gorm:"column:msg_destruct_time;default:604800" json:"msgDestructTime"`
	IsMsgDestruct         bool   `gorm:"column:is_msg_destruct;default:false" json:"isMsgDestruct"`
}

func (LocalConversation) TableName() string {
	return "local_conversations"
}

type LocalConversationUnreadMessage struct {
	ConversationID string `gorm:"column:conversation_id;primary_key;type:char(128)" json:"conversationID"`
	ClientMsgID    string `gorm:"column:client_msg_id;primary_key;type:char(64)" json:"clientMsgID"`
	SendTime       int64  `gorm:"column:send_time" json:"sendTime"`
	Ex             string `gorm:"column:ex;type:varchar(1024)" json:"ex"`
}

type LocalAdminGroupRequest struct {
	LocalGroupRequest
}

type LocalChatLogReactionExtensions struct {
	ClientMsgID             string `gorm:"column:client_msg_id;primary_key;type:char(64)" json:"clientMsgID"`
	LocalReactionExtensions []byte `gorm:"column:local_reaction_extensions" json:"localReactionExtensions"`
}

type NotificationSeqs struct {
	ConversationID string `gorm:"column:conversation_id;primary_key;type:char(128)" json:"conversationID"`
	Seq            int64  `gorm:"column:seq" json:"seq"`
}

func (NotificationSeqs) TableName() string {
	return "local_notification_seqs"
}

type LocalUpload struct {
	PartHash   string `gorm:"column:part_hash;primary_key" json:"partHash"`
	UploadID   string `gorm:"column:upload_id;type:varchar(1000)" json:"uploadID"`
	UploadInfo string `gorm:"column:upload_info;type:varchar(2000)" json:"uploadInfo"`
	ExpireTime int64  `gorm:"column:expire_time" json:"expireTime"`
	CreateTime int64  `gorm:"column:create_time" json:"createTime"`
}

func (LocalUpload) TableName() string {
	return "local_uploads"
}

type LocalStranger struct {
	UserID           string `gorm:"column:user_id;primary_key;type:varchar(64)" json:"userID"`
	Nickname         string `gorm:"column:name;type:varchar(255)" json:"nickname"`
	FaceURL          string `gorm:"column:face_url;type:varchar(255)" json:"faceURL"`
	CreateTime       int64  `gorm:"column:create_time" json:"createTime"`
	AppMangerLevel   int32  `gorm:"column:app_manger_level" json:"-"`
	Ex               string `gorm:"column:ex;type:varchar(1024)" json:"ex"`
	AttachedInfo     string `gorm:"column:attached_info;type:varchar(1024)" json:"attachedInfo"`
	GlobalRecvMsgOpt int32  `gorm:"column:global_recv_msg_opt" json:"globalRecvMsgOpt"`
}

func (LocalStranger) TableName() string {
	return "local_stranger"
}

type LocalSendingMessages struct {
	ConversationID string `gorm:"column:conversation_id;primary_key;type:char(128)" json:"conversationID"`
	ClientMsgID    string `gorm:"column:client_msg_id;primary_key;type:char(64)" json:"clientMsgID"`
	Ex             string `gorm:"column:ex;type:varchar(1024)" json:"ex"`
}

func (LocalSendingMessages) TableName() string {
	return "local_sending_messages"
}

type LocalUserCommand struct {
	UserID     string `gorm:"column:user_id;type:char(128);primary_key" json:"userID"`
	Type       int32  `gorm:"column:type;primary_key" json:"type"`
	Uuid       string `gorm:"column:uuid;type:varchar(255);primary_key" json:"uuid"`
	CreateTime int64  `gorm:"column:create_time" json:"createTime"`
	Value      string `gorm:"column:value;type:varchar(255)" json:"value"`
	Ex         string `gorm:"column:ex;type:varchar(1024)" json:"ex"`
}

func (LocalUserCommand) TableName() string {
	return "local_user_command"
}

type StringArray []string

func (a StringArray) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *StringArray) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errs.New("type assertion to []byte failed").Wrap()
	}
	return json.Unmarshal(b, &a)
}

type LocalVersionSync struct {
	Table      string      `gorm:"column:table_name;type:varchar(255);primary_key" json:"tableName"`
	EntityID   string      `gorm:"column:entity_id;type:varchar(255);primary_key" json:"entityID"`
	VersionID  string      `gorm:"column:version_id" json:"versionID"`
	Version    uint64      `gorm:"column:version" json:"version"`
	CreateTime int64       `gorm:"column:create_time" json:"createTime"`
	UIDList    StringArray `gorm:"column:id_list;type:text" json:"uidList"`
}

func (LocalVersionSync) TableName() string {
	return "local_sync_version"
}

type LocalAppSDKVersion struct {
	Version   string `gorm:"column:version;type:varchar(255);primary_key" json:"version"`
	Installed bool   `gorm:"column:installed" json:"installed"` // Mark whether it has already been loaded
}

func (LocalAppSDKVersion) TableName() string {
	return "local_app_sdk_version"
}
