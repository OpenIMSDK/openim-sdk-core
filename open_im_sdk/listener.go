package open_im_sdk

import (
	"github.com/openimsdk/openim-sdk-core/v3/open_im_sdk_callback"
)

func SetGroupListener(listener open_im_sdk_callback.OnGroupListener) {
	listenerCall(UserForSDK.SetGroupListener, listener)
}

func SetConversationListener(listener open_im_sdk_callback.OnConversationListener) {
	listenerCall(UserForSDK.SetConversationListener, listener)
}

func SetAdvancedMsgListener(listener open_im_sdk_callback.OnAdvancedMsgListener) {
	listenerCall(UserForSDK.SetAdvancedMsgListener, listener)
}

func SetBatchMsgListener(listener open_im_sdk_callback.OnBatchMsgListener) {
	listenerCall(UserForSDK.SetBatchMsgListener, listener)
}

func SetUserListener(listener open_im_sdk_callback.OnUserListener) {
	listenerCall(UserForSDK.SetUserListener, listener)

}

func SetFriendListener(listener open_im_sdk_callback.OnFriendshipListener) {
	listenerCall(UserForSDK.SetFriendshipListener, listener)
}

func SetCustomBusinessListener(listener open_im_sdk_callback.OnCustomBusinessListener) {
	listenerCall(UserForSDK.SetCustomBusinessListener, listener)
}

func SetMessageKvInfoListener(listener open_im_sdk_callback.OnMessageKvInfoListener) {
	listenerCall(UserForSDK.SetMessageKvInfoListener, listener)
}
