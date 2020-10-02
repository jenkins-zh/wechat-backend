package reply

import (
	"math"

	core "github.com/jenkins-zh/wechat-backend/pkg"
)

// UnknownAutoReply unknown auto reply
type UnknownAutoReply struct {
	Request *core.TextRequestBody
}

// Name represent name for current auto reply
func (u *UnknownAutoReply) Name() string {
	return "UnknownReply"
}

// Accept all keywords
func (u *UnknownAutoReply) Accept(request *core.TextRequestBody) bool {
	u.Request = request
	return true
}

// Handle take care of unknown things
func (u *UnknownAutoReply) Handle() (string, error) {
	from := u.Request.ToUserName
	to := u.Request.FromUserName
	commonReply := `您的提的问题已经远远超过了我的智商，请回复"小助手"，社区机器人会把您拉进群里。更多关键字，请回复"帮助"。`

	// try to find a configured reply sentence
	if response, ok := responseManager.GetResponse("unknown"); ok && response != nil {
		if text, ok := response.(core.TextResponseBody); ok {
			commonReply = text.Content
		}
	}

	data, err := makeTextResponseBody(from, to, commonReply)
	return string(data), err
}

// Weight should be the last one
func (u *UnknownAutoReply) Weight() int {
	return math.MaxInt64
}

var _ AutoReply = &UnknownAutoReply{}

func init() {
	Register(func() AutoReply {
		return &UnknownAutoReply{}
	})
}
