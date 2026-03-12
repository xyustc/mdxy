import request from './request'
import type { Tool, ToolForm, ApiResponse } from './types'

export const toolApi = {
  list(category?: string): Promise<ApiResponse<Tool[]>> {
    const params = category ? { category } : {}
    return request.get('/tools', { params })
  },

  adminList(category?: string): Promise<ApiResponse<Tool[]>> {
    const params = category ? { category } : {}
    return request.get('/admin/tools', { params })
  },

  getCategories(): Promise<ApiResponse<string[]>> {
    return request.get('/tools/categories')
  },

  getById(id: number): Promise<ApiResponse<Tool>> {
    return request.get(`/admin/tools/${id}`)
  },

  create(data: ToolForm): Promise<ApiResponse<Tool>> {
    return request.post('/admin/tools', data)
  },

  update(id: number, data: ToolForm): Promise<ApiResponse<Tool>> {
    return request.put(`/admin/tools/${id}`, data)
  },

  delete(id: number): Promise<ApiResponse<null>> {
    return request.delete(`/admin/tools/${id}`)
  }
}
