# Swagger API 文档使用指南

## 访问Swagger UI

启动服务后，在浏览器中访问：

```
http://localhost:8080/swagger/index.html
```

## 使用步骤

### 1. 启动服务

```bash
go run cmd/server/main.go -config config.yaml
```

### 2. 打开Swagger UI

在浏览器中打开：`http://localhost:8080/swagger/index.html`

### 3. 测试API

#### 步骤1：教师登录
1. 找到 `教师认证` 分组
2. 点击 `POST /api/v1/teacher/auth/login`
3. 点击 "Try it out"
4. 输入测试数据：
```json
{
  "email": "teacher1@example.com",
  "password": "teacher123"
}
```
5. 点击 "Execute"
6. 复制返回的 `token` 值

#### 步骤2：设置认证
1. 点击页面右上角的 "Authorize" 按钮
2. 在弹出框的 "Value" 输入框中，**必须输入完整的Bearer token格式**：
   - 格式：`Bearer YOUR_TOKEN`
   - 例如：如果登录返回的token是 `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...`
   - 则输入：`Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...`
   - **注意**：`Bearer` 和 token 之间有一个空格
3. 点击 "Authorize"
4. 点击 "Close"

**重要提示**：由于使用的是Swagger 2.0格式，需要手动输入 `Bearer ` 前缀，Swagger UI不会自动添加。

#### 步骤3：测试需要认证的API
现在可以测试所有需要认证的接口了，例如：
- 创建课程
- 创建章节
- 创建课程
- 创建任务
等等

## 更新Swagger文档

如果修改了API注释，需要重新生成文档：

```bash
swag init -g cmd/server/main.go -o docs
```

## 注意事项

1. 确保在数据库中已插入测试教师账号
2. Token有效期为24小时（可在config.yaml中配置）
3. 所有需要认证的接口都需要先设置Bearer Token

