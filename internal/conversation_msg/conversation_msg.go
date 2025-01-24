package conversation_msg

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"sync"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/api"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/cache"
	"github.com/openimsdk/tools/utils/stringutil"

	"github.com/openimsdk/openim-sdk-core/v3/internal/group"
	"github.com/openimsdk/openim-sdk-core/v3/internal/interaction"
	"github.com/openimsdk/openim-sdk-core/v3/internal/relation"
	"github.com/openimsdk/openim-sdk-core/v3/internal/third/file"
	"github.com/openimsdk/openim-sdk-core/v3/internal/user"
	"github.com/openimsdk/openim-sdk-core/v3/open_im_sdk_callback"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/ccontext"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/common"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/db_interface"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/page"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/syncer"
	pbConversation "github.com/openimsdk/protocol/conversation"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/utils/datautil"

	"sort"
	"time"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/utils"
	"github.com/openimsdk/openim-sdk-core/v3/sdk_struct"

	"github.com/jinzhu/copier"
)

const (
	conversationSyncLimit       int64 = math.MaxInt64
	searchMessageGoroutineLimit       = 10
)

var SearchContentType = []int{constant.Text, constant.AtText, constant.File}

type Conversation struct {
	*interaction.LongConnMgr
	conversationSyncer          *syncer.Syncer[*model_struct.LocalConversation, pbConversation.GetOwnerConversationResp, string]
	db                          db_interface.DataBase
	ConversationListener        func() open_im_sdk_callback.OnConversationListener
	msgListener                 func() open_im_sdk_callback.OnAdvancedMsgListener
	msgKvListener               func() open_im_sdk_callback.OnMessageKvInfoListener
	batchMsgListener            func() open_im_sdk_callback.OnBatchMsgListener
	businessListener            func() open_im_sdk_callback.OnCustomBusinessListener
	recvCH                      chan common.Cmd2Value
	loginUserID                 string
	platformID                  int32
	DataDir                     string
	relation                    *relation.Relation
	group                       *group.Group
	user                        *user.User
	file                        *file.File
	cache                       *cache.Cache[string, *model_struct.LocalConversation]
	maxSeqRecorder              MaxSeqRecorder
	messagePullForwardEndSeqMap *cache.ConversationSeqContextCache
	messagePullReverseEndSeqMap *cache.ConversationSeqContextCache
	IsExternalExtensions        bool
	msgOffset                   int
	progress                    int
	conversationSyncMutex       sync.Mutex

	startTime time.Time

	typing *typing
}

func (c *Conversation) SetMsgListener(msgListener func() open_im_sdk_callback.OnAdvancedMsgListener) {
	c.msgListener = msgListener
}

func (c *Conversation) SetMsgKvListener(msgKvListener func() open_im_sdk_callback.OnMessageKvInfoListener) {
	c.msgKvListener = msgKvListener
}

func (c *Conversation) SetBatchMsgListener(batchMsgListener func() open_im_sdk_callback.OnBatchMsgListener) {
	c.batchMsgListener = batchMsgListener
}

func (c *Conversation) SetBusinessListener(businessListener func() open_im_sdk_callback.OnCustomBusinessListener) {
	c.businessListener = businessListener
}

func NewConversation(ctx context.Context, longConnMgr *interaction.LongConnMgr, db db_interface.DataBase,
	ch chan common.Cmd2Value, relation *relation.Relation, group *group.Group, user *user.User,
	file *file.File) *Conversation {
	info := ccontext.Info(ctx)
	n := &Conversation{db: db,
		LongConnMgr:                 longConnMgr,
		recvCH:                      ch,
		loginUserID:                 info.UserID(),
		platformID:                  info.PlatformID(),
		DataDir:                     info.DataDir(),
		relation:                    relation,
		group:                       group,
		user:                        user,
		file:                        file,
		IsExternalExtensions:        info.IsExternalExtensions(),
		maxSeqRecorder:              NewMaxSeqRecorder(),
		messagePullForwardEndSeqMap: cache.NewConversationSeqContextCache(),
		messagePullReverseEndSeqMap: cache.NewConversationSeqContextCache(),
		msgOffset:                   0,
		progress:                    0,
	}
	n.typing = newTyping(n)
	n.initSyncer()
	n.cache = cache.NewCache[string, *model_struct.LocalConversation]()
	return n
}

func (c *Conversation) initSyncer() {
	c.conversationSyncer = syncer.New2[*model_struct.LocalConversation, pbConversation.GetOwnerConversationResp, string](
		syncer.WithInsert[*model_struct.LocalConversation, pbConversation.GetOwnerConversationResp, string](func(ctx context.Context, value *model_struct.LocalConversation) error {
			if err := c.batchAddFaceURLAndName(ctx, value); err != nil {
				return err
			}
			return c.db.InsertConversation(ctx, value)
		}),
		syncer.WithDelete[*model_struct.LocalConversation, pbConversation.GetOwnerConversationResp, string](func(ctx context.Context, value *model_struct.LocalConversation) error {
			return c.db.DeleteConversation(ctx, value.ConversationID)
		}),
		syncer.WithUpdate[*model_struct.LocalConversation, pbConversation.GetOwnerConversationResp, string](func(ctx context.Context, serverConversation, localConversation *model_struct.LocalConversation) error {
			return c.db.UpdateColumnsConversation(ctx, serverConversation.ConversationID,
				map[string]interface{}{"recv_msg_opt": serverConversation.RecvMsgOpt,
					"is_pinned": serverConversation.IsPinned, "is_private_chat": serverConversation.IsPrivateChat, "burn_duration": serverConversation.BurnDuration,
					"is_not_in_group": serverConversation.IsNotInGroup, "group_at_type": serverConversation.GroupAtType,
					"update_unread_count_time": serverConversation.UpdateUnreadCountTime,
					"attached_info":            serverConversation.AttachedInfo, "ex": serverConversation.Ex, "msg_destruct_time": serverConversation.MsgDestructTime,
					"is_msg_destruct": serverConversation.IsMsgDestruct,
					"max_seq":         serverConversation.MaxSeq, "min_seq": serverConversation.MinSeq})
		}),
		syncer.WithUUID[*model_struct.LocalConversation, pbConversation.GetOwnerConversationResp, string](func(value *model_struct.LocalConversation) string {
			return value.ConversationID
		}),
		syncer.WithEqual[*model_struct.LocalConversation, pbConversation.GetOwnerConversationResp, string](func(server, local *model_struct.LocalConversation) bool {
			if server.RecvMsgOpt != local.RecvMsgOpt ||
				server.IsPinned != local.IsPinned ||
				server.IsPrivateChat != local.IsPrivateChat ||
				server.BurnDuration != local.BurnDuration ||
				server.IsNotInGroup != local.IsNotInGroup ||
				server.GroupAtType != local.GroupAtType ||
				server.UpdateUnreadCountTime != local.UpdateUnreadCountTime ||
				server.AttachedInfo != local.AttachedInfo ||
				server.Ex != local.Ex ||
				server.MaxSeq != local.MaxSeq ||
				server.MinSeq != local.MinSeq ||
				server.MsgDestructTime != local.MsgDestructTime ||
				server.IsMsgDestruct != local.IsMsgDestruct {
				log.ZDebug(context.Background(), "not same", "conversationID", server.ConversationID, "server", server.RecvMsgOpt, "local", local.RecvMsgOpt)
				return false
			}
			return true
		}),
		syncer.WithNotice[*model_struct.LocalConversation, pbConversation.GetOwnerConversationResp, string](func(ctx context.Context, state int, server, local *model_struct.LocalConversation) error {
			if state == syncer.Update || state == syncer.Insert {
				c.doUpdateConversation(common.Cmd2Value{Value: common.UpdateConNode{ConID: server.ConversationID, Action: constant.ConChange, Args: []string{server.ConversationID}}})
			}
			return nil
		}),
		syncer.WithBatchInsert[*model_struct.LocalConversation, pbConversation.GetOwnerConversationResp, string](func(ctx context.Context, values []*model_struct.LocalConversation) error {
			if err := c.batchAddFaceURLAndName(ctx, values...); err != nil {
				return err
			}
			return c.db.BatchInsertConversationList(ctx, values)
		}),
		syncer.WithDeleteAll[*model_struct.LocalConversation, pbConversation.GetOwnerConversationResp, string](func(ctx context.Context, _ string) error {
			return c.db.DeleteAllConversation(ctx)
		}),
		syncer.WithBatchPageReq[*model_struct.LocalConversation, pbConversation.GetOwnerConversationResp, string](func(entityID string) page.PageReq {
			return &pbConversation.GetOwnerConversationReq{UserID: entityID,
				Pagination: &sdkws.RequestPagination{ShowNumber: 300}}
		}),
		syncer.WithBatchPageRespConvertFunc[*model_struct.LocalConversation, pbConversation.GetOwnerConversationResp, string](func(resp *pbConversation.GetOwnerConversationResp) []*model_struct.LocalConversation {
			return datautil.Batch(ServerConversationToLocal, resp.Conversations)
		}),
		syncer.WithReqApiRouter[*model_struct.LocalConversation, pbConversation.GetOwnerConversationResp, string](api.GetOwnerConversation.Route()),
		syncer.WithFullSyncLimit[*model_struct.LocalConversation, pbConversation.GetOwnerConversationResp, string](conversationSyncLimit),
	)

}

func (c *Conversation) GetCh() chan common.Cmd2Value {
	return c.recvCH
}

type onlineMsgKey struct {
	ClientMsgID string
	ServerMsgID string
}

func (c *Conversation) doMsgNew(c2v common.Cmd2Value) {
	allMsg := c2v.Value.(sdk_struct.CmdNewMsgComeToConversation).Msgs
	ctx := c2v.Ctx
	var isTriggerUnReadCount bool
	insertMsg := make(map[string][]*model_struct.LocalChatLog, 10)
	updateMsg := make(map[string][]*model_struct.LocalChatLog, 10)
	var exceptionMsg []*model_struct.LocalErrChatLog
	//var unreadMessages []*model_struct.LocalConversationUnreadMessage
	var newMessages sdk_struct.NewMsgList
	// var reactionMsgModifierList, reactionMsgDeleterList sdk_struct.NewMsgList

	var isUnreadCount, isConversationUpdate, isHistory, isNotPrivate, isSenderConversationUpdate bool

	conversationChangedSet := make(map[string]*model_struct.LocalConversation)
	newConversationSet := make(map[string]*model_struct.LocalConversation)
	conversationSet := make(map[string]*model_struct.LocalConversation)
	phConversationChangedSet := make(map[string]*model_struct.LocalConversation)
	phNewConversationSet := make(map[string]*model_struct.LocalConversation)

	log.ZDebug(ctx, "message come here conversation ch", "conversation length", len(allMsg))
	b := time.Now()

	onlineMap := make(map[onlineMsgKey]struct{})

	for conversationID, msgs := range allMsg {
		log.ZDebug(ctx, "parse message in one conversation", "conversationID",
			conversationID, "message length", len(msgs.Msgs))
		var insertMessage, selfInsertMessage, othersInsertMessage []*model_struct.LocalChatLog
		var updateMessage []*model_struct.LocalChatLog

		for _, v := range msgs.Msgs {
			log.ZDebug(ctx, "parse message ", "conversationID", conversationID, "msg", v)
			isHistory = utils.GetSwitchFromOptions(v.Options, constant.IsHistory)

			isUnreadCount = utils.GetSwitchFromOptions(v.Options, constant.IsUnreadCount)

			isConversationUpdate = utils.GetSwitchFromOptions(v.Options, constant.IsConversationUpdate)

			isNotPrivate = utils.GetSwitchFromOptions(v.Options, constant.IsNotPrivate)

			isSenderConversationUpdate = utils.GetSwitchFromOptions(v.Options, constant.IsSenderConversationUpdate)

			msg := &sdk_struct.MsgStruct{}
			copier.Copy(msg, v)
			msg.Content = string(v.Content)

			var attachedInfo sdk_struct.AttachedInfoElem
			_ = utils.JsonStringToStruct(v.AttachedInfo, &attachedInfo)
			msg.AttachedInfoElem = &attachedInfo

			//When the message has been marked and deleted by the cloud, it is directly inserted locally without any conversation and message update.
			if msg.Status == constant.MsgStatusHasDeleted {
				insertMessage = append(insertMessage, MsgStructToLocalChatLog(msg))
				continue
			}

			msg.Status = constant.MsgStatusSendSuccess

			//De-analyze data
			err := msgHandleByContentType(msg)
			if err != nil {
				log.ZError(ctx, "Parsing data error:", err, "type: ", msg.ContentType, "msg", msg)
				continue
			}

			if !isNotPrivate {
				msg.AttachedInfoElem.IsPrivateChat = true
			}
			if msg.ClientMsgID == "" {
				exceptionMsg = append(exceptionMsg, c.msgStructToLocalErrChatLog(msg))
				continue
			}
			if conversationID == "" {
				log.ZError(ctx, "conversationID is empty", errors.New("conversationID is empty"), "msg", msg)
				continue
			}
			if !isHistory {
				onlineMap[onlineMsgKey{ClientMsgID: v.ClientMsgID, ServerMsgID: v.ServerMsgID}] = struct{}{}
				newMessages = append(newMessages, msg)

			}
			log.ZDebug(ctx, "decode message", "msg", msg)
			if v.SendID == c.loginUserID { //seq
				// Messages sent by myself  //if  sent through  this terminal
				m, err := c.db.GetMessage(ctx, conversationID, msg.ClientMsgID)
				if err == nil {
					log.ZInfo(ctx, "have message", "msg", msg)
					if m.Seq == 0 {
						if !isConversationUpdate {
							msg.Status = constant.MsgStatusFiltered
						}
						updateMessage = append(updateMessage, MsgStructToLocalChatLog(msg))
					} else {
						exceptionMsg = append(exceptionMsg, c.msgStructToLocalErrChatLog(msg))
					}
				} else {
					log.ZInfo(ctx, "sync message", "msg", msg)
					lc := model_struct.LocalConversation{
						ConversationType:  v.SessionType,
						LatestMsg:         utils.StructToJsonString(msg),
						LatestMsgSendTime: msg.SendTime,
						ConversationID:    conversationID,
					}
					switch v.SessionType {
					case constant.SingleChatType:
						lc.UserID = v.RecvID
					case constant.WriteGroupChatType, constant.ReadGroupChatType:
						lc.GroupID = v.GroupID
					}
					if isConversationUpdate {
						if isSenderConversationUpdate {
							log.ZDebug(ctx, "updateConversation msg", "message", v, "conversation", lc)
							c.updateConversation(&lc, conversationSet)
						}
						newMessages = append(newMessages, msg)
					}
					if isHistory {
						selfInsertMessage = append(selfInsertMessage, MsgStructToLocalChatLog(msg))
					}
				}
			} else { //Sent by others
				if _, err := c.db.GetMessage(ctx, conversationID, msg.ClientMsgID); err != nil { //Deduplication operation
					lc := model_struct.LocalConversation{
						ConversationType:  v.SessionType,
						LatestMsg:         utils.StructToJsonString(msg),
						LatestMsgSendTime: msg.SendTime,
						ConversationID:    conversationID,
					}
					switch v.SessionType {
					case constant.SingleChatType:
						lc.UserID = v.SendID
						lc.ShowName = msg.SenderNickname
						lc.FaceURL = msg.SenderFaceURL
					case constant.WriteGroupChatType, constant.ReadGroupChatType:
						lc.GroupID = v.GroupID
					case constant.NotificationChatType:
						lc.UserID = v.SendID
					}
					if isUnreadCount {
						//cacheConversation := c.cache.GetConversation(lc.ConversationID)
						if c.maxSeqRecorder.IsNewMsg(conversationID, msg.Seq) {
							isTriggerUnReadCount = true
							lc.UnreadCount = 1
							c.maxSeqRecorder.Incr(conversationID, 1)
						}
					}
					if isConversationUpdate {
						c.updateConversation(&lc, conversationSet)
						newMessages = append(newMessages, msg)
					}
					if isHistory {
						othersInsertMessage = append(othersInsertMessage, MsgStructToLocalChatLog(msg))
					}

				} else {
					exceptionMsg = append(exceptionMsg, c.msgStructToLocalErrChatLog(msg))
					log.ZWarn(ctx, "Deduplication operation ", nil, "msg", *c.msgStructToLocalErrChatLog(msg))
					msg.Status = constant.MsgStatusFiltered
					msg.ClientMsgID = msg.ClientMsgID + utils.Int64ToString(msg.Seq)
					othersInsertMessage = append(othersInsertMessage, MsgStructToLocalChatLog(msg))
				}
			}
		}
		insertMsg[conversationID] = append(insertMessage, c.faceURLAndNicknameHandle(ctx, selfInsertMessage, othersInsertMessage, conversationID)...)
		if len(updateMessage) > 0 {
			updateMsg[conversationID] = updateMessage

		}
	}

	//todo The lock granularity needs to be optimized to the conversation level.
	c.conversationSyncMutex.Lock()
	defer c.conversationSyncMutex.Unlock()

	list, err := c.db.GetAllConversationListDB(ctx)
	if err != nil {
		log.ZError(ctx, "GetAllConversationListDB", err)
	}

	m := make(map[string]*model_struct.LocalConversation)
	listToMap(list, m)
	log.ZDebug(ctx, "listToMap: ", "local conversation", list, "generated c map",
		string(stringutil.StructToJsonBytes(conversationSet)))

	c.diff(ctx, m, conversationSet, conversationChangedSet, newConversationSet)
	log.ZInfo(ctx, "trigger map is :", "newConversations", string(stringutil.StructToJsonBytes(newConversationSet)),
		"changedConversations", string(stringutil.StructToJsonBytes(conversationChangedSet)))

	//seq sync message update
	if err := c.batchUpdateMessageList(ctx, updateMsg); err != nil {
		log.ZError(ctx, "sync seq normal message err  :", err)
	}

	//Normal message storage
	_ = c.batchInsertMessageList(ctx, insertMsg)

	hList, _ := c.db.GetHiddenConversationList(ctx)
	for _, v := range hList {
		if nc, ok := newConversationSet[v.ConversationID]; ok {
			phConversationChangedSet[v.ConversationID] = nc
			nc.RecvMsgOpt = v.RecvMsgOpt
			nc.GroupAtType = v.GroupAtType
			nc.IsPinned = v.IsPinned
			nc.IsPrivateChat = v.IsPrivateChat
			if nc.IsPrivateChat {
				nc.BurnDuration = v.BurnDuration
			}
			if v.UnreadCount != 0 {
				nc.UnreadCount = v.UnreadCount
			}
			nc.IsNotInGroup = v.IsNotInGroup
			nc.AttachedInfo = v.AttachedInfo
			nc.Ex = v.Ex
			nc.IsMsgDestruct = v.IsMsgDestruct
			nc.MsgDestructTime = v.MsgDestructTime
		}
	}

	for k, v := range newConversationSet {
		if _, ok := phConversationChangedSet[v.ConversationID]; !ok {
			phNewConversationSet[k] = v
		}
	}

	if err := c.db.BatchUpdateConversationList(ctx, append(mapConversationToList(conversationChangedSet), mapConversationToList(phConversationChangedSet)...)); err != nil {
		log.ZError(ctx, "insert changed conversation err :", err)
	}
	//New conversation storage

	if err := c.db.BatchInsertConversationList(ctx, mapConversationToList(phNewConversationSet)); err != nil {
		log.ZError(ctx, "insert new conversation err:", err)
	}
	log.ZDebug(ctx, "before trigger msg", "cost time", time.Since(b).Seconds(), "len", len(allMsg))

	if c.batchMsgListener() != nil {
		c.batchNewMessages(ctx, newMessages, conversationChangedSet, newConversationSet, onlineMap)
	} else {
		c.newMessage(ctx, newMessages, conversationChangedSet, newConversationSet, onlineMap)
	}
	if len(newConversationSet) > 0 {
		c.doUpdateConversation(common.Cmd2Value{Value: common.UpdateConNode{Action: constant.NewConDirect, Args: utils.StructToJsonString(mapConversationToList(newConversationSet))}})
	}
	if len(conversationChangedSet) > 0 {
		c.doUpdateConversation(common.Cmd2Value{Value: common.UpdateConNode{Action: constant.ConChangeDirect, Args: utils.StructToJsonString(mapConversationToList(conversationChangedSet))}})
	}

	if isTriggerUnReadCount {
		c.doUpdateConversation(common.Cmd2Value{Value: common.UpdateConNode{Action: constant.TotalUnreadMessageChanged, Args: ""}})
	}

	for _, msgs := range allMsg {
		for _, msg := range msgs.Msgs {
			if msg.ContentType == constant.Typing {
				c.typing.onNewMsg(ctx, msg)
			}
		}
	}

	log.ZDebug(ctx, "insert msg", "duration", fmt.Sprintf("%dms", time.Since(b)), "len", len(allMsg))
}

func (c *Conversation) doMsgSyncByReinstalled(c2v common.Cmd2Value) {
	allMsg := c2v.Value.(sdk_struct.CmdMsgSyncInReinstall).Msgs
	ctx := c2v.Ctx
	msgLen := len(allMsg)
	c.msgOffset += msgLen
	total := c2v.Value.(sdk_struct.CmdMsgSyncInReinstall).Total

	insertMsg := make(map[string][]*model_struct.LocalChatLog, 10)
	conversationList := make([]*model_struct.LocalConversation, 0)

	log.ZDebug(ctx, "message come here conversation ch in reinstalled", "conversation length", msgLen)
	b := time.Now()

	for conversationID, msgs := range allMsg {
		log.ZDebug(ctx, "parse message in one conversation", "conversationID",
			conversationID, "message length", len(msgs.Msgs))
		var insertMessage, selfInsertMessage, othersInsertMessage []*model_struct.LocalChatLog
		var latestMsg *sdk_struct.MsgStruct
		if len(msgs.Msgs) == 0 {
			log.ZWarn(ctx, "msg.Msgs is empty", errs.New("msg.Msgs is empty"), "conversationID", conversationID)
			continue
		}
		for _, v := range msgs.Msgs {

			log.ZDebug(ctx, "parse message ", "conversationID", conversationID, "msg", v)
			msg := &sdk_struct.MsgStruct{}
			// TODO need replace when after.
			copier.Copy(msg, v)
			msg.Content = string(v.Content)
			var attachedInfo sdk_struct.AttachedInfoElem
			_ = utils.JsonStringToStruct(v.AttachedInfo, &attachedInfo)
			msg.AttachedInfoElem = &attachedInfo

			//When the message has been marked and deleted by the cloud, it is directly inserted locally without any conversation and message update.
			if msg.Status == constant.MsgStatusHasDeleted {
				insertMessage = append(insertMessage, MsgStructToLocalChatLog(msg))
				continue
			}
			msg.Status = constant.MsgStatusSendSuccess

			err := msgHandleByContentType(msg)
			if err != nil {
				log.ZError(ctx, "Parsing data error:", err, "type: ", msg.ContentType, "msg", msg)
				continue
			}

			if conversationID == "" {
				log.ZError(ctx, "conversationID is empty", errors.New("conversationID is empty"), "msg", msg)
				continue
			}

			log.ZDebug(ctx, "decode message", "msg", msg)
			if v.SendID == c.loginUserID {
				// Messages sent by myself  //if  sent through  this terminal
				log.ZInfo(ctx, "sync message in reinstalled", "msg", msg)

				latestMsg = msg

				selfInsertMessage = append(selfInsertMessage, MsgStructToLocalChatLog(msg))
			} else { //Sent by others
				othersInsertMessage = append(othersInsertMessage, MsgStructToLocalChatLog(msg))

				latestMsg = msg
			}
		}

		if latestMsg != nil {
			conversationList = append(conversationList, &model_struct.LocalConversation{
				LatestMsg:         utils.StructToJsonString(latestMsg),
				LatestMsgSendTime: latestMsg.SendTime,
				ConversationID:    conversationID,
			})
		} else {
			log.ZWarn(ctx, "latestMsg is nil", errs.New("latestMsg is nil"), "conversationID", conversationID)
		}

		insertMsg[conversationID] = append(insertMessage, c.faceURLAndNicknameHandle(ctx, selfInsertMessage, othersInsertMessage, conversationID)...)
	}

	// message storage
	_ = c.batchInsertMessageList(ctx, insertMsg)

	// conversation storage
	if err := c.db.BatchUpdateConversationList(ctx, conversationList); err != nil {
		log.ZError(ctx, "insert new conversation err:", err)
	}
	log.ZDebug(ctx, "before trigger msg", "cost time", time.Since(b).Seconds(), "len", len(allMsg))

	// log.ZDebug(ctx, "progress is", "msgLen", msgLen, "msgOffset", c.msgOffset, "total", total, "now progress is", (c.msgOffset*(100-InitSyncProgress))/total + InitSyncProgress)
	c.ConversationListener().OnSyncServerProgress((c.msgOffset*(100-InitSyncProgress))/total + InitSyncProgress)
}

func (c *Conversation) addInitProgress(progress int) {
	c.progress += progress
	if c.progress > 100 {
		c.progress = 100
	}
}

func listToMap(list []*model_struct.LocalConversation, m map[string]*model_struct.LocalConversation) {
	for _, v := range list {
		m[v.ConversationID] = v
	}
}

func (c *Conversation) diff(ctx context.Context, local, generated, cc, nc map[string]*model_struct.LocalConversation) {
	var newConversations []*model_struct.LocalConversation
	for _, v := range generated {
		if localC, ok := local[v.ConversationID]; ok {

			if v.LatestMsgSendTime > localC.LatestMsgSendTime {
				localC.UnreadCount = localC.UnreadCount + v.UnreadCount
				localC.LatestMsg = v.LatestMsg
				localC.LatestMsgSendTime = v.LatestMsgSendTime
				cc[v.ConversationID] = localC
			} else {
				localC.UnreadCount = localC.UnreadCount + v.UnreadCount
				cc[v.ConversationID] = localC
			}

		} else {
			newConversations = append(newConversations, v)
		}
	}
	if err := c.batchAddFaceURLAndName(ctx, newConversations...); err != nil {
		log.ZError(ctx, "batchAddFaceURLAndName err", err, "conversations", newConversations)
	} else {
		for _, v := range newConversations {
			nc[v.ConversationID] = v
		}
	}
}

func (c *Conversation) msgStructToLocalErrChatLog(m *sdk_struct.MsgStruct) *model_struct.LocalErrChatLog {
	var lc model_struct.LocalErrChatLog
	copier.Copy(&lc, m)
	if m.SessionType == constant.WriteGroupChatType || m.SessionType == constant.ReadGroupChatType {
		lc.RecvID = m.GroupID
	}
	return &lc
}

func (c *Conversation) batchUpdateMessageList(ctx context.Context, updateMsg map[string][]*model_struct.LocalChatLog) error {
	if updateMsg == nil {
		return nil
	}
	for conversationID, messages := range updateMsg {
		conversation, err := c.db.GetConversation(ctx, conversationID)
		if err != nil {
			log.ZError(ctx, "GetConversation err", err, "conversationID", conversationID)
			continue
		}
		latestMsg := &sdk_struct.MsgStruct{}
		if err := json.Unmarshal([]byte(conversation.LatestMsg), latestMsg); err != nil {
			log.ZError(ctx, "Unmarshal err", err, "conversationID",
				conversationID, "latestMsg", conversation.LatestMsg, "messages", messages)
			continue
		}
		for _, v := range messages {
			v1 := new(model_struct.LocalChatLog)
			v1.ClientMsgID = v.ClientMsgID
			v1.Seq = v.Seq
			v1.Status = v.Status
			v1.RecvID = v.RecvID
			v1.SessionType = v.SessionType
			v1.ServerMsgID = v.ServerMsgID
			v1.SendTime = v.SendTime
			err := c.db.UpdateMessage(ctx, conversationID, v1)
			if err != nil {
				return errs.WrapMsg(err, "BatchUpdateMessageList failed")
			}
			if latestMsg.ClientMsgID == v.ClientMsgID {
				latestMsg.ServerMsgID = v.ServerMsgID
				latestMsg.Seq = v.Seq
				latestMsg.SendTime = v.SendTime
				latestMsg.Status = v.Status
				conversation.LatestMsg = utils.StructToJsonString(latestMsg)

				c.doUpdateConversation(common.Cmd2Value{Value: common.UpdateConNode{ConID: conversation.ConversationID,
					Action: constant.AddConOrUpLatMsg, Args: *conversation}})

			}
		}

	}
	return nil
}

func (c *Conversation) batchInsertMessageList(ctx context.Context, insertMsg map[string][]*model_struct.LocalChatLog) error {
	if insertMsg == nil {
		return nil
	}
	for conversationID, messages := range insertMsg {
		if len(messages) == 0 {
			continue
		}
		err := c.db.BatchInsertMessageList(ctx, conversationID, messages)
		if err != nil {
			log.ZError(ctx, "insert GetMessage detail err:", err, "conversationID", conversationID, "messages", messages)
			for _, v := range messages {
				e := c.db.InsertMessage(ctx, conversationID, v)
				if e != nil {
					log.ZError(ctx, "InsertMessage err", err, "conversationID", conversationID, "message", v)
				}
			}
		}

	}
	return nil
}

func (c *Conversation) newMessage(ctx context.Context, newMessagesList sdk_struct.NewMsgList, cc, nc map[string]*model_struct.LocalConversation, onlineMsg map[onlineMsgKey]struct{}) {
	sort.Sort(newMessagesList)
	if c.GetBackground() {
		u, err := c.user.GetSelfUserInfo(ctx)
		if err != nil {
			log.ZWarn(ctx, "GetSelfUserInfo err", err)
			return
		}
		if u.GlobalRecvMsgOpt != constant.ReceiveMessage {
			return
		}
		for _, w := range newMessagesList {
			conversationID := utils.GetConversationIDByMsg(w)
			if v, ok := cc[conversationID]; ok && v.RecvMsgOpt == constant.ReceiveMessage {
				c.msgListener().OnRecvOfflineNewMessage(utils.StructToJsonString(w))
			}
			if v, ok := nc[conversationID]; ok && v.RecvMsgOpt == constant.ReceiveMessage {
				c.msgListener().OnRecvOfflineNewMessage(utils.StructToJsonString(w))
			}
		}
	} else {
		for _, w := range newMessagesList {
			if w.ContentType == constant.Typing {
				continue
			}
			if _, ok := onlineMsg[onlineMsgKey{ClientMsgID: w.ClientMsgID, ServerMsgID: w.ServerMsgID}]; ok {
				c.msgListener().OnRecvOnlineOnlyMessage(utils.StructToJsonString(w))
			} else {
				c.msgListener().OnRecvNewMessage(utils.StructToJsonString(w))
			}
		}
	}
}

func (c *Conversation) batchNewMessages(ctx context.Context, newMessagesList sdk_struct.NewMsgList, conversationChanged, newConversation map[string]*model_struct.LocalConversation, onlineMsg map[onlineMsgKey]struct{}) {
	if len(newMessagesList) == 0 {
		log.ZWarn(ctx, "newMessagesList is empty", errs.New("newMessagesList is empty"))
		return
	}

	sort.Sort(newMessagesList)
	var needNotificationMsgList sdk_struct.NewMsgList

	// offline
	if c.GetBackground() {
		u, err := c.user.GetSelfUserInfo(ctx)
		if err != nil {
			log.ZWarn(ctx, "GetSelfUserInfo err", err)
		}

		if u.GlobalRecvMsgOpt != constant.ReceiveMessage {
			return
		}

		for _, w := range newMessagesList {
			conversationID := utils.GetConversationIDByMsg(w)
			if v, ok := conversationChanged[conversationID]; ok && v.RecvMsgOpt == constant.ReceiveMessage {
				needNotificationMsgList = append(needNotificationMsgList, w)
			}
			if v, ok := newConversation[conversationID]; ok && v.RecvMsgOpt == constant.ReceiveMessage {
				needNotificationMsgList = append(needNotificationMsgList, w)
			}
		}

		if len(needNotificationMsgList) != 0 {
			c.batchMsgListener().OnRecvOfflineNewMessages(utils.StructToJsonString(needNotificationMsgList))
		}
	} else { // online
		for _, w := range newMessagesList {
			if w.ContentType == constant.Typing {
				continue
			}

			needNotificationMsgList = append(needNotificationMsgList, w)
		}

		if len(needNotificationMsgList) != 0 {
			c.batchMsgListener().OnRecvNewMessages(utils.StructToJsonString(needNotificationMsgList))
		}
	}
}

func (c *Conversation) updateConversation(lc *model_struct.LocalConversation, cs map[string]*model_struct.LocalConversation) {
	if oldC, ok := cs[lc.ConversationID]; !ok {
		cs[lc.ConversationID] = lc
	} else {
		if lc.LatestMsgSendTime > oldC.LatestMsgSendTime {
			oldC.UnreadCount = oldC.UnreadCount + lc.UnreadCount
			oldC.LatestMsg = lc.LatestMsg
			oldC.LatestMsgSendTime = lc.LatestMsgSendTime
			cs[lc.ConversationID] = oldC
		} else {
			oldC.UnreadCount = oldC.UnreadCount + lc.UnreadCount
			cs[lc.ConversationID] = oldC
		}
	}
}

func mapConversationToList(m map[string]*model_struct.LocalConversation) (cs []*model_struct.LocalConversation) {
	for _, v := range m {
		cs = append(cs, v)
	}
	return cs
}

func (c *Conversation) batchAddFaceURLAndName(ctx context.Context, conversations ...*model_struct.LocalConversation) error {
	if len(conversations) == 0 {
		return nil
	}
	var userIDs, groupIDs []string
	for _, conversation := range conversations {
		if conversation.ConversationType == constant.SingleChatType ||
			conversation.ConversationType == constant.NotificationChatType {
			userIDs = append(userIDs, conversation.UserID)
		} else if conversation.ConversationType == constant.ReadGroupChatType {
			groupIDs = append(groupIDs, conversation.GroupID)
		}
	}

	// if userIDs = nil, return nil, nil
	users, err := c.batchGetUserNameAndFaceURL(ctx, userIDs...)
	if err != nil {
		return err
	}

	groupInfoList, err := c.group.GetSpecifiedGroupsInfo(ctx, groupIDs)
	if err != nil {
		return err
	}

	groups := datautil.SliceToMap(groupInfoList, func(groupInfo *model_struct.LocalGroup) string {
		return groupInfo.GroupID
	})

	for _, conversation := range conversations {
		if conversation.ConversationType == constant.SingleChatType ||
			conversation.ConversationType == constant.NotificationChatType {
			if v, ok := users[conversation.UserID]; ok {
				conversation.FaceURL = v.FaceURL
				conversation.ShowName = v.Nickname
			} else {
				log.ZWarn(ctx, "user info not found", errors.New("user not found"), "userID", conversation.UserID)

				conversation.FaceURL = ""
				conversation.ShowName = "UserNotFound"
			}
		} else if conversation.ConversationType == constant.ReadGroupChatType {
			if v, ok := groups[conversation.GroupID]; ok {
				conversation.FaceURL = v.FaceURL
				conversation.ShowName = v.GroupName
			} else {
				log.ZWarn(ctx, "group info not found", errors.New("group not found"),
					"groupID", conversation.GroupID)
			}

		}
	}

	return nil
}

func (c *Conversation) batchGetUserNameAndFaceURL(ctx context.Context, userIDs ...string) (map[string]*model_struct.LocalUser, error) {
	m := make(map[string]*model_struct.LocalUser)
	var notInFriend []string

	if len(userIDs) == 0 {
		return m, nil
	}

	friendList, err := c.relation.Db().GetFriendInfoList(ctx, userIDs)
	if err != nil {
		log.ZWarn(ctx, "BatchGetUserNameAndFaceURL", err, "userIDs", userIDs)
		notInFriend = userIDs
	} else {
		notInFriend = datautil.SliceSub(userIDs, datautil.Slice(friendList, func(e *model_struct.LocalFriend) string {
			return e.FriendUserID
		}))
	}
	for _, localFriend := range friendList {
		userInfo := &model_struct.LocalUser{UserID: localFriend.FriendUserID, FaceURL: localFriend.FaceURL}
		if localFriend.Remark != "" {
			userInfo.Nickname = localFriend.Remark
		} else {
			userInfo.Nickname = localFriend.Nickname
		}
		m[localFriend.FriendUserID] = userInfo
	}

	usersInfo, err := c.user.GetUsersInfoWithCache(ctx, notInFriend)
	if err != nil {
		return nil, err
	}

	for _, userInfo := range usersInfo {
		m[userInfo.UserID] = userInfo
	}
	return m, nil
}

func (c *Conversation) getUserNameAndFaceURL(ctx context.Context, userID string) (faceURL, name string, err error) {
	friendInfo, err := c.relation.Db().GetFriendInfoByFriendUserID(ctx, userID)
	if err == nil {
		faceURL = friendInfo.FaceURL
		if friendInfo.Remark != "" {
			name = friendInfo.Remark
		} else {
			name = friendInfo.Nickname
		}
		return faceURL, name, nil
	}
	userInfo, err := c.user.GetUserInfoWithCache(ctx, userID)
	if err != nil {
		return "", "", nil
	}
	return userInfo.FaceURL, userInfo.Nickname, nil
}
