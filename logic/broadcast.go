package logic

import (
	"expvar" // Go 內建的變數監控工具，用來監控 message_queue 長度的變數數據
	"fmt"
	"github.com/rorast/go-chatroom/global"
	"log"
)

func init() {
	// 初始化 Expvar 變數監控
	expvar.Publish("message_queue", expvar.Func(calcMessageQueueLen))
}

// 計算訊息佇列的長度
func calcMessageQueueLen() interface{} {
	fmt.Println("++++length++++:", len(Broadcaster.messageChannel))
	return len(Broadcaster.messageChannel)
}

// Broadcaster 廣播器結構體
type broadcaster struct {
	// 所有聊天室用户(透過 make(map[string]*User) 初始化 users，用來管理聊天室用戶。)
	users map[string]*User

	// 所有 channel 統一管理，可以避免外部亂用
	// enteringChannel, leavingChannel, messageChannel 等 chan 變數也都透過 make(chan ...) 初始化，確保可以正常使用。
	enteringChannel chan *User    // 使用者進入聊天室
	leavingChannel  chan *User    // 使用者離開聊天室
	messageChannel  chan *Message // 訊息佇列

	// 判斷該昵稱用戶是否可進入聊天室（重復與否）：true 能，false 不能
	checkUserChannel      chan string
	checkUserCanInChannel chan bool

	// 獲取用戶列表
	requestUsersChannel chan struct{} // 當外部請求用戶列表時，這個通道會收到一個空結構體 struct{} 來觸發查詢。
	usersChannel        chan []*User  // 用來回傳當前所有在線的 User。
}

// Broadcaster 變數：初始化 broadcaster - 單例模式(這裡定義了一個全域變數 Broadcaster，以確保聊天室的 broadcaster 只有一個實例。)
var Broadcaster = &broadcaster{
	users: make(map[string]*User),

	enteringChannel: make(chan *User),
	leavingChannel:  make(chan *User),
	// messageChannel 是唯一有 buffer（global.MessageQueueLen）的 channel，確保訊息佇列不會阻塞。
	messageChannel: make(chan *Message, global.MessageQueueLen), // messageChannel 的容量由 global.MessageQueueLen 控制，避免訊息堆積過多導致崩潰。

	checkUserChannel:      make(chan string),
	checkUserCanInChannel: make(chan bool),

	requestUsersChannel: make(chan struct{}),
	usersChannel:        make(chan []*User),
}

// Start() - 廣播器的核心 - 需要在一个新 goroutine 中運行，因为它不會返回
func (b *broadcaster) Start() {
	fmt.Println("Start() ...  ")
	// 事件驅動的 Goroutine，負責處理不同的聊天室事件。
	for { // 這裡是一個無限循環，不斷地從不同的 channel 中讀取數據。
		select {
		// 新使用者進入聊天室，存入 users，並可能發送離線訊息。
		case user := <-b.enteringChannel:
			fmt.Println("user := <-b.enteringChannel: ...  ")
			// 新用户进入
			b.users[user.NickName] = user

			OfflineProcessor.Send(user)
		// 使用者離開聊天室，從 users 刪除並關閉訊息通道。
		case user := <-b.leavingChannel:
			// 用户离开
			delete(b.users, user.NickName)
			// 避免 goroutine 泄露
			user.CloseMessageChannel()
		// 對所有使用者廣播訊息，但排除發送者自己。
		case msg := <-b.messageChannel:
			// 给所有在线用户发送消息
			for _, user := range b.users {
				if user.UID == msg.User.UID {
					continue
				}
				user.MessageChannel <- msg
				fmt.Println("msg :: ", msg)
			}
			OfflineProcessor.Save(msg)
		// 檢查用戶是否已存在，結果透過 checkUserCanInChannel 回傳。
		case nickname := <-b.checkUserChannel:
			if _, ok := b.users[nickname]; ok {
				b.checkUserCanInChannel <- false
			} else {
				b.checkUserCanInChannel <- true
			}
		// 查詢當前在線使用者並透過 usersChannel 回傳。
		case <-b.requestUsersChannel:
			userList := make([]*User, 0, len(b.users))
			for _, user := range b.users {
				userList = append(userList, user)
			}

			b.usersChannel <- userList
		}
	}
}

/*
UserEntering() 和 UserLeaving() 負責把 User 寫入對應的 channel 來驅動 Start() 內的事件。
*/
// 使用者進入
func (b *broadcaster) UserEntering(u *User) {
	b.enteringChannel <- u
}

// 使用者離開
func (b *broadcaster) UserLeaving(u *User) {
	b.leavingChannel <- u
}

// 訊息廣播
func (b *broadcaster) Broadcast(msg *Message) {
	// 如果 messageChannel 滿了，則記錄錯誤日誌。
	if len(b.messageChannel) >= global.MessageQueueLen {
		log.Println("broadcast queue 滿了")
	}
	b.messageChannel <- msg
}

// 查詢使用者是否能進入 (傳送 nickname，接收 bool 來決定是否能進入。)
func (b *broadcaster) CanEnterRoom(nickname string) bool {
	b.checkUserChannel <- nickname
	return <-b.checkUserCanInChannel
}

// 獲取目前在線使用者
func (b *broadcaster) GetUserList() []*User {
	b.requestUsersChannel <- struct{}{}
	return <-b.usersChannel
}
