package main

import (
	"crypto/sha1"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/linuxsuren/wechat-backend/pkg/config"
	"github.com/linuxsuren/wechat-backend/pkg/github"
)

// WeChat represents WeChat
type WeChat struct {
	Config *config.WeChatConfig
}

func (w *WeChat) makeSignature(timestamp, nonce string) string {
	sl := []string{w.Config.Token, timestamp, nonce}
	sort.Strings(sl)
	s := sha1.New()
	io.WriteString(s, strings.Join(sl, ""))
	return fmt.Sprintf("%x", s.Sum(nil))
}

func (we *WeChat) validateUrl(w http.ResponseWriter, r *http.Request) bool {
	timestamp := strings.Join(r.Form["timestamp"], "")
	nonce := strings.Join(r.Form["nonce"], "")
	signatureGen := we.makeSignature(timestamp, nonce)

	signatureIn := strings.Join(r.Form["signature"], "")
	if signatureGen != signatureIn {
		return false
	}
	echostr := strings.Join(r.Form["echostr"], "")
	fmt.Fprintf(w, echostr)
	return true
}

func (we *WeChat) procRequest(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if !we.validateUrl(w, r) {
		log.Println("Wechat Service: this http request is not from Wechat platform!")
		return
	}
	log.Println("Wechat Service: validateUrl Ok!")

	switch r.Method {
	case http.MethodPost:
		we.wechatRequest(w, r)
	case http.MethodGet:
		we.normalRequest(w, r)
	}
}

func (we *WeChat) normalRequest(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("welcome aboard WeChat."))
}

func (we *WeChat) wechatRequest(w http.ResponseWriter, r *http.Request) {
	textRequestBody := we.parseTextRequestBody(r)
	if textRequestBody != nil {
		fmt.Printf("Wechat Service: Recv [%s] msg [%s] from user [%s]!\n",
			textRequestBody.MsgType,
			textRequestBody.Content,
			textRequestBody.FromUserName)

		if "event" == textRequestBody.MsgType && "subscribe" == textRequestBody.Event {
			resp, err := we.replyResponse("welcome", textRequestBody.ToUserName, textRequestBody.FromUserName)
			if err != nil {
				log.Println("handle welcome replay error:", err)
			} else {
				fmt.Fprintf(w, string(resp))
			}
		} else {
			keyword := textRequestBody.Content
			log.Println(textRequestBody.MsgType, keyword, respMap)
			if "text" == textRequestBody.MsgType {
				resp, err := we.replyResponse(keyword, textRequestBody.ToUserName, textRequestBody.FromUserName)
				if err != nil {
					log.Println("handle auto replay error:", err)
				} else {
					log.Println("response", string(resp))
					fmt.Fprintf(w, string(resp))
				}
			}
		}
	}
}

func (we *WeChat) replyResponse(keyword string, from string, to string) (data []byte, err error) {
	if resp, ok := respMap[keyword]; ok {
		if text, ok := resp.(TextResponseBody); ok {
			data, err = we.makeTextResponseBody(from, to, text.Content)
			if err != nil {
				err = fmt.Errorf("Wechat Service: makeTextResponseBody error: %v", err)
			}
		} else if image, ok := resp.(ImageResponseBody); ok {
			data, err = we.makeImageResponseBody(from, to, image.Image.MediaID)
			if err != nil {
				err = fmt.Errorf("Wechat Service: makeImageResponseBody error: %v", err)
			}
		} else if news, ok := resp.(NewsResponseBody); ok {
			data, err = we.makeNewsResponseBody(from, to, news)
			if err != nil {
				err = fmt.Errorf("Wechat Service: makeNewsResponseBody error: %v", err)
			}
		} else {
			err = fmt.Errorf("type error")
		}
	}
	return
}

func (w *WeChat) makeTextResponseBody(fromUserName, toUserName, content string) ([]byte, error) {
	textResponseBody := &TextResponseBody{}
	textResponseBody.FromUserName = fromUserName
	textResponseBody.ToUserName = toUserName
	textResponseBody.MsgType = "text"
	textResponseBody.Content = content
	textResponseBody.CreateTime = time.Duration(time.Now().Unix())
	return xml.MarshalIndent(textResponseBody, " ", "  ")
}

func (w *WeChat) makeImageResponseBody(fromUserName, toUserName, mediaID string) ([]byte, error) {
	imageResponseBody := &ImageResponseBody{}
	imageResponseBody.FromUserName = fromUserName
	imageResponseBody.ToUserName = toUserName
	imageResponseBody.MsgType = "image"
	imageResponseBody.CreateTime = time.Duration(time.Now().Unix())
	imageResponseBody.Image = Image{
		MediaID: mediaID,
	}
	return xml.MarshalIndent(imageResponseBody, " ", "  ")
}

func (w *WeChat) makeNewsResponseBody(fromUserName, toUserName string, news NewsResponseBody) ([]byte, error) {
	newsResponseBody := &NewsResponseBody{}
	newsResponseBody.FromUserName = fromUserName
	newsResponseBody.ToUserName = toUserName
	newsResponseBody.MsgType = "news"
	newsResponseBody.ArticleCount = 1
	newsResponseBody.Articles = Articles{
		Articles: news.Articles.Articles,
	}
	newsResponseBody.CreateTime = time.Duration(time.Now().Unix())
	return xml.MarshalIndent(newsResponseBody, " ", "  ")
}

func (w *WeChat) parseTextRequestBody(r *http.Request) *TextRequestBody {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	fmt.Println(string(body))
	requestBody := &TextRequestBody{}
	xml.Unmarshal(body, requestBody)
	return requestBody
}

func main() {
	weConfig, err := config.LoadConfig("config/wechat.yaml")
	if err != nil {
		log.Printf("load config error %v\n", err)
	}

	if weConfig == nil {
		weConfig = &config.WeChatConfig{
			ServerPort: 8080,
		}
	}

	if weConfig.ServerPort <= 0 {
		weConfig.ServerPort = 8080
	}

	wechat := WeChat{
		Config: weConfig,
	}
	go func() {
		initCheck(weConfig)
	}()
	createWxMenu()

	http.HandleFunc("/", wechat.procRequest)
	http.HandleFunc("/status", healthHandler)
	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		github.WebhookHandler(w, r, weConfig, initCheck)
	})

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", weConfig.ServerPort), nil))
}
