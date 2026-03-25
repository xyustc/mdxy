import request from './request'
import type { ApiResponse, CheatSheetPublishedPayload, CheatSheetSnapshotSummary, CheatSheetSyncResult } from './types'

export const cheatSheetApi = {
  getPublished(slug = 'claude-code'): Promise<ApiResponse<CheatSheetPublishedPayload>> {
    return request.get(`/cheat-sheets/${slug}`)
  },

  adminListSnapshots(slug = 'claude-code', limit = 10): Promise<ApiResponse<CheatSheetSnapshotSummary[]>> {
    return request.get(`/admin/cheat-sheets/${slug}/snapshots`, {
      params: { limit }
    })
  },

  adminSync(slug = 'claude-code'): Promise<ApiResponse<CheatSheetSyncResult>> {
    return request.post(`/admin/cheat-sheets/${slug}/sync`)
  },

  adminPublish(slug = 'claude-code', id: number): Promise<ApiResponse<CheatSheetSnapshotSummary>> {
    return request.post(`/admin/cheat-sheets/${slug}/publish/${id}`)
  }
}
