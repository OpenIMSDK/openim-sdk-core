package friend

import (
	"context"
)

const (
	LocalFriendSyncMaxNum = 1000
)

func (f *Friend) IncrSyncFriends(ctx context.Context) error {
	//opt := incrversion.Option[*model_struct.LocalFriend, *friend.GetIncrementalFriendsResp]{
	//	Ctx: ctx,
	//	DB:  f.db,
	//	Key: func(localFriend *model_struct.LocalFriend) string {
	//		return localFriend.FriendUserID
	//	},
	//	SyncKey: func() string {
	//		return "friend:" + f.loginUserID
	//	},
	//	Local: func() ([]*model_struct.LocalFriend, error) {
	//		return f.db.GetAllFriendList(ctx)
	//	},
	//	Server: func(version *model_struct.LocalVersionSync) (*friend.GetIncrementalFriendsResp, error) {
	//		return util.CallApi[friend.GetIncrementalFriendsResp](ctx, constant.GetIncrementalFriends, &friend.GetIncrementalFriendsReq{
	//			UserID:    f.loginUserID,
	//			Version:   version.Version,
	//			VersionID: version.VersionID,
	//		})
	//	},
	//	Full: func(resp *friend.GetIncrementalFriendsResp) bool {
	//		return resp.Full
	//	},
	//	Version: func(resp *friend.GetIncrementalFriendsResp) (string, uint64) {
	//		return resp.VersionID, resp.Version
	//	},
	//	DeleteIDs: func(resp *friend.GetIncrementalFriendsResp) []string {
	//		return resp.DeleteUserIds
	//	},
	//	Changes: func(resp *friend.GetIncrementalFriendsResp) []*model_struct.LocalFriend {
	//		return util.Batch(ServerFriendToLocalFriend, resp.Changes)
	//	},
	//	Syncer: func(server, local []*model_struct.LocalFriend) error {
	//		return f.friendSyncer.Sync(ctx, server, local, nil)
	//	},
	//}
	//return opt.Sync()
	return nil
}
