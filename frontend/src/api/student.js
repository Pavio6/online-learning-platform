import request from './request'

// 获取校区列表
export const getBranches = () => {
  return request.get('/student/branches')
}

// 学生注册
export const register = (data) => {
  return request.post('/student/auth/register', data)
}

// 学生登录
export const login = (data) => {
  return request.post('/student/auth/login', data)
}

// 获取个人信息
export const getProfile = () => {
  return request.get('/student/profile')
}

// 获取课程列表
export const getCourses = (params) => {
  return request.get('/student/courses', { params })
}

// 获取课程详情
export const getCourse = (id) => {
  return request.get(`/student/courses/${id}`)
}

// 报名课程
export const enrollCourse = (id) => {
  return request.post(`/student/courses/${id}/enroll`)
}

// 获取课程任务列表
export const getCourseTasks = (id) => {
  return request.get(`/student/courses/${id}/tasks`)
}

// 获取学习进度（用于检查是否已报名）
export const getProgress = (id) => {
  return request.get(`/student/courses/${id}/progress`)
}

// 获取已报名课程列表
export const getEnrolledCourses = (params) => {
  return request.get('/student/courses/enrolled', { params })
}

// 获取任务详情
export const getTask = (id) => {
  return request.get(`/student/tasks/${id}`)
}

// 提交作业（支持文件上传）
export const submitAnswer = (taskId, formData) => {
  return request.post(`/student/tasks/${taskId}/answers`, formData, {
    headers: {
      'Content-Type': 'multipart/form-data'
    },
    timeout: 60000 // 文件上传设置为60秒超时
  })
}

// 获取自己的作业
export const getMyAnswer = (taskId) => {
  return request.get(`/student/tasks/${taskId}/answers`)
}

// 获取课程评论列表
export const getCourseComments = (courseId) => {
  return request.get(`/courses/${courseId}/comments`)
}

// 发表评论
export const addComment = (courseId, data) => {
  return request.post(`/student/courses/${courseId}/comments`, data)
}

