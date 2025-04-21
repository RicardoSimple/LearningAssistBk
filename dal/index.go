package dal

import (
	"context"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"ar-app-api/conf"
	"ar-app-api/dal/schema"
	"ar-app-api/util/log"
)

const GraphDataBaseName = "neo4j"

var DB *gorm.DB
var GraphDB neo4j.DriverWithContext

func Init(ctx context.Context) {
	InitConnection(ctx)
	//InitNeo4jDriver(ctx)
}
func InitConnection(ctx context.Context) {
	cfg := conf.GetConfig().DB
	user := cfg.User
	password := cfg.Pass
	host := cfg.IP
	port := cfg.Port
	dbname := cfg.DbName

	// 构建 DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%v)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, port, dbname)

	log.Info("DSN: %v", dsn)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("连接数据库失败")
	}
	RefreshMigrate(ctx)
}

func RefreshMigrate(ctx context.Context) {
	log.Info("init migrate")
	err := DB.AutoMigrate(&schema.ImageHash{}, &schema.User{}, &schema.Message{}, &schema.Media{}, &schema.ChatGroup{}, &schema.Class{}, &schema.ClassTeacher{}, &schema.Assignment{}, &schema.CourseSubject{})
	if err != nil {
		panic("数据库表更新失败")
	}
}

//func InitNeo4jDriver(ctx context.Context) {
//	neo4jConfig := conf.GetConfig().Neo4j
//	var err error
//	GraphDB, err = neo4j.NewDriverWithContext(
//		neo4jConfig.URI,
//		neo4j.BasicAuth(neo4jConfig.User, neo4jConfig.Pass, ""))
//	err = GraphDB.VerifyConnectivity(ctx)
//	if err != nil {
//		log.Error("neo4j connectivity error: %v", err)
//	}
//}
