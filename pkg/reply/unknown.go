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
	data, err := makeTextResponseBody(from, to, `
	我貌似没有明白你的问题，请回复“聊天室”问问其他人吧。
	`)
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
