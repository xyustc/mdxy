<template>
  <div class="tools-manage">
    <div class="tools-manage__toolbar">
      <el-button type="primary" @click="openDialog()">添加工具</el-button>
    </div>

    <section class="surface-panel tools-summary">
      <article>
        <span class="meta-label">Total</span>
        <strong>{{ tools.length }}</strong>
      </article>
      <article>
        <span class="meta-label">Visible</span>
        <strong>{{ visibleCount }}</strong>
      </article>
      <article>
        <span class="meta-label">Categories</span>
        <strong>{{ categoryCount }}</strong>
      </article>
    </section>

    <el-card class="manage-card" shadow="never">
      <el-table :data="tools" stripe style="width: 100%">
        <el-table-column prop="name" label="名称" width="180" />
        <el-table-column prop="type" label="类型" width="100">
          <template #default="{ row }">
            <el-tag :type="typeTagMap[normalizeToolType(row.type)] || 'info'" size="small">
              {{ normalizeToolType(row.type) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="category" label="分类" width="140" />
        <el-table-column prop="sort_order" label="排序" width="80" />
        <el-table-column label="可见" width="90">
          <template #default="{ row }">
            <el-switch v-model="row.is_visible" @change="handleToggleVisible(row)" />
          </template>
        </el-table-column>
        <el-table-column prop="url" label="链接" show-overflow-tooltip />
        <el-table-column label="操作" width="170" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="openDialog(row)">编辑</el-button>
            <el-button size="small" type="danger" @click="handleDelete(row.id)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑工具' : '添加工具'" width="520px">
      <el-form :model="form" label-width="90px">
        <el-form-item label="名称"><el-input v-model="form.name" /></el-form-item>
        <el-form-item label="描述"><el-input v-model="form.description" type="textarea" :rows="3" /></el-form-item>
        <el-form-item label="类型">
          <el-select v-model="form.type" style="width: 100%">
            <el-option label="应用" value="app" />
            <el-option label="视频" value="video" />
            <el-option label="游戏" value="game" />
            <el-option label="链接" value="link" />
          </el-select>
        </el-form-item>
        <el-form-item label="链接"><el-input v-model="form.url" placeholder="外链或站内路由，例如 /tools/etc-image-obfuscator" /></el-form-item>
        <el-form-item label="图标"><el-input v-model="form.icon" placeholder="emoji 或图标URL" /></el-form-item>
        <el-form-item label="分类"><el-input v-model="form.category" /></el-form-item>
        <el-form-item label="排序"><el-input-number v-model="form.sort_order" :min="0" /></el-form-item>
        <el-form-item label="可见"><el-switch v-model="form.is_visible" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import {
  ElButton,
  ElCard,
  ElDialog,
  ElForm,
  ElFormItem,
  ElInput,
  ElInputNumber,
  ElMessage,
  ElMessageBox,
  ElOption,
  ElSelect,
  ElSwitch,
  ElTable,
  ElTableColumn,
  ElTag
} from 'element-plus'
import { toolApi } from '@/api/tool'
import type { Tool, ToolForm } from '@/api/types'

const tools = ref<Tool[]>([])
const dialogVisible = ref(false)
const isEdit = ref(false)
const editId = ref(0)

const typeTagMap: Record<string, 'primary' | 'danger' | 'warning' | 'success' | 'info'> = {
  app: 'primary',
  video: 'danger',
  game: 'warning',
  link: 'success'
}

const defaultForm = (): ToolForm => ({
  name: '',
  description: '',
  type: 'app',
  url: '',
  icon: '',
  category: '',
  sort_order: 0,
  is_visible: true
})

const form = ref<ToolForm>(defaultForm())

const visibleCount = computed(() => tools.value.filter((item) => item.is_visible).length)
const categoryCount = computed(() => new Set(tools.value.map((item) => item.category).filter(Boolean)).size)

async function loadTools() {
  const res = await toolApi.adminList()
  if (res.success) tools.value = res.data || []
}

function openDialog(tool?: Tool) {
  if (tool) {
    isEdit.value = true
    editId.value = tool.id
    form.value = {
      name: tool.name,
      description: tool.description,
      type: normalizeToolType(tool.type),
      url: tool.url,
      icon: tool.icon,
      category: tool.category,
      sort_order: tool.sort_order,
      is_visible: tool.is_visible
    }
  } else {
    isEdit.value = false
    form.value = defaultForm()
  }
  dialogVisible.value = true
}

async function handleSubmit() {
  const payload: ToolForm = { ...form.value, type: normalizeToolType(form.value.type) }
  const res = isEdit.value ? await toolApi.update(editId.value, payload) : await toolApi.create(payload)
  if (res.success) {
    ElMessage.success(isEdit.value ? '更新成功' : '创建成功')
    dialogVisible.value = false
    loadTools()
  } else {
    ElMessage.error(res.error || '操作失败')
  }
}

async function handleDelete(id: number) {
  try {
    await ElMessageBox.confirm('确定删除该工具？', '提示', { type: 'warning' })
  } catch {
    return
  }

  const res = await toolApi.delete(id)
  if (res.success) {
    ElMessage.success('删除成功')
    loadTools()
  }
}

async function handleToggleVisible(row: Tool) {
  try {
    await toolApi.update(row.id, {
      name: row.name,
      description: row.description,
      type: normalizeToolType(row.type),
      url: row.url,
      icon: row.icon,
      category: row.category,
      sort_order: row.sort_order,
      is_visible: row.is_visible
    })
  } catch {
    ElMessage.error('更新可见性失败')
    loadTools()
  }
}

onMounted(loadTools)

function normalizeToolType(type: string) {
  return type === 'software' ? 'app' : type
}
</script>

<style scoped>
.tools-manage {
  display: grid;
  gap: var(--space-lg);
}

.tools-manage__toolbar {
  display: flex;
  justify-content: flex-end;
}

.tools-summary {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: var(--space-md);
  padding: var(--space-lg);
}

.tools-summary article {
  border: 1px solid var(--border-primary);
  border-radius: var(--radius-lg);
  padding: var(--space-md);
  background: var(--bg-panel-strong);
}

.tools-summary strong {
  display: block;
  margin-top: 0.2rem;
  font-family: var(--font-display);
  font-size: 1.9rem;
  letter-spacing: -0.04em;
}

.manage-card {
  border: 1px solid var(--border-primary);
  border-radius: var(--radius-xl);
  background: var(--bg-panel);
}

@media (max-width: 860px) {
  .tools-manage__toolbar {
    justify-content: stretch;
  }

  .tools-manage__toolbar :deep(.el-button) {
    width: 100%;
  }

  .tools-summary {
    grid-template-columns: 1fr;
  }
}
</style>
