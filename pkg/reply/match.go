package reply

import (
	"fmt"
	"log"

	core "github.com/linuxsuren/wechat-backend/pkg"
	"github.com/linuxsuren/wechat-backend/pkg/article"
)

var responseManager article.ResponseManager

func SetResponseManager(manager article.ResponseManager) {
	responseManager = manager
}

// MatchAutoReply only reply for match
type MatchAutoReply struct {
	Response interface{}
	Request  *core.TextRequestBody
}

var _ AutoReply = &MatchAutoReply{}

func (m *MatchAutoReply) Name() string {
	return "SimpleMatchReply"
}

// Accept consider if it will accept the request
func (m *MatchAutoReply) Accept(request *core.TextRequestBody) (ok bool) {
	m.Request = request
	keyword := request.Content

	fmt.Printf("request is %v\n", request)

	if responseManager == nil || "text" != request.MsgType {
		log.Printf("responseManager is nil or not support msgType %s", request.MsgType)
		return false
	}

	m.Response, ok = responseManager.GetResponse(keyword)
	return ok
}

// Handle hanlde the request then return data
func (m *MatchAutoReply) Handle() (string, error) {
	resp := m.Response
	from := m.Request.ToUserName
	to := m.Request.FromUserName
	var err error

	if resp == nil {
		err = fmt.Errorf("response is nil")
		return "", err
	}

	fmt.Printf("response %v\n", resp)

	var data []byte
	if text, ok := resp.(core.TextResponseBody); ok {
		data, err = makeTextResponseBody(from, to, text.Content)
		fmt.Printf("data %v\n", string(data))
		if err != nil {
			err = fmt.Errorf("Wechat Service: makeTextResponseBody error: %v", err)
		}
	} else if image, ok := resp.(core.ImageResponseBody); ok {
		data, err = makeImageResponseBody(from, to, image.Image.MediaID)
		if err != nil {
			err = fmt.Errorf("Wechat Service: makeImageResponseBody error: %v", err)
		}
	} else if news, ok := resp.(core.NewsResponseBody); ok {
		data, err = makeNewsResponseBody(from, to, news)
		if err != nil {
			err = fmt.Errorf("Wechat Service: makeNewsResponseBody error: %v", err)
		}
	} else {
		err = fmt.Errorf("type error %v", resp)
	}

	return string(data), err
}

func init() {
	Register(func() AutoReply {
		return &MatchAutoReply{}
	})
}
