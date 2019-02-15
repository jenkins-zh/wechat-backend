package token

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/linuxsuren/wechat-backend/pkg/config"
)

//https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=APPID&secret=APPSECRET

type AccessToken struct {
	AccessToken string `json:"access_token"`
	Expires     int    `json:"expires_in"`
}

func GetAccessToken(config *config.WeChatConfig) string {
	resp, err := http.Get(strings.Join([]string{"https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential",
		"&appid=", config.AppID, "&secret=", config.AppSecret}, ""))
	if err != nil {
		return ""
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("search query failed: %s\n", resp.Status)
		return ""
	}

	var result AccessToken
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return ""
	}

	return result.AccessToken
}
