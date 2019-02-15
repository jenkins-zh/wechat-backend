package service

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/linuxsuren/wechat-backend/pkg/config"
)

// HandleConfig handle the config modify
func HandleConfig(w http.ResponseWriter, r *http.Request, weConfig config.WeChatConfigurator) {
	r.ParseForm()

	validStr := strings.Join(r.Form["valid"], "")

	config := weConfig.GetConfig()
	config.Valid = (validStr == "true")

	w.Write([]byte(fmt.Sprintf("WeChat valid: %v", config.Valid)))

	weConfig.SaveConfig()
}
