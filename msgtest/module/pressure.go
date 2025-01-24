package module

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/openimsdk/tools/utils/datautil"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/utils"
	"github.com/openimsdk/openim-sdk-core/v3/sdk_struct"
	"github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/tools/log"
)

var (
	TESTIP        = "127.0.0.1"
	APIADDR       = fmt.Sprintf("http://%v:10002", TESTIP)
	WSADDR        = fmt.Sprintf("ws://%v:10001", TESTIP)
	SECRET        = "openIM123"
	MANAGERUSERID = "imAdmin"

	PLATFORMID = constant.AndroidPlatformID
	LogLevel   = uint32(5)
)

var (
	totalOnlineUserNum    = 200000 // total online users num
	friendMsgSenderNum    = 200    // friend msg sender num
	NotFriendMsgSenderNum = 200    // not friend msg sender num
	groupMsgSenderNum     = 200    // group msg sender num
	msgSenderNumEvreyUser = 100    // the number of messages per user.
	fastenedUserNum       = 600    // fixed number of users

	recvMsgUserNum       = 20 // the number of message recipients, sampled accounts
	singleSampleUserList []string
)

const (
	HundredThousandGroupUserNum = 100000
	TenThousandGroupUserNum     = 10000
	ThousandGroupUserNum        = 1000
	HundredGroupUserNum         = 100
	FiftyGroupUserNum           = 50
	TenGroupUserNum             = 10

	HundredThousandGroupNum = 1
	TenThousandGroupNum     = 2
	ThousandGroupNum        = 5
	HundredGroupNum         = 50
	FiftyGroupNum           = 100
	TenGroupNum             = 1000

	FastenedUserPrefix  = "test_v3_u"
	OfflineUserPrefix   = "o"
	RecvMsgPrefix       = "recv_msg_prefix"
	singleMsgRecvPrefix = "single_msg_recv_prefix"
)

type PressureTester struct {
	friendManager          *TestFriendManager
	userManager            *TestUserManager
	groupManager           *TestGroupManager
	msgSender              map[string]*SendMsgUser
	rw                     sync.RWMutex
	groupRandomSender      map[string][]string
	groupRandomOnlineUsers map[string][]string
	groupOwnerUserID       map[string]string
	groupMemberNum         map[string]int
	timeOffset             int64
	singleSendNum          atomic.Int64

	groupSenderUserIDs, friendSenderUserIDs, notfriendSenderUserIDs []string
	recvMsgUserIDs                                                  []string
	offlineUserIDs                                                  []string

	tenThousandGroupIDs, thousandGroupIDs, hundredGroupUserIDs, fiftyGroupUserIDs []string
}

func (p *PressureTester) SetOfflineUserIDs(offlineUserIDs []string) {
	p.offlineUserIDs = offlineUserIDs
}

func (p *PressureTester) FormatGroupInfo(ctx context.Context) {

	groupsByMemberNum := make(map[int][]string)

	for groupID, memberNum := range p.groupMemberNum {
		groupsByMemberNum[memberNum] = append(groupsByMemberNum[memberNum], groupID)
	}

	if len(p.groupMemberNum) == 0 {
		log.ZWarn(ctx, "no group created", nil)
		return
	}
	fmt.Println("---------------------------")
	for memberNum, groupIDs := range groupsByMemberNum {
		log.ZWarn(ctx, "member count", nil, "memberNum", memberNum)
		log.ZWarn(ctx, "group num", nil, "groupNum", len(groupIDs))
		fmt.Println("---------------------------")
	}
}

func (p *PressureTester) GetSingleSendNum() int64 {
	return p.singleSendNum.Load()
}

func NewPressureTester() (*PressureTester, error) {
	metaManager := NewMetaManager(APIADDR, SECRET, MANAGERUSERID)
	if err := metaManager.initToken(); err != nil {
		return nil, err
	}
	serverTime, err := metaManager.GetServerTime()
	if err != nil {
		return nil, err
	}
	log.ZWarn(context.Background(), "server time is", nil, "serverTime", serverTime, "current time",
		utils.GetCurrentTimestampByMill(), "time offset", serverTime-utils.GetCurrentTimestampByMill())

	return &PressureTester{friendManager: metaManager.NewFriendManager(), userManager: metaManager.NewUserManager(),
		groupManager:      metaManager.NewGroupMananger(),
		msgSender:         make(map[string]*SendMsgUser),
		groupRandomSender: make(map[string][]string), groupOwnerUserID: make(map[string]string),
		groupMemberNum: make(map[string]int),
		timeOffset:     serverTime - utils.GetCurrentTimestampByMill()}, nil
}

func (p *PressureTester) genUserIDs() (userIDs, fastenedUserIDs, recvMsgUserIDs []string) {
	userIDs = p.userManager.GenUserIDs(totalOnlineUserNum - fastenedUserNum)                  // online users
	fastenedUserIDs = p.userManager.GenUserIDsWithPrefix(fastenedUserNum, FastenedUserPrefix) // specify the fixed users for sending messages
	recvMsgUserIDs = p.userManager.GenUserIDsWithPrefix(recvMsgUserNum, RecvMsgPrefix)        // complete SDK for sampling users
	return
}

// selectSample
func (p *PressureTester) SelectSample(total int, percentage float64) (fastenedUserIDs []string,
	sampleReceiver, offlineUserIDs []string, err error) {
	if percentage < 0 || percentage > 1 {
		return nil, nil, nil, fmt.Errorf("percentage must be between 0 and 1")
	}

	fastenedUserIDs = p.userManager.GenUserIDsWithPrefix(total, FastenedUserPrefix)
	offlineUserIDs = p.userManager.GenSEUserIDsWithPrefix(total, 2*total, OfflineUserPrefix)
	step := int(1.0 / percentage)
	for i := 0; i < total; i += step {
		sampleReceiver = append(sampleReceiver, fmt.Sprintf("%s_testv3_%d", FastenedUserPrefix, i))
	}
	singleSampleUserList = sampleReceiver
	return fastenedUserIDs, sampleReceiver, offlineUserIDs, nil

}
func (p *PressureTester) SelectSampleFromStarEnd(start, end int, percentage float64) (fastenedUserIDs []string,
	sampleReceiver, offlineUserIDs []string, err error) {
	if percentage < 0 || percentage > 1 {
		return nil, nil, nil, fmt.Errorf("percentage must be between 0 and 1")
	}
	fastenedUserIDs = p.userManager.GenSEUserIDsWithPrefix(start, end, FastenedUserPrefix)
	offlineUserIDs = p.userManager.GenSEUserIDsWithPrefix(end, 2*end-start, OfflineUserPrefix)
	step := int(1.0 / percentage)
	for i := start; i < end; i += step {
		sampleReceiver = append(sampleReceiver, fmt.Sprintf("%s%d", FastenedUserPrefix, i))
	}
	singleSampleUserList = sampleReceiver
	return fastenedUserIDs, sampleReceiver, offlineUserIDs, nil
}

func (p *PressureTester) SelectStartAndEnd(start, end int) (fastenedUserIDs []string) {
	return p.userManager.GenSEUserIDsWithPrefix(start, end, FastenedUserPrefix)
}

func (p *PressureTester) RegisterUsers(userIDs []string, fastenedUserIDs []string, recvMsgUserIDs []string) error {
	for i := 0; i < len(userIDs); i += 1000 {
		end := i + 1000
		if end > len(userIDs) {
			end = len(userIDs)
		}
		userIDsSlice := userIDs[i:end]
		if err := p.userManager.RegisterUsers(userIDsSlice...); err != nil {
			return err
		}
		if len(userIDsSlice) < 1000 {
			break
		}
	}
	if len(fastenedUserIDs) != 0 {
		if err := p.userManager.RegisterUsers(fastenedUserIDs...); err != nil {
			return err
		}
	}
	if len(recvMsgUserIDs) != 0 {
		if err := p.userManager.RegisterUsers(recvMsgUserIDs...); err != nil {
			return err
		}
	}
	return nil
}

func (p *PressureTester) InitUserConns(userIDs []string) {
	for _, userID := range userIDs {
		token, err := p.userManager.GetToken(userID, int32(PLATFORMID))
		if err != nil {
			log.ZError(context.Background(), "get token failed", err, "userID", userID, "platformID", PLATFORMID)
			continue
		}
		user := NewUser(userID, token, p.timeOffset, p, sdk_struct.IMConfig{WsAddr: WSADDR, ApiAddr: APIADDR, PlatformID: int32(PLATFORMID)})
		p.msgSender[userID] = user

	}

}

func (p *PressureTester) getGroup(fastenedUserIDs []string, groupMemberNum int, groupSenderRate, groupOnlineRate float64) (ownerUserID string,
	userIDs []string, randomSender []string) {
	//get the group online users
	olineUserIDNum := int(float64(groupMemberNum) * groupOnlineRate)

	userIDs = p.Shuffle(fastenedUserIDs, olineUserIDNum)
	//get the group offline users
	offlineUserID := p.Shuffle(p.offlineUserIDs, groupMemberNum-olineUserIDNum)
	ownerUserID = p.Shuffle(userIDs, 1)[0]
	randomSender = p.Shuffle(userIDs, int(float64(groupMemberNum)*groupSenderRate))
	return ownerUserID, append(datautil.DeleteElems(userIDs, ownerUserID), offlineUserID...), randomSender
}

func (p *PressureTester) CreateTestGroups(fastenedUserIDs []string, total int, groupSenderRate, groupOnlineRate float64, hundredThousandGroupNum, tenThousandGroupNum, thousandGroupNum,
	hundredGroupNum, fiftyGroupNum, tenGroupNum int) (err error) {
	// create ten thousand group
	if hundredThousandGroupNum != 0 {
		if total < HundredThousandGroupUserNum {
			return fmt.Errorf("total user num must be greater than 100000")
		}
	}
	if tenThousandGroupNum != 0 {
		if total < TenThousandGroupUserNum {
			return fmt.Errorf("total user num must be greater than 10000")
		}

	}
	if thousandGroupNum != 0 {
		if total < ThousandGroupUserNum {
			return fmt.Errorf("total user num must be greater than 1000")
		}
	}

	if hundredGroupNum != 0 {
		if total < HundredGroupUserNum {
			return fmt.Errorf("total user num must be greater than 100")
		}

	}
	if fiftyGroupNum != 0 {
		if total < FiftyGroupUserNum {
			return fmt.Errorf("total user num must be greater than 50")
		}
	}

	if tenGroupNum != 0 {
		if total < TenGroupUserNum {
			return fmt.Errorf("total user num must be greater than 10")
		}

	}

	f := func(GroupNum int, GroupUserNum int, groupSenderRate, groupOnlineRate float64, groupIDAndNameString string) (err error) {
		for i := 1; i <= GroupNum; i++ {
			if groupOnlineRate < groupSenderRate {
				return fmt.Errorf("group online rate must > group sender rate")
			}
			ownerUserID, memberUserIDs, randomSenderUserIDs := p.getGroup(fastenedUserIDs, GroupUserNum, groupSenderRate, groupOnlineRate)
			groupID := p.groupManager.GenGroupID(fmt.Sprintf(groupIDAndNameString+"_%d", i))
			err = p.groupManager.CreateGroup(groupID, fmt.Sprintf(groupIDAndNameString+"_%d", i), ownerUserID,
				memberUserIDs)
			if err != nil {
				return
			}
			p.groupRandomSender[groupID] = randomSenderUserIDs
			p.groupOwnerUserID[groupID] = ownerUserID
			p.groupMemberNum[groupID] = GroupUserNum
		}
		return nil
	}
	err = f(hundredThousandGroupNum, HundredThousandGroupUserNum, groupSenderRate, groupOnlineRate, "hundredThousandGroupUserNum")
	if err != nil {
		return err
	}
	err = f(tenThousandGroupNum, TenThousandGroupUserNum, groupSenderRate, groupOnlineRate, "tenThousandGroupUserNum")
	if err != nil {
		return err
	}
	err = f(thousandGroupNum, ThousandGroupUserNum, groupSenderRate, groupOnlineRate, "thousandGroupUserNum")
	if err != nil {
		return err
	}
	err = f(hundredGroupNum, HundredGroupUserNum, groupSenderRate, groupOnlineRate, "hundredGroupUserNum")
	if err != nil {
		return err
	}
	err = f(fiftyGroupNum, FiftyGroupUserNum, groupSenderRate, groupOnlineRate, "fiftyGroupUserNum")
	if err != nil {
		return err
	}
	err = f(tenGroupNum, TenGroupUserNum, groupSenderRate, groupOnlineRate, "tenGroupUserNum")
	if err != nil {
		return err
	}

	return nil
}

func (p *PressureTester) SendGroupMessage(ctx context.Context, num int, duration time.Duration) {
	var wg sync.WaitGroup
	log.ZWarn(ctx, "send group message start", nil, "groupNum", len(p.groupOwnerUserID))
	if len(p.groupOwnerUserID) == 0 || num == 0 {
		log.ZWarn(ctx, "send group message over,do not need to send group message", nil)
	}
	defer log.ZWarn(ctx, "send group message over", nil, "groupNum", len(p.groupOwnerUserID))
	for groupID, _ := range p.groupOwnerUserID {
		wg.Add(1)
		go func(groupID string) {
			p.rw.RLock()
			if senderUserIDs, ok := p.groupRandomSender[groupID]; ok {
				p.rw.RUnlock()
				p.sendMessage2Groups(senderUserIDs, groupID, num, duration)
			}
			wg.Done()
		}(groupID)
	}
	wg.Wait()
}

func (p *PressureTester) sendMessage2Groups(senderIDs []string, groupID string, num int, duration time.Duration) {
	var wg sync.WaitGroup
	for _, senderID := range senderIDs {
		wg.Add(1)
		go func(senderID, groupID string) {
			defer wg.Done()
			for i := 0; i < num; i++ {
				if user, ok := p.msgSender[senderID]; ok {
					user.SendGroupMsgWithContext(groupID, i)
				}
				time.Sleep(duration)
			}
		}(senderID, groupID)

	}
	wg.Wait()
}

func (p *PressureTester) SendSingleMessages(ctx context.Context, fastenedUserIDs []string, randomSender []string, randomReceiver, num int, duration time.Duration) {
	log.ZWarn(ctx, "send single message start", nil, "randomSender", len(randomSender), "randomReceiver", randomReceiver)
	if len(randomSender) == 0 || randomReceiver == 0 || num == 0 {
		log.ZWarn(ctx, "send single message over,do not need to send single message", nil)
		return
	}
	defer log.ZWarn(ctx, "send single message over", nil, "randomSender", len(randomSender), "randomReceiver", randomReceiver)
	var wg sync.WaitGroup
	length := len(fastenedUserIDs)
	rand.Seed(time.Now().UnixNano())
	for _, userID := range randomSender {
		counter := 0
		var receiverUserIDs []string
		for counter < randomReceiver {
			index := rand.Intn(length)
			if fastenedUserIDs[index] != userID {
				counter++
				receiverUserIDs = append(receiverUserIDs, fastenedUserIDs[index])
			}
		}
		wg.Add(1)
		go func(receiverUserIDs []string, u string) {
			//log.ZError(context.Background(), "SendSingleMessages", nil, "length", len(receiverUserIDs))
			defer wg.Done()
			user, _ := p.msgSender[u]
			for _, rv := range receiverUserIDs {
				for x := 0; x < num; x++ {
					user.SendMsgWithContext(rv, x)
					p.singleSendNum.Add(1)
					time.Sleep(duration)
				}
			}
		}(receiverUserIDs, userID)
	}
	wg.Wait()
}

// Shuffle gets random userID from fastenedUserIDs and returns a slice of userID with length of needNum.
func (p *PressureTester) Shuffle(fastenedUserIDs []string, needNum int) []string {

	rand.Shuffle(len(fastenedUserIDs), func(i, j int) {
		fastenedUserIDs[i], fastenedUserIDs[j] = fastenedUserIDs[j], fastenedUserIDs[i]
	})

	selectedUserIDs := make([]string, 0, needNum)
	seen := make(map[string]bool)

	for _, userID := range fastenedUserIDs {
		if len(selectedUserIDs) == needNum {
			break
		}

		if !seen[userID] {
			selectedUserIDs = append(selectedUserIDs, userID)
			seen[userID] = true
		}
	}
	return selectedUserIDs
}

func (p *PressureTester) SendSingleMessagesTo(fastenedUserIDs []string, num int, duration time.Duration) {
	var wg sync.WaitGroup
	//length := len(fastenedUserIDs)
	rand.Seed(time.Now().UnixNano())
	for i, userID := range fastenedUserIDs {
		var receiverUserIDs []string
		for i < num {
			receiverUserIDs = append(receiverUserIDs, utils.IntToString(i))
			i++
		}
		wg.Add(1)
		go func(receiverUserIDs []string, u string) {
			defer wg.Done()
			user, _ := p.msgSender[u]
			for j, rv := range receiverUserIDs {
				user.SendMsgWithContext(rv, j)
				time.Sleep(duration)

			}
		}(receiverUserIDs, userID)
	}
	wg.Wait()

}

func (p *PressureTester) SendMessages(sendID, recvID string, msgNum int) {
	var i = 0
	var ws sync.WaitGroup
	user, _ := p.msgSender[sendID]
	for i < msgNum {
		ws.Add(1)
		i++
		go func() {
			defer ws.Done()
			user.SendMsgWithContext(recvID, i)
		}()
	}
	ws.Wait()
}

func (p *PressureTester) importFriends(friendSenderUserIDs, recvMsgUserIDs []string) error {
	for _, recvMsgUserID := range recvMsgUserIDs {
		if err := p.friendManager.ImportFriends(recvMsgUserID, friendSenderUserIDs); err != nil {
			return err
		}
	}
	return nil
}

func (p *PressureTester) CheckMsg(ctx context.Context) {
	log.ZWarn(ctx, "message send finished checking", nil)
	var max, min, latencySum int64
	samepleReceiverFailedMap := make(map[string]*errorValue)
	failedMessageAllMap := make(map[string]*errorValue)
	sendSampleMessageAllMap := make(map[string]*msgValue)
	recvSampleMessageAllMap := make(map[string]*msgValue)
	groupSendSampleNum := make(map[string]int)
	groupSendFailedNum := make(map[string]int)
	groupRecvSampleInfo := make(map[string]*groupMessageValue)
	var sampleSendLength, sampleRecvLength, failedMessageLength int
	for _, user := range p.msgSender {
		if len(user.singleFailedMessageMap) != 0 {
			failedMessageLength += len(user.singleFailedMessageMap)
			for s, value := range user.singleFailedMessageMap {
				failedMessageAllMap[s] = value
				if utils.IsContain(value.RecvID, singleSampleUserList) {
					samepleReceiverFailedMap[s] = value
				}
			}
		}
		if len(user.singleSendSampleMessage) != 0 {
			sampleSendLength += len(user.singleSendSampleMessage)
			for s, value := range user.singleSendSampleMessage {
				sendSampleMessageAllMap[s] = value
			}
		}
		if len(user.singleRecvSampleMessage) != 0 {
			sampleRecvLength += len(user.singleRecvSampleMessage)
			for _, value := range user.singleRecvSampleMessage {
				if min == 0 && max == 0 {
					min = value.Latency
					max = value.Latency
				}
				if value.Latency < min {
					min = value.Latency
				}
				if value.Latency > max {
					max = value.Latency
				}
				latencySum += value.Latency
			}
			for s, value := range user.singleRecvSampleMessage {
				recvSampleMessageAllMap[s] = value
			}
		}

		//group check
		if len(user.groupSendSampleNum) != 0 {
			for group, num := range user.groupSendSampleNum {
				groupSendSampleNum[group] += num
			}
		}
		if len(user.groupFailedMessageMap) != 0 {
			for groupID, errInfo := range user.groupFailedMessageMap {
				groupSendFailedNum[groupID] += len(errInfo)
			}
		}
	}
	for groupID, ownerUserID := range p.groupOwnerUserID {
		if s, ok := p.msgSender[ownerUserID]; ok {
			if info, ok := s.groupRecvSampleInfo[groupID]; ok {
				info.Latency = info.LatencySum / info.Num
				groupRecvSampleInfo[groupID] = info
			}
		}
	}
	var latency string
	if sampleRecvLength != 0 {
		latency = utils.Int64ToString(latencySum/int64(sampleRecvLength)) + " ms"
	} else {
		latency = "0 ms"
	}
	log.ZWarn(context.Background(), "single message check result", nil, "failedMessageLength", failedMessageLength,
		"sampleSendLength", sampleSendLength, "sampleRecvLength", sampleRecvLength, "Average of message latency", latency,
		"max", utils.Int64ToString(max)+" ms",
		"min", utils.Int64ToString(min)+" ms")
	if len(groupSendSampleNum) > 0 {
		log.ZWarn(context.Background(), "group message check result", nil, "failedMessageLength", groupSendFailedNum,
			"sampleSendLength", groupSendSampleNum, "sampleRecvLength", groupRecvSampleInfo)
	}
	if len(failedMessageAllMap) > 0 {
		err := p.saveFailedMessageToFile(failedMessageAllMap, "failedMessageAllMap")
		if err != nil {
			log.ZWarn(ctx, "save failed message to file failed", err)
		}
	}

	if len(samepleReceiverFailedMap) > 0 {
		err := p.saveFailedMessageToFile(failedMessageAllMap, "sampleReceiverFailedMap")
		if err != nil {
			log.ZWarn(ctx, "save sampleReceiverFailedMap message to file failed", err)
		}
	}
	if sampleSendLength != sampleRecvLength {
		recvEx, sendEx := findMapIntersection(sendSampleMessageAllMap, recvSampleMessageAllMap)
		if len(recvEx) != 0 {
			p.saveSuccessMessageToFile(recvEx, "recvAdditional")
		}
		if len(sendEx) != 0 {
			p.saveSuccessMessageToFile(recvEx, "sendAdditional")
		}
	}
	log.ZDebug(ctx, "message send finished start to check message")
	os.Exit(1)
}

func (p *PressureTester) saveFailedMessageToFile(m map[string]*errorValue, filename string) error {
	file, err := os.Create(filename + ".txt")
	if err != nil {
		return err
	}
	defer file.Close()

	for key, value := range m {

		line := fmt.Sprintf("Key: %s, Value: %v\n", key, value)
		_, err := file.WriteString(line)
		if err != nil {
			return err
		}
	}
	return nil
}
func (p *PressureTester) saveSuccessMessageToFile(m map[string]*msgValue, filename string) error {
	file, err := os.Create(filename + ".txt")
	if err != nil {
		return err
	}
	defer file.Close()

	for key, value := range m {

		line := fmt.Sprintf("Key: %s, Value: %v\n", key, value)
		_, err := file.WriteString(line)
		if err != nil {
			return err
		}
	}
	return nil
}

func findMapIntersection(map1, map2 map[string]*msgValue) (map[string]*msgValue, map[string]*msgValue) {
	InMap1NotInMap2 := make(map[string]*msgValue)
	InMap2NotInMap1 := make(map[string]*msgValue)
	for key := range map1 {
		if _, ok := map2[key]; !ok {
			InMap1NotInMap2[key] = map1[key]
		}
	}
	for key := range map2 {
		if _, ok := map1[key]; !ok {
			InMap2NotInMap1[key] = map2[key]
		}
	}
	return InMap1NotInMap2, InMap2NotInMap1
}
