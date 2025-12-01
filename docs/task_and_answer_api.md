image.png# 任务管理和作业提交服务 API 文档

## 概述

本文档说明已实现的任务管理服务（中央服务器）和作业提交服务（分支节点）的所有功能。

## 步骤5.3：任务管理服务（中央服务器）

### ✅ 1. 教师创建任务

**接口：** `POST /api/v1/teacher/lessons/{id}/tasks`

**认证：** 需要 Bearer Token（教师角色）

**请求体：**
```json
{
  "task_title": "任务标题",
  "description": "任务描述",
  "task_type": "essay",  // essay, quiz, upload
  "max_score": 100
}
```

**响应：**
```json
{
  "task_id": 1,
  "lesson_id": 1,
  "task_title": "任务标题",
  "description": "任务描述",
  "task_type": "essay",
  "max_score": 100,
  "created_at": "2025-01-01 12:00:00",
  "updated_at": "2025-01-01 12:00:00"
}
```

**功能说明：**
- 验证课程是否存在且属于该教师
- 支持三种任务类型：essay（作文）、quiz（测验）、upload（上传）
- 默认任务类型为 essay，默认最高分为 100

### ✅ 2. 任务查询（学生和教师）

#### 2.1 获取任务详情

**学生接口：** `GET /api/v1/student/tasks/{id}`  
**教师接口：** `GET /api/v1/teacher/tasks/{id}`

**认证：**
- 学生接口：不需要认证
- 教师接口：需要 Bearer Token（教师角色）

**响应：**
```json
{
  "task_id": 1,
  "lesson_id": 1,
  "task_title": "任务标题",
  "description": "任务描述",
  "task_type": "essay",
  "max_score": 100,
  "created_at": "2025-01-01 12:00:00",
  "updated_at": "2025-01-01 12:00:00"
}
```

#### 2.2 获取课程的所有任务

**学生接口：** `GET /api/v1/student/courses/{id}/tasks`  
**教师接口：** `GET /api/v1/teacher/courses/{id}/tasks`

**认证：**
- 学生接口：不需要认证
- 教师接口：需要 Bearer Token（教师角色）

**响应：**
```json
[
  {
    "task_id": 1,
    "lesson_id": 1,
    "task_title": "任务标题",
    "description": "任务描述",
    "task_type": "essay",
    "max_score": 100,
    "created_at": "2025-01-01 12:00:00",
    "updated_at": "2025-01-01 12:00:00"
  }
]
```

## 步骤5.4：作业提交服务（分支节点）

### ✅ 1. 学生提交作业（上传图片到OSS，分片路由）

**接口：** `POST /api/v1/student/tasks/{id}/answers`

**认证：** 需要 Bearer Token（学生角色）

**请求格式：** `multipart/form-data`

**请求参数：**
- `answer_content` (string, 可选): 作业文本内容
- `type` (string, 可选): 作业类型（text/image_url），默认为 text
- `file` (file, 可选): 作业图片文件

**功能说明：**
- 支持文本提交和图片上传两种方式
- 如果上传图片，会自动上传到 OSS（阿里云对象存储）
- OSS 路径格式：`answers/{branch_id}/{user_id}/{task_id}/{timestamp}_{filename}`
- 如果数据库操作失败，会自动删除已上传的 OSS 文件（回滚机制）
- 如果学生已提交过作业，会更新现有作业而不是创建新记录
- 分片路由：根据 `branch_id` 和 `user_id` 路由到对应的分支数据库

**响应：**
```json
{
  "answer_id": 1,
  "task_id": 1,
  "branch_id": 1,
  "user_id": 1,
  "answer_content": "作业内容或OSS图片URL",
  "type": "text",  // 或 "image_url"
  "score": 0,
  "is_graded": false,
  "graded_by": null,
  "submitted_at": "2025-01-01 12:00:00",
  "created_at": "2025-01-01 12:00:00",
  "updated_at": "2025-01-01 12:00:00"
}
```

### ✅ 2. 作业查询（学生和教师）

#### 2.1 学生查询自己的作业

**接口：** `GET /api/v1/student/tasks/{id}/answers`

**认证：** 需要 Bearer Token（学生角色）

**功能说明：**
- 查询当前学生在指定任务下的作业
- 如果未提交作业，返回 404 错误

**响应：**
```json
{
  "answer_id": 1,
  "task_id": 1,
  "branch_id": 1,
  "user_id": 1,
  "answer_content": "作业内容或OSS图片URL",
  "type": "text",
  "score": 85,
  "is_graded": true,
  "graded_by": 2,
  "submitted_at": "2025-01-01 12:00:00",
  "created_at": "2025-01-01 12:00:00",
  "updated_at": "2025-01-01 12:00:00"
}
```

#### 2.2 教师查看任务的所有作业

**接口：** `GET /api/v1/teacher/tasks/{id}/answers`

**认证：** 需要 Bearer Token（教师角色）

**功能说明：**
- 教师只能查看自己分支的学生作业
- 验证任务是否属于该教师
- 按提交时间倒序排列

**响应：**
```json
[
  {
    "answer_id": 1,
    "task_id": 1,
    "branch_id": 1,
    "user_id": 1,
    "answer_content": "作业内容或OSS图片URL",
    "type": "text",
    "score": 85,
    "is_graded": true,
    "graded_by": 2,
    "submitted_at": "2025-01-01 12:00:00",
    "created_at": "2025-01-01 12:00:00",
    "updated_at": "2025-01-01 12:00:00"
  }
]
```

### ✅ 3. 教师评分

**接口：** `PUT /api/v1/teacher/answers/{id}/grade`

**认证：** 需要 Bearer Token（教师角色）

**请求体：**
```json
{
  "score": 85
}
```

**功能说明：**
- 教师只能对自己分支的作业进行评分
- 评分后自动设置 `is_graded = true` 和 `graded_by = 教师ID`

**响应：**
```json
{
  "answer_id": 1,
  "task_id": 1,
  "branch_id": 1,
  "user_id": 1,
  "answer_content": "作业内容",
  "type": "text",
  "score": 85,
  "is_graded": true,
  "graded_by": 2,
  "submitted_at": "2025-01-01 12:00:00",
  "created_at": "2025-01-01 12:00:00",
  "updated_at": "2025-01-01 12:00:00"
}
```

## 测试步骤

### 1. 教师创建任务

1. 使用教师账号登录，获取 Bearer Token
2. 调用 `POST /api/v1/teacher/lessons/{lesson_id}/tasks`
3. 传入任务信息（标题、描述、类型、最高分）
4. 验证返回的任务信息

### 2. 查询任务

1. **学生查询任务详情：**
   - 调用 `GET /api/v1/student/tasks/{task_id}`
   - 验证返回的任务信息

2. **教师查询任务详情：**
   - 使用教师 Token 调用 `GET /api/v1/teacher/tasks/{task_id}`
   - 验证返回的任务信息

3. **查询课程的所有任务：**
   - 学生：`GET /api/v1/student/courses/{course_id}/tasks`
   - 教师：`GET /api/v1/teacher/courses/{course_id}/tasks`

### 3. 学生提交作业

1. 使用学生账号登录，获取 Bearer Token
2. **提交文本作业：**
   - 调用 `POST /api/v1/student/tasks/{task_id}/answers`
   - Content-Type: `multipart/form-data`
   - 传入 `answer_content` 参数
   - 验证返回的作业信息

3. **提交图片作业：**
   - 调用 `POST /api/v1/student/tasks/{task_id}/answers`
   - Content-Type: `multipart/form-data`
   - 传入 `file` 参数（图片文件）
   - 验证返回的作业信息，确认 `type` 为 `image_url`，`answer_content` 为 OSS URL

4. **更新已提交的作业：**
   - 再次调用提交接口
   - 验证作业被更新而不是创建新记录

### 4. 查询作业

1. **学生查询自己的作业：**
   - 使用学生 Token 调用 `GET /api/v1/student/tasks/{task_id}/answers`
   - 验证返回的作业信息

2. **教师查看任务的所有作业：**
   - 使用教师 Token 调用 `GET /api/v1/teacher/tasks/{task_id}/answers`
   - 验证返回的作业列表（仅限自己分支的学生）

### 5. 教师评分

1. 使用教师账号登录，获取 Bearer Token
2. 调用 `PUT /api/v1/teacher/answers/{answer_id}/grade`
3. 传入分数（0-100）
4. 验证返回的作业信息，确认 `is_graded = true`，`graded_by` 为教师ID

## 错误处理

### 常见错误码

- `3004`: 任务不存在
- `5001`: 作业不存在
- `3005`: 不是课程教师（无权操作）
- `1001`: 参数错误
- `1003`: 未授权（Token 无效或过期）
- `1004`: 禁止访问（权限不足）

## 技术实现细节

### 分片路由

- 学生提交作业时，根据 `branch_id` 和 `user_id` 路由到对应的分支数据库
- 教师查看作业时，只能查看自己分支（`branch_id`）的学生作业

### OSS 上传

- 图片上传到阿里云 OSS
- 路径格式：`answers/{branch_id}/{user_id}/{task_id}/{timestamp}_{filename}`
- 如果数据库操作失败，自动删除已上传的文件（回滚机制）

### 数据存储

- **任务（Tasks）**：存储在中央数据库（central）
- **作业（Answers）**：存储在分支数据库（branch），按 `branch_id` 和 `user_id` 分片

## Swagger 文档

所有接口已在 Swagger 文档中注册，可以通过以下方式访问：

1. 启动后端服务
2. 访问 `http://localhost:8080/swagger/index.html`
3. 查看 "教师任务管理" 和 "学生任务" 标签下的所有接口

