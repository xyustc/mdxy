import { defineStore } from 'pinia'
import { login as loginAPI } from '../api/admin'

export const useAdminStore = defineStore('admin', {
  state: () => ({
    token: localStorage.getItem('admin_token') || null,
    isAuthenticated: !!localStorage.getItem('admin_token')
  }),
  
  actions: {
    async login(password) {
      try {
        const response = await loginAPI(password)
        
        if (response.success) {
          this.token = response.data.token
          this.isAuthenticated = true
          localStorage.setItem('admin_token', this.token)
          return { success: true }
        } else {
          return { success: false, message: response.message }
        }
      } catch (error) {
        return {
          success: false,
          message: error.response?.data?.detail || '登录失败'
        }
      }
    },
    
    logout() {
      this.token = null
      this.isAuthenticated = false
      localStorage.removeItem('admin_token')
    }
  }
})
