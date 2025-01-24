package conversation_msg

import (
	"context"
	"encoding/json"

	"github.com/openimsdk/openim-sdk-core/v3/internal/third/file"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/db_interface"
	"github.com/openimsdk/openim-sdk-core/v3/sdk_struct"

	"github.com/openimsdk/tools/log"
)

func NewUploadFileCallback(ctx context.Context, progress func(progress int), msg *sdk_struct.MsgStruct, conversationID string, db db_interface.DataBase) file.UploadFileCallback {
	if msg.AttachedInfoElem == nil {
		msg.AttachedInfoElem = &sdk_struct.AttachedInfoElem{}
	}
	if msg.AttachedInfoElem.Progress == nil {
		msg.AttachedInfoElem.Progress = &sdk_struct.UploadProgress{}
	}
	return &msgUploadFileCallback{ctx: ctx, progress: progress, msg: msg, db: db, conversationID: conversationID}
}

type msgUploadFileCallback struct {
	ctx            context.Context
	db             db_interface.DataBase
	msg            *sdk_struct.MsgStruct
	conversationID string
	value          int
	progress       func(progress int)
}

func (c *msgUploadFileCallback) Open(size int64) {
}

func (c *msgUploadFileCallback) PartSize(partSize int64, num int) {
}

func (c *msgUploadFileCallback) HashPartProgress(index int, size int64, partHash string) {
}

func (c *msgUploadFileCallback) HashPartComplete(partsHash string, fileHash string) {
}

func (c *msgUploadFileCallback) UploadID(uploadID string) {
	c.msg.AttachedInfoElem.Progress.UploadID = uploadID
	data, err := json.Marshal(c.msg.AttachedInfoElem)
	if err != nil {
		panic(err)
	}
	if err := c.db.UpdateColumnsMessage(c.ctx, c.conversationID, c.msg.ClientMsgID, map[string]any{"attached_info": string(data)}); err != nil {
		log.ZError(c.ctx, "update PutProgress message attached info failed", err)
	}
}

func (c *msgUploadFileCallback) UploadPartComplete(index int, partSize int64, partHash string) {
}

func (c *msgUploadFileCallback) UploadComplete(fileSize int64, streamSize int64, storageSize int64) {
	c.msg.AttachedInfoElem.Progress.Save = storageSize
	c.msg.AttachedInfoElem.Progress.Current = streamSize
	c.msg.AttachedInfoElem.Progress.Total = fileSize
	data, err := json.Marshal(c.msg.AttachedInfoElem)
	if err != nil {
		panic(err)
	}
	if err := c.db.UpdateColumnsMessage(c.ctx, c.conversationID, c.msg.ClientMsgID, map[string]any{"attached_info": string(data)}); err != nil {
		log.ZError(c.ctx, "update PutProgress message attached info failed", err)
	}
	value := int(float64(streamSize) / float64(fileSize) * 100)
	if c.value < value {
		c.value = value
		c.progress(value)
	}
}

func (c *msgUploadFileCallback) Complete(size int64, url string, typ int) {
	if c.value != 100 {
		c.progress(100)
	}
	c.msg.AttachedInfoElem.Progress = nil
	data, err := json.Marshal(c.msg.AttachedInfoElem)
	if err != nil {
		panic(err)
	}
	if err := c.db.UpdateColumnsMessage(c.ctx, c.conversationID, c.msg.ClientMsgID, map[string]any{"attached_info": string(data)}); err != nil {
		log.ZError(c.ctx, "update PutComplete message attached info failed", err)
	}
}
