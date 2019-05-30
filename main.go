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

	core "github.com/jenkins-zh/wechat-backend/pkg"
	"github.com/jenkins-zh/wechat-backend/pkg/article"
	"github.com/jenkins-zh/wechat-backend/pkg/config"
	"github.com/jenkins-zh/wechat-backend/pkg/github"
	"github.com/jenkins-zh/wechat-backend/pkg/reply"
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

	// log.Println("request url", r.URL.String())
	// if strings.HasPrefix(r.URL.String(), "/?signature=") {
	// 	log.Println("just for valid")
	// 	return
	// }

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

func (we *WeChat) wechatRequest(writer http.ResponseWriter, r *http.Request) {
	textRequestBody := we.parseTextRequestBody(r)
	if textRequestBody != nil {
		autoReplyInitChains := reply.AutoReplyChains()
		fmt.Printf("found [%d] autoReply", len(autoReplyInitChains))

		var potentialReplys []reply.AutoReply
		for _, autoReplyInit := range autoReplyInitChains {
			if autoReplyInit == nil {
				fmt.Printf("found a nil autoReply.")
				continue
			}
			autoReply := autoReplyInit()
			if !autoReply.Accept(textRequestBody) {
				continue
			}

			potentialReplys = append(potentialReplys, autoReply)
		}

		sort.Sort(reply.ByWeight(potentialReplys))

		if len(potentialReplys) > 0 {
			autoReply := potentialReplys[0]
			fmt.Printf("going to handle by %s\n", autoReply.Name())

			if data, err := autoReply.Handle(); err != nil {
				fmt.Printf("handle auto replay error: %v\n", err)
			} else if len(data) == 0 {
				fmt.Println("response body is empty.")
			} else {
				fmt.Printf("response:%s\n", data)
				fmt.Fprintf(writer, data)
			}
		} else {
			fmt.Println("should have at least one reply")
		}
	}
}

func (w *WeChat) parseTextRequestBody(r *http.Request) *core.TextRequestBody {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	fmt.Println(string(body))
	requestBody := &core.TextRequestBody{}
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

	defaultRM := article.NewDefaultResponseManager()
	reply.SetResponseManager(defaultRM)

	wechat := WeChat{
		Config: weConfig,
	}
	go func() {
		defaultRM.InitCheck(weConfig)
	}()
	createWxMenu(weConfig)

	http.HandleFunc("/", wechat.procRequest)
	http.HandleFunc("/status", healthHandler)
	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		github.WebhookHandler(w, r, weConfig, defaultRM.InitCheck)
	})

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", weConfig.ServerPort), nil))
}
