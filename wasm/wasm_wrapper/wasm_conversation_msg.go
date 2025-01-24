//go:build js && wasm

package wasm_wrapper

import (
	"syscall/js"

	"github.com/openimsdk/openim-sdk-core/v3/open_im_sdk"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/utils"
	"github.com/openimsdk/openim-sdk-core/v3/wasm/event_listener"
)

// ------------------------------------message---------------------------
type WrapperConMsg struct {
	*WrapperCommon
}

func NewWrapperConMsg(wrapperCommon *WrapperCommon) *WrapperConMsg {
	return &WrapperConMsg{WrapperCommon: wrapperCommon}
}

func (w *WrapperConMsg) CreateTextMessage(_ js.Value, args []js.Value) interface{} {
	return event_listener.NewCaller(open_im_sdk.CreateTextMessage, nil, &args).AsyncCallWithOutCallback()
}

func (w *WrapperConMsg) CreateImageMessage(_ js.Value, args []js.Value) interface{} {
	return event_listener.NewCaller(open_im_sdk.CreateImageMessage, nil, &args).AsyncCallWithOutCallback()
}

func (w *WrapperConMsg) CreateImageMessageByURL(_ js.Value, args []js.Value) interface{} {
	return event_listener.NewCaller(open_im_sdk.CreateImageMessageByURL, nil, &args).AsyncCallWithOutCallback()
}

func (w *WrapperConMsg) CreateSoundMessageByURL(_ js.Value, args []js.Value) interface{} {
	return event_listener.NewCaller(open_im_sdk.CreateSoundMessageByURL, nil, &args).AsyncCallWithOutCallback()
}

func (w *WrapperConMsg) CreateVideoMessageByURL(_ js.Value, args []js.Value) interface{} {
	return event_listener.NewCaller(open_im_sdk.CreateVideoMessageByURL, nil, &args).AsyncCallWithOutCallback()
}

func (w *WrapperConMsg) CreateFileMessageByURL(_ js.Value, args []js.Value) interface{} {
	return event_listener.NewCaller(open_im_sdk.CreateFileMessageByURL, nil, &args).AsyncCallWithOutCallback()
}

func (w *WrapperConMsg) CreateCustomMessage(_ js.Value, args []js.Value) interface{} {
	return event_listener.NewCaller(open_im_sdk.CreateCustomMessage, nil, &args).AsyncCallWithOutCallback()
}

func (w *WrapperConMsg) CreateQuoteMessage(_ js.Value, args []js.Value) interface{} {
	return event_listener.NewCaller(open_im_sdk.CreateQuoteMessage, nil, &args).AsyncCallWithOutCallback()
}

func (w *WrapperConMsg) CreateAdvancedQuoteMessage(_ js.Value, args []js.Value) interface{} {
	return event_listener.NewCaller(open_im_sdk.CreateAdvancedQuoteMessage, nil, &args).AsyncCallWithOutCallback()
}

func (w *WrapperConMsg) CreateAdvancedTextMessage(_ js.Value, args []js.Value) interface{} {
	return event_listener.NewCaller(open_im_sdk.CreateAdvancedTextMessage, nil, &args).AsyncCallWithOutCallback()
}

func (w *WrapperConMsg) CreateCardMessage(_ js.Value, args []js.Value) interface{} {
	return event_listener.NewCaller(open_im_sdk.CreateCardMessage, nil, &args).AsyncCallWithOutCallback()
}

func (w *WrapperConMsg) CreateTextAtMessage(_ js.Value, args []js.Value) interface{} {
	return event_listener.NewCaller(open_im_sdk.CreateTextAtMessage, nil, &args).AsyncCallWithOutCallback()
}

func (w *WrapperConMsg) CreateVideoMessage(_ js.Value, args []js.Value) interface{} {
	return event_listener.NewCaller(open_im_sdk.CreateVideoMessage, nil, &args).AsyncCallWithOutCallback()
}

func (w *WrapperConMsg) CreateFileMessage(_ js.Value, args []js.Value) interface{} {
	return event_listener.NewCaller(open_im_sdk.CreateFileMessage, nil, &args).AsyncCallWithOutCallback()
}

func (w *WrapperConMsg) CreateMergerMessage(_ js.Value, args []js.Value) interface{} {
	return event_listener.NewCaller(open_im_sdk.CreateMergerMessage, nil, &args).AsyncCallWithOutCallback()
}

func (w *WrapperConMsg) CreateFaceMessage(_ js.Value, args []js.Value) interface{} {
	return event_listener.NewCaller(open_im_sdk.CreateFaceMessage, nil, &args).AsyncCallWithOutCallback()
}

func (w *WrapperConMsg) CreateForwardMessage(_ js.Value, args []js.Value) interface{} {
	return event_listener.NewCaller(open_im_sdk.CreateForwardMessage, nil, &args).AsyncCallWithOutCallback()
}

func (w *WrapperConMsg) CreateLocationMessage(_ js.Value, args []js.Value) interface{} {
	return event_listener.NewCaller(open_im_sdk.CreateLocationMessage, nil, &args).AsyncCallWithOutCallback()
}

func (w *WrapperConMsg) CreateVideoMessageFromFullPath(_ js.Value, args []js.Value) interface{} {
	return event_listener.NewCaller(open_im_sdk.CreateVideoMessageFromFullPath, nil, &args).AsyncCallWithOutCallback()
}

func (w *WrapperConMsg) CreateImageMessageFromFullPath(_ js.Value, args []js.Value) interface{} {
	return event_listener.NewCaller(open_im_sdk.CreateImageMessageFromFullPath, nil, &args).AsyncCallWithOutCallback()
}

func (w *WrapperConMsg) CreateSoundMessageFromFullPath(_ js.Value, args []js.Value) interface{} {
	return event_listener.NewCaller(open_im_sdk.CreateSoundMessageFromFullPath, nil, &args).AsyncCallWithOutCallback()
}

func (w *WrapperConMsg) CreateFileMessageFromFullPath(_ js.Value, args []js.Value) interface{} {
	return event_listener.NewCaller(open_im_sdk.CreateFileMessageFromFullPath, nil, &args).AsyncCallWithOutCallback()
}

func (w *WrapperConMsg) CreateSoundMessage(_ js.Value, args []js.Value) interface{} {
	return event_listener.NewCaller(open_im_sdk.CreateSoundMessage, nil, &args).AsyncCallWithOutCallback()
}

func (w *WrapperConMsg) GetAtAllTag(_ js.Value, args []js.Value) interface{} {
	return event_listener.NewCaller(open_im_sdk.GetAtAllTag, nil, &args).AsyncCallWithOutCallback()
}

func (w *WrapperConMsg) MarkConversationMessageAsRead(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.MarkConversationMessageAsRead, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperConMsg) MarkAllConversationMessageAsRead(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.MarkAllConversationMessageAsRead, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperConMsg) MarkMessagesAsReadByMsgID(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.MarkMessagesAsReadByMsgID, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperConMsg) SendMessage(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewSendMessageCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc).SetClientMsgID(&args)
	return event_listener.NewCaller(open_im_sdk.SendMessage, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperConMsg) SendMessageNotOss(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewSendMessageCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc).SetClientMsgID(&args)
	return event_listener.NewCaller(open_im_sdk.SendMessageNotOss, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperConMsg) GetAllConversationList(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.GetAllConversationList, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperConMsg) GetConversationListSplit(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.GetConversationListSplit, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperConMsg) GetOneConversation(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.GetOneConversation, callback, &args).AsyncCallWithCallback()

}

func (w *WrapperConMsg) GetAdvancedHistoryMessageList(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.GetAdvancedHistoryMessageList, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperConMsg) GetAdvancedHistoryMessageListReverse(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.GetAdvancedHistoryMessageListReverse, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperConMsg) GetMultipleConversation(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.GetMultipleConversation, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperConMsg) FindMessageList(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.FindMessageList, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperConMsg) RevokeMessage(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.RevokeMessage, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperConMsg) TypingStatusUpdate(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.TypingStatusUpdate, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperConMsg) DeleteMessageFromLocalStorage(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.DeleteMessageFromLocalStorage, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperConMsg) DeleteMessage(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.DeleteMessage, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperConMsg) HideAllConversations(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.HideAllConversations, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperConMsg) DeleteAllMsgFromLocalAndSvr(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.DeleteAllMsgFromLocalAndSvr, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperConMsg) DeleteAllMsgFromLocal(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.DeleteAllMsgFromLocal, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperConMsg) ClearConversationAndDeleteAllMsg(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.ClearConversationAndDeleteAllMsg, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperConMsg) InsertSingleMessageToLocalStorage(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.InsertSingleMessageToLocalStorage, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperConMsg) InsertGroupMessageToLocalStorage(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.InsertGroupMessageToLocalStorage, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperConMsg) SearchLocalMessages(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.SearchLocalMessages, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperConMsg) SetMessageLocalEx(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.SetMessageLocalEx, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperConMsg) SearchConversation(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.SearchConversation, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperConMsg) DeleteConversationAndDeleteAllMsg(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.DeleteConversationAndDeleteAllMsg, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperConMsg) HideConversation(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.HideConversation, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperConMsg) SetConversation(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.SetConversation, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperConMsg) SetConversationDraft(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.SetConversationDraft, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperConMsg) GetTotalUnreadMsgCount(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.GetTotalUnreadMsgCount, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperConMsg) ChangeInputStates(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.ChangeInputStates, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperConMsg) GetInputStates(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.GetInputStates, callback, &args).AsyncCallWithCallback()
}
