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
