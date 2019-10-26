package config

import (
	"io/ioutil"
	"log"

	yaml "gopkg.in/yaml.v2"
)

// LocalFileConfig implement the WeChatConfigurator by local file system
type LocalFileConfig struct {
	path   string
	config *WeChatConfig
}

// LoadConfig load config from the file
func (l *LocalFileConfig) LoadConfig(configFile string) (
	cfg *WeChatConfig, err error) {
	var content []byte
	content, err = ioutil.ReadFile(configFile)
	if err != nil {
		log.Printf("load config file [%s] error: %v\n", configFile, err)
		return
	}

	cfg = &WeChatConfig{}
	l.config = cfg
	l.path = configFile
	err = yaml.Unmarshal(content, cfg)
	if err != nil {
		log.Printf("parse config file error: %v\n", err)
	}
	return
}

// GetConfig just return exists config object
func (l *LocalFileConfig) GetConfig() *WeChatConfig {
	return l.config
}

// SaveConfig save the config into local file
func (l *LocalFileConfig) SaveConfig() (err error) {
	if l.config == nil {
		return
	}

	var data []byte
	data, err = yaml.Marshal(l.config)
	if err == nil {
		err = ioutil.WriteFile(l.path, data, 0644)
	}
	return
}

var _ WeChatConfigurator = &LocalFileConfig{}
