<template>
  <div class="register-container">
    <div class="register-left">
      <div class="register-content">
        <h1 class="main-title">在线学习平台</h1>
        <h2 class="sub-title">学生注册</h2>
        <p class="description">创建您的账号，开始学习之旅</p>
        <el-form
          ref="registerFormRef"
          :model="registerForm"
          :rules="registerRules"
          label-width="0"
          class="register-form"
        >
          <el-form-item prop="username">
            <el-input
              v-model="registerForm.username"
              placeholder="请输入用户名"
              prefix-icon="User"
              size="large"
            />
          </el-form-item>
          <el-form-item prop="email">
            <el-input
              v-model="registerForm.email"
              placeholder="请输入邮箱"
              prefix-icon="Message"
              size="large"
            />
          </el-form-item>
          <el-form-item prop="password">
            <el-input
              v-model="registerForm.password"
              type="password"
              placeholder="请输入密码"
              prefix-icon="Lock"
              show-password
              size="large"
            />
          </el-form-item>
          <el-form-item prop="confirmPassword">
            <el-input
              v-model="registerForm.confirmPassword"
              type="password"
              placeholder="请再次输入密码"
              prefix-icon="Lock"
              show-password
              size="large"
            />
          </el-form-item>
          <el-form-item prop="branchId">
            <el-select
              v-model="registerForm.branchId"
              placeholder="请选择校区"
              style="width: 100%"
              size="large"
              :loading="branchesLoading"
            >
              <el-option
                v-for="branch in branches"
                :key="branch.branch_id"
                :label="branch.branch_name"
                :value="branch.branch_id"
              />
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button
              type="primary"
              :loading="loading"
              @click="handleRegister"
              size="large"
              style="width: 100%"
            >
              注册
            </el-button>
          </el-form-item>
          <el-form-item>
            <div class="login-link">
              已有账号？
              <el-link type="primary" @click="$router.push('/login')">
                立即登录
              </el-link>
            </div>
          </el-form-item>
        </el-form>
      </div>
    </div>
    <div class="register-right">
      <div class="right-content">
        <h2>加入我们的学习社区</h2>
        <p>与数千名学生一起，开启知识探索之旅</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { register as studentRegister, getBranches } from '../../api/student'

const router = useRouter()
const registerFormRef = ref(null)
const loading = ref(false)
const branchesLoading = ref(false)
const branches = ref([])

const registerForm = reactive({
  username: '',
  email: '',
  password: '',
  confirmPassword: '',
  branchId: null
})

const validateConfirmPassword = (rule, value, callback) => {
  if (value !== registerForm.password) {
    callback(new Error('两次输入的密码不一致'))
  } else {
    callback()
  }
}

const registerRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '用户名长度在3到20个字符', trigger: 'blur' }
  ],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于6位', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请再次输入密码', trigger: 'blur' },
    { validator: validateConfirmPassword, trigger: 'blur' }
  ],
  branchId: [
    { required: true, message: '请选择校区', trigger: 'change' }
  ]
}

const loadBranches = async () => {
  branchesLoading.value = true
  try {
    const res = await getBranches()
    branches.value = res
  } catch (error) {
    console.error('获取校区列表失败:', error)
  } finally {
    branchesLoading.value = false
  }
}

const handleRegister = async () => {
  if (!registerFormRef.value) return
  
  await registerFormRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        const { confirmPassword, branchId, ...formData } = registerForm
        // 将 branchId 转换为 branch_id 以匹配后端API
        const requestData = {
          ...formData,
          branch_id: branchId
        }
        await studentRegister(requestData)
        ElMessage.success('注册成功，请登录')
        router.push('/login')
      } catch (error) {
        console.error('注册失败:', error)
      } finally {
        loading.value = false
      }
    }
  })
}

onMounted(() => {
  loadBranches()
})
</script>

<style scoped>
.register-container {
  display: flex;
  width: 100%;
  height: 100vh;
  margin: 0;
  padding: 0;
  overflow: hidden;
}

.register-left {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #fff;
  padding: 40px;
  overflow-y: auto;
  margin: 0;
}

.register-content {
  width: 100%;
  max-width: 500px;
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

.register-form {
  margin-top: 40px;
}

.login-link {
  text-align: center;
  width: 100%;
  color: #666;
  font-size: 14px;
}

.register-right {
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
  .register-container {
    flex-direction: column;
  }
  
  .register-right {
    display: none;
  }
  
  .register-left {
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

