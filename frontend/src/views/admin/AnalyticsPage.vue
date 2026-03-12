<template>
  <div class="analytics-page">
    <h1 class="page-title">数据统计</h1>

    <!-- 概览卡片 -->
    <el-row :gutter="20" class="overview-row">
      <el-col :span="6">
        <el-card shadow="hover">
          <div class="stat-card">
            <div class="stat-icon"><el-icon :size="32" color="#409eff"><ElIconView /></el-icon></div>
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
            <div class="stat-icon"><el-icon :size="32" color="#67c23a"><ElIconUser /></el-icon></div>
            <div class="stat-content">
              <div class="stat-value">{{ overview.total_uv }}</div>
              <div class="stat-label">总访客数</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <div class="stat-card">
            <div class="stat-icon"><el-icon :size="32" color="#e6a23c"><ElIconView /></el-icon></div>
            <div class="stat-content">
              <div class="stat-value">{{ overview.today_pv }}</div>
              <div class="stat-label">今日PV</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <div class="stat-card">
            <div class="stat-icon"><el-icon :size="32" color="#f56c6c"><ElIconUser /></el-icon></div>
            <div class="stat-content">
              <div class="stat-value">{{ overview.today_uv }}</div>
              <div class="stat-label">今日UV</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- PV/UV 趋势 -->
    <el-card class="chart-card" shadow="hover">
      <template #header>
        <div class="chart-header">
          <span>访问趋势</span>
          <el-radio-group v-model="trendDays" size="small" @change="loadTrends">
            <el-radio-button :value="7">7天</el-radio-button>
            <el-radio-button :value="30">30天</el-radio-button>
            <el-radio-button :value="90">90天</el-radio-button>
          </el-radio-group>
        </div>
      </template>
      <v-chart :option="trendOption" style="height: 350px" autoresize />
    </el-card>

    <!-- 第二行：热门页面 + 设备统计 -->
    <el-row :gutter="20" class="chart-row">
      <el-col :span="14">
        <el-card shadow="hover">
          <template #header><span>热门页面 TOP 10</span></template>
          <el-table :data="popularPages" size="small" stripe>
            <el-table-column type="index" width="50" />
            <el-table-column prop="path" label="路径" show-overflow-tooltip />
            <el-table-column prop="count" label="访问量" width="100" sortable />
          </el-table>
        </el-card>
      </el-col>
      <el-col :span="10">
        <el-card shadow="hover">
          <template #header><span>设备类型分布</span></template>
          <v-chart :option="deviceOption" style="height: 300px" autoresize />
        </el-card>
      </el-col>
    </el-row>

    <!-- 第三行：浏览器 + 地域 -->
    <el-row :gutter="20" class="chart-row">
      <el-col :span="10">
        <el-card shadow="hover">
          <template #header><span>浏览器分布</span></template>
          <v-chart :option="browserOption" style="height: 300px" autoresize />
        </el-card>
      </el-col>
      <el-col :span="14">
        <el-card shadow="hover">
          <template #header><span>地域分布</span></template>
          <el-table :data="geoStats" size="small" stripe>
            <el-table-column type="index" width="50" />
            <el-table-column prop="country" label="国家" width="120" />
            <el-table-column prop="region" label="地区" />
            <el-table-column prop="count" label="访问量" width="100" sortable />
          </el-table>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>


<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import { LineChart, PieChart } from 'echarts/charts'
import { GridComponent, TooltipComponent, LegendComponent, TitleComponent } from 'echarts/components'
import VChart from 'vue-echarts'
import { View as ElIconView, User as ElIconUser } from '@element-plus/icons-vue'
import { analyticsApi } from '@/api/analytics'
import type { AnalyticsOverview, DailyStat, PageStat, GeoStat } from '@/api/types'

use([CanvasRenderer, LineChart, PieChart, GridComponent, TooltipComponent, LegendComponent, TitleComponent])

const overview = ref<AnalyticsOverview>({ total_pv: 0, total_uv: 0, today_pv: 0, today_uv: 0, note_count: 0, tool_count: 0 })
const trendDays = ref(30)
const pvData = ref<DailyStat[]>([])
const uvData = ref<DailyStat[]>([])
const popularPages = ref<PageStat[]>([])
const deviceData = ref<{ device_type: string; count: number }[]>([])
const browserData = ref<{ browser: string; count: number }[]>([])
const geoStats = ref<GeoStat[]>([])

const trendOption = computed(() => ({
  tooltip: { trigger: 'axis' },
  legend: { data: ['PV', 'UV'] },
  grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
  xAxis: { type: 'category', data: pvData.value.map(d => d.date) },
  yAxis: { type: 'value' },
  series: [
    { name: 'PV', type: 'line', smooth: true, data: pvData.value.map(d => d.count), itemStyle: { color: '#409eff' } },
    { name: 'UV', type: 'line', smooth: true, data: uvData.value.map(d => d.count), itemStyle: { color: '#67c23a' } }
  ]
}))

const deviceOption = computed(() => ({
  tooltip: { trigger: 'item', formatter: '{b}: {c} ({d}%)' },
  series: [{
    type: 'pie', radius: ['40%', '70%'],
    data: deviceData.value.map(d => ({ name: d.device_type, value: d.count })),
    emphasis: { itemStyle: { shadowBlur: 10, shadowOffsetX: 0, shadowColor: 'rgba(0,0,0,0.5)' } }
  }]
}))

const browserOption = computed(() => ({
  tooltip: { trigger: 'item', formatter: '{b}: {c} ({d}%)' },
  series: [{
    type: 'pie', radius: '65%',
    data: browserData.value.map(d => ({ name: d.browser, value: d.count }))
  }]
}))

async function loadTrends() {
  const res = await analyticsApi.getTrends(trendDays.value)
  if (res.success && res.data) {
    pvData.value = res.data.pv || []
    uvData.value = res.data.uv || []
  }
}

async function loadAll() {
  const [overviewRes, popularRes, deviceRes, browserRes, geoRes] = await Promise.all([
    analyticsApi.getOverview(),
    analyticsApi.getPopularPages(),
    analyticsApi.getDevices(),
    analyticsApi.getBrowsers(),
    analyticsApi.getGeo()
  ])
  if (overviewRes.success && overviewRes.data) overview.value = overviewRes.data
  if (popularRes.success) popularPages.value = popularRes.data || []
  if (deviceRes.success) deviceData.value = deviceRes.data || []
  if (browserRes.success) browserData.value = browserRes.data || []
  if (geoRes.success) geoStats.value = geoRes.data || []
}

onMounted(() => { loadAll(); loadTrends() })
</script>

<style scoped>
.analytics-page { padding: 20px; }
.page-title { font-size: 24px; font-weight: 600; margin-bottom: 24px; }
.overview-row { margin-bottom: 20px; }
.stat-card { display: flex; align-items: center; gap: 16px; }
.stat-value { font-size: 28px; font-weight: 600; color: #333; margin-bottom: 4px; }
.stat-label { font-size: 14px; color: #666; }
.chart-card { margin-bottom: 20px; }
.chart-header { display: flex; justify-content: space-between; align-items: center; }
.chart-row { margin-bottom: 20px; }
</style>
