package server

import (
	"net/http"
)

func websocketHandleFunc(w http.ResponseWriter, req *http.Request) {
	// Accept 從客戶端接收 WebSocket 握手，升級 HTTP 請求到 WebSocket 請求
	// 如果 Origin 域與主機不同，Accept 會拒絕請求，除非設置了 InsecureSkipVerify 選項(通過第三個參數 AcceptOption 進行設置)
	// 默認不允許跨域請求。如果發生錯誤，Accept 將始終寫入適當的響應
	//conn, err := websocket.Accept(w, req, &websocket.AcceptOptions{InsecureSkipVerify: true})
	//if err != nil {
	//	log.Println("websocket accept error:", err)
	//	return
	//}

	// 1. 創建用戶進來，構建用戶對象
	//token := req.FormValue("token")
	//nickname := req.FormValue("nickname")
	//if l := len(nickname); l < 2 || l > 20 {
	//	log.Println("nickname illegal:", nickname)
	//	wsjson.Write(req.Context(), conn, logic.NewErrorMessage("昵稱長度不合法，昵稱長度：2-20"))
	//	conn.Close(websocket.StatusUnsupportedData, "nickname illegal")
	//	return
	//}

	// 2. 啟動用戶寫入數據的 goroutine

	// 3. 給當前用戶發送歡迎消息

	// 4. 將該用戶加入到廣播器的用戶列表中

	// 5. 接收用戶消息

	// 6. 用戶離開
}
