package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/jenkins-zh/wechat-backend/pkg/config"
)

func ListMedias(w http.ResponseWriter, r *http.Request, cfg config.WeChatConfigurator) {
	r.ParseForm()
	// imageName := r.Form["imageName"]

	if weConfig, err := cfg.LoadConfig("config/wechat.yaml"); err == nil {
		token := GetAccessToken(weConfig)

		url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/material/batchget_material?access_token=%s", token)

		request := ListRequest{
			Type:   "image",
			Offset: 0,
			Count:  20,
		}

		var data []byte
		var err error
		if data, err = json.Marshal(&request); err != nil {
			data = []byte("")
		}

		log.Println("rquest ", string(data))

		postReq, err := http.NewRequest(http.MethodPost,
			url,
			bytes.NewReader(data))
		client := &http.Client{}
		resp, err := client.Do(postReq)
		if err != nil {
			log.Println("rquest media failure", err)
			return
		} else {
			log.Println("media request success")
		}

		data, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("read media request body error", err)
			return
		} else {
			log.Println("read body success, data: ", string(data))
		}

		itemList := MediaItemList{}
		if err = json.Unmarshal(data, &itemList); err == nil {
			for _, item := range itemList.ItemList {
				fmt.Println("name: ", item.Name, ", id: ", item.MediaID, "url: ", item.URL)
			}
		} else {
			log.Printf("read yaml error %v, data: %s", err, string(data))
		}
	}
}

// 素材的类型，图片（image）、视频（video）、语音 （voice）、图文（news）
type ListRequest struct {
	Type   string `json:"type"`
	Offset int    `json:"offset"`
	Count  int    `json:"count"`
}

type MediaItemList struct {
	TotalCouont int         `json:"total_count"`
	ItemCount   int         `json:"item_count"`
	ItemList    []MediaItem `json:"item"`
}

type MediaItem struct {
	MediaID    string `json:"media_id"`
	Name       string `json:"name"`
	UpdateTime string `json:"update_time"`
	URL        string `json:"url"`
}
