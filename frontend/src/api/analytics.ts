import request from './request'
import type { ApiResponse, AnalyticsOverview, DailyStat, PageStat, DeviceStat, BrowserStat, GeoStat } from './types'

export const analyticsApi = {
  getOverview(): Promise<ApiResponse<AnalyticsOverview>> {
    return request.get('/admin/analytics/overview')
  },

  getTrends(days = 30): Promise<ApiResponse<{ pv: DailyStat[]; uv: DailyStat[] }>> {
    return request.get('/admin/analytics/trends', { params: { days } })
  },

  getPopularPages(limit = 10): Promise<ApiResponse<PageStat[]>> {
    return request.get('/admin/analytics/popular', { params: { limit } })
  },

  getDevices(): Promise<ApiResponse<DeviceStat[]>> {
    return request.get('/admin/analytics/devices')
  },

  getBrowsers(): Promise<ApiResponse<BrowserStat[]>> {
    return request.get('/admin/analytics/browsers')
  },

  getGeo(): Promise<ApiResponse<GeoStat[]>> {
    return request.get('/admin/analytics/geo')
  }
}
