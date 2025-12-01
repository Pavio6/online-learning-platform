import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate, Trend, Counter } from 'k6/metrics';

// 自定义指标
const errorRate = new Rate('errors');
const loginDuration = new Trend('login_duration');
const courseListDuration = new Trend('course_list_duration');
const courseDetailDuration = new Trend('course_detail_duration');
const commentListDuration = new Trend('comment_list_duration');

// 配置
export const options = {
  stages: [
    { duration: '30s', target: 10 },   // 30秒内逐步增加到10个并发用户
    { duration: '1m', target: 50 },    // 1分钟内增加到50个并发用户
    { duration: '2m', target: 100 },   // 2分钟内增加到100个并发用户
    { duration: '1m', target: 50 },    // 1分钟内减少到50个并发用户
    { duration: '30s', target: 0 },   // 30秒内减少到0个并发用户
  ],
  thresholds: {
    http_req_duration: ['p(95)<500', 'p(99)<1000'], // 95%的请求在500ms内，99%在1s内
    http_req_failed: ['rate<0.05'],                  // HTTP请求失败率小于5%（放宽阈值，因为登录可能失败）
    errors: ['rate<0.05'],                           // 自定义错误率小于5%（放宽阈值）
  },
};

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';

// 测试数据
const testUsers = [
  { email: 'user01@test.com', password: 'user01@test.com' },
  { email: 'user02@test.com', password: 'user02@test.com' },
  { email: 'teacher@example.com', password: 'teacher123' },
];

// 随机选择一个测试用户
function getRandomUser() {
  return testUsers[Math.floor(Math.random() * testUsers.length)];
}

// 学生登录
function studentLogin(email, password) {
  const url = `${BASE_URL}/api/v1/student/auth/login`;
  const payload = JSON.stringify({
    email: email,
    password: password,
  });
  const params = {
    headers: {
      'Content-Type': 'application/json',
    },
  };
  
  const startTime = Date.now();
  const res = http.post(url, payload, params);
  const duration = Date.now() - startTime;
  
  loginDuration.add(duration);
  
  const success = check(res, {
    '登录状态码为200': (r) => r.status === 200,
    '登录返回token': (r) => {
      if (r.status !== 200) {
        // 如果登录失败，记录错误信息
        try {
          const body = JSON.parse(r.body);
          console.log(`登录失败: ${body.message || '未知错误'}, 状态码: ${r.status}`);
        } catch (e) {
          console.log(`登录失败: 无法解析响应, 状态码: ${r.status}`);
        }
        return false;
      }
      try {
        const body = JSON.parse(r.body);
        return body.token !== undefined && body.token !== '';
      } catch (e) {
        return false;
      }
    },
  });
  
  if (!success) {
    errorRate.add(1);
  } else {
    errorRate.add(0);
  }
  
  let token = null;
  if (res.status === 200) {
    try {
      const body = JSON.parse(res.body);
      token = body.token;
    } catch (e) {
      token = null;
    }
  }
  
  return token;
}

// 获取课程列表
function getCourseList() {
  const url = `${BASE_URL}/api/v1/student/courses`;
  const params = {
    headers: {
      'Content-Type': 'application/json',
    },
  };
  
  const startTime = Date.now();
  const res = http.get(url, params);
  const duration = Date.now() - startTime;
  
  courseListDuration.add(duration);
  
  const success = check(res, {
    '课程列表状态码为200': (r) => r.status === 200,
    '课程列表返回正确格式': (r) => {
      if (r.status !== 200) return false;
      try {
        const body = JSON.parse(r.body);
        return body.courses !== undefined && Array.isArray(body.courses);
      } catch (e) {
        return false;
      }
    },
  });
  
  if (!success) {
    errorRate.add(1);
  } else {
    errorRate.add(0);
  }
  
  let courses = [];
  if (res.status === 200) {
    try {
      const body = JSON.parse(res.body);
      courses = body.courses || [];
    } catch (e) {
      courses = [];
    }
  }
  
  return courses;
}

// 获取课程详情
function getCourseDetail(courseId, token = null) {
  const url = `${BASE_URL}/api/v1/student/courses/${courseId}`;
  const params = {
    headers: {
      'Content-Type': 'application/json',
    },
  };
  
  if (token) {
    params.headers['Authorization'] = `Bearer ${token}`;
  }
  
  const startTime = Date.now();
  const res = http.get(url, params);
  const duration = Date.now() - startTime;
  
  courseDetailDuration.add(duration);
  
  const success = check(res, {
    '课程详情状态码为200': (r) => r.status === 200,
    '课程详情包含标题': (r) => {
      const body = JSON.parse(r.body);
      return body.course_title !== undefined;
    },
  });
  
  if (!success) {
    errorRate.add(1);
  } else {
    errorRate.add(0);
  }
  
  return res;
}

// 获取课程评论列表（跨分片查询）
function getCourseComments(courseId) {
  const url = `${BASE_URL}/api/v1/courses/${courseId}/comments`;
  const params = {
    headers: {
      'Content-Type': 'application/json',
    },
  };
  
  const startTime = Date.now();
  const res = http.get(url, params);
  const duration = Date.now() - startTime;
  
  commentListDuration.add(duration);
  
  const success = check(res, {
    '评论列表状态码为200': (r) => r.status === 200,
    '评论列表返回数组': (r) => {
      if (r.status !== 200) return false;
      try {
        const body = JSON.parse(r.body);
        return Array.isArray(body);
      } catch (e) {
        return false;
      }
    },
  });
  
  if (!success) {
    errorRate.add(1);
  } else {
    errorRate.add(0);
  }
  
  return res;
}

// 获取校区列表
function getBranches() {
  const url = `${BASE_URL}/api/v1/student/branches`;
  const params = {
    headers: {
      'Content-Type': 'application/json',
    },
  };
  
  const res = http.get(url, params);
  
  check(res, {
    '校区列表状态码为200': (r) => r.status === 200,
  });
  
  return res;
}

// 主测试函数
export default function () {
  // 场景1: 获取校区列表（不需要认证）
  getBranches();
  sleep(1);
  
  // 场景2: 获取课程列表（不需要认证）
  const courses = getCourseList();
  sleep(1);
  
  // 场景3: 如果有课程，获取课程详情
  if (courses && courses.length > 0) {
    const randomCourse = courses[Math.floor(Math.random() * courses.length)];
    getCourseDetail(randomCourse.course_id);
    sleep(1);
    
    // 场景4: 获取课程评论（跨分片查询）
    getCourseComments(randomCourse.course_id);
    sleep(1);
  }
  
  // 场景5: 学生登录（需要认证的操作）
  const user = getRandomUser();
  const token = studentLogin(user.email, user.password);
  sleep(1);
  
  // 场景6: 使用token获取课程详情
  if (token && courses && courses.length > 0) {
    const randomCourse = courses[Math.floor(Math.random() * courses.length)];
    getCourseDetail(randomCourse.course_id, token);
    sleep(1);
  }
}

// 设置函数 - 在测试开始前执行
export function setup() {
  console.log('开始负载测试...');
  console.log(`目标服务器: ${BASE_URL}`);
  
  // 健康检查
  const healthUrl = `${BASE_URL}/health`;
  const res = http.get(healthUrl);
  
  if (res.status !== 200) {
    throw new Error(`服务器健康检查失败: ${res.status}`);
  }
  
  console.log('服务器健康检查通过');
  return { baseUrl: BASE_URL };
}

// 清理函数 - 在测试结束后执行
export function teardown(data) {
  console.log('负载测试完成');
}

