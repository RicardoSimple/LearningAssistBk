package graph

import (
	"context"
	"strings"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"

	"learning-assistant/dal"
	"learning-assistant/model"
	"learning-assistant/util/log"
)

func CreateContact(ctx context.Context, user *model.User) ([]*model.User, error) {
	records, err := QueryByUserName(ctx, user.Username)
	if err != nil {
		return nil, err
	}
	if len(records) > 0 {
		log.Info("[GRAPH] User %s already exists", user.Username)
		return nil, nil
	}
	log.Info("[GRAPH] create contact")
	result, err := generalExecute(ctx,
		"MERGE (u:User {username: $username,u_id:$u_id,gender:$gender,address:$address})"+
			formatReturn("u"),
		map[string]any{
			"username": user.Username,
			"u_id":     user.ID,
			"gender":   user.Gender,
			"address":  user.Address,
		})
	if err != nil {
		log.Info("[GRAPH] create error:%v", err)
		return nil, err
	}
	return resultToUser(result), nil
}

func QueryByUserName(ctx context.Context, username string) ([]*neo4j.Record, error) {
	log.Info("[GRAPH] query by user name")
	res, err := generalExecute(ctx,
		"MATCH (u:User{username:$username}) "+
			formatReturn("u"),
		map[string]any{
			"username": username,
		})
	if err != nil {
		log.Info("[GRAPH] query error:%v", err)
		return nil, err
	}
	return res.Records, nil
}

// LineByUserName 建立friend关系
func LineByUserName(ctx context.Context, toName, fromName string) ([]map[string]any, error) {
	log.Info("[GRAPH] line by user name")
	between, err2 := FindRelationBetween(ctx, toName, fromName)
	if err2 != nil {
		return nil, err2
	}
	if between != nil && len(between) >= 0 {
		log.Info("[GRAPH] exists relation")
		return nil, nil
	}
	res, err := generalExecute(ctx,
		"MATCH (tu:User{username:$toName})"+
			"MATCH (fu:User{username:$fName})"+
			"CREATE (tu)-[r:friend{time:$now}]->(fu)"+
			"RETURN tu,r,fu", map[string]any{
			"toName": toName,
			"fName":  fromName,
			"now":    time.Now().Format(time.DateTime),
		})
	if err != nil {
		log.Info("[GRAPH] line by user name error:%v", err)
		return nil, err
	}
	resMap := make([]map[string]any, 0, len(res.Records))
	for i := range res.Records {
		resMap = append(resMap, res.Records[i].AsMap())
	}
	return resMap, nil
}

// FindLinedFriend 查询好友
func FindLinedFriend(ctx context.Context, username string) ([]*model.User, error) {
	log.Info("[GRAPH] find friend")
	res, err := generalExecute(ctx,
		"MATCH (tu:User{username:$username})-[r:friend]-(fu:User) "+
			formatReturn("fu"), map[string]any{
			"username": username,
		})
	if err != nil {
		log.Info("[GRAPH] find friend error:%v", err)
		return nil, err
	}
	return resultToUser(res), nil
}

func FindRelationBetween(ctx context.Context, toName, fromName string) ([]*neo4j.Record, error) {
	result, err := generalExecute(ctx,
		"MATCH (toU:User{username:$toName})-[r:friend]-(fuU:User{username:$fromName})"+
			"RETURN r",
		map[string]any{
			"toName":   toName,
			"fromName": fromName,
		})
	if err != nil {
		log.Info("[GRAPH] find relation between error:%v", err)
		return nil, err
	}
	return result.Records, nil
}

// DeleteByUserName 通过username删除节点及相关关系
func DeleteByUserName(ctx context.Context, username string) error {
	log.Info("[GRAPH] delete by user name")
	_, err := generalExecute(ctx,
		"MATCH (u:User{username:$username})"+
			"DETACH DELETE u",
		map[string]any{
			"username": username,
		})
	if err != nil {
		log.Info("[GRAPH] delete error:%v", err)
		return err
	}
	return nil
}

// DeleteByUId 按照userid删除节点及相关关系
func DeleteByUId(ctx context.Context, uid string) error {
	log.Info("[GRAPH] delete by uid")
	_, err := generalExecute(ctx, "MATCH (u:User{u_id:$u_id})"+
		"DETACH DELETE u",
		map[string]any{
			"uid": uid,
		})
	if err != nil {
		log.Info("[GRAPH] delete error:%v", err)
		return err
	}
	return nil
}

func generalExecute(ctx context.Context, query string, params map[string]any) (*neo4j.EagerResult, error) {
	log.Info("[GRAPH] general execute")
	result, err := neo4j.ExecuteQuery(ctx, dal.GraphDB, query, params, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase(dal.GraphDataBaseName))
	return result, err
}

func formatReturn(rName string) string {
	returnStr := " RETURN *.username AS username,*.u_id as uId,*.gender AS gender"
	return strings.ReplaceAll(returnStr, "*", rName)
}
func resultToUser(result *neo4j.EagerResult) []*model.User {
	users := make([]*model.User, 0, len(result.Records))
	for _, record := range result.Records {
		m := record.AsMap()
		users = append(users, &model.User{
			Username: m["username"].(string),
			ID:       uint(m["uId"].(int64)),
			Gender:   m["gender"].(string),
		})
	}
	return users
}
