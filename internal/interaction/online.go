package interaction

import (
	"context"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
	userPb "github.com/openimsdk/protocol/user"
)

func (c *LongConnMgr) subscribeUsersStatus(ctx context.Context, userIDs []string) ([]*userPb.OnlineStatus, error) {
	if len(userIDs) == 0 {
		return []*userPb.OnlineStatus{}, nil
	}
	res, err := c.GetUserOnlinePlatformIDs(ctx, userIDs)
	if err != nil {
		return nil, err
	}
	status := make([]*userPb.OnlineStatus, 0, len(res))
	for userID, platformIDs := range res {
		value := &userPb.OnlineStatus{
			UserID:      userID,
			PlatformIDs: platformIDs,
		}
		if len(platformIDs) == 0 {
			value.Status = constant.Offline
		} else {
			value.Status = constant.Online
		}
		status = append(status, value)
	}
	return status, nil
}

func (c *LongConnMgr) UnsubscribeUsersStatus(ctx context.Context, userIDs []string) error {
	return c.UnsubscribeUserOnlinePlatformIDs(ctx, userIDs)
}

func (c *LongConnMgr) SubscribeUsersStatus(ctx context.Context, userIDs []string) ([]*userPb.OnlineStatus, error) {
	if len(userIDs) == 0 {
		return []*userPb.OnlineStatus{}, nil
	}
	return c.subscribeUsersStatus(ctx, userIDs)
}

func (c *LongConnMgr) GetSubscribeUsersStatus(ctx context.Context) ([]*userPb.OnlineStatus, error) {
	return c.subscribeUsersStatus(ctx, nil)
}
