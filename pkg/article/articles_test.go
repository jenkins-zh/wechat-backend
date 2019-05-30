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

	responseHandler([]byte(yml))

	resp := respMap["hi"]
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

	responseHandler([]byte(yml))

	resp := respMap["about"]
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
