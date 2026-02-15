<template>
  <div class="dashboard-page">
    <h1 class="page-title">仪表盘</h1>

    <el-row :gutter="20">
      <el-col :span="6">
        <el-card shadow="hover">
          <div class="stat-card">
            <div class="stat-icon">
              <el-icon :size="32" color="#409eff"><el-icon-document /></el-icon>
            </div>
            <div class="stat-content">
              <div class="stat-value">0</div>
              <div class="stat-label">文章总数</div>
            </div>
          </div>
        </el-card>
      </el-col>

      <el-col :span="6">
        <el-card shadow="hover">
          <div class="stat-card">
            <div class="stat-icon">
              <el-icon :size="32" color="#67c23a"><el-icon-view /></el-icon>
            </div>
            <div class="stat-content">
              <div class="stat-value">0</div>
              <div class="stat-label">总访问量</div>
            </div>
          </div>
        </el-card>
      </el-col>

      <el-col :span="6">
        <el-card shadow="hover">
          <div class="stat-card">
            <div class="stat-icon">
              <el-icon :size="32" color="#e6a23c"><el-icon-collection /></el-icon>
            </div>
            <div class="stat-content">
              <div class="stat-value">0</div>
              <div class="stat-label">工具数量</div>
            </div>
          </div>
        </el-card>
      </el-col>

      <el-col :span="6">
        <el-card shadow="hover">
          <div class="stat-card">
            <div class="stat-icon">
              <el-icon :size="32" color="#f56c6c"><el-icon-user /></el-icon>
            </div>
            <div class="stat-content">
              <div class="stat-value">0</div>
              <div class="stat-label">访客数量</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-card class="welcome-card" shadow="hover">
      <h2>欢迎使用个人网站管理系统</h2>
      <p>这是 Phase 1 的基础版本，包含以下功能：</p>
      <ul>
        <li>个人信息管理</li>
        <li>八股笔记浏览</li>
        <li>基础数据统计</li>
      </ul>
      <p>更多功能将在后续版本中添加。</p>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ElRow, ElCol, ElCard, ElIcon } from 'element-plus'
import { Document as ElIconDocument, View as ElIconView, Collection as ElIconCollection, User as ElIconUser } from '@element-plus/icons-vue'
</script>

<style scoped>
.dashboard-page {
  padding: 20px;
}

.page-title {
  font-size: 24px;
  font-weight: 600;
  margin-bottom: 24px;
}

.stat-card {
  display: flex;
  align-items: center;
  gap: 16px;
}

.stat-icon {
  flex-shrink: 0;
}

.stat-content {
  flex: 1;
}

.stat-value {
  font-size: 28px;
  font-weight: 600;
  color: #333;
  margin-bottom: 4px;
}

.stat-label {
  font-size: 14px;
  color: #666;
}

.welcome-card {
  margin-top: 20px;
}

.welcome-card h2 {
  font-size: 20px;
  margin-bottom: 16px;
}

.welcome-card p {
  margin-bottom: 12px;
  line-height: 1.6;
}

.welcome-card ul {
  margin-left: 20px;
  line-height: 1.8;
}
</style>
