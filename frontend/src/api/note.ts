import request from './request'
import type { NoteNode, NoteContent, SearchResult, ApiResponse } from './types'

export const noteApi = {
  // 获取笔记目录树
  getTree(): Promise<ApiResponse<NoteNode[]>> {
    return request.get('/notes/tree')
  },

  // 获取笔记内容
  getContent(path: string): Promise<ApiResponse<NoteContent>> {
    return request.get(`/notes/content/${path}`)
  },

  // 搜索笔记
  search(keyword: string): Promise<ApiResponse<{ results: SearchResult[]; total: number }>> {
    return request.get('/notes/search', { params: { q: keyword } })
  },

  // Admin: 获取目录树
  adminGetTree(): Promise<ApiResponse<NoteNode[]>> {
    return request.get('/admin/notes/tree')
  },

  // Admin: 获取原始内容
  adminGetContent(path: string): Promise<ApiResponse<NoteContent>> {
    return request.get(`/admin/notes/content/${path}`)
  },

  // Admin: 保存笔记
  adminSave(path: string, content: string): Promise<ApiResponse<null>> {
    return request.post('/admin/notes/save', { path, content })
  },

  // Admin: 创建目录
  adminCreateDir(path: string): Promise<ApiResponse<null>> {
    return request.post('/admin/notes/directory', { path })
  },

  // Admin: 设置 featured
  adminSetFeatured(path: string, featured: boolean): Promise<ApiResponse<null>> {
    return request.post('/admin/notes/featured', { path, featured })
  },

  // Admin: 删除笔记
  adminDelete(path: string): Promise<ApiResponse<null>> {
    return request.delete(`/admin/notes/${path}`)
  }
}
