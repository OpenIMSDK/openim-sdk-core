//go:build js && wasm

package indexdb

import (
	"context"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/utils"
	"github.com/openimsdk/openim-sdk-core/v3/wasm/exec"
)

type LocalSendingMessages struct {
}

func NewLocalSendingMessages() *LocalSendingMessages {
	return &LocalSendingMessages{}
}
func (i *LocalSendingMessages) InsertSendingMessage(ctx context.Context, message *model_struct.LocalSendingMessages) error {
	_, err := exec.Exec(utils.StructToJsonString(message))
	return err
}

func (i *LocalSendingMessages) DeleteSendingMessage(ctx context.Context, conversationID, clientMsgID string) error {
	_, err := exec.Exec(conversationID, clientMsgID)
	return err
}
func (i *LocalSendingMessages) GetAllSendingMessages(ctx context.Context) (result []*model_struct.LocalSendingMessages, err error) {
	gList, err := exec.Exec()
	if err != nil {
		return nil, err
	} else {
		if v, ok := gList.(string); ok {
			var temp []model_struct.LocalSendingMessages
			err := utils.JsonStringToStruct(v, &temp)
			if err != nil {
				return nil, err
			}
			for _, v := range temp {
				v1 := v
				result = append(result, &v1)
			}
			return result, err
		} else {
			return nil, exec.ErrType
		}
	}
}
