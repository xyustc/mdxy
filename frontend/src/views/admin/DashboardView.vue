<template>
  <div class="dashboard">
    <el-row :gutter="20" style="margin-bottom: 20px">
      <el-col :span="6">
        <el-card shadow="hover">
          <div class="stat-card">
            <div class="stat-icon" style="background: #409eff">
              <el-icon size="30"><View /></el-icon>
            </div>
            <div class="stat-content">
              <div class="stat-value">{{ stats.total_visits }}</div>
              <div class="stat-label">总访问量</div>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :span="6">
        <el-card shadow="hover">
          <div class="stat-card">
            <div class="stat-icon" style="background: #67c23a">
              <el-icon size="30"><User /></el-icon>
            </div>
            <div class="stat-content">
              <div class="stat-value">{{ stats.unique_visitors }}</div>
              <div class="stat-label">独立访客</div>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :span="6">
        <el-card shadow="hover">
          <div class="stat-card">
            <div class="stat-icon" style="background: #e6a23c">
              <el-icon size="30"><Timer /></el-icon>
            </div>
            <div class="stat-content">
              <div class="stat-value">{{ stats.avg_response_time }}ms</div>
              <div class="stat-label">平均响应时间</div>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :span="6">
        <el-card shadow="hover">
          <div class="stat-card">
            <div class="stat-icon" style="background: #f56c6c">
              <el-icon size="30"><Document /></el-icon>
            </div>
            <div class="stat-content">
              <div class="stat-value">{{ stats.top_pages?.length || 0 }}</div>
              <div class="stat-label">访问页面数</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>
    
    <el-row :gutter="20">
      <el-col :span="16">
        <el-card>
          <template #header>
            <span>访问趋势</span>
          </template>
          <div ref="trendChartRef" style="height: 300px"></div>
        </el-card>
      </el-col>
      
      <el-col :span="8">
        <el-card>
          <template #header>
            <span>设备分布</span>
          </template>
          <div ref="deviceChartRef" style="height: 300px"></div>
        </el-card>
      </el-col>
    </el-row>
    
    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="12">
        <el-card>
          <template #header>
            <span>Top 10 访问页面</span>
          </template>
          <el-table :data="stats.top_pages" max-height="300">
            <el-table-column prop="path" label="路径" />
            <el-table-column prop="count" label="访问次数" width="120" />
          </el-table>
        </el-card>
      </el-col>
      
      <el-col :span="12">
        <el-card>
          <template #header>
            <span>浏览器统计</span>
          </template>
          <el-table :data="stats.browser_stats" max-height="300">
            <el-table-column prop="browser" label="浏览器" />
            <el-table-column prop="count" label="访问次数" width="120" />
          </el-table>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, nextTick } from 'vue'
import { getAnalyticsOverview } from '../../api/analytics'
import { ElMessage } from 'element-plus'
import * as echarts from 'echarts'
import { View, User, Timer, Document } from '@element-plus/icons-vue'

const stats = ref({
  total_visits: 0,
  unique_visitors: 0,
  avg_response_time: 0,
  top_pages: [],
  visitor_trends: [],
  device_stats: {},
  os_stats: [],
  browser_stats: []
})

const trendChartRef = ref(null)
const deviceChartRef = ref(null)

const loadStats = async () => {
  try {
    const response = await getAnalyticsOverview()
    if (response.success) {
      stats.value = response.data
      await nextTick()
      initCharts()
    }
  } catch (error) {
    ElMessage.error('加载统计数据失败')
  }
}

const initCharts = () => {
  // 访问趋势图
  if (trendChartRef.value) {
    const trendChart = echarts.init(trendChartRef.value)
    const dates = stats.value.visitor_trends.map(item => item.date)
    const counts = stats.value.visitor_trends.map(item => item.count)
    
    trendChart.setOption({
      tooltip: {
        trigger: 'axis'
      },
      xAxis: {
        type: 'category',
        data: dates
      },
      yAxis: {
        type: 'value'
      },
      series: [{
        data: counts,
        type: 'line',
        smooth: true,
        areaStyle: {
          color: 'rgba(64, 158, 255, 0.2)'
        },
        itemStyle: {
          color: '#409eff'
        }
      }]
    })
  }
  
  // 设备分布饼图
  if (deviceChartRef.value) {
    const deviceChart = echarts.init(deviceChartRef.value)
    const data = Object.entries(stats.value.device_stats).map(([name, value]) => ({
      name,
      value
    }))
    
    deviceChart.setOption({
      tooltip: {
        trigger: 'item'
      },
      legend: {
        bottom: 'bottom'
      },
      series: [{
        type: 'pie',
        radius: ['40%', '70%'],
        avoidLabelOverlap: false,
        itemStyle: {
          borderRadius: 10,
          borderColor: '#fff',
          borderWidth: 2
        },
        label: {
          show: false,
          position: 'center'
        },
        emphasis: {
          label: {
            show: true,
            fontSize: '18',
            fontWeight: 'bold'
          }
        },
        labelLine: {
          show: false
        },
        data
      }]
    })
  }
}

onMounted(() => {
  loadStats()
  
  // 定时刷新数据
  const timer = setInterval(() => {
    loadStats()
  }, 30000) // 每30秒刷新一次
  
  // 组件卸载时清除定时器
  onUnmounted(() => {
    clearInterval(timer)
  })
})
</script>

<style scoped>
.stat-card {
  display: flex;
  align-items: center;
}

.stat-icon {
  width: 60px;
  height: 60px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 20px;
}

.stat-content {
  flex: 1;
}

.stat-value {
  font-size: 24px;
  font-weight: bold;
  color: #303133;
  margin-bottom: 5px;
}

.stat-label {
  font-size: 14px;
  color: #909399;
}
</style>