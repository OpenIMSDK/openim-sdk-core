package open_im_sdk

import (
	"github.com/openimsdk/openim-sdk-core/v3/open_im_sdk_callback"
)

func GetUsersInfo(callback open_im_sdk_callback.Base, operationID string, userIDs string) {
	call(callback, operationID, UserForSDK.User().GetUsersInfo, userIDs)
}

// SetSelfInfo sets the user's own information.
func SetSelfInfo(callback open_im_sdk_callback.Base, operationID string, userInfo string) {
	call(callback, operationID, UserForSDK.User().SetSelfInfo, userInfo)
}

// GetSelfUserInfo obtains the user's own information.
func GetSelfUserInfo(callback open_im_sdk_callback.Base, operationID string) {
	call(callback, operationID, UserForSDK.User().GetSelfUserInfo)
}

// AddUserCommand add to user's favorite
func AddUserCommand(callback open_im_sdk_callback.Base, operationID string, Type int32, uuid string, value string) {
	call(callback, operationID, UserForSDK.User().ProcessUserCommandAdd, Type, uuid, value)
}

// DeleteUserCommand delete from user's favorite
func DeleteUserCommand(callback open_im_sdk_callback.Base, operationID string, Type int32, uuid string) {
	call(callback, operationID, UserForSDK.User().ProcessUserCommandDelete, Type, uuid)
}

// GetAllUserCommands get user's favorite
func GetAllUserCommands(callback open_im_sdk_callback.Base, operationID string, Type int32) {
	call(callback, operationID, UserForSDK.User().ProcessUserCommandGetAll, Type)
}
