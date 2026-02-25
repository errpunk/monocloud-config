# monocloud-config

## 背景

monocloud 是一个代理服务器机场，它提供了 clash 的配置文件订阅链接。

它的订阅连接下载的内容包括了：

```yaml
port: 7890
socks-port: 7891
allow-lan: false
mode: rule
log-level: info
external-controller: :9090

# 节点列表
proxies:

# 规则组
proxy-groups:
  - name: Proxy
    type: select
    proxies:
      - 🇭🇰 Relay-HK1
      - 🇭🇰 Relay-HK2
      - 🇭🇰 Relay-HK3
      - 🇭🇰 Relay-HK4    

# 规则
rules:
  - 'DOMAIN-KEYWORD,google,Proxy'

```

而我需要它的：
- proxies 节点列表
- proxy-groups 规则组
- rules 规则


## 目标

将 monocloud 的订阅链接下载的内容，提取出 proxies、proxy-groups、rules，然后合并到我本地的 config.yaml 中。


## 步骤

1. 从环境变量中获取 monocloud 的订阅链接。
2. 下载 monocloud 的订阅链接，得到一个 yaml 文件。
3. 从 yaml 文件中提取出 proxies、proxy-groups、rules。
4. 合并到我本地的 config.yaml 中。
5. 验证 config.yaml 是否正确。

## 环境变量

- MONOCLOUD_URL: monocloud 的订阅链接

## 注意

从订阅链接下载规则的时候，要设置 user-agent 为 clash.meta，否则会下载失败。
