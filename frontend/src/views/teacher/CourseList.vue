<template>
  <div class="course-list-container">
    <div class="header">
      <h1>我的课程</h1>
      <div class="header-actions">
        <span class="username" v-if="authStore.userInfo?.username">
          欢迎，{{ authStore.userInfo.username }}
        </span>
        <el-button type="primary" @click="goToCreate">创建课程</el-button>
        <el-button @click="handleLogout">退出登录</el-button>
      </div>
    </div>
    
    <div class="content">
      <div v-loading="loading" class="courses-grid">
        <el-card
          v-for="course in courses"
          :key="course.course_id"
          class="course-card"
          shadow="hover"
          @click="goToCourseDetail(course.course_id)"
        >
          <div class="course-header">
            <h3>{{ course.course_title }}</h3>
            <el-tag :type="course.status === 'active' ? 'success' : 'info'">
              {{ course.status === 'active' ? '进行中' : '已结束' }}
            </el-tag>
          </div>
          <p class="course-description">{{ course.description || '暂无描述' }}</p>
          <div class="course-footer">
            <span class="course-date">
              {{ formatDate(course.start_date) }} - {{ formatDate(course.end_date) }}
            </span>
          </div>
        </el-card>
      </div>
      
      <el-empty v-if="!loading && courses.length === 0" description="暂无课程，点击创建课程开始" />
      
      <el-pagination
        v-if="total > 0"
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :total="total"
        :page-sizes="[10, 20, 50]"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="handleSizeChange"
        @current-change="handlePageChange"
        class="pagination"
      />
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { getCourses } from '../../api/teacher'
import { useAuthStore } from '../../stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const loading = ref(false)
const courses = ref([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)

const loadCourses = async () => {
  loading.value = true
  try {
    const res = await getCourses({
      page: currentPage.value,
      page_size: pageSize.value
    })
    courses.value = res.courses || []
    total.value = res.total || 0
  } catch (error) {
    console.error('获取课程列表失败:', error)
    ElMessage.error('获取课程列表失败')
  } finally {
    loading.value = false
  }
}

const handleSizeChange = (val) => {
  pageSize.value = val
  currentPage.value = 1
  loadCourses()
}

const handlePageChange = (val) => {
  currentPage.value = val
  loadCourses()
}

const goToCourseDetail = (courseId) => {
  router.push(`/courses/${courseId}`)
}

const goToCreate = () => {
  router.push('/courses/create')
}

const handleLogout = () => {
  authStore.logout()
  router.push('/login')
}

const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  const seconds = String(date.getSeconds()).padStart(2, '0')
  return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`
}

onMounted(() => {
  loadCourses()
})
</script>

<style scoped>
.course-list-container {
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

.header-actions {
  display: flex;
  align-items: center;
  gap: 15px;
}

.username {
  color: #666;
  font-size: 14px;
  margin-right: 10px;
}

.content {
  max-width: 1200px;
  margin: 0 auto;
  padding: 30px 20px;
}

.courses-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 20px;
  margin-bottom: 30px;
}

.course-card {
  cursor: pointer;
  transition: transform 0.2s;
}

.course-card:hover {
  transform: translateY(-5px);
}

.course-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 15px;
}

.course-header h3 {
  margin: 0;
  font-size: 18px;
  color: #333;
  flex: 1;
  margin-right: 10px;
}

.course-description {
  color: #666;
  font-size: 14px;
  line-height: 1.6;
  margin-bottom: 15px;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.course-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 15px;
  border-top: 1px solid #eee;
}

.course-date {
  font-size: 12px;
  color: #999;
}

.pagination {
  margin-top: 30px;
  display: flex;
  justify-content: center;
}
</style>

