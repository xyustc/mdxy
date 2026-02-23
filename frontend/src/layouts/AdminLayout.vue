<template>
  <el-container class="admin-layout">
    <el-aside width="200px" class="sidebar">
      <div class="logo">管理后台</div>
      <el-menu
        :default-active="activeMenu"
        router
        background-color="#001529"
        text-color="#fff"
        active-text-color="#1890ff"
      >
        <el-menu-item index="/admin/dashboard">
          <el-icon><ElIconHouse /></el-icon>
          <span>仪表盘</span>
        </el-menu-item>
        <el-menu-item index="/admin/profile">
          <el-icon><ElIconUser /></el-icon>
          <span>个人信息</span>
        </el-menu-item>
        <el-menu-item index="/admin/analytics">
          <el-icon><ElIconDataAnalysis /></el-icon>
          <span>数据统计</span>
        </el-menu-item>
        <el-menu-item index="/admin/tools">
          <el-icon><ElIconSuitcase /></el-icon>
          <span>工具管理</span>
        </el-menu-item>
      </el-menu>
    </el-aside>

    <el-container>
      <el-header class="header">
        <div class="header-content">
          <span>欢迎回来</span>
          <el-button type="danger" @click="handleLogout">退出登录</el-button>
        </div>
      </el-header>

      <el-main class="main">
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElContainer, ElAside, ElHeader, ElMain, ElMenu, ElMenuItem, ElButton, ElIcon, ElMessage } from 'element-plus'
import { House as ElIconHouse, User as ElIconUser, DataAnalysis as ElIconDataAnalysis, Suitcase as ElIconSuitcase } from '@element-plus/icons-vue'
import { useAdminStore } from '@/stores/admin'

const route = useRoute()
const router = useRouter()
const adminStore = useAdminStore()

const activeMenu = computed(() => route.path)

const handleLogout = () => {
  adminStore.logout()
  ElMessage.success('已退出登录')
  router.push('/admin/login')
}
</script>

<style scoped>
.admin-layout {
  height: 100vh;
}

.sidebar {
  background: #001529;
  overflow-y: auto;
}

.logo {
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  font-size: 18px;
  font-weight: 600;
}

.header {
  background: #fff;
  border-bottom: 1px solid #f0f0f0;
  display: flex;
  align-items: center;
}

.header-content {
  width: 100%;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.main {
  background: #f0f2f5;
  padding: 24px;
}
</style>
