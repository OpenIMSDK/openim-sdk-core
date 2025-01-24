//go:build js && wasm

package indexdb

import (
	"context"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/utils"
	"github.com/openimsdk/openim-sdk-core/v3/wasm/exec"
	"github.com/openimsdk/openim-sdk-core/v3/wasm/indexdb/temp_struct"
)

type FriendRequest struct {
	loginUserID string
}

func NewFriendRequest(loginUserID string) *FriendRequest {
	return &FriendRequest{loginUserID: loginUserID}
}

func (i FriendRequest) InsertFriendRequest(ctx context.Context, friendRequest *model_struct.LocalFriendRequest) error {
	_, err := exec.Exec(utils.StructToJsonString(friendRequest))
	return err
}

func (i FriendRequest) DeleteFriendRequestBothUserID(ctx context.Context, fromUserID, toUserID string) error {
	_, err := exec.Exec(fromUserID, toUserID)
	return err
}

func (i FriendRequest) UpdateFriendRequest(ctx context.Context, friendRequest *model_struct.LocalFriendRequest) error {
	tempLocalFriendRequest := temp_struct.LocalFriendRequest{
		FromUserID:    friendRequest.FromUserID,
		FromNickname:  friendRequest.FromNickname,
		FromFaceURL:   friendRequest.FromFaceURL,
		ToUserID:      friendRequest.ToUserID,
		ToNickname:    friendRequest.ToNickname,
		ToFaceURL:     friendRequest.ToFaceURL,
		HandleResult:  friendRequest.HandleResult,
		ReqMsg:        friendRequest.ReqMsg,
		CreateTime:    friendRequest.CreateTime,
		HandlerUserID: friendRequest.HandlerUserID,
		HandleMsg:     friendRequest.HandleMsg,
		HandleTime:    friendRequest.HandleTime,
		Ex:            friendRequest.Ex,
		AttachedInfo:  friendRequest.AttachedInfo,
	}
	_, err := exec.Exec(utils.StructToJsonString(tempLocalFriendRequest))
	return err
}

func (i FriendRequest) GetRecvFriendApplication(ctx context.Context) (result []*model_struct.LocalFriendRequest, err error) {
	gList, err := exec.Exec(i.loginUserID)
	if err != nil {
		return nil, err
	} else {
		if v, ok := gList.(string); ok {
			var temp []model_struct.LocalFriendRequest
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

func (i FriendRequest) GetSendFriendApplication(ctx context.Context) (result []*model_struct.LocalFriendRequest, err error) {
	gList, err := exec.Exec(i.loginUserID)
	if err != nil {
		return nil, err
	} else {
		if v, ok := gList.(string); ok {
			var temp []model_struct.LocalFriendRequest
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

func (i FriendRequest) GetFriendApplicationByBothID(ctx context.Context, fromUserID, toUserID string) (*model_struct.LocalFriendRequest, error) {
	c, err := exec.Exec(fromUserID, toUserID)
	if err != nil {
		return nil, err
	} else {
		if v, ok := c.(string); ok {
			result := model_struct.LocalFriendRequest{}
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

func (i FriendRequest) GetBothFriendReq(ctx context.Context, fromUserID, toUserID string) (result []*model_struct.LocalFriendRequest, err error) {
	gList, err := exec.Exec(fromUserID, toUserID)
	if err != nil {
		return nil, err
	} else {
		if v, ok := gList.(string); ok {
			var temp []model_struct.LocalFriendRequest
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
