import AdminLayout from '../views/admin/LayoutView.vue'
import AdminLogin from '../views/admin/LoginView.vue'
import AdminDashboard from '../views/admin/DashboardView.vue'
import FileManager from '../views/admin/FileManagerView.vue'
import Analytics from '../views/admin/AnalyticsView.vue'

// 管理端路由配置
export const adminRoutes = [
  {
    path: '/admin/login',
    name: 'admin-login',
    component: AdminLogin
  },
  {
    path: '/admin',
    component: AdminLayout,
    meta: { requiresAuth: true },
    children: [
      {
        path: '',
        redirect: '/admin/dashboard'
      },
      {
        path: 'dashboard',
        name: 'admin-dashboard',
        component: AdminDashboard
      },
      {
        path: 'files',
        name: 'admin-files',
        component: FileManager
      },
      {
        path: 'analytics',
        name: 'admin-analytics',
        component: Analytics
      }
    ]
  }
]