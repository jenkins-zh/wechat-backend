package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

// WeChatConfig represents WeChat config
type WeChatConfig struct {
	GitURL              string `yaml:"git_url"`
	GitBranch           string `yaml:"git_branch"`
	GitHubWebHookSecret string `yaml:"github_webhook_secret"`

	ServerPort int `yaml:"server_port"`

	AppID     string `yaml:"appID"`
	AppSecret string `yaml:"appSecret"`
	Token     string `yaml:"token"`
}

// LoadConfig load config
func LoadConfig(configFile string) (config *WeChatConfig, err error) {
	var content []byte
	content, err = ioutil.ReadFile(configFile)
	if err != nil {
		log.Printf("load config file [%s] error: %v\n", configFile, err)
		return
	}

	config = &WeChatConfig{}
	err = yaml.Unmarshal(content, config)
	if err != nil {
		log.Printf("parse config file error: %v\n", err)
	}
	return
}
