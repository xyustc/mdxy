<template>
  <div class="analytics-page">
    <section class="metrics-grid">
      <article v-for="metric in overviewMetrics" :key="metric.label" class="surface-panel metric-card">
        <component :is="metric.icon" class="metric-card__icon" />
        <div>
          <strong>{{ metric.value }}</strong>
          <span>{{ metric.label }}</span>
        </div>
      </article>
    </section>

    <el-card class="analytics-card" shadow="never">
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
      <v-chart :option="trendOption" class="chart-view chart-view--large" autoresize />
    </el-card>

    <section class="analytics-grid">
      <el-card class="analytics-card" shadow="never">
        <template #header><span>热门页面 TOP 10</span></template>
        <el-table :data="popularPages" size="small" stripe>
          <el-table-column type="index" width="50" />
          <el-table-column prop="path" label="路径" show-overflow-tooltip />
          <el-table-column prop="count" label="访问量" width="100" sortable />
        </el-table>
      </el-card>

      <el-card class="analytics-card" shadow="never">
        <template #header><span>设备类型分布</span></template>
        <v-chart :option="deviceOption" class="chart-view" autoresize />
      </el-card>
    </section>

    <section class="analytics-grid">
      <el-card class="analytics-card" shadow="never">
        <template #header><span>浏览器分布</span></template>
        <v-chart :option="browserOption" class="chart-view" autoresize />
      </el-card>

      <el-card class="analytics-card" shadow="never">
        <template #header><span>地域分布</span></template>
        <el-table :data="geoStats" size="small" stripe>
          <el-table-column type="index" width="50" />
          <el-table-column prop="country" label="国家" width="120" />
          <el-table-column prop="region" label="地区" />
          <el-table-column prop="count" label="访问量" width="100" sortable />
        </el-table>
      </el-card>
    </section>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import { LineChart, PieChart } from 'echarts/charts'
import { GridComponent, TooltipComponent, LegendComponent, TitleComponent } from 'echarts/components'
import VChart from 'vue-echarts'
import { ElCard, ElRadioGroup, ElRadioButton, ElTable, ElTableColumn } from 'element-plus'
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

const overviewMetrics = computed(() => [
  { label: '总访问量', value: overview.value.total_pv, icon: ElIconView },
  { label: '总访客数', value: overview.value.total_uv, icon: ElIconUser },
  { label: '今日 PV', value: overview.value.today_pv, icon: ElIconView },
  { label: '今日 UV', value: overview.value.today_uv, icon: ElIconUser }
])

const trendOption = computed(() => ({
  tooltip: { trigger: 'axis' },
  legend: { data: ['PV', 'UV'] },
  grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
  xAxis: { type: 'category', data: pvData.value.map((d) => d.date) },
  yAxis: { type: 'value' },
  series: [
    { name: 'PV', type: 'line', smooth: true, data: pvData.value.map((d) => d.count), itemStyle: { color: '#355f54' } },
    { name: 'UV', type: 'line', smooth: true, data: uvData.value.map((d) => d.count), itemStyle: { color: '#8b5e3c' } }
  ]
}))

const deviceOption = computed(() => ({
  tooltip: { trigger: 'item', formatter: '{b}: {c} ({d}%)' },
  series: [
    {
      type: 'pie',
      radius: ['40%', '70%'],
      data: deviceData.value.map((d) => ({ name: d.device_type, value: d.count })),
      emphasis: { itemStyle: { shadowBlur: 10, shadowOffsetX: 0, shadowColor: 'rgba(0,0,0,0.35)' } }
    }
  ]
}))

const browserOption = computed(() => ({
  tooltip: { trigger: 'item', formatter: '{b}: {c} ({d}%)' },
  series: [{ type: 'pie', radius: '65%', data: browserData.value.map((d) => ({ name: d.browser, value: d.count })) }]
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

onMounted(() => {
  loadAll()
  loadTrends()
})
</script>

<style scoped>
.analytics-page {
  display: grid;
  gap: var(--space-xl);
}

.metrics-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: var(--space-md);
}

.metric-card {
  padding: var(--space-lg);
  display: flex;
  align-items: center;
  gap: var(--space-md);
}

.metric-card__icon {
  width: 1.35rem;
  height: 1.35rem;
  color: var(--accent-primary);
}

.metric-card strong {
  display: block;
  font-family: var(--font-display);
  font-size: 1.85rem;
  letter-spacing: -0.04em;
}

.metric-card span {
  color: var(--text-secondary);
  font-size: 0.9rem;
}

.analytics-card {
  border: 1px solid var(--border-primary);
  border-radius: var(--radius-xl);
  background: var(--bg-panel);
}

.analytics-card :deep(.el-card__header) {
  border-bottom-color: var(--border-primary);
}

.chart-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.chart-view {
  height: 300px;
}

.chart-view--large {
  height: 360px;
}

.analytics-grid {
  display: grid;
  grid-template-columns: 1.2fr 1fr;
  gap: var(--space-lg);
}

@media (max-width: 1100px) {
  .metrics-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .analytics-grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 640px) {
  .metrics-grid {
    grid-template-columns: 1fr;
  }

  .chart-header {
    flex-direction: column;
    align-items: flex-start;
    gap: var(--space-sm);
  }
}
</style>
