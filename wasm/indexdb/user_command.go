//go:build js && wasm

package indexdb

import (
	"context"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/utils"
	"github.com/openimsdk/openim-sdk-core/v3/wasm/exec"
)

type LocalUserCommand struct{}

func NewLocalUserCommand() *LocalUserCommand {
	return &LocalUserCommand{}
}

func (i *LocalUserCommand) ProcessUserCommandGetAll(ctx context.Context) ([]*model_struct.LocalUserCommand, error) {
	c, err := exec.Exec()
	if err != nil {
		return nil, err
	} else {
		if v, ok := c.(string); ok {
			result := []*model_struct.LocalUserCommand{}
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

func (i *LocalUserCommand) ProcessUserCommandAdd(ctx context.Context, command *model_struct.LocalUserCommand) error {
	_, err := exec.Exec(utils.StructToJsonString(command))
	return err
}

func (i *LocalUserCommand) ProcessUserCommandUpdate(ctx context.Context, command *model_struct.LocalUserCommand) error {
	_, err := exec.Exec(utils.StructToJsonString(command))
	return err
}
func (i *LocalUserCommand) ProcessUserCommandDelete(ctx context.Context, command *model_struct.LocalUserCommand) error {
	_, err := exec.Exec(utils.StructToJsonString(command))
	return err
}
