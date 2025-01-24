//go:build js && wasm

package indexdb

import (
	"context"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/utils"
	"github.com/openimsdk/openim-sdk-core/v3/wasm/exec"
)

type LocalTableMaster struct {
}

func NewLocalTableMaster() *LocalTableMaster {
	return &LocalTableMaster{}
}

func (i *LocalTableMaster) GetExistTables(ctx context.Context) (result []string, err error) {
	nameList, err := exec.Exec()
	if err != nil {
		return nil, err
	} else {
		if v, ok := nameList.(string); ok {
			err := utils.JsonStringToStruct(v, &result)
			if err != nil {
				return nil, err
			}
			return result, err
		} else {
			return nil, exec.ErrType
		}
	}
}
