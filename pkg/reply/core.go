package reply

import (
	"encoding/xml"
	"time"

	core "github.com/linuxsuren/wechat-backend/pkg"
)

// AutoReply represent auto reply interface
type AutoReply interface {
	Accept(request *core.TextRequestBody) bool
	Handle() (string, error)
	Name() string
	Weight() int
}

type Init func() AutoReply

func makeTextResponseBody(fromUserName, toUserName string, content string) ([]byte, error) {
	textResponseBody := &core.TextResponseBody{}
	textResponseBody.FromUserName = fromUserName
	textResponseBody.ToUserName = toUserName
	textResponseBody.MsgType = "text"
	textResponseBody.Content = content
	textResponseBody.CreateTime = time.Duration(time.Now().Unix())
	return marshal(textResponseBody)
}

func makeImageResponseBody(fromUserName, toUserName, mediaID string) ([]byte, error) {
	imageResponseBody := &core.ImageResponseBody{}
	imageResponseBody.FromUserName = fromUserName
	imageResponseBody.ToUserName = toUserName
	imageResponseBody.MsgType = "image"
	imageResponseBody.CreateTime = time.Duration(time.Now().Unix())
	imageResponseBody.Image = core.Image{
		MediaID: mediaID,
	}
	return marshal(imageResponseBody)
}

func makeNewsResponseBody(fromUserName, toUserName string, news core.NewsResponseBody) ([]byte, error) {
	newsResponseBody := &core.NewsResponseBody{}
	newsResponseBody.FromUserName = fromUserName
	newsResponseBody.ToUserName = toUserName
	newsResponseBody.MsgType = "news"
	newsResponseBody.ArticleCount = 1
	newsResponseBody.Articles = core.Articles{
		Articles: news.Articles.Articles,
	}
	newsResponseBody.CreateTime = time.Duration(time.Now().Unix())
	return marshal(newsResponseBody)
}

func marshal(response interface{}) ([]byte, error) {
	return xml.MarshalIndent(response, " ", "  ")
}
