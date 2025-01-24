package open_im_sdk

import (
	"github.com/openimsdk/openim-sdk-core/v3/internal/third/file"
	"github.com/openimsdk/openim-sdk-core/v3/open_im_sdk_callback"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/sdkerrs"
)

func UpdateFcmToken(callback open_im_sdk_callback.Base, operationID, fcmToken string, expireTime int64) {
	call(callback, operationID, UserForSDK.Third().UpdateFcmToken, fcmToken, expireTime)
}

func SetAppBadge(callback open_im_sdk_callback.Base, operationID string, appUnreadCount int32) {
	call(callback, operationID, UserForSDK.Third().SetAppBadge, appUnreadCount)
}

func UploadLogs(callback open_im_sdk_callback.Base, operationID string, line int, ex string, progress open_im_sdk_callback.UploadLogProgress) {
	call(callback, operationID, UserForSDK.Third().UploadLogs, line, ex, progress)
}

func Logs(callback open_im_sdk_callback.Base, operationID string, logLevel int, file string, line int, msgs string, err string, keyAndValue string) {
	if UserForSDK == nil || UserForSDK.Third() == nil {
		callback.OnError(sdkerrs.SdkInternalError, "sdk not init")
		return
	}
	call(callback, operationID, UserForSDK.Third().Log, logLevel, file, line, msgs, err, keyAndValue)
}

func UploadFile(callback open_im_sdk_callback.Base, operationID string, req string, progress open_im_sdk_callback.UploadFileCallback) {
	call(callback, operationID, UserForSDK.File().UploadFile, req, file.UploadFileCallback(progress))
}
