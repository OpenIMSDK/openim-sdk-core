//go:build js && wasm

package indexdb

import (
	"context"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/utils"
	"github.com/openimsdk/openim-sdk-core/v3/wasm/exec"
)

type LocalGroupRequest struct {
}

func NewLocalGroupRequest() *LocalGroupRequest {
	return &LocalGroupRequest{}
}

func (i *LocalGroupRequest) InsertGroupRequest(ctx context.Context, groupRequest *model_struct.LocalGroupRequest) error {
	_, err := exec.Exec(utils.StructToJsonString(groupRequest))
	return err
}

func (i *LocalGroupRequest) DeleteGroupRequest(ctx context.Context, groupID, userID string) error {
	_, err := exec.Exec(groupID, userID)
	return err
}

func (i *LocalGroupRequest) UpdateGroupRequest(ctx context.Context, groupRequest *model_struct.LocalGroupRequest) error {
	_, err := exec.Exec(utils.StructToJsonString(groupRequest))
	return err
}

func (i *LocalGroupRequest) GetSendGroupApplication(ctx context.Context) ([]*model_struct.LocalGroupRequest, error) {
	result, err := exec.Exec()
	if err != nil {
		return nil, err
	}
	if v, ok := result.(string); ok {
		var request []*model_struct.LocalGroupRequest
		if err := utils.JsonStringToStruct(v, &request); err != nil {
			return nil, err
		}
		return request, nil
	} else {
		return nil, exec.ErrType
	}
}

func (i *LocalGroupRequest) GetAdminGroupApplication(ctx context.Context) ([]*model_struct.LocalAdminGroupRequest, error) {
	result, err := exec.Exec()
	if err != nil {
		return nil, err
	}
	if v, ok := result.(string); ok {
		var request []*model_struct.LocalAdminGroupRequest
		if err := utils.JsonStringToStruct(v, &request); err != nil {
			return nil, err
		}
		return request, nil
	} else {
		return nil, exec.ErrType
	}
}

func (i *LocalGroupRequest) InsertAdminGroupRequest(ctx context.Context, groupRequest *model_struct.LocalAdminGroupRequest) error {
	_, err := exec.Exec(utils.StructToJsonString(groupRequest))
	return err
}

func (i *LocalGroupRequest) DeleteAdminGroupRequest(ctx context.Context, groupID, userID string) error {
	_, err := exec.Exec(groupID, userID)
	return err
}

func (i *LocalGroupRequest) UpdateAdminGroupRequest(ctx context.Context, groupRequest *model_struct.LocalAdminGroupRequest) error {
	_, err := exec.Exec(utils.StructToJsonString(groupRequest))
	return err
}
