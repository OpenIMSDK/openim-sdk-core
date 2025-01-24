//go:build js && wasm

package indexdb

import (
	"context"

	"github.com/openimsdk/openim-sdk-core/v3/wasm/exec"
)

import (
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/utils"
)

type LocalConversationUnreadMessages struct {
}

func NewLocalConversationUnreadMessages() *LocalConversationUnreadMessages {
	return &LocalConversationUnreadMessages{}
}

func (i *LocalConversationUnreadMessages) BatchInsertConversationUnreadMessageList(ctx context.Context, messageList []*model_struct.LocalConversationUnreadMessage) error {
	if messageList == nil {
		return nil
	}
	_, err := exec.Exec(utils.StructToJsonString(messageList))
	return err
}

func (i *LocalConversationUnreadMessages) DeleteConversationUnreadMessageList(ctx context.Context, conversationID string, sendTime int64) int64 {
	deleteRows, err := exec.Exec(conversationID, sendTime)
	if err != nil {
		return 0
	} else {
		if v, ok := deleteRows.(float64); ok {
			var result int64
			result = int64(v)
			return result
		} else {
			return 0
		}
	}
}
