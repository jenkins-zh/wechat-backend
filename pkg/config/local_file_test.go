package config

import "testing"

func TestLoadConfig(t *testing.T) {
	cfg := &LocalFileConfig{}

	config, err := cfg.LoadConfig("config.yaml")
	if err != nil {
		t.Errorf("load config error %v", err)
	}

	if config.Token != "Token" || config.GitURL != "GitURL" ||
		config.GitBranch != "GitBranch" ||
		config.GitHubWebHookSecret != "GitHubWebHookSecret" ||
		config.ServerPort != 80 ||
		config.AppID != "appID" ||
		config.AppSecret != "appSecret" {
		t.Errorf("parse error, config %#v", config)
	}
}
