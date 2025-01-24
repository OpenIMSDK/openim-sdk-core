package group

import (
	"context"
	"sync"

	"github.com/openimsdk/openim-sdk-core/v3/open_im_sdk_callback"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/api"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/cache"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/common"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/datafetcher"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/db_interface"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/page"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/sdkerrs"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/syncer"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/utils"
	"github.com/openimsdk/protocol/group"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/utils/datautil"
)

const (
	groupSyncLimit       = 1000
	groupMemberSyncLimit = 1000
)

func NewGroup(loginUserID string, db db_interface.DataBase,
	conversationCh chan common.Cmd2Value) *Group {
	g := &Group{
		loginUserID:    loginUserID,
		db:             db,
		conversationCh: conversationCh,
	}
	g.initSyncer()
	g.groupMemberCache = cache.NewCache[string, *model_struct.LocalGroupMember]()
	return g
}

type Group struct {
	listener                func() open_im_sdk_callback.OnGroupListener
	loginUserID             string
	db                      db_interface.DataBase
	groupSyncer             *syncer.Syncer[*model_struct.LocalGroup, group.GetJoinedGroupListResp, string]
	groupMemberSyncer       *syncer.Syncer[*model_struct.LocalGroupMember, group.GetGroupMemberListResp, [2]string]
	groupRequestSyncer      *syncer.Syncer[*model_struct.LocalGroupRequest, syncer.NoResp, [2]string]
	groupAdminRequestSyncer *syncer.Syncer[*model_struct.LocalAdminGroupRequest, syncer.NoResp, [2]string]

	conversationCh chan common.Cmd2Value

	groupSyncMutex     sync.Mutex
	listenerForService open_im_sdk_callback.OnListenerForService

	groupMemberCache *cache.Cache[string, *model_struct.LocalGroupMember]
}

func (g *Group) initSyncer() {
	g.groupSyncer = syncer.New2[*model_struct.LocalGroup, group.GetJoinedGroupListResp, string](
		syncer.WithInsert[*model_struct.LocalGroup, group.GetJoinedGroupListResp, string](func(ctx context.Context, value *model_struct.LocalGroup) error {
			return g.db.InsertGroup(ctx, value)
		}),
		syncer.WithDelete[*model_struct.LocalGroup, group.GetJoinedGroupListResp, string](func(ctx context.Context, value *model_struct.LocalGroup) error {
			if err := g.db.DeleteGroupAllMembers(ctx, value.GroupID); err != nil {
				return err
			}
			if err := g.db.DeleteVersionSync(ctx, g.groupAndMemberVersionTableName(), value.GroupID); err != nil {
				return err
			}
			return g.db.DeleteGroup(ctx, value.GroupID)
		}),
		syncer.WithUpdate[*model_struct.LocalGroup, group.GetJoinedGroupListResp, string](func(ctx context.Context, server, local *model_struct.LocalGroup) error {
			log.ZInfo(ctx, "groupSyncer trigger update function", "groupID", server.GroupID, "server", server, "local", local)
			return g.db.UpdateGroup(ctx, server)
		}),
		syncer.WithUUID[*model_struct.LocalGroup, group.GetJoinedGroupListResp, string](func(value *model_struct.LocalGroup) string {
			return value.GroupID
		}),
		syncer.WithNotice[*model_struct.LocalGroup, group.GetJoinedGroupListResp, string](func(ctx context.Context, state int, server, local *model_struct.LocalGroup) error {
			switch state {
			case syncer.Insert:
				// when a user kicked to the group and invited to the group again, group info maybe updated,
				// so conversation info need to be updated
				g.listener().OnJoinedGroupAdded(utils.StructToJsonString(server))
				_ = common.TriggerCmdUpdateConversation(ctx, common.UpdateConNode{
					Action: constant.UpdateConFaceUrlAndNickName,
					Args: common.SourceIDAndSessionType{
						SourceID: server.GroupID, SessionType: constant.ReadGroupChatType,
						FaceURL: server.FaceURL, Nickname: server.GroupName,
					},
				}, g.conversationCh)
			case syncer.Delete:
				local.MemberCount = 0
				g.listener().OnJoinedGroupDeleted(utils.StructToJsonString(local))
			case syncer.Update:
				log.ZInfo(ctx, "groupSyncer trigger update", "groupID",
					server.GroupID, "data", server, "isDismissed", server.Status == constant.GroupStatusDismissed)
				if server.Status == constant.GroupStatusDismissed {
					if err := g.db.DeleteGroupAllMembers(ctx, server.GroupID); err != nil {
						log.ZError(ctx, "delete group all members failed", err)
					}
					g.listener().OnGroupDismissed(utils.StructToJsonString(server))
				} else {
					g.listener().OnGroupInfoChanged(utils.StructToJsonString(server))
					if server.GroupName != local.GroupName || local.FaceURL != server.FaceURL {
						_ = common.TriggerCmdUpdateConversation(ctx, common.UpdateConNode{
							Action: constant.UpdateConFaceUrlAndNickName,
							Args: common.SourceIDAndSessionType{
								SourceID: server.GroupID, SessionType: constant.ReadGroupChatType,
								FaceURL: server.FaceURL, Nickname: server.GroupName,
							},
						}, g.conversationCh)
					}
				}
			}
			return nil
		}),

		syncer.WithBatchInsert[*model_struct.LocalGroup, group.GetJoinedGroupListResp, string](func(ctx context.Context, values []*model_struct.LocalGroup) error {
			return g.db.BatchInsertGroup(ctx, values)
		}),
		syncer.WithDeleteAll[*model_struct.LocalGroup, group.GetJoinedGroupListResp, string](func(ctx context.Context, _ string) error {
			return g.db.DeleteAllGroup(ctx)
		}),
		syncer.WithBatchPageReq[*model_struct.LocalGroup, group.GetJoinedGroupListResp, string](func(entityID string) page.PageReq {
			return &group.GetJoinedGroupListReq{FromUserID: entityID,
				Pagination: &sdkws.RequestPagination{ShowNumber: 100}}
		}),
		syncer.WithBatchPageRespConvertFunc[*model_struct.LocalGroup, group.GetJoinedGroupListResp, string](func(resp *group.GetJoinedGroupListResp) []*model_struct.LocalGroup {
			return datautil.Batch(ServerGroupToLocalGroup, resp.Groups)
		}),
		syncer.WithReqApiRouter[*model_struct.LocalGroup, group.GetJoinedGroupListResp, string](api.GetJoinedGroupList.Route()),
		syncer.WithFullSyncLimit[*model_struct.LocalGroup, group.GetJoinedGroupListResp, string](groupSyncLimit),
	)
	g.groupMemberSyncer = syncer.New2[*model_struct.LocalGroupMember, group.GetGroupMemberListResp, [2]string](
		syncer.WithInsert[*model_struct.LocalGroupMember, group.GetGroupMemberListResp, [2]string](func(ctx context.Context, value *model_struct.LocalGroupMember) error {
			return g.db.InsertGroupMember(ctx, value)
		}),
		syncer.WithDelete[*model_struct.LocalGroupMember, group.GetGroupMemberListResp, [2]string](func(ctx context.Context, value *model_struct.LocalGroupMember) error {
			return g.db.DeleteGroupMember(ctx, value.GroupID, value.UserID)
		}),
		syncer.WithUpdate[*model_struct.LocalGroupMember, group.GetGroupMemberListResp, [2]string](func(ctx context.Context, server, local *model_struct.LocalGroupMember) error {
			g.groupMemberCache.Delete(g.buildGroupMemberKey(server.GroupID, server.UserID))
			return g.db.UpdateGroupMember(ctx, server)
		}),
		syncer.WithUUID[*model_struct.LocalGroupMember, group.GetGroupMemberListResp, [2]string](func(value *model_struct.LocalGroupMember) [2]string {
			return [...]string{value.GroupID, value.UserID}
		}),
		syncer.WithNotice[*model_struct.LocalGroupMember, group.GetGroupMemberListResp, [2]string](func(ctx context.Context, state int, server, local *model_struct.LocalGroupMember) error {
			switch state {
			case syncer.Insert:
				g.listener().OnGroupMemberAdded(utils.StructToJsonString(server))
				// When a user is kicked and invited to the group again, group member info will be updated.
				_ = common.TriggerCmdUpdateMessage(ctx,
					common.UpdateMessageNode{
						Action: constant.UpdateMsgFaceUrlAndNickName,
						Args: common.UpdateMessageInfo{
							SessionType: constant.ReadGroupChatType, UserID: server.UserID, FaceURL: server.FaceURL,
							Nickname: server.Nickname, GroupID: server.GroupID,
						},
					}, g.conversationCh)
			case syncer.Delete:
				g.listener().OnGroupMemberDeleted(utils.StructToJsonString(local))
			case syncer.Update:
				g.listener().OnGroupMemberInfoChanged(utils.StructToJsonString(server))
				if server.Nickname != local.Nickname || server.FaceURL != local.FaceURL {
					_ = common.TriggerCmdUpdateMessage(ctx,
						common.UpdateMessageNode{
							Action: constant.UpdateMsgFaceUrlAndNickName,
							Args: common.UpdateMessageInfo{
								SessionType: constant.ReadGroupChatType, UserID: server.UserID, FaceURL: server.FaceURL,
								Nickname: server.Nickname, GroupID: server.GroupID,
							},
						}, g.conversationCh)
					_ = common.TriggerCmdUpdateConversation(ctx, common.UpdateConNode{Action: constant.UpdateLatestMessageFaceUrlAndNickName, Args: common.UpdateMessageInfo{
						SessionType: constant.ReadGroupChatType, UserID: server.UserID, FaceURL: server.FaceURL,
						Nickname: server.Nickname, GroupID: server.GroupID,
					}}, g.conversationCh)
				}
			}
			return nil
		}),
		syncer.WithBatchInsert[*model_struct.LocalGroupMember, group.GetGroupMemberListResp, [2]string](func(ctx context.Context, values []*model_struct.LocalGroupMember) error {
			return g.db.BatchInsertGroupMember(ctx, values)
		}),
		syncer.WithDeleteAll[*model_struct.LocalGroupMember, group.GetGroupMemberListResp, [2]string](func(ctx context.Context, groupID string) error {
			return g.db.DeleteGroupAllMembers(ctx, groupID)
		}),
		syncer.WithBatchPageReq[*model_struct.LocalGroupMember, group.GetGroupMemberListResp, [2]string](func(entityID string) page.PageReq {
			return &group.GetGroupMemberListReq{GroupID: entityID, Pagination: &sdkws.RequestPagination{ShowNumber: 100}}
		}),
		syncer.WithBatchPageRespConvertFunc[*model_struct.LocalGroupMember, group.GetGroupMemberListResp, [2]string](func(resp *group.GetGroupMemberListResp) []*model_struct.LocalGroupMember {
			return datautil.Batch(ServerGroupMemberToLocalGroupMember, resp.Members)
		}),
		syncer.WithReqApiRouter[*model_struct.LocalGroupMember, group.GetGroupMemberListResp, [2]string](api.GetGroupMemberList.Route()),
		syncer.WithFullSyncLimit[*model_struct.LocalGroupMember, group.GetGroupMemberListResp, [2]string](groupMemberSyncLimit),
	)

	g.groupRequestSyncer = syncer.New[*model_struct.LocalGroupRequest, syncer.NoResp, [2]string](func(ctx context.Context, value *model_struct.LocalGroupRequest) error {
		return g.db.InsertGroupRequest(ctx, value)
	}, func(ctx context.Context, value *model_struct.LocalGroupRequest) error {
		return g.db.DeleteGroupRequest(ctx, value.GroupID, value.UserID)
	}, func(ctx context.Context, server, local *model_struct.LocalGroupRequest) error {
		return g.db.UpdateGroupRequest(ctx, server)
	}, func(value *model_struct.LocalGroupRequest) [2]string {
		return [...]string{value.GroupID, value.UserID}
	}, nil, func(ctx context.Context, state int, server, local *model_struct.LocalGroupRequest) error {
		switch state {
		case syncer.Insert:
			g.listener().OnGroupApplicationAdded(utils.StructToJsonString(server))
		case syncer.Update:
			switch server.HandleResult {
			case constant.FriendResponseAgree:
				g.listener().OnGroupApplicationAccepted(utils.StructToJsonString(server))
			case constant.FriendResponseRefuse:
				g.listener().OnGroupApplicationRejected(utils.StructToJsonString(server))
			default:
				g.listener().OnGroupApplicationAdded(utils.StructToJsonString(server))
			}
		}
		return nil
	})

	g.groupAdminRequestSyncer = syncer.New[*model_struct.LocalAdminGroupRequest, syncer.NoResp, [2]string](func(ctx context.Context, value *model_struct.LocalAdminGroupRequest) error {
		return g.db.InsertAdminGroupRequest(ctx, value)
	}, func(ctx context.Context, value *model_struct.LocalAdminGroupRequest) error {
		return g.db.DeleteAdminGroupRequest(ctx, value.GroupID, value.UserID)
	}, func(ctx context.Context, server, local *model_struct.LocalAdminGroupRequest) error {
		return g.db.UpdateAdminGroupRequest(ctx, server)
	}, func(value *model_struct.LocalAdminGroupRequest) [2]string {
		return [...]string{value.GroupID, value.UserID}
	}, nil, func(ctx context.Context, state int, server, local *model_struct.LocalAdminGroupRequest) error {
		switch state {
		case syncer.Insert:
			g.listener().OnGroupApplicationAdded(utils.StructToJsonString(server))
		case syncer.Update:
			switch server.HandleResult {
			case constant.FriendResponseAgree:
				g.listener().OnGroupApplicationAccepted(utils.StructToJsonString(server))
			case constant.FriendResponseRefuse:
				g.listener().OnGroupApplicationRejected(utils.StructToJsonString(server))
			default:
				g.listener().OnGroupApplicationAdded(utils.StructToJsonString(server))
			}
		}
		return nil
	})

}

func (g *Group) SetGroupListener(listener func() open_im_sdk_callback.OnGroupListener) {
	g.listener = listener
}

func (g *Group) SetListenerForService(listener open_im_sdk_callback.OnListenerForService) {
	g.listenerForService = listener
}

func (g *Group) FetchGroupOrError(ctx context.Context, groupID string) (*model_struct.LocalGroup, error) {
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
	groups, err := dataFetcher.FetchMissingAndCombineLocal(ctx, []string{groupID})
	if err != nil {
		return nil, err
	}
	if len(groups) == 0 {
		return nil, sdkerrs.ErrGroupIDNotFound.WrapMsg("sdk and server not this group")
	}
	return groups[0], nil
}

func (g *Group) delLocalGroupRequest(ctx context.Context, groupID, userID string) error {
	return g.db.DeleteGroupRequest(ctx, groupID, userID)
}
