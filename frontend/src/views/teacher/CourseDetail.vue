<template>
  <div class="course-detail-container">
    <div class="header">
      <el-button @click="goBack" icon="ArrowLeft">返回</el-button>
      <h1 v-if="course">{{ course.course_title }}</h1>
      <div></div>
    </div>
    
    <div v-loading="loading" class="content">
      <div v-if="course">
        <el-card class="info-card">
          <div class="course-meta">
            <p><strong>课程描述：</strong>{{ course.description || '暂无描述' }}</p>
            <p><strong>状态：</strong>
              <el-tag :type="course.status === 'active' ? 'success' : 'info'">
                {{ course.status === 'active' ? '进行中' : '已结束' }}
              </el-tag>
            </p>
            <p><strong>开课时间：</strong>{{ formatDate(course.start_date) }}</p>
            <p><strong>结课时间：</strong>{{ formatDate(course.end_date) }}</p>
          </div>
        </el-card>
        
        <el-card class="chapters-card">
          <template #header>
            <div class="card-header">
              <h2>课程章节</h2>
              <el-button type="primary" size="small" @click="openAddChapterDialog">
                添加章节
              </el-button>
            </div>
          </template>
          
          <el-collapse v-if="course.chapters && course.chapters.length > 0">
            <el-collapse-item
              v-for="chapter in course.chapters"
              :key="chapter.chapter_id"
              :title="`${chapter.chapter_order || ''}. ${chapter.chapter_title}`"
            >
              <div class="chapter-actions">
                <el-button type="primary" size="small" @click="openAddLessonDialog(chapter)">
                  添加课时
                </el-button>
              </div>
              
              <div v-if="chapter.lessons && chapter.lessons.length > 0" class="lessons-list">
                <div
                  v-for="lesson in chapter.lessons"
                  :key="lesson.lesson_id"
                  class="lesson-item"
                >
                  <div class="lesson-info">
                    <span class="lesson-order">{{ lesson.lesson_order }}</span>
                    <span class="lesson-title">{{ lesson.lesson_title }}</span>
                  </div>
                  <div class="lesson-actions">
                    <el-button
                      v-if="lesson.content_url"
                      type="primary"
                      size="small"
                      @click="playVideo(lesson.content_url)"
                    >
                      {{ lesson.lesson_type === 'video' ? '查看视频' : '查看文档' }}
                    </el-button>
                    <el-button
                      v-if="!hasTask(lesson.lesson_id)"
                      type="success"
                      size="small"
                      @click="openAddTaskDialog(lesson)"
                    >
                      添加任务
                    </el-button>
                    <el-button
                      v-if="hasTask(lesson.lesson_id)"
                      type="warning"
                      size="small"
                      @click="viewTaskAnswersByLesson(lesson)"
                    >
                      查看学生作业
                    </el-button>
                  </div>
                </div>
              </div>
              <p v-else class="empty-lessons">暂无课时</p>
            </el-collapse-item>
          </el-collapse>
          <p v-else class="empty-chapters">暂无章节，点击添加章节开始</p>
        </el-card>

        <!-- 课程评论 -->
        <el-card class="comments-card" style="margin-top: 20px">
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
              <el-empty v-else description="暂无评论" />
            </div>
          </div>
        </el-card>
      </div>
      
      <el-empty v-if="!loading && !course" description="课程不存在" />
    </div>
    
    <!-- 添加章节对话框 -->
    <el-dialog v-model="showAddChapterDialog" title="添加章节" width="600px">
      <el-form ref="chapterFormRef" :model="chapterForm" :rules="chapterRules" label-width="100px">
        <el-form-item label="章节标题" prop="chapter_title">
          <el-input v-model="chapterForm.chapter_title" placeholder="请输入章节标题" />
        </el-form-item>
        <el-form-item label="章节顺序" prop="chapter_order">
          <el-input-number v-model="chapterForm.chapter_order" :min="1" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showAddChapterDialog = false">取消</el-button>
        <el-button type="primary" @click="handleAddChapter" :loading="addingChapter">
          确定
        </el-button>
      </template>
    </el-dialog>
    
    <!-- 添加任务对话框 -->
    <el-dialog v-model="showAddTaskDialog" title="添加任务" width="600px">
      <el-form ref="taskFormRef" :model="taskForm" :rules="taskRules" label-width="100px">
        <el-form-item label="任务标题" prop="task_title">
          <el-input v-model="taskForm.task_title" placeholder="请输入任务标题" />
        </el-form-item>
        <el-form-item label="任务描述" prop="description">
          <el-input
            v-model="taskForm.description"
            type="textarea"
            :rows="4"
            placeholder="请输入任务描述"
          />
        </el-form-item>
        <el-form-item label="任务类型" prop="task_type">
          <el-select v-model="taskForm.task_type" placeholder="选择任务类型" style="width: 100%">
            <el-option label="作文 (essay) - 学生提交文本作业" value="essay" />
            <el-option label="上传 (upload) - 学生上传文件（图片/文档）" value="upload" />
          </el-select>
          <div style="margin-top: 5px; font-size: 12px; color: #909399">
            <p style="margin: 2px 0;">• <strong>作文</strong>：学生提交文本形式的作业内容</p>
            <p style="margin: 2px 0;">• <strong>上传</strong>：学生上传图片或文档文件作为作业（支持 jpg、png、gif 图片格式，单个文件不超过 10MB）</p>
          </div>
        </el-form-item>
        <el-form-item label="最高分" prop="max_score">
          <el-input-number v-model="taskForm.max_score" :min="1" :max="1000" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showAddTaskDialog = false">取消</el-button>
        <el-button type="primary" @click="handleAddTask" :loading="addingTask">
          确定
        </el-button>
      </template>
    </el-dialog>

    <!-- 查看作业对话框 -->
    <el-dialog v-model="showAnswersDialog" :title="currentTaskInfo ? `作业列表 - ${currentTaskInfo.task_title}` : '作业列表'" width="900px">
      <div v-loading="loadingAnswers">
        <!-- 显示任务信息（满分） -->
        <div v-if="currentTaskInfo" class="task-info-bar">
          <el-alert
            :title="`满分：${currentTaskInfo.max_score} 分`"
            type="info"
            :closable="false"
            show-icon
            style="margin-bottom: 20px;"
          />
        </div>
        <div v-if="answersList.length > 0">
          <!-- 未评分作业 -->
          <div v-if="ungradedAnswers.length > 0" class="answer-section">
            <h3 style="margin: 0 0 15px 0; color: #e6a23c; font-size: 16px;">
              <el-icon><Warning /></el-icon>
              待评分 ({{ ungradedAnswers.length }})
            </h3>
            <el-table :data="ungradedAnswers" style="width: 100%" border>
              <el-table-column label="学生姓名" width="150">
                <template #default="{ row }">
                  <span>{{ getStudentName(row) }}</span>
                </template>
              </el-table-column>
              <el-table-column prop="answer_content" label="作业内容" min-width="200">
                <template #default="{ row }">
                  <div v-if="row.type === 'image_url' || row.type === 'upload'">
                    <el-link
                      :underline="false"
                      type="primary"
                      @click="downloadFile(row.answer_content)"
                      style="cursor: pointer;"
                    >
                      <el-icon style="margin-right: 4px;"><Download /></el-icon>
                      下载文件
                    </el-link>
                    <div style="margin-top: 4px; font-size: 12px; color: #909399;">
                      {{ getFileNameFromUrl(row.answer_content) }}
                    </div>
                  </div>
                  <div v-else class="answer-text-content">
                    {{ row.answer_content }}
                  </div>
                </template>
              </el-table-column>
              <el-table-column label="操作" width="280" fixed="right">
                <template #default="{ row }">
                  <div class="grade-actions">
                    <div style="display: flex; align-items: center; gap: 8px;">
                      <span style="font-size: 12px; color: #909399;">分数：</span>
                      <el-input-number
                        v-model="row.gradeScore"
                        :min="0"
                        :max="currentTaskInfo ? currentTaskInfo.max_score : 1000"
                        size="small"
                        style="width: 100px;"
                        @change="(val) => handleScoreChange(row, val)"
                      />
                      <span style="font-size: 12px; color: #909399;">/ {{ currentTaskInfo ? currentTaskInfo.max_score : 100 }}</span>
                    </div>
                    <el-button
                      type="primary"
                      size="small"
                      @click="handleGradeAnswer(row.answer_id, row.branch_id, row.gradeScore)"
                      :disabled="!row.gradeScore || row.gradeScore <= 0"
                      style="margin-top: 8px;"
                    >
                      评分
                    </el-button>
                  </div>
                </template>
              </el-table-column>
            </el-table>
          </div>

          <!-- 已评分作业 -->
          <div v-if="gradedAnswers.length > 0" class="answer-section" :style="{ marginTop: ungradedAnswers.length > 0 ? '30px' : '0' }">
            <h3 style="margin: 0 0 15px 0; color: #67c23a; font-size: 16px;">
              <el-icon><CircleCheck /></el-icon>
              已评分 ({{ gradedAnswers.length }})
            </h3>
            <el-table :data="gradedAnswers" style="width: 100%" border>
              <el-table-column label="学生姓名" width="150">
                <template #default="{ row }">
                  <span>{{ getStudentName(row) }}</span>
                </template>
              </el-table-column>
              <el-table-column prop="answer_content" label="作业内容" min-width="200">
                <template #default="{ row }">
                  <div v-if="row.type === 'image_url' || row.type === 'upload'">
                    <el-link
                      :underline="false"
                      type="primary"
                      @click="downloadFile(row.answer_content)"
                      style="cursor: pointer;"
                    >
                      <el-icon style="margin-right: 4px;"><Download /></el-icon>
                      下载文件
                    </el-link>
                    <div style="margin-top: 4px; font-size: 12px; color: #909399;">
                      {{ getFileNameFromUrl(row.answer_content) }}
                    </div>
                  </div>
                  <div v-else class="answer-text-content">
                    {{ row.answer_content }}
                  </div>
                </template>
              </el-table-column>
              <el-table-column prop="score" label="分数" width="120">
                <template #default="{ row }">
                  <span style="color: #67c23a; font-weight: bold;">{{ row.score }}</span>
                  <span style="color: #909399; font-size: 12px; margin-left: 4px;">
                    / {{ currentTaskInfo ? currentTaskInfo.max_score : 100 }}
                  </span>
                </template>
              </el-table-column>
              <el-table-column label="状态" width="100">
                <template #default="{ row }">
                  <el-tag type="success" size="small">已评分</el-tag>
                </template>
              </el-table-column>
            </el-table>
          </div>
        </div>
        <el-empty v-else description="暂无作业提交" />
      </div>
    </el-dialog>

    <!-- 添加课时对话框 -->
    <el-dialog v-model="showAddLessonDialog" title="添加课时" width="600px">
      <el-form ref="lessonFormRef" :model="lessonForm" :rules="lessonRules" label-width="100px">
        <el-form-item label="课时标题" prop="lesson_title">
          <el-input v-model="lessonForm.lesson_title" placeholder="请输入课时标题" />
        </el-form-item>
        <el-form-item label="课时顺序" prop="lesson_order">
          <el-input-number v-model="lessonForm.lesson_order" :min="1" />
        </el-form-item>
        <el-form-item label="课时类型" prop="lesson_type">
          <el-select v-model="lessonForm.lesson_type" placeholder="选择课时类型" style="width: 100%">
            <el-option label="视频" value="video" />
            <el-option label="文档" value="document" />
          </el-select>
        </el-form-item>
        <el-form-item label="文件" prop="video_file">
          <el-upload
            :auto-upload="false"
            :on-change="handleVideoChange"
            :file-list="videoFileList"
            :accept="lessonForm.lesson_type === 'video' ? 'video/*' : '.pdf,.doc,.docx,.txt'"
          >
            <el-button type="primary">
              {{ lessonForm.lesson_type === 'video' ? '选择视频' : '选择文档' }}
            </el-button>
          </el-upload>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showAddLessonDialog = false">取消</el-button>
        <el-button type="primary" @click="handleAddLesson" :loading="addingLesson">
          确定
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, reactive, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Warning, CircleCheck, Download, ChatLineRound } from '@element-plus/icons-vue'
import { getCourse, createChapter, createLesson, createTask, getCourseTasks, getTaskAnswers, gradeAnswer, getCourseComments, addComment } from '../../api/teacher'

const route = useRoute()
const router = useRouter()

const loading = ref(false)
const course = ref(null)
const showAddChapterDialog = ref(false)
const showAddLessonDialog = ref(false)
const showAddTaskDialog = ref(false)
const showAnswersDialog = ref(false)
const addingChapter = ref(false)
const addingLesson = ref(false)
const addingTask = ref(false)
const loadingAnswers = ref(false)
const currentChapter = ref(null)
const currentLesson = ref(null)
const currentTaskId = ref(null)
const currentTaskInfo = ref(null)
const videoFileList = ref([])
const answersList = ref([])
const lessonTaskMap = ref({}) // 课时ID到任务的映射

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

const chapterFormRef = ref(null)
const lessonFormRef = ref(null)
const taskFormRef = ref(null)

const chapterForm = reactive({
  chapter_title: '',
  chapter_order: 1
})

const lessonForm = reactive({
  lesson_title: '',
  lesson_order: 1,
  lesson_type: 'video',
  video_file: null
})

const chapterRules = {
  chapter_title: [{ required: true, message: '请输入章节标题', trigger: 'blur' }],
  chapter_order: [{ required: true, message: '请输入章节顺序', trigger: 'blur' }]
}

const lessonRules = {
  lesson_title: [{ required: true, message: '请输入课时标题', trigger: 'blur' }],
  lesson_order: [{ required: true, message: '请输入课时顺序', trigger: 'blur' }],
  lesson_type: [{ required: true, message: '请选择课时类型', trigger: 'change' }]
}

const taskForm = reactive({
  task_title: '',
  description: '',
  task_type: 'essay',
  max_score: 100
})

const taskRules = {
  task_title: [{ required: true, message: '请输入任务标题', trigger: 'blur' }],
  task_type: [{ required: true, message: '请选择任务类型', trigger: 'change' }],
  max_score: [{ required: true, message: '请输入最高分', trigger: 'blur' }]
}

const loadCourse = async () => {
  loading.value = true
  try {
    const courseId = route.params.id
    const res = await getCourse(courseId)
    course.value = res
    
    // 加载课程的所有任务，建立课时ID到任务的映射
    await loadCourseTasks(courseId)
    // 加载评论
    await loadComments(courseId)
  } catch (error) {
    console.error('获取课程详情失败:', error)
    ElMessage.error('获取课程详情失败')
  } finally {
    loading.value = false
  }
}

const loadCourseTasks = async (courseId) => {
  try {
    const tasks = await getCourseTasks(courseId)
    // 建立课时ID到任务的映射
    const map = {}
    if (tasks && Array.isArray(tasks)) {
      tasks.forEach(task => {
        map[task.lesson_id] = task
      })
    }
    lessonTaskMap.value = map
  } catch (error) {
    console.warn('获取任务列表失败:', error)
    lessonTaskMap.value = {}
  }
}

const hasTask = (lessonId) => {
  return !!lessonTaskMap.value[lessonId]
}

const openAddChapterDialog = () => {
  // 自动计算下一个章节顺序
  if (course.value && course.value.chapters && course.value.chapters.length > 0) {
    const maxOrder = Math.max(...course.value.chapters.map(ch => ch.chapter_order || 0))
    chapterForm.chapter_order = maxOrder + 1
  } else {
    chapterForm.chapter_order = 1
  }
  chapterForm.chapter_title = ''
  showAddChapterDialog.value = true
}

const handleAddChapter = async () => {
  if (!chapterFormRef.value) return
  
  await chapterFormRef.value.validate(async (valid) => {
    if (valid) {
      addingChapter.value = true
      try {
        await createChapter(route.params.id, chapterForm)
        ElMessage.success('章节添加成功')
        showAddChapterDialog.value = false
        chapterForm.chapter_title = ''
        chapterForm.chapter_order = 1
        loadCourse()
      } catch (error) {
        console.error('添加章节失败:', error)
      } finally {
        addingChapter.value = false
      }
    }
  })
}

const openAddLessonDialog = (chapter) => {
  currentChapter.value = chapter
  showAddLessonDialog.value = true
  lessonForm.lesson_title = ''
  
  // 自动计算下一个课时顺序
  if (chapter.lessons && chapter.lessons.length > 0) {
    const maxOrder = Math.max(...chapter.lessons.map(lesson => lesson.lesson_order || 0))
    lessonForm.lesson_order = maxOrder + 1
  } else {
    lessonForm.lesson_order = 1
  }
  
  lessonForm.lesson_type = 'video'
  lessonForm.video_file = null
  videoFileList.value = []
}

const handleVideoChange = (file) => {
  lessonForm.video_file = file.raw
}

const handleAddLesson = async () => {
  if (!lessonFormRef.value) return
  
  await lessonFormRef.value.validate(async (valid) => {
    if (valid) {
      if (!lessonForm.video_file) {
        ElMessage.warning('请选择视频文件')
        return
      }
      
      addingLesson.value = true
      try {
        const formData = new FormData()
        formData.append('lesson_title', lessonForm.lesson_title)
        formData.append('lesson_order', lessonForm.lesson_order)
        formData.append('lesson_type', lessonForm.lesson_type)
        formData.append('video_file', lessonForm.video_file)
        
        await createLesson(route.params.id, currentChapter.value.chapter_id, formData)
        ElMessage.success('课时添加成功')
        showAddLessonDialog.value = false
        lessonForm.lesson_title = ''
        lessonForm.lesson_order = 1
        lessonForm.lesson_type = 'video'
        lessonForm.video_file = null
        videoFileList.value = []
        loadCourse()
      } catch (error) {
        console.error('添加课时失败:', error)
      } finally {
        addingLesson.value = false
      }
    }
  })
}

const playVideo = (url) => {
  window.open(url, '_blank')
}

const goBack = () => {
  router.push('/courses')
}

const openAddTaskDialog = (lesson) => {
  currentLesson.value = lesson
  showAddTaskDialog.value = true
  taskForm.task_title = ''
  taskForm.description = ''
  taskForm.task_type = 'essay'
  taskForm.max_score = 100
}

const handleAddTask = async () => {
  if (!taskFormRef.value) return
  
  await taskFormRef.value.validate(async (valid) => {
    if (valid) {
      addingTask.value = true
      try {
        await createTask(currentLesson.value.lesson_id, taskForm)
        ElMessage.success('任务创建成功')
        showAddTaskDialog.value = false
        taskForm.task_title = ''
        taskForm.description = ''
        taskForm.task_type = 'essay'
        taskForm.max_score = 100
        // 重新加载任务列表
        await loadCourseTasks(route.params.id)
      } catch (error) {
        console.error('创建任务失败:', error)
        if (error.response?.data?.code === 3006) {
          ElMessage.error('该课时已有任务，每个课时只能有一个任务')
        } else {
          ElMessage.error(error.response?.data?.message || '创建任务失败')
        }
      } finally {
        addingTask.value = false
      }
    }
  })
}

// 根据课时直接查看学生作业（跳过任务列表）
const viewTaskAnswersByLesson = async (lesson) => {
  currentLesson.value = lesson
  loadingAnswers.value = true
  
  try {
    // 获取该课时的任务
    const task = lessonTaskMap.value[lesson.lesson_id]
    if (!task) {
      ElMessage.warning('该课时没有任务')
      return
    }
    
    currentTaskId.value = task.task_id
    currentTaskInfo.value = task
    showAnswersDialog.value = true
    
    // 获取学生提交记录
    const res = await getTaskAnswers(task.task_id)
    answersList.value = res.map(answer => ({
      ...answer,
      gradeScore: answer.is_graded ? answer.score : 0 // 已评分的显示原分数，未评分的显示0
    }))
  } catch (error) {
    console.error('获取作业列表失败:', error)
    ElMessage.error('获取作业列表失败')
    answersList.value = []
  } finally {
    loadingAnswers.value = false
  }
}

// 计算属性：未评分的作业
const ungradedAnswers = computed(() => {
  return answersList.value.filter(answer => !answer.is_graded)
})

// 计算属性：已评分的作业
const gradedAnswers = computed(() => {
  return answersList.value.filter(answer => answer.is_graded)
})

// 处理分数变化，确保不超过最大值
const handleScoreChange = (row, val) => {
  if (!currentTaskInfo) return
  const maxScore = currentTaskInfo.max_score
  if (val > maxScore) {
    row.gradeScore = maxScore
    ElMessage.warning(`分数不能超过满分 ${maxScore} 分，已自动调整为 ${maxScore} 分`)
  }
}

const handleGradeAnswer = async (answerId, branchId, score) => {
  if (!score || score < 0) {
    ElMessage.warning('请输入有效的分数')
    return
  }
  
  // 验证分数不超过满分
  if (currentTaskInfo && score > currentTaskInfo.max_score) {
    ElMessage.warning(`分数不能超过满分 ${currentTaskInfo.max_score} 分`)
    return
  }
  
  try {
    await gradeAnswer(answerId, branchId, score)
    ElMessage.success('评分成功')
    // 更新本地数据
    const answer = answersList.value.find(a => a.answer_id === answerId)
    if (answer) {
      answer.score = score
      answer.is_graded = true
      answer.gradeScore = score // 同步更新gradeScore
    }
    // 重新获取作业列表以更新排序
    if (currentTaskId.value) {
      const res = await getTaskAnswers(currentTaskId.value)
      answersList.value = res.map(a => ({
        ...a,
        gradeScore: a.is_graded ? a.score : 0
      }))
    }
  } catch (error) {
    console.error('评分失败:', error)
    ElMessage.error(error.response?.data?.message || '评分失败')
  }
}

const getTaskTypeName = (type) => {
  const typeMap = {
    essay: '作文',
    upload: '上传'
  }
  return typeMap[type] || type
}

const getStudentName = (answer) => {
  if (answer.student_first_name || answer.student_last_name) {
    const firstName = answer.student_first_name || ''
    const lastName = answer.student_last_name || ''
    return `${lastName}${firstName}`.trim() || `学生ID: ${answer.user_id}`
  }
  return `学生ID: ${answer.user_id}`
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

// 加载评论列表
const loadComments = async (courseId) => {
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

.course-meta p {
  margin: 10px 0;
  color: #666;
}

.chapters-card {
  margin-top: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-header h2 {
  margin: 0;
}

.chapter-actions {
  margin-bottom: 15px;
}

.lessons-list {
  margin-top: 15px;
}

.lesson-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 0;
  border-bottom: 1px solid #eee;
}

.lesson-actions {
  display: flex;
  gap: 10px;
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

.empty-lessons,
.empty-chapters {
  color: #999;
  text-align: center;
  padding: 20px;
}

.lesson-actions {
  display: flex;
  gap: 10px;
}

.grade-actions {
  display: flex;
  flex-direction: column;
  gap: 8px;
  align-items: flex-start;
}

.grade-actions .el-input-number {
  flex-shrink: 0;
}

.grade-actions .el-button {
  flex-shrink: 0;
  margin: 0;
}

.task-info-bar {
  margin-bottom: 10px;
}

.answer-section {
  margin-bottom: 20px;
}

.answer-section h3 {
  display: flex;
  align-items: center;
  gap: 8px;
}

.answer-text-content {
  max-width: 400px;
  word-wrap: break-word;
  word-break: break-all;
  white-space: pre-wrap;
  line-height: 1.6;
  color: #333;
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
</style>

