package test

import (
	"testing"
	"time"

	"github.com/openimsdk/openim-sdk-core/v3/open_im_sdk"
	"github.com/openimsdk/tools/log"
)

func Test_Empty(t *testing.T) {
	for {
		time.Sleep(time.Second * 1)
	}
}

func Test_ChangeInputState(t *testing.T) {
	for {
		err := open_im_sdk.UserForSDK.Conversation().ChangeInputStates(ctx, "sg_2309860938", true)
		if err != nil {
			log.ZError(ctx, "ChangeInputStates", err)
		}
		time.Sleep(time.Second * 1)
	}
}

func Test_RunWait(t *testing.T) {
	time.Sleep(time.Second * 10)
}

func Test_OnlineState(t *testing.T) {
	defer func() { select {} }()
	userIDs := []string{
		//"3611802798",
		"2110910952",
	}
	for i := 1; ; i++ {
		time.Sleep(time.Second)
		res, err := open_im_sdk.UserForSDK.LongConnMgr().GetUserOnlinePlatformIDs(ctx, userIDs)
		if err != nil {
			t.Logf("@@@@@@@@@@@@=====> <%d> error %s", i, err)
			continue
		}
		t.Logf("@@@@@@@@@@@@=====> <%d> success %+v", i, res)
	}
}
