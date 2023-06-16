// Copyright © 2023 OpenIM SDK. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package conversation_msg

import (
	"context"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/sdk_struct"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
)

// Marks a set of messages in a specified session as read and sends a request for the mark to the server.
func (c *Conversation) markMsgAsRead2Svr(ctx context.Context, conversationID string, seqs []int64) error {
	var req sdkws.MarkMsgsAsReadReq
	req.UserID = c.loginUserID
	req.ConversationID = conversationID
	req.Seqs = seqs
	err := c.SendReqWaitResp(ctx, req, constant.HasReadMsg, nil)
	if err != nil {
		return err
	}
	return nil
}

// Marks the specified session as read and sends the mark request to the server.
func (c *Conversation) markConversationAsReadSvr(ctx context.Context, conversationID string, hasReadSeq int64) error {
	var markReadReq sdkws.MarkReadReq
	markReadReq.UserID = c.loginUserID
	markReadReq.ConversationID = conversationID
	markReadReq.HasReadSeq = hasReadSeq
	err := c.SendReqWaitResp(ctx, &markReadReq, constant.HasReadMsg, nil)
	if err != nil {
		log.ZError(ctx, "markConversationAsReadSvr failed", err, "req", markReadReq)
		return err
	}
	return nil
}

// Sets the read sequence number for a given session to the specified value and stores that value in the database.
func (c *Conversation) setConversationHasReadSeq(ctx context.Context, conversationID string, hasReadSeq int64) error {
	// req := &pbMsg.SetConversationHasReadSeqReq{UserID: c.loginUserID, ConversationID: conversationID, HasReadSeq: hasReadSeq}
	var req sdkws.SetConversationHasReadSeqReq
	req.UserID = c.loginUserID
	req.ConversationID = conversationID
	req.HasReadSeq = hasReadSeq
	err := c.LongConnMgr.SendReqWaitResp(ctx, req, constant.SetConversationHasReadSeq, nil)
	if err != nil {
		log.ZError(ctx, "setConversationHasReadSeq failed", err, "conversationID", conversationID, "hasReadSeq", hasReadSeq)
		return err
	}
	return nil
}

func (c *Conversation) getConversationMaxSeqAndSetHasRead(ctx context.Context, conversationID string) error {
	maxSeq, err := c.db.GetConversationNormalMsgSeq(ctx, conversationID)
	if err != nil {
		return err
	}
	if err := c.setConversationHasReadSeq(ctx, conversationID, maxSeq); err != nil {
		return err
	}
	if err := c.db.UpdateColumnsConversation(ctx, conversationID, map[string]interface{}{"has_read_seq": maxSeq}); err != nil {
		return err
	}
	return nil
}

// mark a conversation's all message as read
func (c *Conversation) markConversationMessageAsRead(ctx context.Context, conversationID string) error {
	_, err := c.db.GetConversation(ctx, conversationID)
	if err != nil {
		return err
	}
	peerUserMaxSeq, err := c.db.GetConversationPeerNormalMsgSeq(ctx, conversationID)
	if err != nil {
		return err
	}
	maxSeq, err := c.db.GetConversationNormalMsgSeq(ctx, conversationID)
	if err != nil {
		return err
	}
	msgs, err := c.db.GetUnreadMessage(ctx, conversationID)
	if err != nil {
		return err
	}
	log.ZDebug(ctx, "get unread message", "msgs", len(msgs))
	msgIDs, seqs := c.getAsReadMsgMapAndList(ctx, msgs)
	log.ZDebug(ctx, "markConversationMessageAsRead", "conversationID", conversationID, "seqs", seqs, "peerUserMaxSeq", peerUserMaxSeq, "maxSeq", maxSeq)
	if err := c.markConversationAsReadSvr(ctx, conversationID, maxSeq); err != nil {
		return err
	}
	_, err = c.db.MarkConversationMessageAsRead(ctx, conversationID, msgIDs)
	if err != nil {
		log.ZWarn(ctx, "MarkConversationMessageAsRead err", err, "conversationID", conversationID, "msgIDs", msgIDs)
	}
	if err := c.db.UpdateColumnsConversation(ctx, conversationID, map[string]interface{}{"unread_count": 0}); err != nil {
		log.ZError(ctx, "UpdateColumnsConversation err", err, "conversationID", conversationID)
	}
	log.ZDebug(ctx, "update columns sucess")
	c.unreadChangeTrigger(ctx, conversationID, peerUserMaxSeq == maxSeq)
	return nil
}

// mark a conversation's message as read by seqs
func (c *Conversation) markMessagesAsReadByMsgID(ctx context.Context, conversationID string, msgIDs []string) error {
	_, err := c.db.GetConversation(ctx, conversationID)
	if err != nil {
		return err
	}
	msgs, err := c.db.GetMessagesByClientMsgIDs(ctx, conversationID, msgIDs)
	if err != nil {
		return err
	}
	if len(msgs) == 0 {
		return nil
	}
	var hasReadSeq = msgs[0].Seq
	maxSeq, err := c.db.GetConversationNormalMsgSeq(ctx, conversationID)
	if err != nil {
		return err
	}
	markAsReadMsgIDs, seqs := c.getAsReadMsgMapAndList(ctx, msgs)
	if err := c.markMsgAsRead2Svr(ctx, conversationID, seqs); err != nil {
		return err
	}
	decrCount, err := c.db.MarkConversationMessageAsRead(ctx, conversationID, markAsReadMsgIDs)
	if err != nil {
		return err
	}
	if err := c.db.DecrConversationUnreadCount(ctx, conversationID, decrCount); err != nil {
		log.ZError(ctx, "decrConversationUnreadCount err", err, "conversationID", conversationID, "decrCount", decrCount)
	}
	c.unreadChangeTrigger(ctx, conversationID, hasReadSeq == maxSeq && msgs[0].SendID != c.loginUserID)
	return nil
}

func (c *Conversation) getAsReadMsgMapAndList(ctx context.Context, msgs []*model_struct.LocalChatLog) (asReadMsgIDs []string, seqs []int64) {
	for _, msg := range msgs {
		if !msg.IsRead && msg.SendID != c.loginUserID {
			asReadMsgIDs = append(asReadMsgIDs, msg.ClientMsgID)
			seqs = append(seqs, msg.Seq)
		} else {
			log.ZWarn(ctx, "msg can't marked as read", nil, "msg", msg)
		}
	}
	return
}

func (c *Conversation) unreadChangeTrigger(ctx context.Context, conversationID string, latestMsgIsRead bool) {
	if latestMsgIsRead {
		c.doUpdateConversation(common.Cmd2Value{Value: common.UpdateConNode{ConID: conversationID, Action: constant.UpdateLatestMessageChange, Args: []string{conversationID}}, Ctx: ctx})
	}
	c.doUpdateConversation(common.Cmd2Value{Value: common.UpdateConNode{ConID: conversationID, Action: constant.ConChange, Args: []string{conversationID}}, Ctx: ctx})
	c.doUpdateConversation(common.Cmd2Value{Value: common.UpdateConNode{Action: constant.TotalUnreadMessageChanged}, Ctx: ctx})
}

// Updates the number of unread messages in the specified session as the difference
// between the current read message sequence number and the read message sequence number in the database
func (c *Conversation) doUnreadCount(ctx context.Context, conversationID string, hasReadSeq int64) {
	conversation, err := c.db.GetConversation(ctx, conversationID)
	if err != nil {
		log.ZError(ctx, "GetConversation err", err, "conversationID", conversationID)
		return
	}
	var seqs []int64
	if hasReadSeq > conversation.HasReadSeq {
		for i := conversation.HasReadSeq + 1; i <= hasReadSeq; i++ {
			seqs = append(seqs, i)
		}
		_, err := c.db.MarkConversationMessageAsReadBySeqs(ctx, conversationID, seqs)
		if err != nil {
			log.ZWarn(ctx, "MarkConversationMessageAsReadBySeqs err", err, "conversationID", conversationID, "seqs", seqs)
		}
		if err := c.db.DecrConversationUnreadCount(ctx, conversationID, int64(len(seqs))); err != nil {
			log.ZError(ctx, "decrConversationUnreadCount err", err, "conversationID", conversationID, "decrCount", int64(len(seqs)))
		}
		if err := c.db.UpdateColumnsConversation(ctx, conversationID, map[string]interface{}{"has_read_seq": hasReadSeq}); err != nil {
			log.ZError(ctx, "UpdateColumnsConversation err", err, "conversationID", conversationID)
		}
		c.doUpdateConversation(common.Cmd2Value{Value: common.UpdateConNode{ConID: conversationID, Action: constant.ConChange, Args: []string{conversationID}}})
		c.doUpdateConversation(common.Cmd2Value{Value: common.UpdateConNode{Action: constant.TotalUnreadMessageChanged}})
	} else {
		log.ZWarn(ctx, "hasReadSeq <= conversation.HasReadSeq", nil, "hasReadSeq", hasReadSeq, "conversation.HasReadSeq", conversation.HasReadSeq)
	}
}

func (c *Conversation) doReadDrawing(ctx context.Context, msg *sdkws.MsgData) {
	tips := &sdkws.MarkAsReadTips{}
	utils.UnmarshalNotificationElem(msg.Content, tips)

	if tips.MarkAsReadUserID == c.loginUserID {
		log.ZDebug(ctx, "do unread count", "tips", tips)
		c.doUnreadCount(ctx, tips.ConversationID, tips.HasReadSeq)
		return
	}

	log.ZDebug(ctx, "do readDrawing", "tips", tips)

	conversation, err := c.db.GetConversation(ctx, tips.ConversationID)
	if err != nil {
		log.ZError(ctx, "GetConversation err", err, "conversationID", tips.ConversationID)
		return
	}

	messages, err := c.db.GetMessagesBySeqs(ctx, tips.ConversationID, tips.Seqs)
	if err != nil {
		log.ZError(ctx, "GetMessagesBySeqs err", err, "conversationID", tips.ConversationID, "seqs", tips.Seqs)
		return
	}

	var successMsgIDs []string

	for _, message := range messages {
		attachInfo := sdk_struct.AttachedInfoElem{}
		_ = utils.JsonStringToStruct(message.AttachedInfo, &attachInfo)

		attachInfo.HasReadTime = msg.SendTime

		if conversation.ConversationType == constant.SingleChatType {
			message.IsRead = true
		} else if conversation.ConversationType == constant.SuperGroupChatType {
			attachInfo.GroupHasReadInfo.HasReadUserIDList = utils.RemoveRepeatedStringInList(append(attachInfo.GroupHasReadInfo.HasReadUserIDList, tips.MarkAsReadUserID))
			attachInfo.GroupHasReadInfo.HasReadCount = int32(len(attachInfo.GroupHasReadInfo.HasReadUserIDList))
		}

		message.AttachedInfo = utils.StructToJsonString(attachInfo)

		if err = c.db.UpdateMessage(ctx, tips.ConversationID, message); err != nil {
			log.ZError(ctx, "UpdateMessage err", err, "conversationID", tips.ConversationID, "message", message)
		} else {
			successMsgIDs = append(successMsgIDs, message.ClientMsgID)
		}
	}

	var messageReceiptResp []*sdk_struct.MessageReceipt

	if conversation.ConversationType == constant.SingleChatType {
		messageReceiptResp = []*sdk_struct.MessageReceipt{{UserID: tips.MarkAsReadUserID, MsgIDList: successMsgIDs, SessionType: conversation.ConversationType, ReadTime: msg.SendTime}}
		c.msgListener.OnRecvC2CReadReceipt(utils.StructToJsonString(messageReceiptResp))
	} else if conversation.ConversationType == constant.SuperGroupChatType {
		messageReceiptResp = []*sdk_struct.MessageReceipt{{GroupID: conversation.GroupID, MsgIDList: successMsgIDs, SessionType: conversation.ConversationType, ReadTime: msg.SendTime}}
		c.msgListener.OnRecvGroupReadReceipt(utils.StructToJsonString(messageReceiptResp))
	}
}
