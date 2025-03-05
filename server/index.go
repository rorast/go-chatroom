package server

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/rorast/go-chatroom/global"
)

func indexHandleFunc(w http.ResponseWriter, req *http.Request) {
	tpl, err := template.ParseFiles(global.RootDir + "/template/index.html")
	if err != nil {
		fmt.Fprint(w, "模板解析錯誤！")
		return
	}

	err = tpl.Execute(w, nil)
	if err != nil {
		fmt.Fprint(w, "模板執行錯誤！")
		return
	}
}

func userHandleFunc(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	//userList := logic.Broadcaster.GetUserList()
	//b, err := json.Marshal(userList)
	//
	//if err != nil {
	//	fmt.Fprint(w, `[]`)
	//} else {
	//	fmt.Fprint(w, string(b))
	//}
}
