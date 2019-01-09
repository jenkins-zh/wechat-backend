package reply

import "fmt"

// MatchAutoReply only reply for match
type MatchAutoReply struct {
	ResponseMap map[string]interface{}
	Response    interface{}
	Request     *TextRequestBody
}

func InitMatchAutoReply() AutoReply {
	return &MatchAutoReply{}
}

var _ AutoReply = &MatchAutoReply{}

// Accept consider if it will accept the request
func (m *MatchAutoReply) Accept(request *TextRequestBody) (ok bool) {
	m.Request = request
	keyword := request.Content

	if "text" != request.MsgType {
		return false
	}

	m.Response, ok = m.ResponseMap[keyword]
	return ok
}

// Handle hanlde the request then return data
func (m *MatchAutoReply) Handle() (data []byte, err error) {
	resp := m.Response
	from := m.Request.FromUserName
	to := m.Request.ToUserName

	if text, ok := resp.(TextResponseBody); ok {
		data, err = makeTextResponseBody(from, to, text.Content)
		if err != nil {
			err = fmt.Errorf("Wechat Service: makeTextResponseBody error: %v", err)
		}
	} else if image, ok := resp.(ImageResponseBody); ok {
		data, err = makeImageResponseBody(from, to, image.Image.MediaID)
		if err != nil {
			err = fmt.Errorf("Wechat Service: makeImageResponseBody error: %v", err)
		}
	} else if news, ok := resp.(NewsResponseBody); ok {
		data, err = makeNewsResponseBody(from, to, news)
		if err != nil {
			err = fmt.Errorf("Wechat Service: makeNewsResponseBody error: %v", err)
		}
	} else {
		err = fmt.Errorf("type error")
	}
	return
}
