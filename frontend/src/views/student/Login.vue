<template>
  <div class="login-container">
    <div class="login-left">
      <div class="login-content">
        <h1 class="main-title">在线学习平台</h1>
        <h2 class="sub-title">学生登录</h2>
        <p class="description">欢迎回来，请登录您的账号</p>
        <el-form
          ref="loginFormRef"
          :model="loginForm"
          :rules="loginRules"
          label-width="0"
          class="login-form"
        >
          <el-form-item prop="email">
            <el-input
              v-model="loginForm.email"
              placeholder="请输入邮箱"
              prefix-icon="Message"
              size="large"
            />
          </el-form-item>
          <el-form-item prop="password">
            <el-input
              v-model="loginForm.password"
              type="password"
              placeholder="请输入密码"
              prefix-icon="Lock"
              show-password
              size="large"
              @keyup.enter="handleLogin"
            />
          </el-form-item>
          <el-form-item>
            <el-button
              type="primary"
              :loading="loading"
              @click="handleLogin"
              size="large"
              style="width: 100%"
            >
              登录
            </el-button>
          </el-form-item>
          <el-form-item>
            <div class="register-link">
              还没有账号？
              <el-link type="primary" @click="$router.push('/register')">
                立即注册
              </el-link>
            </div>
          </el-form-item>
        </el-form>
      </div>
    </div>
    <div class="login-right">
      <div class="right-content">
        <h2>开始您的学习之旅</h2>
        <p>丰富的课程资源，专业的教学团队</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { login as studentLogin } from '../../api/student'
import { useAuthStore } from '../../stores/auth'

const router = useRouter()
const authStore = useAuthStore()
const loginFormRef = ref(null)
const loading = ref(false)

const loginForm = reactive({
  email: '',
  password: ''
})

const loginRules = {
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于6位', trigger: 'blur' }
  ]
}

const handleLogin = async () => {
  if (!loginFormRef.value) return
  
  await loginFormRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        const res = await studentLogin(loginForm)
        authStore.setToken(res.token)
        // 后端返回的是 LoginResponse，直接包含 username 等字段
        authStore.setUserInfo({
          user_id: res.user_id,
          username: res.username,
          email: res.email,
          role: res.role,
          branch_id: res.branch_id
        })
        authStore.setRole('student')
        ElMessage.success('登录成功')
        router.push('/courses')
      } catch (error) {
        console.error('登录失败:', error)
      } finally {
        loading.value = false
      }
    }
  })
}
</script>

<style scoped>
.login-container {
  display: flex;
  width: 100%;
  height: 100vh;
  margin: 0;
  padding: 0;
  overflow: hidden;
}

.login-left {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #fff;
  padding: 40px;
  margin: 0;
}

.login-content {
  width: 100%;
  max-width: 450px;
}

.main-title {
  font-size: 36px;
  font-weight: 700;
  color: #333;
  margin-bottom: 10px;
}

.sub-title {
  font-size: 28px;
  font-weight: 600;
  color: #333;
  margin-bottom: 10px;
}

.description {
  font-size: 16px;
  color: #666;
  margin-bottom: 40px;
}

.login-form {
  margin-top: 40px;
}

.register-link {
  text-align: center;
  width: 100%;
  color: #666;
  font-size: 14px;
}

.login-right {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 40px;
  position: relative;
  overflow: hidden;
  margin: 0;
  height: 100vh;
}

.right-content {
  color: white;
  text-align: center;
  z-index: 1;
}

.right-content h2 {
  font-size: 42px;
  font-weight: 700;
  margin-bottom: 20px;
}

.right-content p {
  font-size: 20px;
  opacity: 0.9;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .login-container {
    flex-direction: column;
  }
  
  .login-right {
    display: none;
  }
  
  .login-left {
    padding: 30px 20px;
  }
  
  .main-title {
    font-size: 28px;
  }
  
  .sub-title {
    font-size: 24px;
  }
}
</style>

