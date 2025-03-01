package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
	"sync"
	"time"
)

var (
	// 新用戶到來，通過該 channel 進行登記
	enteringChannel = make(chan *User)
	// 用戶離開，通過該 channel 進行登記
	leavingChannel = make(chan *User)
	// 广播用戶發送的消息 channel ，緩沖是盡可能避免出現異常情況阻塞，這裡設置為 8，可以根據實際情況進行調整
	messageChannel = make(chan Message, 8)
)

type User struct {
	ID             int
	Addr           string
	EnterAt        time.Time
	MessageChannel chan string
}

func (u *User) String() string {
	return u.Addr + ", UID:" + strconv.Itoa(u.ID) + ", EnterAt:" + u.EnterAt.Format("2025-03-01 15:04:05")
}

// 給用戶發送的消息
type Message struct {
	OwnerID int
	Content string
}

func main() {
	listener, err := net.Listen("tcp", ":2020")
	if err != nil {
		panic(err)
	}

	go broadcaster()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go handleConn(conn)
	}
}

// broadcaster 用于記錄聊天室用戶，並進行消息廣播:
// 1。新用戶加入: 2。用戶普通消息: 3。用戶退出
func broadcaster() {
	users := make(map[*User]struct{})

	for {
		select {
		case user := <-enteringChannel:
			// 新用戶加入
			users[user] = struct{}{}
		case msg := <-messageChannel:
			// 向所有用戶發送消息
			for user := range users {
				//user.MessageChannel <- msg
				if user.ID == msg.OwnerID {
					continue
				}
				user.MessageChannel <- msg.Content
			}
		case user := <-leavingChannel:
			// 用戶退出
			delete(users, user)
			// 關閉用戶的消息通道 避免 goroutine 泄漏
			close(user.MessageChannel)
		}
	}
}

// handleConn 用於處理用戶連接:
func handleConn(conn net.Conn) {
	defer conn.Close()

	// 1. 新用戶加入, 構建該用戶的實例
	user := &User{
		ID:             GenUserID(),
		Addr:           conn.RemoteAddr().String(),
		EnterAt:        time.Now(),
		MessageChannel: make(chan string, 8),
	}

	// 2. 由於當前是在一個新的 goroutine 中進行讀寫操作的，所以這裡需要進行讀寫 goroutine 的啟動
	// 讀寫 goroutine 用於處理用戶發送的消息和向用戶發送消息 通過 channel 進行通信
	go sendMessage(conn, user.MessageChannel)

	// 3. 向當前用戶發送消息: 歡迎用戶加入
	user.MessageChannel <- "Welcome, " + user.String()
	msg := Message{
		OwnerID: user.ID,
		Content: "user:`" + strconv.Itoa(user.ID) + "` has enter",
	}
	messageChannel <- msg

	// 4. 將當前用戶加入到全局用戶列表中，這裡是通過 channel 進行通知的, 避免用鎖
	enteringChannel <- user
	// 踢出超時用戶
	var userActive = make(chan struct{})
	go func() {
		d := 1 * time.Minute
		timer := time.NewTimer(d)
		for {
			select {
			case <-timer.C:
				conn.Close()
			case <-userActive:
				timer.Reset(d)
			}
		}
	}()

	// 5. 循環讀取用戶的輸入
	input := bufio.NewScanner(conn)
	for input.Scan() {
		//messageChannel <- strconv.Itoa(user.ID) + ":" + input.Text()
		msg.Content = strconv.Itoa(user.ID) + ":" + input.Text()
		messageChannel <- msg

		// 用戶活躍
		userActive <- struct{}{}
	}

	if err := input.Err(); err != nil {
		log.Println("讀取錯誤:", err)
	}

	// 6. 用戶退出
	leavingChannel <- user
	msg.Content = "user:`" + strconv.Itoa(user.ID) + "` has left"
	messageChannel <- msg
}

func sendMessage(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		//conn.Write([]byte(msg + "\n"))
		fmt.Fprintln(conn, msg)
	}
}

// 生成用戶 ID
var (
	globalID int
	idLocker sync.Mutex
)

func GenUserID() int {
	idLocker.Lock()
	defer idLocker.Unlock()

	globalID++
	return globalID
}
