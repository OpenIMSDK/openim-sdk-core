//go:build js && wasm

package indexdb

import (
	"context"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/utils"
	"github.com/openimsdk/openim-sdk-core/v3/wasm/exec"
	"github.com/openimsdk/openim-sdk-core/v3/wasm/indexdb/temp_struct"
)

type Black struct {
	loginUserID string
}

func NewBlack(loginUserID string) *Black {
	return &Black{loginUserID: loginUserID}
}

// GetBlackListDB gets the blacklist list from the database
func (i Black) GetBlackListDB(ctx context.Context) (result []*model_struct.LocalBlack, err error) {
	gList, err := exec.Exec()
	if err != nil {
		return nil, err
	} else {
		if v, ok := gList.(string); ok {
			var temp []model_struct.LocalBlack
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

// GetBlackListUserID gets the list of blocked user IDs
func (i Black) GetBlackListUserID(ctx context.Context) (result []string, err error) {
	gList, err := exec.Exec()
	if err != nil {
		return nil, err
	} else {
		if v, ok := gList.(string); ok {
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

// GetBlackInfoByBlockUserID gets the information of a blocked user by their user ID
func (i Black) GetBlackInfoByBlockUserID(ctx context.Context, blockUserID string) (result *model_struct.LocalBlack, err error) {
	gList, err := exec.Exec(blockUserID, i.loginUserID)
	if err != nil {
		return nil, err
	} else {
		if v, ok := gList.(string); ok {
			var temp model_struct.LocalBlack
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

// GetBlackInfoList gets the information of multiple blocked users by their user IDs
func (i Black) GetBlackInfoList(ctx context.Context, blockUserIDList []string) (result []*model_struct.LocalBlack, err error) {
	gList, err := exec.Exec(utils.StructToJsonString(blockUserIDList))
	if err != nil {
		return nil, err
	} else {
		if v, ok := gList.(string); ok {
			var temp []model_struct.LocalBlack
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

// InsertBlack inserts a new blocked user into the database
func (i Black) InsertBlack(ctx context.Context, black *model_struct.LocalBlack) error {
	_, err := exec.Exec(utils.StructToJsonString(black))
	return err
}

// UpdateBlack updates the information of a blocked user in the database
func (i Black) UpdateBlack(ctx context.Context, black *model_struct.LocalBlack) error {
	tempLocalBlack := temp_struct.LocalBlack{
		Nickname:       black.Nickname,
		FaceURL:        black.FaceURL,
		CreateTime:     black.CreateTime,
		AddSource:      black.AddSource,
		OperatorUserID: black.OperatorUserID,
		Ex:             black.Ex,
		AttachedInfo:   black.AttachedInfo,
	}
	_, err := exec.Exec(black.OwnerUserID, black.BlockUserID, utils.StructToJsonString(tempLocalBlack))
	return err
}

// DeleteBlack removes a blocked user from the database
func (i Black) DeleteBlack(ctx context.Context, blockUserID string) error {
	_, err := exec.Exec(blockUserID, i.loginUserID)
	return err
}
