import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import NoteView from '../views/NoteView.vue'

const routes = [
  {
    path: '/',
    name: 'home',
    component: HomeView
  },
  {
    path: '/note/:path(.*)',
    name: 'note',
    component: NoteView  // 改为直接导入，避免首次点击卡顿
  },
  {
    path: '/search',
    name: 'search',
    component: () => import('../views/SearchView.vue')  // 搜索页面不常用，保持懒加载
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router
