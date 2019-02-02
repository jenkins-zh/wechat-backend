package service

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/linuxsuren/wechat-backend/pkg/config"
)

// HandleConfig handle the config modify
func HandleConfig(w http.ResponseWriter, r *http.Request, weConfig *config.WeChatConfig) {
	r.ParseForm()

	validStr := strings.Join(r.Form["valid"], "")
	weConfig.Valid = (validStr == "true")

	w.Write([]byte(fmt.Sprintf("WeChat valid: %v", weConfig.Valid)))

	weConfig.SettingConfig()
}
