import request from './request'
import type { Profile, ApiResponse } from './types'

export const profileApi = {
  // 获取个人信息
  get(): Promise<ApiResponse<Profile>> {
    return request.get('/profile')
  },

  // 更新个人信息
  update(data: Partial<Profile>): Promise<ApiResponse<Profile>> {
    return request.put('/admin/profile', data)
  }
}
