# go-chatroom
A nhooyr.io/websocket package backend

## 0、思路流程
1. 建立專案目錄並初始化 Go 模組
-   mkdir chatroom
-   cd chatroom
-   go mod init chatroom
-   mkdir chatroom 創建專案目錄
-   go mod init chatroom 初始化 Go 模組，讓 Go 可以管理依賴
2. 建立 cmd/chatroom/main.go
-   mkdir -p cmd/chatroom
-   touch cmd/chatroom/main.go
-   cmd/chatroom/main.go 是應用程式的進入點
-   建立 config 的設定及 global 的全域變數
-   這裡將負責啟動 HTTP 伺服器，載入 WebSocket 服務，並初始化其他必要的服務
3. 建立 server 資料夾與處理邏輯
-   mkdir server
-   touch server/home.go server/handle.go server/websocket.go
-   server/home.go 處理首頁請求
-   server/handle.go 處理 WebSocket 連線
-   server/websocket.go 管理 WebSocket 事件
4. 建立 logic 資料夾，負責業務邏輯
-   mkdir logic
-   touch logic/broadcast.go logic/message.go logic/user.go
-   logic/broadcast.go 管理廣播訊息
-   logic/message.go 處理訊息結構
-   logic/user.go 處理使用者狀態
5. 建立 template 目錄並加上 home.html
-   mkdir template
-   touch template/home.html
-   home.html 為聊天室前端頁面，包含 WebSocket 連線的 JavaScript


---

這樣你的專案就按照正確的順序建立好了！  
之所以這樣安排，是因為：
1. **先建立 `go mod`**：確保有 Go 模組管理
2. **建立 `cmd/chatroom/main.go`**：專案的進入點
3. **建立 `server` 處理 HTTP & WebSocket**
4. **建立 `logic` 處理聊天邏輯**
5. **建立 `template` 作為前端**
6. **最後補充 `README.md`**，方便記錄專案資訊

這樣的順序能確保每個部分都能正常組合，讓聊天室系統順利運行！


## 1、安装 wiresshark 抓包工具
- 官網下載安裝 : https://www.wireshark.org/
- 使用參考 : https://iter01.com/564644.html

## 2、安裝 nhoyr.io/websocket  (照常理安裝第一個即可)
- go get -u nhooyr.io/websocket
- go get -u nhooyr.io/websocket/wsjson
- go mod tidy

## 3、安裝 gorilla/websocket  (使用目前最受歡迎的套件進行server端重構)
- go get -u github.com/gorilla/websocket
- go mod tidy

## 4、安裝 viper (安裝完會一併安裝 fsnotify，如沒有再進行下一項安裝)
- viper 是一個 Go 語言的設定管理庫，支持多種格式的設定文件，如 JSON、TOML、YAML、HCL、envfile 和 Java properties config files。
- 安裝 viper：
- go get -u github.com/spf13/viper@v1.4.0
- go mod tidy

## 5、安裝 fsnotify (文件變更監控) 跨平台文件系統監聽事件庫
- fsnotify 是一個 Go 語言的文件變更監控庫，支持跨平台文件系統監聽事件。
- 安裝 fsnotify：
- go get -u golang.org/x/sys/...
- go get -u github.com/fsnotify/fsnotify
- go mod tidy

## 6、進行測試 (window要安裝 GCC（C 編譯器），mac 下不用 )
-  安裝有二種 : 方法 1：安裝 MinGW-w64（推薦），方法 2：使用 TDM-GCC（替代方案）
-  本次使用方法 2 : https://jmeubank.github.io/tdm-gcc/download/
-  安裝完在PowerShell下執行 : $env:Path += ";C:\TDM-GCC-64\bin"
-  主程式運行 :  go build -o go-chat.exe .\cmd\chatroom\main.go
-  在PowerShell 下 : $env:CGO_ENABLED="1"
-  壓力測試 : go run -race .\cmd\benchmark\main.go -u 100 -m 20s -l 0 
