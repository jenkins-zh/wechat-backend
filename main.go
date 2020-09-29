package main

import (
	"crypto/sha1"
	"encoding/xml"
	"fmt"
	"github.com/jenkins-zh/wechat-backend/pkg/api"
	"github.com/jenkins-zh/wechat-backend/pkg/health"
	"github.com/jenkins-zh/wechat-backend/pkg/menu"
	"github.com/jenkins-zh/wechat-backend/pkg/service"
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

func (we *WeChat) makeSignature(timestamp, nonce string) string {
	sl := []string{we.Config.Token, timestamp, nonce}
	sort.Strings(sl)
	s := sha1.New()
	io.WriteString(s, strings.Join(sl, ""))
	return fmt.Sprintf("%x", s.Sum(nil))
}

func (we *WeChat) validateURL(w http.ResponseWriter, r *http.Request) bool {
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
	if !we.validateURL(w, r) {
		we.normalRequest(w, r)
		log.Println("Wechat Service: this http request is not from Wechat platform!")
		return
	}
	log.Println("Wechat Service: validateURL Ok!")

	if we.Config.Valid {
		log.Println("request url", r.URL.String())
		if strings.HasPrefix(r.URL.String(), "/?signature=") {
			log.Println("just for valid")
			return
		}
	}

	switch r.Method {
	case http.MethodPost:
		we.wechatRequest(w, r)
	}
}

func (we *WeChat) normalRequest(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("Welcome aboard Jenkins WeChat.")); err != nil {
		fmt.Printf("got error when response normal request, %v", err)
	}
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

func (we *WeChat) parseTextRequestBody(r *http.Request) *core.TextRequestBody {
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
	configurator := &config.LocalFileConfig{}
	weConfig, err := configurator.LoadConfig(core.ConfigPath)
	if err != nil {
		log.Printf("load config error %v\n", err)
	}

	// TODO this should be handle by config function
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
	menu.CreateWxMenu(weConfig)

	http.HandleFunc("/", wechat.procRequest)
	http.HandleFunc("/status", health.SimpleHealthHandler)
	//http.HandleFunc("/status", healthHandler)
	http.HandleFunc("/medias", func(w http.ResponseWriter, r *http.Request) {
		api.ListMedias(w, r, configurator)
	})
	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		github.WebhookHandler(w, r, weConfig, defaultRM.InitCheck)
	})
	http.HandleFunc("/config", func(w http.ResponseWriter, r *http.Request) {
		service.HandleConfig(w, r, configurator)
	})

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", weConfig.ServerPort), nil))
}
