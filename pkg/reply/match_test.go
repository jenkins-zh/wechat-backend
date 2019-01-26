package reply

import (
	"testing"

	"github.com/golang/mock/gomock"
	core "github.com/linuxsuren/wechat-backend/pkg"
	mArticle "github.com/linuxsuren/wechat-backend/pkg/mock/article"
)

func TestAccept(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	var reply AutoReply
	reply = &MatchAutoReply{}

	if reply.Accept(&core.TextRequestBody{}) {
		t.Errorf("should not accept")
	}

	m := mArticle.NewMockResponseManager(ctrl)
	m.EXPECT().GetResponse("hello").
		Return(&core.TextResponseBody{}, true)
	SetResponseManager(m)

	if !reply.Accept(&core.TextRequestBody{
		MsgType: "text",
		Content: "hello",
	}) {
		t.Errorf("should accept")
	}
}

func TestHandle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	reply := &MatchAutoReply{}

	m := mArticle.NewMockResponseManager(ctrl)
	m.EXPECT().GetResponse("hello").
		Return(core.TextResponseBody{
			ResponseBody: core.ResponseBody{
				MsgType: "text",
			},
			Content: "hello",
		}, true)

	SetResponseManager(m)
	if !reply.Accept(&core.TextRequestBody{
		MsgType: "text",
		Content: "hello",
	}) {
		t.Errorf("should accept")
	}

	data, err := reply.Handle()
	if err != nil {
		t.Errorf("should not error %v", err)
	} else if string(data) != "hello" {
		t.Errorf("got an error content: %s", string(data))
	}
}
