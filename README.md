# 分布式在线学习平台

这是一个基于Go语言开发的分布式在线学习平台，作为分布式数据库课程作业项目。

## 项目特性

- **分布式架构**：采用分片架构，中央服务器存储课程数据，分支节点存储用户数据
- **数据同步**：实现РОК（主从复制）和РОК+КД（主从复制+数据整合）机制
- **分片路由**：支持按branch_id和user_id进行数据路由
- **跨分片查询**：支持跨分片聚合查询

## 技术栈

- **语言**：Go 1.21+
- **数据库**：PostgreSQL 14+
- **ORM**：GORM
- **API框架**：Gin
- **OSS**：阿里云OSS SDK
- **配置管理**：Viper
- **数据库迁移**：golang-migrate
- **日志**：logrus
- **JWT**：golang-jwt/jwt
- **定时任务**：robfig/cron

## 项目结构

```
online-learning-platform/
├── cmd/server/          # 应用入口
├── internal/
│   ├── models/          # 数据模型
│   ├── database/        # 数据库连接和分片路由
│   ├── service/         # 业务逻辑层
│   ├── api/             # API路由和处理器
│   │   ├── student/     # 学生端API
│   │   └── teacher/     # 教师端API
│   ├── config/          # 配置管理
│   └── oss/             # OSS客户端
├── migrations/          # 数据库迁移脚本
│   ├── central/         # 中央服务器迁移
│   └── branch/          # 分支节点迁移
├── docs/                # 文档
├── tests/               # 测试文件
└── scripts/             # 部署脚本
```

## 快速开始

### 1. 安装依赖

```bash
go mod download
```

### 2. 配置环境

复制配置文件模板：

```bash
cp config.yaml.example config.yaml
```

编辑 `config.yaml`，配置数据库、OSS等信息。

### 3. 数据库初始化

目前提供了两份聚合 SQL，直接在数据库里执行即可完成建表：

```bash
# 中央服务器（课程/章节/课程/任务）
psql -h <central-host> -U <user> -d learning_central -f migrations/central/central_schema.sql

# 分支节点（用户/作业/评论/学习进度）
psql -h <branch-host> -U <user> -d learning_branch1 -f migrations/branch/branch_schema.sql
```

> 如果有多个分支节点，只需在各自数据库上执行 `branch_schema.sql`。

### 4. 运行应用

```bash
go run cmd/server/main.go
```

## API文档

### 学生端API

- `POST /api/v1/student/auth/register` - 学生注册
- `POST /api/v1/student/auth/login` - 学生登录
- `GET /api/v1/student/courses` - 查看课程列表
- `POST /api/v1/student/courses/:id/enroll` - 报名课程
- `POST /api/v1/student/tasks/:id/answers` - 提交作业
- `POST /api/v1/student/courses/:id/comments` - 发表评论

### 教师端API

- `POST /api/v1/teacher/auth/login` - 教师登录
- `POST /api/v1/teacher/courses` - 创建课程
- `POST /api/v1/teacher/courses/:id/chapters` - 创建章节
- `POST /api/v1/teacher/courses/:id/chapters/:chapter_id/lessons` - 创建课程
- `POST /api/v1/teacher/lessons/:id/tasks` - 创建任务
- `PUT /api/v1/teacher/answers/:id/grade` - 评分作业

## 分片策略

- **中央服务器**：Courses, Chapters, Lessons, Tasks
- **分支节点**：Users, Answers, Comments, Learning, Branches

## 数据同步

- **РОК（主从复制）**：中央服务器 → 分支节点，每日同步课程数据
- **РОК+КД（主从复制+数据整合）**：分支节点 → 中央服务器，每日整合统计数据

## 开发计划

1. ✅ 项目初始化与架构设计
2. ⏳ 数据库层实现
3. ⏳ 核心业务服务实现
4. ⏳ API层实现
5. ⏳ 分布式特性实现
6. ⏳ 测试与优化

## License

MIT

