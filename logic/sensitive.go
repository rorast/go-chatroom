package logic

import (
	"strings"

	"github.com/rorast/go-chatroom/global"
)

// global.SensitiveWords：來自 global 套件的全局變數，是 []string 類型，儲存所有的敏感詞。
func FilterSensitive(content string) string {
	// range 遍歷 slice (切片) ， _ 代表省略索引，只關心值 word。
	for _, word := range global.SensitiveWords {
		content = strings.ReplaceAll(content, word, "**")
	}

	return content
}
