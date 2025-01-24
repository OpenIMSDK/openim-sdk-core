//go:build js && wasm

package indexdb

import (
	"context"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/utils"
	"github.com/openimsdk/openim-sdk-core/v3/wasm/exec"
	"github.com/openimsdk/tools/errs"
)

type LocalVersionSync struct {
}

func NewLocalVersionSync() *LocalVersionSync {
	return &LocalVersionSync{}
}

func (i *LocalVersionSync) GetVersionSync(ctx context.Context, tableName, entityID string) (*model_struct.LocalVersionSync, error) {
	version, err := exec.Exec(tableName, entityID)
	if err != nil {
		if err == errs.ErrRecordNotFound {
			return &model_struct.LocalVersionSync{}, err
		}
		return nil, err
	}
	if v, ok := version.(string); ok {
		var temp model_struct.LocalVersionSync
		if err := utils.JsonStringToStruct(v, &temp); err != nil {
			return nil, err
		}
		return &temp, err
	} else {
		return nil, exec.ErrType
	}
}

func (i *LocalVersionSync) SetVersionSync(ctx context.Context, lv *model_struct.LocalVersionSync) error {
	_, err := exec.Exec(utils.StructToJsonString(lv))
	return err
}

func (i *LocalVersionSync) DeleteVersionSync(ctx context.Context, tableName, entityID string) error {
	_, err := exec.Exec(tableName, entityID)
	return err
}
