package test

import (
	"context"
	"testing"
	"time"

	"github.com/openimsdk/openim-sdk-core/v3/open_im_sdk"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/cache"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/sdk_params_callback"
	"github.com/openimsdk/openim-sdk-core/v3/sdk_struct"
	"github.com/openimsdk/protocol/conversation"
	"github.com/openimsdk/tools/log"
)

func Test_GetAllConversationList(t *testing.T) {
	conversations, err := open_im_sdk.UserForSDK.Conversation().GetAllConversationList(ctx)
	if err != nil {
		t.Fatal(err)
	}
	for _, conversation := range conversations {
		t.Log(conversation)
	}
	t.Log(len(conversations))
	time.Sleep(time.Second * 100)
}

func Test_GetConversationListSplit(t *testing.T) {
	conversations, err := open_im_sdk.UserForSDK.Conversation().GetConversationListSplit(ctx, 0, 20)
	if err != nil {
		t.Fatal(err)
	}
	for _, conversation := range conversations {
		t.Log(conversation)
	}
}

func Test_HideConversation(t *testing.T) {
	err := open_im_sdk.UserForSDK.Conversation().HideConversation(ctx, "asdasd")
	if err != nil {
		t.Fatal(err)
	}
}

func Test_GetGlobalRecvMessageOpt(t *testing.T) {
	opt, err := open_im_sdk.UserForSDK.Conversation().GetOneConversation(ctx, 2, "1772958501")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(*opt)
}

func Test_GetGetMultipleConversation(t *testing.T) {
	conversations, err := open_im_sdk.UserForSDK.Conversation().GetMultipleConversation(ctx, []string{"asdasd"})
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range conversations {
		t.Log(v)
	}
}

func Test_SetConversationDraft(t *testing.T) {
	err := open_im_sdk.UserForSDK.Conversation().SetConversationDraft(ctx, "group_17729585012", "draft")
	if err != nil {
		t.Fatal(err)
	}
}

func Test_SetConversation(t *testing.T) {
	err := open_im_sdk.UserForSDK.Conversation().SetConversation(ctx, "group_17729585012", &conversation.ConversationReq{})
	if err != nil {
		t.Fatal(err)
	}
}

func Test_GetTotalUnreadMsgCount(t *testing.T) {
	count, err := open_im_sdk.UserForSDK.Conversation().GetTotalUnreadMsgCount(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(count)
}

func Test_SendMessage(t *testing.T) {
	ctx = context.WithValue(ctx, "callback", TestSendMsg{})
	msg, _ := open_im_sdk.UserForSDK.Conversation().CreateTextMessage(ctx, "textMsg")
	_, err := open_im_sdk.UserForSDK.Conversation().SendMessage(ctx, msg, "3411008330", "", nil, false)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_SendMessageNotOss(t *testing.T) {
	ctx = context.WithValue(ctx, "callback", TestSendMsg{})
	msg, _ := open_im_sdk.UserForSDK.Conversation().CreateTextMessage(ctx, "textMsg")
	_, err := open_im_sdk.UserForSDK.Conversation().SendMessageNotOss(ctx, msg, "3411008330", "", nil, false)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_FindMessageList(t *testing.T) {
	msgs, err := open_im_sdk.UserForSDK.Conversation().FindMessageList(ctx, []*sdk_params_callback.ConversationArgs{})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(msgs.TotalCount)
	for _, v := range msgs.FindResultItems {
		t.Log(v)
	}
}

func Test_GetAdvancedHistoryMessageList(t *testing.T) {
	msgs, err := open_im_sdk.UserForSDK.Conversation().GetAdvancedHistoryMessageList(ctx, sdk_params_callback.GetAdvancedHistoryMessageListParams{})
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range msgs.MessageList {
		t.Log(v)
	}
}

func Test_GetAdvancedHistoryMessageListReverse(t *testing.T) {
	msgs, err := open_im_sdk.UserForSDK.Conversation().GetAdvancedHistoryMessageListReverse(ctx, sdk_params_callback.GetAdvancedHistoryMessageListParams{
		ConversationID:   "si_3325086438_5054969402",
		StartClientMsgID: "91e40552a05a60494a56e86d36c497ce",
		Count:            20,
		ViewType:         cache.ViewHistory,
	})
	if err != nil {
		t.Fatal(err)
	}
	log.ZDebug(context.Background(), "GetAdvancedHistoryMessageListReverse Resp", "resp", msgs)
	for _, v := range msgs.MessageList {
		t.Log(v)
	}
}

func Test_InsertSingleMessageToLocalStorage(t *testing.T) {
	_, err := open_im_sdk.UserForSDK.Conversation().InsertSingleMessageToLocalStorage(ctx, &sdk_struct.MsgStruct{}, "3411008330", "")
	if err != nil {
		t.Fatal(err)
	}
}

func Test_InsertGroupMessageToLocalStorage(t *testing.T) {
	_, err := open_im_sdk.UserForSDK.Conversation().InsertGroupMessageToLocalStorage(ctx, &sdk_struct.MsgStruct{}, "group_17729585012", "")
	if err != nil {
		t.Fatal(err)
	}
}

func Test_SearchLocalMessages(t *testing.T) {
	req := &sdk_params_callback.SearchLocalMessagesParams{
		Count:            20,
		KeywordList:      []string{"1"},
		MessageTypeList:  []int{105},
		PageIndex:        1,
		SenderUserIDList: []string{},
	}
	msgs, err := open_im_sdk.UserForSDK.Conversation().SearchLocalMessages(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range msgs.SearchResultItems {
		t.Log(v)
	}
}

func Test_SetMessageLocalEx(t *testing.T) {
	err := open_im_sdk.UserForSDK.Conversation().SetMessageLocalEx(ctx, "si_2975755104_6386894923", "53ca4b3be29f7ea231a5e82e7af8a43f", "{key,value}")
	if err != nil {
		t.Fatal(err)
	}
}

func Test_DeleteAllMsgFromLocalAndSvr(t *testing.T) {
	err := open_im_sdk.UserForSDK.Conversation().DeleteAllMsgFromLocalAndServer(ctx)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_DeleteAllMessageFromLocalStorage(t *testing.T) {
	err := open_im_sdk.UserForSDK.Conversation().DeleteAllMessageFromLocalStorage(ctx)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_DeleteMessage(t *testing.T) {
	err := open_im_sdk.UserForSDK.Conversation().DeleteMessage(ctx, "si_1695766238_5099153716", "8b67803979bce9c6daf82fb64dbffc5f")
	if err != nil {
		t.Fatal(err)
	}
}

func Test_ClearConversationAndDeleteAllMsg(t *testing.T) {
	err := open_im_sdk.UserForSDK.Conversation().ClearConversationAndDeleteAllMsg(ctx, "si_3271407977_7152307910")
	if err != nil {
		t.Fatal(err)
	}
}

func Test_MarkConversationMessageAsRead(t *testing.T) {
	err := open_im_sdk.UserForSDK.Conversation().MarkConversationMessageAsRead(ctx, "si_2688118337_7249315132")
	if err != nil {
		t.Fatal(err)
	}
}

func Test_MarkAllConversationMessageAsRead(t *testing.T) {
	err := open_im_sdk.UserForSDK.Conversation().MarkAllConversationMessageAsRead(ctx)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_MarkMsgsAsRead(t *testing.T) {
	err := open_im_sdk.UserForSDK.Conversation().MarkMessagesAsReadByMsgID(ctx, "si_2688118337_7249315132",
		[]string{"fb56ed151b675e0837ed3af79dbf66b1",
			"635715c539be2e7812a0fc802f0cdc54", "1aba3fae3dc3f61c17e8eb09519cf8e1"})
	if err != nil {
		t.Fatal(err)
	}
}

func Test_SendImgMsg(t *testing.T) {
	ctx = context.WithValue(ctx, "callback", TestSendMsg{})
	msg, err := open_im_sdk.UserForSDK.Conversation().CreateImageMessage(ctx, "C:\\Users\\Admin\\Desktop\\test.png")
	if err != nil {
		t.Fatal(err)
	}
	res, err := open_im_sdk.UserForSDK.Conversation().SendMessage(ctx, msg, "1919501984", "", nil, false)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("send smg => %+v\n", res)
}

func Test_SearchConversation(t *testing.T) {
	result, err := open_im_sdk.UserForSDK.Conversation().SearchConversation(ctx, "a")
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range result {
		t.Log(v)
	}
}
