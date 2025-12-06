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
  },
  // 管理后台路由
  {
    path: '/admin/login',
    name: 'admin-login',
    component: () => import('../views/admin/LoginView.vue')
  },
  {
    path: '/admin',
    component: () => import('../views/admin/LayoutView.vue'),
    meta: { requiresAuth: true },
    children: [
      {
        path: '',
        redirect: '/admin/dashboard'
      },
      {
        path: 'dashboard',
        name: 'admin-dashboard',
        component: () => import('../views/admin/DashboardView.vue')
      },
      {
        path: 'files',
        name: 'admin-files',
        component: () => import('../views/admin/FileManagerView.vue')
      },
      {
        path: 'analytics',
        name: 'admin-analytics',
        component: () => import('../views/admin/AnalyticsView.vue')
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫 - 验证管理员认证
router.beforeEach((to, from, next) => {
  if (to.meta.requiresAuth) {
    const token = localStorage.getItem('admin_token')
    if (!token) {
      next({ name: 'admin-login' })
    } else {
      next()
    }
  } else {
    next()
  }
})

export default router
