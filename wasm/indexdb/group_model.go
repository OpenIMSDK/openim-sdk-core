//go:build js && wasm

package indexdb

import (
	"context"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/utils"
	"github.com/openimsdk/openim-sdk-core/v3/wasm/exec"
)

type LocalGroups struct{}

func NewLocalGroups() *LocalGroups {
	return &LocalGroups{}
}

func (i *LocalGroups) InsertGroup(ctx context.Context, groupInfo *model_struct.LocalGroup) error {
	_, err := exec.Exec(utils.StructToJsonString(groupInfo))
	return err
}

func (i *LocalGroups) DeleteGroup(ctx context.Context, groupID string) error {
	_, err := exec.Exec(groupID)
	return err
}

// this method need full update
func (i *LocalGroups) UpdateGroup(ctx context.Context, groupInfo *model_struct.LocalGroup) error {
	_, err := exec.Exec(groupInfo.GroupID, utils.StructToJsonString(groupInfo))
	return err
}

func (i *LocalGroups) BatchInsertGroup(ctx context.Context, groupList []*model_struct.LocalGroup) error {
	_, err := exec.Exec(utils.StructToJsonString(groupList))
	return err
}
func (i *LocalGroups) DeleteAllGroup(ctx context.Context) error {
	_, err := exec.Exec()
	return err
}

func (i *LocalGroups) GetJoinedGroupListDB(ctx context.Context) (result []*model_struct.LocalGroup, err error) {
	gList, err := exec.Exec()
	if err != nil {
		return nil, err
	} else {
		if v, ok := gList.(string); ok {
			var temp []model_struct.LocalGroup
			err := utils.JsonStringToStruct(v, &temp)
			if err != nil {
				return nil, err
			}
			for _, v := range temp {
				v1 := v
				result = append(result, &v1)
			}
			return result, err
		} else {
			return nil, exec.ErrType
		}
	}
}

func (i *LocalGroups) GetGroups(ctx context.Context, groupIDs []string) (result []*model_struct.LocalGroup, err error) {
	gList, err := exec.Exec(utils.StructToJsonString(groupIDs))
	if err != nil {
		return nil, err
	} else {
		if v, ok := gList.(string); ok {
			var temp []model_struct.LocalGroup
			err := utils.JsonStringToStruct(v, &temp)
			if err != nil {
				return nil, err
			}
			for _, v := range temp {
				v1 := v
				result = append(result, &v1)
			}
			return result, err
		} else {
			return nil, exec.ErrType
		}
	}
}

func (i *LocalGroups) GetGroupInfoByGroupID(ctx context.Context, groupID string) (*model_struct.LocalGroup, error) {
	c, err := exec.Exec(groupID)
	if err != nil {
		return nil, err
	} else {
		if v, ok := c.(string); ok {
			result := model_struct.LocalGroup{}
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

func (i *LocalGroups) GetAllGroupInfoByGroupIDOrGroupName(ctx context.Context, keyword string, isSearchGroupID bool, isSearchGroupName bool) (result []*model_struct.LocalGroup, err error) {
	gList, err := exec.Exec(keyword, isSearchGroupID, isSearchGroupName)
	if err != nil {
		return nil, err
	} else {
		if v, ok := gList.(string); ok {
			var temp []model_struct.LocalGroup
			err := utils.JsonStringToStruct(v, &temp)
			if err != nil {
				return nil, err
			}
			for _, v := range temp {
				v1 := v
				result = append(result, &v1)
			}
			return result, err
		} else {
			return nil, exec.ErrType
		}
	}
}
