//go:build !js
// +build !js

package db

import (
	"context"
	"errors"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	"github.com/openimsdk/tools/errs"
)

func (d *DataBase) InsertGroupRequest(ctx context.Context, groupRequest *model_struct.LocalGroupRequest) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return errs.WrapMsg(d.conn.WithContext(ctx).Create(groupRequest).Error, "InsertGroupRequest failed")
}

func (d *DataBase) DeleteGroupRequest(ctx context.Context, groupID, userID string) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return errs.WrapMsg(d.conn.WithContext(ctx).Where("group_id=? and user_id=?", groupID, userID).Delete(&model_struct.LocalGroupRequest{}).Error, "DeleteGroupRequest failed")
}

func (d *DataBase) UpdateGroupRequest(ctx context.Context, groupRequest *model_struct.LocalGroupRequest) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	t := d.conn.WithContext(ctx).Model(groupRequest).Select("*").Updates(*groupRequest)
	if t.RowsAffected == 0 {
		return errs.WrapMsg(errors.New("RowsAffected == 0"), "no update")
	}
	return errs.Wrap(t.Error)
}

func (d *DataBase) GetSendGroupApplication(ctx context.Context) ([]*model_struct.LocalGroupRequest, error) {
	d.mRWMutex.RLock()
	defer d.mRWMutex.RUnlock()
	var groupRequestList []*model_struct.LocalGroupRequest
	return groupRequestList, errs.Wrap(d.conn.WithContext(ctx).Order("create_time DESC").Find(&groupRequestList).Error)
}
