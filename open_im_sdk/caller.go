package open_im_sdk

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/pkg/errors"

	"github.com/openimsdk/openim-sdk-core/v3/open_im_sdk_callback"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/ccontext"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/sdkerrs"

	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/mw"

	"github.com/openimsdk/tools/errs"
)

func isNumeric(kind reflect.Kind) bool {
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return true
	default:
		return false
	}
}

func setNumeric(in interface{}, out interface{}) {
	inValue := reflect.ValueOf(in)
	outValue := reflect.ValueOf(out)
	outElem := outValue.Elem()
	outType := outElem.Type()
	inType := inValue.Type()
	if outType.AssignableTo(inType) {
		outElem.Set(inValue)
		return
	}
	inKind := inValue.Kind()
	outKind := outElem.Kind()
	switch inKind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch outKind {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			outElem.SetInt(inValue.Int())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			outElem.SetUint(uint64(inValue.Int()))
		case reflect.Float32, reflect.Float64:
			outElem.SetFloat(float64(inValue.Int()))
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		switch outKind {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			outElem.SetInt(int64(inValue.Uint()))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			outElem.SetUint(inValue.Uint())
		case reflect.Float32, reflect.Float64:
			outElem.SetFloat(float64(inValue.Uint()))
		}
	case reflect.Float32, reflect.Float64:
		switch outKind {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			outElem.SetInt(int64(inValue.Float()))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			outElem.SetUint(uint64(inValue.Float()))
		case reflect.Float32, reflect.Float64:
			outElem.SetFloat(inValue.Float())
		}
	}
}

func call_(operationID string, fn any, args ...any) (res any, err error) {
	t := time.Now()
	funcPtr := reflect.ValueOf(fn).Pointer()
	funcName := runtime.FuncForPC(funcPtr).Name()
	if operationID == "" {
		return nil, sdkerrs.ErrArgs.WrapMsg("call function operationID is empty")
	}
	if err := CheckResourceLoad(UserForSDK, funcName); err != nil {
		return nil, sdkerrs.ErrResourceLoad.WrapMsg("not load resource")
	}
	ctx := ccontext.WithOperationID(UserForSDK.Context(), operationID)

	defer func(start time.Time) {
		if r := recover(); r != nil {
			p := fmt.Sprintf("panic: %+v\n%s", r, debug.Stack())
			err = fmt.Errorf("call panic: %+v", p)
		} else {
			elapsed := time.Since(start).Milliseconds()
			if err == nil {
				log.ZInfo(ctx, "fn call success", "function name", funcName, "cost time", fmt.Sprintf("%d ms", elapsed), "resp", res)
			} else {
				log.ZError(ctx, "fn call error", mw.FormatError(err), "function name", funcName, "cost time", fmt.Sprintf("%d ms", elapsed))

			}

		}
	}(t)

	log.ZInfo(ctx, "func call req", "function name", funcName, "args", args)
	fnv := reflect.ValueOf(fn)
	if fnv.Kind() != reflect.Func {
		return nil, sdkerrs.ErrSdkInternal.WrapMsg(fmt.Sprintf("call function fn is not function, is %T", fn))
	}

	fnt := fnv.Type()
	nin := fnt.NumIn()

	if len(args)+1 != nin {
		return nil, sdkerrs.ErrSdkInternal.WrapMsg("go code error: fn in args num is not match")
	}

	ins := make([]reflect.Value, 0, nin)
	ins = append(ins, reflect.ValueOf(ctx))
	for i := 0; i < len(args); i++ {
		inFnField := fnt.In(i + 1)
		arg := reflect.TypeOf(args[i])
		if arg.String() == inFnField.String() || inFnField.Kind() == reflect.Interface {
			ins = append(ins, reflect.ValueOf(args[i]))
			continue
		}
		if arg.Kind() == reflect.String { // json
			var ptr int
			for inFnField.Kind() == reflect.Ptr {
				inFnField = inFnField.Elem()
				ptr++
			}
			switch inFnField.Kind() {
			case reflect.Struct, reflect.Slice, reflect.Array, reflect.Map:
				v := reflect.New(inFnField)
				if err := json.Unmarshal([]byte(args[i].(string)), v.Interface()); err != nil {
					return nil, sdkerrs.ErrSdkInternal.WrapMsg(fmt.Sprintf("go call json.Unmarshal error: %s", err))
				}
				if ptr == 0 {
					v = v.Elem()
				} else if ptr != 1 {
					for i := ptr - 1; i > 0; i-- {
						temp := reflect.New(v.Type())
						temp.Elem().Set(v)
						v = temp
					}
				}
				ins = append(ins, v)
				continue
			}
		}
		return nil, sdkerrs.ErrSdkInternal.WrapMsg("go code error: fn in args type is not match")
	}

	outs := fnv.Call(ins)
	if len(outs) == 0 {
		return "", nil
	}

	if fnt.Out(len(outs) - 1).Implements(reflect.ValueOf(new(error)).Elem().Type()) {
		if errValueOf := outs[len(outs)-1]; !errValueOf.IsNil() {
			return nil, errValueOf.Interface().(error)
		}
		if len(outs) == 1 {
			return "", nil
		}
		outs = outs[:len(outs)-1]
	}

	for i := 0; i < len(outs); i++ {
		out := outs[i]
		switch out.Kind() {
		case reflect.Map:
			if out.IsNil() {
				outs[i] = reflect.MakeMap(out.Type())
			}
		case reflect.Slice:
			if out.IsNil() {
				outs[i] = reflect.MakeSlice(out.Type(), 0, 0)
			}
		}
	}

	if len(outs) == 1 {
		return outs[0].Interface(), nil
	}

	val := make([]any, 0, len(outs))
	for i := range outs {
		val = append(val, outs[i].Interface())
	}

	return val, nil
}

func call(callback open_im_sdk_callback.Base, operationID string, fn any, args ...any) {
	if callback == nil {
		log.ZWarn(context.Background(), "callback is nil", nil)
		return
	}
	go func() {
		res, err := call_(operationID, fn, args...)
		if err != nil {
			if code, ok := err.(errs.CodeError); ok {
				callback.OnError(int32(code.Code()), code.Error())
			} else {
				callback.OnError(sdkerrs.UnknownCode, fmt.Sprintf("error %T not implement CodeError: %s", err, err))
			}
			return
		}
		data, err := json.Marshal(res)
		if err != nil {
			callback.OnError(sdkerrs.SdkInternalError, fmt.Sprintf("function res json.Marshal error: %s", err))
			return
		}
		callback.OnSuccess(string(data))
	}()
}

func syncCall(operationID string, fn any, args ...any) (res string) {
	err := error(nil)
	if operationID == "" {
		return ""
	}
	fnv := reflect.ValueOf(fn)
	if fnv.Kind() != reflect.Func {
		err = errs.ErrRecordNotFound
		return ""
	}
	funcPtr := reflect.ValueOf(fn).Pointer()
	funcName := runtime.FuncForPC(funcPtr).Name()
	if err = CheckResourceLoad(UserForSDK, funcName); err != nil {
		return ""
	}
	fnt := fnv.Type()
	numIn := fnt.NumIn()
	if len(args)+1 != numIn {
		err = errors.New("go code error: fn in args num is not match")
		return ""
	}
	ins := make([]reflect.Value, 0, numIn)

	ctx := ccontext.WithOperationID(UserForSDK.Context(), operationID)
	t := time.Now()
	defer func(start time.Time) {
		if r := recover(); r != nil {
			fmt.Printf("panic: %+v\n%s", r, debug.Stack())
		} else {
			elapsed := time.Since(start).Milliseconds()
			if err == nil {
				log.ZInfo(ctx, "fn call success", "function name", funcName, "resp", res, "cost time", fmt.Sprintf("%d ms", elapsed))
			} else {
				log.ZError(ctx, "fn call error", err, "function name", funcName, "cost time", fmt.Sprintf("%d ms", elapsed))
			}

		}
	}(t)
	log.ZInfo(ctx, "func call req", "function name", funcName, "args", args)
	ins = append(ins, reflect.ValueOf(ctx))
	for i := 0; i < len(args); i++ {
		tag := fnt.In(i + 1)
		arg := reflect.TypeOf(args[i])
		if arg.String() == tag.String() || tag.Kind() == reflect.Interface {
			ins = append(ins, reflect.ValueOf(args[i]))
			continue
		}
		if arg.Kind() == reflect.String { // json
			switch tag.Kind() {
			case reflect.Struct, reflect.Slice, reflect.Array, reflect.Map, reflect.Ptr:
				v := reflect.New(tag)
				if args[i].(string) != "" {
					if err = json.Unmarshal([]byte(args[i].(string)), v.Interface()); err != nil {
						return ""
					}
				}

				ins = append(ins, v.Elem())
				continue
			}
		}
		if isNumeric(arg.Kind()) && isNumeric(tag.Kind()) {
			v := reflect.Zero(tag).Interface()
			setNumeric(args[i], &v)
			ins = append(ins, reflect.ValueOf(v))
			continue
		}
		err = errors.New("err args type")
		return ""
	}
	var lastErr bool
	if numOut := fnt.NumOut(); numOut > 0 {
		lastErr = fnt.Out(numOut - 1).Implements(reflect.TypeOf((*error)(nil)).Elem())
	}
	outs := fnv.Call(ins)
	if len(outs) == 0 {
		err = errors.New("err res type")
		return ""
	}
	outVals := make([]any, 0, len(outs))
	for i := 0; i < len(outs); i++ {
		outVals = append(outVals, outs[i].Interface())
	}
	if lastErr {
		if last := outVals[len(outVals)-1]; last != nil {
			err = last.(error)
			return ""
		}
		if len(outs) == 1 {
			return ""
		}
		outVals = outVals[:len(outVals)-1]
	}
	// Convert nil maps and slices to non-nil
	for i := 0; i < len(outVals); i++ {
		switch outs[i].Kind() {
		case reflect.Map:
			if outs[i].IsNil() {
				outVals[i] = reflect.MakeMap(outs[i].Type()).Interface()
			}
		case reflect.Slice:
			if outs[i].IsNil() {
				outVals[i] = reflect.MakeSlice(outs[i].Type(), 0, 0).Interface()
			}
		}
	}
	var jsonVal any
	if len(outVals) == 1 {
		jsonVal = outVals[0]
	} else {
		jsonVal = outVals
	}
	jsonData, err := json.Marshal(jsonVal)
	if err != nil {
		err = errors.New("json marshal error")
		return ""
	}
	return string(jsonData)
}
func messageCall(callback open_im_sdk_callback.SendMsgCallBack, operationID string, fn any, args ...any) {
	if callback == nil {
		log.ZWarn(context.Background(), "callback is nil", nil)
		return
	}
	go messageCall_(callback, operationID, fn, args...)
}
func messageCall_(callback open_im_sdk_callback.SendMsgCallBack, operationID string, fn any, args ...any) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(" panic err:", r, string(debug.Stack()))
			callback.OnError(sdkerrs.SdkInternalError, fmt.Sprintf("recover: %+v", r))
			return
		}
	}()
	if operationID == "" {
		callback.OnError(sdkerrs.ArgsError, sdkerrs.ErrArgs.WrapMsg("operationID is empty").Error())
		return
	}
	if err := CheckResourceLoad(UserForSDK, ""); err != nil {
		callback.OnError(sdkerrs.ResourceLoadNotCompleteError, "resource load error: "+err.Error())
		return
	}
	fnv := reflect.ValueOf(fn)
	if fnv.Kind() != reflect.Func {
		callback.OnError(sdkerrs.SdkInternalError, "go code error: fn is not function")
		return
	}
	fnt := fnv.Type()
	numIn := fnt.NumIn()
	if len(args)+1 != numIn {
		callback.OnError(sdkerrs.SdkInternalError, "go code error: fn in args num is not match")
		return
	}

	t := time.Now()
	ins := make([]reflect.Value, 0, numIn)
	ctx := ccontext.WithOperationID(UserForSDK.Context(), operationID)
	ctx = ccontext.WithSendMessageCallback(ctx, callback)
	funcPtr := reflect.ValueOf(fn).Pointer()
	funcName := runtime.FuncForPC(funcPtr).Name()
	log.ZInfo(ctx, "input req", "function name", funcName, "args", args)

	ins = append(ins, reflect.ValueOf(ctx))
	for i := 0; i < len(args); i++ { // callback open_im_sdk_callback.Base, operationID string, ...
		tag := fnt.In(i + 1) // ctx context.Context, ...
		arg := reflect.TypeOf(args[i])
		if arg.String() == tag.String() || tag.Kind() == reflect.Interface {
			ins = append(ins, reflect.ValueOf(args[i]))
			continue
		}
		if arg.Kind() == reflect.String { // json
			switch tag.Kind() {
			case reflect.Struct, reflect.Slice, reflect.Array, reflect.Map, reflect.Ptr:
				v := reflect.New(tag)
				if err := json.Unmarshal([]byte(args[i].(string)), v.Interface()); err != nil {
					callback.OnError(sdkerrs.ArgsError, err.Error())
					return
				}
				ins = append(ins, v.Elem())
				continue
			}
		}
		if isNumeric(arg.Kind()) && isNumeric(tag.Kind()) {
			v := reflect.Zero(tag).Interface()
			setNumeric(args[i], &v)
			ins = append(ins, reflect.ValueOf(v))
			continue
		}
		callback.OnError(sdkerrs.ArgsError, "go code error: fn in args type is not match")
		return
	}
	var lastErr bool
	if numOut := fnt.NumOut(); numOut > 0 {
		lastErr = fnt.Out(numOut - 1).Implements(reflect.ValueOf(new(error)).Elem().Type())
	}
	outs := fnv.Call(ins)

	outVals := make([]any, 0, len(outs))
	for i := 0; i < len(outs); i++ {
		outVals = append(outVals, outs[i].Interface())
	}
	if lastErr {
		if last := outVals[len(outVals)-1]; last != nil {
			if code, ok := last.(error).(errs.CodeError); ok {
				callback.OnError(int32(code.Code()), code.Error())
			} else {
				callback.OnError(sdkerrs.UnknownCode, fmt.Sprintf("error %T not implement CodeError: %s", last.(error), last.(error).Error()))
			}
			return
		}

		outVals = outVals[:len(outVals)-1]
	}
	// Convert nil maps and slices to non-nil
	for i := 0; i < len(outVals); i++ {
		switch outs[i].Kind() {
		case reflect.Map:
			if outs[i].IsNil() {
				outVals[i] = reflect.MakeMap(outs[i].Type()).Interface()
			}
		case reflect.Slice:
			if outs[i].IsNil() {
				outVals[i] = reflect.MakeSlice(outs[i].Type(), 0, 0).Interface()
			}
		}
	}
	var jsonVal any
	if len(outVals) == 1 {
		jsonVal = outVals[0]
	} else {
		jsonVal = outVals
	}
	jsonData, err := json.Marshal(jsonVal)
	if err != nil {
		callback.OnError(sdkerrs.ArgsError, err.Error())
		return
	}
	log.ZInfo(ctx, "output resp", "function name", funcName, "resp", jsonVal, "cost time", time.Since(t))
	callback.OnSuccess(string(jsonData))
}

func listenerCall(fn any, listener any) {
	ctx := context.Background()
	if UserForSDK == nil {
		log.ZWarn(ctx, "UserForSDK is nil,set listener is invalid", nil)
		return
	}
	fnv := reflect.ValueOf(fn)
	if fnv.Kind() != reflect.Func {
		log.ZWarn(ctx, "fn is error,set listener is invalid", nil)
		return
	}
	args := reflect.ValueOf(listener)
	fnv.Call([]reflect.Value{args})
}
