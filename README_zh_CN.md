# 微信公众号自动应答

作为一个机器人🤖️，我帮您做一些简单事情。

# 功能

* 欢迎👏新成员
* 自动回复即代码
* 支持通过 WebHook 自动更新配置（GitHub）

# Docker

一条简单的 Docker 命令就可以把微信公众号的自动应答程序运行起来：

`docker run -t -p 12345:8080 -v /var/wechat/config:/config surenpi/jenkins-wechat`

示例配置文件 config.yaml:

```
token: wechat-token
git_url: https://github.com/jenkins-zh/wechat
git_branch: master
github_webhook_secret: github-secret
appID: wechat-appid
appSecret: wechat-appsecret
server_port: 8080
```
