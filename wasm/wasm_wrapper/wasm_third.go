//go:build js && wasm

package wasm_wrapper

import (
	"syscall/js"

	"github.com/openimsdk/openim-sdk-core/v3/open_im_sdk"
	"github.com/openimsdk/openim-sdk-core/v3/open_im_sdk_callback"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/utils"
	"github.com/openimsdk/openim-sdk-core/v3/wasm/event_listener"
)

// ------------------------------------third---------------------------
type WrapperThird struct {
	*WrapperCommon
}

func NewWrapperThird(wrapperCommon *WrapperCommon) *WrapperThird {
	return &WrapperThird{WrapperCommon: wrapperCommon}
}
func (w *WrapperThird) UpdateFcmToken(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.UpdateFcmToken, callback, &args).AsyncCallWithCallback()
}
func (w *WrapperThird) UploadFile(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewUploadFileCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc).SetUuid(&args)
	return event_listener.NewCaller(UploadFile, callback, &args).AsyncCallWithCallback()
}

var _ open_im_sdk_callback.Base = (*TempBase)(nil)

type TempBase struct {
	u event_listener.UploadInterface
}

func NewTempBase(u event_listener.UploadInterface) *TempBase {
	return &TempBase{u: u}
}

func (t TempBase) OnError(errCode int32, errMsg string) {
	t.u.OnError(errCode, errMsg)
}

func (t TempBase) OnSuccess(data string) {
	t.u.OnSuccess(data)
}

var _ open_im_sdk_callback.UploadFileCallback = (*TempUploadFile)(nil)

type TempUploadFile struct {
	u event_listener.UploadInterface
}

func NewTempUploadFile(u event_listener.UploadInterface) *TempUploadFile {
	return &TempUploadFile{u: u}
}

func (t TempUploadFile) Open(size int64) {
	t.u.Open(size)
}

func (t TempUploadFile) PartSize(partSize int64, num int) {
	t.u.PartSize(partSize, num)
}

func (t TempUploadFile) HashPartProgress(index int, size int64, partHash string) {
	t.u.HashPartProgress(index, size, partHash)
}

func (t TempUploadFile) HashPartComplete(partsHash string, fileHash string) {
	t.u.HashPartComplete(partsHash, fileHash)
}

func (t TempUploadFile) UploadID(uploadID string) {
	t.u.UploadID(uploadID)
}

func (t TempUploadFile) UploadPartComplete(index int, partSize int64, partHash string) {
	t.u.UploadPartComplete(index, partSize, partHash)
}

func (t TempUploadFile) UploadComplete(fileSize int64, streamSize int64, storageSize int64) {
	t.u.UploadComplete(fileSize, streamSize, storageSize)
}

func (t TempUploadFile) Complete(size int64, url string, typ int) {
	t.u.Complete(size, url, typ)
}

func UploadFile(callback event_listener.UploadInterface, operationID string, req string) {
	b := NewTempBase(callback)
	t := NewTempUploadFile(callback)
	open_im_sdk.UploadFile(b, operationID, req, t)
}
