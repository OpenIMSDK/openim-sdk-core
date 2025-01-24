package sdk_params_callback

import (
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	"github.com/openimsdk/protocol/wrapperspb"
)

type ProcessFriendApplicationParams struct {
	ToUserID  string `json:"toUserID" validate:"required"`
	HandleMsg string `json:"handleMsg"`
}

type SearchFriendsParam struct {
	KeywordList      []string `json:"keywordList"`
	IsSearchUserID   bool     `json:"isSearchUserID"`
	IsSearchNickname bool     `json:"isSearchNickname"`
	IsSearchRemark   bool     `json:"isSearchRemark"`
}

type SearchFriendItem struct {
	model_struct.LocalFriend
	Relationship int `json:"relationship"`
}

type SetFriendRemarkParams struct {
	ToUserID string `json:"toUserID" validate:"required"`
	Remark   string `json:"remark" validate:"required"`
}
type SetFriendPinParams struct {
	ToUserIDs []string              `json:"toUserIDs" validate:"required"`
	IsPinned  *wrapperspb.BoolValue `json:"isPinned" validate:"required"`
}
