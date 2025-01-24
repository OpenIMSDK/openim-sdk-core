//go:build !js
// +build !js

package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	"github.com/openimsdk/tools/errs"

	"gorm.io/gorm"
)

func (d *DataBase) InsertGroup(ctx context.Context, groupInfo *model_struct.LocalGroup) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return errs.WrapMsg(d.conn.WithContext(ctx).Create(groupInfo).Error, "InsertGroup failed")
}

func (d *DataBase) DeleteGroup(ctx context.Context, groupID string) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	localGroup := model_struct.LocalGroup{GroupID: groupID}
	return errs.WrapMsg(d.conn.WithContext(ctx).Delete(&localGroup).Error, "DeleteGroup failed")
}

func (d *DataBase) UpdateGroup(ctx context.Context, groupInfo *model_struct.LocalGroup) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()

	t := d.conn.WithContext(ctx).Model(groupInfo).Select("*").Updates(*groupInfo)
	if t.RowsAffected == 0 {
		return errs.WrapMsg(errors.New("RowsAffected == 0"), "no update")
	}
	return errs.Wrap(t.Error)
}

func (d *DataBase) BatchInsertGroup(ctx context.Context, groupList []*model_struct.LocalGroup) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return errs.WrapMsg(d.conn.WithContext(ctx).Create(groupList).Error, "BatchInsertGroup failed")
}

func (d *DataBase) DeleteAllGroup(ctx context.Context) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return errs.WrapMsg(d.conn.WithContext(ctx).Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&model_struct.LocalGroup{}).Error, "DeleteAllGroup failed")
}

func (d *DataBase) GetJoinedGroupListDB(ctx context.Context) ([]*model_struct.LocalGroup, error) {
	d.mRWMutex.RLock()
	defer d.mRWMutex.RUnlock()
	var groupList []*model_struct.LocalGroup
	err := d.conn.WithContext(ctx).Find(&groupList).Error
	return groupList, errs.WrapMsg(err, "GetJoinedGroupList failed ")
}

func (d *DataBase) GetGroups(ctx context.Context, groupIDs []string) ([]*model_struct.LocalGroup, error) {
	d.mRWMutex.RLock()
	defer d.mRWMutex.RUnlock()
	var groupList []*model_struct.LocalGroup
	err := d.conn.WithContext(ctx).Where("group_id in (?)", groupIDs).Find(&groupList).Error
	return groupList, errs.WrapMsg(err, "GetGroups failed ")
}

func (d *DataBase) GetGroupInfoByGroupID(ctx context.Context, groupID string) (*model_struct.LocalGroup, error) {
	d.mRWMutex.RLock()
	defer d.mRWMutex.RUnlock()
	var g model_struct.LocalGroup
	return &g, errs.WrapMsg(d.conn.WithContext(ctx).Where("group_id = ?", groupID).Take(&g).Error, "GetGroupList failed")
}

func (d *DataBase) GetAllGroupInfoByGroupIDOrGroupName(ctx context.Context, keyword string, isSearchGroupID bool, isSearchGroupName bool) ([]*model_struct.LocalGroup, error) {
	d.mRWMutex.RLock()
	defer d.mRWMutex.RUnlock()

	var groupList []*model_struct.LocalGroup
	var condition string
	if isSearchGroupID {
		if isSearchGroupName {
			condition = fmt.Sprintf("group_id like %q or name like %q", "%"+keyword+"%", "%"+keyword+"%")
		} else {
			condition = fmt.Sprintf("group_id like %q ", "%"+keyword+"%")
		}
	} else {
		condition = fmt.Sprintf("name like %q ", "%"+keyword+"%")
	}
	err := d.conn.WithContext(ctx).Where(condition).Order("create_time DESC").Find(&groupList).Error
	return groupList, errs.WrapMsg(err, "GetAllGroupInfoByGroupIDOrGroupName failed ")
}
