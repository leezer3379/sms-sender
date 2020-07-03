# sms-sender

Nightingale的理念，是将告警事件扔到redis里就不管了，接下来由各种sender来读取redis里的事件并发送，毕竟发送报警的方式太多了，适配起来比较费劲，希望社区同仁能够共建。

这里提供一个钉钉的sender，参考了[https://github.com/n9e/wechat-sender](https://github.com/n9e/wechat-sender) 及 [https://github.com/wulorn/dingtalk](https://github.com/wulorn/dingtalk)，具体如何获取钉钉机器人token，也可以参看钉钉官网

## compile

```bash
cd $GOPATH/src
mkdir -p github.com/n9e
cd github.com/n9e
git clone https://github.com/leezer3379/sms-sender.git
cd sms-sender
./control build
```

如上编译完就可以拿到二进制了。

## configuration

直接修改etc/sms-sender.yml即可

## 注意

sms-sender仅支持通过http请求调用自己封装的短信接口，即传递数据为body数据如下：
```bash
{"message":"test message from n9e at 2020-07-03 15:51:03.778382 +0800 CST m=+0.001693823","mobile":"10086"}

成功返回 ok字符串即可

```

需monapi.yaml设置里的notify添加sms告警，如下：

```yaml
notify:
  p1: ["sms"]
  p2: ["sms"]
  p3: ["sms"]
```

在web后台添加用户 (此用户作为报警群虚拟用户，仅用来指定对应的钉钉群告警

## pack

编译完成之后可以打个包扔到线上去跑，将二进制和配置文件打包即可：

```bash
tar zcvf sms-sender.tar.gz sms-sender etc/sms-sender.yml etc/sms.tpl
```

## sms-sender.yml

```yaml
sms:
  openurl: "http://127.0.0.1:8001/api-develop/test/"  # 配置自己的短信接口
```
## test

配置etc/sms-sender.yml，相关配置修改好，我们先来测试一下是否好使， `./sms-sender -p phone`，token为钉钉群机器人的token值，程序会自动读取etc目录下的配置文件，发一个测试消息给钉钉群`token`

## run

如果测试发送没问题，扔到线上跑吧，使用systemd或者supervisor之类的托管起来，systemd的配置实例：


```
$ cat sms-sender.service
[Unit]
Description=Nightingale sms sender
After=network-online.target
Wants=network-online.target

[Service]
User=root
Group=root

Type=simple
ExecStart=/home/n9e/sms-sender
WorkingDirectory=/home/n9e

Restart=always
RestartSec=1
StartLimitInterval=0

[Install]
WantedBy=multi-user.target
```