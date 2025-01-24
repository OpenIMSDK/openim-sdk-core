//go:build js && wasm

package event_listener

import (
	"bytes"
	"context"
	"errors"
	"reflect"
	"strconv"
	"strings"
	"syscall/js"

	"github.com/openimsdk/tools/errs"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/utils"
	"github.com/openimsdk/openim-sdk-core/v3/wasm/exec"
	"github.com/openimsdk/tools/log"
)

type Caller interface {
	// AsyncCallWithCallback has promise object return
	AsyncCallWithCallback() interface{}
	// AsyncCallWithOutCallback has promise object return
	AsyncCallWithOutCallback() interface{}
	// SyncCall has not promise
	SyncCall() (result []interface{})
}

type FuncLogic func()

var ErrNotSetCallback = errors.New("not set callback to call")
var ErrNotSetFunc = errors.New("not set func to call")

type ReflectCall struct {
	funcName  interface{}
	callback  CallbackWriter
	arguments []js.Value
}

func NewCaller(funcName interface{}, callback CallbackWriter, arguments *[]js.Value) Caller {
	return &ReflectCall{funcName: funcName, callback: callback, arguments: *arguments}
}

func (r *ReflectCall) AsyncCallWithCallback() interface{} {
	return r.callback.HandlerFunc(r.asyncCallWithCallback)

}
func (r *ReflectCall) asyncCallWithCallback() {
	defer func() {
		if rc := recover(); rc != nil {
			r.ErrHandle(rc)
		}
	}()
	ctx := context.Background()
	var funcName reflect.Value
	var typeFuncName reflect.Type
	var hasCallback bool
	var temp int
	if r.funcName == nil {
		panic(ErrNotSetFunc)
	} else {
		funcName = reflect.ValueOf(r.funcName)
		typeFuncName = reflect.TypeOf(r.funcName)
	}
	var values []reflect.Value
	if r.callback != nil {
		hasCallback = true
		r.callback.SetOperationID(r.arguments[0].String())
		values = append(values, reflect.ValueOf(r.callback))
	} else {
		log.ZDebug(ctx, "AsyncCallWithCallback not set callback")
		panic(ErrNotSetCallback)
	}
	funcFieldsNum := typeFuncName.NumIn()
	if funcFieldsNum-len(r.arguments) > 1 {
		r.arguments = append(r.arguments, js.Value{})
	}
	for i := 0; i < len(r.arguments); i++ {
		if hasCallback {
			temp++
		} else {
			temp = i
		}
		//log.NewDebug(r.callback.GetOperationID(), "type is ", typeFuncName.In(temp).Kind(), r.arguments[i].IsNaN())
		switch typeFuncName.In(temp).Kind() {
		case reflect.String:
			convertValue := r.arguments[i].String()
			if !strings.HasPrefix(convertValue, "<number: ") {
				values = append(values, reflect.ValueOf(convertValue))
			} else {
				log.ZError(ctx, "AsyncCallWithCallback", nil, "input args type err index:",
					utils.IntToString(i))
				panic("input args type err index:" + utils.IntToString(i))
			}
		case reflect.Int:
			values = append(values, reflect.ValueOf(r.arguments[i].Int()))
		case reflect.Int32:
			values = append(values, reflect.ValueOf(int32(r.arguments[i].Int())))
		case reflect.Bool:
			values = append(values, reflect.ValueOf(r.arguments[i].Bool()))
		case reflect.Int64:
			values = append(values, reflect.ValueOf(int64(r.arguments[i].Int())))
		case reflect.Ptr:
			values = append(values, reflect.ValueOf(bytes.NewBuffer(exec.ExtractArrayBuffer(r.arguments[i]))))
		default:
			log.ZError(ctx, "AsyncCallWithCallback", nil,
				"input args type not support:", strconv.Itoa(int(typeFuncName.In(temp).Kind())))
			panic("input args type not support:" + strconv.Itoa(int(typeFuncName.In(temp).Kind())))
		}
	}
	funcName.Call(values)

}
func (r *ReflectCall) AsyncCallWithOutCallback() interface{} {
	if r.callback == nil {
		r.callback = NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), nil)
	}
	return r.callback.HandlerFunc(r.asyncCallWithOutCallback)
}
func (r *ReflectCall) asyncCallWithOutCallback() {
	defer func() {
		if rc := recover(); rc != nil {
			r.ErrHandle(rc)
		}
	}()
	ctx := context.Background()
	var funcName reflect.Value
	var typeFuncName reflect.Type
	if r.funcName == nil {
		panic(ErrNotSetFunc)
	} else {
		funcName = reflect.ValueOf(r.funcName)
		typeFuncName = reflect.TypeOf(r.funcName)
	}
	var values []reflect.Value
	if r.callback == nil {
		r.callback = NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), nil)
	}
	log.ZWarn(ctx, "asyncCall", nil, "asyncCallWithOutCallback", len(r.arguments))

	r.callback.SetOperationID(r.arguments[0].String())
	//strings.SplitAfter()
	for i := 0; i < len(r.arguments); i++ {
		//log.NewDebug(r.callback.GetOperationID(), "type is ", typeFuncName.In(temp).Kind(), r.arguments[i].IsNaN())
		switch typeFuncName.In(i).Kind() {
		case reflect.String:
			convertValue := r.arguments[i].String()
			if !strings.HasPrefix(convertValue, "<number: ") {
				values = append(values, reflect.ValueOf(convertValue))
			} else {
				panic("input args type err index:" + utils.IntToString(i))
			}
		case reflect.Int:
			values = append(values, reflect.ValueOf(r.arguments[i].Int()))
		case reflect.Int32:
			values = append(values, reflect.ValueOf(int32(r.arguments[i].Int())))
		case reflect.Bool:
			values = append(values, reflect.ValueOf(r.arguments[i].Bool()))
		case reflect.Float64:
			values = append(values, reflect.ValueOf(r.arguments[i].Float()))
		default:
			panic("input args type not support:" + strconv.Itoa(int(typeFuncName.In(i).Kind())))
		}
	}
	go func() {

		returnValues := funcName.Call(values)
		if len(returnValues) != 0 {
			var result []interface{}
			for _, v := range returnValues {
				switch v.Kind() {
				case reflect.String:
					result = append(result, v.String())
				case reflect.Bool:
					result = append(result, v.Bool())
				case reflect.Int32:
					result = append(result, v.Int())
				case reflect.Int:
					result = append(result, v.Int())
				default:
					panic("not support type")
				}
			}
			r.callback.SetData(result).SendMessage()
		} else {
			r.callback.SetErrCode(200).SetErrMsg(errors.New("null string").Error()).SendMessage()
		}
	}()

}
func (r *ReflectCall) SyncCall() (result []interface{}) {
	defer func() {
		if rc := recover(); rc != nil {
			temp := r.ErrHandle(rc)
			for _, v := range temp {
				result = append(result, v)
			}
		}
	}()
	var funcName reflect.Value
	var typeFuncName reflect.Type
	var hasCallback bool
	var temp int
	if r.funcName == nil {
		return nil
	} else {
		funcName = reflect.ValueOf(r.funcName)
		typeFuncName = reflect.TypeOf(r.funcName)
	}
	var values []reflect.Value
	if r.callback != nil {
		hasCallback = true
		r.callback.SetOperationID(r.arguments[0].String())
		values = append(values, reflect.ValueOf(r.callback))
	}
	for i := 0; i < len(r.arguments); i++ {
		if hasCallback {
			temp++
		} else {
			temp = i
		}
		//log.NewDebug(r.callback.GetOperationID(), "type is ", typeFuncName.In(temp).Kind(), r.arguments[i].IsNaN())
		switch typeFuncName.In(temp).Kind() {
		case reflect.String:
			convertValue := r.arguments[i].String()
			if !strings.HasPrefix(convertValue, "<number: ") {
				values = append(values, reflect.ValueOf(convertValue))
			} else {
				panic("input args type err index:" + utils.IntToString(i))
			}
		case reflect.Int:
			values = append(values, reflect.ValueOf(r.arguments[i].Int()))
		case reflect.Int32:
			values = append(values, reflect.ValueOf(int32(r.arguments[i].Int())))
		case reflect.Bool:
			values = append(values, reflect.ValueOf(r.arguments[i].Bool()))
		default:
			panic("implement me")
		}
	}
	returnValues := funcName.Call(values)
	if len(returnValues) != 0 {
		for _, v := range returnValues {
			switch v.Kind() {
			case reflect.String:
				result = append(result, v.String())
			case reflect.Bool:
				result = append(result, v.Bool())
			default:
				panic("not support type")
			}
		}
		return result

	} else {
		return nil
	}

}
func (r *ReflectCall) ErrHandle(recover interface{}) []string {
	ctx := context.Background()
	var temp string
	switch x := recover.(type) {
	case string:
		log.ZError(ctx, "STRINGERR", nil, "r", x)
		temp = errs.WrapMsg(errors.New(x), "").Error()
	case error:
		//buf := make([]byte, 1<<20)
		//runtime.Stack(buf, true)
		log.ZError(ctx, "ERR", x, "r", x.Error())
		temp = x.Error()
	default:
		log.ZError(ctx, "unknown panic", nil, "r", x)
		temp = errs.WrapMsg(errors.New("unknown panic"), "").Error()
	}
	if r.callback != nil {
		r.callback.SetErrCode(100).SetErrMsg(temp).SendMessage()
		return []string{}
	} else {
		return []string{temp}
	}
}
