# WeChat Backend

As a robot, I can take care of some simple works for you.

# Fetures

* Welcome new members
* Auto Replay as Code
* WebHook (GitHub)

# Docker

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

One simple command could bring the Jenkins wechat backend up:

`docker run -t -p 12345:8080 -v /var/wechat/config:/config surenpi/jenkins-wechat`

docker run -t -p 45678:18080 -v /var/wechat/config:/config surenpi/jenkins-wechat:dev--LinuxSuRen-gmail.com

## API

Below are the APIs of WeChat bot:

| Path | Description |
|---|---|
| GET `/medias` | Returns the media file list |
