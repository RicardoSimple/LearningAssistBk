package conf

import (
	"context"
)

var cfg Config

type Config struct {
	Jwt JwtConfig `json:"jwt"`
	DB  DBConfig  `json:"db"`
	W   WechatConfig
	Git GitHubConfig
	//Neo4j Neo4jConfig
}

type LarkConfig struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectURL  string `json:"redirect_url"`
}

type JwtConfig struct {
	Secret        string `json:"secret"`
	Expire        int    `json:"expire"` // 超时时间(小时)
	RefreshSecret string `json:"refresh_secret"`
	RefreshExpire int    `json:"refresh_expire"` // refresh_token 超时时间(小时)
	AesKey        string `json:"aes_key"`
}

type GitHubConfig struct {
	UserName string `json:"user_name"`
	Repo     string `json:"repo"`
	Token    string `json:"token"`
	Message  string `json:"message"`
}
type DBConfig struct {
	IP     string `json:"ip"`
	Port   int    `json:"port"`
	User   string `json:"user"`
	Pass   string `json:"pass"`
	DbName string `json:"db_name"`
}

type Neo4jConfig struct {
	URI  string `json:"uri"`
	User string `json:"user"`
	Pass string `json:"pass"`
}

type WechatConfig struct {
	AppId          string `json:"appid"`
	AppSecret      string `json:"secret"`
	EncodingAESKey string `json:"encodingAESKey"`
	Token          string `json:"token"`
}

func Init(ctx context.Context) {
	cfg = Config{
		Jwt: JwtConfig{Secret: "*#3507", Expire: 24, RefreshSecret: "refresh_*#3507", RefreshExpire: 720, AesKey: "QWERTYUIOPASDFGH"},
		DB:  DBConfig{IP: "111.229.121.171", Port: 3306, User: "learningAssist", Pass: "zAEkGG6L8z3EpawG", DbName: "learningassist"},
		//Git: GitHubConfig{
		//	UserName: "RicardoSimple",
		//	Repo:     "pic-store",
		//	Token:    "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",
		//	Message:  "add image",
		//},
		//Neo4j: Neo4jConfig{
		//	URI:  "neo4j://1.117.229.251:7687",
		//	User: "neo4j",
		//	Pass: "everest-nikita-cabinet-welcome-galaxy-4397",
		//},
	}
}

func GetConfig() Config {
	return cfg
}
