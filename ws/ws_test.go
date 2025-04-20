package ws

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestWs(t *testing.T) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/ws/ping"}
	log.Printf("connecting to %s", u.String())

	token := "AjyghT+SS2uetlm8Zfxplhn/PY4jLMoYfQy6oCgmn6TLldu1ZEIW/j2MgYDEX5JZp2Igcd4Ip1OMfgTavFN0KDoJ+bXZ7+rKmVlqiTsoqAiGzcuyTvL7NGnMCH9n5GMNP1gGUxT4PaRw5Rx5ptNCF+btt3cLsBWESNmp4Q0lF7yKOlIyHac0mOvCnSZNAxgrdRDPVeXIq84J4qO4donVlUcezXbsx8wEpxiQn4gLi3nUzJ5CL+eU6a0ZUPdWqYA="
	// 使用子协议 "chat"
	c, _, err := websocket.DefaultDialer.Dial(u.String(), http.Header{"Sec-WebSocket-Protocol": {token}})
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
			if err != nil {
				log.Println("write:", err)
				return
			}
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
