//go:build !js
// +build !js

package db

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/utils"
	"github.com/openimsdk/openim-sdk-core/v3/sdk_struct"

	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
)

func (d *DataBase) initChatLog(ctx context.Context, conversationID string) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	tableName := utils.GetTableName(conversationID)
	if !d.tableChecker.HasTable(tableName) {
		createTableSQL := fmt.Sprintf(`
            CREATE TABLE "%s" (
                client_msg_id CHAR(64),
                server_msg_id CHAR(64),
                send_id CHAR(64),
                recv_id CHAR(64),
                sender_platform_id INTEGER,
                sender_nick_name VARCHAR(255),
                sender_face_url VARCHAR(255),
                session_type INTEGER,
                msg_from INTEGER,
                content_type INTEGER,
                content VARCHAR(1000),
                is_read NUMERIC,
                status INTEGER,
                seq INTEGER DEFAULT 0,
                send_time INTEGER,
                create_time INTEGER,
                attached_info VARCHAR(1024),
                ex VARCHAR(1024),
                local_ex VARCHAR(1024),
                is_react NUMERIC,
                is_external_extensions NUMERIC,
                msg_first_modify_time INTEGER,
                PRIMARY KEY (client_msg_id)
            );`, tableName)

		if result := d.conn.Exec(createTableSQL); result.Error != nil {
			return errs.WrapMsg(result.Error, "Create table failed", "table", tableName)
		}
		result := d.conn.Exec(fmt.Sprintf("CREATE INDEX `%s` ON `%s` (seq)", "index_seq_"+conversationID, tableName))
		if result.Error != nil {
			return errs.WrapMsg(result.Error, "Create index_seq failed", "table", tableName, "index", "index_seq_"+conversationID)
		}
		result = d.conn.Exec(fmt.Sprintf("CREATE INDEX `%s` ON `%s` (send_time)", "index_send_time_"+conversationID, tableName))
		if result.Error != nil {
			return errs.WrapMsg(result.Error, "Create index_send_time failed", "table", tableName, "index", "index_send_time_"+conversationID)
		}
		d.tableChecker.UpdateTable(tableName)
	}
	return nil
}

func (d *DataBase) UpdateMessage(ctx context.Context, conversationID string, c *model_struct.LocalChatLog) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	t := d.conn.WithContext(ctx).Table(utils.GetTableName(conversationID)).Updates(c)
	if t.RowsAffected == 0 {
		return errs.WrapMsg(errors.New("RowsAffected == 0"), "no update ")
	}
	return errs.WrapMsg(t.Error, "UpdateMessage failed")
}

func (d *DataBase) UpdateMessageBySeq(ctx context.Context, conversationID string, c *model_struct.LocalChatLog) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return errs.WrapMsg(d.conn.WithContext(ctx).Table(utils.GetTableName(conversationID)).Where("seq=?", c.Seq).Updates(c).Error, "UpdateMessage failed")
}

func (d *DataBase) BatchInsertMessageList(ctx context.Context, conversationID string, MessageList []*model_struct.LocalChatLog) error {
	err := d.initChatLog(ctx, conversationID)
	if err != nil {
		log.ZWarn(ctx, "initChatLog err", err)
		return err
	}

	if MessageList == nil {
		return nil
	}
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return errs.WrapMsg(d.conn.WithContext(ctx).Table(utils.GetTableName(conversationID)).Create(MessageList).Error, "BatchInsertMessageList failed")
}

func (d *DataBase) InsertMessage(ctx context.Context, conversationID string, Message *model_struct.LocalChatLog) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return errs.WrapMsg(d.conn.WithContext(ctx).Table(utils.GetTableName(conversationID)).Create(Message).Error, "InsertMessage failed")
}

func (d *DataBase) GetMessage(ctx context.Context, conversationID string, clientMsgID string) (*model_struct.LocalChatLog, error) {
	err := d.initChatLog(ctx, conversationID)
	if err != nil {
		log.ZWarn(ctx, "initChatLog err", err)
		return nil, err
	}
	d.mRWMutex.RLock()
	defer d.mRWMutex.RUnlock()
	var c model_struct.LocalChatLog
	return &c, errs.WrapMsg(d.conn.WithContext(ctx).Table(utils.GetTableName(conversationID)).Where("client_msg_id = ?",
		clientMsgID).Take(&c).Error, "GetMessage failed")
}

func (d *DataBase) GetMessageBySeq(ctx context.Context, conversationID string, seq int64) (*model_struct.LocalChatLog, error) {
	err := d.initChatLog(ctx, conversationID)
	if err != nil {
		log.ZWarn(ctx, "initChatLog err", err)
		return nil, err
	}
	d.mRWMutex.RLock()
	defer d.mRWMutex.RUnlock()
	var c model_struct.LocalChatLog
	return &c, errs.WrapMsg(d.conn.WithContext(ctx).Table(utils.GetTableName(conversationID)).Where("seq = ?",
		seq).Take(&c).Error, "GetMessage failed")
}

func (d *DataBase) UpdateMessageTimeAndStatus(ctx context.Context, conversationID, clientMsgID string, serverMsgID string, sendTime int64, status int32) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return errs.WrapMsg(d.conn.WithContext(ctx).Table(utils.GetTableName(conversationID)).Model(model_struct.LocalChatLog{}).Where("client_msg_id=? And seq=?", clientMsgID, 0).
		Updates(model_struct.LocalChatLog{Status: status, SendTime: sendTime, ServerMsgID: serverMsgID}).Error, "UpdateMessageStatusBySourceID failed")
}

func (d *DataBase) GetMessageList(ctx context.Context, conversationID string, count int, startTime int64, isReverse bool) (result []*model_struct.LocalChatLog, err error) {
	if err = d.initChatLog(ctx, conversationID); err != nil {
		log.ZWarn(ctx, "initChatLog err", err)
		return nil, err
	}
	d.mRWMutex.RLock()
	defer d.mRWMutex.RUnlock()
	var condition, timeOrder, timeSymbol string
	if isReverse {
		timeOrder = "send_time ASC"
		timeSymbol = ">"
	} else {
		timeOrder = "send_time DESC"
		timeSymbol = "<"
	}
	if startTime > 0 {
		condition = "send_time " + timeSymbol + " ?"
		err = errs.WrapMsg(d.conn.WithContext(ctx).Table(utils.GetTableName(conversationID)).Where(condition, startTime).
			Order(timeOrder).Offset(0).Limit(count).Find(&result).Error, "GetMessageList failed")
		if err != nil {
			return nil, err
		}
	} else {
		err = errs.WrapMsg(d.conn.WithContext(ctx).Table(utils.GetTableName(conversationID)).Order(timeOrder).
			Offset(0).Limit(count).Find(&result).Error, "GetMessageList failed")
		if err != nil {
			return nil, err
		}
	}
	return result, err
}

func (d *DataBase) DeleteConversationAllMessages(ctx context.Context, conversationID string) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return errs.WrapMsg(d.conn.WithContext(ctx).Table(utils.GetTableName(conversationID)).Where("1 = 1").Delete(model_struct.LocalChatLog{}).Error, "DeleteConversationAllMessages failed")
}

func (d *DataBase) MarkDeleteConversationAllMessages(ctx context.Context, conversationID string) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return errs.WrapMsg(d.conn.WithContext(ctx).Table(utils.GetTableName(conversationID)).Where("1 = 1").Updates(model_struct.LocalChatLog{Status: constant.MsgStatusHasDeleted}).Error, "DeleteConversationAllMessages failed")
}

func (d *DataBase) SearchMessageByContentType(ctx context.Context, contentType []int, conversationID string, startTime, endTime int64, offset, count int) (result []*model_struct.LocalChatLog, err error) {
	d.mRWMutex.RLock()
	defer d.mRWMutex.RUnlock()
	condition := fmt.Sprintf("send_time between %d and %d AND status <=%d And content_type IN ?", startTime, endTime, constant.MsgStatusSendFailed)
	err = errs.WrapMsg(d.conn.WithContext(ctx).Table(utils.GetTableName(conversationID)).Where(condition, contentType).Order("send_time DESC").Offset(offset).Limit(count).Find(&result).Error, "SearchMessage failed")
	return result, err
}

func (d *DataBase) SearchMessageByKeyword(ctx context.Context, contentType []int, keywordList []string, keywordListMatchType int, conversationID string, startTime, endTime int64, offset, count int) (result []*model_struct.LocalChatLog, err error) {
	d.mRWMutex.RLock()
	defer d.mRWMutex.RUnlock()
	var condition string
	var subCondition string
	if keywordListMatchType == constant.KeywordMatchOr {
		for i := 0; i < len(keywordList); i++ {
			if i == 0 {
				subCondition += "And ("
			}
			if i+1 >= len(keywordList) {
				subCondition += "content like " + "'%" + keywordList[i] + "%') "
			} else {
				subCondition += "content like " + "'%" + keywordList[i] + "%' " + "or "

			}
		}
	} else {
		for i := 0; i < len(keywordList); i++ {
			if i == 0 {
				subCondition += "And ("
			}
			if i+1 >= len(keywordList) {
				subCondition += "content like " + "'%" + keywordList[i] + "%') "
			} else {
				subCondition += "content like " + "'%" + keywordList[i] + "%' " + "and "
			}
		}
	}
	condition = fmt.Sprintf(" send_time  between %d and %d AND status <=%d  And content_type IN ? ", startTime, endTime, constant.MsgStatusSendFailed)
	condition += subCondition
	err = errs.WrapMsg(d.conn.WithContext(ctx).Table(utils.GetTableName(conversationID)).Where(condition, contentType).Order("send_time DESC").Offset(offset).Limit(count).Find(&result).Error, "SearchMessage failed")
	return result, err
}

// SearchMessageByContentTypeAndKeyword searches for messages in the database that match specified content types and keywords within a given time range.
func (d *DataBase) SearchMessageByContentTypeAndKeyword(ctx context.Context, contentType []int, conversationID string, keywordList []string, keywordListMatchType int, startTime, endTime int64) (result []*model_struct.LocalChatLog, err error) {
	d.mRWMutex.RLock()
	defer d.mRWMutex.RUnlock()
	var condition string
	var subCondition string

	// Construct a sub-condition for SQL query based on keyword list and match type
	if keywordListMatchType == constant.KeywordMatchOr {
		// Use OR logic if keywordListMatchType is KeywordMatchOr
		for i := 0; i < len(keywordList); i++ {
			if i == 0 {
				subCondition += "And ("
			}
			if i+1 >= len(keywordList) {
				subCondition += "content like " + "'%" + keywordList[i] + "%') "
			} else {
				subCondition += "content like " + "'%" + keywordList[i] + "%' " + "or "
			}
		}
	} else {
		// Use AND logic for other keywordListMatchType
		for i := 0; i < len(keywordList); i++ {
			if i == 0 {
				subCondition += "And ("
			}
			if i+1 >= len(keywordList) {
				subCondition += "content like " + "'%" + keywordList[i] + "%') "
			} else {
				subCondition += "content like " + "'%" + keywordList[i] + "%' " + "and "
			}
		}
	}

	// Construct the main SQL condition string
	condition = fmt.Sprintf("send_time between %d and %d AND status <=%d  And content_type IN ? ", startTime, endTime, constant.MsgStatusSendFailed)
	condition += subCondition

	// Execute the query using the constructed condition and handle errors
	err = errs.WrapMsg(d.conn.WithContext(ctx).Table(utils.GetTableName(conversationID)).Where(condition, contentType).Order("send_time DESC").Find(&result).Error, "SearchMessage failed")

	return result, err
}

func (d *DataBase) UpdateMsgSenderFaceURLAndSenderNickname(ctx context.Context, conversationID, sendID, faceURL, nickname string) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return errs.WrapMsg(d.conn.WithContext(ctx).Table(utils.GetTableName(conversationID)).Model(model_struct.LocalChatLog{}).Where(
		"send_id = ?", sendID).Updates(
		map[string]interface{}{"sender_face_url": faceURL, "sender_nick_name": nickname}).Error, utils.GetSelfFuncName()+" failed")
}

func (d *DataBase) UpdateColumnsMessage(ctx context.Context, conversationID, ClientMsgID string, args map[string]interface{}) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	c := model_struct.LocalChatLog{ClientMsgID: ClientMsgID}
	t := d.conn.WithContext(ctx).Table(utils.GetConversationTableName(conversationID)).Model(&c).Updates(args)
	if t.RowsAffected == 0 {
		return errs.WrapMsg(errors.New("RowsAffected == 0"), "no update")
	}
	return errs.WrapMsg(t.Error, "UpdateColumnsConversation failed")
}

func (d *DataBase) SearchAllMessageByContentType(ctx context.Context, conversationID string, contentType int) (result []*model_struct.LocalChatLog, err error) {
	d.mRWMutex.RLock()
	defer d.mRWMutex.RUnlock()
	err = d.conn.WithContext(ctx).Table(utils.GetTableName(conversationID)).Model(&model_struct.LocalChatLog{}).Where("content_type = ?", contentType).Find(&result).Error
	return result, err
}

func (d *DataBase) GetUnreadMessage(ctx context.Context, conversationID string) (msgs []*model_struct.LocalChatLog, err error) {
	d.mRWMutex.RLock()
	defer d.mRWMutex.RUnlock()
	err = errs.WrapMsg(d.conn.WithContext(ctx).Table(utils.GetConversationTableName(conversationID)).Debug().Where("send_id != ? AND is_read = ?", d.loginUserID, constant.NotRead).Find(&msgs).Error, "GetMessageList failed")
	return msgs, err
}

func (d *DataBase) MarkConversationMessageAsReadBySeqs(ctx context.Context, conversationID string, seqs []int64) (rowsAffected int64, err error) {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	t := d.conn.WithContext(ctx).Table(utils.GetConversationTableName(conversationID)).Where("seq in ? AND send_id != ?", seqs, d.loginUserID).Update("is_read", constant.HasRead)
	if t.RowsAffected == 0 {
		return 0, errs.WrapMsg(errors.New("RowsAffected == 0"), "no update")
	}
	return t.RowsAffected, errs.WrapMsg(t.Error, "UpdateMessageStatusBySourceID failed")
}

func (d *DataBase) MarkConversationMessageAsReadDB(ctx context.Context, conversationID string, msgIDs []string) (rowsAffected int64, err error) {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	var msgs []*model_struct.LocalChatLog
	if err := d.conn.WithContext(ctx).Table(utils.GetConversationTableName(conversationID)).Where("client_msg_id in ? AND send_id != ?", msgIDs, d.loginUserID).Find(&msgs).Error; err != nil {
		return 0, errs.WrapMsg(err, "MarkConversationMessageAsReadDB failed")
	}
	for _, msg := range msgs {
		var attachedInfo sdk_struct.AttachedInfoElem
		utils.JsonStringToStruct(msg.AttachedInfo, &attachedInfo)
		attachedInfo.HasReadTime = utils.GetCurrentTimestampByMill()
		msg.IsRead = true
		msg.AttachedInfo = utils.StructToJsonString(attachedInfo)
		if err := d.conn.WithContext(ctx).Table(utils.GetConversationTableName(conversationID)).Where("client_msg_id = ?", msg.ClientMsgID).Updates(msg).Error; err != nil {
			log.ZError(ctx, "MarkConversationMessageAsReadDB failed", err, "msg", msg)
		} else {
			rowsAffected++
		}
	}
	// t := d.conn.WithContext(ctx).Table(utils.GetConversationTableName(conversationID)).Where("client_msg_id in ? AND send_id != ?", msgIDs, d.loginUserID).Update("is_read", constant.HasRead)
	// if t.RowsAffected == 0 {
	// 	return 0, errs.WrapMsg(errors.New("RowsAffected == 0"), "no update")
	// }
	return rowsAffected, nil
}

func (d *DataBase) GetMessagesByClientMsgIDs(ctx context.Context, conversationID string, msgIDs []string) (msgs []*model_struct.LocalChatLog, err error) {
	d.mRWMutex.RLock()
	defer d.mRWMutex.RUnlock()
	err = errs.WrapMsg(d.conn.WithContext(ctx).Table(utils.GetConversationTableName(conversationID)).Where("client_msg_id IN ?", msgIDs).Order("send_time DESC").Find(&msgs).Error, "GetMessagesByClientMsgIDs error")
	return msgs, err
}

func (d *DataBase) GetMessagesBySeqs(ctx context.Context, conversationID string, seqs []int64) (msgs []*model_struct.LocalChatLog, err error) {
	d.mRWMutex.RLock()
	defer d.mRWMutex.RUnlock()
	err = errs.WrapMsg(d.conn.WithContext(ctx).Table(utils.GetConversationTableName(conversationID)).Where("seq IN ?", seqs).Order("send_time DESC").Find(&msgs).Error, "GetMessagesBySeqs error")
	return msgs, err
}

func (d *DataBase) GetConversationNormalMsgSeq(ctx context.Context, conversationID string) (int64, error) {
	err := d.initChatLog(ctx, conversationID)
	if err != nil {
		log.ZWarn(ctx, "initChatLog err", err)
		return 0, err
	}
	d.mRWMutex.RLock()
	defer d.mRWMutex.RUnlock()
	var seq int64
	err = d.conn.WithContext(ctx).Table(utils.GetConversationTableName(conversationID)).Select("IFNULL(max(seq),0)").Find(&seq).Error
	return seq, errs.WrapMsg(err, "GetConversationNormalMsgSeq")
}

func (d *DataBase) CheckConversationNormalMsgSeq(ctx context.Context, conversationID string) (int64, error) {
	var seq int64
	d.mRWMutex.RLock()
	defer d.mRWMutex.RUnlock()
	if d.tableChecker.HasTable(utils.GetConversationTableName(conversationID)) {
		err := d.conn.WithContext(ctx).Table(utils.GetConversationTableName(conversationID)).Select("IFNULL(max(seq),0)").Find(&seq).Error
		return seq, errs.Wrap(err)
	}
	return 0, nil
}

func (d *DataBase) GetConversationPeerNormalMsgSeq(ctx context.Context, conversationID string) (int64, error) {
	d.mRWMutex.RLock()
	defer d.mRWMutex.RUnlock()
	var seq int64
	err := d.conn.WithContext(ctx).Table(utils.GetConversationTableName(conversationID)).Select("IFNULL(max(seq),0)").Where("send_id != ?", d.loginUserID).Find(&seq).Error
	return seq, errs.WrapMsg(err, "GetConversationPeerNormalMsgSeq")
}

func (d *DataBase) GetLatestActiveMessage(ctx context.Context, conversationID string, isReverse bool) (result []*model_struct.LocalChatLog, err error) {
	if err = d.initChatLog(ctx, conversationID); err != nil {
		log.ZWarn(ctx, "initChatLog err", err)
		return nil, err
	}

	d.mRWMutex.RLock()
	defer d.mRWMutex.RUnlock()

	var timeOrder string
	if isReverse {
		timeOrder = "send_time ASC"
	} else {
		timeOrder = "send_time DESC"
	}

	// only get status < 4(NotHasDeleted) Msg
	err = errs.WrapMsg(d.conn.WithContext(ctx).Table(utils.GetTableName(conversationID)).Where("status < ?", constant.MsgStatusHasDeleted).Order(timeOrder).Offset(0).Limit(1).Find(&result).Error, "GetLatestActiveMessage failed")
	if err != nil {
		return nil, err
	}

	return result, err
}

func (d *DataBase) GetLatestValidServerMessage(ctx context.Context, conversationID string, startTime int64, isReverse bool) (*model_struct.LocalChatLog, error) {
	d.mRWMutex.RLock()
	defer d.mRWMutex.RUnlock()

	var condition, timeOrder, timeSymbol string
	var result model_struct.LocalChatLog

	if isReverse {
		timeOrder = "send_time DESC"
		timeSymbol = "<"
	} else {
		timeOrder = "send_time ASC"
		timeSymbol = ">"
	}

	condition = "send_time " + timeSymbol + " ? AND seq != ?"

	err := d.conn.WithContext(ctx).Table(utils.GetTableName(conversationID)).
		Where(condition, startTime, 0).
		Order(timeOrder).
		Limit(1).
		First(&result).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, errs.WrapMsg(err, "GetLatestValidServerMessage failed")
	}

	return &result, nil
}
