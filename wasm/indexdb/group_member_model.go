//go:build js && wasm

package indexdb

import (
	"context"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/utils"
	"github.com/openimsdk/openim-sdk-core/v3/wasm/exec"
)

type LocalGroupMember struct {
}

func NewLocalGroupMember() *LocalGroupMember {
	return &LocalGroupMember{}
}

func (i *LocalGroupMember) GetGroupMemberInfoByGroupIDUserID(ctx context.Context, groupID, userID string) (*model_struct.LocalGroupMember, error) {
	member, err := exec.Exec(groupID, userID)
	if err != nil {
		return nil, err
	} else {
		if v, ok := member.(string); ok {
			var temp model_struct.LocalGroupMember
			err := utils.JsonStringToStruct(v, &temp)
			if err != nil {
				return nil, err
			}
			return &temp, err
		} else {
			return nil, exec.ErrType
		}
	}
}

func (i *LocalGroupMember) GetAllGroupMemberUserIDList(ctx context.Context) ([]model_struct.LocalGroupMember, error) {
	member, err := exec.Exec()
	if err != nil {
		return nil, err
	} else {
		if v, ok := member.(string); ok {
			var temp []model_struct.LocalGroupMember
			err := utils.JsonStringToStruct(v, &temp)
			if err != nil {
				return nil, err
			}
			return temp, err
		} else {
			return nil, exec.ErrType
		}
	}
}

func (i *LocalGroupMember) GetGroupMemberCount(ctx context.Context, groupID string) (int32, error) {
	count, err := exec.Exec(groupID)
	if err != nil {
		return 0, err
	}
	if v, ok := count.(float64); ok {
		return int32(v), nil
	}
	return 0, exec.ErrType
}

func (i *LocalGroupMember) GetGroupSomeMemberInfo(ctx context.Context, groupID string, userIDList []string) ([]*model_struct.LocalGroupMember, error) {
	member, err := exec.Exec(groupID, utils.StructToJsonString(userIDList))
	if err != nil {
		return nil, err
	} else {
		if v, ok := member.(string); ok {
			var temp []*model_struct.LocalGroupMember
			err := utils.JsonStringToStruct(v, &temp)
			if err != nil {
				return nil, err
			}
			return temp, err
		} else {
			return nil, exec.ErrType
		}
	}
}

func (i *LocalGroupMember) GetGroupMemberListByGroupID(ctx context.Context, groupID string) ([]*model_struct.LocalGroupMember, error) {
	member, err := exec.Exec(groupID)
	if err != nil {
		return nil, err
	} else {
		if v, ok := member.(string); ok {
			var temp []*model_struct.LocalGroupMember
			err := utils.JsonStringToStruct(v, &temp)
			if err != nil {
				return nil, err
			}
			return temp, err
		} else {
			return nil, exec.ErrType
		}
	}
}

func (i *LocalGroupMember) GetGroupMemberListByUserIDs(ctx context.Context, groupID string, filter int32, userIDs []string) ([]*model_struct.LocalGroupMember, error) {
	member, err := exec.Exec(groupID, filter, utils.StructToJsonString(userIDs))
	if err != nil {
		return nil, err
	} else {
		if v, ok := member.(string); ok {
			var temp []*model_struct.LocalGroupMember
			if err := utils.JsonStringToStruct(v, &temp); err != nil {
				return nil, err
			}
			return temp, err
		} else {
			return nil, exec.ErrType
		}
	}
}

func (i *LocalGroupMember) GetGroupMemberListSplit(ctx context.Context, groupID string, filter int32, offset, count int) ([]*model_struct.LocalGroupMember, error) {
	member, err := exec.Exec(groupID, filter, offset, count)
	if err != nil {
		return nil, err
	} else {
		if v, ok := member.(string); ok {
			var temp []*model_struct.LocalGroupMember
			err := utils.JsonStringToStruct(v, &temp)
			if err != nil {
				return nil, err
			}
			return temp, err
		} else {
			return nil, exec.ErrType
		}
	}
}

func (i *LocalGroupMember) GetGroupMemberOwnerAndAdminDB(ctx context.Context, groupID string) ([]*model_struct.LocalGroupMember, error) {
	member, err := exec.Exec(groupID)
	if err != nil {
		return nil, err
	} else {
		if v, ok := member.(string); ok {
			var temp []*model_struct.LocalGroupMember
			err := utils.JsonStringToStruct(v, &temp)
			if err != nil {
				return nil, err
			}
			return temp, err
		} else {
			return nil, exec.ErrType
		}
	}
}

func (i *LocalGroupMember) GetGroupMemberListSplitByJoinTimeFilter(ctx context.Context, groupID string, offset, count int, joinTimeBegin, joinTimeEnd int64, userIDList []string) ([]*model_struct.LocalGroupMember, error) {
	member, err := exec.Exec(groupID, offset, count, joinTimeBegin, joinTimeEnd, utils.StructToJsonString(userIDList))
	if err != nil {
		return nil, err
	} else {
		if v, ok := member.(string); ok {
			var temp []*model_struct.LocalGroupMember
			err := utils.JsonStringToStruct(v, &temp)
			if err != nil {
				return nil, err
			}
			return temp, err
		} else {
			return nil, exec.ErrType
		}
	}
}

func (i *LocalGroupMember) InsertGroupMember(ctx context.Context, groupMember *model_struct.LocalGroupMember) error {
	_, err := exec.Exec(utils.StructToJsonString(groupMember))
	return err
}

func (i *LocalGroupMember) BatchInsertGroupMember(ctx context.Context, groupMemberList []*model_struct.LocalGroupMember) error {
	_, err := exec.Exec(utils.StructToJsonString(groupMemberList))
	return err
}

func (i *LocalGroupMember) DeleteGroupMember(ctx context.Context, groupID, userID string) error {
	_, err := exec.Exec(groupID, userID)
	return err
}

func (i *LocalGroupMember) DeleteGroupAllMembers(ctx context.Context, groupID string) error {
	_, err := exec.Exec(groupID)
	return err
}

func (i *LocalGroupMember) UpdateGroupMember(ctx context.Context, groupMember *model_struct.LocalGroupMember) error {
	_, err := exec.Exec(utils.StructToJsonString(groupMember))
	return err
}

func (i *LocalGroupMember) SearchGroupMembersDB(ctx context.Context, keyword string, groupID string, isSearchMemberNickname, isSearchUserID bool, offset, count int) (result []*model_struct.LocalGroupMember, err error) {
	member, err := exec.Exec(keyword, groupID, isSearchMemberNickname, isSearchUserID, offset, count)
	if err != nil {
		return nil, err
	} else {
		if v, ok := member.(string); ok {
			var temp []*model_struct.LocalGroupMember
			err := utils.JsonStringToStruct(v, &temp)
			if err != nil {
				return nil, err
			}
			return temp, err
		} else {
			return nil, exec.ErrType
		}
	}
}

func (i *LocalGroupMember) GetUserJoinedGroupIDs(ctx context.Context, userID string) (result []string, err error) {
	IDList, err := exec.Exec(userID)
	if err != nil {
		return nil, err
	}
	if v, ok := IDList.(string); ok {
		err := utils.JsonStringToStruct(v, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	}
	return nil, exec.ErrType
}
