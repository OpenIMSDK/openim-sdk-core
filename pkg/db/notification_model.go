//go:build !js
// +build !js

package db

import (
	"context"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	"github.com/openimsdk/tools/errs"
)

func (d *DataBase) SetNotificationSeq(ctx context.Context, conversationID string, seq int64) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	cursor := d.conn.WithContext(ctx).Model(&model_struct.NotificationSeqs{}).Where("conversation_id = ?", conversationID).Updates(map[string]interface{}{"seq": seq})
	if cursor.Error != nil {
		return errs.WrapMsg(cursor.Error, "Updates failed")
	}
	if cursor.RowsAffected == 0 {
		return errs.WrapMsg(d.conn.WithContext(ctx).Create(&model_struct.NotificationSeqs{ConversationID: conversationID, Seq: seq}).Error, "Create failed")
	}
	return nil
}

func (d *DataBase) BatchInsertNotificationSeq(ctx context.Context, notificationSeqs []*model_struct.NotificationSeqs) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return errs.WrapMsg(d.conn.WithContext(ctx).Create(notificationSeqs).Error, "BatchInsertNotificationSeq failed")
}

func (d *DataBase) GetNotificationAllSeqs(ctx context.Context) ([]*model_struct.NotificationSeqs, error) {
	d.mRWMutex.RLock()
	defer d.mRWMutex.RUnlock()
	var seqs []*model_struct.NotificationSeqs
	return seqs, errs.WrapMsg(d.conn.WithContext(ctx).Where("1=1").Find(&seqs).Error, "GetNotificationAllSeqs failed")
}
