<template>
  <div class="dashboard-page">
    <h1 class="page-title">仪表盘</h1>

    <el-row :gutter="20">
      <el-col :span="6">
        <el-card shadow="hover">
          <div class="stat-card">
            <div class="stat-icon">
              <el-icon :size="32" color="#409eff"><ElIconDocument /></el-icon>
            </div>
            <div class="stat-content">
              <div class="stat-value">{{ overview.note_count }}</div>
              <div class="stat-label">笔记总数</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <div class="stat-card">
            <div class="stat-icon">
              <el-icon :size="32" color="#67c23a"><ElIconView /></el-icon>
            </div>
            <div class="stat-content">
              <div class="stat-value">{{ overview.total_pv }}</div>
              <div class="stat-label">总访问量</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <div class="stat-card">
            <div class="stat-icon">
              <el-icon :size="32" color="#e6a23c"><ElIconCollection /></el-icon>
            </div>
            <div class="stat-content">
              <div class="stat-value">{{ overview.tool_count }}</div>
              <div class="stat-label">工具数量</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <div class="stat-card">
            <div class="stat-icon">
              <el-icon :size="32" color="#f56c6c"><ElIconUser /></el-icon>
            </div>
            <div class="stat-content">
              <div class="stat-value">{{ overview.total_uv }}</div>
              <div class="stat-label">访客数量</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-card class="welcome-card" shadow="hover">
      <h2>欢迎使用个人网站管理系统</h2>
      <p>Phase 2 已上线，包含以下功能：</p>
      <ul>
        <li>工具箱管理</li>
        <li>数据统计分析</li>
        <li>全局搜索优化</li>
      </ul>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElRow, ElCol, ElCard, ElIcon } from 'element-plus'
import { Document as ElIconDocument, View as ElIconView, Collection as ElIconCollection, User as ElIconUser } from '@element-plus/icons-vue'
import { analyticsApi } from '@/api/analytics'
import type { AnalyticsOverview } from '@/api/types'

const overview = ref<AnalyticsOverview>({ total_pv: 0, total_uv: 0, today_pv: 0, today_uv: 0, note_count: 0, tool_count: 0 })

onMounted(async () => {
  try {
    const res = await analyticsApi.getOverview()
    if (res.success && res.data) overview.value = res.data
  } catch { /* ignore */ }
})
</script>

<style scoped>
.dashboard-page { padding: 20px; }
.page-title { font-size: 24px; font-weight: 600; margin-bottom: 24px; }
.stat-card { display: flex; align-items: center; gap: 16px; }
.stat-icon { flex-shrink: 0; }
.stat-value { font-size: 28px; font-weight: 600; color: #333; margin-bottom: 4px; }
.stat-label { font-size: 14px; color: #666; }
.welcome-card { margin-top: 20px; }
.welcome-card h2 { font-size: 20px; margin-bottom: 16px; }
.welcome-card p { margin-bottom: 12px; line-height: 1.6; }
.welcome-card ul { margin-left: 20px; line-height: 1.8; }
</style>
