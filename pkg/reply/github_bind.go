package reply

import (
	"fmt"
	"os"
	"strings"

	"io/ioutil"

	core "github.com/linuxsuren/wechat-backend/pkg"
	yaml "gopkg.in/yaml.v2"
)

// GitHubBindAutoReply only reply for match
type GitHubBindAutoReply struct {
	Request    *core.TextRequestBody
	GitHubBind GitHubBind
	Event      string
	Keyword    string
}

var _ AutoReply = &GitHubBindAutoReply{}

const (
	GitHubEventRegister   = "注册"
	GitHubEventUnregister = "注销"
)

// Name indicate reply's name
func (m *GitHubBindAutoReply) Name() string {
	return "GitHubBindAutoReply"
}

// Weight weight for order
func (m *GitHubBindAutoReply) Weight() int {
	return 0
}

// Accept consider if it will accept the request
func (m *GitHubBindAutoReply) Accept(request *core.TextRequestBody) bool {
	m.Request = request
	m.Keyword = request.Content

	if "text" != request.MsgType {
		return false
	}

	if strings.HasPrefix(request.Content, GitHubEventRegister) {
		m.Event = GitHubEventRegister
		m.Keyword = strings.TrimLeft(m.Keyword, fmt.Sprintf("%s ", GitHubEventRegister))
	} else if strings.HasPrefix(request.Content, GitHubEventUnregister) {
		m.Event = GitHubEventUnregister
		m.Keyword = strings.TrimLeft(m.Keyword, fmt.Sprintf("%s ", GitHubEventUnregister))
	} else {
		return false
	}

	return true
}

// Handle hanlde the request then return data
func (m *GitHubBindAutoReply) Handle() (string, error) {
	from := m.Request.ToUserName
	to := m.Request.FromUserName
	var reply string

	if m.Keyword != "" {
		switch m.Event {
		case GitHubEventRegister:
			m.GitHubBind.Add(GitHubBindData{
				WeChatID: to,
				GitHubID: m.Keyword,
			})
			reply = "register success"
		case GitHubEventUnregister:
			m.GitHubBind.Remove(to)
			reply = "unregister success"
		default:
			reply = "unknow event"
		}
	} else {
		reply = "need your github id"
	}

	var err error
	var data []byte
	data, err = makeTextResponseBody(from, to, reply)
	if err != nil {
		err = fmt.Errorf("Wechat Service: makeTextResponseBody error: %v", err)
	}
	return string(data), err
}

type GitHubBind interface {
	Add(GitHubBindData) error
	Update(GitHubBindData) error
	Remove(string)
	Exists(string) bool
	Find(string) *GitHubBindData
	Count() int
}

type GitHubBindData struct {
	WeChatID string
	GitHubID string
}

type GitHubBindDataList []GitHubBindData

type GitHubBinder struct {
	File     string
	DataList GitHubBindDataList
}

func (g *GitHubBinder) Read() (err error) {
	if _, err = os.Stat(g.File); os.IsNotExist(err) {
		g.DataList = GitHubBindDataList{}
		return nil
	}

	var content []byte
	if content, err = ioutil.ReadFile(g.File); err != nil {
		return
	}

	g.DataList = GitHubBindDataList{}
	err = yaml.Unmarshal(content, &g.DataList)
	return
}

func (g *GitHubBinder) Add(bindData GitHubBindData) (err error) {
	g.Read()

	if g.Exists(bindData.WeChatID) {
		return
	}

	g.DataList = append(g.DataList, bindData)
	var data []byte
	if data, err = yaml.Marshal(g.DataList); err == nil {
		err = ioutil.WriteFile(g.File, data, 0644)
	}
	return
}

func (g *GitHubBinder) Update(bindData GitHubBindData) (err error) {
	return
}

func (g *GitHubBinder) Remove(wechatID string) {
}

func (g *GitHubBinder) Exists(wechatID string) bool {
	g.Read()
	if g.DataList == nil {
		return false
	}

	for _, item := range g.DataList {
		if item.WeChatID == wechatID {
			return true
		}
	}
	return false
}

func (g *GitHubBinder) Find(wechatID string) *GitHubBindData {
	g.Read()
	if g.DataList == nil {
		return nil
	}

	for _, item := range g.DataList {
		if item.WeChatID == wechatID {
			return &item
		}
	}
	return nil
}

func (g *GitHubBinder) Count() int {
	g.Read()
	return len(g.DataList)
}

func init() {
	Register(func() AutoReply {
		return &GitHubBindAutoReply{
			GitHubBind: &GitHubBinder{
				File: "config/github_bind.yaml",
			},
		}
	})
}
