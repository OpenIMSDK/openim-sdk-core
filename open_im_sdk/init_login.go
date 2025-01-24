package open_im_sdk

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/openimsdk/openim-sdk-core/v3/open_im_sdk_callback"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/ccontext"
	pbConstant "github.com/openimsdk/protocol/constant"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
	"github.com/openimsdk/openim-sdk-core/v3/sdk_struct"
	"github.com/openimsdk/openim-sdk-core/v3/version"

	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/mcontext"
)

func GetSdkVersion() string {
	return version.Version
}

const (
	rotateCount  uint = 1
	rotationTime uint = 24
)

func InitSDK(listener open_im_sdk_callback.OnConnListener, operationID string, config string) bool {
	if UserForSDK != nil {
		fmt.Println(operationID, "Initialize multiple times, use the existing ", UserForSDK, " Previous configuration ", UserForSDK.ImConfig(), " now configuration: ", config)
		return true
	}
	var configArgs sdk_struct.IMConfig
	if err := json.Unmarshal([]byte(config), &configArgs); err != nil {
		fmt.Println(operationID, "Unmarshal failed ", err.Error(), config)
		return false
	}
	if configArgs.PlatformID == 0 {
		return false
	}
	if err := log.InitLoggerFromConfig("open-im-sdk-core", "", configArgs.SystemType, pbConstant.PlatformID2Name[int(configArgs.PlatformID)], int(configArgs.LogLevel), configArgs.IsLogStandardOutput, false, configArgs.LogFilePath, rotateCount, rotationTime, version.Version, true); err != nil {
		fmt.Println(operationID, "log init failed ", err.Error())
	}
	fmt.Println("init log success")
	ctx := mcontext.NewCtx(operationID)
	if !strings.Contains(configArgs.ApiAddr, "http") {
		log.ZError(ctx, "api is http protocol, api format is invalid", nil)
		return false
	}
	if !strings.Contains(configArgs.WsAddr, "ws") {
		log.ZError(ctx, "ws is ws protocol, ws format is invalid", nil)
		return false
	}

	log.ZInfo(ctx, "InitSDK info", "config", configArgs)
	if listener == nil || config == "" {
		log.ZError(ctx, "listener or config is nil", nil)
		return false
	}
	UserForSDK = new(LoginMgr)
	return UserForSDK.InitSDK(configArgs, listener)
}
func UnInitSDK(operationID string) {
	if UserForSDK == nil {
		fmt.Println(operationID, "UserForSDK is nil,")
		return
	}
	UserForSDK.UnInitSDK()
	UserForSDK = nil
}

func Login(callback open_im_sdk_callback.Base, operationID string, userID, token string) {
	call(callback, operationID, UserForSDK.Login, userID, token)
}

func Logout(callback open_im_sdk_callback.Base, operationID string) {
	call(callback, operationID, UserForSDK.Logout)
}

func SetAppBackgroundStatus(callback open_im_sdk_callback.Base, operationID string, isBackground bool) {
	call(callback, operationID, UserForSDK.SetAppBackgroundStatus, isBackground)
}
func NetworkStatusChanged(callback open_im_sdk_callback.Base, operationID string) {
	call(callback, operationID, UserForSDK.NetworkStatusChanged)
}

func GetLoginStatus(operationID string) int {
	if UserForSDK == nil {
		return constant.Uninitialized
	}
	return UserForSDK.GetLoginStatus(ccontext.WithOperationID(context.Background(), operationID))
}

func GetLoginUserID() string {
	if UserForSDK == nil {
		return ""
	}
	return UserForSDK.GetLoginUserID()
}

func (u *LoginMgr) Login(ctx context.Context, userID, token string) error {
	return u.login(ctx, userID, token)
}

func (u *LoginMgr) Logout(ctx context.Context) error {
	return u.logout(ctx, false)
}

func (u *LoginMgr) SetAppBackgroundStatus(ctx context.Context, isBackground bool) error {
	return u.setAppBackgroundStatus(ctx, isBackground)
}
func (u *LoginMgr) NetworkStatusChanged(ctx context.Context) {
	u.longConnMgr.Close(ctx)
}

func (u *LoginMgr) GetLoginStatus(ctx context.Context) int {
	return u.getLoginStatus(ctx)
}
