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

// 响应拦截器 - 处理token过期
api.interceptors.response.use(
  response => response,
  error => {
    if (error.response?.status === 401) {
      // token过期，清除并跳转登录
      localStorage.removeItem('admin_token')
      if (window.location.pathname !== '/admin/login') {
        window.location.href = '/admin/login'
      }
    }
    return Promise.reject(error)
  }
)

/**
 * 管理员登录
 */
export async function login(password) {
  const response = await api.post('/admin/login', { password })
  return response.data
}

/**
 * 验证token
 */
export async function verifyToken() {
  const response = await api.get('/admin/verify')
  return response.data
}

/**
 * 上传笔记
 */
export async function uploadNote(file, path = '') {
  const formData = new FormData()
  formData.append('file', file)
  formData.append('path', path)
  
  const response = await api.post('/admin/notes/upload', formData, {
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
  return response.data
}

/**
 * 删除笔记
 */
export async function deleteNote(path) {
  const response = await api.delete(`/admin/notes/${path}`)
  return response.data
}

/**
 * 创建目录
 */
export async function createDirectory(path) {
  const response = await api.post('/admin/notes/mkdir', { path })
  return response.data
}

/**
 * 移动/重命名文件
 */
export async function moveNote(oldPath, newPath) {
  const response = await api.put('/admin/notes/move', {
    old_path: oldPath,
    new_path: newPath
  })
  return response.data
}