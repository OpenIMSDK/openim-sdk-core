package conversation_msg

import (
	"context"
	"errors"
	"sort"
	"sync"
	"time"

	"github.com/jinzhu/copier"
	"golang.org/x/sync/errgroup"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/api"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/cache"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	sdk "github.com/openimsdk/openim-sdk-core/v3/pkg/sdk_params_callback"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/sdkerrs"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/utils"
	"github.com/openimsdk/openim-sdk-core/v3/sdk_struct"

	"github.com/openimsdk/tools/log"

	"github.com/openimsdk/protocol/sdkws"

	pbConversation "github.com/openimsdk/protocol/conversation"
)

const MaxRecursionDepth = 3

func (c *Conversation) setConversation(ctx context.Context, apiReq *pbConversation.SetConversationsReq, localConversation *model_struct.LocalConversation) error {
	apiReq.Conversation.ConversationID = localConversation.ConversationID
	apiReq.Conversation.ConversationType = localConversation.ConversationType
	apiReq.Conversation.UserID = localConversation.UserID
	apiReq.Conversation.GroupID = localConversation.GroupID
	apiReq.UserIDs = []string{c.loginUserID}

	return api.SetConversations.Execute(ctx, apiReq)
}

func (c *Conversation) getAdvancedHistoryMessageList(ctx context.Context, req sdk.GetAdvancedHistoryMessageListParams, isReverse bool) (*sdk.GetAdvancedHistoryMessageListCallback, error) {
	t := time.Now()
	var messageListCallback sdk.GetAdvancedHistoryMessageListCallback
	var conversationID string
	var startTime int64
	var err error
	var messageList sdk_struct.NewMsgList
	conversationID = req.ConversationID
	if len(req.StartClientMsgID) > 0 {
		m, err := c.db.GetMessage(ctx, conversationID, req.StartClientMsgID)
		if err != nil {
			return nil, err
		}
		startTime = m.SendTime
		err = c.handleEndSeq(ctx, req, isReverse, m)
		if err != nil {
			return nil, err
		}
	} else {
		// Clear both maps when the user enters the conversation
		c.messagePullForwardEndSeqMap.Delete(conversationID, req.ViewType)
		c.messagePullReverseEndSeqMap.Delete(conversationID, req.ViewType)
	}

	log.ZDebug(ctx, "Assembly conversation parameters", "cost time", time.Since(t), "conversationID",
		conversationID, "startTime:", startTime, "count:", req.Count, "startTime", startTime)
	list, err := c.fetchMessagesWithGapCheck(ctx, conversationID, req.Count, startTime, isReverse, req.ViewType, &messageListCallback)
	if err != nil {
		return nil, err
	}
	log.ZDebug(ctx, "pull message", "pull cost time", time.Since(t).Milliseconds())
	t = time.Now()

	messageList = c.LocalChatLog2MsgStruct(list)
	log.ZDebug(ctx, "message convert and unmarshal", "unmarshal cost time", time.Since(t))
	t = time.Now()
	if !isReverse {
		sort.Sort(messageList)
	}
	log.ZDebug(ctx, "sort", "sort cost time", time.Since(t))
	messageListCallback.MessageList = messageList
	return &messageListCallback, nil
}
func (c *Conversation) handleEndSeq(ctx context.Context, req sdk.GetAdvancedHistoryMessageListParams, isReverse bool, startMessage *model_struct.LocalChatLog) error {
	if isReverse {
		if _, ok := c.messagePullReverseEndSeqMap.Load(req.ConversationID, req.ViewType); !ok {
			if startMessage.Seq != 0 {
				c.messagePullReverseEndSeqMap.Store(req.ConversationID, req.ViewType, startMessage.Seq)
			} else {
				validServerMessage, err := c.db.GetLatestValidServerMessage(ctx, req.ConversationID, startMessage.SendTime, isReverse)
				if err != nil {
					return err
				}
				if validServerMessage != nil {
					c.messagePullReverseEndSeqMap.Store(req.ConversationID, req.ViewType, validServerMessage.Seq)
				} else {
					log.ZDebug(ctx, "no valid server message", "conversationID", req.ConversationID, "startTime", startMessage.SendTime)
				}
			}
		}
	} else {
		if _, ok := c.messagePullForwardEndSeqMap.Load(req.ConversationID, req.ViewType); !ok {
			if startMessage.Seq != 0 {
				c.messagePullForwardEndSeqMap.Store(req.ConversationID, req.ViewType, startMessage.Seq)
			} else {
				validServerMessage, err := c.db.GetLatestValidServerMessage(ctx, req.ConversationID, startMessage.SendTime, isReverse)
				if err != nil {
					return err
				}
				if validServerMessage != nil {
					c.messagePullForwardEndSeqMap.Store(req.ConversationID, req.ViewType, validServerMessage.Seq)
				} else {
					log.ZDebug(ctx, "no valid server message", "conversationID", req.ConversationID, "startTime", startMessage.SendTime)
				}
			}
		}
	}
	return nil
}

func (c *Conversation) fetchMessagesWithGapCheck(ctx context.Context, conversationID string,
	count int, startTime int64, isReverse bool, viewType int, messageListCallback *sdk.GetAdvancedHistoryMessageListCallback) ([]*model_struct.LocalChatLog, error) {

	var list, validMessages []*model_struct.LocalChatLog

	// Get the number of invalid messages in this batch to recursive fetching from earlier points.
	shouldFetchMoreMessagesNum := func(messages []*model_struct.LocalChatLog) int {
		var thisEndSeq int64
		// Represents the number of valid messages in the batch
		validateMessageNum := 0
		for _, msg := range messages {
			if msg.Seq != 0 && thisEndSeq == 0 {
				thisEndSeq = msg.Seq
			}
			if isReverse {
				if msg.Seq > thisEndSeq && thisEndSeq != 0 {
					thisEndSeq = msg.Seq
				}

			} else {
				if msg.Seq < thisEndSeq && msg.Seq != 0 {
					thisEndSeq = msg.Seq
				}
			}
			if msg.Status >= constant.MsgStatusHasDeleted {
				log.ZDebug(ctx, "this message has been deleted or exception message", "msg", msg)
				continue
			}

			validateMessageNum++
			validMessages = append(validMessages, msg)

		}
		if !isReverse {
			if thisEndSeq != 0 {
				c.messagePullForwardEndSeqMap.StoreWithFunc(conversationID, viewType, thisEndSeq, func(_ string, value int64) bool {
					lastEndSeq, _ := c.messagePullForwardEndSeqMap.Load(conversationID, viewType)
					if value < lastEndSeq || lastEndSeq == 0 {
						log.ZDebug(ctx, "update the end sequence of the message", "lastEndSeq", lastEndSeq, "thisEndSeq", value)
						return true
					}
					log.ZWarn(ctx, "The end sequence number of the message is more than the last end sequence number",
						nil, "conversationID", conversationID, "value", value, "lastEndSeq", lastEndSeq)
					return false
				})
			}
		} else {
			if thisEndSeq != 0 {
				c.messagePullReverseEndSeqMap.StoreWithFunc(conversationID, viewType, thisEndSeq, func(_ string, value int64) bool {
					lastEndSeq, _ := c.messagePullReverseEndSeqMap.Load(conversationID, viewType)
					if value > lastEndSeq || lastEndSeq == 0 {
						log.ZDebug(ctx, "update the end sequence of the message", "lastEndSeq", lastEndSeq, "thisEndSeq", value)
						return true
					}
					log.ZWarn(ctx, "The end sequence number of the message is less than the last end sequence number",
						nil, "conversationID", conversationID, "value", value, "lastEndSeq", lastEndSeq)
					return false
				})
			}
		}

		return count - validateMessageNum
	}
	getNewStartTime := func(messages []*model_struct.LocalChatLog) int64 {
		if len(messages) == 0 {
			return 0
		}
		// Returns the SendTime of the last element in the message list
		return messages[len(messages)-1].SendTime
	}

	t := time.Now()
	list, err := c.db.GetMessageList(ctx, conversationID, count, startTime, isReverse)
	log.ZDebug(ctx, "db get messageList", "cost time", time.Since(t), "len", len(list), "err",
		err, "conversationID", conversationID)

	if err != nil {
		return nil, err
	}
	t = time.Now()
	thisStartSeq := c.validateAndFillInternalGaps(ctx, conversationID, isReverse,
		count, startTime, &list, messageListCallback)
	log.ZDebug(ctx, "internal continuity check", "cost time", time.Since(t), "thisStartSeq", thisStartSeq)
	t = time.Now()
	c.validateAndFillInterBlockGaps(ctx, thisStartSeq, conversationID,
		isReverse, viewType, count, startTime, &list, messageListCallback)
	log.ZDebug(ctx, "between continuity check", "cost time", time.Since(t), "thisStartSeq", thisStartSeq)
	t = time.Now()
	c.validateAndFillEndBlockContinuity(ctx, conversationID, isReverse, viewType,
		count, startTime, &list, messageListCallback)
	log.ZDebug(ctx, "end continuity check", "cost time", time.Since(t))
	// If the number of valid messages retrieved is less than the count,
	// continue fetching recursively until the valid messages are sufficient or all messages have been fetched.
	missingCount := shouldFetchMoreMessagesNum(list)
	if missingCount > 0 && !messageListCallback.IsEnd {
		newStartTime := getNewStartTime(list)
		log.ZDebug(ctx, "fetch more messages", "missingCount", missingCount, "conversationID",
			conversationID, "newStartTime", newStartTime)
		missingMessages, err := c.fetchMessagesWithGapCheck(ctx, conversationID, missingCount, newStartTime, isReverse, viewType, messageListCallback)
		if err != nil {
			return nil, err
		}
		log.ZDebug(ctx, "fetch more messages", "missingMessages", missingMessages)
		return append(validMessages, missingMessages...), nil
	}

	return validMessages, nil
}

func (c *Conversation) LocalChatLog2MsgStruct(list []*model_struct.LocalChatLog) []*sdk_struct.MsgStruct {
	messageList := make([]*sdk_struct.MsgStruct, 0, len(list))
	for _, v := range list {
		temp := LocalChatLogToMsgStruct(v)

		if temp.AttachedInfoElem.IsPrivateChat && temp.SendTime+int64(temp.AttachedInfoElem.BurnDuration) < time.Now().Unix() {
			continue
		}
		messageList = append(messageList, temp)
	}
	return messageList
}

func (c *Conversation) typingStatusUpdate(ctx context.Context, recvID, msgTip string) error {
	if recvID == "" {
		return sdkerrs.ErrArgs
	}
	s := sdk_struct.MsgStruct{}
	err := c.initBasicInfo(ctx, &s, constant.UserMsgType, constant.Typing)
	if err != nil {
		return err
	}
	s.RecvID = recvID
	s.SessionType = constant.SingleChatType
	typingElem := sdk_struct.TypingElem{}
	typingElem.MsgTips = msgTip
	s.Content = utils.StructToJsonString(typingElem)
	options := make(map[string]bool, 6)
	utils.SetSwitchFromOptions(options, constant.IsHistory, false)
	utils.SetSwitchFromOptions(options, constant.IsPersistent, false)
	utils.SetSwitchFromOptions(options, constant.IsSenderSync, false)
	utils.SetSwitchFromOptions(options, constant.IsConversationUpdate, false)
	utils.SetSwitchFromOptions(options, constant.IsSenderConversationUpdate, false)
	utils.SetSwitchFromOptions(options, constant.IsUnreadCount, false)
	utils.SetSwitchFromOptions(options, constant.IsOfflinePush, false)

	var wsMsgData sdkws.MsgData
	copier.Copy(&wsMsgData, s)
	wsMsgData.Content = []byte(s.Content)
	wsMsgData.CreateTime = s.CreateTime
	wsMsgData.Options = options
	var sendMsgResp sdkws.UserSendMsgResp
	err = c.LongConnMgr.SendReqWaitResp(ctx, &wsMsgData, constant.SendMsg, &sendMsgResp)
	if err != nil {
		log.ZError(ctx, "send msg to server failed", err, "message", s)
		return err
	}
	return nil
}

func (c *Conversation) insertMessageToLocalStorage(ctx context.Context, conversationID string, s *model_struct.LocalChatLog) error {
	return c.db.InsertMessage(ctx, conversationID, s)
}

func (c *Conversation) judgeMultipleSubString(keywordList []string, main string, keywordListMatchType int) bool {
	if len(keywordList) == 0 {
		return true
	}
	if keywordListMatchType == constant.KeywordMatchOr {
		for _, v := range keywordList {
			if utils.KMP(main, v) {
				return true
			}
		}
		return false
	} else {
		for _, v := range keywordList {
			if !utils.KMP(main, v) {
				return false
			}
		}
	}
	return true
}

// searchLocalMessages searches for local messages based on the given search parameters.
func (c *Conversation) searchLocalMessages(ctx context.Context, searchParam *sdk.SearchLocalMessagesParams) (*sdk.SearchLocalMessagesCallback, error) {
	var r sdk.SearchLocalMessagesCallback                                   // Initialize the result structure
	var startTime, endTime int64                                            // Variables to hold start and end times for the search
	var list []*model_struct.LocalChatLog                                   // Slice to store the search results
	conversationMap := make(map[string]*sdk.SearchByConversationResult, 10) // Map to store results grouped by conversation, with initial capacity of 10
	var err error                                                           // Variable to store any errors encountered
	var conversationID string                                               // Variable to store the current conversation ID

	// Clear the sequence cache for pull-up and pull-down operations in the search view,
	// to prevent the completion operations from the previous round from affecting the next round
	c.messagePullForwardEndSeqMap.DeleteByViewType(cache.ViewSearch)
	c.messagePullReverseEndSeqMap.DeleteByViewType(cache.ViewSearch)

	// Set the end time for the search; if SearchTimePosition is 0, use the current timestamp
	if searchParam.SearchTimePosition == 0 {
		endTime = time.Now().Unix()
	} else {
		endTime = searchParam.SearchTimePosition
	}

	// Set the start time based on the specified time period
	if searchParam.SearchTimePeriod != 0 {
		startTime = endTime - searchParam.SearchTimePeriod
	}
	// Convert start and end times to milliseconds
	startTime = utils.UnixSecondToTime(startTime).UnixNano() / 1e6
	endTime = utils.UnixSecondToTime(endTime).UnixNano() / 1e6

	// Validate that either keyword list or message type list is provided
	if len(searchParam.KeywordList) == 0 && len(searchParam.MessageTypeList) == 0 {
		return nil, errors.New("keywordlist and messageTypelist all null")
	}

	// Search in a specific conversation if ConversationID is provided
	if searchParam.ConversationID != "" {
		// Validate pagination parameters
		if searchParam.PageIndex < 1 || searchParam.Count < 1 {
			return nil, errors.New("page or count is null")
		}
		offset := (searchParam.PageIndex - 1) * searchParam.Count
		_, err := c.db.GetConversation(ctx, searchParam.ConversationID)
		if err != nil {
			return nil, err
		}
		// Search by content type or keyword based on provided parameters
		if len(searchParam.MessageTypeList) != 0 && len(searchParam.KeywordList) == 0 {
			list, err = c.db.SearchMessageByContentType(ctx, searchParam.MessageTypeList, searchParam.ConversationID, startTime, endTime, offset, searchParam.Count)
		} else {
			newContentTypeList := func(list []int) (result []int) {
				for _, v := range list {
					if utils.IsContainInt(v, SearchContentType) {
						result = append(result, v)
					}
				}
				return result
			}(searchParam.MessageTypeList)
			if len(newContentTypeList) == 0 {
				newContentTypeList = SearchContentType
			}
			list, err = c.db.SearchMessageByKeyword(ctx, newContentTypeList, searchParam.KeywordList, searchParam.KeywordListMatchType,
				searchParam.ConversationID, startTime, endTime, offset, searchParam.Count)
		}
	} else {
		// Comprehensive search across all conversations
		if len(searchParam.MessageTypeList) == 0 {
			searchParam.MessageTypeList = SearchContentType
		}
		list, err = c.searchMessageByContentTypeAndKeyword(ctx, searchParam.MessageTypeList, searchParam.KeywordList, searchParam.KeywordListMatchType, startTime, endTime)
	}

	// Handle any errors encountered during the search
	if err != nil {
		return nil, err
	}

	// Logging and processing each message in the search results
	log.ZDebug(ctx, "get raw data length is", "len", len(list))

	for _, v := range list {
		temp := LocalChatLogToMsgStruct(v)
		if c.filterMsg(temp, searchParam) {
			continue
		}
		// Determine the conversation ID based on the session type
		switch temp.SessionType {
		case constant.SingleChatType:
			if temp.SendID == c.loginUserID {
				conversationID = c.getConversationIDBySessionType(temp.RecvID, constant.SingleChatType)
			} else {
				conversationID = c.getConversationIDBySessionType(temp.SendID, constant.SingleChatType)
			}
		case constant.WriteGroupChatType:
			conversationID = c.getConversationIDBySessionType(temp.GroupID, constant.WriteGroupChatType)
		case constant.ReadGroupChatType:
			conversationID = c.getConversationIDBySessionType(temp.GroupID, constant.ReadGroupChatType)
		}
		// Populate the conversationMap with search results
		if oldItem, ok := conversationMap[conversationID]; !ok {
			searchResultItem := sdk.SearchByConversationResult{}
			localConversation, err := c.db.GetConversation(ctx, conversationID)
			if err != nil {
				// log.Error("", "get conversation err ", err.Error(), conversationID)
				continue
			}
			searchResultItem.ConversationID = conversationID
			searchResultItem.FaceURL = localConversation.FaceURL
			searchResultItem.ShowName = localConversation.ShowName
			searchResultItem.LatestMsgSendTime = localConversation.LatestMsgSendTime
			searchResultItem.ConversationType = localConversation.ConversationType
			searchResultItem.MessageList = append(searchResultItem.MessageList, temp)
			searchResultItem.MessageCount++
			conversationMap[conversationID] = &searchResultItem
		} else {
			oldItem.MessageCount++
			oldItem.MessageList = append(oldItem.MessageList, temp)
			conversationMap[conversationID] = oldItem
		}
	}

	// Compile the results from the conversationMap into the response structure
	for _, v := range conversationMap {
		r.SearchResultItems = append(r.SearchResultItems, v)
		r.TotalCount += v.MessageCount
	}

	// Sort the search results based on the latest message send time
	sort.Slice(r.SearchResultItems, func(i, j int) bool {
		return r.SearchResultItems[i].LatestMsgSendTime > r.SearchResultItems[j].LatestMsgSendTime
	})

	return &r, nil // Return the final search results
}

func (c *Conversation) searchMessageByContentTypeAndKeyword(ctx context.Context, contentType []int, keywordList []string,
	keywordListMatchType int, startTime, endTime int64) (result []*model_struct.LocalChatLog, err error) {
	var list []*model_struct.LocalChatLog
	conversationIDList, err := c.db.GetAllConversationIDList(ctx)
	if err != nil {
		return nil, err
	}

	var mu sync.Mutex
	g, _ := errgroup.WithContext(ctx)
	g.SetLimit(searchMessageGoroutineLimit)
	for _, v := range conversationIDList {
		conversationID := v
		g.Go(func() error {
			sList, err := c.db.SearchMessageByContentTypeAndKeyword(ctx, contentType, conversationID, keywordList, keywordListMatchType, startTime, endTime)
			if err != nil {
				log.ZWarn(ctx, "search conversation message", err, "conversationID", conversationID)
				return nil
			}

			mu.Lock()
			list = append(list, sList...)
			mu.Unlock()
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return list, nil
}

// true is filter, false is not filter
func (c *Conversation) filterMsg(temp *sdk_struct.MsgStruct, searchParam *sdk.SearchLocalMessagesParams) bool {
	switch temp.ContentType {
	case constant.Text:
		return !c.judgeMultipleSubString(searchParam.KeywordList, temp.TextElem.Content,
			searchParam.KeywordListMatchType)
	case constant.AtText:
		return !c.judgeMultipleSubString(searchParam.KeywordList, temp.AtTextElem.Text,
			searchParam.KeywordListMatchType)
	case constant.File:
		return !c.judgeMultipleSubString(searchParam.KeywordList, temp.FileElem.FileName,
			searchParam.KeywordListMatchType)
	case constant.Merger:
		if !c.judgeMultipleSubString(searchParam.KeywordList, temp.MergeElem.Title, searchParam.KeywordListMatchType) {
			for _, msgStruct := range temp.MergeElem.MultiMessage {
				if c.filterMsg(msgStruct, searchParam) {
					continue
				} else {
					break
				}
			}
		}
	case constant.Card:
		return !c.judgeMultipleSubString(searchParam.KeywordList, temp.CardElem.Nickname,
			searchParam.KeywordListMatchType)
	case constant.Location:
		return !c.judgeMultipleSubString(searchParam.KeywordList, temp.LocationElem.Description,
			searchParam.KeywordListMatchType)
	case constant.Custom:
		return !c.judgeMultipleSubString(searchParam.KeywordList, temp.CustomElem.Description,
			searchParam.KeywordListMatchType)
	case constant.Quote:
		if !c.judgeMultipleSubString(searchParam.KeywordList, temp.QuoteElem.Text, searchParam.KeywordListMatchType) {
			return c.filterMsg(temp.QuoteElem.QuoteMessage, searchParam)
		}
	case constant.Picture:
		fallthrough
	case constant.Sound:
		fallthrough
	case constant.Video:
		if len(searchParam.KeywordList) == 0 {
			return false
		} else {
			return true
		}
	default:
		return true
	}
	return false
}
