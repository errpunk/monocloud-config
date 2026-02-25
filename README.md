# monocloud-config

从 monocloud 订阅链接下载 Clash 配置，提取 `proxies`、`proxy-groups`、`rules`，自动合并到本地 `config.yaml`，并作为后台守护进程定时更新。

## 快速开始

### 使用 Docker Compose（推荐）

```bash
# 1. 复制环境变量模板
cp .env.example .env

# 2. 填写订阅链接
#    打开 .env，设置 MONOCLOUD_URL=https://your-subscription-link

# 3. 启动服务
docker compose up -d
```

`monocloud-config` 会与 `mihomo` 共享同一个配置目录，启动时立即更新，之后按 `UPDATE_INTERVAL` 定时刷新。

参考配置见 [config-reference/](config-reference/)。

---

### 手动运行（二进制）

从 [Releases](https://github.com/errpunk/monocloud-config/releases) 或 [GitHub Package](https://ghcr.io/errpunk/monocloud-config) 下载，也可自行编译：

```bash
make build          # 编译当前平台
make build-linux    # 交叉编译 linux/amd64 + arm64（输出到 bin/）
```

运行：

```bash
export MONOCLOUD_URL="https://your-subscription-link"
export CONFIG_PATH="/etc/mihomo/config.yaml"   # 默认 config.yaml
export UPDATE_INTERVAL="1h"                    # 默认 1h，支持 30m / 6h 等

./bin/monocloud-config          # 持续运行，定时更新
./bin/monocloud-config -v       # 打印版本号
```

---

## 环境变量

| 变量 | 必填 | 默认值 | 说明 |
|------|------|--------|------|
| `MONOCLOUD_URL` | ✅ | — | monocloud 订阅链接 |
| `CONFIG_PATH` | — | `config.yaml` | 本地 Clash 配置文件路径 |
| `UPDATE_INTERVAL` | — | `1h` | 更新间隔，Go duration 格式（如 `30m`、`6h`） |

## Docker 镜像

```bash
docker pull ghcr.io/errpunk/monocloud-config:latest
```

支持 `linux/amd64` 和 `linux/arm64`。

## 开发

```bash
make test          # 运行单元测试
make docker-build  # 构建本地镜像
make docker-push   # 构建并推送多平台镜像到 ghcr.io
```
