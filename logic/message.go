package logic

import (
	"github.com/spf13/cast"
	"time"
)

/*
📌 知識點：iota
iota 是 Go 的常數計數器，從 0 開始，每定義一個新常數，值會自動 +1。
這裡 MsgTypeNormal = 0，MsgTypeWelcome = 1，以此類推。
這種方式讓程式碼更 簡潔 且 易於擴展，如果未來要新增訊息類型，只需在 const 區塊內加一行即可。
*/
const (
	MsgTypeNormal    = iota // 普通 用戶訊息
	MsgTypeWelcome          // 當前用户歡迎訊息
	MsgTypeUserEnter        // 用戶進入聊天室
	MsgTypeUserLeave        // 用戶離開聊天室
	MsgTypeError            // 錯誤消息
)

// 給用戶發送的消息
type Message struct {
	// 哪個用戶發送的消息
	User    *User     `json:"user"`
	Type    int       `json:"type"`
	Content string    `json:"content"`
	MsgTime time.Time `json:"msg_time"`

	ClientSendTime time.Time `json:"client_send_time"`

	// 消息 @ 了誰
	Ats []string `json:"ats"`

	// 用戶列表不通過 WebSocket 下發
	// Users []*User `json:"users"`
}

// NewMessage 創建消息
func NewMessage(user *User, content, clientTime string) *Message {
	message := &Message{
		User:    user,
		Type:    MsgTypeNormal,
		Content: content,
		MsgTime: time.Now(),
	}
	if clientTime != "" {
		message.ClientSendTime = time.Unix(0, cast.ToInt64(clientTime))
	}
	return message
}

func NewWelcomeMessage(user *User) *Message {
	return &Message{
		User:    user,
		Type:    MsgTypeWelcome,
		Content: user.NickName + " 您好，歡迎加入聊天室！",
		MsgTime: time.Now(),
	}
}

func NewUserEnterMessage(user *User) *Message {
	return &Message{
		User:    user,
		Type:    MsgTypeUserEnter,
		Content: user.NickName + " 加入了聊天室",
		MsgTime: time.Now(),
	}
}

func NewUserLeaveMessage(user *User) *Message {
	return &Message{
		User:    user,
		Type:    MsgTypeUserLeave,
		Content: user.NickName + " 離開了聊天室",
		MsgTime: time.Now(),
	}
}

func NewErrorMessage(content string) *Message {
	return &Message{
		User:    System,
		Type:    MsgTypeError,
		Content: content,
		MsgTime: time.Now(),
	}
}
