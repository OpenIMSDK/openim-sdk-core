//go:build js && wasm

package wasm_wrapper

// ------------------------------------group---------------------------
type WrapperSignaling struct {
	*WrapperCommon
}

func NewWrapperSignaling(wrapperCommon *WrapperCommon) *WrapperSignaling {
	return &WrapperSignaling{WrapperCommon: wrapperCommon}
}
