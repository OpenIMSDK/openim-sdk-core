package conversation_msg

import (
	"context"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/common"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/sdkerrs"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/utils"
	"github.com/openimsdk/openim-sdk-core/v3/sdk_struct"

	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/log"
)

// Delete the local and server
// Delete the local, do not change the server data
// To delete the server, you need to change the local message status to delete
func (c *Conversation) clearConversationFromLocalAndServer(ctx context.Context, conversationID string, f func(ctx context.Context, conversationID string) error) error {
	_, err := c.db.GetConversation(ctx, conversationID)
	if err != nil {
		return err
	}
	// Use conversationID to remove conversations and messages from the server first
	err = c.clearConversationMsgFromServer(ctx, conversationID)
	if err != nil {
		return err
	}
	if err := c.clearConversationAndDeleteAllMsg(ctx, conversationID, false, f); err != nil {
		return err
	}
	c.doUpdateConversation(common.Cmd2Value{Value: common.UpdateConNode{Action: constant.ConChange, Args: []string{conversationID}}})
	c.doUpdateConversation(common.Cmd2Value{Value: common.UpdateConNode{Action: constant.TotalUnreadMessageChanged}})
	return nil
}

func (c *Conversation) clearConversationAndDeleteAllMsg(ctx context.Context, conversationID string, markDelete bool, f func(ctx context.Context, conversationID string) error) error {
	err := c.getConversationMaxSeqAndSetHasRead(ctx, conversationID)
	if err != nil {
		return err
	}
	if markDelete {
		err = c.db.MarkDeleteConversationAllMessages(ctx, conversationID)
	} else {
		err = c.db.DeleteConversationAllMessages(ctx, conversationID)
	}
	if err != nil {
		return err
	}
	log.ZDebug(ctx, "reset conversation", "conversationID", conversationID)
	err = f(ctx, conversationID)
	if err != nil {
		return err
	}
	return nil
}

// Delete all messages
func (c *Conversation) deleteAllMsgFromLocalAndServer(ctx context.Context) error {
	// Delete the server first (high error rate), then delete it.
	err := c.deleteAllMessageFromServer(ctx)
	if err != nil {
		return err
	}
	err = c.deleteAllMsgFromLocal(ctx, false)
	if err != nil {
		return err
	}
	c.doUpdateConversation(common.Cmd2Value{Value: common.UpdateConNode{Action: constant.TotalUnreadMessageChanged}})
	return nil
}

// Delete all messages from the local
func (c *Conversation) deleteAllMsgFromLocal(ctx context.Context, markDelete bool) error {
	conversations, err := c.db.GetAllConversationListDB(ctx)
	if err != nil {
		return err
	}
	var successCids []string
	log.ZDebug(ctx, "deleteAllMsgFromLocal", "conversations", conversations, "markDelete", markDelete)
	for _, v := range conversations {
		if err := c.clearConversationAndDeleteAllMsg(ctx, v.ConversationID, markDelete, c.db.ClearConversation); err != nil {
			log.ZError(ctx, "clearConversation err", err, "conversationID", v.ConversationID)
			continue
		}
		successCids = append(successCids, v.ConversationID)
	}
	c.doUpdateConversation(common.Cmd2Value{Value: common.UpdateConNode{Action: constant.ConChange, Args: successCids}})
	c.doUpdateConversation(common.Cmd2Value{Value: common.UpdateConNode{Action: constant.TotalUnreadMessageChanged}})
	return nil

}

// Delete a message from the local
func (c *Conversation) deleteMessage(ctx context.Context, conversationID string, clientMsgID string) error {
	_, err := c.db.GetMessage(ctx, conversationID, clientMsgID)
	if err != nil {
		return err
	}
	localMessage, err := c.db.GetMessage(ctx, conversationID, clientMsgID)
	if err != nil {
		return err
	}
	if localMessage.Status == constant.MsgStatusSendFailed {
		log.ZInfo(ctx, "delete msg status is send failed, do not need delete", "msg", localMessage)
		return nil
	}
	if localMessage.Seq == 0 {
		log.ZInfo(ctx, "delete msg seq is 0, try again", "msg", localMessage)
		return sdkerrs.ErrMsgHasNoSeq
	}
	err = c.deleteMessagesFromServer(ctx, conversationID, []int64{localMessage.Seq})
	if err != nil {
		return err
	}

	return c.deleteMessageFromLocal(ctx, conversationID, clientMsgID)
}

// Delete messages from local
func (c *Conversation) deleteMessageFromLocal(ctx context.Context, conversationID string, clientMsgID string) error {
	s, err := c.db.GetMessage(ctx, conversationID, clientMsgID)
	if err != nil {
		return err
	}

	if err := c.db.UpdateColumnsMessage(ctx, conversationID, clientMsgID, map[string]interface{}{"status": constant.MsgStatusHasDeleted}); err != nil {
		return err
	}

	if !s.IsRead && s.SendID != c.loginUserID {
		if err := c.db.DecrConversationUnreadCount(ctx, conversationID, 1); err != nil {
			return err
		}
		c.doUpdateConversation(common.Cmd2Value{Value: common.UpdateConNode{ConID: conversationID, Action: constant.ConChange, Args: []string{conversationID}}})
		c.doUpdateConversation(common.Cmd2Value{Value: common.UpdateConNode{Action: constant.TotalUnreadMessageChanged}})
	}

	conversation, err := c.db.GetConversation(ctx, conversationID)
	if err != nil {
		return err
	}

	var latestMsg sdk_struct.MsgStruct
	// Convert the latest message in the conversation table.
	utils.JsonStringToStruct(conversation.LatestMsg, &latestMsg)

	if latestMsg.ClientMsgID == clientMsgID {
		log.ZDebug(ctx, "latestMsg deleted", "seq", latestMsg.Seq, "clientMsgID", latestMsg.ClientMsgID)
		msg, err := c.db.GetLatestActiveMessage(ctx, conversationID, false)
		if err != nil {
			return err
		}

		latestMsgSendTime := latestMsg.SendTime
		latestMsgStr := ""
		if len(msg) > 0 {
			latestMsg = *LocalChatLogToMsgStruct(msg[0])

			latestMsgStr = utils.StructToJsonString(latestMsg)
			latestMsgSendTime = latestMsg.SendTime
		}
		if err := c.db.UpdateColumnsConversation(ctx, conversationID, map[string]interface{}{"latest_msg": latestMsgStr, "latest_msg_send_time": latestMsgSendTime}); err != nil {
			return err
		}
		c.doUpdateConversation(common.Cmd2Value{Value: common.UpdateConNode{Action: constant.ConChange, Args: []string{conversationID}}})
	}
	c.msgListener().OnMsgDeleted(utils.StructToJsonString(s))
	return nil
}

func (c *Conversation) doDeleteMsgs(ctx context.Context, msg *sdkws.MsgData) error {
	tips := sdkws.DeleteMsgsTips{}
	utils.UnmarshalNotificationElem(msg.Content, &tips)
	log.ZDebug(ctx, "doDeleteMsgs", "seqs", tips.Seqs)
	for _, v := range tips.Seqs {
		msg, err := c.db.GetMessageBySeq(ctx, tips.ConversationID, v)
		if err != nil {
			log.ZWarn(ctx, "GetMessageBySeq err", err, "conversationID", tips.ConversationID, "seq", v)
			continue
		}
		if err := c.deleteMessageFromLocal(ctx, tips.ConversationID, msg.ClientMsgID); err != nil {
			log.ZWarn(ctx, "deleteMessageFromLocal err", err, "conversationID", tips.ConversationID, "seq", v)
			return err
		}
	}
	return nil
}

func (c *Conversation) doClearConversations(ctx context.Context, msg *sdkws.MsgData) error {
	tips := &sdkws.ClearConversationTips{}
	err := utils.UnmarshalNotificationElem(msg.Content, tips)
	if err != nil {
		return err
	}

	log.ZDebug(ctx, "doClearConversations", "tips", tips)
	for _, v := range tips.ConversationIDs {
		if err := c.clearConversationAndDeleteAllMsg(ctx, v, false, c.db.ClearConversation); err != nil {
			log.ZWarn(ctx, "clearConversation err", err, "conversationID", v)
			return err
		}
	}
	c.doUpdateConversation(common.Cmd2Value{Value: common.UpdateConNode{Action: constant.ConChange, Args: tips.ConversationIDs}})
	c.doUpdateConversation(common.Cmd2Value{Value: common.UpdateConNode{Action: constant.TotalUnreadMessageChanged}})
	return nil
}
