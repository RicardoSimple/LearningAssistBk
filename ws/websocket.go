package ws

import (
	"net/http"
	"strconv"
	"time"

	"learning-assistant/util/log"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"learning-assistant/handler/aerrors"
	"learning-assistant/handler/basic"
	"learning-assistant/util"
)

var upgrader = websocket.Upgrader{
	// 这个是校验请求来源
	// 在这里我们不做校验，直接return true
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func PingWS(context *gin.Context) {
	// 将普通的http GET请求升级为websocket请求
	client, err := upgrader.Upgrade(context.Writer, context.Request, nil)
	if err != nil {
		log.Info("[WEBSOCKET PING]Failed to upgrade to WebSocket:", err)
		return
	}
	go func() {
		for {
			// 每隔两秒给前端推送一句消息“hello, WebSocket”
			err := client.WriteMessage(websocket.TextMessage, []byte("hello, WebSocket"))
			if err != nil {
				log.Error("[WEBSOCKET PING] err:%s", err.Error())
			}
			time.Sleep(time.Second * 2)
		}
	}()
	go func() {
		for {
			// 读取客户端发送的消息
			_, message, err := client.ReadMessage()
			if err != nil {
				log.Info("[WEBSOCKET PING ]Failed to read message:%s", err.Error())
				break
			}
			log.Info("[WEBSOCKET PING]Received message from client: %s\n", message)

			// 在这里可以对接收到的消息进行处理，比如广播给其他客户端

			// 示例：将接收到的消息原样返回给客户端
			err = client.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Info("[WEBSOCKET PING]Failed to send message to client:", err)
				break
			}
		}
	}()
}

// api:/getPushNews接口处理函数
func GetPushNews(context *gin.Context) {
	userInfo, err := util.GetUserFromGinContext(context)
	if err != nil {
		basic.AuthFailure(context)
		return
	}
	id := strconv.Itoa(int(userInfo.ID))
	log.Info("[WEBSOCKET] 建立连接:id: %s", id)
	// 升级为websocket长链接
	WsHandler(context.Writer, context.Request, id)
}

// api:/deleteClient接口处理函数
func DeleteClient(context *gin.Context) {
	// 关闭websocket链接
	userInfo, err := util.GetUserFromGinContext(context)
	if err != nil {
		basic.AuthFailure(context)
		return
	}
	id := strconv.Itoa(int(userInfo.ID))
	conn, exist := getClient(id)
	if exist {
		conn.Close()
		deleteClient(id)
	} else {
		basic.RequestFailureWithCode(context, "无客户端", aerrors.RecordNotFind)
	}
	// 关闭其消息通道
	_, exist = getNewsChannel(id)
	if exist {
		deleteNewsChannel(id)
	}
}
