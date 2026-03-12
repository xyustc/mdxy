import request from './request'
import type { LoginRequest, LoginResponse, ApiResponse } from './types'

export const adminApi = {
  // 登录
  login(data: LoginRequest): Promise<ApiResponse<LoginResponse>> {
    return request.post('/admin/login', data)
  },

  // 登出
  logout() {
    localStorage.removeItem('admin_token')
  },

  // 检查是否已登录
  isAuthenticated(): boolean {
    return !!localStorage.getItem('admin_token')
  }
}
