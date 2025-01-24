package sdkerrs

import "github.com/openimsdk/tools/errs"

func New(code int, msg string, dtl string) errs.CodeError {
	return errs.NewCodeError(code, msg).WithDetail(dtl)
}

var (
	Wrap    = errs.Wrap
	WrapMsg = errs.WrapMsg
)
