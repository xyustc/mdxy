import request from './request'
import type { ApiResponse, NoteNode, Tool } from './types'

export interface HomeData {
  featured_notes: NoteNode[]
  note_count: number
  featured_tools: Tool[]
  tool_count: number
}

export const homeApi = {
  get(): Promise<ApiResponse<HomeData>> {
    return request.get('/home')
  }
}
