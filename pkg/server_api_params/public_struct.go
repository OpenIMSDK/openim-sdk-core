// Copyright © 2023 OpenIM SDK. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package server_api_params

import (
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
)

type PublicUser struct {
	UserID     string `json:"userID"`
	Nickname   string `json:"nickname"`
	FaceURL    string `json:"faceURL"`
	Ex         string `json:"ex"`
	CreateTime int64  `json:"createTime"`
}

type FullUserInfo struct {
	PublicInfo *PublicUser               `json:"publicInfo"`
	FriendInfo *model_struct.LocalFriend `json:"friendInfo"`
	BlackInfo  *model_struct.LocalBlack  `json:"blackInfo"`
}

type FullUserInfoWithCache struct {
	PublicInfo      *PublicUser                    `json:"publicInfo"`
	FriendInfo      *model_struct.LocalFriend      `json:"friendInfo"`
	BlackInfo       *model_struct.LocalBlack       `json:"blackInfo"`
	GroupMemberInfo *model_struct.LocalGroupMember `json:"groupMemberInfo"`
}
