package reply

import "testing"

func TestAccept(t *testing.T) {
	var reply AutoReply

	reply = &MatchAutoReply{
		ResponseMap: map[string]interface{}{
			"hello": "",
		},
	}

	if reply.Accept(&TextRequestBody{}) {
		t.Errorf("should not accept")
	}

	if !reply.Accept(&TextRequestBody{
		Content: "hello",
	}) {
		t.Errorf("should accept")
	}
}

func TestHandle(t *testing.T) {
	var reply AutoReply

	reply = &MatchAutoReply{
		ResponseMap: map[string]interface{}{
			"hello": TextResponseBody{},
		},
	}

	if !reply.Accept(&TextRequestBody{
		Content: "hello",
	}) {
		t.Errorf("should accept")
	}

	data, err := reply.Handle()
	if err != nil {
		t.Errorf("should not error %v", err)
	} else if data == nil {
		t.Errorf("not have data")
	}
}
