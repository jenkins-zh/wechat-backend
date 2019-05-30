package reply

import (
	"testing"

	core "github.com/jenkins-zh/wechat-backend/pkg"
)

func TestSearch(t *testing.T) {
	reply := SearchAutoReply{}

	reply.Accept(&core.TextRequestBody{})

	data, err := reply.Handle()
	if err != nil {
		t.Errorf("error %v", err)
	}

	t.Errorf("%s", string(data))
}
