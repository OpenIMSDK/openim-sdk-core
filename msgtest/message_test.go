package msgtest

import (
	"context"
	"testing"

	"github.com/openimsdk/openim-sdk-core/v3/msgtest/module"
	"github.com/openimsdk/openim-sdk-core/v3/msgtest/sdk_user_simulator"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/ccontext"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/utils"
	"github.com/openimsdk/tools/log"
)

func Test_SimulateMultiOnline(t *testing.T) {
	ctx := ccontext.WithOperationID(context.Background(), "TEST_ROOT")
	userIDList := []string{"1", "2"}
	metaManager := module.NewMetaManager(APIADDR, SECRET, MANAGERUSERID)
	userManager := metaManager.NewUserManager()
	serverTime, err := metaManager.GetServerTime()
	if err != nil {
		t.Fatal(err)
	}
	offset := serverTime - utils.GetCurrentTimestampByMill()
	sdk_user_simulator.SetServerTimeOffset(offset)
	for _, userID := range userIDList {
		token, err := userManager.GetToken(userID, int32(PLATFORMID))
		if err != nil {
			log.ZError(ctx, "get token failed, err: %v", err, "userID", userID)
			continue
		}
		err = sdk_user_simulator.InitSDKAndLogin(userID, token)
		if err != nil {
			log.ZError(ctx, "login failed, err: %v", err, "userID", userID)
		} else {
			log.ZDebug(ctx, "login success, userID: %v", "userID", userID)
		}
	}

}
