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

package test

import (
	"context"
	"fmt"
	"github.com/openimsdk/openim-sdk-core/v3/internal/util"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/ccontext"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/common"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/log"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/network"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/server_api_params"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/utils"
	"github.com/openimsdk/openim-sdk-core/v3/sdk_struct"
	"net"
	"os"
	"strconv"
	"sync"
	"time"

	authPB "github.com/OpenIMSDK/protocol/auth"
	"github.com/OpenIMSDK/protocol/sdkws"
	userPB "github.com/OpenIMSDK/protocol/user"
	imLog "github.com/OpenIMSDK/tools/log"
)

func GenUid(uid int, prefix string) string {
	if getMyIP() == "" {
		log.Error("", "getMyIP() failed, exit ")
		os.Exit(1)
	}
	UidPrefix := getMyIP() + "_" + prefix + "_"
	return UidPrefix + strconv.FormatInt(int64(uid), 10)
}

func RegisterOnlineAccounts(number int) {
	var wg sync.WaitGroup
	wg.Add(number)
	for i := 0; i < number; i++ {
		go func(t int) {
			userID := GenUid(t, "online")
			register(userID)
			log.Info("register ", userID)
			wg.Done()
		}(i)

	}
	wg.Wait()
	log.Info("", "RegisterAccounts finish ", number)
}

type GetTokenReq struct {
	Secret   string `json:"secret"`
	Platform int    `json:"platform"`
	Uid      string `json:"uid"`
}

type RegisterReq struct {
	Secret   string `json:"secret"`
	Platform int    `json:"platform"`
	Uid      string `json:"uid"`
	Name     string `json:"name"`
}

type ResToken struct {
	Data struct {
		ExpiredTime int64  `json:"expiredTime"`
		Token       string `json:"token"`
		Uid         string `json:"uid"`
	}
	ErrCode int    `json:"errCode"`
	ErrMsg  string `json:"errMsg"`
}

var AdminToken = ""

func init() {
	AdminToken = getToken("openIM123456")
	if err := imLog.InitFromConfig("open-im-sdk-core", "", int(LogLevel), IsLogStandardOutput, false, LogFilePath, 0, 24); err != nil {
		fmt.Println("123456", "log init failed ", err.Error())
	}
}

var ctx context.Context
var cctx context.Context

func register(uid string) error {
	if ctx.Err() == context.Canceled {
		return errors.New("register aborted due to context cancellation")
	}
	...
}
	ctx = ccontext.WithInfo(context.Background(), &ccontext.GlobalConfig{
		UserID: uid,
		Token:  AdminToken,
		IMConfig: sdk_struct.IMConfig{
			PlatformID: PlatformID,
			ApiAddr:    APIADDR,
			WsAddr:     WSADDR,
			LogLevel:   LogLevel,
		},
	})
	ctx = ccontext.WithOperationID(ctx, "123456")

	//ACCOUNTCHECK
	var getAccountCheckReq userPB.AccountCheckReq
	var getAccountCheckResp userPB.AccountCheckResp
	getAccountCheckReq.CheckUserIDs = []string{uid}

	for {
		if ctx.Err() == context.Canceled {
			return errors.New("ApiPost aborted due to context cancellation")
		}
		err := util.ApiPost(ctx, "/user/account_check", &getAccountCheckReq, &getAccountCheckResp)
		if err != nil {
		...
	}
			return err
		}
		if len(getAccountCheckResp.Results) == 1 &&
			getAccountCheckResp.Results[0].AccountStatus == "registered" {
			log.Warn(getAccountCheckReq.CheckUserIDs[0], "Already registered ", uid, getAccountCheckResp)
			userLock.Lock()
			allUserID = append(allUserID, uid)
			userLock.Unlock()
			return nil
		} else if len(getAccountCheckResp.Results) == 1 &&
			getAccountCheckResp.Results[0].AccountStatus == "unregistered" {
			log.Info(getAccountCheckReq.CheckUserIDs[0], "not registered ", uid, getAccountCheckResp)
			break
		} else {
			log.Error(getAccountCheckReq.CheckUserIDs[0], " failed, continue ", err, REGISTERADDR, getAccountCheckReq)
			continue
		}
	}

	var rreq userPB.UserRegisterReq
	rreq.Users = []*sdkws.UserInfo{{UserID: uid}}

	for {
		if ctx.Err() == context.Canceled {
			return errors.New("ApiPost aborted due to context cancellation")
		}
		err := util.ApiPost(ctx, "/auth/user_register", &rreq, nil)
		if err != nil {
		...
	}
			log.Error("post failed ,continue ", err.Error(), REGISTERADDR, getAccountCheckReq)
			time.Sleep(100 * time.Millisecond)
			continue
		} else {
			log.Info("register ok ", REGISTERADDR, getAccountCheckReq)
			userLock.Lock()
			allUserID = append(allUserID, uid)
			userLock.Unlock()
			return nil
		}
	}
}

func getToken(uid string) string {
	if ctx.Err() == context.Canceled {
		return errors.New("getToken aborted due to context cancellation")
	}
	...
}

func RunGetToken(strMyUid string) string {
	if ctx.Err() == context.Canceled {
		return errors.New("RunGetToken aborted due to context cancellation")
	}
	var token string
	for true {
		token = getToken(strMyUid)
		...
	}
}
		if token == "" {
			time.Sleep(time.Duration(100) * time.Millisecond)
			continue
		} else {
			break
		}
	}
	return token
}

func getMyIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Error("", "InterfaceAddrs failed ", err.Error())
		os.Exit(1)
		return ""
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

func RegisterReliabilityUser(id int, timeStamp string) {
	if ctx.Err() == context.Canceled {
		return errors.New("RegisterReliabilityUser aborted due to context cancellation")
	}
	...
}

func WorkGroupRegisterReliabilityUser(id int) {
	if ctx.Err() == context.Canceled {
		return errors.New("WorkGroupRegisterReliabilityUser aborted due to context cancellation")
	}
	...
}

func RegisterPressUser(id int) {
	if ctx.Err() == context.Canceled {
		return errors.New("RegisterPressUser aborted due to context cancellation")
	}
	...
}

func GetGroupMemberNum(groupID string) uint32 {
	if ctx.Err() == context.Canceled {
		return errors.New("GetGroupMemberNum aborted due to context cancellation")
	}
	...
}
