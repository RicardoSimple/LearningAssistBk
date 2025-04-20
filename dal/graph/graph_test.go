package graph

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"

	"ar-app-api/conf"
	"ar-app-api/dal"
	"ar-app-api/util/log"
)

type GraphTest struct {
	suite.Suite
}

func (t *GraphTest) SetupTest() {
	// init neo4j
	ctx := context.Background()
	log.InitLogging("", log.DEBUG)
	conf.Init(ctx)
	dal.InitNeo4jDriver(ctx)
}

func TestGraphTest(t *testing.T) {
	suite.Run(t, new(GraphTest))
}
func (t *GraphTest) TestLineByName() {
	resMap, err := LineByUserName(context.Background(), "rrrr", "Alice")
	t.Nil(err)
	fmt.Println(len(resMap))
}
func (t *GraphTest) TestFindLinedFriend() {
	friend, err := FindLinedFriend(context.Background(), "rrrr")
	t.Nil(err)
	for _, f := range friend {
		fmt.Printf("username:%s\n,id:%d,gender:%s", f.Username, f.ID, f.Gender)
	}
}
