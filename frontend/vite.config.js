import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { resolve } from 'path'
import fs from 'fs'

// https://vite.dev/config/
export default defineConfig(({ mode }) => {
  const isStudent = mode === 'student'
  const isTeacher = mode === 'teacher'
  
  // 根据模式选择对应的 HTML 文件
  const htmlFile = isStudent ? 'index-student.html' : 'index-teacher.html'
  
  // 创建一个插件，在开发模式下将 index.html 指向正确的文件
  const htmlPlugin = () => {
    return {
      name: 'html-transform',
      configureServer(server) {
        server.middlewares.use((req, res, next) => {
          const url = req.url.split('?')[0] // 移除查询参数
          
          // 跳过 API 请求
          if (url.startsWith('/api')) {
            next()
            return
          }
          
          // 跳过 Vite 内部请求
          if (url.startsWith('/@')) {
            next()
            return
          }
          
          // 跳过静态资源（有文件扩展名且不是 .html）
          const hasExtension = /\.\w+$/.test(url)
          if (hasExtension && !url.endsWith('.html')) {
            next()
            return
          }
          
          // 对于所有其他请求（包括前端路由），返回对应的 HTML 文件
          // 这样 Vue Router 就能处理这些路由了
          const htmlPath = resolve(__dirname, htmlFile)
          if (fs.existsSync(htmlPath)) {
            const html = fs.readFileSync(htmlPath, 'utf-8')
            res.setHeader('Content-Type', 'text/html')
            res.end(html)
            return
          }
          
          next()
        })
      }
    }
  }
  
  return {
    plugins: [vue(), htmlPlugin()],
    root: '.',
    publicDir: 'public',
    build: {
      rollupOptions: {
        input: {
          student: resolve(__dirname, 'index-student.html'),
          teacher: resolve(__dirname, 'index-teacher.html')
        }
      }
    },
    server: {
      port: isStudent ? 3000 : 3001,
      proxy: {
        '/api': {
          target: 'http://localhost:8080',
          changeOrigin: true
        }
      }
    }
  }
})
