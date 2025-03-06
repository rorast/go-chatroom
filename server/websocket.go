package server

import (
	"github.com/rorast/go-chatroom/logic"
	"log"
	"net/http"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

func websocketHandleFunc(w http.ResponseWriter, req *http.Request) {
	// Accept 從客戶端接收 WebSocket 握手，升級 HTTP 請求到 WebSocket 請求
	// 如果 Origin 域與主機不同，Accept 會拒絕請求，除非設置了 InsecureSkipVerify 選項(通過第三個參數 AcceptOption 進行設置)
	// 默認不允許跨域請求。如果發生錯誤，Accept 將始終寫入適當的響應
	conn, err := websocket.Accept(w, req, &websocket.AcceptOptions{InsecureSkipVerify: true})
	if err != nil {
		log.Println("websocket accept error:", err)
		return
	}

	// 1. 創建用戶進來，構建用戶對象
	token := req.FormValue("token")
	nickname := req.FormValue("nickname")
	if l := len(nickname); l < 2 || l > 20 {
		log.Println("nickname illegal:", nickname)
		wsjson.Write(req.Context(), conn, logic.NewErrorMessage("昵稱長度不合法，昵稱長度：2-20"))
		conn.Close(websocket.StatusUnsupportedData, "nickname illegal")
		return
	}

	userHasToken := logic.NewUser(conn, token, nickname, req.RemoteAddr)

	// 2. 啟動用戶寫入數據的 goroutine
	go userHasToken.SendMessage(req.Context())

	// 3. 給當前用戶發送歡迎消息
	userHasToken.MessageChannel <- logic.NewWelcomeMessage(userHasToken)

	// 避免 token 泄露
	tmpUser := *userHasToken
	user := &tmpUser
	user.Token = ""

	// 發給所有用戶歡迎新用戶的進入
	msg := logic.NewUserEnterMessage(user)
	logic.Broadcaster.Broadcast(msg)

	// 4. 將該用戶加入到廣播器的用戶列表中
	logic.Broadcaster.UserEntering(user)
	log.Println("user:", nickname, "joins chat")

	// 5. 接收用戶消息
	err = user.ReceiveMessage(req.Context())

	// 6. 用戶離開
	logic.Broadcaster.UserLeaving(user)
	msg = logic.NewUserLeaveMessage(user)
	logic.Broadcaster.Broadcast(msg)
	log.Println("user:", nickname, "Leaves Chat")

	// 根據取到的錯誤執行同的 Close
	if err == nil {
		conn.Close(websocket.StatusNormalClosure, "")
	} else {
		log.Println("Read from client error:", err)
		conn.Close(websocket.StatusInternalError, "Read from client error")
	}
}
