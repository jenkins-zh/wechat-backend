package article

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"strings"

	core "github.com/jenkins-zh/wechat-backend/pkg"
	"github.com/jenkins-zh/wechat-backend/pkg/config"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/yaml.v2"
)

const (
	CONFIG = "wechat"
)

type ResponseManager interface {
	GetResponse(string) (interface{}, bool)
	InitCheck(weConfig *config.WeChatConfig)
}

type DefaultResponseManager struct {
	ResponseMap map[string]interface{}
}

// NewDefaultResponseManager should always call this method to get a object
func NewDefaultResponseManager() (mgr *DefaultResponseManager) {
	mgr = &DefaultResponseManager{
		ResponseMap: make(map[string]interface{}, 10),
	}
	return
}

func (drm *DefaultResponseManager) GetResponse(keyword string) (interface{}, bool) {
	res, ok := drm.ResponseMap[keyword]
	return res, ok
}

func (drm *DefaultResponseManager) InitCheck(weConfig *config.WeChatConfig) {
	var err error

	_, err = os.Stat(CONFIG)
	if err != nil {
		if os.IsNotExist(err) {
			_, err = git.PlainClone(CONFIG, false, &git.CloneOptions{
				URL:      weConfig.GitURL,
				Progress: os.Stdout,
			})
			if err != nil {
				log.Println("clone failure", err)
				return
			}
			log.Println("the clone progress is done")
		} else {
			r, err := git.PlainOpen(CONFIG)
			if err == nil {
				w, err := r.Worktree()
				if err == nil {
					err = w.Pull(&git.PullOptions{
						RemoteName: "origin",
					})
					if err != nil {
						log.Println(err)
						// return
					}
				} else {
					log.Println("open work tree with git error", err)
					os.Remove(CONFIG)
					drm.InitCheck(weConfig)
				}
			} else {
				log.Println("open dir with git error", err)
				os.Remove(CONFIG)
				drm.InitCheck(weConfig)
			}
		}
	} else {
		log.Println("can't get config dir status", err)

		if os.RemoveAll(CONFIG) == nil {
			drm.InitCheck(weConfig)
		}
	}

	if err == nil {
		log.Println("going to update the cache.")
		drm.update()
	}
}

func (drm *DefaultResponseManager) responseHandler(yamlContent []byte) {
	reps := core.ResponseBody{}
	err := yaml.Unmarshal(yamlContent, &reps)
	if err == nil {
		log.Println(reps.MsgType, reps.Keyword, reps)

		switch reps.MsgType {
		case "text":
			text := core.TextResponseBody{}
			yaml.Unmarshal(yamlContent, &text)
			drm.ResponseMap[reps.Keyword] = text
		case "image":
			image := core.ImageResponseBody{}
			yaml.Unmarshal(yamlContent, &image)
			drm.ResponseMap[reps.Keyword] = image
		case "news":
			news := core.NewsResponseBody{}
			yaml.Unmarshal(yamlContent, &news)
			drm.ResponseMap[reps.Keyword] = news
		case "random": // TODO this not the regular way
			random := core.RandomResponseBody{}
			yaml.Unmarshal(yamlContent, &random)
			drm.ResponseMap[reps.Keyword] = random
		default:
			log.Println("unknow type", reps.MsgType)
		}
	} else {
		fmt.Println(err)
	}
}

func (drm *DefaultResponseManager) update() {
	root := CONFIG + "/management/auto-reply"
	files, err := ioutil.ReadDir(root)
	if err != nil {
		log.Fatal(err)
		return
	}
	for _, file := range files {
		if !strings.Contains(file.Name(), "keywords") {
			continue
		}

		content, err := ioutil.ReadFile(root + "/" + file.Name())
		if err == nil {
			drm.responseHandler(content)
		} else {
			log.Println("Can't read file ", file.Name())
		}
	}
}
