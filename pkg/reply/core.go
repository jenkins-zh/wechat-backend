package reply

import (
	"encoding/xml"
	"time"
)

type AutoReply interface {
	Accept(request *TextRequestBody) bool
	Handle() ([]byte, error)
}

type Init func() AutoReply

func makeTextResponseBody(fromUserName, toUserName, content string) ([]byte, error) {
	textResponseBody := &TextResponseBody{}
	textResponseBody.FromUserName = fromUserName
	textResponseBody.ToUserName = toUserName
	textResponseBody.MsgType = "text"
	textResponseBody.Content = content
	textResponseBody.CreateTime = time.Duration(time.Now().Unix())
	return xml.MarshalIndent(textResponseBody, " ", "  ")
}

func makeImageResponseBody(fromUserName, toUserName, mediaID string) ([]byte, error) {
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

func makeNewsResponseBody(fromUserName, toUserName string, news NewsResponseBody) ([]byte, error) {
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
