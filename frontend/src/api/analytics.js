import axios from 'axios'
import { ensureVisitorId } from '../utils/visitorId'

const api = axios.create({
  baseURL: '/api',
  timeout: 10000
})

// 请求拦截器 - 添加token和用户标识
api.interceptors.request.use(
  config => {
    const token = localStorage.getItem('admin_token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    
    const visitorId = ensureVisitorId();
    if (visitorId) {
      config.headers['X-Visitor-ID'] = visitorId;
    }
    
    return config
  },
  error => {
    return Promise.reject(error)
  }
)

/**
 * 获取访问日志列表
 */
export async function getAnalyticsLogs(params) {
  const response = await api.get('/admin/analytics/logs', { params })
  return response.data
}

/**
 * 获取统计概览
 */
export async function getAnalyticsOverview(params) {
  const response = await api.get('/admin/analytics/overview', { params })
  return response.data
}