import axios from 'axios'
import { ElMessage } from 'element-plus'

// 开发环境使用相对路径通过Vite代理，生产环境使用完整URL
const baseURL = import.meta.env.DEV 
  ? '/api/v1'  // 开发环境使用代理
  : (import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1')

const request = axios.create({
  baseURL: baseURL,
  timeout: 10000
})

// 请求拦截器
request.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 响应拦截器
request.interceptors.response.use(
  (response) => {
    return response.data
  },
  (error) => {
    if (error.response) {
      const { status, data } = error.response
      if (status === 401) {
        localStorage.removeItem('token')
        localStorage.removeItem('userInfo')
        localStorage.removeItem('role')
        ElMessage.error('登录已过期，请重新登录')
        // 根据当前路径判断是学生端还是教师端
        const role = localStorage.getItem('role') || 'student'
        window.location.href = role === 'teacher' ? '/teacher/login' : '/login'
      } else if (data?.code === 4001) {
        // 错误码 4001 表示"未报名课程"，这是正常状态，不显示错误提示
        // 静默处理，让调用方自己决定如何处理
      } else {
        ElMessage.error(data?.message || '请求失败')
      }
    } else if (error.code === 'ECONNABORTED' || error.message?.includes('timeout')) {
      // 超时错误
      console.error('请求超时:', error.message)
      ElMessage.error('请求超时，文件可能较大，请稍后刷新页面查看结果')
    } else if (error.request) {
      // 请求已发出但没有收到响应
      console.error('请求失败，未收到响应:', error.request)
      ElMessage.error('无法连接到服务器，请检查后端服务是否运行（http://localhost:8080）')
    } else {
      // 请求配置错误
      console.error('请求配置错误:', error.message)
      ElMessage.error('网络错误: ' + (error.message || '请检查网络连接'))
    }
    return Promise.reject(error)
  }
)

export default request

