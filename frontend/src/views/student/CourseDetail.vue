<template>
  <div class="course-detail-container">
    <div class="header">
      <el-button @click="goBack" icon="ArrowLeft">返回</el-button>
      <h1 v-if="course">{{ course.course_title }}</h1>
      <div></div>
    </div>
    
    <div v-loading="loading" class="content">
      <div v-if="course" class="course-info">
        <el-card class="info-card">
          <div class="course-meta">
            <p><strong>课程描述：</strong>{{ course.description || '暂无描述' }}</p>
            <p><strong>授课教师：</strong>{{ getInstructorName(course) }}</p>
            <p><strong>状态：</strong>
              <el-tag :type="course.status === 'active' ? 'success' : 'info'">
                {{ course.status === 'active' ? '进行中' : '已结束' }}
              </el-tag>
            </p>
            <p><strong>开课时间：</strong>{{ formatDate(course.start_date) }}</p>
            <p><strong>结课时间：</strong>{{ formatDate(course.end_date) }}</p>
            
            <!-- 未报名提示 -->
            <div v-if="authStore.isAuthenticated() && !isEnrolled" class="enrollment-notice">
              <p class="enrollment-text">学生暂未报名该课程</p>
            </div>
          </div>
          
          <div class="actions" v-if="authStore.isAuthenticated()">
            <el-button
              v-if="!isEnrolled"
              type="primary"
              size="large"
              @click="handleEnroll"
              :loading="enrolling"
            >
              报名课程
            </el-button>
            <el-button
              v-else
              type="success"
              size="large"
              disabled
            >
              已报名
            </el-button>
          </div>
        </el-card>
        
        <!-- 任务列表（仅已报名学生可见） -->
        <el-card v-if="isEnrolled" class="tasks-card" style="margin-top: 20px">
          <template #header>
            <h2>课程任务</h2>
          </template>
          <div v-loading="loadingTasks">
            <div v-if="tasksList.length > 0">
              <!-- 按课时分组显示任务 -->
              <el-collapse v-model="activeTaskChapter" accordion>
                <el-collapse-item
                  v-for="chapterGroup in groupedTasks"
                  :key="chapterGroup.chapterId"
                  :name="chapterGroup.chapterId"
                  :title="`${chapterGroup.chapterOrder}. ${chapterGroup.chapterTitle}`"
                >
                  <div v-for="lessonGroup in chapterGroup.lessons" :key="lessonGroup.lessonId" class="lesson-task-group">
                    <div class="lesson-task-header">
                      <span class="lesson-title-text">
                        {{ lessonGroup.lessonOrder }}. {{ lessonGroup.lessonTitle }}
                      </span>
                    </div>
                    <el-table :data="lessonGroup.tasks" style="width: 100%" border>
                      <el-table-column prop="task_title" label="任务标题" min-width="200" />
                      <el-table-column prop="task_type" label="类型" width="100">
                        <template #default="{ row }">
                          <el-tag size="small">{{ getTaskTypeName(row.task_type) }}</el-tag>
                        </template>
                      </el-table-column>
                      <el-table-column prop="max_score" label="最高分" width="100" />
                      <el-table-column label="状态" width="120">
                        <template #default="{ row }">
                          <el-tag
                            :type="getTaskStatus(row.task_id) === 'submitted' ? 'success' : 'info'"
                            size="small"
                          >
                            {{ getTaskStatus(row.task_id) === 'submitted' ? '已提交' : '未提交' }}
                          </el-tag>
                        </template>
                      </el-table-column>
                      <el-table-column label="操作" width="150" fixed="right">
                        <template #default="{ row }">
                          <el-button
                            type="primary"
                            size="small"
                            @click="openSubmitAnswerDialog(row)"
                          >
                            {{ getTaskStatus(row.task_id) === 'submitted' ? '查看作业' : '提交作业' }}
                          </el-button>
                        </template>
                      </el-table-column>
                    </el-table>
                  </div>
                </el-collapse-item>
              </el-collapse>
            </div>
            <el-empty v-else description="该课程暂无任务" />
          </div>
        </el-card>

        <el-card class="chapters-card" v-if="course.chapters && course.chapters.length > 0">
          <template #header>
            <h2>课程章节</h2>
          </template>
          <!-- 已报名的学生：可以展开查看课时 -->
          <div v-if="isEnrolled">
            <el-collapse>
              <el-collapse-item
                v-for="chapter in course.chapters"
                :key="chapter.chapter_id"
                :title="`${chapter.chapter_order}. ${chapter.chapter_title}`"
              >
                <div v-if="chapter.lessons && chapter.lessons.length > 0">
                  <div
                    v-for="lesson in chapter.lessons"
                    :key="lesson.lesson_id"
                    class="lesson-item"
                  >
                    <div class="lesson-info">
                      <span class="lesson-order">{{ lesson.lesson_order }}</span>
                      <span class="lesson-title">{{ lesson.lesson_title }}</span>
                    </div>
                    <el-button
                      v-if="lesson.content_url"
                      type="primary"
                      size="small"
                      @click="openMediaViewer(lesson)"
                    >
                      {{ lesson.lesson_type === 'video' ? '观看视频' : '查看文档' }}
                    </el-button>
                  </div>
                </div>
                <p v-else class="empty-lessons">暂无课程</p>
              </el-collapse-item>
            </el-collapse>
          </div>
          <!-- 未报名时：只显示章节标题，不可展开 -->
          <div v-else class="chapters-list-unenrolled">
            <div
              v-for="chapter in course.chapters"
              :key="chapter.chapter_id"
              class="chapter-item-unenrolled"
            >
              <span class="chapter-title">{{ chapter.chapter_order }}. {{ chapter.chapter_title }}</span>
            </div>
          </div>
        </el-card>

        <!-- 课程评论（仅已报名学生可见） -->
        <el-card v-if="isEnrolled" class="comments-card" style="margin-top: 20px">
          <template #header>
            <h2>课程评论</h2>
          </template>
          <div v-loading="loadingComments">
            <!-- 发表评论表单 -->
            <div class="comment-form-section">
              <el-form :model="commentForm" @submit.prevent="handleSubmitComment">
                <el-form-item>
                  <el-input
                    v-model="commentForm.content"
                    type="textarea"
                    :rows="4"
                    placeholder="请输入评论内容..."
                    maxlength="500"
                    show-word-limit
                  />
                </el-form-item>
                <el-form-item>
                  <el-button
                    type="primary"
                    @click="handleSubmitComment"
                    :loading="submittingComment"
                    :disabled="!commentForm.content.trim()"
                  >
                    发表评论
                  </el-button>
                </el-form-item>
              </el-form>
            </div>

            <!-- 评论列表 -->
            <div class="comments-list">
              <div v-if="commentsList.length > 0">
                <div
                  v-for="comment in topLevelComments"
                  :key="comment.comment_id"
                  class="comment-item"
                >
                  <!-- 主评论 -->
                  <div class="comment-main">
                    <div class="comment-header">
                      <span class="comment-author">{{ comment.username }}</span>
                      <span class="comment-time">{{ formatCommentTime(comment.created_at) }}</span>
                    </div>
                    <div class="comment-content">{{ comment.comment_content }}</div>
                    <div class="comment-actions">
                      <el-button
                        type="text"
                        size="small"
                        @click="toggleReply(comment.comment_id)"
                      >
                        <el-icon style="margin-right: 4px;"><ChatLineRound /></el-icon>
                        {{ replyingTo === comment.comment_id ? '取消回复' : '回复' }}
                      </el-button>
                    </div>
                  </div>

                  <!-- 回复表单 -->
                  <div v-if="replyingTo === comment.comment_id" class="reply-form">
                    <el-input
                      v-model="replyForm.content"
                      type="textarea"
                      :rows="2"
                      placeholder="请输入回复内容..."
                      maxlength="500"
                      show-word-limit
                      style="margin-bottom: 10px;"
                    />
                    <div>
                      <el-button
                        type="primary"
                        size="small"
                        @click="handleSubmitReply(comment.comment_id)"
                        :loading="submittingReply"
                        :disabled="!replyForm.content.trim()"
                      >
                        提交回复
                      </el-button>
                      <el-button
                        size="small"
                        @click="cancelReply"
                      >
                        取消
                      </el-button>
                    </div>
                  </div>

                  <!-- 回复列表 -->
                  <div v-if="getReplies(comment.comment_id).length > 0" class="replies-list">
                    <div class="replies-title">回复 ({{ getReplies(comment.comment_id).length }})</div>
                    <div
                      v-for="reply in getReplies(comment.comment_id)"
                      :key="reply.comment_id"
                      class="reply-item"
                    >
                      <div class="comment-header">
                        <span class="comment-author">{{ reply.username }}</span>
                        <span class="comment-time">{{ formatCommentTime(reply.created_at) }}</span>
                      </div>
                      <div class="comment-content">{{ reply.comment_content }}</div>
                    </div>
                  </div>
                </div>
              </div>
              <el-empty v-else description="暂无评论，快来发表第一条评论吧！" />
            </div>
          </div>
        </el-card>
      </div>
      
      <el-empty v-if="!loading && !course" description="课程不存在" />
    </div>

    <!-- 提交作业对话框 -->
    <el-dialog 
      :model-value="showSubmitAnswerDialog" 
      :title="isAnswerSubmitted ? '查看作业' : '提交作业'" 
      width="700px"
      @close="showSubmitAnswerDialog = false"
    >
      <el-form ref="answerFormRef" :model="answerForm" label-width="100px">
        <!-- 任务详情 -->
        <el-card shadow="never" style="margin-bottom: 20px; background: #f5f7fa;">
          <h3 style="margin-top: 0; margin-bottom: 15px;">{{ currentTask?.task_title }}</h3>
          <div style="margin-bottom: 10px;">
            <el-tag style="margin-right: 10px;">{{ getTaskTypeName(currentTask?.task_type) }}</el-tag>
            <span style="color: #909399; font-size: 14px;">最高分：{{ currentTask?.max_score }}</span>
          </div>
          <div v-if="currentTask?.description" style="margin-top: 15px;">
            <p style="margin: 0; color: #606266; line-height: 1.6;">{{ currentTask.description }}</p>
          </div>
          <div v-else style="margin-top: 15px; color: #909399; font-style: italic;">
            暂无任务描述
          </div>
        </el-card>
        
        <!-- 已提交状态：只显示内容，不允许修改 -->
        <div v-if="isAnswerSubmitted">
          <el-form-item label="提交状态">
            <el-tag type="success">已提交</el-tag>
            <span v-if="myAnswer?.is_graded" style="margin-left: 15px; color: #67c23a; font-weight: bold;">
              分数：{{ myAnswer.score }} / {{ currentTask?.max_score }}
            </span>
            <span v-else style="margin-left: 15px; color: #909399;">
              待教师评分
            </span>
          </el-form-item>
          
          <!-- 作文类型：显示已提交的文本内容 -->
          <el-form-item v-if="currentTask?.task_type === 'essay'" label="我的作业">
            <div style="padding: 15px; background: #f5f7fa; border-radius: 4px; border: 1px solid #e4e7ed; min-height: 100px;">
              <p style="margin: 0; white-space: pre-wrap; word-wrap: break-word;">{{ myAnswer?.answer_content || '无内容' }}</p>
            </div>
          </el-form-item>
          
          <!-- 上传类型：显示文件名和下载链接 -->
          <el-form-item v-if="currentTask?.task_type === 'upload'" label="我的作业">
            <div v-if="myAnswer?.answer_content" style="padding: 15px; background: #f5f7fa; border-radius: 4px; border: 1px solid #e4e7ed;">
              <el-link
                :underline="false"
                type="primary"
                @click="downloadFile(myAnswer.answer_content)"
                style="cursor: pointer; font-size: 14px;"
              >
                <el-icon style="margin-right: 4px; vertical-align: middle;"><Download /></el-icon>
                {{ getFileNameFromUrl(myAnswer.answer_content) }}
              </el-link>
            </div>
            <div v-else style="padding: 15px; background: #f5f7fa; border-radius: 4px; border: 1px solid #e4e7ed; color: #909399;">
              无文件
            </div>
          </el-form-item>
        </div>
        
        <!-- 未提交状态：显示输入框 -->
        <div v-else>
          <!-- 作文类型：只显示文本输入框 -->
          <el-form-item
            v-if="currentTask?.task_type === 'essay'"
            label="作业内容"
            prop="answer_content"
          >
            <el-input
              v-model="answerForm.answer_content"
              type="textarea"
              :rows="8"
              placeholder="请输入作业内容"
            />
          </el-form-item>
          
          <!-- 上传类型：只显示文件上传按钮 -->
          <el-form-item
            v-if="currentTask?.task_type === 'upload'"
            label="上传文件"
            prop="file"
          >
            <el-upload
              :auto-upload="false"
              :on-change="handleFileChange"
              :file-list="fileList"
              accept="image/jpeg,image/jpg,image/png,image/gif"
              :limit="1"
            >
              <el-button type="primary">选择文件</el-button>
              <template #tip>
                <div class="el-upload__tip">
                  <p style="margin: 5px 0; color: #909399; font-size: 12px;">
                    <strong>文件要求：</strong>
                  </p>
                  <ul style="margin: 5px 0; padding-left: 20px; color: #909399; font-size: 12px;">
                    <li>支持格式：jpg、png、gif（图片格式）</li>
                    <li>文件数量：1 个文件</li>
                    <li>文件大小：单个文件不超过 10MB</li>
                  </ul>
                </div>
              </template>
            </el-upload>
          </el-form-item>
        </div>
      </el-form>
      <template #footer>
        <el-button @click="showSubmitAnswerDialog = false">{{ isAnswerSubmitted ? '关闭' : '取消' }}</el-button>
        <el-button 
          v-if="!isAnswerSubmitted"
          type="primary" 
          @click="handleSubmitAnswer" 
          :loading="submitting"
        >
          提交
        </el-button>
      </template>
    </el-dialog>

    <!-- 视频/文档查看对话框 -->
    <el-dialog
      :model-value="showMediaDialog"
      :title="currentLesson ? (currentLesson.lesson_type === 'video' ? '观看视频' : '查看文档') : ''"
      width="90%"
      :close-on-click-modal="false"
      @close="closeMediaViewer"
    >
      <div v-if="currentLesson" class="media-viewer-container">
        <!-- 视频播放器 -->
        <div v-if="currentLesson.lesson_type === 'video'" class="video-container">
          <video
            :src="currentLesson.content_url"
            controls
            autoplay
            style="width: 100%; max-height: 70vh;"
          >
            您的浏览器不支持视频播放。
          </video>
        </div>
        
        <!-- 文档查看器 -->
        <div v-else class="document-container">
          <!-- Office 文档（docx, doc, xlsx, xls, pptx, ppt）使用在线查看器 -->
          <div v-if="isOfficeDocument(currentLesson.content_url)" class="office-document-viewer">
            <iframe
              :src="getOfficeViewerUrl(currentLesson.content_url)"
              style="width: 100%; height: 70vh; border: 1px solid #e4e7ed; border-radius: 4px;"
              frameborder="0"
            >
              您的浏览器不支持iframe，请
              <a :href="currentLesson.content_url" target="_blank">点击这里</a>
              下载文档。
            </iframe>
            <!-- 提示信息放在下方 -->
            <div class="viewer-tip-bottom">
              <el-alert
                title="提示"
                type="info"
                :closable="false"
                show-icon
              >
                <template #default>
                  <p style="margin: 0;">如果文档无法正常显示，可
                    <el-button
                      type="primary"
                      size="small"
                      @click="downloadFile(currentLesson.content_url)"
                      style="margin: 0 5px;"
                    >
                      <el-icon style="margin-right: 4px;"><Download /></el-icon>
                      点击下载按钮查看
                    </el-button>
                  </p>
                </template>
              </el-alert>
            </div>
          </div>
          
          <!-- PDF 文档直接显示 -->
          <div v-else-if="isPdfDocument(currentLesson.content_url)" class="pdf-viewer">
            <!-- 使用 Google Docs Viewer 在线查看 PDF（避免下载问题） -->
            <iframe
              :src="getPdfViewerUrl(currentLesson.content_url)"
              style="width: 100%; height: 70vh; border: 1px solid #e4e7ed; border-radius: 4px;"
              frameborder="0"
            >
              您的浏览器不支持PDF预览，请
              <a :href="currentLesson.content_url" target="_blank">点击这里</a>
              下载PDF文件。
            </iframe>
            <!-- 提示信息放在下方 -->
            <div class="viewer-tip-bottom">
              <el-alert
                title="提示"
                type="info"
                :closable="false"
                show-icon
              >
                <template #default>
                  <p style="margin: 0;">如果文档无法正常显示，可
                    <el-button
                      type="primary"
                      size="small"
                      @click="downloadFile(currentLesson.content_url)"
                      style="margin: 0 5px;"
                    >
                      <el-icon style="margin-right: 4px;"><Download /></el-icon>
                      点击下载按钮查看
                    </el-button>
                  </p>
                </template>
              </el-alert>
            </div>
          </div>
          
          <!-- 图片直接显示 -->
          <div v-else-if="isImageDocument(currentLesson.content_url)" class="image-viewer">
            <img
              :src="currentLesson.content_url"
              alt="文档图片"
              style="max-width: 100%; max-height: 70vh; border: 1px solid #e4e7ed; border-radius: 4px;"
            />
          </div>
          
          <!-- 其他格式，提供下载链接 -->
          <div v-else class="unsupported-document">
            <el-empty description="该文档格式不支持在线预览">
              <el-button type="primary" @click="downloadFile(currentLesson.content_url)">
                <el-icon style="margin-right: 4px;"><Download /></el-icon>
                下载文档
              </el-button>
            </el-empty>
          </div>
        </div>
      </div>
      <template #footer>
        <el-button @click="closeMediaViewer">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, computed, reactive } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Download, ChatLineRound } from '@element-plus/icons-vue'
import { getCourse, enrollCourse, getProgress, getCourseTasks, getTask, submitAnswer, getMyAnswer, getCourseComments, addComment } from '../../api/student'
import { useAuthStore } from '../../stores/auth'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

const loading = ref(false)
const course = ref(null)
const enrolling = ref(false)
const isEnrolled = ref(false)
const loadingTasks = ref(false)
const tasksList = ref([])
const showSubmitAnswerDialog = ref(false)
const currentTask = ref(null)
const myAnswer = ref(null)
const submitting = ref(false)
const fileList = ref([])
const answerFormRef = ref(null)
const isAnswerSubmitted = ref(false) // 是否已提交作业

const answerForm = reactive({
  submitType: 'text',
  answer_content: '',
  file: null
})

const taskStatusMap = ref({}) // 记录每个任务的提交状态
const activeTaskChapter = ref([]) // 展开的章节（用于手风琴模式）

// 视频/文档查看相关
const showMediaDialog = ref(false)
const currentLesson = ref(null)

// 评论相关
const loadingComments = ref(false)
const commentsList = ref([])
const submittingComment = ref(false)
const submittingReply = ref(false)
const replyingTo = ref(null)
const commentForm = reactive({
  content: ''
})
const replyForm = reactive({
  content: ''
})

const loadCourse = async () => {
  loading.value = true
  try {
    const courseId = route.params.id
    const res = await getCourse(courseId)
    course.value = res
    
    // 检查是否已报名（只有登录用户才检查）
    if (authStore.isAuthenticated()) {
      await checkEnrollment(courseId)
      // 如果已报名，加载任务列表和评论
      if (isEnrolled.value) {
        await loadTasks(courseId)
        await loadComments(courseId)
      }
    } else {
      isEnrolled.value = false
    }
  } catch (error) {
    console.error('获取课程详情失败:', error)
    ElMessage.error('获取课程详情失败')
  } finally {
    loading.value = false
  }
}

// 检查是否已报名
const checkEnrollment = async (courseId) => {
  try {
    const res = await getProgress(courseId)
    // 检查返回的数据
    if (res.enrolled === false) {
      // 明确返回了未报名状态
      isEnrolled.value = false
    } else {
      // 返回了学习进度数据，说明已报名
      isEnrolled.value = true
    }
  } catch (error) {
    // 如果请求失败（网络错误等），默认设置为未报名
    console.warn('检查报名状态失败:', error)
    isEnrolled.value = false
  }
}

const handleEnroll = async () => {
  if (!authStore.isAuthenticated()) {
    ElMessage.warning('请先登录')
    router.push('/login')
    return
  }
  
  enrolling.value = true
  try {
    await enrollCourse(route.params.id)
    ElMessage.success('报名成功')
    isEnrolled.value = true
    // 报名成功后，重新加载课程详情以显示章节内容
    isEnrolled.value = true
    await loadCourse()
    // 加载任务列表和评论
    await loadTasks(route.params.id)
    await loadComments(route.params.id)
  } catch (error) {
    console.error('报名失败:', error)
    ElMessage.error(error.response?.data?.message || '报名失败')
  } finally {
    enrolling.value = false
  }
}

// 打开视频/文档查看器
const openMediaViewer = (lesson) => {
  currentLesson.value = lesson
  showMediaDialog.value = true
}

// 关闭视频/文档查看器
const closeMediaViewer = () => {
  showMediaDialog.value = false
  currentLesson.value = null
}

// 检测是否为 Office 文档
const isOfficeDocument = (url) => {
  if (!url) return false
  const lowerUrl = url.toLowerCase()
  return lowerUrl.endsWith('.docx') || 
         lowerUrl.endsWith('.doc') || 
         lowerUrl.endsWith('.xlsx') || 
         lowerUrl.endsWith('.xls') || 
         lowerUrl.endsWith('.pptx') || 
         lowerUrl.endsWith('.ppt')
}

// 检测是否为 PDF 文档
const isPdfDocument = (url) => {
  if (!url) return false
  return url.toLowerCase().endsWith('.pdf')
}

// 检测是否为图片
const isImageDocument = (url) => {
  if (!url) return false
  const lowerUrl = url.toLowerCase()
  return lowerUrl.endsWith('.jpg') || 
         lowerUrl.endsWith('.jpeg') || 
         lowerUrl.endsWith('.png') || 
         lowerUrl.endsWith('.gif') || 
         lowerUrl.endsWith('.bmp') || 
         lowerUrl.endsWith('.webp') ||
         lowerUrl.endsWith('.svg')
}

// 获取 Office 文档在线查看器 URL
const getOfficeViewerUrl = (url) => {
  // 使用 Microsoft Office Online Viewer
  // 注意：需要文档是公开可访问的，或者使用其他在线查看器服务
  const encodedUrl = encodeURIComponent(url)
  return `https://view.officeapps.live.com/op/embed.aspx?src=${encodedUrl}`
}

// 获取 PDF 在线查看器 URL（Google Docs Viewer）
const getPdfViewerUrl = (url) => {
  const encodedUrl = encodeURIComponent(url)
  return `https://docs.google.com/viewer?url=${encodedUrl}&embedded=true`
}

const goBack = () => {
  router.push('/courses')
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

const loadTasks = async (courseId) => {
  if (!authStore.isAuthenticated()) return
  
  loadingTasks.value = true
  try {
    const res = await getCourseTasks(courseId)
    tasksList.value = res || []
    // 检查每个任务的提交状态
    for (const task of tasksList.value) {
      await checkTaskStatus(task.task_id)
    }
  } catch (error) {
    console.error('获取任务列表失败:', error)
    // 不显示错误，因为可能课程没有任务
  } finally {
    loadingTasks.value = false
  }
}

const checkTaskStatus = async (taskId) => {
  try {
    const res = await getMyAnswer(taskId)
    // 检查是否是未提交状态
    if (res && res.submitted === false) {
      taskStatusMap.value[taskId] = 'not_submitted'
    } else if (res && res.answer_id) {
      taskStatusMap.value[taskId] = 'submitted'
    } else {
      taskStatusMap.value[taskId] = 'not_submitted'
    }
  } catch (error) {
    // 如果获取失败，说明未提交
    taskStatusMap.value[taskId] = 'not_submitted'
  }
}

const getTaskStatus = (taskId) => {
  return taskStatusMap.value[taskId] || 'not_submitted'
}


const openSubmitAnswerDialog = async (task) => {
  // 先获取完整的任务详情
  try {
    const taskDetail = await getTask(task.task_id)
    currentTask.value = taskDetail
  } catch (error) {
    console.error('获取任务详情失败:', error)
    ElMessage.error('获取任务详情失败')
    return
  }
  
  showSubmitAnswerDialog.value = true
  answerForm.answer_content = ''
  answerForm.file = null
  fileList.value = []
  isAnswerSubmitted.value = false
  myAnswer.value = null
  
  // 检查是否已提交作业
  try {
    const answer = await getMyAnswer(task.task_id)
    // 检查是否是未提交状态
    if (answer && answer.submitted === false) {
      // 未提交
      isAnswerSubmitted.value = false
      myAnswer.value = null
    } else if (answer && answer.answer_id) {
      // 已提交
      isAnswerSubmitted.value = true
      myAnswer.value = answer
    }
  } catch (error) {
    // 如果仍然有错误，记录但不显示
    console.warn('获取作业信息失败:', error)
    isAnswerSubmitted.value = false
    myAnswer.value = null
  }
}

const handleFileChange = (file) => {
  answerForm.file = file.raw
}

const handleSubmitAnswer = async () => {
  if (!answerFormRef.value) return
  
  // 根据任务类型进行验证
  if (currentTask.value.task_type === 'essay') {
    // 作文类型：验证文本内容
    if (!answerForm.answer_content.trim()) {
      ElMessage.warning('请输入作业内容')
      return
    }
  } else if (currentTask.value.task_type === 'upload') {
    // 上传类型：验证文件
    if (!answerForm.file) {
      ElMessage.warning('请选择要上传的文件')
      return
    }
  }
  
  submitting.value = true
  try {
    const formData = new FormData()
    if (currentTask.value.task_type === 'essay') {
      // 作文类型：提交文本
      formData.append('answer_content', answerForm.answer_content)
      formData.append('type', 'text')
    } else if (currentTask.value.task_type === 'upload') {
      // 上传类型：提交文件
      formData.append('file', answerForm.file)
      formData.append('type', 'image_url')
    }
    
    await submitAnswer(currentTask.value.task_id, formData)
    ElMessage.success('作业提交成功')
    
    // 更新任务状态
    taskStatusMap.value[currentTask.value.task_id] = 'submitted'
    
    // 重新获取作业信息，切换到查看模式
    try {
      const answer = await getMyAnswer(currentTask.value.task_id)
      if (answer && answer.answer_id) {
        isAnswerSubmitted.value = true
        myAnswer.value = answer
      }
    } catch (error) {
      console.warn('获取作业信息失败:', error)
    }
    
    // 重新加载任务列表（更新按钮文本）
    await loadTasks(route.params.id)
  } catch (error) {
    console.error('提交作业失败:', error)
    ElMessage.error(error.response?.data?.message || '提交作业失败')
  } finally {
    submitting.value = false
  }
}

const getTaskTypeName = (type) => {
  const typeMap = {
    essay: '作文',
    upload: '上传'
  }
  return typeMap[type] || type
}

// 按章节和课时分组任务
const groupedTasks = computed(() => {
  const groups = {}
  
  tasksList.value.forEach(task => {
    const chapterId = task.chapter_id || 0
    const chapterTitle = task.chapter_title || '未分类'
    const chapterOrder = task.chapter_order || 0
    const lessonId = task.lesson_id || 0
    const lessonTitle = task.lesson_title || '未知课时'
    const lessonOrder = task.lesson_order || 0
    
    if (!groups[chapterId]) {
      groups[chapterId] = {
        chapterId,
        chapterTitle,
        chapterOrder,
        lessons: {}
      }
    }
    
    if (!groups[chapterId].lessons[lessonId]) {
      groups[chapterId].lessons[lessonId] = {
        lessonId,
        lessonTitle,
        lessonOrder,
        tasks: []
      }
    }
    
    groups[chapterId].lessons[lessonId].tasks.push(task)
  })
  
  // 转换为数组并排序
  return Object.values(groups)
    .map(chapter => ({
      ...chapter,
      lessons: Object.values(chapter.lessons)
        .sort((a, b) => a.lessonOrder - b.lessonOrder)
    }))
    .sort((a, b) => a.chapterOrder - b.chapterOrder)
})

const getInstructorName = (course) => {
  if (course.instructor_first_name || course.instructor_last_name) {
    const firstName = course.instructor_first_name || ''
    const lastName = course.instructor_last_name || ''
    return `${lastName}${firstName}`.trim() || '未知教师'
  }
  return '未知教师'
}

// 从URL中提取文件名
const getFileNameFromUrl = (url) => {
  if (!url) return '文件'
  try {
    const urlObj = new URL(url)
    const pathname = urlObj.pathname
    const fileName = pathname.split('/').pop() || '文件'
    // 解码URL编码的文件名
    return decodeURIComponent(fileName)
  } catch (e) {
    // 如果不是有效的URL，尝试从路径中提取
    const parts = url.split('/')
    return parts[parts.length - 1] || '文件'
  }
}

// 下载文件
const downloadFile = async (url) => {
  try {
    // 创建一个临时的a标签来触发下载
    const link = document.createElement('a')
    link.href = url
    link.target = '_blank'
    link.download = getFileNameFromUrl(url)
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
  } catch (error) {
    console.error('下载文件失败:', error)
    // 如果直接下载失败，则在新窗口打开
    window.open(url, '_blank')
  }
}

// 加载评论列表
const loadComments = async (courseId) => {
  if (!authStore.isAuthenticated()) return
  
  loadingComments.value = true
  try {
    const res = await getCourseComments(courseId)
    commentsList.value = res || []
  } catch (error) {
    console.error('获取评论列表失败:', error)
    ElMessage.error('获取评论列表失败')
    commentsList.value = []
  } finally {
    loadingComments.value = false
  }
}

// 获取顶级评论（没有父评论的评论）
const topLevelComments = computed(() => {
  return commentsList.value.filter(comment => !comment.parent_comment_id)
    .sort((a, b) => new Date(b.created_at) - new Date(a.created_at))
})

// 获取某个评论的回复
const getReplies = (parentCommentId) => {
  return commentsList.value
    .filter(comment => comment.parent_comment_id === parentCommentId)
    .sort((a, b) => new Date(b.created_at) - new Date(a.created_at))
}

// 发表评论
const handleSubmitComment = async () => {
  if (!commentForm.content.trim()) {
    ElMessage.warning('请输入评论内容')
    return
  }
  
  if (!authStore.isAuthenticated()) {
    ElMessage.warning('请先登录')
    router.push('/login')
    return
  }
  
  submittingComment.value = true
  try {
    await addComment(route.params.id, {
      content: commentForm.content.trim(),
      parent_comment_id: null
    })
    ElMessage.success('评论发表成功')
    commentForm.content = ''
    // 重新加载评论列表
    await loadComments(route.params.id)
  } catch (error) {
    console.error('发表评论失败:', error)
    ElMessage.error(error.response?.data?.message || '发表评论失败')
  } finally {
    submittingComment.value = false
  }
}

// 切换回复表单显示
const toggleReply = (commentId) => {
  if (replyingTo.value === commentId) {
    replyingTo.value = null
    replyForm.content = ''
  } else {
    replyingTo.value = commentId
    replyForm.content = ''
  }
}

// 取消回复
const cancelReply = () => {
  replyingTo.value = null
  replyForm.content = ''
}

// 提交回复
const handleSubmitReply = async (parentCommentId) => {
  if (!replyForm.content.trim()) {
    ElMessage.warning('请输入回复内容')
    return
  }
  
  if (!authStore.isAuthenticated()) {
    ElMessage.warning('请先登录')
    router.push('/login')
    return
  }
  
  submittingReply.value = true
  try {
    await addComment(route.params.id, {
      content: replyForm.content.trim(),
      parent_comment_id: parentCommentId
    })
    ElMessage.success('回复成功')
    replyForm.content = ''
    replyingTo.value = null
    // 重新加载评论列表
    await loadComments(route.params.id)
  } catch (error) {
    console.error('回复失败:', error)
    ElMessage.error(error.response?.data?.message || '回复失败')
  } finally {
    submittingReply.value = false
  }
}

// 格式化评论时间（精确到时分秒）
const formatCommentTime = (timeStr) => {
  if (!timeStr) return ''
  const date = new Date(timeStr)
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  const seconds = String(date.getSeconds()).padStart(2, '0')
  return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`
}

onMounted(() => {
  loadCourse()
})
</script>

<style scoped>
.course-detail-container {
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

.info-card {
  margin-bottom: 20px;
}

.course-meta {
  margin-bottom: 20px;
}

.course-meta p {
  margin: 10px 0;
  color: #666;
}

.actions {
  margin-top: 20px;
  text-align: center;
}

.chapters-card {
  margin-top: 20px;
}

.lesson-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 0;
  border-bottom: 1px solid #eee;
}

.lesson-item:last-child {
  border-bottom: none;
}

.lesson-info {
  display: flex;
  align-items: center;
  gap: 15px;
  flex: 1;
}

.lesson-order {
  font-weight: bold;
  color: #409eff;
  min-width: 30px;
}

.lesson-title {
  flex: 1;
  color: #333;
}

.empty-lessons {
  color: #999;
  text-align: center;
  padding: 20px;
}

.enrollment-notice {
  margin-top: 15px;
}

.enrollment-text {
  color: #909399;
  font-size: 14px;
  margin: 0;
  padding: 8px 12px;
  background-color: #f4f4f5;
  border-radius: 4px;
}

.chapters-list-unenrolled {
  padding: 10px 0;
}

.chapter-item-unenrolled {
  padding: 12px 0;
  border-bottom: 1px solid #eee;
  color: #666;
}

.chapter-item-unenrolled:last-child {
  border-bottom: none;
}

.chapter-title {
  font-size: 14px;
  color: #666;
}

.tasks-card {
  margin-top: 20px;
}

.lesson-task-group {
  margin-bottom: 20px;
}

.lesson-task-group:last-child {
  margin-bottom: 0;
}

.lesson-task-header {
  padding: 12px 0;
  margin-bottom: 10px;
  border-bottom: 2px solid #409eff;
}

.lesson-title-text {
  font-size: 16px;
  font-weight: 600;
  color: #409eff;
}

.task-actions {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.task-actions .el-button {
  margin: 0;
}

/* 评论相关样式 */
.comments-card {
  margin-top: 20px;
}

.comments-card .el-card__header {
  padding: 18px 20px;
  border-bottom: 1px solid #ebeef5;
}

.comments-card .el-card__header h2 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: #303133;
}

.comments-card .el-card__body {
  padding: 20px;
}

.comment-form-section {
  margin-bottom: 30px;
  padding-bottom: 20px;
  border-bottom: 1px solid #e4e7ed;
}

.comment-form-section .el-form-item {
  margin-bottom: 15px;
}

.comment-form-section .el-form-item:last-child {
  margin-bottom: 0;
}

.comments-list {
  margin-top: 20px;
}

.comment-item {
  margin-bottom: 24px;
  padding-bottom: 20px;
  border-bottom: 1px solid #f0f0f0;
  transition: all 0.3s;
}

.comment-item:last-child {
  border-bottom: none;
  margin-bottom: 0;
}

.comment-item:hover {
  background-color: #fafafa;
  padding: 12px;
  margin: -12px -12px 12px -12px;
  border-radius: 6px;
}

.comment-main {
  margin-bottom: 12px;
}

.comment-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 10px;
  flex-wrap: wrap;
}

.comment-author {
  font-weight: 600;
  color: #409eff;
  font-size: 15px;
  display: flex;
  align-items: center;
}

.comment-time {
  font-size: 13px;
  color: #909399;
  white-space: nowrap;
}

.comment-content {
  color: #606266;
  line-height: 1.8;
  margin-bottom: 12px;
  white-space: pre-wrap;
  word-wrap: break-word;
  padding: 12px;
  background-color: #f8f9fa;
  border-radius: 4px;
  border-left: 3px solid #409eff;
}

.comment-actions {
  margin-top: 8px;
  display: flex;
  align-items: center;
}

.comment-actions .el-button {
  padding: 0;
  font-size: 13px;
}

.reply-form {
  margin-top: 15px;
  margin-left: 40px;
  padding: 15px;
  background: #f5f7fa;
  border-radius: 6px;
  border: 1px solid #e4e7ed;
}

.reply-form .el-input {
  margin-bottom: 10px;
}

.reply-form .el-button {
  margin-right: 8px;
}

.replies-list {
  margin-top: 15px;
  margin-left: 40px;
}

.replies-title {
  font-size: 13px;
  color: #909399;
  margin-bottom: 10px;
  font-weight: 500;
}

.reply-item {
  padding: 12px 15px;
  margin-bottom: 12px;
  background: #f8f9fa;
  border-radius: 6px;
  border-left: 3px solid #67c23a;
  transition: all 0.2s;
}

.reply-item:hover {
  background: #f0f2f5;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
}

.reply-item:last-child {
  margin-bottom: 0;
}

.reply-item .comment-header {
  margin-bottom: 8px;
}

.reply-item .comment-author {
  font-size: 14px;
  color: #67c23a;
}

.reply-item .comment-time {
  font-size: 12px;
}

.reply-item .comment-content {
  background-color: transparent;
  border-left: none;
  padding: 0;
  margin-bottom: 0;
  font-size: 14px;
  color: #606266;
}

/* 视频/文档查看器样式 */
.media-viewer-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 400px;
}

.video-container {
  width: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
}

.document-container {
  width: 100%;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
}

.office-document-viewer {
  width: 100%;
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.viewer-tip-bottom {
  width: 100%;
  margin-top: 0;
}

.pdf-viewer {
  width: 100%;
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.viewer-tip-bottom {
  width: 100%;
  margin-top: 0;
}

.image-viewer {
  width: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
}

.unsupported-document {
  width: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 400px;
}
</style>

