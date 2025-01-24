package group

import (
	"context"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/utils/datautil"
)

func (g *Group) buildGroupMemberKey(groupID string, userID string) string {
	return groupID + ":" + userID
}

func (g *Group) GetGroupMembersInfoFunc(ctx context.Context, groupID string, userIDs []string,
	fetchFunc func(ctx context.Context, missingKeys []string) ([]*model_struct.LocalGroupMember, error)) (map[string]*model_struct.LocalGroupMember, error) {
	var (
		res         = make(map[string]*model_struct.LocalGroupMember)
		missingKeys []string
	)

	for _, userID := range userIDs {
		key := g.buildGroupMemberKey(groupID, userID)
		if member, ok := g.groupMemberCache.Load(key); ok {
			res[userID] = member
		} else {
			missingKeys = append(missingKeys, userIDs...)
		}
	}

	log.ZDebug(ctx, "GetGroupMembersInfoFunc fetch", "missingKeys", missingKeys)
	fetchData, err := fetchFunc(ctx, missingKeys)
	if err != nil {
		return nil, err
	}

	for i, data := range fetchData {
		key := g.buildGroupMemberKey(groupID, data.UserID)
		res[data.UserID] = fetchData[i]
		g.groupMemberCache.Store(key, fetchData[i])
	}

	return res, nil
}

func (g *Group) GetGroupMembersInfo(ctx context.Context, groupID string, userIDs []string) (map[string]*model_struct.LocalGroupMember, error) {
	return g.GetGroupMembersInfoFunc(ctx, groupID, userIDs, func(ctx context.Context, dbKeys []string) ([]*model_struct.LocalGroupMember, error) {
		if len(dbKeys) == 0 {
			return nil, nil
		}
		dbData, err := g.db.GetGroupSomeMemberInfo(ctx, groupID, dbKeys)
		if err != nil {
			return nil, err
		}
		queryKeys := datautil.SliceSubAny(dbKeys, dbData, func(t *model_struct.LocalGroupMember) string {
			return t.UserID
		})
		if len(queryKeys) != 0 {
			queryData, err := g.GetDesignatedGroupMembers(ctx, groupID, queryKeys)
			if err != nil {
				return nil, err
			}

			dbData = append(dbData, datautil.Batch(ServerGroupMemberToLocalGroupMember, queryData)...)
		}
		return dbData, nil
	})
}
