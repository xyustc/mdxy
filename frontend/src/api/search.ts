import request from './request'
import type { ApiResponse, UnifiedSearchResult, PopularKeyword } from './types'

export const searchApi = {
  search(q: string, searchType = 'all'): Promise<ApiResponse<UnifiedSearchResult>> {
    return request.get('/search', { params: { q, type: searchType } })
  },

  getPopular(): Promise<ApiResponse<PopularKeyword[]>> {
    return request.get('/search/popular')
  }
}
