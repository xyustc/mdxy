<template>
  <div class="admin-layout">
    <aside class="admin-sidebar">
      <router-link to="/" class="admin-brand">
        <span class="admin-brand__index">Editor Desk</span>
        <strong>MDXY Admin</strong>
      </router-link>

      <nav class="admin-nav" aria-label="后台导航">
        <router-link
          v-for="item in navItems"
          :key="item.to"
          :to="item.to"
          class="admin-nav__item"
          :class="{ 'admin-nav__item--active': activeMenu === item.to }"
        >
          <component :is="item.icon" class="admin-nav__icon" />
          <div>
            <strong>{{ item.label }}</strong>
            <span>{{ item.caption }}</span>
          </div>
        </router-link>
      </nav>
    </aside>

    <div class="admin-main">
      <header class="admin-header">
        <div>
          <p class="section-kicker">Workspace</p>
          <h1>{{ currentPage.label }}</h1>
          <p>{{ currentPage.caption }}</p>
        </div>

        <div class="admin-header__actions">
          <router-link to="/" class="admin-header__action admin-header__action--secondary">返回前台</router-link>
          <button type="button" class="admin-header__action admin-header__action--primary" @click="handleLogout">退出登录</button>
        </div>
      </header>

      <main class="admin-content">
        <router-view />
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import {
  House as ElIconHouse,
  User as ElIconUser,
  DataAnalysis as ElIconDataAnalysis,
  Suitcase as ElIconSuitcase
} from '@element-plus/icons-vue'
import { useAdminStore } from '@/stores/admin'

const route = useRoute()
const router = useRouter()
const adminStore = useAdminStore()

const navItems = [
  { to: '/admin/dashboard', label: '仪表盘', caption: '站点概览与最近活动', icon: ElIconHouse },
  { to: '/admin/profile', label: '个人信息', caption: '编辑公开资料与栏目内容', icon: ElIconUser },
  { to: '/admin/analytics', label: '数据统计', caption: '查看访问趋势与内容表现', icon: ElIconDataAnalysis },
  { to: '/admin/tools', label: '工具管理', caption: '维护工具箱与展示顺序', icon: ElIconSuitcase }
]

const activeMenu = computed(() => navItems.find((item) => route.path.startsWith(item.to))?.to || '/admin/dashboard')
const currentPage = computed(() => navItems.find((item) => item.to === activeMenu.value) || navItems[0])

const handleLogout = () => {
  adminStore.logout()
  ElMessage.success('已退出登录')
  router.push('/admin/login')
}
</script>

<style scoped>
.admin-layout {
  min-height: 100vh;
  display: grid;
  grid-template-columns: 280px minmax(0, 1fr);
  background: var(--bg-page);
}

.admin-sidebar {
  padding: var(--space-xl);
  border-right: 1px solid var(--border-primary);
  background: linear-gradient(180deg, rgba(15, 18, 20, 0.94), rgba(22, 25, 28, 0.96));
  color: var(--text-inverse);
  display: flex;
  flex-direction: column;
  gap: var(--space-2xl);
}

.admin-brand {
  display: flex;
  flex-direction: column;
  gap: 0.45rem;
  color: var(--text-inverse);
}

.admin-brand__index {
  font-family: var(--font-mono);
  font-size: 0.72rem;
  letter-spacing: 0.16em;
  text-transform: uppercase;
  color: rgba(245, 239, 229, 0.6);
}

.admin-brand strong {
  font-family: var(--font-display);
  font-size: 2rem;
  letter-spacing: -0.04em;
}

.admin-nav {
  display: grid;
  gap: 0.5rem;
}

.admin-nav__item {
  display: grid;
  grid-template-columns: 1rem 1fr;
  gap: var(--space-md);
  padding: 1rem 1rem 1rem 1.1rem;
  border-radius: var(--radius-lg);
  color: rgba(245, 239, 229, 0.72);
  border: 1px solid transparent;
  transition: var(--transition);
}

.admin-nav__item strong {
  display: block;
  font-size: 1rem;
  font-weight: 600;
}

.admin-nav__item span {
  display: block;
  margin-top: 0.2rem;
  font-size: 0.85rem;
  color: rgba(245, 239, 229, 0.5);
  line-height: 1.55;
}

.admin-nav__icon {
  width: 1rem;
  height: 1rem;
  margin-top: 0.2rem;
}

.admin-nav__item:hover,
.admin-nav__item--active {
  background: rgba(255, 255, 255, 0.06);
  border-color: rgba(255, 255, 255, 0.08);
  color: var(--text-inverse);
}

.admin-main {
  min-width: 0;
  display: flex;
  flex-direction: column;
}

.admin-header {
  padding: var(--space-xl) clamp(1.2rem, 3vw, 2.4rem);
  border-bottom: 1px solid var(--border-primary);
  display: flex;
  justify-content: space-between;
  gap: var(--space-lg);
  align-items: flex-end;
  backdrop-filter: blur(16px);
  -webkit-backdrop-filter: blur(16px);
}

.admin-header h1 {
  margin-top: var(--space-sm);
  font-family: var(--font-display);
  font-size: clamp(1.8rem, 3vw, 2.8rem);
  letter-spacing: -0.04em;
}

.admin-header p:last-child {
  margin-top: 0.35rem;
  color: var(--text-secondary);
}

.admin-header__actions {
  display: flex;
  align-items: center;
  gap: var(--space-md);
}

.admin-header__action {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-height: 3rem;
  padding: 0 1.15rem;
  border-radius: var(--radius-pill);
  border: 1px solid var(--border-primary);
  font-weight: 600;
  letter-spacing: -0.01em;
  transition: var(--transition);
}

.admin-header__action--secondary {
  background: var(--bg-panel);
  color: var(--text-primary);
}

.admin-header__action--secondary:hover {
  border-color: var(--border-strong);
  color: var(--text-primary);
}

.admin-header__action--primary {
  background: var(--bg-contrast);
  border-color: var(--bg-contrast);
  color: var(--text-inverse);
}

.admin-header__action--primary:hover {
  opacity: 0.92;
  color: var(--text-inverse);
}

.admin-content {
  padding: clamp(1rem, 3vw, 2rem);
}

@media (max-width: 980px) {
  .admin-layout {
    grid-template-columns: 1fr;
  }

  .admin-sidebar {
    border-right: none;
    border-bottom: 1px solid var(--border-primary);
  }

  .admin-header {
    align-items: flex-start;
    flex-direction: column;
  }
}

@media (max-width: 640px) {
  .admin-header__actions {
    width: 100%;
    justify-content: space-between;
  }

  .admin-nav__item {
    grid-template-columns: 1fr;
  }

  .admin-nav__icon {
    margin-top: 0;
  }
}
</style>
