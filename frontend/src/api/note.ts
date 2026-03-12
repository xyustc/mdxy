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
  }
}
