package logic

/*
這段程式碼實現了一個離線消息處理系統，其核心邏輯如下：

1、recentRing（全局環形緩存）：存儲最近 n 條消息，所有用戶重新上線時都能收到這些消息。
2、userRing（用戶專屬環形緩存）：用來存放某些特定用戶（被 @）的離線消息，讓被標記的用戶上線後可以看到這些對話。
3、消息儲存 (Save)：
   - 普通消息存入 recentRing。
   - 如果消息包含 @，則存入 userRing。
4、消息發送 (Send)：
   - 用戶上線時，先發送 recentRing 的歷史消息。
   - 如果該用戶曾被 @，則發送 userRing 的專屬離線消息，然後刪除記錄。
這樣的設計能夠確保：

   - 用戶不會錯過聊天室的最新對話。
   - @某人的消息不會被忽略，讓用戶重新上線時能補足未讀消息。
*/

import (
	"container/ring"
	"github.com/spf13/viper"
) // 這是一個標準庫，提供**環形緩存（Ring Buffer）**結構，可用於儲存固定數量的最近消息，當超過容量時會自動覆蓋最舊的數據。

type offlineProcessor struct {
	n int // 這是一個標準庫，提供**環形緩存（Ring Buffer）**結構，可用於儲存固定數量的最近消息，當超過容量時會自動覆蓋最舊的數據。

	// 這是一個環形緩衝區（ring.Ring），用於存放所有用戶最近的 n 條消息。
	recentRing *ring.Ring

	// 這是一個映射（map），key 為用戶名稱（string），value 是 ring.Ring，用來存放該用戶的個人離線消息（最多 n 條）。
	userRing map[string]*ring.Ring
}

var OfflineProcessor = newOfflineProcessor()

func newOfflineProcessor() *offlineProcessor {
	n := viper.GetInt("offline-num") // 從設定檔中讀取 offline-num，確定環形緩存的大小（n）。

	return &offlineProcessor{
		n:          n,
		recentRing: ring.New(n),                 // 建立一個大小為 n 的 recentRing，用來存最近的 n 條全局消息。
		userRing:   make(map[string]*ring.Ring), // 初始化 userRing，用於存放特定用戶的個人消息。
	}
}

// 儲存離線消息
func (o *offlineProcessor) Save(msg *Message) {
	// 負責存儲新的聊天消息，但只儲存普通類型（MsgTypeNormal）的消息，其他類型的消息會被忽略。
	if msg.Type != MsgTypeNormal {
		return
	}
	o.recentRing.Value = msg           // 將 msg 儲存在目前的 ring 節點中。
	o.recentRing = o.recentRing.Next() // 移動到下一個節點，這樣當緩存滿時，最舊的數據會被覆蓋。

	// 這段程式碼處理「@提及某個用戶」的情況，msg.Ats 代表消息中所有被 @ 提及的用戶列表。
	for _, nickname := range msg.Ats {
		nickname = nickname[1:] // 去掉 @ 符号，得到真正的用户名
		var (
			r  *ring.Ring
			ok bool
		)
		if r, ok = o.userRing[nickname]; !ok {
			r = ring.New(o.n)
		}
		r.Value = msg
		o.userRing[nickname] = r.Next()
	}
}

// 發送離線消息
func (o *offlineProcessor) Send(user *User) {
	// 這個方法在用戶重新連接聊天室時執行，它會發送該用戶應該接收到的離線消息。
	// 這段程式碼會遍歷 recentRing 中的所有消息，然後逐條發送到 user.MessageChannel，讓用戶收到這些歷史消息。
	o.recentRing.Do(func(value interface{}) {
		if value != nil {
			user.MessageChannel <- value.(*Message)
		}
	})

	// 如果用戶是新加入的 (isNew == true)，則不需要發送私人歷史消息，直接返回。
	if user.isNew {
		return
	}

	// 如果該用戶曾被 @ 過，則從 userRing 取出所有消息，發送到 user.MessageChannel。
	// 發送完後，刪除該用戶的記錄，避免重複發送。
	if r, ok := o.userRing[user.NickName]; ok {
		r.Do(func(value interface{}) {
			if value != nil {
				user.MessageChannel <- value.(*Message)
			}
		})

		delete(o.userRing, user.NickName)
	}
}
