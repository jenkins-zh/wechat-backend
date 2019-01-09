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

	"github.com/linuxsuren/wechat-backend/pkg/config"
	"github.com/linuxsuren/wechat-backend/pkg/github"
	"github.com/linuxsuren/wechat-backend/pkg/reply"
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

		for _, autoReplyInit := range autoReplyInitChains {
			autoReply := autoReplyInit()
			if !autoReply.Accept(textRequestBody) {
				continue
			}

			var data []byte
			var err error
			if data, err = autoReply.Handle(); err != nil {
				log.Println("handle auto replay error:", err)
			}
			fmt.Fprintf(w, string(data))
			break
		}
	}
}

func (w *WeChat) parseTextRequestBody(r *http.Request) *reply.TextRequestBody {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	fmt.Println(string(body))
	requestBody := &reply.TextRequestBody{}
	xml.Unmarshal(body, requestBody)
	return requestBody
}

var autoReplyInitChains []reply.Init

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

	autoReplyInitChains = make([]reply.Init, 1)
	autoReplyInitChains = append(autoReplyInitChains, reply.InitMatchAutoReply)

	http.HandleFunc("/", wechat.procRequest)
	http.HandleFunc("/status", healthHandler)
	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		github.WebhookHandler(w, r, weConfig, initCheck)
	})

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", weConfig.ServerPort), nil))
}
