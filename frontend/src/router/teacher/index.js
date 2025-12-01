import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../../stores/auth'

const routes = [
  {
    path: '/login',
    name: 'TeacherLogin',
    component: () => import('../../views/teacher/Login.vue'),
    meta: { requiresAuth: false }
  },
  {
    path: '/courses',
    name: 'TeacherCourseList',
    component: () => import('../../views/teacher/CourseList.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/courses/create',
    name: 'CreateCourse',
    component: () => import('../../views/teacher/CreateCourse.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/courses/:id',
    name: 'TeacherCourseDetail',
    component: () => import('../../views/teacher/CourseDetail.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/',
    redirect: '/courses'
  }
]

// 教师端路由使用根路径
const router = createRouter({
  history: createWebHistory('/'),
  routes
})

// 路由守卫
router.beforeEach((to, from, next) => {
  // 如果路由不需要认证，直接通过
  if (to.meta.requiresAuth === false) {
    next()
    return
  }
  
  // 需要认证的路由
  const authStore = useAuthStore()
  if (!authStore.isAuthenticated()) {
    next('/login')
  } else {
    next()
  }
})

export default router

