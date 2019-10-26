package config

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

// WeChatConfigurator represent the spec for configuration reader
type WeChatConfigurator interface {
	LoadConfig(string) (*WeChatConfig, error)
	GetConfig() *WeChatConfig
	SaveConfig() error
}
