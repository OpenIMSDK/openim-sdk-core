package conversation_msg

import (
	"context"
	"errors"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/common"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/utils"
	"github.com/openimsdk/openim-sdk-core/v3/sdk_struct"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/utils/timeutil"

	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/log"
)

func (c *Conversation) doRevokeMsg(ctx context.Context, msg *sdkws.MsgData) error {
	var tips sdkws.RevokeMsgTips
	if err := utils.UnmarshalNotificationElem(msg.Content, &tips); err != nil {
		log.ZWarn(ctx, "unmarshal failed", err, "msg", msg)
		return errs.Wrap(err)
	}
	log.ZDebug(ctx, "do revokeMessage", "tips", &tips)
	return c.revokeMessage(ctx, &tips)
}

func (c *Conversation) revokeMessage(ctx context.Context, tips *sdkws.RevokeMsgTips) error {
	revokedMsg, err := c.db.GetMessageBySeq(ctx, tips.ConversationID, tips.Seq)
	if err != nil {
		log.ZError(ctx, "GetMessageBySeq failed", err, "tips", &tips)
		return errs.Wrap(err)
	}

	var revokerRole int32
	var revokerNickname string
	if tips.IsAdminRevoke || tips.SesstionType == constant.SingleChatType {
		_, userName, err := c.getUserNameAndFaceURL(ctx, tips.RevokerUserID)
		if err != nil {
			log.ZError(ctx, "GetUserNameAndFaceURL failed", err, "tips", &tips)
			return errs.Wrap(err)
		} else {
			log.ZDebug(ctx, "revoker user name", "userName", userName)
		}

		revokerNickname = userName
	} else if tips.SesstionType == constant.ReadGroupChatType {
		conversation, err := c.db.GetConversation(ctx, tips.ConversationID)
		if err != nil {
			log.ZError(ctx, "GetConversation failed", err, "conversationID", tips.ConversationID)
			return errs.Wrap(err)
		}

		groupMember, err := c.db.GetGroupMemberInfoByGroupIDUserID(ctx, conversation.GroupID, tips.RevokerUserID)
		if err != nil {
			log.ZError(ctx, "GetGroupMemberInfoByGroupIDUserID failed", err, "tips", &tips)
			return errs.Wrap(err)
		} else {
			log.ZDebug(ctx, "revoker member name", "groupMember", groupMember)
			revokerRole = groupMember.RoleLevel
			revokerNickname = groupMember.Nickname
		}
	}
	m := sdk_struct.MessageRevoked{
		RevokerID:                   tips.RevokerUserID,
		RevokerRole:                 revokerRole,
		ClientMsgID:                 revokedMsg.ClientMsgID,
		RevokerNickname:             revokerNickname,
		RevokeTime:                  tips.RevokeTime,
		SourceMessageSendTime:       revokedMsg.SendTime,
		SourceMessageSendID:         revokedMsg.SendID,
		SourceMessageSenderNickname: revokedMsg.SenderNickname,
		SessionType:                 tips.SesstionType,
		Seq:                         tips.Seq,
		Ex:                          revokedMsg.Ex,
		IsAdminRevoke:               tips.IsAdminRevoke,
	}
	// log.ZDebug(ctx, "callback revokeMessage", "m", m)
	var n sdk_struct.NotificationElem
	n.Detail = utils.StructToJsonString(m)
	if err := c.db.UpdateMessageBySeq(ctx, tips.ConversationID, &model_struct.LocalChatLog{Seq: tips.Seq,
		Content: utils.StructToJsonString(n), ContentType: constant.RevokeNotification}); err != nil {
		log.ZError(ctx, "UpdateMessageBySeq failed", err, "tips", &tips)
		return errs.Wrap(err)
	}
	conversation, err := c.db.GetConversation(ctx, tips.ConversationID)
	if err != nil {
		log.ZError(ctx, "GetConversation failed", err, "tips", &tips)
		return errs.Wrap(err)
	}
	var latestMsg sdk_struct.MsgStruct
	utils.JsonStringToStruct(conversation.LatestMsg, &latestMsg)
	log.ZDebug(ctx, "latestMsg", "latestMsg", &latestMsg, "seq", tips.Seq)
	if latestMsg.Seq <= tips.Seq {
		var newLatestMsg sdk_struct.MsgStruct
		msgs, err := c.db.GetMessageList(ctx, tips.ConversationID, 1, 0, false)
		if err != nil || len(msgs) == 0 {
			log.ZError(ctx, "GetMessageListNoTime failed", err, "tips", &tips)
			return errs.Wrap(err)
		}
		log.ZDebug(ctx, "latestMsg is revoked", "seq", tips.Seq, "msg", msgs[0])
		newLatestMsg = *LocalChatLogToMsgStruct(msgs[0])
		log.ZDebug(ctx, "revoke update conversatoin", "msg", utils.StructToJsonString(newLatestMsg))
		if err := c.db.UpdateColumnsConversation(ctx, tips.ConversationID, map[string]interface{}{"latest_msg": utils.StructToJsonString(newLatestMsg),
			"latest_msg_send_time": newLatestMsg.SendTime}); err != nil {
			log.ZError(ctx, "UpdateColumnsConversation failed", err, "newLatestMsg", newLatestMsg)
		} else {
			c.doUpdateConversation(common.Cmd2Value{Value: common.UpdateConNode{Action: constant.ConChange, Args: []string{tips.ConversationID}}})
		}

	}
	c.msgListener().OnNewRecvMessageRevoked(utils.StructToJsonString(m))
	msgList, err := c.db.SearchAllMessageByContentType(ctx, conversation.ConversationID, constant.Quote)
	if err != nil {
		log.ZError(ctx, "SearchAllMessageByContentType failed", err, "tips", &tips)
		return errs.Wrap(err)
	}

	for _, v := range msgList {
		err = c.quoteMsgRevokeHandle(ctx, tips.ConversationID, v, m)
		if err != nil {
			log.ZError(ctx, "quote Msg Revoke Handle failed.", err, "chat Log content", v)
		}
	}
	return errs.Wrap(err)
}

func (c *Conversation) quoteMsgRevokeHandle(ctx context.Context, conversationID string, v *model_struct.LocalChatLog, revokedMsg sdk_struct.MessageRevoked) error {
	s := sdk_struct.QuoteElem{}
	if v.Content == "" {
		return errs.New("Chat Log Content not found")
	}

	if err := utils.JsonStringToStruct(v.Content, &s); err != nil {
		return errs.New("ChatLog content transfer failed.")
	}

	if s.QuoteMessage == nil {
		return errs.New("QuoteMessage is nil").Wrap()
	}
	if s.QuoteMessage.ClientMsgID != revokedMsg.ClientMsgID {
		return nil
	}

	s.QuoteMessage.Content = utils.StructToJsonString(revokedMsg)
	s.QuoteMessage.ContentType = constant.RevokeNotification
	v.Content = utils.StructToJsonString(s)
	if err := c.db.UpdateMessageBySeq(ctx, conversationID, v); err != nil {
		log.ZError(ctx, "UpdateMessage failed", err, "v", v)
		return errs.Wrap(err)
	}
	return nil
}

func (c *Conversation) revokeOneMessage(ctx context.Context, conversationID, clientMsgID string) error {
	conversation, err := c.db.GetConversation(ctx, conversationID)
	if err != nil {
		return err
	}
	message, err := c.db.GetMessage(ctx, conversationID, clientMsgID)
	if err != nil {
		return err
	}
	if message.Status != constant.MsgStatusSendSuccess {
		return errors.New("only send success message can be revoked")
	}
	switch conversation.ConversationType {
	case constant.SingleChatType:
		if message.SendID != c.loginUserID {
			return errors.New("only send by yourself message can be revoked")
		}
	case constant.ReadGroupChatType:
		if message.SendID != c.loginUserID {
			groupAdmins, err := c.db.GetGroupMemberOwnerAndAdminDB(ctx, conversation.GroupID)
			if err != nil {
				return err
			}
			var isAdmin bool
			for _, member := range groupAdmins {
				if member.UserID == c.loginUserID {
					isAdmin = true
					break
				}
			}
			if !isAdmin {
				return errors.New("only group admin can revoke message")
			}
		}
	}

	err = c.revokeMessageFromServer(ctx, conversationID, message.Seq)
	if err != nil {
		return err
	}

	c.revokeMessage(ctx, &sdkws.RevokeMsgTips{
		ConversationID: conversationID,
		Seq:            message.Seq,
		RevokerUserID:  c.loginUserID,
		RevokeTime:     timeutil.GetCurrentTimestampBySecond(),
		SesstionType:   conversation.ConversationType,
		ClientMsgID:    clientMsgID,
	})
	return nil
}
