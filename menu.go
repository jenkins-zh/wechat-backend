package main

import (
	"bytes"
	"log"
	"net/http"
	"strings"

	"github.com/jenkins-zh/wechat-backend/pkg/config"
)

func PushWxMenuCreate(accessToken string, menuJsonBytes []byte) error {
	postReq, err := http.NewRequest(http.MethodPost,
		strings.Join([]string{"https://api.weixin.qq.com/cgi-bin/menu/create", "?access_token=", accessToken}, ""),
		bytes.NewReader(menuJsonBytes))

	if err != nil {
		log.Println("向微信发送菜单建立请求失败", err)
		return err
	}

	postReq.Header.Set("Content-Type", "application/json; encoding=utf-8")

	client := &http.Client{}
	resp, err := client.Do(postReq)
	if err != nil {
		log.Println("client向微信发送菜单建立请求失败", err)
		return err
	}

	defer resp.Body.Close()
	log.Println("向微信发送菜单建立成功")

	return nil
}

func createWxMenu(config *config.WeChatConfig) {

	//btn1 := models.Btn{Name: "进入商城", Url: "http://www.baidu.com/", Btype: "view"}
	//btn2 := models.Btn{Name: "会员中心", Key: "molan_user_center", Btype: "click"}
	//btn3 := models.Btn{Name: "我的", Url: "http://www.baidu.com/user_view", Btype: "view"}
	//
	//btns := []models.Btn{btn1, btn2, btn3}
	//wxMenu := models.WxMenu{Button: btns}
	//menuJsonBytes, err := json.Marshal(wxMenu)

	menuStr := `{
            "button": [
            {
                "name": "进入商城",
                "type": "view",
                "url": "http://www.baidu.com/"
            },
            {

                "name":"管理中心",
                 "sub_button":[
                        {
                        "name": "用户中心",
                        "type": "click",
                        "key": "molan_user_center"
                        },
                        {
                        "name": "公告",
                        "type": "click",
                        "key": "molan_institution"
                        }]
            },
            {
                "name": "资料修改",
                "type": "view",
                "url": "http://www.baidu.com/user_view"
            }
            ]
        }`

	//if err == nil {

	//fmt.Println("生成的菜单json--->", menuStr)

	//发送建立菜单的post请求
	token := getAccessToken(config)
	PushWxMenuCreate(token, []byte(menuStr))
	//} else {
	//  logUtils.GetLog().Error("微信菜单json转换错误", err)
	//}
}
