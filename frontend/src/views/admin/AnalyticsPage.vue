<template>
  <div class="analytics-page">
    <h1 class="page-title">数据统计</h1>

    <el-card>
      <p>数据统计功能将在后续版本中完善。</p>
      <p>计划包含：</p>
      <ul>
        <li>访问量趋势图表</li>
        <li>热门文章排行</li>
        <li>访客地域分布</li>
        <li>设备类型统计</li>
      </ul>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ElCard } from 'element-plus'
</script>

<style scoped>
.analytics-page {
  padding: 20px;
}

.page-title {
  font-size: 24px;
  font-weight: 600;
  margin-bottom: 24px;
}

ul {
  margin-left: 20px;
  line-height: 1.8;
}
</style>
