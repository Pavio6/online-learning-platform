<template>
  <div class="create-course-container">
    <div class="header">
      <el-button @click="goBack" icon="ArrowLeft">返回</el-button>
      <h1>创建课程</h1>
      <div></div>
    </div>
    
    <div class="content">
      <el-card>
        <el-form
          ref="formRef"
          :model="form"
          :rules="rules"
          label-width="120px"
          style="max-width: 800px; margin: 0 auto;"
        >
          <el-form-item label="课程标题" prop="course_title">
            <el-input v-model="form.course_title" placeholder="请输入课程标题" />
          </el-form-item>
          
          <el-form-item label="课程描述" prop="description">
            <el-input
              v-model="form.description"
              type="textarea"
              :rows="4"
              placeholder="请输入课程描述"
            />
          </el-form-item>
          
          <el-form-item label="开课时间" prop="start_date">
            <el-date-picker
              v-model="form.start_date"
              type="datetime"
              placeholder="选择开课时间"
              format="YYYY-MM-DD HH:mm:ss"
              value-format="YYYY-MM-DD HH:mm:ss"
              style="width: 100%"
            />
          </el-form-item>
          
          <el-form-item label="结课时间" prop="end_date">
            <el-date-picker
              v-model="form.end_date"
              type="datetime"
              placeholder="选择结课时间"
              format="YYYY-MM-DD HH:mm:ss"
              value-format="YYYY-MM-DD HH:mm:ss"
              :disabled-date="disabledEndDate"
              style="width: 100%"
            />
          </el-form-item>
          
          <el-form-item label="课程状态" prop="status">
            <el-select v-model="form.status" placeholder="选择课程状态" style="width: 100%">
              <el-option label="进行中" value="active" />
              <el-option label="已结束" value="completed" />
            </el-select>
          </el-form-item>
          
          <el-form-item>
            <el-button type="primary" @click="handleSubmit" :loading="submitting">
              创建课程
            </el-button>
            <el-button @click="goBack">取消</el-button>
          </el-form-item>
        </el-form>
      </el-card>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, watch } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { createCourse } from '../../api/teacher'

const router = useRouter()
const formRef = ref(null)
const submitting = ref(false)

const form = reactive({
  course_title: '',
  description: '',
  start_date: '',
  end_date: '',
  status: 'active'
})

// 监听开始日期变化，如果结束日期小于新的开始日期，则清空结束日期
watch(() => form.start_date, (newStartDate) => {
  if (newStartDate && form.end_date) {
    const startDate = new Date(newStartDate)
    const endDate = new Date(form.end_date)
    if (endDate <= startDate) {
      form.end_date = ''
      ElMessage.warning('开课时间已更新，请重新选择结课时间')
    }
  }
})

// 验证结束日期必须大于开始日期
const validateEndDate = (rule, value, callback) => {
  if (!value) {
    callback(new Error('请选择结课时间'))
    return
  }
  if (!form.start_date) {
    callback(new Error('请先选择开课时间'))
    return
  }
  const startDate = new Date(form.start_date)
  const endDate = new Date(value)
  if (endDate <= startDate) {
    callback(new Error('结课时间必须大于开课时间'))
    return
  }
  callback()
}

// 禁用结束日期选择器中早于开始日期的日期
const disabledEndDate = (time) => {
  if (!form.start_date) {
    return false
  }
  const startDate = new Date(form.start_date)
  return time.getTime() <= startDate.getTime()
}

const rules = {
  course_title: [
    { required: true, message: '请输入课程标题', trigger: 'blur' }
  ],
  start_date: [
    { required: true, message: '请选择开课时间', trigger: 'change' }
  ],
  end_date: [
    { required: true, message: '请选择结课时间', trigger: 'change' },
    { validator: validateEndDate, trigger: 'change' }
  ],
  status: [
    { required: true, message: '请选择课程状态', trigger: 'change' }
  ]
}

const handleSubmit = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (valid) {
      submitting.value = true
      try {
        // 构建请求数据，只包含有值的字段
        const requestData = {
          course_title: form.course_title,
          description: form.description,
          status: form.status
        }
        
        // 只添加非空的日期字段
        if (form.start_date) {
          requestData.start_date = form.start_date
        }
        if (form.end_date) {
          requestData.end_date = form.end_date
        }
        
        console.log('发送的请求数据:', requestData)
        
        await createCourse(requestData)
        ElMessage.success('课程创建成功')
        router.push('/courses')
      } catch (error) {
        console.error('创建课程失败:', error)
        ElMessage.error(error.response?.data?.message || '创建课程失败')
      } finally {
        submitting.value = false
      }
    }
  })
}

const goBack = () => {
  router.push('/courses')
}
</script>

<style scoped>
.create-course-container {
  min-height: 100vh;
  background: #f5f5f5;
}

.header {
  background: white;
  padding: 20px 40px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.header h1 {
  margin: 0;
  font-size: 24px;
  color: #333;
}

.content {
  max-width: 1200px;
  margin: 0 auto;
  padding: 30px 20px;
}
</style>

