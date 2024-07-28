package sdk_user_simulator

import (
	"github.com/openimsdk/openim-sdk-core/v3/open_im_sdk"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/utils"
	"github.com/openimsdk/openim-sdk-core/v3/sdk_struct"
	"sync"
)

var (
	MapLock        sync.Mutex
	UserMessageMap = make(map[string]*MsgListenerCallBak)
	timeOffset     int64
)

func GetRelativeServerTime() int64 {
	return utils.GetCurrentTimestampByMill() + timeOffset
}

func InitSDK(userID, token string, cf sdk_struct.IMConfig) (*open_im_sdk.LoginMgr, error) {
	userForSDK := open_im_sdk.NewLoginMgr()
	var testConnListener testConnListener
	userForSDK.InitSDK(cf, &testConnListener)

	SetListener(userForSDK, userID)
	userForSDK.SetToken(userID, token)
	return userForSDK, nil
}

func SetListener(userForSDK *open_im_sdk.LoginMgr, userID string) {
	var testConversation conversationCallBack
	userForSDK.SetConversationListener(&testConversation)
	var testUser userCallback
	userForSDK.SetUserListener(testUser)

	msgCallBack := NewMsgListenerCallBak(userID)
	MapLock.Lock()
	UserMessageMap[userID] = msgCallBack
	MapLock.Unlock()
	userForSDK.SetAdvancedMsgListener(msgCallBack)

	var friendListener testFriendListener
	userForSDK.SetFriendListener(friendListener)

	var groupListener testGroupListener
	userForSDK.SetGroupListener(groupListener)
}
