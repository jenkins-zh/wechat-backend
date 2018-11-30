package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

//https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=APPID&secret=APPSECRET

type AccessToken struct {
	AccessToken string `json:"access_token"`
	Expires     int    `json:"expires_in"`
}

func getAccessToken() string {
	resp, err := http.Get(strings.Join([]string{"https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential", "&appid=", "wxdd2924d2e2598bff", "&secret=", "8db8b2ee78059e7f4cb7e38020e2f566"}, ""))
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
