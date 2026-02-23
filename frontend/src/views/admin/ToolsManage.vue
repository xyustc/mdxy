<template>
  <div class="tools-manage">
    <div class="page-header">
      <h1 class="page-title">工具管理</h1>
      <el-button type="primary" @click="openDialog()">添加工具</el-button>
    </div>

    <el-table :data="tools" stripe style="width: 100%">
      <el-table-column prop="name" label="名称" width="160" />
      <el-table-column prop="type" label="类型" width="100">
        <template #default="{ row }">
          <el-tag :type="typeTagMap[row.type] || 'info'" size="small">{{ row.type }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="category" label="分类" width="120" />
      <el-table-column prop="sort_order" label="排序" width="80" />
      <el-table-column label="可见" width="80">
        <template #default="{ row }">
          <el-switch v-model="row.is_visible" @change="handleToggleVisible(row)" />
        </template>
      </el-table-column>
      <el-table-column prop="url" label="链接" show-overflow-tooltip />
      <el-table-column label="操作" width="160" fixed="right">
        <template #default="{ row }">
          <el-button size="small" @click="openDialog(row)">编辑</el-button>
          <el-button size="small" type="danger" @click="handleDelete(row.id)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑工具' : '添加工具'" width="500px">
      <el-form :model="form" label-width="80px">
        <el-form-item label="名称">
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="form.description" type="textarea" :rows="3" />
        </el-form-item>
        <el-form-item label="类型">
          <el-select v-model="form.type" style="width: 100%">
            <el-option label="软件" value="software" />
            <el-option label="视频" value="video" />
            <el-option label="游戏" value="game" />
            <el-option label="链接" value="link" />
          </el-select>
        </el-form-item>
        <el-form-item label="链接">
          <el-input v-model="form.url" />
        </el-form-item>
        <el-form-item label="图标">
          <el-input v-model="form.icon" placeholder="emoji 或图标URL" />
        </el-form-item>
        <el-form-item label="分类">
          <el-input v-model="form.category" />
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="form.sort_order" :min="0" />
        </el-form-item>
        <el-form-item label="可见">
          <el-switch v-model="form.is_visible" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { toolApi } from '@/api/tool'
import type { Tool, ToolForm } from '@/api/types'

const tools = ref<Tool[]>([])
const dialogVisible = ref(false)
const isEdit = ref(false)
const editId = ref(0)

const typeTagMap: Record<string, string> = {
  software: 'primary',
  video: 'danger',
  game: 'warning',
  link: 'success'
}

const defaultForm = (): ToolForm => ({
  name: '', description: '', type: 'link', url: '', icon: '', category: '', sort_order: 0, is_visible: true
})
const form = ref<ToolForm>(defaultForm())

async function loadTools() {
  const res = await toolApi.list()
  if (res.success) tools.value = res.data || []
}

function openDialog(tool?: Tool) {
  if (tool) {
    isEdit.value = true
    editId.value = tool.id
    form.value = { name: tool.name, description: tool.description, type: tool.type, url: tool.url, icon: tool.icon, category: tool.category, sort_order: tool.sort_order, is_visible: tool.is_visible }
  } else {
    isEdit.value = false
    form.value = defaultForm()
  }
  dialogVisible.value = true
}

async function handleSubmit() {
  const res = isEdit.value
    ? await toolApi.update(editId.value, form.value)
    : await toolApi.create(form.value)
  if (res.success) {
    ElMessage.success(isEdit.value ? '更新成功' : '创建成功')
    dialogVisible.value = false
    loadTools()
  } else {
    ElMessage.error(res.error || '操作失败')
  }
}

async function handleDelete(id: number) {
  await ElMessageBox.confirm('确定删除该工具？', '提示', { type: 'warning' })
  const res = await toolApi.delete(id)
  if (res.success) {
    ElMessage.success('删除成功')
    loadTools()
  }
}

async function handleToggleVisible(row: Tool) {
  await toolApi.update(row.id, {
    name: row.name, description: row.description, type: row.type,
    url: row.url, icon: row.icon, category: row.category,
    sort_order: row.sort_order, is_visible: row.is_visible
  })
}

onMounted(loadTools)
</script>

<style scoped>
.tools-manage { padding: 20px; }
.page-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px; }
.page-title { font-size: 24px; font-weight: 600; }
</style>
