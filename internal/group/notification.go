package group

import (
	"context"
	"fmt"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/utils"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/utils/datautil"

	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/log"
)

const (
	groupSortIDUnchanged = 0
	groupSortIDChanged   = 1
)

func (g *Group) DoNotification(ctx context.Context, msg *sdkws.MsgData) {
	go func() {
		if err := g.doNotification(ctx, msg); err != nil {
			log.ZError(ctx, "DoGroupNotification failed", err)
		}
	}()
}

func (g *Group) doNotification(ctx context.Context, msg *sdkws.MsgData) error {
	switch msg.ContentType {
	case constant.GroupApplicationAcceptedNotification: // 1505
		var detail sdkws.GroupApplicationAcceptedTips
		if err := utils.UnmarshalNotificationElem(msg.Content, &detail); err != nil {
			return err
		}
		switch detail.ReceiverAs {
		case constant.ApplicantReceiver:
			return g.SyncAllSelfGroupApplication(ctx)
		case constant.AdminReceiver:
			return g.SyncAdminGroupApplications(ctx, detail.Group.GroupID)
		default:
			return errs.New(fmt.Sprintf("GroupApplicationAcceptedNotification ReceiverAs unknown %d", detail.ReceiverAs)).Wrap()
		}

	case constant.GroupApplicationRejectedNotification: // 1506
		var detail sdkws.GroupApplicationRejectedTips
		if err := utils.UnmarshalNotificationElem(msg.Content, &detail); err != nil {
			return err
		}
		switch detail.ReceiverAs {
		case 0:
			return g.SyncAllSelfGroupApplication(ctx)
		case 1:
			return g.SyncAdminGroupApplications(ctx, detail.Group.GroupID)
		default:
			return errs.New(fmt.Sprintf("GroupApplicationRejectedNotification ReceiverAs unknown %d", detail.ReceiverAs)).Wrap()
		}

	default:
		g.groupSyncMutex.Lock()
		defer g.groupSyncMutex.Unlock()
		switch msg.ContentType {
		case constant.GroupCreatedNotification: // 1501
			var detail sdkws.GroupCreatedTips
			if err := utils.UnmarshalNotificationElem(msg.Content, &detail); err != nil {
				return err
			}

			if err := g.IncrSyncJoinGroup(ctx); err != nil {
				return err
			}
			return g.IncrSyncGroupAndMember(ctx, detail.Group.GroupID)

		case constant.GroupInfoSetNotification: // 1502
			var detail sdkws.GroupInfoSetTips
			if err := utils.UnmarshalNotificationElem(msg.Content, &detail); err != nil {
				return err
			}
			return g.onlineSyncGroupAndMember(ctx, detail.Group.GroupID, nil,
				nil, nil, detail.Group, groupSortIDUnchanged, detail.GroupMemberVersion, detail.GroupMemberVersionID)
		case constant.JoinGroupApplicationNotification: // 1503
			var detail sdkws.JoinGroupApplicationTips
			if err := utils.UnmarshalNotificationElem(msg.Content, &detail); err != nil {
				return err
			}
			if detail.Applicant.UserID == g.loginUserID {
				return g.SyncSelfGroupApplications(ctx, detail.Group.GroupID)
			} else {
				return g.SyncAdminGroupApplications(ctx, detail.Group.GroupID)
			}
		case constant.MemberQuitNotification: // 1504
			var detail sdkws.MemberQuitTips
			if err := utils.UnmarshalNotificationElem(msg.Content, &detail); err != nil {
				return err
			}
			if detail.QuitUser.UserID == g.loginUserID {
				return g.IncrSyncJoinGroup(ctx)
			} else {
				return g.onlineSyncGroupAndMember(ctx, detail.Group.GroupID, []*sdkws.GroupMemberFullInfo{detail.QuitUser}, nil,
					nil, detail.Group, groupSortIDUnchanged, detail.GroupMemberVersion, detail.GroupMemberVersionID)
			}
		case constant.GroupOwnerTransferredNotification: // 1507
			var detail sdkws.GroupOwnerTransferredTips
			if err := utils.UnmarshalNotificationElem(msg.Content, &detail); err != nil {
				return err
			}
			if detail.Group == nil {
				return errs.New(fmt.Sprintf("group is nil, groupID: %s", detail.Group.GroupID)).Wrap()
			}
			if detail.NewGroupOwner.RoleLevel < constant.GroupAdmin && detail.OldGroupOwner == g.loginUserID {
				if err := g.delLocalGroupRequest(ctx, detail.Group.GroupID, g.loginUserID); err != nil {
					return err
				}
			}
			return g.onlineSyncGroupAndMember(ctx, detail.Group.GroupID, nil,
				[]*sdkws.GroupMemberFullInfo{detail.NewGroupOwner, detail.OldGroupOwnerInfo}, nil,
				detail.Group, groupSortIDChanged, detail.GroupMemberVersion, detail.GroupMemberVersionID)
		case constant.MemberKickedNotification: // 1508
			var detail sdkws.MemberKickedTips
			if err := utils.UnmarshalNotificationElem(msg.Content, &detail); err != nil {
				return err
			}
			var self bool
			for _, info := range detail.KickedUserList {
				if info.UserID == g.loginUserID {
					self = true
					break
				}
			}
			if self {
				if err := g.delLocalGroupRequest(ctx, detail.Group.GroupID, g.loginUserID); err != nil {
					return err
				}
				return g.IncrSyncJoinGroup(ctx)
			} else {
				return g.onlineSyncGroupAndMember(ctx, detail.Group.GroupID, detail.KickedUserList, nil,
					nil, detail.Group, groupSortIDUnchanged, detail.GroupMemberVersion, detail.GroupMemberVersionID)
			}
		case constant.MemberInvitedNotification: // 1509
			var detail sdkws.MemberInvitedTips
			if err := utils.UnmarshalNotificationElem(msg.Content, &detail); err != nil {
				return err
			}
			userIDMap := datautil.SliceSetAny(detail.InvitedUserList, func(e *sdkws.GroupMemberFullInfo) string {
				return e.UserID
			})
			//Also invited as a member
			if _, ok := userIDMap[g.loginUserID]; ok {
				if err := g.IncrSyncJoinGroup(ctx); err != nil {
					return err
				}
				return g.IncrSyncGroupAndMember(ctx, detail.Group.GroupID)
			} else {
				return g.onlineSyncGroupAndMember(ctx, detail.Group.GroupID, nil, nil,
					detail.InvitedUserList, detail.Group, groupSortIDUnchanged, detail.GroupMemberVersion, detail.GroupMemberVersionID)
			}
		case constant.MemberEnterNotification: // 1510
			var detail sdkws.MemberEnterTips
			if err := utils.UnmarshalNotificationElem(msg.Content, &detail); err != nil {
				return err
			}
			if detail.EntrantUser.UserID == g.loginUserID {
				if err := g.IncrSyncJoinGroup(ctx); err != nil {
					return err
				}
				return g.IncrSyncGroupAndMember(ctx, detail.Group.GroupID)
			} else {
				return g.onlineSyncGroupAndMember(ctx, detail.Group.GroupID, nil, nil,
					[]*sdkws.GroupMemberFullInfo{detail.EntrantUser}, detail.Group, groupSortIDUnchanged, detail.GroupMemberVersion, detail.GroupMemberVersionID)
			}
		case constant.GroupDismissedNotification: // 1511
			var detail sdkws.GroupDismissedTips
			if err := utils.UnmarshalNotificationElem(msg.Content, &detail); err != nil {
				return err
			}
			g.listener().OnGroupDismissed(utils.StructToJsonString(detail.Group))

			return g.IncrSyncJoinGroup(ctx)
		case constant.GroupMemberMutedNotification: // 1512
			var detail sdkws.GroupMemberMutedTips
			if err := utils.UnmarshalNotificationElem(msg.Content, &detail); err != nil {
				return err
			}
			return g.onlineSyncGroupAndMember(ctx, detail.Group.GroupID, nil,
				[]*sdkws.GroupMemberFullInfo{detail.MutedUser}, nil, nil,
				groupSortIDUnchanged, detail.GroupMemberVersion, detail.GroupMemberVersionID)
		case constant.GroupMemberCancelMutedNotification: // 1513
			var detail sdkws.GroupMemberCancelMutedTips
			if err := utils.UnmarshalNotificationElem(msg.Content, &detail); err != nil {
				return err
			}
			return g.onlineSyncGroupAndMember(ctx, detail.Group.GroupID, nil,
				[]*sdkws.GroupMemberFullInfo{detail.MutedUser}, nil, nil,
				groupSortIDUnchanged, detail.GroupMemberVersion, detail.GroupMemberVersionID)
		case constant.GroupMutedNotification: // 1514
			var detail sdkws.GroupMutedTips
			if err := utils.UnmarshalNotificationElem(msg.Content, &detail); err != nil {
				return err
			}
			return g.onlineSyncGroupAndMember(ctx, detail.Group.GroupID, nil, nil,
				nil, detail.Group, groupSortIDUnchanged, detail.GroupMemberVersion, detail.GroupMemberVersionID)
		case constant.GroupCancelMutedNotification: // 1515
			var detail sdkws.GroupCancelMutedTips
			if err := utils.UnmarshalNotificationElem(msg.Content, &detail); err != nil {
				return err
			}
			return g.onlineSyncGroupAndMember(ctx, detail.Group.GroupID, nil, nil,
				nil, detail.Group, groupSortIDUnchanged, detail.GroupMemberVersion, detail.GroupMemberVersionID)
		case constant.GroupMemberInfoSetNotification: // 1516
			var detail sdkws.GroupMemberInfoSetTips
			if err := utils.UnmarshalNotificationElem(msg.Content, &detail); err != nil {
				return err
			}
			if detail.ChangedUser.RoleLevel < constant.GroupAdmin && detail.ChangedUser.UserID == g.loginUserID {
				if err := g.delLocalGroupRequest(ctx, detail.Group.GroupID, g.loginUserID); err != nil {
					return err
				}
			}
			return g.onlineSyncGroupAndMember(ctx, detail.Group.GroupID, nil,
				[]*sdkws.GroupMemberFullInfo{detail.ChangedUser}, nil, nil,
				detail.GroupSortVersion, detail.GroupMemberVersion, detail.GroupMemberVersionID)
		case constant.GroupMemberSetToAdminNotification: // 1517
			var detail sdkws.GroupMemberInfoSetTips
			if err := utils.UnmarshalNotificationElem(msg.Content, &detail); err != nil {
				return err
			}
			return g.onlineSyncGroupAndMember(ctx, detail.Group.GroupID, nil,
				[]*sdkws.GroupMemberFullInfo{detail.ChangedUser}, nil, nil,
				detail.GroupSortVersion, detail.GroupMemberVersion, detail.GroupMemberVersionID)
		case constant.GroupMemberSetToOrdinaryUserNotification: // 1518
			var detail sdkws.GroupMemberInfoSetTips
			if err := utils.UnmarshalNotificationElem(msg.Content, &detail); err != nil {
				return err
			}
			if detail.ChangedUser.UserID == g.loginUserID {
				if err := g.delLocalGroupRequest(ctx, detail.Group.GroupID, g.loginUserID); err != nil {
					return err
				}
			}
			return g.onlineSyncGroupAndMember(ctx, detail.Group.GroupID, nil,
				[]*sdkws.GroupMemberFullInfo{detail.ChangedUser}, nil, nil,
				detail.GroupSortVersion, detail.GroupMemberVersion, detail.GroupMemberVersionID)
		case constant.GroupInfoSetAnnouncementNotification: // 1519
			var detail sdkws.GroupInfoSetAnnouncementTips //
			if err := utils.UnmarshalNotificationElem(msg.Content, &detail); err != nil {
				return err
			}
			return g.onlineSyncGroupAndMember(ctx, detail.Group.GroupID, nil, nil,
				nil, detail.Group, groupSortIDUnchanged, detail.GroupMemberVersion, detail.GroupMemberVersionID)
		case constant.GroupInfoSetNameNotification: // 1520
			var detail sdkws.GroupInfoSetNameTips //
			if err := utils.UnmarshalNotificationElem(msg.Content, &detail); err != nil {
				return err
			}
			return g.onlineSyncGroupAndMember(ctx, detail.Group.GroupID, nil, nil,
				nil, detail.Group, groupSortIDUnchanged, detail.GroupMemberVersion, detail.GroupMemberVersionID)
		default:
			return errs.New("unknown tips type", "contentType", msg.ContentType).Wrap()
		}
	}
}
