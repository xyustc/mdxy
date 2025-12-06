<template>
  <div class="analytics">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>访问记录</span>
          <div>
            <el-button @click="loadLogs">
              <el-icon><Refresh /></el-icon>
              刷新
            </el-button>
          </div>
        </div>
      </template>
      
      <!-- 筛选条件 -->
      <el-form :inline="true" :model="filters" class="filter-form">
        <el-form-item label="IP地址">
          <el-input
            v-model="filters.ip"
            placeholder="输入IP地址"
            clearable
            style="width: 200px"
          />
        </el-form-item>
        
        <el-form-item label="路径">
          <el-input
            v-model="filters.path"
            placeholder="输入访问路径"
            clearable
            style="width: 200px"
          />
        </el-form-item>
        
        <el-form-item label="日期范围">
          <el-date-picker
            v-model="dateRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            format="YYYY-MM-DD"
            value-format="YYYY-MM-DD"
            style="width: 300px"
          />
        </el-form-item>
        
        <el-form-item>
          <el-button type="primary" @click="handleSearch">
            <el-icon><Search /></el-icon>
            查询
          </el-button>
          <el-button @click="handleReset">
            <el-icon><RefreshLeft /></el-icon>
            重置
          </el-button>
        </el-form-item>
      </el-form>
      
      <!-- 数据表格 -->
      <el-table
        v-loading="loading"
        :data="logs"
        style="width: 100%"
      >
        <el-table-column prop="ip_address" label="IP地址" width="150" />
        <el-table-column prop="path" label="访问路径" min-width="200" />
        <el-table-column prop="method" label="方法" width="80" />
        <el-table-column prop="status_code" label="状态码" width="90">
          <template #default="{ row }">
            <el-tag
              :type="row.status_code < 400 ? 'success' : 'danger'"
              size="small"
            >
              {{ row.status_code }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="response_time" label="响应时间" width="120">
          <template #default="{ row }">
            {{ row.response_time }}ms
          </template>
        </el-table-column>
        <el-table-column prop="device_type" label="设备" width="100" />
        <el-table-column prop="os" label="操作系统" width="150" />
        <el-table-column prop="browser" label="浏览器" width="150" />
        <el-table-column prop="created_at" label="访问时间" width="180">
          <template #default="{ row }">
            {{ formatDateTime(row.created_at) }}
          </template>
        </el-table-column>
      </el-table>
      
      <!-- 分页 -->
      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.limit"
        :total="pagination.total"
        :page-sizes="[20, 50, 100]"
        layout="total, sizes, prev, pager, next, jumper"
        style="margin-top: 20px; justify-content: center"
        @size-change="loadLogs"
        @current-change="loadLogs"
      />
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getAnalyticsLogs } from '../../api/analytics'
import { ElMessage } from 'element-plus'
import { Refresh, Search, RefreshLeft } from '@element-plus/icons-vue'

const loading = ref(false)
const logs = ref([])
const dateRange = ref([])
const filters = ref({
  ip: '',
  path: ''
})
const pagination = ref({
  page: 1,
  limit: 50,
  total: 0
})

const loadLogs = async () => {
  loading.value = true
  
  try {
    const params = {
      page: pagination.value.page,
      limit: pagination.value.limit
    }
    
    if (filters.value.ip) {
      params.ip = filters.value.ip
    }
    
    if (filters.value.path) {
      params.path = filters.value.path
    }
    
    if (dateRange.value && dateRange.value.length === 2) {
      params.start_date = dateRange.value[0]
      params.end_date = dateRange.value[1]
    }
    
    const response = await getAnalyticsLogs(params)
    
    if (response.success) {
      logs.value = response.data
      pagination.value.total = response.total
    }
  } catch (error) {
    ElMessage.error('加载访问记录失败')
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.value.page = 1
  loadLogs()
}

const handleReset = () => {
  filters.value = {
    ip: '',
    path: ''
  }
  dateRange.value = []
  pagination.value.page = 1
  loadLogs()
}

const formatDateTime = (dateStr) => {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit'
  })
}

onMounted(() => {
  loadLogs()
})
</script>

<style scoped>
.analytics {
  width: 100%;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.filter-form {
  margin-bottom: 20px;
}
</style>
