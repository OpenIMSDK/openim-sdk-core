package group

import (
	"context"
	"sync"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/syncer"
	constantpb "github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/protocol/group"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/utils/datautil"
)

func (g *Group) groupAndMemberVersionTableName() string {
	return "local_group_entities_version"
}

func (g *Group) groupTableName() string {
	return model_struct.LocalGroup{}.TableName()
}

func (g *Group) IncrSyncJoinGroupMember(ctx context.Context) error {
	groups, err := g.db.GetJoinedGroupListDB(ctx)
	if err != nil {
		return err
	}
	groupIDs := datautil.Slice(groups, func(e *model_struct.LocalGroup) string {
		return e.GroupID
	})
	return g.IncrSyncGroupAndMember(ctx, groupIDs...)
}

func (g *Group) IncrSyncGroupAndMember(ctx context.Context, groupIDs ...string) error {
	var wg sync.WaitGroup
	if len(groupIDs) == 0 {
		return nil
	}
	const maxSyncNum = constantpb.MaxSyncPullNumber
	groupIDSet := datautil.SliceSet(groupIDs)
	var groups []*group.GetIncrementalGroupMemberReq
	if len(groupIDs) > maxSyncNum {
		groups = make([]*group.GetIncrementalGroupMemberReq, 0, maxSyncNum)
	} else {
		groups = make([]*group.GetIncrementalGroupMemberReq, 0, len(groupIDs))
	}
	for {
		if len(groupIDSet) == 0 {
			return nil
		}
		for groupID := range groupIDSet {
			if len(groups) == cap(groups) {
				break
			}
			req := group.GetIncrementalGroupMemberReq{
				GroupID: groupID,
			}

			lvs, err := g.db.GetVersionSync(ctx, g.groupAndMemberVersionTableName(), groupID)
			if err == nil {
				req.VersionID = lvs.VersionID
				req.Version = lvs.Version
			} else if !errs.ErrRecordNotFound.Is(err) {
				return err
			}
			groups = append(groups, &req)
		}
		groupVersion, err := g.getIncrementalGroupMemberBatch(ctx, groups)
		if err != nil {
			return err
		}
		groups = groups[:0]
		for groupID, resp := range groupVersion {
			tempResp := resp
			tempGroupID := groupID
			wg.Add(1)
			go func() error {
				defer wg.Done()
				if err := g.syncGroupAndMember(ctx, tempGroupID, tempResp); err != nil {
					log.ZError(ctx, "sync Group And Member error", errs.Wrap(err))
					return errs.Wrap(err)
				}
				return nil
			}()
			delete(groupIDSet, tempGroupID)
		}
		wg.Wait()
	}
}

func (g *Group) syncGroupAndMember(ctx context.Context, groupID string, resp *group.GetIncrementalGroupMemberResp) error {
	groupMemberSyncer := syncer.VersionSynchronizer[*model_struct.LocalGroupMember, *group.GetIncrementalGroupMemberResp]{
		Ctx:       ctx,
		DB:        g.db,
		TableName: g.groupAndMemberVersionTableName(),
		EntityID:  groupID,
		Key: func(localGroupMember *model_struct.LocalGroupMember) string {
			return localGroupMember.UserID
		},
		Local: func() ([]*model_struct.LocalGroupMember, error) {
			return g.db.GetGroupMemberListByGroupID(ctx, groupID)
		},
		ServerVersion: func() *group.GetIncrementalGroupMemberResp {
			return resp
		},
		Full: func(resp *group.GetIncrementalGroupMemberResp) bool {
			return resp.Full
		},
		Version: func(resp *group.GetIncrementalGroupMemberResp) (string, uint64) {
			return resp.VersionID, resp.Version
		},
		Delete: func(resp *group.GetIncrementalGroupMemberResp) []string {
			return resp.Delete
		},
		Update: func(resp *group.GetIncrementalGroupMemberResp) []*model_struct.LocalGroupMember {
			return datautil.Batch(ServerGroupMemberToLocalGroupMember, resp.Update)
		},
		Insert: func(resp *group.GetIncrementalGroupMemberResp) []*model_struct.LocalGroupMember {
			return datautil.Batch(ServerGroupMemberToLocalGroupMember, resp.Insert)
		},
		ExtraData: func(resp *group.GetIncrementalGroupMemberResp) any {
			return resp.Group
		},
		ExtraDataProcessor: func(ctx context.Context, data any) error {
			groupInfo, ok := data.(*sdkws.GroupInfo)
			if !ok {
				return errs.New("group info type error")
			}
			if groupInfo == nil {
				return nil
			}
			local, err := g.db.GetJoinedGroupListDB(ctx)
			if err != nil {
				return err
			}
			log.ZDebug(ctx, "group info", "groupInfo", groupInfo)
			changes := datautil.Batch(ServerGroupToLocalGroup, []*sdkws.GroupInfo{groupInfo})
			kv := datautil.SliceToMapAny(local, func(e *model_struct.LocalGroup) (string, *model_struct.LocalGroup) {
				return e.GroupID, e
			})
			for i, change := range changes {
				key := change.GroupID
				kv[key] = changes[i]
			}
			server := datautil.Values(kv)
			return g.groupSyncer.Sync(ctx, server, local, nil)
		},
		Syncer: func(server, local []*model_struct.LocalGroupMember) error {
			return g.groupMemberSyncer.Sync(ctx, server, local, nil)
		},
		FullSyncer: func(ctx context.Context) error {
			return g.groupMemberSyncer.FullSync(ctx, groupID)
		},
		FullID: func(ctx context.Context) ([]string, error) {
			resp, err := g.getFullGroupMemberUserIDs(ctx, &group.GetFullGroupMemberUserIDsReq{
				GroupID: groupID,
			})
			if err != nil {
				return nil, err
			}
			return resp.UserIDs, nil
		},
		IDOrderChanged: func(resp *group.GetIncrementalGroupMemberResp) bool {
			if resp.SortVersion > 0 {
				return true
			}
			return false
		},
	}
	return groupMemberSyncer.IncrementalSync()
}

func (g *Group) onlineSyncGroupAndMember(ctx context.Context, groupID string, deleteGroupMembers, updateGroupMembers, insertGroupMembers []*sdkws.GroupMemberFullInfo,
	updateGroup *sdkws.GroupInfo, sortVersion uint64, version uint64, versionID string) error {
	groupMemberSyncer := syncer.VersionSynchronizer[*model_struct.LocalGroupMember, *group.GetIncrementalGroupMemberResp]{
		Ctx:       ctx,
		DB:        g.db,
		TableName: g.groupAndMemberVersionTableName(),
		EntityID:  groupID,
		Key: func(localGroupMember *model_struct.LocalGroupMember) string {
			return localGroupMember.UserID
		},
		Local: func() ([]*model_struct.LocalGroupMember, error) {
			return g.db.GetGroupMemberListByGroupID(ctx, groupID)
		},
		ServerVersion: func() *group.GetIncrementalGroupMemberResp {
			return &group.GetIncrementalGroupMemberResp{
				Version:   version,
				VersionID: versionID,
				Full:      false,
				Delete: datautil.Slice(deleteGroupMembers, func(e *sdkws.GroupMemberFullInfo) string {
					return e.UserID
				}),
				Insert:      insertGroupMembers,
				Update:      updateGroupMembers,
				Group:       updateGroup,
				SortVersion: sortVersion,
			}
		},
		Server: func(version *model_struct.LocalVersionSync) (*group.GetIncrementalGroupMemberResp, error) {
			singleGroupReq := &group.GetIncrementalGroupMemberReq{
				GroupID:   groupID,
				VersionID: version.VersionID,
				Version:   version.Version,
			}

			resp, err := g.getIncrementalGroupMemberBatch(ctx, []*group.GetIncrementalGroupMemberReq{singleGroupReq})
			if err != nil {
				return nil, err
			}
			if resp != nil {
				if singleGroupResp, ok := resp[groupID]; ok {
					return singleGroupResp, nil
				}
			}
			return nil, errs.New("group member version record not found")

		},
		Full: func(resp *group.GetIncrementalGroupMemberResp) bool {
			return resp.Full
		},
		Version: func(resp *group.GetIncrementalGroupMemberResp) (string, uint64) {
			return resp.VersionID, resp.Version
		},
		Delete: func(resp *group.GetIncrementalGroupMemberResp) []string {
			return resp.Delete
		},
		Update: func(resp *group.GetIncrementalGroupMemberResp) []*model_struct.LocalGroupMember {
			return datautil.Batch(ServerGroupMemberToLocalGroupMember, resp.Update)
		},
		Insert: func(resp *group.GetIncrementalGroupMemberResp) []*model_struct.LocalGroupMember {
			return datautil.Batch(ServerGroupMemberToLocalGroupMember, resp.Insert)
		},
		ExtraData: func(resp *group.GetIncrementalGroupMemberResp) any {
			return resp.Group
		},
		ExtraDataProcessor: func(ctx context.Context, data any) error {
			groupInfo, ok := data.(*sdkws.GroupInfo)
			if !ok {
				return errs.New("group info type error")
			}
			if groupInfo == nil {
				return nil
			}
			local, err := g.db.GetJoinedGroupListDB(ctx)
			if err != nil {
				return err
			}
			log.ZDebug(ctx, "group info", "groupInfo", groupInfo)
			changes := datautil.Batch(ServerGroupToLocalGroup, []*sdkws.GroupInfo{groupInfo})
			kv := datautil.SliceToMapAny(local, func(e *model_struct.LocalGroup) (string, *model_struct.LocalGroup) {
				return e.GroupID, e
			})
			for i, change := range changes {
				key := change.GroupID
				kv[key] = changes[i]
			}
			server := datautil.Values(kv)
			return g.groupSyncer.Sync(ctx, server, local, nil)
		},
		Syncer: func(server, local []*model_struct.LocalGroupMember) error {
			return g.groupMemberSyncer.Sync(ctx, server, local, nil)
		},
		FullSyncer: func(ctx context.Context) error {
			return g.groupMemberSyncer.FullSync(ctx, groupID)
		},
		FullID: func(ctx context.Context) ([]string, error) {
			resp, err := g.getFullGroupMemberUserIDs(ctx, &group.GetFullGroupMemberUserIDsReq{
				GroupID: groupID,
			})
			if err != nil {
				return nil, err
			}
			return resp.UserIDs, nil
		},
		IDOrderChanged: func(resp *group.GetIncrementalGroupMemberResp) bool {
			if resp.SortVersion > 0 {
				return true
			}
			return false
		},
	}
	return groupMemberSyncer.CheckVersionSync()
}

func (g *Group) IncrSyncJoinGroup(ctx context.Context) error {
	joinedGroupSyncer := syncer.VersionSynchronizer[*model_struct.LocalGroup, *group.GetIncrementalJoinGroupResp]{
		Ctx:       ctx,
		DB:        g.db,
		TableName: g.groupTableName(),
		EntityID:  g.loginUserID,
		Key: func(LocalGroup *model_struct.LocalGroup) string {
			return LocalGroup.GroupID
		},
		Local: func() ([]*model_struct.LocalGroup, error) {
			return g.db.GetJoinedGroupListDB(ctx)
		},
		Server: func(version *model_struct.LocalVersionSync) (*group.GetIncrementalJoinGroupResp, error) {
			return g.getIncrementalJoinGroup(ctx, &group.GetIncrementalJoinGroupReq{
				UserID:    g.loginUserID,
				Version:   version.Version,
				VersionID: version.VersionID,
			})
		},
		Full: func(resp *group.GetIncrementalJoinGroupResp) bool {
			return resp.Full
		},
		Version: func(resp *group.GetIncrementalJoinGroupResp) (string, uint64) {
			return resp.VersionID, resp.Version
		},
		Delete: func(resp *group.GetIncrementalJoinGroupResp) []string {
			return resp.Delete
		},
		Update: func(resp *group.GetIncrementalJoinGroupResp) []*model_struct.LocalGroup {
			return datautil.Batch(ServerGroupToLocalGroup, resp.Update)
		},
		Insert: func(resp *group.GetIncrementalJoinGroupResp) []*model_struct.LocalGroup {
			return datautil.Batch(ServerGroupToLocalGroup, resp.Insert)
		},
		Syncer: func(server, local []*model_struct.LocalGroup) error {
			return g.groupSyncer.Sync(ctx, server, local, nil)
		},
		FullSyncer: func(ctx context.Context) error {
			return g.groupSyncer.FullSync(ctx, g.loginUserID)
		},
		FullID: func(ctx context.Context) ([]string, error) {
			resp, err := g.getFullJoinGroupIDs(ctx, &group.GetFullJoinGroupIDsReq{
				UserID: g.loginUserID,
			})
			if err != nil {
				return nil, err
			}
			return resp.GroupIDs, nil

		},
		IDOrderChanged: func(resp *group.GetIncrementalJoinGroupResp) bool {
			if resp.SortVersion > 0 {
				return true
			}
			return false
		},
	}
	return joinedGroupSyncer.IncrementalSync()
}
