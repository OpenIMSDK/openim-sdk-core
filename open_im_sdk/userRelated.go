package open_im_sdk

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime/debug"
	"strings"
	"sync"
	"time"
	"unsafe"

	"github.com/openimsdk/openim-sdk-core/v3/internal/flagconst"
	"github.com/openimsdk/openim-sdk-core/v3/internal/third/file"

	"github.com/openimsdk/openim-sdk-core/v3/internal/relation"

	conv "github.com/openimsdk/openim-sdk-core/v3/internal/conversation_msg"
	"github.com/openimsdk/openim-sdk-core/v3/internal/group"
	"github.com/openimsdk/openim-sdk-core/v3/internal/interaction"
	"github.com/openimsdk/openim-sdk-core/v3/internal/third"
	"github.com/openimsdk/openim-sdk-core/v3/internal/user"
	"github.com/openimsdk/openim-sdk-core/v3/open_im_sdk_callback"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/ccontext"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/common"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/db_interface"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/sdkerrs"
	"github.com/openimsdk/openim-sdk-core/v3/sdk_struct"
	"github.com/openimsdk/protocol/push"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/utils/jsonutil"
)

const (
	LogoutStatus = iota + 1
	Logging
	Logged
)

const (
	LogoutTips = "js sdk socket close"
)

var (
	// UserForSDK Client-independent user class
	UserForSDK *LoginMgr
)

// CheckResourceLoad checks the SDK is resource load status.
func CheckResourceLoad(uSDK *LoginMgr, funcName string) error {
	if uSDK == nil {
		return errs.New("SDK not initialized,userForSDK is nil", "funcName", funcName).Wrap()
	}
	if funcName == "" {
		return nil
	}
	parts := strings.Split(funcName, ".")
	if parts[len(parts)-1] == "Login-fm" {
		return nil
	}
	if uSDK.getLoginStatus(context.Background()) != Logged {
		return errs.New("SDK not logged in", "funcName", funcName).Wrap()
	}
	return nil
}

type LoginMgr struct {
	relation     *relation.Relation
	group        *group.Group
	conversation *conv.Conversation
	user         *user.User
	file         *file.File

	db           db_interface.DataBase
	longConnMgr  *interaction.LongConnMgr
	msgSyncer    *interaction.MsgSyncer
	third        *third.Third
	token        string
	loginUserID  string
	connListener open_im_sdk_callback.OnConnListener

	justOnceFlag bool

	w           sync.Mutex
	loginStatus int

	groupListener        open_im_sdk_callback.OnGroupListener
	friendshipListener   open_im_sdk_callback.OnFriendshipListener
	conversationListener open_im_sdk_callback.OnConversationListener
	advancedMsgListener  open_im_sdk_callback.OnAdvancedMsgListener
	batchMsgListener     open_im_sdk_callback.OnBatchMsgListener
	userListener         open_im_sdk_callback.OnUserListener
	signalingListener    open_im_sdk_callback.OnSignalingListener
	businessListener     open_im_sdk_callback.OnCustomBusinessListener
	msgKvListener        open_im_sdk_callback.OnMessageKvInfoListener

	conversationCh     chan common.Cmd2Value
	cmdWsCh            chan common.Cmd2Value
	pushMsgAndMaxSeqCh chan common.Cmd2Value
	loginMgrCh         chan common.Cmd2Value

	ctx       context.Context
	cancel    context.CancelFunc
	info      *ccontext.GlobalConfig
	id2MinSeq map[string]int64
}

func (u *LoginMgr) GroupListener() open_im_sdk_callback.OnGroupListener {
	return u.groupListener
}

func (u *LoginMgr) FriendshipListener() open_im_sdk_callback.OnFriendshipListener {
	return u.friendshipListener
}

func (u *LoginMgr) ConversationListener() open_im_sdk_callback.OnConversationListener {
	return u.conversationListener
}

func (u *LoginMgr) AdvancedMsgListener() open_im_sdk_callback.OnAdvancedMsgListener {
	return u.advancedMsgListener
}

func (u *LoginMgr) BatchMsgListener() open_im_sdk_callback.OnBatchMsgListener {
	return u.batchMsgListener
}

func (u *LoginMgr) UserListener() open_im_sdk_callback.OnUserListener {
	return u.userListener
}

func (u *LoginMgr) SignalingListener() open_im_sdk_callback.OnSignalingListener {
	return u.signalingListener
}

func (u *LoginMgr) BusinessListener() open_im_sdk_callback.OnCustomBusinessListener {
	return u.businessListener
}

func (u *LoginMgr) MsgKvListener() open_im_sdk_callback.OnMessageKvInfoListener {
	return u.msgKvListener
}

func (u *LoginMgr) Exit() {
	u.cancel()
}

func (u *LoginMgr) GetToken() string {
	return u.token
}

func (u *LoginMgr) Third() *third.Third {
	return u.third
}

func (u *LoginMgr) ImConfig() sdk_struct.IMConfig {
	return sdk_struct.IMConfig{
		PlatformID:           u.info.PlatformID,
		ApiAddr:              u.info.ApiAddr,
		WsAddr:               u.info.WsAddr,
		DataDir:              u.info.DataDir,
		LogLevel:             u.info.LogLevel,
		IsExternalExtensions: u.info.IsExternalExtensions,
	}
}

func (u *LoginMgr) Conversation() *conv.Conversation {
	return u.conversation
}

func (u *LoginMgr) User() *user.User {
	return u.user
}

func (u *LoginMgr) File() *file.File {
	return u.file
}

func (u *LoginMgr) Group() *group.Group {
	return u.group
}

func (u *LoginMgr) Relation() *relation.Relation {
	return u.relation
}

func (u *LoginMgr) SetConversationListener(conversationListener open_im_sdk_callback.OnConversationListener) {
	u.conversationListener = conversationListener
}

func (u *LoginMgr) SetAdvancedMsgListener(advancedMsgListener open_im_sdk_callback.OnAdvancedMsgListener) {
	u.advancedMsgListener = advancedMsgListener
}

func (u *LoginMgr) SetMessageKvInfoListener(messageKvInfoListener open_im_sdk_callback.OnMessageKvInfoListener) {
	u.msgKvListener = messageKvInfoListener
}

func (u *LoginMgr) SetBatchMsgListener(batchMsgListener open_im_sdk_callback.OnBatchMsgListener) {
	u.batchMsgListener = batchMsgListener
}

func (u *LoginMgr) SetFriendshipListener(friendshipListener open_im_sdk_callback.OnFriendshipListener) {
	u.friendshipListener = friendshipListener
}

func (u *LoginMgr) SetGroupListener(groupListener open_im_sdk_callback.OnGroupListener) {
	u.groupListener = groupListener
}

func (u *LoginMgr) SetUserListener(userListener open_im_sdk_callback.OnUserListener) {
	u.userListener = userListener
}

func (u *LoginMgr) SetCustomBusinessListener(listener open_im_sdk_callback.OnCustomBusinessListener) {
	u.businessListener = listener
}
func (u *LoginMgr) GetLoginUserID() string {
	return u.loginUserID
}
func (u *LoginMgr) logoutListener(ctx context.Context) {
	defer func() {
		if r := recover(); r != nil {
			err := fmt.Sprintf("panic: %+v\n%s", r, debug.Stack())

			log.ZWarn(ctx, "logoutListener panic", nil, "panic info", err)
		}
	}()

	for {
		select {
		case <-u.loginMgrCh:
			log.ZDebug(ctx, "logoutListener exit")
			err := u.logout(ctx, true)
			if err != nil {
				log.ZError(ctx, "logout error", err)
			}
		case <-ctx.Done():
			log.ZInfo(ctx, "logoutListener done sdk logout.....")
			return
		}
	}

}

func NewLoginMgr() *LoginMgr {
	return &LoginMgr{
		info: &ccontext.GlobalConfig{},
	}
}

func (u *LoginMgr) getLoginStatus(_ context.Context) int {
	u.w.Lock()
	defer u.w.Unlock()
	return u.loginStatus
}

func (u *LoginMgr) setLoginStatus(status int) {
	u.w.Lock()
	defer u.w.Unlock()
	u.loginStatus = status
}

func (u *LoginMgr) checkSendingMessage(ctx context.Context) {
	sendingMessages, err := u.db.GetAllSendingMessages(ctx)
	if err != nil {
		log.ZError(ctx, "GetAllSendingMessages failed", err)
	}
	for _, message := range sendingMessages {
		if err := u.handlerSendingMsg(ctx, message); err != nil {
			log.ZError(ctx, "handlerSendingMsg failed", err, "message", message)
		}
		if err := u.db.DeleteSendingMessage(ctx, message.ConversationID, message.ClientMsgID); err != nil {
			log.ZError(ctx, "DeleteSendingMessage failed", err, "conversationID", message.ConversationID, "clientMsgID", message.ClientMsgID)
		}
	}
}

func (u *LoginMgr) handlerSendingMsg(ctx context.Context, sendingMsg *model_struct.LocalSendingMessages) error {
	tableMessage, err := u.db.GetMessage(ctx, sendingMsg.ConversationID, sendingMsg.ClientMsgID)
	if err != nil {
		return err
	}
	if tableMessage.Status != constant.MsgStatusSending {
		return nil
	}
	err = u.db.UpdateMessage(ctx, sendingMsg.ConversationID, &model_struct.LocalChatLog{ClientMsgID: sendingMsg.ClientMsgID, Status: constant.MsgStatusSendFailed})
	if err != nil {
		return err
	}
	conversation, err := u.db.GetConversation(ctx, sendingMsg.ConversationID)
	if err != nil {
		return err
	}
	latestMsg := &sdk_struct.MsgStruct{}
	if err := json.Unmarshal([]byte(conversation.LatestMsg), &latestMsg); err != nil {
		return err
	}
	if latestMsg.ClientMsgID == sendingMsg.ClientMsgID {
		latestMsg.Status = constant.MsgStatusSendFailed
		conversation.LatestMsg = jsonutil.StructToJsonString(latestMsg)
		return u.db.UpdateConversation(ctx, conversation)
	}
	return nil
}

func (u *LoginMgr) login(ctx context.Context, userID, token string) error {
	if u.getLoginStatus(ctx) == Logged {
		return sdkerrs.ErrLoginRepeat
	}
	u.setLoginStatus(Logging)
	log.ZDebug(ctx, "login start... ", "userID", userID, "token", token)
	t1 := time.Now()

	u.info.UserID = userID
	u.info.Token = token
	u.token = token
	u.loginUserID = userID
	var err error
	u.db, err = db.NewDataBase(ctx, userID, u.info.DataDir, int(u.info.LogLevel))
	if err != nil {
		return sdkerrs.ErrSdkInternal.WrapMsg("init database " + err.Error())
	}
	u.checkSendingMessage(ctx)
	log.ZDebug(ctx, "NewDataBase ok", "userID", userID, "dataDir", u.info.DataDir, "login cost time", time.Since(t1))
	u.user = user.NewUser(u.db, u.loginUserID, u.conversationCh)
	u.file = file.NewFile(u.db, u.loginUserID)
	u.relation = relation.NewFriend(u.loginUserID, u.db, u.user, u.conversationCh)

	u.group = group.NewGroup(u.loginUserID, u.db, u.conversationCh)
	u.third = third.NewThird(u.info.PlatformID, u.loginUserID, u.info.SystemType, u.info.LogFilePath, u.file)
	log.ZDebug(ctx, "forcedSynchronization success...", "login cost time: ", time.Since(t1))

	u.msgSyncer, _ = interaction.NewMsgSyncer(ctx, u.conversationCh, u.pushMsgAndMaxSeqCh, u.loginUserID, u.longConnMgr, u.db, 0)
	u.conversation = conv.NewConversation(ctx, u.longConnMgr, u.db, u.conversationCh,
		u.relation, u.group, u.user, u.file)
	u.setListener(ctx)

	u.run(ctx)
	u.setLoginStatus(Logged)
	log.ZDebug(ctx, "login success...", "login cost time: ", time.Since(t1))
	return nil
}

func (u *LoginMgr) setListener(ctx context.Context) {
	setListener(ctx, &u.userListener, u.UserListener, u.user.SetListener, newEmptyUserListener)
	setListener(ctx, &u.friendshipListener, u.FriendshipListener, u.relation.SetListener, newEmptyFriendshipListener)
	setListener(ctx, &u.groupListener, u.GroupListener, u.group.SetGroupListener, newEmptyGroupListener)
	setListener(ctx, &u.conversationListener, u.ConversationListener, u.conversation.SetConversationListener, newEmptyConversationListener)
	setListener(ctx, &u.advancedMsgListener, u.AdvancedMsgListener, u.conversation.SetMsgListener, newEmptyAdvancedMsgListener)
	setListener(ctx, &u.batchMsgListener, u.BatchMsgListener, u.conversation.SetBatchMsgListener, nil)
	setListener(ctx, &u.businessListener, u.BusinessListener, u.conversation.SetBusinessListener, newEmptyCustomBusinessListener)
}

func setListener[T any](ctx context.Context, listener *T, getter func() T, setFunc func(listener func() T), newFunc func(context.Context) T) {
	if *(*unsafe.Pointer)(unsafe.Pointer(listener)) == nil && newFunc != nil {
		*listener = newFunc(ctx)
	}
	setFunc(getter)
}

func (u *LoginMgr) run(ctx context.Context) {
	u.longConnMgr.Run(ctx)
	go u.msgSyncer.DoListener(ctx)
	go common.DoListener(u.ctx, u.conversation)
	go u.logoutListener(ctx)
}

func (u *LoginMgr) InitSDK(config sdk_struct.IMConfig, listener open_im_sdk_callback.OnConnListener) bool {
	if listener == nil {
		return false
	}
	u.info = &ccontext.GlobalConfig{}
	u.info.IMConfig = config
	u.connListener = listener
	u.initResources()
	return true
}

func (u *LoginMgr) Context() context.Context {
	return u.ctx
}

func (u *LoginMgr) initResources() {
	ctx := ccontext.WithInfo(context.Background(), u.info)
	u.ctx, u.cancel = context.WithCancel(ctx)
	var convChanLen int
	if flagconst.TestMode {
		convChanLen = 100000
	} else {
		convChanLen = 1000
	}
	u.conversationCh = make(chan common.Cmd2Value, convChanLen)
	u.pushMsgAndMaxSeqCh = make(chan common.Cmd2Value, 1000)
	u.loginMgrCh = make(chan common.Cmd2Value, 1)
	u.longConnMgr = interaction.NewLongConnMgr(u.ctx, u.connListener, u.userOnlineStatusChange, u.pushMsgAndMaxSeqCh, u.loginMgrCh)
	u.ctx = ccontext.WithApiErrCode(u.ctx, &apiErrCallback{loginMgrCh: u.loginMgrCh, listener: u.connListener})
	u.setLoginStatus(LogoutStatus)
}

func (u *LoginMgr) userOnlineStatusChange(users map[string][]int32) {
	u.User().UserOnlineStatusChange(users)
}

func (u *LoginMgr) UnInitSDK() {
	if u.getLoginStatus(context.Background()) == Logged {
		fmt.Println("sdk not logout, please logout first")
		return
	}
	u.info = nil
	u.setLoginStatus(0)
}

// token error recycle recourse, kicked not recycle
func (u *LoginMgr) logout(ctx context.Context, isTokenValid bool) error {
	if ccontext.Info(ctx).OperationID() == LogoutTips {
		isTokenValid = true
	}
	if !isTokenValid {
		ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
		defer cancel()
		err := u.longConnMgr.SendReqWaitResp(ctx, &push.DelUserPushTokenReq{UserID: u.info.UserID, PlatformID: u.info.PlatformID}, constant.LogoutMsg, &push.DelUserPushTokenResp{})
		if err != nil {
			log.ZWarn(ctx, "TriggerCmdLogout server recycle resources failed...", err)
		} else {
			log.ZDebug(ctx, "TriggerCmdLogout server recycle resources success...")
		}
	}
	u.Exit()
	err := u.db.Close(u.ctx)
	if err != nil {
		log.ZWarn(ctx, "TriggerCmdLogout db recycle resources failed...", err)
	}
	// user object must be rest  when user logout
	u.initResources()
	log.ZDebug(ctx, "TriggerCmdLogout client success...",
		"isTokenValid", isTokenValid)
	return nil
}

func (u *LoginMgr) setAppBackgroundStatus(ctx context.Context, isBackground bool) error {
	if u.longConnMgr.GetConnectionStatus() == interaction.DefaultNotConnect {
		u.longConnMgr.SetBackground(isBackground)
		return nil
	}
	var resp sdkws.SetAppBackgroundStatusResp
	err := u.longConnMgr.SendReqWaitResp(ctx, &sdkws.SetAppBackgroundStatusReq{UserID: u.loginUserID, IsBackground: isBackground}, constant.SetBackgroundStatus, &resp)
	if err != nil {
		return err
	} else {
		u.longConnMgr.SetBackground(isBackground)
		if !isBackground {
			_ = common.TriggerCmdWakeUpDataSync(ctx, u.pushMsgAndMaxSeqCh)
		}

		return nil
	}
}

func (u *LoginMgr) LongConnMgr() *interaction.LongConnMgr {
	return u.longConnMgr
}
