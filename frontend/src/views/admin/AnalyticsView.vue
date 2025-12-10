<template>
  <div class="analytics">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>访问记录</span>
          <div>
            <el-button :icon="Refresh" @click="loadLogs" :loading="loading">
              刷新
            </el-button>
          </div>
        </div>
      </template>
      
      <!-- 筛选表单 -->
      <el-form :inline="true" :model="filters" class="filter-form" label-width="80px">
        <el-form-item label="时间范围">
          <el-date-picker
            v-model="dateRange"
            type="datetimerange"
            range-separator="至"
            start-placeholder="开始时间"
            end-placeholder="结束时间"
            format="YYYY-MM-DD HH:mm:ss"
            value-format="YYYY-MM-DDTHH:mm:ss"
          />
        </el-form-item>
        
        <el-form-item label="IP地址">
          <el-input v-model="filters.ip" placeholder="请输入IP地址" />
        </el-form-item>
        
        <el-form-item label="访问路径">
          <el-input v-model="filters.path" placeholder="请输入路径" />
        </el-form-item>
        
        <el-form-item>
          <el-button type="primary" :icon="Search" @click="handleSearch">
            搜索
          </el-button>
          <el-button :icon="RefreshLeft" @click="handleReset">
            重置
          </el-button>
        </el-form-item>
      </el-form>
      
      <!-- 访问记录表格 -->
      <el-table :data="logs" border stripe v-loading="loading" style="width: 100%">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="ip_address" label="IP地址" width="150" />
        <el-table-column prop="visitor_id" label="用户标识" width="200">
          <template #default="{ row }">
            <el-tooltip 
              v-if="row.visitor_id" 
              :content="row.visitor_id" 
              placement="top"
            >
              <span>{{ row.visitor_id.substring(0, 8) }}...</span>
            </el-tooltip>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column prop="path" label="访问路径" show-overflow-tooltip>
          <template #default="{ row }">
            <el-link 
              :href="'#/note' + row.path" 
              target="_blank" 
              type="primary"
              v-if="row.path.startsWith('/')"
            >
              {{ row.path }}
            </el-link>
            <span v-else>{{ row.path }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="method" label="方法" width="80" />
        <el-table-column prop="status_code" label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="getStatusCodeType(row.status_code)">
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
  if (!dateStr) return '-';
  // 解析UTC时间字符串
  const date = new Date(dateStr);
  // 转换为本地时间并格式化
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
    timeZoneName: 'short'
  });
}

const getStatusCodeType = (statusCode) => {
  if (statusCode >= 200 && statusCode < 300) {
    return 'success'
  } else if (statusCode >= 300 && statusCode < 400) {
    return 'warning'
  } else if (statusCode >= 400) {
    return 'danger'
  }
  return ''
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