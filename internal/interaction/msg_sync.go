package interaction

import (
	"context"
	"errors"
	"fmt"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/common"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/db_interface"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	"github.com/openimsdk/openim-sdk-core/v3/sdk_struct"
	"github.com/openimsdk/protocol/msg"
	"github.com/openimsdk/tools/errs"

	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/log"
)

const (
	connectPullNums = 1
	defaultPullNums = 10
	SplitPullMsgNum = 100

	pullMsgGoroutineLimit = 10
)

// MsgSyncer is a central hub for message relay, responsible for sequential message gap pulling,
// handling network events, and managing app foreground and background events.
type MsgSyncer struct {
	loginUserID        string                // login user ID
	longConnMgr        *LongConnMgr          // long connection manager
	PushMsgAndMaxSeqCh chan common.Cmd2Value // channel for receiving push messages and the maximum SEQ number
	conversationCh     chan common.Cmd2Value // storage and session triggering
	syncedMaxSeqs      map[string]int64      // map of the maximum synced SEQ numbers for all group IDs
	syncedMaxSeqsLock  sync.RWMutex          // syncedMaxSeqs map lock
	db                 db_interface.DataBase // data store
	syncTimes          int                   // times of sync
	ctx                context.Context       // context
	reinstalled        bool                  //true if the app was uninstalled and reinstalled
	isSyncing          bool                  // indicates whether data is being synced
	isSyncingLock      sync.Mutex            // lock for syncing state

}

// NewMsgSyncer creates a new instance of the message synchronizer.
func NewMsgSyncer(ctx context.Context, conversationCh, PushMsgAndMaxSeqCh chan common.Cmd2Value,
	loginUserID string, longConnMgr *LongConnMgr, db db_interface.DataBase, syncTimes int) (*MsgSyncer, error) {
	m := &MsgSyncer{
		loginUserID:        loginUserID,
		longConnMgr:        longConnMgr,
		PushMsgAndMaxSeqCh: PushMsgAndMaxSeqCh,
		conversationCh:     conversationCh,
		ctx:                ctx,
		syncedMaxSeqs:      make(map[string]int64),
		db:                 db,
		syncTimes:          syncTimes,
	}
	if err := m.loadSeq(ctx); err != nil {
		log.ZError(ctx, "loadSeq err", err)
		return nil, err
	}
	return m, nil
}

// seq The db reads the data to the memory,set syncedMaxSeqs
func (m *MsgSyncer) loadSeq(ctx context.Context) error {
	conversationIDList, err := m.db.GetAllConversationIDList(ctx)
	if err != nil {
		log.ZError(ctx, "get conversation id list failed", err)
		return err
	}

	if len(conversationIDList) == 0 {
		version, err := m.db.GetAppSDKVersion(ctx)
		if err != nil && !errors.Is(err, errs.ErrRecordNotFound) {
			return err
		}
		if version == nil || !version.Installed {
			m.reinstalled = true
		}
	}

	// TODO With a large number of sessions(10w), this could potentially cause blocking and needs optimization.

	type SyncedSeq struct {
		ConversationID string
		MaxSyncedSeq   int64
		Err            error
	}

	partSize := 20
	currency := (len(conversationIDList)-1)/partSize + 1
	if len(conversationIDList) == 0 {
		currency = 0
	}
	var wg sync.WaitGroup
	resultMaps := make([]map[string]SyncedSeq, currency)

	for i := 0; i < currency; i++ {
		wg.Add(1)
		start := i * partSize
		end := start + partSize
		if i == currency-1 {
			end = len(conversationIDList)
		}

		resultMaps[i] = make(map[string]SyncedSeq)

		go func(i, start, end int) {
			defer wg.Done()
			for _, v := range conversationIDList[start:end] {
				maxSyncedSeq, err := m.db.CheckConversationNormalMsgSeq(ctx, v)
				resultMaps[i][v] = SyncedSeq{
					ConversationID: v,
					MaxSyncedSeq:   maxSyncedSeq,
					Err:            err,
				}
			}
		}(i, start, end)
	}

	wg.Wait()

	// merge map
	for _, resultMap := range resultMaps {
		for k, v := range resultMap {
			if v.Err != nil {
				log.ZError(ctx, "get group normal seq failed", errs.Wrap(v.Err), "conversationID", k)
				continue
			}
			m.syncedMaxSeqs[k] = v.MaxSyncedSeq
		}
	}
	notificationSeqs, err := m.db.GetNotificationAllSeqs(ctx)
	if err != nil {
		log.ZError(ctx, "get notification seq failed", err)
		return err
	}
	for _, notificationSeq := range notificationSeqs {
		m.syncedMaxSeqs[notificationSeq.ConversationID] = notificationSeq.Seq
	}
	log.ZDebug(ctx, "loadSeq", "syncedMaxSeqs", m.syncedMaxSeqs)
	return nil
}

// DoListener Listen to the message pipe of the message synchronizer
// and process received and pushed messages
func (m *MsgSyncer) DoListener(ctx context.Context) {
	defer func() {
		if r := recover(); r != nil {
			err := fmt.Sprintf("panic: %+v\n%s", r, debug.Stack())

			log.ZWarn(ctx, "DoListener panic", nil, "panic info", err)
		}
	}()
	for {
		select {
		case cmd := <-m.PushMsgAndMaxSeqCh:
			m.handlePushMsgAndEvent(cmd)
		case <-ctx.Done():
			log.ZInfo(m.ctx, "msg syncer done, sdk logout.....")
			return
		}
	}
}

// get seqs need sync interval
func (m *MsgSyncer) getSeqsNeedSync(syncedMaxSeq, maxSeq int64) []int64 {
	var seqs []int64
	for i := syncedMaxSeq + 1; i <= maxSeq; i++ {
		seqs = append(seqs, i)
	}
	return seqs
}

// recv msg from
func (m *MsgSyncer) handlePushMsgAndEvent(cmd common.Cmd2Value) {
	switch cmd.Cmd {
	case constant.CmdConnSuccesss:
		log.ZInfo(cmd.Ctx, "recv long conn mgr connected", "cmd", cmd.Cmd, "value", cmd.Value)
		if m.startSync() {
			m.doConnected(cmd.Ctx)
		} else {
			log.ZWarn(cmd.Ctx, "syncing, ignore connected event", nil, "cmd", cmd.Cmd, "value", cmd.Value)
		}
	case constant.CmdWakeUpDataSync:
		log.ZInfo(cmd.Ctx, "app wake up, start sync msgs", "cmd", cmd.Cmd, "value", cmd.Value)
		if m.startSync() {
			m.doWakeupDataSync(cmd.Ctx)
		} else {
			log.ZWarn(cmd.Ctx, "syncing, ignore wake up event", nil, "cmd", cmd.Cmd, "value", cmd.Value)

		}
	case constant.CmdPushMsg:
		m.doPushMsg(cmd.Ctx, cmd.Value.(*sdkws.PushMessages))
	}
}

func (m *MsgSyncer) compareSeqsAndBatchSync(ctx context.Context, maxSeqToSync map[string]int64, pullNums int64) {
	needSyncSeqMap := make(map[string][2]int64)
	//when app reinstalled do not pull notifications messages.
	if m.reinstalled {
		notificationsSeqMap := make(map[string]int64)
		messagesSeqMap := make(map[string]int64)
		for conversationID, seq := range maxSeqToSync {
			if IsNotification(conversationID) {
				if seq != 0 { // seq is 0, no need to sync
					notificationsSeqMap[conversationID] = seq
				}
			} else {
				messagesSeqMap[conversationID] = seq
			}
		}

		var notificationSeqs []*model_struct.NotificationSeqs

		for conversationID, seq := range notificationsSeqMap {
			notificationSeqs = append(notificationSeqs, &model_struct.NotificationSeqs{
				ConversationID: conversationID,
				Seq:            seq,
			})
			m.syncedMaxSeqs[conversationID] = seq
		}

		if len(notificationSeqs) > 0 {
			err := m.db.BatchInsertNotificationSeq(ctx, notificationSeqs)
			if err != nil {
				log.ZWarn(ctx, "BatchInsertNotificationSeq err", err)
			}
		}

		for conversationID, maxSeq := range messagesSeqMap {
			if syncedMaxSeq, ok := m.syncedMaxSeqs[conversationID]; ok {
				if maxSeq > syncedMaxSeq {
					needSyncSeqMap[conversationID] = [2]int64{syncedMaxSeq + 1, maxSeq}
				}
			} else {
				needSyncSeqMap[conversationID] = [2]int64{0, maxSeq}
			}
		}
		defer func() {
			if err := m.db.SetAppSDKVersion(ctx, &model_struct.LocalAppSDKVersion{
				Installed: true,
			}); err != nil {
				log.ZError(ctx, "SetAppSDKVersion err", err)
			}
			m.reinstalled = false
		}()
		_ = m.syncAndTriggerReinstallMsgs(m.ctx, needSyncSeqMap, pullNums)
	} else {
		for conversationID, maxSeq := range maxSeqToSync {
			if syncedMaxSeq, ok := m.syncedMaxSeqs[conversationID]; ok {
				if maxSeq > syncedMaxSeq {
					needSyncSeqMap[conversationID] = [2]int64{syncedMaxSeq + 1, maxSeq}
				}
			} else {
				if maxSeq != 0 { // seq is 0, no need to sync
					needSyncSeqMap[conversationID] = [2]int64{0, maxSeq}
				}
			}
		}
		_ = m.syncAndTriggerMsgs(m.ctx, needSyncSeqMap, pullNums)
	}
}

// startSync checks if the sync is already in progress.
// If syncing is in progress, it returns false. Otherwise, it starts syncing and returns true.
func (m *MsgSyncer) startSync() bool {
	m.isSyncingLock.Lock()
	defer m.isSyncingLock.Unlock()

	if m.isSyncing {
		// If already syncing, return false
		return false
	}

	// Set syncing to true and start the sync
	m.isSyncing = true

	// Create a goroutine that waits for 5 seconds and then sets isSyncing to false
	go func() {
		time.Sleep(5 * time.Second)
		m.isSyncingLock.Lock()
		m.isSyncing = false
		m.isSyncingLock.Unlock()
	}()

	return true
}

func (m *MsgSyncer) doPushMsg(ctx context.Context, push *sdkws.PushMessages) {
	log.ZDebug(ctx, "push msgs", "push", push, "syncedMaxSeqs", m.syncedMaxSeqs)
	m.pushTriggerAndSync(ctx, push.Msgs, m.triggerConversation)
	m.pushTriggerAndSync(ctx, push.NotificationMsgs, m.triggerNotification)
}

func (m *MsgSyncer) pushTriggerAndSync(ctx context.Context, pushMessages map[string]*sdkws.PullMsgs, triggerFunc func(ctx context.Context, msgs map[string]*sdkws.PullMsgs) error) {
	if len(pushMessages) == 0 {
		return
	}
	needSyncSeqMap := make(map[string][2]int64)
	var lastSeq int64
	var storageMsgs []*sdkws.MsgData
	for conversationID, msgs := range pushMessages {
		for _, msg := range msgs.Msgs {
			if msg.Seq == 0 {
				_ = triggerFunc(ctx, map[string]*sdkws.PullMsgs{conversationID: {Msgs: []*sdkws.MsgData{msg}}})
				continue
			}
			lastSeq = msg.Seq
			storageMsgs = append(storageMsgs, msg)
		}
		if lastSeq == m.syncedMaxSeqs[conversationID]+int64(len(storageMsgs)) && lastSeq != 0 {
			log.ZDebug(ctx, "trigger msgs", "msgs", storageMsgs)
			_ = triggerFunc(ctx, map[string]*sdkws.PullMsgs{conversationID: {Msgs: storageMsgs}})
			m.syncedMaxSeqs[conversationID] = lastSeq
		} else if lastSeq != 0 && lastSeq > m.syncedMaxSeqs[conversationID] {
			//must pull message when message type is notification
			needSyncSeqMap[conversationID] = [2]int64{m.syncedMaxSeqs[conversationID] + 1, lastSeq}
		}
	}
	m.syncAndTriggerMsgs(ctx, needSyncSeqMap, defaultPullNums)
}

// Called after successful reconnection to synchronize the latest message
func (m *MsgSyncer) doConnected(ctx context.Context) {
	reinstalled := m.reinstalled
	if reinstalled {
		common.TriggerCmdSyncFlag(m.ctx, constant.AppDataSyncStart, m.conversationCh)
	} else {
		common.TriggerCmdSyncFlag(m.ctx, constant.MsgSyncBegin, m.conversationCh)
	}
	var resp sdkws.GetMaxSeqResp
	if err := m.longConnMgr.SendReqWaitResp(m.ctx, &sdkws.GetMaxSeqReq{UserID: m.loginUserID}, constant.GetNewestSeq, &resp); err != nil {
		log.ZError(m.ctx, "get max seq error", err)
		common.TriggerCmdSyncFlag(m.ctx, constant.MsgSyncFailed, m.conversationCh)
		return
	} else {
		log.ZDebug(m.ctx, "get max seq success", "resp", resp.MaxSeqs)
	}
	m.compareSeqsAndBatchSync(ctx, resp.MaxSeqs, connectPullNums)
	if reinstalled {
		common.TriggerCmdSyncFlag(m.ctx, constant.AppDataSyncFinish, m.conversationCh)
	} else {
		common.TriggerCmdSyncFlag(m.ctx, constant.MsgSyncEnd, m.conversationCh)
	}
}

func (m *MsgSyncer) doWakeupDataSync(ctx context.Context) {
	common.TriggerCmdSyncData(ctx, m.conversationCh)
	var resp sdkws.GetMaxSeqResp
	if err := m.longConnMgr.SendReqWaitResp(m.ctx, &sdkws.GetMaxSeqReq{UserID: m.loginUserID}, constant.GetNewestSeq, &resp); err != nil {
		log.ZError(m.ctx, "get max seq error", err)
		return
	} else {
		log.ZDebug(m.ctx, "get max seq success", "resp", resp.MaxSeqs)
	}
	m.compareSeqsAndBatchSync(ctx, resp.MaxSeqs, defaultPullNums)
}

func IsNotification(conversationID string) bool {
	return strings.HasPrefix(conversationID, "n_")
}

func (m *MsgSyncer) syncAndTriggerMsgs(ctx context.Context, seqMap map[string][2]int64, syncMsgNum int64) error {
	if len(seqMap) == 0 {
		log.ZDebug(ctx, "nothing to sync", "syncMsgNum", syncMsgNum)
		return nil
	}

	log.ZDebug(ctx, "current sync seqMap", "seqMap", seqMap)
	var (
		tempSeqMap = make(map[string][2]int64, 50)
		msgNum     = 0
	)

	for k, v := range seqMap {
		oneConversationSyncNum := v[1] - v[0] + 1
		tempSeqMap[k] = v
		// For notification conversations, use oneConversationSyncNum directly
		if IsNotification(k) {
			msgNum += int(oneConversationSyncNum)
		} else {
			// For regular conversations, ensure msgNum is the minimum of oneConversationSyncNum and syncMsgNum
			currentSyncMsgNum := int64(0)
			if oneConversationSyncNum > syncMsgNum {
				currentSyncMsgNum = syncMsgNum
			} else {
				currentSyncMsgNum = oneConversationSyncNum
			}
			msgNum += int(currentSyncMsgNum)
		}

		// If accumulated msgNum reaches SplitPullMsgNum, trigger a batch pull
		if msgNum >= SplitPullMsgNum {
			resp, err := m.pullMsgBySeqRange(ctx, tempSeqMap, syncMsgNum)
			if err != nil {
				log.ZError(ctx, "syncMsgFromServer error", err, "tempSeqMap", tempSeqMap)
				return err
			}
			_ = m.triggerConversation(ctx, resp.Msgs)
			_ = m.triggerNotification(ctx, resp.NotificationMsgs)
			for conversationID, seqs := range tempSeqMap {
				m.syncedMaxSeqs[conversationID] = seqs[1]
			}
			// Reset tempSeqMap and msgNum to handle the next batch
			tempSeqMap = make(map[string][2]int64, 50)
			msgNum = 0
		}
	}

	// Handle remaining messages to ensure all are synced
	if len(tempSeqMap) > 0 {
		resp, err := m.pullMsgBySeqRange(ctx, tempSeqMap, syncMsgNum)
		if err != nil {
			log.ZError(ctx, "syncMsgFromServer error", err, "tempSeqMap", tempSeqMap)
			return err
		}
		_ = m.triggerConversation(ctx, resp.Msgs)
		_ = m.triggerNotification(ctx, resp.NotificationMsgs)
		for conversationID, seqs := range tempSeqMap {
			m.syncedMaxSeqs[conversationID] = seqs[1]
		}
	}

	return nil
}

// Fragment synchronization message, seq refresh after successful trigger
func (m *MsgSyncer) syncAndTriggerReinstallMsgs(ctx context.Context, seqMap map[string][2]int64, syncMsgNum int64) error {
	if len(seqMap) > 0 {
		log.ZDebug(ctx, "current sync seqMap", "seqMap", seqMap)
		var (
			tempSeqMap = make(map[string][2]int64, 50)
			msgNum     = 0
			total      = len(seqMap)
			gr         *errgroup.Group
		)
		gr, _ = errgroup.WithContext(ctx)
		gr.SetLimit(pullMsgGoroutineLimit)
		for k, v := range seqMap {
			oneConversationSyncNum := min(v[1]-v[0]+1, syncMsgNum)
			tempSeqMap[k] = v
			if oneConversationSyncNum > 0 {
				msgNum += int(oneConversationSyncNum)
			}
			if msgNum >= SplitPullMsgNum {
				tpSeqMap := make(map[string][2]int64, len(tempSeqMap))
				for k, v := range tempSeqMap {
					tpSeqMap[k] = v
				}

				gr.Go(func() error {
					resp, err := m.pullMsgBySeqRange(ctx, tpSeqMap, syncMsgNum)
					if err != nil {
						log.ZError(ctx, "syncMsgFromServer err", err, "tempSeqMap", tpSeqMap)
						return err
					}
					m.checkMessagesAndGetLastMessage(ctx, resp.Msgs)
					_ = m.triggerReinstallConversation(ctx, resp.Msgs, total)
					for conversationID, seqs := range tpSeqMap {
						m.syncedMaxSeqsLock.Lock()
						m.syncedMaxSeqs[conversationID] = seqs[1]
						m.syncedMaxSeqsLock.Unlock()
					}
					return nil
				})

				tempSeqMap = make(map[string][2]int64, 50)
				msgNum = 0
			}
		}
		gr.Go(func() error {
			resp, err := m.pullMsgBySeqRange(ctx, tempSeqMap, syncMsgNum)
			if err != nil {
				log.ZError(ctx, "syncMsgFromServer err", err, "seqMap", seqMap)
				return err
			}
			m.checkMessagesAndGetLastMessage(ctx, resp.Msgs)
			_ = m.triggerReinstallConversation(ctx, resp.Msgs, total)
			for conversationID, seqs := range seqMap {
				m.syncedMaxSeqsLock.Lock()
				m.syncedMaxSeqs[conversationID] = seqs[1]
				m.syncedMaxSeqsLock.Unlock()
			}
			return nil
		})
		if err := gr.Wait(); err != nil {
			return err
		}

	} else {
		log.ZDebug(ctx, "noting conversation to sync", "syncMsgNum", syncMsgNum)
	}
	return nil
}
func (m *MsgSyncer) checkMessagesAndGetLastMessage(ctx context.Context, messages map[string]*sdkws.PullMsgs) {
	var conversationIDs []string

	for conversationID, message := range messages {
		allInValid := true
		for _, data := range message.Msgs {
			if data.Status < constant.MsgStatusHasDeleted {
				allInValid = false
				break
			}
		}
		if allInValid {
			conversationIDs = append(conversationIDs, conversationID)
		}
	}
	if len(conversationIDs) > 0 {
		resp, err := m.fetchLatestValidMessages(ctx, conversationIDs)
		if err != nil {
			log.ZError(ctx, "fetchLatestValidMessages", err, "conversationIDs", conversationIDs)
			return
		}
		for conversationID, message := range resp.Msgs {
			messages[conversationID] = &sdkws.PullMsgs{Msgs: []*sdkws.MsgData{message}}
		}
	}

}

func (m *MsgSyncer) splitSeqs(split int, seqsNeedSync []int64) (splitSeqs [][]int64) {
	if len(seqsNeedSync) <= split {
		splitSeqs = append(splitSeqs, seqsNeedSync)
		return
	}
	for i := 0; i < len(seqsNeedSync); i += split {
		end := i + split
		if end > len(seqsNeedSync) {
			end = len(seqsNeedSync)
		}
		splitSeqs = append(splitSeqs, seqsNeedSync[i:end])
	}
	return
}

func (m *MsgSyncer) pullMsgBySeqRange(ctx context.Context, seqMap map[string][2]int64, syncMsgNum int64) (resp *sdkws.PullMessageBySeqsResp, err error) {
	log.ZDebug(ctx, "pullMsgBySeqRange", "seqMap", seqMap, "syncMsgNum", syncMsgNum)

	req := sdkws.PullMessageBySeqsReq{UserID: m.loginUserID}
	for conversationID, seqs := range seqMap {
		req.SeqRanges = append(req.SeqRanges, &sdkws.SeqRange{
			ConversationID: conversationID,
			Begin:          seqs[0],
			End:            seqs[1],
			Num:            syncMsgNum,
		})
	}
	resp = &sdkws.PullMessageBySeqsResp{}
	if err := m.longConnMgr.SendReqWaitResp(ctx, &req, constant.PullMsgByRange, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *MsgSyncer) fetchLatestValidMessages(ctx context.Context, conversationID []string) (resp *msg.GetLastMessageResp, err error) {
	log.ZDebug(ctx, "fetchLatestValidMessages", "conversationID", conversationID)

	req := msg.GetLastMessageReq{
		UserID:          m.loginUserID,
		ConversationIDs: conversationID,
	}
	resp = &msg.GetLastMessageResp{}
	if err := m.longConnMgr.SendReqWaitResp(ctx, &req, constant.PullConvLastMessage, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// synchronizes messages by SEQs.
func (m *MsgSyncer) syncMsgBySeqs(ctx context.Context, conversationID string, seqsNeedSync []int64) (allMsgs []*sdkws.MsgData, err error) {
	pullMsgReq := sdkws.PullMessageBySeqsReq{}
	pullMsgReq.UserID = m.loginUserID
	split := constant.SplitPullMsgNum
	seqsList := m.splitSeqs(split, seqsNeedSync)
	for i := 0; i < len(seqsList); {
		var pullMsgResp sdkws.PullMessageBySeqsResp
		err := m.longConnMgr.SendReqWaitResp(ctx, &pullMsgReq, constant.PullMsgByRange, &pullMsgResp)
		if err != nil {
			log.ZError(ctx, "syncMsgFromServerSplit err", err, "pullMsgReq", pullMsgReq)
			continue
		}
		i++
		allMsgs = append(allMsgs, pullMsgResp.Msgs[conversationID].Msgs...)
	}
	return allMsgs, nil
}

// triggers a conversation with a new message.
func (m *MsgSyncer) triggerConversation(ctx context.Context, msgs map[string]*sdkws.PullMsgs) error {
	if len(msgs) > 0 {
		err := common.TriggerCmdNewMsgCome(ctx, sdk_struct.CmdNewMsgComeToConversation{Msgs: msgs}, m.conversationCh)
		if err != nil {
			log.ZError(ctx, "triggerCmdNewMsgCome err", err, "msgs", msgs)
		}
		log.ZDebug(ctx, "triggerConversation", "msgs", msgs)
		return err
	} else {
		log.ZDebug(ctx, "triggerConversation is nil", "msgs", msgs)
	}
	return nil
}

// triggers a conversation with a new message.
func (m *MsgSyncer) triggerReinstallConversation(ctx context.Context, msgs map[string]*sdkws.PullMsgs, total int) (err error) {
	if len(msgs) > 0 {
		err = common.TriggerCmdMsgSyncInReinstall(ctx, sdk_struct.CmdMsgSyncInReinstall{
			Msgs:  msgs,
			Total: total,
		}, m.conversationCh)
		if err != nil {
			log.ZError(ctx, "triggerCmdNewMsgCome err", err, "msgs", msgs)
		}
		log.ZDebug(ctx, "triggerConversation", "msgs", msgs)
		return err
	} else {
		log.ZDebug(ctx, "triggerConversation is nil", "msgs", msgs)
	}
	return nil
}

func (m *MsgSyncer) triggerNotification(ctx context.Context, msgs map[string]*sdkws.PullMsgs) error {
	if len(msgs) > 0 {
		common.TriggerCmdNotification(ctx, sdk_struct.CmdNewMsgComeToConversation{Msgs: msgs}, m.conversationCh)
	} else {
		log.ZDebug(ctx, "triggerNotification is nil", "notifications", msgs)
	}
	return nil

}
