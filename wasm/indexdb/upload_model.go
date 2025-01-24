//go:build js && wasm

package indexdb

import (
	"context"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/utils"
	"github.com/openimsdk/openim-sdk-core/v3/wasm/exec"
)

type LocalUpload struct{}

func NewLocalUpload() *LocalUpload {
	return &LocalUpload{}
}

func (i *LocalUpload) GetUpload(ctx context.Context, partHash string) (*model_struct.LocalUpload, error) {
	c, err := exec.Exec(partHash)
	if err != nil {
		return nil, err
	} else {
		if v, ok := c.(string); ok {
			result := model_struct.LocalUpload{}
			err := utils.JsonStringToStruct(v, &result)
			if err != nil {
				return nil, err
			}
			return &result, err
		} else {
			return nil, exec.ErrType
		}
	}
}

func (i *LocalUpload) InsertUpload(ctx context.Context, upload *model_struct.LocalUpload) error {
	_, err := exec.Exec(utils.StructToJsonString(upload))
	return err
}

func (i *LocalUpload) DeleteUpload(ctx context.Context, partHash string) error {
	_, err := exec.Exec(partHash)
	return err
}
func (i *LocalUpload) UpdateUpload(ctx context.Context, upload *model_struct.LocalUpload) error {
	_, err := exec.Exec(utils.StructToJsonString(upload))
	return err
}

func (i *LocalUpload) DeleteExpireUpload(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}
