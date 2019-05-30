package reply

import (
	"testing"

	"github.com/golang/mock/gomock"
	core "github.com/jenkins-zh/wechat-backend/pkg"
	mArticle "github.com/jenkins-zh/wechat-backend/pkg/mock/article"
)

func TestWelcome(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	reply := WelcomeReply{}
	request := &core.TextRequestBody{
		MsgType: "event",
		Event:   "subscribe",
	}

	m := mArticle.NewMockResponseManager(ctrl)
	m.EXPECT().GetResponse("welcome").
		Return(core.TextResponseBody{
			ResponseBody: core.ResponseBody{
				MsgType: "text",
			},
			Content: "welcome",
		}, true)

	SetResponseManager(m)

	if !reply.Accept(request) {
		t.Errorf("should accept")
	}

	if _, err := reply.Handle(); err != nil {
		t.Errorf("should not error: %v", err)
	}
}
