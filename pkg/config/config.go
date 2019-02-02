package config

import (
	"io/ioutil"
	"log"

	core "github.com/linuxsuren/wechat-backend/pkg"
	yaml "gopkg.in/yaml.v2"
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

	Valid bool `yaml:"valid"`
}

// NewConfig new config instance
func NewConfig() *WeChatConfig {
	return &WeChatConfig{}
}

// LoadConfig load config
func LoadConfig(configFile string) (config *WeChatConfig, err error) {
	var content []byte
	content, err = ioutil.ReadFile(configFile)
	if err != nil {
		log.Printf("load config file [%s] error: %v\n", configFile, err)
		return
	}

	config = NewConfig()
	err = yaml.Unmarshal(content, config)
	if err != nil {
		log.Printf("parse config file error: %v\n", err)
	}
	return
}

// SettingConfig set the WeChat config
func (wc *WeChatConfig) SettingConfig() (err error) {
	var data []byte
	data, err = yaml.Marshal(wc)
	if err == nil {
		err = ioutil.WriteFile(core.ConfigPath, data, 0644)
	}
	return
}
