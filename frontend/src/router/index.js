import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    name: 'home',
    component: () => import('../views/HomeView.vue')
  },
  {
    path: '/note/:path(.*)',
    name: 'note',
    component: () => import('../views/NoteView.vue')
  },
  {
    path: '/search',
    name: 'search',
    component: () => import('../views/SearchView.vue')
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router
