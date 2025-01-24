//go:build js && wasm

package indexdb

import (
	"context"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/utils"
	"github.com/openimsdk/openim-sdk-core/v3/wasm/exec"
)

type LocalAppSDKVersion struct {
}

func NewLocalAppSDKVersion() *LocalAppSDKVersion {
	return &LocalAppSDKVersion{}
}

func (i *LocalAppSDKVersion) GetAppSDKVersion(ctx context.Context) (*model_struct.LocalAppSDKVersion, error) {
	sdkVersion, err := exec.Exec()
	if err != nil {
		return nil, err
	} else {
		if v, ok := sdkVersion.(string); ok {
			var temp model_struct.LocalAppSDKVersion
			if err := utils.JsonStringToStruct(v, &temp); err != nil {
				return nil, err
			}
			return &temp, err
		} else {
			return nil, exec.ErrType
		}
	}
}

func (i *LocalAppSDKVersion) SetAppSDKVersion(ctx context.Context, appVersion *model_struct.LocalAppSDKVersion) error {
	_, err := exec.Exec(utils.StructToJsonString(appVersion))
	return err
}
