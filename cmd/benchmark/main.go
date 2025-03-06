package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/rorast/go-chatroom/logic"
	"log"
	_ "net/http/pprof"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
	"strconv"
	"time"
)

var (
	userNum       int           // 用户數
	loginInterval time.Duration // 用户登入時間間隔
	msgInterval   time.Duration // 單個用戶發送訊息的時間間隔
)

func init() {
	// 透過 -u 參數設定用戶數，預設值為 500。
	flag.IntVar(&userNum, "u", 500, "登錄用戶數")
	// 5e9 表示 5 秒 (5 * 10^9 奈秒)。
	flag.DurationVar(&loginInterval, "l", 5e9, "用戶陸續登錄時間間隔")
	// 預設 1 分鐘發送一則訊息。
	flag.DurationVar(&msgInterval, "m", 1*time.Second, "用戶發送消息時間間隔")
}

func main() {
	// flag.Parse() 解析命令列參數，允許 go run main.go -u 1000 -l 1s -m 500ms 這類輸入。
	flag.Parse()

	for i := 0; i < userNum; i++ {
		go UserConnect("user" + strconv.Itoa(i))
		time.Sleep(loginInterval)
	}

	// 防止 main 函式退出，select {} 會 永久阻塞，讓 Goroutines 繼續執行。
	select {}
}

func UserConnect(nickname string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	conn, _, err := websocket.Dial(ctx, "ws://127.0.0.1:2066/ws?nickname="+nickname, nil)
	if err != nil {
		log.Println("Dial error:", err)
		return
	}
	defer conn.Close(websocket.StatusInternalError, "內部錯誤！")

	go sendMessage(conn, nickname)

	ctx = context.Background()

	for {
		var message logic.Message
		err = wsjson.Read(ctx, conn, &message)
		if err != nil {
			log.Println("receive msg error:", err)
			continue
		}

		if message.ClientSendTime.IsZero() {
			continue
		}
		if d := time.Now().Sub(message.ClientSendTime); d > 1*time.Second {
			fmt.Printf("接收到服務器響應(%d)：%#v\n", d.Milliseconds(), message)
		}
	}

	conn.Close(websocket.StatusNormalClosure, "")
}

func sendMessage(conn *websocket.Conn, nickname string) {
	ctx := context.Background()
	i := 1
	for {
		msg := map[string]string{
			"content":   "來自" + nickname + "的消息:" + strconv.Itoa(i),
			"send_time": strconv.FormatInt(time.Now().UnixNano(), 10),
		}
		err := wsjson.Write(ctx, conn, msg)
		if err != nil {
			log.Println("send msg error:", err, "nickname:", nickname, "no:", i)
		}
		i++

		time.Sleep(msgInterval)
	}
}
