//go:build !js
// +build !js

package db

import (
	"context"
	"errors"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	"github.com/openimsdk/tools/errs"
)

func (d *DataBase) InsertFriendRequest(ctx context.Context, friendRequest *model_struct.LocalFriendRequest) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return errs.WrapMsg(d.conn.WithContext(ctx).Create(friendRequest).Error, "InsertFriendRequest failed")
}

func (d *DataBase) DeleteFriendRequestBothUserID(ctx context.Context, fromUserID, toUserID string) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return errs.WrapMsg(d.conn.WithContext(ctx).Where("from_user_id=? and to_user_id=?", fromUserID, toUserID).Delete(&model_struct.LocalFriendRequest{}).Error, "DeleteFriendRequestBothUserID failed")
}

func (d *DataBase) UpdateFriendRequest(ctx context.Context, friendRequest *model_struct.LocalFriendRequest) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	t := d.conn.WithContext(ctx).Model(friendRequest).Select("*").Updates(*friendRequest)
	if t.RowsAffected == 0 {
		return errs.WrapMsg(errors.New("RowsAffected == 0"), "no update")
	}
	return errs.Wrap(t.Error)
}

func (d *DataBase) GetRecvFriendApplication(ctx context.Context) ([]*model_struct.LocalFriendRequest, error) {
	d.mRWMutex.RLock()
	defer d.mRWMutex.RUnlock()
	var friendRequestList []*model_struct.LocalFriendRequest

	return friendRequestList, errs.WrapMsg(d.conn.WithContext(ctx).Where("to_user_id = ?", d.loginUserID).Order("create_time DESC").Find(&friendRequestList).Error, "GetRecvFriendApplication failed")
}

func (d *DataBase) GetSendFriendApplication(ctx context.Context) ([]*model_struct.LocalFriendRequest, error) {
	d.mRWMutex.RLock()
	defer d.mRWMutex.RUnlock()
	var friendRequestList []*model_struct.LocalFriendRequest
	return friendRequestList, errs.WrapMsg(d.conn.WithContext(ctx).Where("from_user_id = ?", d.loginUserID).Order("create_time DESC").Find(&friendRequestList).Error, "GetSendFriendApplication failed")
}

func (d *DataBase) GetBothFriendReq(ctx context.Context, fromUserID, toUserID string) (friendRequests []*model_struct.LocalFriendRequest, err error) {
	d.mRWMutex.RLock()
	defer d.mRWMutex.RUnlock()
	return friendRequests, errs.WrapMsg(d.conn.WithContext(ctx).Where("(from_user_id = ? AND to_user_id = ?) OR (from_user_id = ? AND to_user_id = ?)", fromUserID, toUserID, toUserID, fromUserID).Find(&friendRequests).Error, "GetFriendApplicationByBothID failed")
}
