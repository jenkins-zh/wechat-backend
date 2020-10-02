[![Docker Pulls](https://img.shields.io/docker/pulls/jenkinszh/wechat-backend.svg)](https://hub.docker.com/r/jenkinszh/wechat-backend/tags)
[![HitCount](http://hits.dwyl.com/jenkins-zh/wechat-backend.svg)](http://hits.dwyl.com/jenkins-zh/wechat-backend)

# WeChat Backend
As a robot, I can take care of some simple works for you.

# Features
* Welcome new members
* Auto Replay as Code
* Update backend configuration via WebHook (GitHub)

# Auto-reply
This robot will auto replay by configured files. You should put those config files into a fixed directory: `/management/auto-reply`.
It will reply a fixed sentence if there's no matched word. But you can give it a config file which contains the keyword `unknown`.

The structure of the config file like blew:
```
keyword: join
msgType: text
content: There is the sentence which will be replyed
```

See some examples from [here](https://github.com/jenkins-zh/wechat-bot-config/tree/master/management/auto-reply).

# Docker
One simple command could bring the WeChat backend up:

`docker run -t -p 12345:8080 -v /var/wechat/config:/config jenkinszh/wechat-backend:v0.0.1`

Sample config.yaml:

```
token: wechat-token
git_url: https://github.com/jenkins-zh/wechat
git_branch: master
github_webhook_secret: github-secret
appID: wechat-appid
appSecret: wechat-appsecret
server_port: 8080
```
