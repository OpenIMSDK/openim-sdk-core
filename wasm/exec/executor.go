//go:build js && wasm

package exec

import (
	"context"
	"errors"
	"runtime"
	"syscall/js"
	"time"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/utils"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
)

type CallbackData struct {
	ErrCode int32       `json:"errCode"`
	ErrMsg  string      `json:"errMsg"`
	Data    interface{} `json:"data"`
}

const TIMEOUT = 5
const JSNOTFOUND = 10002

var ErrType = errors.New("from javascript data type err")
var PrimaryKeyNull = errors.New("primary key is null err")

var ErrTimoutFromJavaScript = errors.New("invoke javascript timeout, maybe should check  function from javascript")
var jsErr = js.Global().Get("Error")

func Exec(args ...interface{}) (output interface{}, err error) {
	ctx := context.Background()
	defer func() {
		if r := recover(); r != nil {
			switch x := r.(type) {
			case string:
				err = errs.WrapMsg(errors.New(x), "")
			case error:
				err = x
			default:
				err = errs.WrapMsg(errors.New("unknown panic"), "")
			}
		}
	}()
	thenChannel := make(chan []js.Value)
	defer close(thenChannel)
	catchChannel := make(chan []js.Value)
	defer close(catchChannel)
	pc, _, _, _ := runtime.Caller(1)
	funcName := utils.CleanUpfuncName(runtime.FuncForPC(pc).Name())
	data := CallbackData{}
	thenFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		defer func() {
			if r := recover(); r != nil {
				switch x := r.(type) {
				case string:
					err = errs.WrapMsg(errors.New(x), "")
				case error:
					err = x
				default:
					err = errs.WrapMsg(errors.New("unknown panic"), "")
				}
			}
		}()
		log.ZDebug(ctx, "js then function", "=> (main go context) "+funcName+" "+
			"with response ", args[0].String())
		thenChannel <- args
		return nil
	})
	defer thenFunc.Release()
	catchFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		defer func() {
			if r := recover(); r != nil {
				switch x := r.(type) {
				case string:
					err = errs.WrapMsg(errors.New(x), "")
				case error:
					err = x
				default:
					err = errs.WrapMsg(errors.New("unknown panic"), "")
				}
			}
		}()
		log.ZDebug(ctx, "js catch function", "=> (main go context) "+funcName+" with respone ", args[0].String())
		catchChannel <- args
		return nil
	})
	defer catchFunc.Release()
	js.Global().Call(utils.FirstLower(funcName), args...).Call("then", thenFunc).Call("catch", catchFunc)
	select {
	case result := <-thenChannel:
		if len(result) > 0 {
			switch result[0].Type() {
			case js.TypeString:
				interErr := utils.JsonStringToStruct(result[0].String(), &data)
				if interErr != nil {
					err = errs.WrapMsg(err, "return json unmarshal err from javascript")
				}
			case js.TypeObject:
				return result[0], nil

			default:
				err = errors.New("unknown return type from javascript")
			}

		} else {
			err = errors.New("args err,length is 0")
		}

	case catch := <-catchChannel:
		if catch[0].InstanceOf(jsErr) {
			return nil, js.Error{Value: catch[0]}
		} else {
			panic("unknown javascript exception")
		}
	case <-time.After(TIMEOUT * time.Second):
		panic(ErrTimoutFromJavaScript)
	}
	if data.ErrCode != 0 {
		if data.ErrCode == JSNOTFOUND {
			return nil, errs.ErrRecordNotFound
		}
		return "", errors.New(data.ErrMsg)
	}
	return data.Data, err
}

func ExtractArrayBuffer(arrayBuffer js.Value) []byte {
	uint8Array := js.Global().Get("Uint8Array").New(arrayBuffer)
	dst := make([]byte, uint8Array.Length())
	js.CopyBytesToGo(dst, uint8Array)
	return dst
}
