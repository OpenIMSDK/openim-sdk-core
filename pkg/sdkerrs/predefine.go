package sdkerrs

import "github.com/openimsdk/tools/errs"

var (
	// Common errors
	ErrArgs           = errs.NewCodeError(ArgsError, "Invalid input arguments")
	ErrCtxDeadline    = errs.NewCodeError(CtxDeadlineExceededError, "Context deadline exceeded")
	ErrSdkInternal    = errs.NewCodeError(SdkInternalError, "Internal SDK error")
	ErrNetwork        = errs.NewCodeError(NetworkError, "Network error")
	ErrNetworkTimeOut = errs.NewCodeError(NetworkTimeoutError, "Network timeout error")

	ErrGroupIDNotFound = errs.NewCodeError(GroupIDNotFoundError, "Group ID not found")
	ErrUserIDNotFound  = errs.NewCodeError(UserIDNotFoundError, "User ID not found")

	ErrResourceLoad = errs.NewCodeError(ResourceLoadNotCompleteError, "Resource initialization incomplete")

	// Message-related errors
	ErrMsgDecodeBinaryWs        = errs.NewCodeError(MsgDecodeBinaryWsError, "Message binary WebSocket decoding failed")
	ErrMsgDeCompression         = errs.NewCodeError(MsgDeCompressionError, "Message decompression failed")
	ErrMsgBinaryTypeNotSupport  = errs.NewCodeError(MsgBinaryTypeNotSupportError, "Message type not supported")
	ErrMsgRepeated              = errs.NewCodeError(MsgRepeatError, "Only failed messages can be resent")
	ErrMsgContentTypeNotSupport = errs.NewCodeError(MsgContentTypeNotSupportError, "Message content type not supported")
	ErrMsgHasNoSeq              = errs.NewCodeError(MsgHasNoSeqError, "Message has no sequence number")

	// Group-related errors
	ErrGroupType = errs.NewCodeError(GroupTypeErr, "Invalid group type")

	ErrLoginOut    = errs.NewCodeError(LoginOutError, "User has logged out")
	ErrLoginRepeat = errs.NewCodeError(LoginRepeatError, "User has logged in repeatedly")
)
