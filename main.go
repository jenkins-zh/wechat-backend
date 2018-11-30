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
)

const (
	token = "wechat4go"
)

func makeSignature(timestamp, nonce string) string {
	sl := []string{token, timestamp, nonce}
	sort.Strings(sl)
	s := sha1.New()
	io.WriteString(s, strings.Join(sl, ""))
	return fmt.Sprintf("%x", s.Sum(nil))
}

func validateUrl(w http.ResponseWriter, r *http.Request) bool {
	timestamp := strings.Join(r.Form["timestamp"], "")
	nonce := strings.Join(r.Form["nonce"], "")
	signatureGen := makeSignature(timestamp, nonce)

	signatureIn := strings.Join(r.Form["signature"], "")
	if signatureGen != signatureIn {
		return false
	}
	echostr := strings.Join(r.Form["echostr"], "")
	fmt.Fprintf(w, echostr)
	return true
}

func procRequest(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if !validateUrl(w, r) {
		log.Println("Wechat Service: this http request is not from Wechat platform!")
		return
	}
	log.Println("Wechat Service: validateUrl Ok!")

	if r.Method == "POST" {
		textRequestBody := parseTextRequestBody(r)
		if textRequestBody != nil {
			fmt.Printf("Wechat Service: Recv text msg [%s] from user [%s]!",
				textRequestBody.Content,
				textRequestBody.FromUserName)

			if "event" == textRequestBody.MsgType && "subscribe" == textRequestBody.Event {
				resp, err := makeWelcomeResponseBody(textRequestBody.ToUserName, textRequestBody.FromUserName)
				if err != nil {
					log.Println("Wechat Service: makeTextResponseBody error: ", err)
					return
				}
				fmt.Fprintf(w, string(resp))
			} else {
				keyword := textRequestBody.Content
				fmt.Println(textRequestBody.MsgType, keyword, respMap)
				if "text" == textRequestBody.MsgType {
					if resp, ok := respMap[keyword]; ok {
						if text, ok := resp.(TextResponseBody); ok {
							textResp, err := makeTextResponseBody(textRequestBody.ToUserName, textRequestBody.FromUserName, text.Content)
							if err != nil {
								log.Println("Wechat Service: makeTextResponseBody error: ", err)
								return
							}
							fmt.Fprintf(w, string(textResp))
							return
						} else {
							log.Println("type error", ok)
						}
					} else {
						log.Printf("can't find keyword %s\n", keyword)
					}
				}

				responseTextBody, err := makeNewsResponseBody(textRequestBody.ToUserName,
					textRequestBody.FromUserName,
					"Hello, "+textRequestBody.FromUserName)
				if err != nil {
					log.Println("Wechat Service: makeTextResponseBody error: ", err)
					return
				}
				fmt.Fprintf(w, string(responseTextBody))
			}
		}
	}
}

func makeTextResponseBody(fromUserName, toUserName, content string) ([]byte, error) {
	textResponseBody := &TextResponseBody{}
	textResponseBody.FromUserName = fromUserName
	textResponseBody.ToUserName = toUserName
	textResponseBody.MsgType = "text"
	textResponseBody.Content = content
	textResponseBody.CreateTime = time.Duration(time.Now().Unix())
	return xml.MarshalIndent(textResponseBody, " ", "  ")
}

func makeWelcomeResponseBody(fromUserName string, toUserName string) ([]byte, error) {
	return makeTextResponseBody(fromUserName, toUserName, "welcome")
}

func makeNewsResponseBody(fromUserName, toUserName, content string) ([]byte, error) {
	newsResponseBody := &NewsResponseBody{}
	newsResponseBody.FromUserName = fromUserName
	newsResponseBody.ToUserName = toUserName
	newsResponseBody.MsgType = "news"
	newsResponseBody.ArticleCount = 1
	newsResponseBody.Articles = Articles{
		Articles: []Article{{
			Title:       "来自 Jenkins 官方的消息",
			Description: "来自 Jenkins 官方的消息",
			PicUrl:      "https://raw.githubusercontent.com/jenkins-infra/wechat/master/images/vscode-pipeline-linter/example1.gif",
			Url:         "https://mp.weixin.qq.com/s/4pktvfQ3tJZgqY--VgNgZQ",
		}},
	}
	newsResponseBody.CreateTime = time.Duration(time.Now().Unix())
	return xml.MarshalIndent(newsResponseBody, " ", "  ")
}

func parseTextRequestBody(r *http.Request) *TextRequestBody {
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

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func main() {
	initCheck()
	createWxMenu()
	http.HandleFunc("/", procRequest)
	http.HandleFunc("/status", healthHandler)

	log.Fatal(http.ListenAndServe(":18080", nil))
}
