package ws

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"ar-app-api/util/log"

	"github.com/gorilla/websocket"
)

var (
	// 消息通道
	news = make(map[string]chan interface{})
	// websocket客户端链接池
	client = make(map[string]*websocket.Conn)
	// 互斥锁，防止程序对统一资源同时进行读写
	mux sync.Mutex
)

// 在 ws 包中添加 SendMessageToUser 方法
func SendMessageToUser(userID string, message interface{}) error {
	// 获取指定用户的消息通道
	m, exist := getNewsChannel(userID)
	if !exist {
		return fmt.Errorf("[WEBSOCKET] user not found or disconnected")
	}
	// 向用户的消息通道发送消息
	select {
	case m <- message:
		return nil
	default:
		// 如果发送失败，可能是用户已经断开连接或消息通道已关闭
		return fmt.Errorf("[WEBSOCEKT] failed to send message to user")
	}
}

// websocket Upgrader
var wsupgrader = websocket.Upgrader{
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
	HandshakeTimeout: 5 * time.Second,
	// 取消ws跨域校验
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// WsHandler 处理ws请求
func WsHandler(w http.ResponseWriter, r *http.Request, id string) {
	var conn *websocket.Conn
	var err error
	var exist bool
	// 创建一个定时器用于服务端心跳
	pingTicker := time.NewTicker(time.Second * 10)
	conn, err = wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error(err.Error())
		return
	}
	// 把与客户端的链接添加到客户端链接池中
	addClient(id, conn)

	// 获取该客户端的消息通道
	m, exist := getNewsChannel(id)
	if !exist {
		m = make(chan interface{})
		addNewsChannel(id, m)
	}

	// 设置客户端关闭ws链接回调函数
	conn.SetCloseHandler(func(code int, text string) error {
		deleteClient(id)
		log.Info("[WEBSOCKET] DELETE CLIENT:%d", code)
		return nil
	})
	go func() {
		// 发送消息+心跳
		for {
			select {
			case content, _ := <-m:
				// 从消息通道接收消息，然后推送给前端
				err = conn.WriteJSON(content)
				if err != nil {
					log.Error(err.Error())
					conn.Close()
					deleteClient(id)
					return
				}
			case <-pingTicker.C:
				// 服务端心跳:每20秒ping一次客户端，查看其是否在线
				conn.SetWriteDeadline(time.Now().Add(time.Second * 20))
				err = conn.WriteMessage(websocket.PingMessage, []byte{})
				if err != nil {
					log.Error("[WEBSOCKET] send ping err:%s", err.Error())
					conn.Close()
					deleteClient(id)
					return
				}
			}
		}
	}()
	go func() {
		// 接收消息
		for {
			// 读取客户端发送的消息
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Info("[WEBSOCKET] Failed to read message: %v", err)
				break
			}
			log.Info("[WEBSOCKET] Received message from client %s: %s", id, message)

			// 在这里可以对接收到的消息进行处理，比如广播给其他客户端
			go messageCenter.MessageControl(context.Background(), string(message))
		}
	}()
}

// 将客户端添加到客户端链接池
func addClient(id string, conn *websocket.Conn) {
	mux.Lock()
	client[id] = conn
	mux.Unlock()
}

// 获取指定客户端链接
func getClient(id string) (conn *websocket.Conn, exist bool) {
	mux.Lock()
	conn, exist = client[id]
	mux.Unlock()
	return
}

// 删除客户端链接
func deleteClient(id string) {
	mux.Lock()
	delete(client, id)
	log.Info("[WEBSOCKET] client 关闭:%s", id)
	mux.Unlock()
}

// 添加用户消息通道
func addNewsChannel(id string, m chan interface{}) {
	mux.Lock()
	news[id] = m
	mux.Unlock()
}

// 获取指定用户消息通道
func getNewsChannel(id string) (m chan interface{}, exist bool) {
	mux.Lock()
	m, exist = news[id]
	mux.Unlock()
	return
}

// 删除指定消息通道
func deleteNewsChannel(id string) {
	mux.Lock()
	if m, ok := news[id]; ok {
		close(m)
		delete(news, id)
	}
	mux.Unlock()
}
