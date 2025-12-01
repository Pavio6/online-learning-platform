# K6 负载测试

本目录包含使用 k6 进行负载测试的脚本和配置。

## 安装 k6

### macOS
```bash
brew install k6
```

### Linux
```bash
sudo gpg -k
sudo gpg --no-default-keyring --keyring /usr/share/keyrings/k6-archive-keyring.gpg --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys C5AD17C747E3415A3642D57D77C6C491D6AC1D9
echo "deb [signed-by=/usr/share/keyrings/k6-archive-keyring.gpg] https://dl.k6.io/deb stable main" | sudo tee /etc/apt/sources.list.d/k6.list
sudo apt-get update
sudo apt-get install k6
```

### Windows
下载安装包：https://github.com/grafana/k6/releases

## 测试场景

测试脚本包含以下测试场景：

1. **获取校区列表** - 测试基础查询性能
2. **获取课程列表** - 测试列表查询性能
3. **获取课程详情** - 测试详情查询性能
4. **获取课程评论** - 测试跨分片查询性能（分布式特性）
5. **学生登录** - 测试认证性能
6. **带认证的课程详情** - 测试认证后的查询性能

## 运行测试

### 基础测试
```bash
k6 run k6/load_test.js
```

### 指定服务器地址
```bash
BASE_URL=http://localhost:8080 k6 run k6/load_test.js
```

### 自定义负载配置
修改 `load_test.js` 中的 `options.stages` 来调整负载模式：

```javascript
export const options = {
  stages: [
    { duration: '30s', target: 10 },   // 30秒内增加到10个并发
    { duration: '1m', target: 50 },     // 1分钟内增加到50个并发
    { duration: '2m', target: 100 },    // 2分钟内增加到100个并发
    { duration: '1m', target: 50 },    // 1分钟内减少到50个并发
    { duration: '30s', target: 0 },     // 30秒内减少到0个并发
  ],
};
```

### 输出结果到文件
```bash
k6 run --out json=results.json k6/load_test.js
```

### 生成HTML报告
```bash
k6 run --out json=results.json k6/load_test.js
# 然后使用 k6-to-influxdb 或其他工具生成报告
```

## 性能指标

测试会监控以下指标：

- **HTTP请求持续时间** - 请求响应时间
- **HTTP请求失败率** - 请求失败的比例
- **自定义错误率** - 业务逻辑错误率
- **登录耗时** - 登录操作的响应时间
- **课程列表耗时** - 获取课程列表的响应时间
- **课程详情耗时** - 获取课程详情的响应时间
- **评论列表耗时** - 获取评论列表的响应时间（跨分片查询）

## 性能阈值

默认性能阈值：

- 95%的请求在500ms内完成
- 99%的请求在1s内完成
- 错误率小于1%

可以在 `load_test.js` 的 `options.thresholds` 中修改这些阈值。

## 测试数据准备

在运行测试前，确保：

1. 后端服务已启动并运行在指定端口（默认 localhost:8080）
2. 数据库中已有测试数据：
   - 至少有一个课程
   - 至少有一个学生账号（email: user01@test.com, password: user01@test.com）
   - 至少有一个教师账号（email: teacher@example.com, password: teacher123）

## 测试结果解读

测试完成后，k6 会输出详细的统计信息：

- **http_req_duration**: HTTP请求持续时间统计
- **http_req_failed**: HTTP请求失败率
- **errors**: 自定义错误率
- **login_duration**: 登录操作耗时
- **course_list_duration**: 课程列表查询耗时
- **course_detail_duration**: 课程详情查询耗时
- **comment_list_duration**: 评论列表查询耗时（跨分片）

## 注意事项

1. 测试前确保服务器有足够的资源处理负载
2. 建议在测试环境中运行，避免影响生产环境
3. 可以根据实际情况调整并发用户数和测试持续时间
4. 跨分片查询（如评论列表）的性能可能受网络延迟影响

