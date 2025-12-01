import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('token') || '')
  
  // 安全地解析 userInfo
  let parsedUserInfo = null
  try {
    const userInfoStr = localStorage.getItem('userInfo')
    if (userInfoStr && userInfoStr !== 'undefined' && userInfoStr !== 'null') {
      parsedUserInfo = JSON.parse(userInfoStr)
    }
  } catch (e) {
    console.warn('Failed to parse userInfo from localStorage:', e)
    parsedUserInfo = null
  }
  const userInfo = ref(parsedUserInfo)
  
  const role = ref(localStorage.getItem('role') || '')

  const setToken = (newToken) => {
    token.value = newToken
    localStorage.setItem('token', newToken)
  }

  const setUserInfo = (info) => {
    userInfo.value = info
    localStorage.setItem('userInfo', JSON.stringify(info))
  }

  const setRole = (newRole) => {
    role.value = newRole
    localStorage.setItem('role', newRole)
  }

  const logout = () => {
    token.value = ''
    userInfo.value = null
    role.value = ''
    localStorage.removeItem('token')
    localStorage.removeItem('userInfo')
    localStorage.removeItem('role')
  }

  const isAuthenticated = () => {
    return !!token.value
  }

  return {
    token,
    userInfo,
    role,
    setToken,
    setUserInfo,
    setRole,
    logout,
    isAuthenticated
  }
})

