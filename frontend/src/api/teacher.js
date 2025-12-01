import request from './request'

// 教师登录
export const login = (data) => {
  return request.post('/teacher/auth/login', data)
}

// 获取个人信息
export const getProfile = () => {
  return request.get('/teacher/profile')
}

// 获取课程列表
export const getCourses = (params) => {
  return request.get('/teacher/courses', { params })
}

// 获取课程详情
export const getCourse = (id) => {
  return request.get(`/teacher/courses/${id}`)
}

// 创建课程
export const createCourse = (data) => {
  return request.post('/teacher/courses', data)
}

// 创建章节
export const createChapter = (courseId, data) => {
  return request.post(`/teacher/courses/${courseId}/chapters`, data)
}

// 创建课时（文件上传需要更长的超时时间）
export const createLesson = (courseId, chapterId, formData) => {
  return request.post(`/teacher/courses/${courseId}/chapters/${chapterId}/lessons`, formData, {
    headers: {
      'Content-Type': 'multipart/form-data'
    },
    timeout: 60000 // 文件上传设置为60秒超时
  })
}

// 创建任务
export const createTask = (lessonId, data) => {
  return request.post(`/teacher/lessons/${lessonId}/tasks`, data)
}

// 获取任务详情
export const getTask = (id) => {
  return request.get(`/teacher/tasks/${id}`)
}

// 获取课程的所有任务
export const getCourseTasks = (courseId) => {
  return request.get(`/teacher/courses/${courseId}/tasks`)
}

// 获取任务的所有作业
export const getTaskAnswers = (taskId) => {
  return request.get(`/teacher/tasks/${taskId}/answers`)
}

// 教师评分
export const gradeAnswer = (answerId, branchId, score) => {
  return request.put(`/teacher/answers/${answerId}/grade`, { branch_id: branchId, score })
}

// 获取课程评论列表（公共接口，教师也可以使用）
export const getCourseComments = (courseId) => {
  return request.get(`/courses/${courseId}/comments`)
}

// 教师发表评论
export const addComment = (courseId, data) => {
  return request.post(`/teacher/courses/${courseId}/comments`, data)
}

