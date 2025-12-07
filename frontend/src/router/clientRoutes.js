import ClientLayout from '../views/client/LayoutView.vue'
import HomeView from '../views/client/HomeView.vue'
import NoteView from '../views/client/NoteView.vue'
import SearchView from '../views/client/SearchView.vue'

// 客户端路由配置
export const clientRoutes = [
  {
    path: '/',
    component: ClientLayout,
    children: [
      {
        path: '',
        name: 'home',
        component: HomeView
      },
      {
        path: 'note/:path(.*)',
        name: 'note',
        component: NoteView
      },
      {
        path: 'search',
        name: 'search',
        component: SearchView
      }
    ]
  }
]