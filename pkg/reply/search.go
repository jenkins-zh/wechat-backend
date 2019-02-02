package reply

import (
	"fmt"
	"strings"

	core "github.com/linuxsuren/wechat-backend/pkg"
	"github.com/linuxsuren/wechat-backend/pkg/article"
)

// SearchAutoReply only reply for match
type SearchAutoReply struct {
	ResponseMap map[string]interface{}
	Response    interface{}
	Request     *core.TextRequestBody
	Keyword     string
}

var _ AutoReply = &SearchAutoReply{}

func (m *SearchAutoReply) Name() string {
	return "SearchAutoReply"
}

func (m *SearchAutoReply) Weight() int {
	return 0
}

// Accept consider if it will accept the request
func (m *SearchAutoReply) Accept(request *core.TextRequestBody) (ok bool) {
	m.Request = request
	m.Keyword = request.Content
	m.Keyword = strings.TrimLeft(m.Keyword, "search ")
	m.Keyword = strings.TrimLeft(m.Keyword, "搜索 ")

	if "text" != request.MsgType {
		return false
	}

	return strings.HasPrefix(request.Content, "搜索") ||
		strings.HasPrefix(request.Content, "search")
}

// Handle hanlde the request then return data
func (m *SearchAutoReply) Handle() (string, error) {
	from := m.Request.ToUserName
	to := m.Request.FromUserName
	var err error

	reader := &article.ArticleReader{
		API: "https://jenkins-zh.github.io/index.json",
	}

	var data []byte
	articles, err := reader.FindByTitle(m.Keyword)
	if err != nil {
		return "", err
	}

	fmt.Printf("found aritcle count [%d]\n", len(articles))
	var targetArticle article.Article
	if len(articles) == 0 {
		targetArticle = article.Article{
			Title:       "404",
			Description: "没有找到相关的文章，给我们留言，或者直接发 PR 过来！",
			URI:         "https://jenkins-zh.github.io",
		}
	} else {
		targetArticle = articles[0]
	}

	news := core.NewsResponseBody{
		Articles: core.Articles{
			Articles: []core.Article{
				{
					Title:       targetArticle.Title,
					Description: targetArticle.Description,
					PicUrl:      "https://mmbiz.qpic.cn/mmbiz_jpg/nh8sibXrvHrvicMyefXop7qwrnWQc5gtBgia05BicxFCWdjPkee3Ku9FLwBZR3JJVDwvVDL25p90BLPOTOWUCrribLA/0?wx_fmt=jpeg",
					Url:         targetArticle.URI,
				},
			},
		},
	}
	data, err = makeNewsResponseBody(from, to, news)
	if err != nil {
		err = fmt.Errorf("Wechat Service: makeNewsResponseBody error: %v", err)
	}
	return string(data), err
}

func init() {
	Register(func() AutoReply {
		return &SearchAutoReply{}
	})
}
