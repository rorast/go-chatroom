# go-chatroom
A nhooyr.io/websocket package backend

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