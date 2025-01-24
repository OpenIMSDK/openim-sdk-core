package sdk_struct

import (
	"github.com/openimsdk/protocol/sdkws"
)

type MessageReceipt struct {
	GroupID     string   `json:"groupID"`
	UserID      string   `json:"userID"`
	MsgIDList   []string `json:"msgIDList"`
	ReadTime    int64    `json:"readTime"`
	MsgFrom     int32    `json:"msgFrom"`
	ContentType int32    `json:"contentType"`
	SessionType int32    `json:"sessionType"`
}

type MessageRevoked struct {
	RevokerID                   string `json:"revokerID"`
	RevokerRole                 int32  `json:"revokerRole"`
	ClientMsgID                 string `json:"clientMsgID"`
	RevokerNickname             string `json:"revokerNickname"`
	RevokeTime                  int64  `json:"revokeTime"`
	SourceMessageSendTime       int64  `json:"sourceMessageSendTime"`
	SourceMessageSendID         string `json:"sourceMessageSendID"`
	SourceMessageSenderNickname string `json:"sourceMessageSenderNickname"`
	SessionType                 int32  `json:"sessionType"`
	Seq                         int64  `json:"seq"`
	Ex                          string `json:"ex"`
	IsAdminRevoke               bool   `json:"isAdminRevoke"`
}

type MessageReaction struct {
	ClientMsgID  string `json:"clientMsgID"`
	ReactionType int    `json:"reactionType"`
	Counter      int32  `json:"counter,omitempty"`
	UserID       string `json:"userID"`
	GroupID      string `json:"groupID"`
	SessionType  int32  `json:"sessionType"`
	Info         string `json:"info,omitempty"`
}

type ImageInfo struct {
	Width  int32  `json:"x"`
	Height int32  `json:"y"`
	Type   string `json:"type,omitempty"`
	Size   int64  `json:"size"`
}

type PictureBaseInfo struct {
	UUID   string `json:"uuid,omitempty"`
	Type   string `json:"type,omitempty"`
	Size   int64  `json:"size"`
	Width  int32  `json:"width"`
	Height int32  `json:"height"`
	Url    string `json:"url,omitempty"`
}

type SoundBaseInfo struct {
	UUID      string `json:"uuid,omitempty"`
	SoundPath string `json:"soundPath,omitempty"`
	SourceURL string `json:"sourceUrl,omitempty"`
	DataSize  int64  `json:"dataSize"`
	Duration  int64  `json:"duration"`
	SoundType string `json:"soundType,omitempty"`
}

type VideoBaseInfo struct {
	VideoPath      string `json:"videoPath,omitempty"`
	VideoUUID      string `json:"videoUUID,omitempty"`
	VideoURL       string `json:"videoUrl,omitempty"`
	VideoType      string `json:"videoType,omitempty"`
	VideoSize      int64  `json:"videoSize"`
	Duration       int64  `json:"duration"`
	SnapshotPath   string `json:"snapshotPath,omitempty"`
	SnapshotUUID   string `json:"snapshotUUID,omitempty"`
	SnapshotSize   int64  `json:"snapshotSize"`
	SnapshotURL    string `json:"snapshotUrl,omitempty"`
	SnapshotWidth  int32  `json:"snapshotWidth"`
	SnapshotHeight int32  `json:"snapshotHeight"`
	SnapshotType   string `json:"snapshotType,omitempty"`
}

type FileBaseInfo struct {
	FilePath  string `json:"filePath,omitempty"`
	UUID      string `json:"uuid,omitempty"`
	SourceURL string `json:"sourceUrl,omitempty"`
	FileName  string `json:"fileName,omitempty"`
	FileSize  int64  `json:"fileSize"`
	FileType  string `json:"fileType,omitempty"`
}

type TextElem struct {
	Content string `json:"content"`
}

type CardElem struct {
	UserID   string `json:"userID"`
	Nickname string `json:"nickname"`
	FaceURL  string `json:"faceURL"`
	Ex       string `json:"ex"`
}

type PictureElem struct {
	SourcePath      string           `json:"sourcePath,omitempty"`
	SourcePicture   *PictureBaseInfo `json:"sourcePicture,omitempty"`
	BigPicture      *PictureBaseInfo `json:"bigPicture,omitempty"`
	SnapshotPicture *PictureBaseInfo `json:"snapshotPicture,omitempty"`
}

type SoundElem struct {
	UUID      string `json:"uuid,omitempty"`
	SoundPath string `json:"soundPath,omitempty"`
	SourceURL string `json:"sourceUrl,omitempty"`
	DataSize  int64  `json:"dataSize"`
	Duration  int64  `json:"duration"`
	SoundType string `json:"soundType,omitempty"`
}

type VideoElem struct {
	VideoPath      string `json:"videoPath,omitempty"`
	VideoUUID      string `json:"videoUUID,omitempty"`
	VideoURL       string `json:"videoUrl,omitempty"`
	VideoType      string `json:"videoType,omitempty"`
	VideoSize      int64  `json:"videoSize"`
	Duration       int64  `json:"duration"`
	SnapshotPath   string `json:"snapshotPath,omitempty"`
	SnapshotUUID   string `json:"snapshotUUID,omitempty"`
	SnapshotSize   int64  `json:"snapshotSize"`
	SnapshotURL    string `json:"snapshotUrl,omitempty"`
	SnapshotWidth  int32  `json:"snapshotWidth"`
	SnapshotHeight int32  `json:"snapshotHeight"`
	SnapshotType   string `json:"snapshotType,omitempty"`
}

type FileElem struct {
	FilePath  string `json:"filePath,omitempty"`
	UUID      string `json:"uuid,omitempty"`
	SourceURL string `json:"sourceUrl,omitempty"`
	FileName  string `json:"fileName,omitempty"`
	FileSize  int64  `json:"fileSize"`
	FileType  string `json:"fileType,omitempty"`
}

type MergeElem struct {
	Title             string           `json:"title,omitempty"`
	AbstractList      []string         `json:"abstractList,omitempty"`
	MultiMessage      []*MsgStruct     `json:"multiMessage,omitempty"`
	MessageEntityList []*MessageEntity `json:"messageEntityList,omitempty"`
}

type AtTextElem struct {
	Text         string     `json:"text,omitempty"`
	AtUserList   []string   `json:"atUserList,omitempty"`
	AtUsersInfo  []*AtInfo  `json:"atUsersInfo,omitempty"`
	QuoteMessage *MsgStruct `json:"quoteMessage,omitempty"`
	IsAtSelf     bool       `json:"isAtSelf"`
}

type FaceElem struct {
	Index int    `json:"index"`
	Data  string `json:"data,omitempty"`
}

type LocationElem struct {
	Description string  `json:"description,omitempty"`
	Longitude   float64 `json:"longitude"`
	Latitude    float64 `json:"latitude"`
}

type CustomElem struct {
	Data        string `json:"data,omitempty"`
	Description string `json:"description,omitempty"`
	Extension   string `json:"extension,omitempty"`
}

type QuoteElem struct {
	Text              string           `json:"text,omitempty"`
	QuoteMessage      *MsgStruct       `json:"quoteMessage,omitempty"`
	MessageEntityList []*MessageEntity `json:"messageEntityList,omitempty"`
}

type NotificationElem struct {
	Detail string `json:"detail,omitempty"`
}

type AdvancedTextElem struct {
	Text              string           `json:"text,omitempty"`
	MessageEntityList []*MessageEntity `json:"messageEntityList,omitempty"`
}

type TypingElem struct {
	MsgTips string `json:"msgTips,omitempty"`
}

type MsgStruct struct {
	ClientMsgID      string                 `json:"clientMsgID,omitempty"`
	ServerMsgID      string                 `json:"serverMsgID,omitempty"`
	CreateTime       int64                  `json:"createTime"`
	SendTime         int64                  `json:"sendTime"`
	SessionType      int32                  `json:"sessionType"`
	SendID           string                 `json:"sendID,omitempty"`
	RecvID           string                 `json:"recvID,omitempty"`
	MsgFrom          int32                  `json:"msgFrom"`
	ContentType      int32                  `json:"contentType"`
	SenderPlatformID int32                  `json:"senderPlatformID"`
	SenderNickname   string                 `json:"senderNickname,omitempty"`
	SenderFaceURL    string                 `json:"senderFaceUrl,omitempty"`
	GroupID          string                 `json:"groupID,omitempty"`
	Content          string                 `json:"content,omitempty"`
	Seq              int64                  `json:"seq"`
	IsRead           bool                   `json:"isRead"`
	Status           int32                  `json:"status"`
	OfflinePush      *sdkws.OfflinePushInfo `json:"offlinePush,omitempty"`
	AttachedInfo     string                 `json:"attachedInfo,omitempty"`
	Ex               string                 `json:"ex,omitempty"`
	LocalEx          string                 `json:"localEx,omitempty"`
	TextElem         *TextElem              `json:"textElem,omitempty"`
	CardElem         *CardElem              `json:"cardElem,omitempty"`
	PictureElem      *PictureElem           `json:"pictureElem,omitempty"`
	SoundElem        *SoundElem             `json:"soundElem,omitempty"`
	VideoElem        *VideoElem             `json:"videoElem,omitempty"`
	FileElem         *FileElem              `json:"fileElem,omitempty"`
	MergeElem        *MergeElem             `json:"mergeElem,omitempty"`
	AtTextElem       *AtTextElem            `json:"atTextElem,omitempty"`
	FaceElem         *FaceElem              `json:"faceElem,omitempty"`
	LocationElem     *LocationElem          `json:"locationElem,omitempty"`
	CustomElem       *CustomElem            `json:"customElem,omitempty"`
	QuoteElem        *QuoteElem             `json:"quoteElem,omitempty"`
	NotificationElem *NotificationElem      `json:"notificationElem,omitempty"`
	AdvancedTextElem *AdvancedTextElem      `json:"advancedTextElem,omitempty"`
	TypingElem       *TypingElem            `json:"typingElem,omitempty"`
	AttachedInfoElem *AttachedInfoElem      `json:"attachedInfoElem,omitempty"`
}

type AtInfo struct {
	AtUserID      string `json:"atUserID,omitempty"`
	GroupNickname string `json:"groupNickname,omitempty"`
}

type AttachedInfoElem struct {
	GroupHasReadInfo  GroupHasReadInfo `json:"groupHasReadInfo,omitempty"`
	IsPrivateChat     bool             `json:"isPrivateChat"`
	BurnDuration      int32            `json:"burnDuration"`
	HasReadTime       int64            `json:"hasReadTime"`
	MessageEntityList []*MessageEntity `json:"messageEntityList,omitempty"`
	IsEncryption      bool             `json:"isEncryption"`
	InEncryptStatus   bool             `json:"inEncryptStatus"`
	Progress          *UploadProgress  `json:"uploadProgress,omitempty"`
}

type UploadProgress struct {
	Total    int64  `json:"total"`
	Save     int64  `json:"save"`
	Current  int64  `json:"current"`
	UploadID string `json:"uploadID"`
}

type ReactionElem struct {
	Counter          int32               `json:"counter,omitempty"`
	Type             int                 `json:"type,omitempty"`
	UserReactionList []*UserReactionElem `json:"userReactionList,omitempty"`
	CanRepeat        bool                `json:"canRepeat,omitempty"`
	Info             string              `json:"info,omitempty"`
}
type UserReactionElem struct {
	UserID  string `json:"userID,omitempty"`
	Counter int32  `json:"counter,omitempty"`
	Info    string `json:"info,omitempty"`
}
type MessageEntity struct {
	Type   string `json:"type,omitempty"`
	Offset int32  `json:"offset"`
	Length int32  `json:"length"`
	Url    string `json:"url,omitempty"`
	Ex     string `json:"ex,omitempty"`
}
type GroupHasReadInfo struct {
	HasReadUserIDList []string `json:"hasReadUserIDList,omitempty"`
	HasReadCount      int32    `json:"hasReadCount"`
	GroupMemberCount  int32    `json:"groupMemberCount"`
}
type NewMsgList []*MsgStruct

// Implement the sort.Interface interface to get the number of elements method
func (n NewMsgList) Len() int {
	return len(n)
}

// Implement the sort.Interface interface comparison element method
func (n NewMsgList) Less(i, j int) bool {
	return n[i].SendTime < n[j].SendTime
}

// Implement the sort.Interface interface exchange element method
func (n NewMsgList) Swap(i, j int) {
	n[i], n[j] = n[j], n[i]
}

type IMConfig struct {
	SystemType           string `json:"systemType"`
	PlatformID           int32  `json:"platformID"`
	ApiAddr              string `json:"apiAddr"`
	WsAddr               string `json:"wsAddr"`
	DataDir              string `json:"dataDir"`
	LogLevel             uint32 `json:"logLevel"`
	IsLogStandardOutput  bool   `json:"isLogStandardOutput"`
	LogFilePath          string `json:"logFilePath"`
	IsExternalExtensions bool   `json:"isExternalExtensions"`
}

type CmdNewMsgComeToConversation struct {
	Msgs     map[string]*sdkws.PullMsgs
	SyncFlag int
}

type CmdMsgSyncInReinstall struct {
	Msgs  map[string]*sdkws.PullMsgs
	Total int
}

type BasicInfo struct {
	Nickname string
	FaceURL  string
}

type PublicUser struct {
	UserID     string `json:"userID"`
	Nickname   string `json:"nickname"`
	FaceURL    string `json:"faceURL"`
	Ex         string `json:"ex"`
	CreateTime int64  `json:"createTime"`
}
