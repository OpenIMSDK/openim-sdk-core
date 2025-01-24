package test

import (
	"testing"

	"github.com/openimsdk/openim-sdk-core/v3/open_im_sdk"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/sdk_params_callback"
	"github.com/openimsdk/protocol/group"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/protocol/wrapperspb"
)

func Test_CreateGroupV2(t *testing.T) {
	req := &group.CreateGroupReq{
		MemberUserIDs: []string{"7299270930"},
		AdminUserIDs:  []string{"1"},
		OwnerUserID:   UserID,
		GroupInfo: &sdkws.GroupInfo{
			GroupName: "test",
			GroupType: 2,
		},
	}
	info, err := open_im_sdk.UserForSDK.Group().CreateGroup(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("group info: %s", info.String())
}

func Test_JoinGroup(t *testing.T) {
	err := open_im_sdk.UserForSDK.Group().JoinGroup(ctx, "3889561099", "1234", 3, "ex")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("InviteUserToGroup success")
}

func Test_QuitGroup(t *testing.T) {
	err := open_im_sdk.UserForSDK.Group().QuitGroup(ctx, "xadxwr24")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("QuitGroup success")
}

func Test_DismissGroup(t *testing.T) {
	err := open_im_sdk.UserForSDK.Group().DismissGroup(ctx, "1728503199")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("DismissGroup success")
}

func Test_ChangeGroupMute(t *testing.T) {
	err := open_im_sdk.UserForSDK.Group().ChangeGroupMute(ctx, "3459296007", true)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("ChangeGroupMute success", ctx.Value("operationID"))
}

func Test_CancelMuteGroup(t *testing.T) {
	err := open_im_sdk.UserForSDK.Group().ChangeGroupMute(ctx, "3459296007", false)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("ChangeGroupMute success", ctx.Value("operationID"))
}

func Test_ChangeGroupMemberMute(t *testing.T) {
	err := open_im_sdk.UserForSDK.Group().ChangeGroupMemberMute(ctx, "3459296007", UserID, 10000)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("ChangeGroupMute success", ctx.Value("operationID"))
}

func Test_CancelChangeGroupMemberMute(t *testing.T) {
	err := open_im_sdk.UserForSDK.Group().ChangeGroupMemberMute(ctx, "3459296007", UserID, 0)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("CancelChangeGroupMemberMute success", ctx.Value("operationID"))
}

func Test_SetGroupMemberInfo(t *testing.T) {
	err := open_im_sdk.UserForSDK.Group().SetGroupMemberInfo(ctx, &group.SetGroupMemberInfo{
		GroupID:  "3889561099",
		UserID:   UserID,
		FaceURL:  wrapperspb.String("https://doc.rentsoft.cn/images/logo.png"),
		Nickname: wrapperspb.String("testupdatename"),
		Ex:       wrapperspb.String("a"),
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log("SetGroupMemberNickname success", ctx.Value("operationID"))
}

func Test_GetJoinedGroupList(t *testing.T) {
	info, err := open_im_sdk.UserForSDK.Group().GetJoinedGroupList(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("GetJoinedGroupList: %d\n", len(info))
	for _, localGroup := range info {
		t.Logf("%#v", localGroup)
	}
}

func Test_GetSpecifiedGroupsInfo(t *testing.T) {
	info, err := open_im_sdk.UserForSDK.Group().GetSpecifiedGroupsInfo(ctx, []string{"test"})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("GetGroupsInfo: %d\n", len(info))
	for _, localGroup := range info {
		t.Logf("%#v", localGroup)
	}
}

func Test_GetGroupApplicationListAsRecipient(t *testing.T) {
	info, err := open_im_sdk.UserForSDK.Group().GetGroupApplicationListAsRecipient(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("GetRecvGroupApplicationList: %d\n", len(info))
	for _, localGroup := range info {
		t.Logf("%#v", localGroup)
	}
}

func Test_GetGroupApplicationListAsApplicant(t *testing.T) {
	info, err := open_im_sdk.UserForSDK.Group().GetGroupApplicationListAsApplicant(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("GetSendGroupApplicationList: %d\n", len(info))
	for _, localGroup := range info {
		t.Logf("%#v", localGroup)
	}
}

func Test_AcceptGroupApplication(t *testing.T) {
	err := open_im_sdk.UserForSDK.Group().AcceptGroupApplication(ctx, "3459296007", "863454357", "test accept")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("AcceptGroupApplication success", ctx.Value("operationID"))
}

func Test_RefuseGroupApplication(t *testing.T) {
	t.Log("operationID:", ctx.Value("operationID"))
	err := open_im_sdk.UserForSDK.Group().RefuseGroupApplication(ctx, "3459296007", "863454357", "test refuse")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("AcceptGroupApplication success")
}

func Test_HandlerGroupApplication(t *testing.T) {
	err := open_im_sdk.UserForSDK.Group().HandlerGroupApplication(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("AcceptGroupApplication success", ctx.Value("operationID"))
}

func Test_SearchGroupMembers(t *testing.T) {
	info, err := open_im_sdk.UserForSDK.Group().SearchGroupMembers(ctx, &sdk_params_callback.SearchGroupMembersParam{
		GroupID:                "3459296007",
		KeywordList:            []string{""},
		IsSearchUserID:         false,
		IsSearchMemberNickname: false,
		Offset:                 0,
		Count:                  10,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("SearchGroupMembers: %d\n", len(info))
	for _, localGroup := range info {
		t.Logf("%#v", localGroup)
	}
}

func Test_KickGroupMember(t *testing.T) {
	err := open_im_sdk.UserForSDK.Group().KickGroupMember(ctx, "3459296007", "test", []string{"863454357"})
	if err != nil {
		t.Fatal(err)
	}
	t.Log("InviteUserToGroup success", ctx.Value("operationID"))
}

func Test_TransferGroupOwner(t *testing.T) {
	err := open_im_sdk.UserForSDK.Group().TransferGroupOwner(ctx, "1728503199", "5226390099")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("TransferGroupOwner success", ctx.Value("operationID"))
}

func Test_InviteUserToGroup(t *testing.T) {
	err := open_im_sdk.UserForSDK.Group().InviteUserToGroup(ctx, "3459296007", "test", []string{"45644221123"})
	if err != nil {
		t.Fatal(err)
	}
	t.Log("InviteUserToGroup success", ctx.Value("operationID"))
}

func Test_GetGroup(t *testing.T) {
	t.Log("--------------------------")
	infos, err := open_im_sdk.UserForSDK.Group().GetSpecifiedGroupsInfo(ctx, []string{"3179997540"})
	if err != nil {
		t.Fatal(err)
	}
	for i, info := range infos {
		t.Logf("%d: %#v", i, info)
	}
}
func Test_GetGroupApplicantsList(t *testing.T) {
	t.Log("--------------------------")
	infos, err := open_im_sdk.UserForSDK.Group().GetGroupApplicationListAsRecipient(ctx)
	if err != nil {
		t.Fatal(err)
	}
	for i, info := range infos {
		t.Logf("%d: %#v", i, info)
	}
	// time.Sleep(time.Second * 100000)
}
func Test_IsJoinGroup(t *testing.T) {
	t.Log("--------------------------")
	join, err := open_im_sdk.UserForSDK.Group().IsJoinGroup(ctx, "3889561099")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("join:", join)
}

func Test_GetGroupMemberList(t *testing.T) {
	t.Log("--------------------------")
	m := map[int32]string{
		constant.GroupOwner:         "Group Owner",
		constant.GroupAdmin:         "Administrator",
		constant.GroupOrdinaryUsers: "Members",
	}

	members, err := open_im_sdk.UserForSDK.Group().GetGroupMemberList(ctx, "3889561099", 0, 0, 9999999)
	if err != nil {
		panic(err)
	}
	for i, member := range members {
		name := m[member.RoleLevel]
		t.Log(i, member.UserID, member.Nickname, name)
	}

	t.Log("--------------------------")
}

func Test_SetGroupInfo(t *testing.T) {
	err := open_im_sdk.UserForSDK.Group().SetGroupInfo(ctx, &group.SetGroupInfoExReq{
		GroupID: "3889561099",
		Ex:      &wrapperspb.StringValue{Value: "groupex"},
	})
	if err != nil {
		t.Fatal(err)
	}
}

func Test_GetJoinedGroupListPage(t *testing.T) {
	t.Log("-----------------------")
	info, err := open_im_sdk.UserForSDK.Group().GetJoinedGroupListPage(ctx, 0, 10)
	t.Log("-----------------------")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(info)
}
