package manager

import (
	"context"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/config"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/pkg/decorator"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/pkg/reerrgroup"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/pkg/utils"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/sdk"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/vars"
)

type TestGroupManager struct {
	*MetaManager
}

func NewGroupManager(m *MetaManager) *TestGroupManager {
	return &TestGroupManager{m}
}

// CreateGroups creates group chats. It needs to create both large group chats and regular group chats.
// The number of large group chats to be created is specified by vars.LargeGroupNum, and the group owner cycles from 0 to vars.UserNum.
// Every user creates regular group chats, and the number of regular group chats to be created is specified by vars.CommonGroupNum.
func (m *TestGroupManager) CreateGroups(ctx context.Context) error {
	defer decorator.FuncLog(ctx)()

	gr := reerrgroup.NewGroup(config.ErrGroupCommonLimit)
	m.createLargeGroups(ctx, gr)
	m.createCommonGroups(ctx, gr)
	return gr.Wait()
}

// createLargeGroups see CreateGroups
func (m *TestGroupManager) createLargeGroups(ctx context.Context, gr *reerrgroup.Group) {
	userNum := 0
	for i := 0; i < vars.LargeGroupNum; i++ {
		ctx := vars.Contexts[userNum]
		testSDK := sdk.TestSDKs[userNum]
		gr.Go(func() error {
			_, err := testSDK.CreateLargeGroup(ctx)
			if err != nil {
				return err
			}
			return nil
		})
		userNum = utils.NextNum(userNum)
	}
	return
}

// createLargeGroups see CreateGroups
func (m *TestGroupManager) createCommonGroups(ctx context.Context, gr *reerrgroup.Group) {
	for userNum := 0; userNum < vars.UserNum; userNum++ {
		ctx := vars.Contexts[userNum]
		testSDK := sdk.TestSDKs[userNum]
		gr.Go(func() error {
			for i := 0; i < vars.CommonGroupNum; i++ {
				_, err := testSDK.CreateCommonGroup(ctx, vars.CommonGroupMemberNum)
				if err != nil {
					return err
				}
			}
			return nil
		})
	}
	return
}
