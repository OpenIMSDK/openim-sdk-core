//go:build !js
// +build !js

package db

import (
	"context"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	"github.com/openimsdk/tools/errs"
)

func (d *DataBase) InsertSendingMessage(ctx context.Context, message *model_struct.LocalSendingMessages) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return errs.WrapMsg(d.conn.WithContext(ctx).Create(message).Error, "InsertSendingMessage failed")
}

func (d *DataBase) DeleteSendingMessage(ctx context.Context, conversationID, clientMsgID string) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	localSendingMessage := model_struct.LocalSendingMessages{ConversationID: conversationID, ClientMsgID: clientMsgID}
	return errs.WrapMsg(d.conn.WithContext(ctx).Delete(&localSendingMessage).Error, "DeleteSendingMessage failed")
}

func (d *DataBase) GetAllSendingMessages(ctx context.Context) (friendRequests []*model_struct.LocalSendingMessages, err error) {
	d.mRWMutex.RLock()
	defer d.mRWMutex.RUnlock()
	return friendRequests, errs.WrapMsg(d.conn.WithContext(ctx).Find(&friendRequests).Error, "GetAllSendingMessages failed")
}
