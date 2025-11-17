# 配置说明

项目的所有配置均通过 `config.yaml`（或命令行指定的配置文件）进行管理。下面对各配置项进行说明。

## 1. app

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `name` | string | 应用名称 |
| `port` | int | HTTP服务监听端口 |
| `env` | string | 运行环境，`development` / `production` |
| `log_level` | string | 日志级别，支持 `debug`、`info`、`warn`、`error` |

## 2. jwt

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `secret` | string | JWT签名密钥（生产环境请更换） |
| `expiration` | duration | Token有效期，例如 `24h` |

## 3. database（中央服务器）

中央服务器保存所有“只读”或“需要集中管理”的数据，例如课程、章节、课程、任务。所有分支节点都会按日从中央服务器复制这些表，因此中央库必须先部署好并保证对外提供只读访问。

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `host` | string | 数据库主机 |
| `port` | int | 数据库端口 |
| `user` | string | 数据库用户名 |
| `password` | string | 数据库密码 |
| `dbname` | string | 数据库名称 |
| `sslmode` | string | SSL模式，开发环境可设为 `disable` |
| `max_open_conns` | int | 最大连接数 |
| `max_idle_conns` | int | 最大空闲连接数 |
| `conn_max_lifetime` | duration | 连接最大生命周期，例如 `300s` |

## 4. branches（分支节点）

每个分支节点对应一个完全独立的数据库实例，用来保存本地可写数据（用户信息、作业、评论、学习进度等）。在配置文件中以数组的形式列出所有分支。服务启动时会依次连接这些数据库，并按 `branch_id` 进行分片路由。

> 如果后期新增分支，只需在 `config.yaml` 中追加一条配置并重启服务即可。

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `branch_id` | int | 分支ID，对应于业务中的校区/节点 |
| `name` | string | 分支名称 |
| `host` / `port` / `user` / `password` / `dbname` / `sslmode` | 同中央数据库 |
| `max_open_conns` / `max_idle_conns` / `conn_max_lifetime` | 同中央数据库 |

> 注意：分支配置使用 `squash` 标签注入到 `DBSettings` 中，因此字段与中央库保持一致。

## 5. oss

阿里云OSS配置：

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `endpoint` | string | OSS服务Endpoint，例如 `oss-cn-hangzhou.aliyuncs.com` |
| `access_key_id` | string | OSS访问密钥ID |
| `access_key_secret` | string | OSS访问密钥Secret |
| `bucket_name` | string | 需要操作的bucket名称 |
| `region` | string | 区域信息（可选，用于标注部署区域） |

初始化成功后，可使用 `internal/oss` 包提供的接口上传课程视频、作业图片、生成签名URL等。

## 6. sync

数据同步配置：

### replication

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `enabled` | bool | 是否启用中央服务器到分支节点的复制 |
| `schedule` | string | Cron表达式，控制同步频率 |
| `tables` | []string | 需要复制的表列表，如 `courses`、`chapters` 等 |

### consolidation

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `enabled` | bool | 是否启用分支节点到中央服务器的数据整合 |
| `schedule` | string | Cron表达式，控制整合任务执行时间 |

---

### 使用步骤

1. 复制模板：`cp config.yaml.example config.yaml`
2. 根据部署环境修改数据库、OSS、同步等配置
3. 启动服务：`go run cmd/server/main.go -config config.yaml`

