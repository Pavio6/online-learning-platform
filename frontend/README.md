# 在线学习平台前端

## 技术栈

- Vue 3
- Vite
- Element Plus
- Vue Router
- Pinia
- Axios

## 项目结构

本项目包含两个独立的前端应用：
- **学生端**：端口 3000
- **教师端**：端口 3001

## 启动方式

### 启动学生端
```bash
npm run dev:student
```
访问：http://localhost:3000

### 启动教师端
```bash
npm run dev:teacher
```
访问：http://localhost:3001

### 同时启动两个服务

需要打开两个终端窗口，分别运行：
```bash
# 终端1
npm run dev:student

# 终端2
npm run dev:teacher
```

## 环境配置

确保后端服务运行在 `http://localhost:8080`

环境变量配置在 `.env.development` 文件中。

## 功能说明

### 学生端
- ✅ 学生注册（选择校区）
- ✅ 学生登录
- ⏳ 课程列表/详情
- ⏳ 报名课程
- ⏳ 提交作业
- ⏳ 查看评论

### 教师端
- ✅ 教师登录
- ⏳ 课程管理
- ⏳ 创建课程/章节/课时
- ⏳ 批改作业
