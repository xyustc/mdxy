import { createRouter, createWebHistory } from 'vue-router'
import { clientRoutes } from './clientRoutes'
import { adminRoutes } from './adminRoutes'

// 合并所有路由
const routes = [
  ...clientRoutes,
  ...adminRoutes
]

const router = createRouter({
  history: createWebHistory(),  // 确保使用相同的 history 模式
  routes
})

// 路由守卫 - 验证管理员登录状态
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