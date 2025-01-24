//go:build !js
// +build !js

package db

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	"github.com/openimsdk/tools/errs"
)

func (d *DataBase) GetLoginUser(ctx context.Context, userID string) (*model_struct.LocalUser, error) {
	d.mRWMutex.RLock()
	defer d.mRWMutex.RUnlock()
	var user model_struct.LocalUser
	err := d.conn.WithContext(ctx).Where("user_id = ? ", userID).Take(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errs.ErrRecordNotFound.Wrap()
		}
		return nil, errs.Wrap(err)
	}
	return &user, nil
}

func (d *DataBase) UpdateLoginUser(ctx context.Context, user *model_struct.LocalUser) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	t := d.conn.WithContext(ctx).Model(user).Select("*").Updates(user)
	if t.RowsAffected == 0 {
		return errs.WrapMsg(errors.New("RowsAffected == 0"), "no update")
	}
	return errs.WrapMsg(t.Error, "UpdateLoginUser failed")
}

func (d *DataBase) UpdateLoginUserByMap(ctx context.Context, user *model_struct.LocalUser, args map[string]interface{}) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	t := d.conn.WithContext(ctx).Model(&user).Updates(args)
	if t.RowsAffected == 0 {
		return errs.WrapMsg(errors.New("RowsAffected == 0"), "no update")
	}
	return errs.WrapMsg(t.Error, "UpdateColumnsConversation failed")
}

func (d *DataBase) InsertLoginUser(ctx context.Context, user *model_struct.LocalUser) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return errs.WrapMsg(d.conn.WithContext(ctx).Create(user).Error, "InsertLoginUser failed")
}
