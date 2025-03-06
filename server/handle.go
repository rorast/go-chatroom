package server

import (
	"github.com/rorast/go-chatroom/logic"
	"net/http"
)

func RegisterHandle() {
	// 廣播消息處理
	go logic.Broadcaster.Start()

	// 聊天室服務器處理路由
	http.HandleFunc("/", indexHandleFunc)
	http.HandleFunc("/users", userHandleFunc)
	http.HandleFunc("/ws", websocketHandleFunc)
}
