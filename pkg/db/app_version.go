//go:build !js
// +build !js

package db

import (
	"context"

	"gorm.io/gorm"

	"github.com/openimsdk/tools/errs"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
)

func (d *DataBase) GetAppSDKVersion(ctx context.Context) (*model_struct.LocalAppSDKVersion, error) {
	var appVersion model_struct.LocalAppSDKVersion
	err := d.conn.WithContext(ctx).Take(&appVersion).Error
	if err == gorm.ErrRecordNotFound {
		err = errs.ErrRecordNotFound
	}
	return &appVersion, errs.Wrap(err)
}

func (d *DataBase) SetAppSDKVersion(ctx context.Context, appVersion *model_struct.LocalAppSDKVersion) error {
	var exist model_struct.LocalAppSDKVersion
	err := d.conn.WithContext(ctx).First(&exist).Error
	if err == gorm.ErrRecordNotFound {
		if createErr := d.conn.WithContext(ctx).Create(appVersion).Error; createErr != nil {
			return errs.Wrap(createErr)
		}
		return nil
	} else if err != nil {
		return errs.Wrap(err)
	}

	if updateErr := d.conn.WithContext(ctx).Model(&exist).Updates(appVersion).Error; updateErr != nil {
		return errs.Wrap(updateErr)
	}

	return nil
}
