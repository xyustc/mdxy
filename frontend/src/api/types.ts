export interface Profile {
  id: number
  name: string
  title: string
  bio: string
  avatar: string
  email: string
  github: string
  linkedin: string
  twitter: string
  website: string
  skills: string
  created_at: string
  updated_at: string
}

export interface NoteNode {
  name: string
  type: 'file' | 'directory'
  path: string
  children?: NoteNode[]
}

export interface NoteContent {
  path: string
  content: string
}

export interface SearchResult {
  name: string
  path: string
  context: string
}

export interface LoginRequest {
  username: string
  password: string
}

export interface LoginResponse {
  token: string
}

export interface ApiResponse<T = any> {
  success: boolean
  data?: T
  message?: string
  error?: string
}

// Tool 相关
export interface Tool {
  id: number
  name: string
  description: string
  type: string
  url: string
  icon: string
  category: string
  sort_order: number
  is_visible: boolean
  metadata?: string // JSON 扩展字段
  created_at: string
  updated_at: string
}

export interface ToolForm {
  name: string
  description: string
  type: string
  url: string
  icon: string
  category: string
  sort_order: number
  is_visible: boolean
  metadata?: string
}

// Analytics 相关
export interface AnalyticsOverview {
  total_pv: number
  total_uv: number
  today_pv: number
  today_uv: number
  note_count: number
  tool_count: number
}

export interface DailyStat {
  date: string
  count: number
}

export interface PageStat {
  path: string
  count: number
}

export interface DeviceStat {
  device_type: string
  count: number
}

export interface BrowserStat {
  browser: string
  count: number
}

export interface GeoStat {
  country: string
  region: string
  count: number
}

// Search 相关
export interface SearchResultItem {
  type: string
  title: string
  path: string
  context: string
  url?: string
}

export interface SearchGroup {
  type: string
  label: string
  results: SearchResultItem[]
}

export interface UnifiedSearchResult {
  groups: SearchGroup[]
  total: number
}

export interface PopularKeyword {
  keyword: string
  count: number
}
