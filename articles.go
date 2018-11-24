package main

import (
	"log"
	"os"

	"gopkg.in/src-d/go-git.v4"
)

func initCheck() {
	var err error

	_, err = os.Stat("wechat")
	if err != nil {
		if os.IsNotExist(err) {
			_, err = git.PlainClone("wechat", false, &git.CloneOptions{
				URL:      "https://github.com/jenkins-infra/wechat",
				Progress: os.Stdout,
			})
			if err != nil {
				log.Println(err)
				return
			}
		}
	} else {
		r, err := git.PlainOpen("wechat")
		if err == nil {
			w, err := r.Worktree()
			if err == nil {
				err = w.Pull(&git.PullOptions{
					RemoteName: "origin",
				})
				if err != nil {
					log.Println(err)
					return
				}
			}
		}
	}
}

func update() {

}

func getWelcome() string {
	return ""
}

func getKeywords() map[string]string {
	return nil
}
