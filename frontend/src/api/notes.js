import axios from 'axios'
import { ensureVisitorId } from '../utils/visitorId'

const api = axios.create({
  baseURL: '/api',
  timeout: 10000
})

// 请求拦截器 - 添加用户标识
api.interceptors.request.use(
  config => {
    const visitorId = ensureVisitorId();
    if (visitorId) {
      config.headers['X-Visitor-ID'] = visitorId;
    }
    return config;
  },
  error => {
    return Promise.reject(error);
  }
);

/**
 * 获取笔记目录树
 */
export async function getNoteTree() {
  const response = await api.get('/notes')
  return response.data
}

/**
 * 获取笔记内容
 * @param {string} path - 笔记路径
 */
export async function getNoteContent(path) {
  const response = await api.get(`/notes/${path}`)
  return response.data
}

/**
 * 搜索笔记
 * @param {string} keyword - 搜索关键词
 */
export async function searchNotes(keyword) {
  const response = await api.get('/notes/search', { params: { q: keyword } })
  return response.data
}
