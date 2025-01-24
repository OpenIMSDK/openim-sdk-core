package sdkerrs

// Common error codes
const (
	NetworkError             = 10000
	NetworkTimeoutError      = 10001
	ArgsError                = 10002 // Invalid input parameters
	CtxDeadlineExceededError = 10003 // Context deadline exceeded

	ResourceLoadNotCompleteError = 10004 // Resource initialization incomplete
	UnknownCode                  = 10005 // Unrecognized code
	SdkInternalError             = 10006 // SDK internal error

	UserIDNotFoundError = 10100 // UserID not found or not registered
	LoginOutError       = 10101 // User has logged out
	LoginRepeatError    = 10102 // User logged in repeatedly

	// Message-related errors
	MsgDeCompressionError         = 10201 // Message decompression failed
	MsgDecodeBinaryWsError        = 10202 // Message decoding failed
	MsgBinaryTypeNotSupportError  = 10203 // Message type not supported
	MsgRepeatError                = 10204 // Message repeated
	MsgContentTypeNotSupportError = 10205 // Message content type not supported
	MsgHasNoSeqError              = 10206 // Message does not have a sequence number

	// Group-related errors
	GroupIDNotFoundError = 10400 // GroupID not found
	GroupTypeErr         = 10401 // Invalid group type
)
