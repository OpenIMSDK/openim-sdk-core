package user

import (
	"context"
	"fmt"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/sdkerrs"
	"github.com/openimsdk/tools/log"

	"github.com/openimsdk/openim-sdk-core/v3/open_im_sdk_callback"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/cache"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/common"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/db_interface"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/syncer"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/utils"
	"github.com/openimsdk/tools/utils/datautil"
)

// NewUser creates a new User object.
func NewUser(dataBase db_interface.DataBase, loginUserID string, conversationCh chan common.Cmd2Value) *User {
	user := &User{DataBase: dataBase, loginUserID: loginUserID, conversationCh: conversationCh}
	user.initSyncer()
	//user.OnlineStatusCache = cache.NewCache[string, *userPb.OnlineStatus]()
	user.UserCache = cache.NewUserCache[string, *model_struct.LocalUser](
		func(value *model_struct.LocalUser) string { return value.UserID },
		nil,
		user.GetLoginUser,
		user.GetUsersInfoFromServer,
	)
	return user
}

// User is a struct that represents a user in the system.
type User struct {
	db_interface.DataBase
	loginUserID    string
	listener       func() open_im_sdk_callback.OnUserListener
	userSyncer     *syncer.Syncer[*model_struct.LocalUser, syncer.NoResp, string]
	commandSyncer  *syncer.Syncer[*model_struct.LocalUserCommand, syncer.NoResp, string]
	conversationCh chan common.Cmd2Value
	UserCache      *cache.UserCache[string, *model_struct.LocalUser]

	//OnlineStatusCache *cache.Cache[string, *userPb.OnlineStatus]
}

// SetListener sets the user's listener.
func (u *User) SetListener(listener func() open_im_sdk_callback.OnUserListener) {
	u.listener = listener
}

func (u *User) initSyncer() {
	u.userSyncer = syncer.New[*model_struct.LocalUser, syncer.NoResp, string](
		func(ctx context.Context, value *model_struct.LocalUser) error {
			return u.InsertLoginUser(ctx, value)
		},
		func(ctx context.Context, value *model_struct.LocalUser) error {
			return fmt.Errorf("not support delete user %s", value.UserID)
		},
		func(ctx context.Context, serverUser, localUser *model_struct.LocalUser) error {
			u.UserCache.Delete(localUser.UserID)
			return u.DataBase.UpdateLoginUser(context.Background(), serverUser)
		},
		func(user *model_struct.LocalUser) string {
			return user.UserID
		},
		nil,
		func(ctx context.Context, state int, server, local *model_struct.LocalUser) error {
			switch state {
			case syncer.Update:
				u.listener().OnSelfInfoUpdated(utils.StructToJsonString(server))
				if server.Nickname != local.Nickname || server.FaceURL != local.FaceURL {
					_ = common.TriggerCmdUpdateMessage(ctx, common.UpdateMessageNode{Action: constant.UpdateMsgFaceUrlAndNickName,
						Args: common.UpdateMessageInfo{SessionType: constant.SingleChatType, UserID: server.UserID, FaceURL: server.FaceURL, Nickname: server.Nickname}}, u.conversationCh)
				}
			}
			return nil
		},
	)
	u.commandSyncer = syncer.New[*model_struct.LocalUserCommand, syncer.NoResp, string](
		func(ctx context.Context, command *model_struct.LocalUserCommand) error {
			// Logic to insert a command
			return u.DataBase.ProcessUserCommandAdd(ctx, command)
		},
		func(ctx context.Context, command *model_struct.LocalUserCommand) error {
			// Logic to delete a command
			return u.DataBase.ProcessUserCommandDelete(ctx, command)
		},
		func(ctx context.Context, serverCommand *model_struct.LocalUserCommand, localCommand *model_struct.LocalUserCommand) error {
			// Logic to update a command
			if serverCommand == nil || localCommand == nil {
				return fmt.Errorf("nil command reference")
			}
			return u.DataBase.ProcessUserCommandUpdate(ctx, serverCommand)
		},
		func(command *model_struct.LocalUserCommand) string {
			// Return a unique identifier for the command
			if command == nil {
				return ""
			}
			return command.Uuid
		},
		func(a *model_struct.LocalUserCommand, b *model_struct.LocalUserCommand) bool {
			// Compare two commands to check if they are equal
			if a == nil || b == nil {
				return false
			}
			return a.Uuid == b.Uuid && a.Type == b.Type && a.Value == b.Value
		},
		func(ctx context.Context, state int, serverCommand *model_struct.LocalUserCommand, localCommand *model_struct.LocalUserCommand) error {
			if u.listener == nil {
				return nil
			}
			switch state {
			case syncer.Delete:
				u.listener().OnUserCommandDelete(utils.StructToJsonString(serverCommand))
			case syncer.Update:
				u.listener().OnUserCommandUpdate(utils.StructToJsonString(serverCommand))
			case syncer.Insert:
				u.listener().OnUserCommandAdd(utils.StructToJsonString(serverCommand))
			}
			return nil
		},
	)
}

func (u *User) GetUserInfoWithCache(ctx context.Context, cacheKey string) (*model_struct.LocalUser, error) {
	return u.UserCache.Fetch(ctx, cacheKey)
}

func (u *User) GetUsersInfoWithCache(ctx context.Context, cacheKeys []string) ([]*model_struct.LocalUser, error) {
	m, err := u.UserCache.BatchFetch(ctx, cacheKeys)
	if err != nil {
		return nil, err
	}
	return datautil.Values(m), nil
}

// GetSingleUserFromServer retrieves user information from the server.
func (u *User) GetSingleUserFromServer(ctx context.Context, userID string) (*model_struct.LocalUser, error) {
	users, err := u.getUsersInfo(ctx, []string{userID})
	if err != nil {
		return nil, err
	}
	if len(users) > 0 {
		return ServerUserToLocalUser(users[0]), nil
	}
	return nil, sdkerrs.ErrUserIDNotFound.WrapMsg(fmt.Sprintf("getSelfUserInfo failed, userID: %s not exist", userID))
}

// GetUsersInfoFromServer retrieves user information from the server.
func (u *User) GetUsersInfoFromServer(ctx context.Context, userIDs []string) ([]*model_struct.LocalUser, error) {
	var err error

	serverUsersInfo, err := u.getUsersInfo(ctx, userIDs)
	if err != nil {
		return nil, err
	}

	if len(serverUsersInfo) == 0 {
		log.ZError(ctx, "serverUsersInfo is empty", err, "userIDs", userIDs)
		return nil, err
	}

	return datautil.Batch(ServerUserToLocalUser, serverUsersInfo), nil
}
