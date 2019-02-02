package reply

import (
	"fmt"

	core "github.com/linuxsuren/wechat-backend/pkg"
)

// WelcomeReply for welcome event
type WelcomeReply struct {
	AutoReply
}

var _ AutoReply = &WelcomeReply{}

// Name represents the name for reply
func (m *WelcomeReply) Name() string {
	return "WelcomeReply"
}

// Accept consider if it will accept the request
func (m *WelcomeReply) Accept(request *core.TextRequestBody) (ok bool) {
	if "event" == request.MsgType && "subscribe" == request.Event {
		request.Content = "welcome"
		request.MsgType = "text"
		m.AutoReply = &MatchAutoReply{}
		ok = m.AutoReply.Accept(request)
	}
	return
}

func (m *WelcomeReply) Weight() int {
	return 0
}

func init() {
	fmt.Println("register for welcome")
	Register(func() AutoReply {
		return &WelcomeReply{}
	})
}
