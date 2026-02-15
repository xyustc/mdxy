import { defineStore } from 'pinia'
import { ref } from 'vue'
import { adminApi } from '@/api/admin'

export const useAdminStore = defineStore('admin', () => {
  const isAuthenticated = ref(adminApi.isAuthenticated())
  const username = ref('')

  const login = async (loginData: { username: string; password: string }) => {
    const response = await adminApi.login(loginData)
    if (response.success && response.data) {
      localStorage.setItem('admin_token', response.data.token)
      isAuthenticated.value = true
      username.value = loginData.username
      return true
    }
    return false
  }

  const logout = () => {
    adminApi.logout()
    isAuthenticated.value = false
    username.value = ''
  }

  return {
    isAuthenticated,
    username,
    login,
    logout
  }
})
