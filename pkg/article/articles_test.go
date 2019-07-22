package article

import (
	"testing"

	core "github.com/jenkins-zh/wechat-backend/pkg"
	"github.com/stretchr/testify/assert"
)

func TestImageResponseBody(t *testing.T) {
	yml := `
msgType: image

keyword: hi
content: say hello from jenkins.
image:
  mediaID: mediaId
`

	mgr := NewDefaultResponseManager()
	mgr.responseHandler([]byte(yml))

	resp := mgr.ResponseMap["hi"]
	if resp == nil {
		t.Error("Can't find response by keyword: hi.")
	}

	imageResp, ok := resp.(core.ImageResponseBody)
	if !ok {
		t.Error("Get the wrong type, should be ImageResponseBody.")
	}
	assert.Equal(t, imageResp.Image.MediaID, "mediaId", "ImageResponseBody parse error, can't find the correct mediaId: ", imageResp.Image.MediaID)
}

func TestNewsResponseBody(t *testing.T) {
	yml := `
keyword: about

msgType: news
articleCount: 1
articles:
- title: "title"
  description: "desc"
  picUrl: "http://pic.com"
  url: "http://blog.com"
`

	mgr := NewDefaultResponseManager()
	mgr.responseHandler([]byte(yml))

	resp := mgr.ResponseMap["about"]
	if resp == nil {
		t.Error("Can't find response by keyword: about.")
		return
	}

	newsResp, ok := resp.(core.NewsResponseBody)
	if !ok {
		t.Error("Get the wrong type, should be NewsResponseBody.")
	}
	assert.Equal(t, newsResp.Articles.Articles[0].Title, "title", "title parse error.")
}

func TestRandomResponseBody(t *testing.T) {
	yml := `
keyword: weixin
msgType: random
items:
- abc
- def
`

	mgr := NewDefaultResponseManager()
	mgr.responseHandler([]byte(yml))

	resp := mgr.ResponseMap["weixin"]
	if resp == nil {
		t.Error("Can't find response by keyword: weixin.")
		return
	}

	newsResp, ok := resp.(core.RandomResponseBody)
	if !ok {
		t.Error("Get the wrong type, should be RandomResponseBody.")
	}
	assert.Equal(t, len(newsResp.Items), 2, "can not parse items")
	assert.Equal(t, newsResp.Items[0], "abc")
}
