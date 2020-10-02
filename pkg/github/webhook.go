package github

import (
	"log"
	"net/http"

	"github.com/jenkins-zh/wechat-backend/pkg/config"
	"gopkg.in/go-playground/webhooks.v5/github"
)

func WebhookHandler(w http.ResponseWriter, r *http.Request, weConfig *config.WeChatConfig, initCheck func(*config.WeChatConfig)) {
	hook, _ := github.New(github.Options.Secret("secret"))

	payload, err := hook.Parse(r, github.PushEvent)
	if err != nil {
		if err == github.ErrEventNotFound {
			// ok event wasn;t one of the ones asked to be parsed
			return
		}
	}

	switch payload.(type) {
	case github.PushPayload:
		push := payload.(github.PushPayload)
		// Do whatever you want from here...
		log.Printf("push ref is %s.\n", push.Ref)

		if push.Ref != "refs/heads/master" {
			return
		}
	}

	log.Println("Going to update wechat config.")
	initCheck(weConfig)
}
