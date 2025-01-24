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

func (d *DataBase) GetVersionSync(ctx context.Context, tableName, entityID string) (*model_struct.LocalVersionSync, error) {
	d.mRWMutex.RLock()
	defer d.mRWMutex.RUnlock()
	var res model_struct.LocalVersionSync
	err := d.conn.WithContext(ctx).Where("`table_name` = ? and `entity_id` = ?", tableName, entityID).Take(&res).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &model_struct.LocalVersionSync{}, errs.ErrRecordNotFound.Wrap()
		}
		return nil, errs.Wrap(err)
	}
	return &res, errs.Wrap(err)
}

func (d *DataBase) SetVersionSync(ctx context.Context, lv *model_struct.LocalVersionSync) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()

	var existing model_struct.LocalVersionSync
	err := d.conn.WithContext(ctx).Where("`table_name` = ? AND `entity_id` = ?", lv.Table, lv.EntityID).First(&existing).Error

	if err == gorm.ErrRecordNotFound {
		if createErr := d.conn.WithContext(ctx).Create(lv).Error; createErr != nil {
			return errs.Wrap(createErr)
		}
		return nil
	} else if err != nil {
		return errs.Wrap(err)
	}

	if updateErr := d.conn.WithContext(ctx).Model(&existing).Updates(lv).Error; updateErr != nil {
		return errs.Wrap(updateErr)
	}
	return nil
}

func (d *DataBase) DeleteVersionSync(ctx context.Context, tableName, entityID string) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	localVersionSync := model_struct.LocalVersionSync{Table: tableName, EntityID: entityID}
	return errs.WrapMsg(d.conn.WithContext(ctx).Delete(&localVersionSync).Error, "DeleteVersionSync failed")
}
