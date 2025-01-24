package test

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/api"
	"github.com/openimsdk/openim-sdk-core/v3/sdk_struct"
	"github.com/openimsdk/protocol/auth"

	"github.com/openimsdk/openim-sdk-core/v3/open_im_sdk"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/ccontext"
)

var ctx context.Context

func init() {
	fmt.Println("------------------------>>>>>>>>>>>>>>>>>>> test init func <<<<<<<<<<<<<<<<<<<------------------------")
	rand.Seed(time.Now().UnixNano())
	listner := &OnConnListener{}
	config := getConf(APIADDR, WSADDR)
	config.DataDir = "./"
	configData, err := json.Marshal(config)
	if err != nil {
		panic(err)
	}
	isInit := open_im_sdk.InitSDK(listner, "test", string(configData))
	if !isInit {
		panic("init sdk failed")
	}
	ctx = open_im_sdk.UserForSDK.Context()
	ctx = ccontext.WithOperationID(ctx, "initOperationID_"+strconv.Itoa(int(time.Now().UnixMilli())))
	token, err := GetUserToken(ctx, UserID, PlatformID, Secret, config)
	if err != nil {
		panic(err)
	}
	if err := open_im_sdk.UserForSDK.Login(ctx, UserID, token); err != nil {
		panic(err)
	}
	ch := make(chan error)
	open_im_sdk.UserForSDK.SetConversationListener(&onConversationListener{ctx: ctx, ch: ch})
	open_im_sdk.UserForSDK.SetGroupListener(&onGroupListener{ctx: ctx})
	open_im_sdk.UserForSDK.SetAdvancedMsgListener(&onAdvancedMsgListener{ctx: ctx})
	open_im_sdk.UserForSDK.SetFriendshipListener(&onFriendshipListener{ctx: ctx})
	open_im_sdk.UserForSDK.SetUserListener(&onUserListener{ctx: ctx})
	if err := <-ch; err != nil {
		panic(err)
	}
}

func getConf(APIADDR, WSADDR string) sdk_struct.IMConfig {
	var cf sdk_struct.IMConfig
	cf.ApiAddr = APIADDR
	cf.WsAddr = WSADDR
	cf.DataDir = "./"
	cf.LogLevel = 6
	cf.IsExternalExtensions = true
	cf.PlatformID = PlatformID
	cf.LogFilePath = "./"
	cf.IsLogStandardOutput = true
	return cf
}

func GetUserToken(ctx context.Context, userID string, platformID int32, secret string, imConf sdk_struct.IMConfig) (string, error) {
	adminReq := &auth.GetAdminTokenReq{
		UserID: AdminUserID,
		Secret: secret,
	}
	adminToken, err := api.ExtractField(ctx, api.GetAdminToken.Invoke, adminReq, (*auth.GetAdminTokenResp).GetToken)
	if err != nil {
		return "", err
	}
	userReq := &auth.GetUserTokenReq{
		UserID:     userID,
		PlatformID: platformID,
	}
	ctx = ccontext.WithInfo(ctx, &ccontext.GlobalConfig{
		UserID:   userID,
		Token:    adminToken,
		IMConfig: imConf,
	})
	return api.ExtractField(ctx, api.GetUsersToken.Invoke, userReq, (*auth.GetUserTokenResp).GetToken)
}
