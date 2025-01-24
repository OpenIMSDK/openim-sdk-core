package group

import (
	"context"
	"time"

	"github.com/openimsdk/tools/errs"

	"github.com/openimsdk/tools/utils/datautil"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/datafetcher"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/sdk_params_callback"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/sdkerrs"

	"github.com/openimsdk/protocol/group"
	"github.com/openimsdk/protocol/sdkws"
)

func (g *Group) CreateGroup(ctx context.Context, req *group.CreateGroupReq) (*sdkws.GroupInfo, error) {
	if req.OwnerUserID == "" {
		req.OwnerUserID = g.loginUserID
	}
	if req.GroupInfo.GroupType != constant.WorkingGroup {
		return nil, sdkerrs.ErrGroupType
	}
	req.GroupInfo.CreatorUserID = g.loginUserID
	resp, err := g.createGroup(ctx, req)
	if err != nil {
		return nil, err
	}

	g.groupSyncMutex.Lock()
	defer g.groupSyncMutex.Unlock()

	if err := g.IncrSyncJoinGroup(ctx); err != nil {
		return nil, err
	}
	if err := g.IncrSyncGroupAndMember(ctx, resp.GroupInfo.GroupID); err != nil {
		return nil, err
	}
	return resp.GroupInfo, nil
}

func (g *Group) JoinGroup(ctx context.Context, groupID, reqMsg string, joinSource int32, ex string) error {
	req := &group.JoinGroupReq{GroupID: groupID, ReqMessage: reqMsg, JoinSource: joinSource, InviterUserID: g.loginUserID, Ex: ex}
	if err := g.joinGroup(ctx, req); err != nil {
		return err
	}
	if err := g.SyncSelfGroupApplications(ctx, groupID); err != nil {
		return err
	}
	return nil
}

func (g *Group) QuitGroup(ctx context.Context, groupID string) error {
	if err := g.quitGroup(ctx, groupID); err != nil {
		return err
	}
	g.groupSyncMutex.Lock()
	defer g.groupSyncMutex.Unlock()

	if err := g.IncrSyncJoinGroup(ctx); err != nil {
		return err
	}
	return nil
}

func (g *Group) DismissGroup(ctx context.Context, groupID string) error {
	if err := g.dismissGroup(ctx, groupID); err != nil {
		return err
	}
	g.groupSyncMutex.Lock()
	defer g.groupSyncMutex.Unlock()

	if err := g.IncrSyncJoinGroup(ctx); err != nil {
		return err
	}
	return nil
}

func (g *Group) ChangeGroupMute(ctx context.Context, groupID string, isMute bool) (err error) {
	if isMute {
		err = g.muteGroup(ctx, groupID)
	} else {
		err = g.cancelMuteGroup(ctx, groupID)
	}
	if err != nil {
		return err
	}

	g.groupSyncMutex.Lock()
	defer g.groupSyncMutex.Unlock()

	if err := g.IncrSyncGroupAndMember(ctx, groupID); err != nil {
		return err
	}
	return nil
}

func (g *Group) ChangeGroupMemberMute(ctx context.Context, groupID, userID string, mutedSeconds int) error {
	if mutedSeconds == 0 {
		return g.cancelMuteGroupMember(ctx, &group.CancelMuteGroupMemberReq{GroupID: groupID, UserID: userID})
	} else {
		return g.muteGroupMember(ctx, &group.MuteGroupMemberReq{GroupID: groupID, UserID: userID, MutedSeconds: uint32(mutedSeconds)})
	}
}

func (g *Group) TransferGroupOwner(ctx context.Context, groupID, newOwnerUserID string) error {
	req := &group.TransferGroupOwnerReq{GroupID: groupID, OldOwnerUserID: g.loginUserID, NewOwnerUserID: newOwnerUserID}
	if err := g.transferGroup(ctx, req); err != nil {
		return err
	}
	g.groupSyncMutex.Lock()
	defer g.groupSyncMutex.Unlock()

	if err := g.IncrSyncGroupAndMember(ctx, groupID); err != nil {
		return err
	}
	return nil
}

func (g *Group) KickGroupMember(ctx context.Context, groupID string, reason string, userIDList []string) error {
	req := &group.KickGroupMemberReq{GroupID: groupID, KickedUserIDs: userIDList, Reason: reason}
	if err := g.kickGroupMember(ctx, req); err != nil {
		return err
	}

	g.groupSyncMutex.Lock()
	defer g.groupSyncMutex.Unlock()

	return g.IncrSyncGroupAndMember(ctx, groupID)
}

func (g *Group) SetGroupInfo(ctx context.Context, groupInfo *group.SetGroupInfoExReq) error {
	if err := g.setGroupInfo(ctx, groupInfo); err != nil {
		return err
	}

	g.groupSyncMutex.Lock()
	defer g.groupSyncMutex.Unlock()

	return g.IncrSyncJoinGroup(ctx)
}

func (g *Group) SetGroupMemberInfo(ctx context.Context, groupMemberInfo *group.SetGroupMemberInfo) error {
	req := &group.SetGroupMemberInfoReq{Members: []*group.SetGroupMemberInfo{groupMemberInfo}}
	if err := g.setGroupMemberInfo(ctx, req); err != nil {
		return err
	}

	g.groupSyncMutex.Lock()
	defer g.groupSyncMutex.Unlock()

	return g.IncrSyncGroupAndMember(ctx, groupMemberInfo.GroupID)
}

func (g *Group) GetJoinedGroupList(ctx context.Context) ([]*model_struct.LocalGroup, error) {
	return g.db.GetJoinedGroupListDB(ctx)
}

func (g *Group) GetJoinedGroupListPage(ctx context.Context, offset, count int32) ([]*model_struct.LocalGroup, error) {
	dataFetcher := datafetcher.NewDataFetcher(
		g.db,
		g.groupTableName(),
		g.loginUserID,
		func(localGroup *model_struct.LocalGroup) string {
			return localGroup.GroupID
		},
		func(ctx context.Context, values []*model_struct.LocalGroup) error {
			return g.db.BatchInsertGroup(ctx, values)
		},
		func(ctx context.Context, groupIDs []string) ([]*model_struct.LocalGroup, bool, error) {
			localGroups, err := g.db.GetGroups(ctx, groupIDs)
			return localGroups, true, err
		},
		func(ctx context.Context, groupIDs []string) ([]*model_struct.LocalGroup, error) {
			serverGroupInfo, err := g.getGroupsInfoFromServer(ctx, groupIDs)
			if err != nil {
				return nil, err
			}
			return datautil.Batch(ServerGroupToLocalGroup, serverGroupInfo), nil
		},
	)
	return dataFetcher.FetchWithPagination(ctx, int(offset), int(count))
}

func (g *Group) GetSpecifiedGroupsInfo(ctx context.Context, groupIDs []string) ([]*model_struct.LocalGroup, error) {
	dataFetcher := datafetcher.NewDataFetcher(
		g.db,
		g.groupTableName(),
		g.loginUserID,
		func(localGroup *model_struct.LocalGroup) string {
			return localGroup.GroupID
		},
		func(ctx context.Context, values []*model_struct.LocalGroup) error {
			return g.db.BatchInsertGroup(ctx, values)
		},
		func(ctx context.Context, groupIDs []string) ([]*model_struct.LocalGroup, bool, error) {
			localGroups, err := g.db.GetGroups(ctx, groupIDs)
			return localGroups, true, err
		},
		func(ctx context.Context, groupIDs []string) ([]*model_struct.LocalGroup, error) {
			serverGroupInfo, err := g.getGroupsInfoFromServer(ctx, groupIDs)
			if err != nil {
				return nil, err
			}
			return datautil.Batch(ServerGroupToLocalGroup, serverGroupInfo), nil
		},
	)
	return dataFetcher.FetchMissingAndCombineLocal(ctx, groupIDs)
}

func (g *Group) SearchGroups(ctx context.Context, param sdk_params_callback.SearchGroupsParam) ([]*model_struct.LocalGroup, error) {
	if len(param.KeywordList) == 0 || (!param.IsSearchGroupName && !param.IsSearchGroupID) {
		return nil, sdkerrs.ErrArgs.WrapMsg("keyword is null or search field all false")
	}
	groups, err := g.db.GetAllGroupInfoByGroupIDOrGroupName(ctx, param.KeywordList[0], param.IsSearchGroupID, param.IsSearchGroupName) // todo	param.KeywordList[0]
	if err != nil {
		return nil, err
	}
	return groups, nil
}

func (g *Group) GetGroupMemberOwnerAndAdmin(ctx context.Context, groupID string) ([]*model_struct.LocalGroupMember, error) {
	return g.db.GetGroupMemberOwnerAndAdminDB(ctx, groupID)
}

func (g *Group) GetGroupMemberListByJoinTimeFilter(ctx context.Context, groupID string, offset, count int32, joinTimeBegin, joinTimeEnd int64, userIDs []string) ([]*model_struct.LocalGroupMember, error) {
	if joinTimeEnd == 0 {
		joinTimeEnd = time.Now().UnixMilli()
	}

	dataFetcher := datafetcher.NewDataFetcher(
		g.db,
		g.groupAndMemberVersionTableName(),
		groupID,
		func(localGroupMember *model_struct.LocalGroupMember) string {
			return localGroupMember.UserID
		},
		func(ctx context.Context, values []*model_struct.LocalGroupMember) error {
			return g.db.BatchInsertGroupMember(ctx, values)
		},
		func(ctx context.Context, userIDs []string) ([]*model_struct.LocalGroupMember, bool, error) {
			localGroupMembers, err := g.db.GetGroupMemberListSplitByJoinTimeFilter(ctx, groupID, int(offset), int(count), joinTimeBegin, joinTimeEnd, userIDs)
			return localGroupMembers, true, err
		},
		func(ctx context.Context, userIDs []string) ([]*model_struct.LocalGroupMember, error) {
			serverGroupMember, err := g.GetDesignatedGroupMembers(ctx, groupID, userIDs)
			if err != nil {
				return nil, err
			}
			return datautil.Batch(ServerGroupMemberToLocalGroupMember, serverGroupMember), nil
		},
	)

	return dataFetcher.FetchWithPagination(ctx, int(offset), int(count))
}

func (g *Group) GetSpecifiedGroupMembersInfo(ctx context.Context, groupID string, userIDList []string) ([]*model_struct.LocalGroupMember, error) {
	g.groupSyncMutex.Lock()
	defer g.groupSyncMutex.Unlock()

	lvs, err := g.db.GetVersionSync(ctx, g.groupTableName(), g.loginUserID)
	if err != nil {
		return nil, err
	}
	if datautil.Contain(groupID, lvs.UIDList...) {

		_, err := g.db.GetVersionSync(ctx, g.groupAndMemberVersionTableName(), groupID)
		if err != nil {
			if !errs.ErrRecordNotFound.Is(err) {
				return nil, err
			}
			err := g.IncrSyncGroupAndMember(ctx, groupID)
			if err != nil {
				return nil, err
			}
		}
	} else { // If the user is no longer in the group, return nil immediately
		return nil, nil
	}
	dataFetcher := datafetcher.NewDataFetcher(
		g.db,
		g.groupAndMemberVersionTableName(),
		groupID,
		func(localGroupMember *model_struct.LocalGroupMember) string {
			return localGroupMember.UserID
		},
		func(ctx context.Context, values []*model_struct.LocalGroupMember) error {
			return g.db.BatchInsertGroupMember(ctx, values)
		},
		func(ctx context.Context, userIDs []string) ([]*model_struct.LocalGroupMember, bool, error) {
			localGroupMembers, err := g.db.GetGroupSomeMemberInfo(ctx, groupID, userIDList)
			if err != nil {
				return nil, false, err
			}
			localGroup, err := g.db.GetGroupInfoByGroupID(ctx, groupID)
			if err != nil {
				return nil, false, err
			}
			if localGroup.MemberCount < groupMemberSyncLimit {
				return localGroupMembers, false, nil
			}
			return localGroupMembers, true, nil
		},
		func(ctx context.Context, userIDs []string) ([]*model_struct.LocalGroupMember, error) {
			serverGroupMember, err := g.GetDesignatedGroupMembers(ctx, groupID, userIDs)
			if err != nil {
				return nil, err
			}
			if len(serverGroupMember) == 0 {
				return nil, nil
			}
			return datautil.Batch(ServerGroupMemberToLocalGroupMember, serverGroupMember), nil
		},
	)
	return dataFetcher.FetchMissingAndFillLocal(ctx, userIDList)
}

func (g *Group) GetGroupMemberList(ctx context.Context, groupID string, filter, offset, count int32) ([]*model_struct.LocalGroupMember, error) {
	g.groupSyncMutex.Lock()
	defer g.groupSyncMutex.Unlock()

	lvs, err := g.db.GetVersionSync(ctx, g.groupTableName(), g.loginUserID)
	if err != nil {
		return nil, err
	}
	if datautil.Contain(groupID, lvs.UIDList...) {

		_, err := g.db.GetVersionSync(ctx, g.groupAndMemberVersionTableName(), groupID)
		if err != nil {
			if !errs.ErrRecordNotFound.Is(err) {
				return nil, err
			}
			err := g.IncrSyncGroupAndMember(ctx, groupID)
			if err != nil {
				return nil, err
			}
		}
	} else { // If the user is no longer in the group, return nil immediately
		return nil, nil
	}

	dataFetcher := datafetcher.NewDataFetcher(
		g.db,
		g.groupAndMemberVersionTableName(),
		groupID,
		func(localGroupMember *model_struct.LocalGroupMember) string {
			return localGroupMember.UserID
		},
		func(ctx context.Context, values []*model_struct.LocalGroupMember) error {
			return g.db.BatchInsertGroupMember(ctx, values)
		},
		func(ctx context.Context, userIDs []string) ([]*model_struct.LocalGroupMember, bool, error) {
			localGroupMembers, err := g.db.GetGroupMemberListByUserIDs(ctx, groupID, filter, userIDs)
			if err != nil {
				return nil, false, err
			}
			switch filter {
			case constant.GroupFilterOwner:
				fallthrough
			case constant.GroupFilterAdmin:
				fallthrough
			case constant.GroupFilterOwnerAndAdmin:
				return localGroupMembers, false, nil
			case constant.GroupFilterAll:
				fallthrough
			case constant.GroupFilterOrdinaryUsers:
				fallthrough
			case constant.GroupFilterAdminAndOrdinaryUsers:
				return localGroupMembers, true, nil
			}
			return nil, false, sdkerrs.ErrArgs
		},
		func(ctx context.Context, userIDs []string) ([]*model_struct.LocalGroupMember, error) {
			serverGroupMember, err := g.GetDesignatedGroupMembers(ctx, groupID, userIDs)
			if err != nil {
				return nil, err
			}
			return datautil.Batch(ServerGroupMemberToLocalGroupMember, serverGroupMember), nil
		},
	)
	switch filter {
	case constant.GroupFilterOrdinaryUsers:
		groupOwnerAndGroupMember, err := g.db.GetGroupMemberListSplit(ctx, groupID, constant.GroupFilterOwnerAndAdmin, 0, 100)
		if err != nil {
			return nil, err
		}
		offset = offset + int32(len(groupOwnerAndGroupMember))
	case constant.GroupFilterAdminAndOrdinaryUsers:
		groupOwnerAndGroupMember, err := g.db.GetGroupMemberListSplit(ctx, groupID, constant.GroupFilterOwner, 0, 100)
		if err != nil {
			return nil, err
		}
		offset = offset + int32(len(groupOwnerAndGroupMember))
	}
	return dataFetcher.FetchWithPagination(ctx, int(offset), int(count))
}

func (g *Group) GetGroupApplicationListAsRecipient(ctx context.Context) ([]*model_struct.LocalAdminGroupRequest, error) {
	return g.db.GetAdminGroupApplication(ctx)
}

func (g *Group) GetGroupApplicationListAsApplicant(ctx context.Context) ([]*model_struct.LocalGroupRequest, error) {
	return g.db.GetSendGroupApplication(ctx)
}

func (g *Group) SearchGroupMembers(ctx context.Context, searchParam *sdk_params_callback.SearchGroupMembersParam) ([]*model_struct.LocalGroupMember, error) {
	return g.db.SearchGroupMembersDB(ctx, searchParam.KeywordList[0], searchParam.GroupID, searchParam.IsSearchMemberNickname, searchParam.IsSearchUserID, searchParam.Offset, searchParam.Count)
}

func (g *Group) IsJoinGroup(ctx context.Context, groupID string) (bool, error) {
	g.groupSyncMutex.Lock()
	defer g.groupSyncMutex.Unlock()

	lvs, err := g.db.GetVersionSync(ctx, g.groupTableName(), g.loginUserID)
	if err != nil {
		return false, err
	}
	if datautil.Contain(groupID, lvs.UIDList...) {
		return true, nil
	}
	return false, nil
}

func (g *Group) GetUsersInGroup(ctx context.Context, groupID string, userIDList []string) ([]string, error) {
	g.groupSyncMutex.Lock()
	defer g.groupSyncMutex.Unlock()

	lvs, err := g.db.GetVersionSync(ctx, g.groupTableName(), g.loginUserID)
	if err != nil {
		return nil, err
	}
	if !datautil.Contain(groupID, lvs.UIDList...) {
		return nil, nil
	}
	lvs, err = g.db.GetVersionSync(ctx, g.groupAndMemberVersionTableName(), groupID)
	if err != nil {
		return nil, err
	}

	groupMembersMap := datautil.SliceSetAny(lvs.UIDList, func(e string) string {
		return e
	})

	var usersInGroup []string
	for _, userID := range userIDList {
		if _, exists := groupMembersMap[userID]; exists {
			usersInGroup = append(usersInGroup, userID)
		}
	}

	return usersInGroup, nil
}

func (g *Group) InviteUserToGroup(ctx context.Context, groupID, reason string, userIDList []string) error {
	req := &group.InviteUserToGroupReq{GroupID: groupID, Reason: reason, InvitedUserIDs: userIDList}
	if err := g.inviteUserToGroup(ctx, req); err != nil {
		return err
	}

	g.groupSyncMutex.Lock()
	defer g.groupSyncMutex.Unlock()

	if err := g.IncrSyncGroupAndMember(ctx, groupID); err != nil {
		return err
	}
	return nil
}

func (g *Group) AcceptGroupApplication(ctx context.Context, groupID, fromUserID, handleMsg string) error {
	return g.HandlerGroupApplication(ctx, &group.GroupApplicationResponseReq{GroupID: groupID, FromUserID: fromUserID, HandledMsg: handleMsg, HandleResult: constant.GroupResponseAgree})
}

func (g *Group) RefuseGroupApplication(ctx context.Context, groupID, fromUserID, handleMsg string) error {
	return g.HandlerGroupApplication(ctx, &group.GroupApplicationResponseReq{GroupID: groupID, FromUserID: fromUserID, HandledMsg: handleMsg, HandleResult: constant.GroupResponseRefuse})
}

func (g *Group) HandlerGroupApplication(ctx context.Context, req *group.GroupApplicationResponseReq) error {
	if err := g.handlerGroupApplication(ctx, req); err != nil {
		return err
	}
	return nil
}

func (g *Group) GetGroupMemberNameAndFaceURL(ctx context.Context, groupID string, userIDs []string) (map[string]*model_struct.LocalGroupMember, error) {
	return g.GetGroupMembersInfo(ctx, groupID, userIDs)
}
