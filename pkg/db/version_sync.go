//go:build !js
// +build !js

package db

import (
	"context"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	"github.com/openimsdk/tools/errs"
	"gorm.io/gorm"
)

func (d *DataBase) GetVersionSync(ctx context.Context, tableName, entityID string) (*model_struct.LocalVersionSync, error) {
	d.versionMtx.RLock()
	defer d.versionMtx.RUnlock()
	var res model_struct.LocalVersionSync
	return &res, errs.Wrap(d.conn.WithContext(ctx).Where("`table_name` = ? and `entity_id` = ?", tableName, entityID).Take(&res).Error)
}

func (d *DataBase) SetVersionSync(ctx context.Context, lv *model_struct.LocalVersionSync) error {
	d.versionMtx.Lock()
	defer d.versionMtx.Unlock()

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
	d.versionMtx.Lock()
	defer d.versionMtx.Unlock()
	localVersionSync := model_struct.LocalVersionSync{Table: tableName, EntityID: entityID}
	return errs.WrapMsg(d.conn.WithContext(ctx).Delete(&localVersionSync).Error, "DeleteVersionSync failed")
}
