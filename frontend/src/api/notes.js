import axios from 'axios'

const api = axios.create({
  baseURL: '/api',
  timeout: 10000
})

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
