# WeChat Backend
As a robot, I can take care of some simple works for you.

# Features
* Welcome new members
* Auto Replay as Code
* Update backend configuration via WebHook (GitHub)

# Keywords
This robot will auto replay by configured files. But you can give it a config file which contains the keyword `unknown`.

The structure of config file like blew:
```
keyword: join
msgType: text
content: There is the sentence which will be replyed
```

# Docker
One simple command could bring the Jenkins wechat backend up:

`docker run -t -p 12345:8080 -v /var/wechat/config:/config surenpi/jenkins-wechat`

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
