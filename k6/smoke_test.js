import http from 'k6/http';
import { check } from 'k6';

// 轻量级冒烟测试 - 验证系统基本功能
export const options = {
  vus: 1,        // 1个虚拟用户
  duration: '30s', // 持续30秒
};

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';

export default function () {
  // 健康检查
  const healthRes = http.get(`${BASE_URL}/health`);
  check(healthRes, {
    '健康检查通过': (r) => r.status === 200,
  });

  // 获取校区列表
  const branchesRes = http.get(`${BASE_URL}/api/v1/student/branches`);
  check(branchesRes, {
    '校区列表状态码为200': (r) => r.status === 200,
  });

  // 获取课程列表
  const coursesRes = http.get(`${BASE_URL}/api/v1/student/courses`);
  check(coursesRes, {
    '课程列表状态码为200': (r) => r.status === 200,
  });

  // 学生登录
  const loginPayload = JSON.stringify({
    email: 'user01@test.com',
    password: 'user01@test.com',
  });
  const loginRes = http.post(
    `${BASE_URL}/api/v1/student/auth/login`,
    loginPayload,
    { headers: { 'Content-Type': 'application/json' } }
  );
  check(loginRes, {
    '登录成功': (r) => r.status === 200,
    '返回token': (r) => {
      const body = JSON.parse(r.body);
      return body.token !== undefined;
    },
  });
}

