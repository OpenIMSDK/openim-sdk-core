package group

import (
	"context"

	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/utils/datautil"
)

func (g *Group) SyncAllJoinedGroupsAndMembers(ctx context.Context) error {
	if err := g.IncrSyncJoinGroup(ctx); err != nil {
		return err
	}
	return g.IncrSyncJoinGroupMember(ctx)
}

func (g *Group) SyncAllSelfGroupApplication(ctx context.Context) error {
	list, err := g.GetServerSelfGroupApplication(ctx)
	if err != nil {
		return err
	}
	localData, err := g.db.GetSendGroupApplication(ctx)
	if err != nil {
		return err
	}
	if err := g.groupRequestSyncer.Sync(ctx, datautil.Batch(ServerGroupRequestToLocalGroupRequest, list), localData, nil); err != nil {
		return err
	}
	return nil
}

func (g *Group) SyncAllSelfGroupApplicationWithoutNotice(ctx context.Context) error {
	list, err := g.GetServerSelfGroupApplication(ctx)
	if err != nil {
		return err
	}
	localData, err := g.db.GetSendGroupApplication(ctx)
	if err != nil {
		return err
	}
	if err := g.groupRequestSyncer.Sync(ctx, datautil.Batch(ServerGroupRequestToLocalGroupRequest, list), localData, nil, false, true); err != nil {
		return err
	}
	return nil
}

func (g *Group) SyncSelfGroupApplications(ctx context.Context, groupIDs ...string) error {
	return g.SyncAllSelfGroupApplication(ctx)
}

func (g *Group) SyncAllAdminGroupApplication(ctx context.Context) error {
	requests, err := g.GetServerAdminGroupApplicationList(ctx)
	if err != nil {
		return err
	}
	localData, err := g.db.GetAdminGroupApplication(ctx)
	if err != nil {
		return err
	}
	return g.groupAdminRequestSyncer.Sync(ctx, datautil.Batch(ServerGroupRequestToLocalAdminGroupRequest, requests), localData, nil)
}

func (g *Group) SyncAllAdminGroupApplicationWithoutNotice(ctx context.Context) error {
	requests, err := g.GetServerAdminGroupApplicationList(ctx)
	if err != nil {
		return err
	}
	localData, err := g.db.GetAdminGroupApplication(ctx)
	if err != nil {
		return err
	}
	return g.groupAdminRequestSyncer.Sync(ctx, datautil.Batch(ServerGroupRequestToLocalAdminGroupRequest, requests), localData, nil, false, true)
}

func (g *Group) SyncAdminGroupApplications(ctx context.Context, groupIDs ...string) error {
	return g.SyncAllAdminGroupApplication(ctx)
}

func (g *Group) GetServerJoinGroup(ctx context.Context) ([]*sdkws.GroupInfo, error) {
	return g.getServerJoinGroup(ctx)
}

func (g *Group) GetServerAdminGroupApplicationList(ctx context.Context) ([]*sdkws.GroupRequest, error) {
	return g.getServerAdminGroupApplicationList(ctx)
}

func (g *Group) GetServerSelfGroupApplication(ctx context.Context) ([]*sdkws.GroupRequest, error) {
	return g.getServerSelfGroupApplication(ctx)
}

func (g *Group) GetDesignatedGroupMembers(ctx context.Context, groupID string, userIDs []string) ([]*sdkws.GroupMemberFullInfo, error) {
	return g.getDesignatedGroupMembers(ctx, groupID, userIDs)
}
