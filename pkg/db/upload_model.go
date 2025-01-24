//go:build !js
// +build !js

package db

import (
	"context"
	"time"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"

	"github.com/openimsdk/tools/errs"
)

func (d *DataBase) GetUpload(ctx context.Context, partHash string) (*model_struct.LocalUpload, error) {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	var upload model_struct.LocalUpload
	err := d.conn.WithContext(ctx).Where("part_hash = ?", partHash).Take(&upload).Error
	if err != nil {
		return nil, errs.Wrap(err)
	}
	return &upload, nil
}

func (d *DataBase) InsertUpload(ctx context.Context, upload *model_struct.LocalUpload) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return errs.Wrap(d.conn.WithContext(ctx).Create(upload).Error)
}

func (d *DataBase) deleteUpload(ctx context.Context, partHash string) error {
	return errs.Wrap(d.conn.WithContext(ctx).Where("part_hash = ?", partHash).Delete(&model_struct.LocalUpload{}).Error)
}

func (d *DataBase) UpdateUpload(ctx context.Context, upload *model_struct.LocalUpload) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return errs.Wrap(d.conn.WithContext(ctx).Updates(upload).Error)
}

func (d *DataBase) DeleteUpload(ctx context.Context, partHash string) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return d.deleteUpload(ctx, partHash)
}

func (d *DataBase) DeleteExpireUpload(ctx context.Context) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	var uploads []*model_struct.LocalUpload
	err := d.conn.WithContext(ctx).Where("expire_time <= ?", time.Now().UnixMilli()).Find(&uploads).Error
	if err != nil {
		return errs.Wrap(err)
	}
	for _, upload := range uploads {
		if err := d.deleteUpload(ctx, upload.PartHash); err != nil {
			return err
		}
	}
	return nil
}
