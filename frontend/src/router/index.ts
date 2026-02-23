import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router'
import { adminApi } from '@/api/admin'

// 前台路由
const clientRoutes: RouteRecordRaw[] = [
  {
    path: '/',
    component: () => import('@/layouts/ClientLayout.vue'),
    children: [
      {
        path: '',
        name: 'home',
        component: () => import('@/views/client/HomePage.vue')
      },
      {
        path: 'notes',
        component: () => import('@/views/client/NotesPage.vue'),
        children: [
          {
            path: '',
            name: 'notes',
            component: () => import('@/views/client/NotesWelcome.vue')
          },
          {
            path: ':path(.*)',
            name: 'note-detail',
            component: () => import('@/views/client/NoteDetail.vue')
          }
        ]
      },
      {
        path: 'tools',
        name: 'tools',
        component: () => import('@/views/client/ToolsPage.vue')
      }
    ]
  }
]

// 后台路由
const adminRoutes: RouteRecordRaw[] = [
  {
    path: '/admin/login',
    name: 'admin-login',
    component: () => import('@/views/admin/LoginPage.vue')
  },
  {
    path: '/admin',
    component: () => import('@/layouts/AdminLayout.vue'),
    meta: { requiresAuth: true },
    children: [
      {
        path: '',
        redirect: '/admin/dashboard'
      },
      {
        path: 'dashboard',
        name: 'admin-dashboard',
        component: () => import('@/views/admin/DashboardPage.vue')
      },
      {
        path: 'profile',
        name: 'admin-profile',
        component: () => import('@/views/admin/ProfileEdit.vue')
      },
      {
        path: 'analytics',
        name: 'admin-analytics',
        component: () => import('@/views/admin/AnalyticsPage.vue')
      },
      {
        path: 'tools',
        name: 'admin-tools',
        component: () => import('@/views/admin/ToolsManage.vue')
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes: [...clientRoutes, ...adminRoutes]
})

// 路由守卫
router.beforeEach((to, from, next) => {
  if (to.meta.requiresAuth && !adminApi.isAuthenticated()) {
    next({ name: 'admin-login', query: { redirect: to.fullPath } })
  } else {
    next()
  }
})

export default router
