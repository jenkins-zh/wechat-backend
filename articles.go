package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"strings"

	"github.com/linuxsuren/wechat-backend/pkg/config"
	"github.com/linuxsuren/wechat-backend/pkg/reply"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/yaml.v2"
)

const (
	CONFIG = "wechat"
)

func initCheck(weConfig *config.WeChatConfig) {
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
					initCheck(weConfig)
				}
			} else {
				log.Println("open dir with git error", err)
				os.Remove(CONFIG)
				initCheck(weConfig)
			}
		}
	} else {
		log.Println("can't get config dir status", err)

		if os.RemoveAll(CONFIG) == nil {
			initCheck(weConfig)
		}
	}

	if err == nil {
		log.Println("going to update the cache.")
		update()
	}
}

var respMap = make(map[string]interface{})

func responseHandler(yamlContent []byte) {
	reps := reply.ResponseBody{}
	err := yaml.Unmarshal(yamlContent, &reps)
	if err == nil {
		log.Println(reps.MsgType, reps.Keyword, reps)
		// reps.MsgType = reps.Kind

		switch reps.MsgType {
		case "text":
			text := reply.TextResponseBody{}
			yaml.Unmarshal(yamlContent, &text)
			respMap[reps.Keyword] = text
		case "image":
			image := reply.ImageResponseBody{}
			yaml.Unmarshal(yamlContent, &image)
			respMap[reps.Keyword] = image
		case "news":
			news := reply.NewsResponseBody{}
			yaml.Unmarshal(yamlContent, &news)
			respMap[reps.Keyword] = news
		default:
			log.Println("unknow type", reps.MsgType)
		}
	} else {
		fmt.Println(err)
	}
}

func update() {
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
			responseHandler(content)
		} else {
			log.Println("Can't read file ", file.Name())
		}
	}
}

func getKeywords() map[string]string {
	return nil
}
