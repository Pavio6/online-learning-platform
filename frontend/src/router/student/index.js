import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../../stores/auth'

const routes = [
  {
    path: '/login',
    name: 'StudentLogin',
    component: () => import('../../views/student/Login.vue'),
    meta: { requiresAuth: false }
  },
  {
    path: '/register',
    name: 'StudentRegister',
    component: () => import('../../views/student/Register.vue'),
    meta: { requiresAuth: false }
  },
  {
    path: '/courses',
    name: 'CourseList',
    component: () => import('../../views/student/CourseList.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/my-courses',
    name: 'MyCourses',
    component: () => import('../../views/student/MyCourses.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/courses/:id',
    name: 'CourseDetail',
    component: () => import('../../views/student/CourseDetail.vue'),
    meta: { requiresAuth: false }
  },
  {
    path: '/',
    redirect: '/courses'
  }
]

// 学生端路由使用根路径
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

