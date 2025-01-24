//go:build js && wasm

package indexdb

import (
	"context"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/utils"
	"github.com/openimsdk/openim-sdk-core/v3/wasm/exec"
)

type NotificationSeqs struct {
}

func NewNotificationSeqs() *NotificationSeqs {
	return &NotificationSeqs{}
}

func (i *NotificationSeqs) SetNotificationSeq(ctx context.Context, conversationID string, seq int64) error {
	_, err := exec.Exec(conversationID, seq)
	return err
}

func (i *NotificationSeqs) BatchInsertNotificationSeq(ctx context.Context, notificationSeqs []*model_struct.NotificationSeqs) error {
	_, err := exec.Exec(utils.StructToJsonString(notificationSeqs))
	return err
}

func (i *NotificationSeqs) GetNotificationAllSeqs(ctx context.Context) (result []*model_struct.NotificationSeqs, err error) {
	gList, err := exec.Exec()
	if err != nil {
		return nil, err
	} else {
		if v, ok := gList.(string); ok {
			err := utils.JsonStringToStruct(v, &result)
			if err != nil {
				return nil, err
			}
			return result, err
		} else {
			return nil, exec.ErrType
		}
	}
}
