package reply

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	core "github.com/linuxsuren/wechat-backend/pkg"
)

// GitterAutoReply only reply for match
type GitterAutoReply struct {
	Request  *core.TextRequestBody
	Keyword  string
	Callback string
}

var _ AutoReply = &GitterAutoReply{}

// Name indicate reply's name
func (m *GitterAutoReply) Name() string {
	return "GitterAutoReply"
}

// Weight weight for order
func (m *GitterAutoReply) Weight() int {
	return 0
}

// Accept consider if it will accept the request
func (m *GitterAutoReply) Accept(request *core.TextRequestBody) (ok bool) {
	m.Request = request
	m.Keyword = request.Content
	m.Keyword = strings.TrimLeft(m.Keyword, "问 ")
	m.Keyword = strings.TrimLeft(m.Keyword, "q ")

	if "text" != request.MsgType {
		return false
	}

	return strings.HasPrefix(request.Content, "问") ||
		strings.HasPrefix(request.Content, "q")
}

// Handle hanlde the request then return data
func (m *GitterAutoReply) Handle() (string, error) {
	from := m.Request.ToUserName
	to := m.Request.FromUserName
	var err error
	var data []byte

	binder := &GitHubBinder{
		File: "config/github_bind.yaml",
	}

	var sender string
	gitHubBindData := binder.Find(to)
	if gitHubBindData == nil {
		sender = "anonymous"
	} else {
		sender = gitHubBindData.GitHubID
	}

	sendMsg(m.Callback, fmt.Sprintf("@%s %s", sender, m.Keyword))
	data, err = makeTextResponseBody(from, to, "sent")
	if err != nil {
		err = fmt.Errorf("Wechat Service: makeTextResponseBody error: %v", err)
	}
	return string(data), err
}

func init() {
	Register(func() AutoReply {
		return &GitterAutoReply{
			Callback: "https://webhooks.gitter.im/e/911738f12cb4ca5d3c41",
		}
	})
}

// SendMsg send message to server
func sendMsg(server, message string) {
	value := url.Values{
		"message": {message},
	}

	http.PostForm(server, value)
}
