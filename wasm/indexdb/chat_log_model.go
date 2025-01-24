//go:build js && wasm

package indexdb

import (
	"context"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/utils"
	"github.com/openimsdk/openim-sdk-core/v3/wasm/exec"
	"github.com/openimsdk/openim-sdk-core/v3/wasm/indexdb/temp_struct"
	"github.com/openimsdk/tools/errs"
)

type LocalChatLogs struct {
	loginUserID string
}

// NewLocalChatLogs creates a new LocalChatLogs
func NewLocalChatLogs(loginUserID string) *LocalChatLogs {
	return &LocalChatLogs{loginUserID: loginUserID}
}

// GetMessage gets the message from the database
func (i *LocalChatLogs) GetMessage(ctx context.Context, conversationID, clientMsgID string) (*model_struct.LocalChatLog, error) {
	msg, err := exec.Exec(conversationID, clientMsgID)
	if err != nil {
		return nil, err
	} else {
		if v, ok := msg.(string); ok {
			result := model_struct.LocalChatLog{}
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

// UpdateMessage updates the message in the database
func (i *LocalChatLogs) UpdateMessage(ctx context.Context, conversationID string, c *model_struct.LocalChatLog) error {
	if c.ClientMsgID == "" {
		return exec.PrimaryKeyNull
	}
	tempLocalChatLog := temp_struct.LocalChatLog{
		ServerMsgID:      c.ServerMsgID,
		SendID:           c.SendID,
		RecvID:           c.RecvID,
		SenderPlatformID: c.SenderPlatformID,
		SenderNickname:   c.SenderNickname,
		SenderFaceURL:    c.SenderFaceURL,
		SessionType:      c.SessionType,
		MsgFrom:          c.MsgFrom,
		ContentType:      c.ContentType,
		Content:          c.Content,
		IsRead:           c.IsRead,
		Status:           c.Status,
		Seq:              c.Seq,
		SendTime:         c.SendTime,
		CreateTime:       c.CreateTime,
		AttachedInfo:     c.AttachedInfo,
		Ex:               c.Ex,
	}
	_, err := exec.Exec(conversationID, c.ClientMsgID, utils.StructToJsonString(tempLocalChatLog))
	return err
}

// UpdateMessageStatus updates the message status in the database
func (i *LocalChatLogs) BatchInsertMessageList(ctx context.Context, conversationID string, messageList []*model_struct.LocalChatLog) error {
	_, err := exec.Exec(conversationID, utils.StructToJsonString(messageList))
	return err
}

// InsertMessage inserts a message into the local chat log.
func (i *LocalChatLogs) InsertMessage(ctx context.Context, conversationID string, message *model_struct.LocalChatLog) error {
	_, err := exec.Exec(conversationID, utils.StructToJsonString(message))
	return err
}

// UpdateColumnsMessageList updates multiple columns of a message in the local chat log.
func (i *LocalChatLogs) UpdateColumnsMessageList(ctx context.Context, clientMsgIDList []string, args map[string]interface{}) error {
	_, err := exec.Exec(utils.StructToJsonString(clientMsgIDList), args)
	return err
}

// UpdateColumnsMessage updates a column of a message in the local chat log.
func (i *LocalChatLogs) UpdateColumnsMessage(ctx context.Context, conversationID, clientMsgID string, args map[string]interface{}) error {
	_, err := exec.Exec(conversationID, clientMsgID, utils.StructToJsonString(args))
	return err
}

// DeleteAllMessage deletes all messages from the local chat log.
func (i *LocalChatLogs) DeleteAllMessage(ctx context.Context) error {
	_, err := exec.Exec()
	return err
}

// UpdateMessageStatusBySourceID updates the status of a message in the local chat log by its source ID.
func (i *LocalChatLogs) UpdateMessageStatusBySourceID(ctx context.Context, sourceID string, status, sessionType int32) error {
	_, err := exec.Exec(sourceID, status, sessionType, i.loginUserID)
	return err
}

// UpdateMessageTimeAndStatus updates the time and status of a message in the local chat log.
func (i *LocalChatLogs) UpdateMessageTimeAndStatus(ctx context.Context, conversationID, clientMsgID string, serverMsgID string, sendTime int64, status int32) error {
	_, err := exec.Exec(conversationID, clientMsgID, serverMsgID, sendTime, status)
	return err
}

// GetMessageList retrieves a list of messages from the local chat log.
func (i *LocalChatLogs) GetMessageList(ctx context.Context, conversationID string, count int, startTime int64, isReverse bool) (result []*model_struct.LocalChatLog, err error) {
	msgList, err := exec.Exec(conversationID, count, startTime, isReverse, i.loginUserID)
	if err != nil {
		return nil, err
	} else {
		if v, ok := msgList.(string); ok {
			var temp []model_struct.LocalChatLog
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

func (i *LocalChatLogs) GetLatestActiveMessage(ctx context.Context, conversationID string, isReverse bool) (result []*model_struct.LocalChatLog, err error) {
	msg, err := exec.Exec(conversationID, isReverse)
	if err != nil {
		return nil, err
	} else {
		if v, ok := msg.(string); ok {
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

func (i *LocalChatLogs) GetLatestValidServerMessage(ctx context.Context, conversationID string, startTime int64, isReverse bool) (*model_struct.LocalChatLog, error) {
	msg, err := exec.Exec(conversationID, startTime, isReverse)
	if err != nil {
		if errs.ErrRecordNotFound.Is(err) {
			return nil, nil
		}
		return nil, err
	} else {
		if v, ok := msg.(string); ok {
			result := model_struct.LocalChatLog{}
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

// UpdateSingleMessageHasRead updates the hasRead field of a single message in the local chat log.
func (i *LocalChatLogs) UpdateSingleMessageHasRead(ctx context.Context, sendID string, msgIDList []string) error {
	_, err := exec.Exec(sendID, utils.StructToJsonString(msgIDList))
	return err
}

// SearchMessageByContentType searches for messages in the local chat log by content type.
func (i *LocalChatLogs) SearchMessageByContentType(ctx context.Context, contentType []int, conversationID string, startTime, endTime int64, offset, count int) (messages []*model_struct.LocalChatLog, err error) {
	msgList, err := exec.Exec(conversationID, utils.StructToJsonString(contentType), startTime, endTime, offset, count)
	if err != nil {
		return nil, err
	} else {
		if v, ok := msgList.(string); ok {
			var temp []model_struct.LocalChatLog
			err := utils.JsonStringToStruct(v, &temp)
			if err != nil {
				return nil, err
			}
			for _, v := range temp {
				v1 := v
				messages = append(messages, &v1)
			}
			return messages, err
		} else {
			return nil, exec.ErrType
		}
	}
}

// SearchMessageByKeyword searches for messages in the local chat log by keyword.
func (i *LocalChatLogs) SearchMessageByContentTypeAndKeyword(ctx context.Context, contentType []int, conversationID string, keywordList []string, keywordListMatchType int, startTime, endTime int64) (result []*model_struct.LocalChatLog, err error) {
	msgList, err := exec.Exec(conversationID, utils.StructToJsonString(contentType), utils.StructToJsonString(keywordList), keywordListMatchType, startTime, endTime)
	if err != nil {
		return nil, err
	} else {
		if v, ok := msgList.(string); ok {
			var temp []model_struct.LocalChatLog
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

// MessageIfExists check if message exists
func (i *LocalChatLogs) MessageIfExists(ctx context.Context, clientMsgID string) (bool, error) {
	isExist, err := exec.Exec(clientMsgID)
	if err != nil {
		return false, err
	} else {
		if v, ok := isExist.(bool); ok {
			return v, nil
		} else {
			return false, exec.ErrType
		}
	}
}

// IsExistsInErrChatLogBySeq check if message exists in error chat log by seq
func (i *LocalChatLogs) IsExistsInErrChatLogBySeq(ctx context.Context, seq int64) bool {
	isExist, err := exec.Exec(seq)
	if err != nil {
		return false
	} else {
		if v, ok := isExist.(bool); ok {
			return v
		} else {
			return false
		}
	}
}

// MessageIfExistsBySeq check if message exists by seq
func (i *LocalChatLogs) MessageIfExistsBySeq(ctx context.Context, seq int64) (bool, error) {
	isExist, err := exec.Exec(seq)
	if err != nil {
		return false, err
	} else {
		if v, ok := isExist.(bool); ok {
			return v, nil
		} else {
			return false, exec.ErrType
		}
	}
}

// GetMultipleMessage gets multiple messages from the local chat log.
func (i *LocalChatLogs) GetMultipleMessage(ctx context.Context, msgIDList []string) (result []*model_struct.LocalChatLog, err error) {
	msgList, err := exec.Exec(utils.StructToJsonString(msgIDList))
	if err != nil {
		return nil, err
	} else {
		if v, ok := msgList.(string); ok {
			var temp []model_struct.LocalChatLog
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

// GetLostMsgSeqList gets lost message seq list.
func (i *LocalChatLogs) GetLostMsgSeqList(ctx context.Context, minSeqInSvr uint32) (result []uint32, err error) {
	l, err := exec.Exec(minSeqInSvr)
	if err != nil {
		return nil, err
	} else {
		if v, ok := l.(string); ok {
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

// GetTestMessage gets test message.
func (i *LocalChatLogs) GetTestMessage(ctx context.Context, seq uint32) (*model_struct.LocalChatLog, error) {
	msg, err := exec.Exec(seq)
	if err != nil {
		return nil, err
	} else {
		if v, ok := msg.(model_struct.LocalChatLog); ok {
			return &v, err
		} else {
			return nil, exec.ErrType
		}
	}
}

// Update the sender's face URL and nickname in the chat logs
func (i *LocalChatLogs) UpdateMsgSenderFaceURLAndSenderNickname(ctx context.Context, conversationID, sendID, faceURL, nickname string) error {
	_, err := exec.Exec(conversationID, sendID, faceURL, nickname)
	return err
}

// Get the message sequence number by client message ID
func (i *LocalChatLogs) GetMsgSeqByClientMsgID(ctx context.Context, clientMsgID string) (uint32, error) {
	result, err := exec.Exec(clientMsgID)
	if err != nil {
		return 0, err
	}
	if v, ok := result.(float64); ok {
		return uint32(v), nil
	}
	return 0, exec.ErrType
}

// Search all messages by content type
func (i *LocalChatLogs) SearchAllMessageByContentType(ctx context.Context, conversationID string, contentType int) (result []*model_struct.LocalChatLog, err error) {
	msgList, err := exec.Exec(conversationID, contentType)
	if err != nil {
		return nil, err
	} else {
		if v, ok := msgList.(string); ok {
			var temp []*model_struct.LocalChatLog
			err := utils.JsonStringToStruct(v, &temp)
			if err != nil {
				return nil, err
			}
			for _, v := range temp {
				v1 := v
				result = append(result, v1)
			}
			return result, err
		} else {
			return nil, exec.ErrType
		}
	}
}

// Get the message sequence number list by group ID
func (i *LocalChatLogs) GetMsgSeqListByGroupID(ctx context.Context, groupID string) (result []uint32, err error) {
	l, err := exec.Exec(groupID)
	if err != nil {
		return nil, err
	} else {
		if v, ok := l.(string); ok {
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

// Get the message sequence number list by peer user ID
func (i *LocalChatLogs) GetMsgSeqListByPeerUserID(ctx context.Context, userID string) (result []uint32, err error) {
	l, err := exec.Exec(userID)
	if err != nil {
		return nil, err
	} else {
		if v, ok := l.(string); ok {
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

// Get the message sequence number list by self user ID
func (i *LocalChatLogs) GetMsgSeqListBySelfUserID(ctx context.Context, userID string) (result []uint32, err error) {
	l, err := exec.Exec(userID)
	if err != nil {
		return nil, err
	} else {
		if v, ok := l.(string); ok {
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

// Get the abnormal message sequence number
func (i *LocalChatLogs) GetAbnormalMsgSeq(ctx context.Context) (int64, error) {
	result, err := exec.Exec()
	if err != nil {
		return 0, err
	}
	if v, ok := result.(float64); ok {
		return int64(v), nil
	}
	return 0, exec.ErrType
}

// Get the list of abnormal message sequence numbers
func (i *LocalChatLogs) GetAbnormalMsgSeqList(ctx context.Context) (result []int64, err error) {
	l, err := exec.Exec()
	if err != nil {
		return nil, err
	} else {
		if v, ok := l.(string); ok {
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

// Batch insert exception messages into the chat logs
func (i *LocalChatLogs) BatchInsertExceptionMsg(ctx context.Context, MessageList []*model_struct.LocalErrChatLog) error {
	_, err := exec.Exec(utils.StructToJsonString(MessageList))
	return err
}

// Update the message status to read in the chat logs
func (i *LocalChatLogs) UpdateGroupMessageHasRead(ctx context.Context, msgIDList []string, sessionType int32) error {
	_, err := exec.Exec(utils.StructToJsonString(msgIDList), sessionType)
	return err
}

// Get the message by message ID
func (i *LocalChatLogs) SearchMessageByKeyword(ctx context.Context, contentType []int, keywordList []string, keywordListMatchType int, conversationID string, startTime, endTime int64, offset, count int) (result []*model_struct.LocalChatLog, err error) {
	msgList, err := exec.Exec(conversationID, utils.StructToJsonString(contentType), utils.StructToJsonString(keywordList), keywordListMatchType, startTime, endTime, offset, count)
	if err != nil {
		return nil, err
	} else {
		if v, ok := msgList.(string); ok {
			var temp []model_struct.LocalChatLog
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

// GetSuperGroupAbnormalMsgSeq get super group abnormal msg seq
func (i *LocalChatLogs) GetSuperGroupAbnormalMsgSeq(ctx context.Context, groupID string) (uint32, error) {
	isExist, err := exec.Exec(groupID)
	if err != nil {
		return 0, err
	} else {
		if v, ok := isExist.(uint32); ok {
			return v, nil
		} else {
			return 0, exec.ErrType
		}
	}
}

// GetAlreadyExistSeqList get already exist seq list
func (i *LocalChatLogs) GetAlreadyExistSeqList(ctx context.Context, conversationID string, lostSeqList []int64) (result []int64, err error) {
	seqList, err := exec.Exec(conversationID, utils.StructToJsonString(lostSeqList))
	if err != nil {
		return nil, err
	} else {
		if v, ok := seqList.(string); ok {
			err := utils.JsonStringToStruct(v, &result)
			if err != nil {
				return nil, err
			}
			return result, nil
		} else {
			return nil, exec.ErrType
		}
	}
}

// GetMessagesBySeq get message by seq
func (i *LocalChatLogs) GetMessageBySeq(ctx context.Context, conversationID string, seq int64) (*model_struct.LocalChatLog, error) {
	msg, err := exec.Exec(conversationID, seq)
	if err != nil {
		return nil, err
	} else {
		if v, ok := msg.(string); ok {
			result := model_struct.LocalChatLog{}
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

// UpdateMessageBySeq update message
func (i *LocalChatLogs) UpdateMessageBySeq(ctx context.Context, conversationID string, c *model_struct.LocalChatLog) error {
	if c.Seq == 0 {
		return exec.PrimaryKeyNull
	}
	tempLocalChatLog := temp_struct.LocalChatLog{
		ServerMsgID:      c.ServerMsgID,
		SendID:           c.SendID,
		RecvID:           c.RecvID,
		SenderPlatformID: c.SenderPlatformID,
		SenderNickname:   c.SenderNickname,
		SenderFaceURL:    c.SenderFaceURL,
		SessionType:      c.SessionType,
		MsgFrom:          c.MsgFrom,
		ContentType:      c.ContentType,
		Content:          c.Content,
		IsRead:           c.IsRead,
		Status:           c.Status,
		Seq:              c.Seq,
		SendTime:         c.SendTime,
		CreateTime:       c.CreateTime,
		AttachedInfo:     c.AttachedInfo,
		Ex:               c.Ex,
	}
	_, err := exec.Exec(conversationID, c.Seq, utils.StructToJsonString(tempLocalChatLog))
	return err
}

func (i *LocalChatLogs) MarkConversationMessageAsReadDB(ctx context.Context, conversationID string, msgIDs []string) (rowsAffected int64, err error) {
	rows, err := exec.Exec(conversationID, utils.StructToJsonString(msgIDs), i.loginUserID)
	if err != nil {
		return 0, err
	} else {
		if v, ok := rows.(float64); ok {
			rowsAffected = int64(v)
			return rowsAffected, err
		} else {
			return 0, exec.ErrType
		}
	}
}

func (i *LocalChatLogs) MarkConversationMessageAsReadBySeqs(ctx context.Context, conversationID string, seqs []int64) (rowsAffected int64, err error) {
	rows, err := exec.Exec(conversationID, utils.StructToJsonString(seqs), i.loginUserID)
	if err != nil {
		return 0, err
	} else {
		if v, ok := rows.(float64); ok {
			rowsAffected = int64(v)
			return rowsAffected, err
		} else {
			return 0, exec.ErrType
		}
	}
}

func (i *LocalChatLogs) MarkConversationAllMessageAsRead(ctx context.Context, conversationID string) (rowsAffected int64, err error) {
	rows, err := exec.Exec(conversationID, i.loginUserID)
	if err != nil {
		return 0, err
	} else {
		if v, ok := rows.(float64); ok {
			rowsAffected = int64(v)
			return rowsAffected, err
		} else {
			return 0, exec.ErrType
		}
	}
}

func (i *LocalChatLogs) GetUnreadMessage(ctx context.Context, conversationID string) (result []*model_struct.LocalChatLog, err error) {
	msgs, err := exec.Exec(conversationID, i.loginUserID)
	if err != nil {
		return nil, err
	} else {
		if v, ok := msgs.(string); ok {
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

func (i *LocalChatLogs) GetMessagesByClientMsgIDs(ctx context.Context, conversationID string, msgIDs []string) (result []*model_struct.LocalChatLog, err error) {
	msgs, err := exec.Exec(conversationID, utils.StructToJsonString(msgIDs))
	if err != nil {
		return nil, err
	} else {
		if v, ok := msgs.(string); ok {
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

// GetMessagesBySeqs gets messages by seqs
func (i *LocalChatLogs) GetMessagesBySeqs(ctx context.Context, conversationID string, seqs []int64) (result []*model_struct.LocalChatLog, err error) {
	msgs, err := exec.Exec(conversationID, utils.StructToJsonString(seqs))
	if err != nil {
		return nil, err
	} else {
		if v, ok := msgs.(string); ok {
			var temp []model_struct.LocalChatLog
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

// GetConversationNormalMsgSeq gets the maximum seq of the session
func (i *LocalChatLogs) GetConversationNormalMsgSeq(ctx context.Context, conversationID string) (int64, error) {
	seq, err := exec.Exec(conversationID)
	if err != nil {
		return 0, err
	} else {
		if v, ok := seq.(float64); ok {
			var result int64
			result = int64(v)
			return result, err
		} else {
			return 0, exec.ErrType
		}
	}
}

func (i *LocalChatLogs) CheckConversationNormalMsgSeq(ctx context.Context, conversationID string) (int64, error) {
	seq, err := exec.Exec(conversationID)
	if err != nil {
		return 0, err
	} else {
		if v, ok := seq.(float64); ok {
			var result int64
			result = int64(v)
			return result, err
		} else {
			return 0, exec.ErrType
		}
	}
}

// GetConversationPeerNormalMsgSeq gets the maximum seq of the peer in the session
func (i *LocalChatLogs) GetConversationPeerNormalMsgSeq(ctx context.Context, conversationID string) (int64, error) {
	seq, err := exec.Exec(conversationID)
	if err != nil {
		return 0, err
	} else {
		if v, ok := seq.(float64); ok {
			var result int64
			result = int64(v)
			return result, nil
		} else {
			return 0, exec.ErrType
		}
	}
}

// GetConversationAbnormalMsgSeq gets the maximum abnormal seq of the session
func (i *LocalChatLogs) GetConversationAbnormalMsgSeq(ctx context.Context, groupID string) (int64, error) {
	//TODO implement me
	panic("implement me")
}

// DeleteConversationAllMessages deletes all messages of the session
func (i *LocalChatLogs) DeleteConversationAllMessages(ctx context.Context, conversationID string) error {
	_, err := exec.Exec(conversationID)
	return err
}

// MarkDeleteConversationAllMessages marks all messages of the session as deleted
func (i *LocalChatLogs) MarkDeleteConversationAllMessages(ctx context.Context, conversationID string) error {
	_, err := exec.Exec(conversationID)
	return err
}

// DeleteConversationMsgs deletes messages of the session
func (i *LocalChatLogs) DeleteConversationMsgs(ctx context.Context, conversationID string, msgIDs []string) error {
	_, err := exec.Exec(conversationID, utils.StructToJsonString(msgIDs))
	return err
}

// DeleteConversationMsgsBySeqs deletes messages of the session
func (i *LocalChatLogs) DeleteConversationMsgsBySeqs(ctx context.Context, conversationID string, seqs []int64) error {
	_, err := exec.Exec(conversationID, utils.StructToJsonString(seqs))
	return err
}
