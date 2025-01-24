package test

import (
	"github.com/openimsdk/openim-sdk-core/v3/open_im_sdk"
	"github.com/openimsdk/openim-sdk-core/v3/sdk_struct"
	"testing"
)

func Test_CreateTextMessage(t *testing.T) {
	message, err := open_im_sdk.UserForSDK.Conversation().CreateTextMessage(ctx, "textMsg")
	if err != nil {
		t.Error(err)
	}
	t.Log(message)
}

func Test_CreateAdvancedTextMessage(t *testing.T) {
	message, err := open_im_sdk.UserForSDK.Conversation().CreateAdvancedTextMessage(ctx, "textAdMsg", []*sdk_struct.MessageEntity{})
	if err != nil {
		t.Error(err)
	}
	t.Log(message)
}

func Test_CreateTextAtMessage(t *testing.T) {
	message, err := open_im_sdk.UserForSDK.Conversation().CreateTextAtMessage(ctx, "textATtsg", []string{}, []*sdk_struct.AtInfo{}, &sdk_struct.MsgStruct{})
	if err != nil {
		t.Error(err)
	}
	t.Log(message)
}

func Test_CreateQuoteMessage(t *testing.T) {
	message, err := open_im_sdk.UserForSDK.Conversation().CreateQuoteMessage(ctx, "textATtsg", &sdk_struct.MsgStruct{})
	if err != nil {
		t.Error(err)
	}
	t.Log(message)
}

func Test_CreateAdvancedQuoteMessage(t *testing.T) {
	message, err := open_im_sdk.UserForSDK.Conversation().CreateAdvancedQuoteMessage(ctx, "textATtsg", &sdk_struct.MsgStruct{}, []*sdk_struct.MessageEntity{})
	if err != nil {
		t.Error(err)
	}
	t.Log(message)
}

func Test_CreateVideoMessageFromFullPath(t *testing.T) {
	message, err := open_im_sdk.UserForSDK.Conversation().CreateVideoMessageFromFullPath(ctx, ".\\test.png", "mp4", 10, ".\\test.png")
	if err != nil {
		t.Error(err)
	}
	t.Log(message)
}

func Test_CreateCardMessage(t *testing.T) {
	message, err := open_im_sdk.UserForSDK.Conversation().CreateCardMessage(ctx, &sdk_struct.CardElem{
		UserID:   "123456",
		Nickname: "testname",
		FaceURL:  "",
	})
	if err != nil {
		t.Error(err)
	}
	t.Log(message)
}

func Test_CreateImageMessage(t *testing.T) {
	message, err := open_im_sdk.UserForSDK.Conversation().CreateImageMessage(ctx, ".\\test.png")
	if err != nil {
		t.Error(err)
	}
	t.Log(message)
}

func Test_CreateImageMessageByURL(t *testing.T) {
	message, err := open_im_sdk.UserForSDK.Conversation().CreateImageMessageByURL(ctx, "",
		sdk_struct.PictureBaseInfo{}, sdk_struct.PictureBaseInfo{}, sdk_struct.PictureBaseInfo{})
	if err != nil {
		t.Error(err)
	}
	t.Log(message)
}

func Test_CreateSoundMessageByURL(t *testing.T) {
	message, err := open_im_sdk.UserForSDK.Conversation().CreateSoundMessageByURL(ctx, &sdk_struct.SoundBaseInfo{})
	if err != nil {
		t.Error(err)
	}
	t.Log(message)
}

func Test_CreateSoundMessage(t *testing.T) {
	message, err := open_im_sdk.UserForSDK.Conversation().CreateSoundMessage(ctx, ".\\test.png", 20)
	if err != nil {
		t.Error(err)
	}
	t.Log(message)
}

func Test_CreateVideoMessage(t *testing.T) {
	message, err := open_im_sdk.UserForSDK.Conversation().CreateVideoMessage(ctx, ".\\test.png", "mp4", 10, "")
	if err != nil {
		t.Error(err)
	}
	t.Log(message)
}

func Test_CreateVideoMessageByURL(t *testing.T) {
	message, err := open_im_sdk.UserForSDK.Conversation().CreateVideoMessageByURL(ctx, sdk_struct.VideoBaseInfo{})
	if err != nil {
		t.Error(err)
	}
	t.Log(message)
}

func Test_CreateFileMessage(t *testing.T) {
	message, err := open_im_sdk.UserForSDK.Conversation().CreateFileMessage(ctx, ".\\test.png", "png")
	if err != nil {
		t.Error(err)
	}
	t.Log(message)
}

func Test_CreateFileMessageByURL(t *testing.T) {
	message, err := open_im_sdk.UserForSDK.Conversation().CreateFileMessageByURL(ctx, sdk_struct.FileBaseInfo{})
	if err != nil {
		t.Error(err)
	}
	t.Log(message)
}

func Test_CreateLocationMessage(t *testing.T) {
	message, err := open_im_sdk.UserForSDK.Conversation().CreateLocationMessage(ctx, "", 0, 0)
	if err != nil {
		t.Error(err)
	}
	t.Log(message)
}

func Test_CreateCustomMessage(t *testing.T) {
	message, err := open_im_sdk.UserForSDK.Conversation().CreateCustomMessage(ctx, "", "", "")
	if err != nil {
		t.Error(err)
	}
	t.Log(message)
}

func Test_CreateMergerMessage(t *testing.T) {
	message, err := open_im_sdk.UserForSDK.Conversation().CreateMergerMessage(ctx, []*sdk_struct.MsgStruct{}, "title", []string{"summary"})
	if err != nil {
		t.Error(err)
	}
	t.Log(message)
}

func Test_CreateFaceMessage(t *testing.T) {
	message, err := open_im_sdk.UserForSDK.Conversation().CreateFaceMessage(ctx, 0, "www.faceURL.com")
	if err != nil {
		t.Error(err)
	}
	t.Log(message)
}

func Test_CreateForwardMessage(t *testing.T) {
	message, err := open_im_sdk.UserForSDK.Conversation().CreateForwardMessage(ctx, &sdk_struct.MsgStruct{})
	if err != nil {
		t.Error(err)
	}
	t.Log(message)
}
