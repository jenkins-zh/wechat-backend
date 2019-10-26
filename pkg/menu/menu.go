package menu

import (
	"bytes"
	"log"
	"net/http"
	"strings"

<<<<<<< HEAD:pkg/menu/menu.go
	"github.com/linuxsuren/wechat-backend/pkg/config"
	"github.com/linuxsuren/wechat-backend/pkg/token"
=======
	"github.com/jenkins-zh/wechat-backend/pkg/api"
	"github.com/jenkins-zh/wechat-backend/pkg/config"
>>>>>>> master:menu.go
)

func pushWxMenuCreate(accessToken string, menuJsonBytes []byte) error {
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

// CreateWxMenu create wechat menu
func CreateWxMenu(config *config.WeChatConfig) {
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

	//发送建立菜单的post请求
<<<<<<< HEAD:pkg/menu/menu.go
	token := token.GetAccessToken(config)
	pushWxMenuCreate(token, []byte(menuStr))
=======
	token := api.GetAccessToken(config)
	PushWxMenuCreate(token, []byte(menuStr))
	//} else {
	//  logUtils.GetLog().Error("微信菜单json转换错误", err)
	//}
>>>>>>> master:menu.go
}
