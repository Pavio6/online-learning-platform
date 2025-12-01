import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate } from 'k6/metrics';

const errorRate = new Rate('errors');
const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';

// 压力测试配置 - 快速增加到高负载
export const options = {
  stages: [
    { duration: '10s', target: 50 },   // 10秒内增加到50个并发
    { duration: '30s', target: 200 },  // 30秒内增加到200个并发
    { duration: '1m', target: 500 },   // 1分钟内增加到500个并发
    { duration: '30s', target: 200 },  // 30秒内减少到200个并发
    { duration: '10s', target: 0 },    // 10秒内减少到0个并发
  ],
  thresholds: {
    http_req_duration: ['p(95)<2000', 'p(99)<5000'], // 压力测试下放宽阈值
    http_req_failed: ['rate<0.05'],                  // 错误率小于5%
    errors: ['rate<0.05'],
  },
};

function studentLogin(email, password) {
  const url = `${BASE_URL}/api/v1/student/auth/login`;
  const payload = JSON.stringify({ email, password });
  const params = { headers: { 'Content-Type': 'application/json' } };
  
  const res = http.post(url, payload, params);
  const success = check(res, {
    '登录成功': (r) => r.status === 200,
  });
  
  if (!success) {
    errorRate.add(1);
  } else {
    errorRate.add(0);
  }
  
  return res.status === 200 ? JSON.parse(res.body).token : null;
}

function getCourseList() {
  const res = http.get(`${BASE_URL}/api/v1/student/courses`);
  const success = check(res, {
    '课程列表获取成功': (r) => r.status === 200,
  });
  
  if (!success) {
    errorRate.add(1);
  } else {
    errorRate.add(0);
  }
  
  return res.status === 200 ? JSON.parse(res.body) : [];
}

function getCourseDetail(courseId) {
  const res = http.get(`${BASE_URL}/api/v1/student/courses/${courseId}`);
  const success = check(res, {
    '课程详情获取成功': (r) => r.status === 200,
  });
  
  if (!success) {
    errorRate.add(1);
  } else {
    errorRate.add(0);
  }
}

export default function () {
  // 获取课程列表
  const courses = getCourseList();
  sleep(0.5);
  
  // 如果有课程，获取详情
  if (courses && courses.length > 0) {
    const randomCourse = courses[Math.floor(Math.random() * courses.length)];
    getCourseDetail(randomCourse.course_id);
    sleep(0.5);
  }
  
  // 登录（部分用户）
  if (Math.random() > 0.7) { // 30%的用户执行登录
    studentLogin('user01@test.com', 'user01@test.com');
    sleep(0.5);
  }
}

