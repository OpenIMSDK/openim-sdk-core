package db

import (
	"context"
	"errors"

	"github.com/openimsdk/openim-sdk-core/v3/wasm/exec"
	"github.com/openimsdk/openim-sdk-core/v3/wasm/indexdb"
)

var ErrType = errors.New("from javascript data type err")

type IndexDB struct {
	*indexdb.LocalUsers
	*indexdb.LocalConversations
	*indexdb.LocalChatLogs
	*indexdb.LocalConversationUnreadMessages
	*indexdb.LocalGroups
	*indexdb.LocalGroupMember
	*indexdb.FriendRequest
	*indexdb.Black
	*indexdb.Friend
	*indexdb.LocalGroupRequest
	*indexdb.LocalChatLogReactionExtensions
	*indexdb.NotificationSeqs
	*indexdb.LocalUpload
	*indexdb.LocalSendingMessages
	*indexdb.LocalUserCommand
	*indexdb.LocalVersionSync
	*indexdb.LocalAppSDKVersion
	*indexdb.LocalTableMaster
	loginUserID string
}

func (i IndexDB) Close(ctx context.Context) error {
	_, err := exec.Exec()
	return err
}

func (i IndexDB) InitDB(ctx context.Context, userID string, dataDir string) error {
	_, err := exec.Exec(userID, dataDir)
	return err
}

func NewDataBase(ctx context.Context, loginUserID string, dbDir string, logLevel int) (*IndexDB, error) {
	i := &IndexDB{
		LocalUsers:                      indexdb.NewLocalUsers(),
		LocalConversations:              indexdb.NewLocalConversations(),
		LocalChatLogs:                   indexdb.NewLocalChatLogs(loginUserID),
		LocalConversationUnreadMessages: indexdb.NewLocalConversationUnreadMessages(),
		LocalGroups:                     indexdb.NewLocalGroups(),
		LocalGroupMember:                indexdb.NewLocalGroupMember(),
		FriendRequest:                   indexdb.NewFriendRequest(loginUserID),
		Black:                           indexdb.NewBlack(loginUserID),
		Friend:                          indexdb.NewFriend(loginUserID),
		LocalGroupRequest:               indexdb.NewLocalGroupRequest(),
		LocalChatLogReactionExtensions:  indexdb.NewLocalChatLogReactionExtensions(),
		NotificationSeqs:                indexdb.NewNotificationSeqs(),
		LocalUpload:                     indexdb.NewLocalUpload(),
		LocalSendingMessages:            indexdb.NewLocalSendingMessages(),
		LocalUserCommand:                indexdb.NewLocalUserCommand(),
		LocalVersionSync:                indexdb.NewLocalVersionSync(),
		LocalAppSDKVersion:              indexdb.NewLocalAppSDKVersion(),
		LocalTableMaster:                indexdb.NewLocalTableMaster(),
		loginUserID:                     loginUserID,
	}
	err := i.InitDB(ctx, loginUserID, dbDir)
	if err != nil {
		return nil, err
	}
	return i, nil
}
