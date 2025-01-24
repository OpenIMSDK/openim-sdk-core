//go:build js && wasm

package indexdb

import (
	"context"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/utils"
	"github.com/openimsdk/openim-sdk-core/v3/wasm/exec"
)

type LocalChatLogReactionExtensions struct {
	ExtKey  string `json:"ext_key"`
	ExtVal  string `json:"ext_val"`
	ExtKey2 string `json:"ext_key2"`
	ExtVal2 string `json:"ext_val2"`
}

func NewLocalChatLogReactionExtensions() *LocalChatLogReactionExtensions {
	return &LocalChatLogReactionExtensions{}
}

func (i *LocalChatLogReactionExtensions) GetMessageReactionExtension(ctx context.Context, clientMsgID string) (result *model_struct.LocalChatLogReactionExtensions, err error) {
	msg, err := exec.Exec(clientMsgID)
	if err != nil {
		return nil, err
	} else {
		if v, ok := msg.(string); ok {
			result := model_struct.LocalChatLogReactionExtensions{}
			err := utils.JsonStringToStruct(v, &result)
			if err != nil {
				return nil, err
			}
			return &result, err
		} else {
			return nil, exec.ErrType
		}
	}
}

func (i *LocalChatLogReactionExtensions) InsertMessageReactionExtension(ctx context.Context, messageReactionExtension *model_struct.LocalChatLogReactionExtensions) error {
	_, err := exec.Exec(utils.StructToJsonString(messageReactionExtension))
	return err
}

//	func (i *LocalChatLogReactionExtensions) GetAndUpdateMessageReactionExtension(ctx context.Context, clientMsgID string, m map[string]*sdkws.KeyValue) error {
//		_, err := exec.Exec(clientMsgID, utils.StructToJsonString(m))
//		return err
//	}
//
//	func (i *LocalChatLogReactionExtensions) DeleteAndUpdateMessageReactionExtension(ctx context.Context, clientMsgID string, m map[string]*sdkws.KeyValue) error {
//		_, err := exec.Exec(clientMsgID, utils.StructToJsonString(m))
//		return err
//	}
func (i *LocalChatLogReactionExtensions) GetMultipleMessageReactionExtension(ctx context.Context, msgIDList []string) (result []*model_struct.LocalChatLogReactionExtensions, err error) {
	msgReactionExtensionList, err := exec.Exec(utils.StructToJsonString(msgIDList))
	if err != nil {
		return nil, err
	} else {
		if v, ok := msgReactionExtensionList.(string); ok {
			var temp []model_struct.LocalChatLogReactionExtensions
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
func (i *LocalChatLogReactionExtensions) DeleteMessageReactionExtension(ctx context.Context, msgID string) error {
	_, err := exec.Exec(msgID)
	return err
}
func (i *LocalChatLogReactionExtensions) UpdateMessageReactionExtension(ctx context.Context, c *model_struct.LocalChatLogReactionExtensions) error {
	_, err := exec.Exec(c.ClientMsgID, utils.StructToJsonString(c))
	return err
}
