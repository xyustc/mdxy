<template>
  <div class="notes-manage">
    <!-- Left: tree panel -->
    <aside class="notes-tree-panel">
      <div class="notes-tree-toolbar">
        <el-button size="small" @click="showNewNoteDialog">新建笔记</el-button>
        <el-button size="small" @click="showNewDirDialog">新建目录</el-button>
      </div>
      <el-tree
        :data="treeData"
        :props="treeProps"
        node-key="path"
        highlight-current
        :expand-on-click-node="false"
        class="notes-tree"
        @node-click="handleNodeClick"
      >
        <template #default="{ node, data }">
          <span class="tree-node">
            <span class="tree-node__label">
              <span v-if="data.featured" class="tree-node__star" title="Featured">★</span>
              {{ node.label }}
            </span>
            <span v-if="data.type === 'file'" class="tree-node__actions">
              <el-switch
                :model-value="data.featured"
                size="small"
                title="精选首页"
                @click.stop
                @change="(val: boolean) => setFeatured(data, val)"
              />
              <el-button link size="small" type="danger" title="删除" @click.stop="deleteNote(data)">✕</el-button>
            </span>
          </span>
        </template>
      </el-tree>
    </aside>

    <!-- Right: editor panel -->
    <div class="notes-editor-panel">
      <div v-if="!currentPath" class="notes-editor-empty">
        <p>从左侧选择一个笔记开始编辑</p>
      </div>
      <template v-else>
        <div class="notes-editor-toolbar">
          <span class="notes-editor-path">{{ currentPath }}</span>
          <div class="notes-editor-toolbar__actions">
            <el-button type="primary" size="small" :disabled="!isDirty" @click="saveNote">保存</el-button>
          </div>
        </div>
        <div class="notes-editor-body">
          <textarea
            v-model="editorContent"
            class="notes-editor-textarea"
            spellcheck="false"
            @input="isDirty = true"
          />
          <div class="notes-editor-preview prose" v-html="renderedContent" />
        </div>
      </template>
    </div>

    <!-- New note dialog -->
    <el-dialog v-model="newNoteVisible" title="新建笔记" width="420px">
      <el-form @submit.prevent="createNote">
        <el-form-item label="路径（相对于 notes/）">
          <el-input v-model="newNotePath" placeholder="例：folder/my-note.md" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="newNoteVisible = false">取消</el-button>
        <el-button type="primary" @click="createNote">创建</el-button>
      </template>
    </el-dialog>

    <!-- New dir dialog -->
    <el-dialog v-model="newDirVisible" title="新建目录" width="420px">
      <el-form @submit.prevent="createDir">
        <el-form-item label="路径（相对于 notes/）">
          <el-input v-model="newDirPath" placeholder="例：folder/sub" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="newDirVisible = false">取消</el-button>
        <el-button type="primary" @click="createDir">创建</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onBeforeUnmount } from 'vue'
import { onBeforeRouteLeave } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import MarkdownIt from 'markdown-it'
import { markdownHighlight } from '@/utils/highlight'
import 'highlight.js/styles/github.css'
import { noteApi } from '@/api/note'
import { parseFrontmatter } from '@/utils/frontmatter'
import { useSiteStore } from '@/stores/site'
import type { NoteNode } from '@/api/types'

const siteStore = useSiteStore()

const treeData = ref<NoteNode[]>([])
const treeProps = { label: 'name', children: 'children' }

const currentPath = ref('')
const editorContent = ref('')
const isDirty = ref(false)

const newNoteVisible = ref(false)
const newNotePath = ref('')
const newDirVisible = ref(false)
const newDirPath = ref('')

const md = new MarkdownIt({ html: false, linkify: true, typographer: true, highlight: markdownHighlight })

const renderedContent = computed(() => {
  const { body } = parseFrontmatter(editorContent.value)
  return md.render(body)
})

async function loadTree() {
  const res = await noteApi.adminGetTree()
  if (res.success && res.data) treeData.value = res.data
}

loadTree()

async function handleNodeClick(data: NoteNode) {
  if (data.type !== 'file') return
  if (isDirty.value) {
    try {
      await ElMessageBox.confirm('有未保存的更改，是否放弃？', '提示', { type: 'warning' })
    } catch { return }
  }
  const res = await noteApi.adminGetContent(data.path)
  if (!res.success || !res.data) return
  currentPath.value = data.path
  editorContent.value = res.data.content
  isDirty.value = false
}

async function saveNote() {
  if (!currentPath.value) return
  const res = await noteApi.adminSave(currentPath.value, editorContent.value)
  if (res.success) {
    isDirty.value = false
    ElMessage.success('已保存')
  } else {
    ElMessage.error('保存失败')
  }
}

async function setFeatured(data: NoteNode, val: boolean) {
  const res = await noteApi.adminSetFeatured(data.path, val)
  if (res.success) {
    ElMessage.success(val ? '已设为精选' : '已取消精选')
    siteStore.clearHomeCache()
    await loadTree()
  } else {
    ElMessage.error('操作失败')
  }
}

async function deleteNote(data: NoteNode) {
  try {
    await ElMessageBox.confirm(`确认删除 "${data.name}"？此操作不可撤销。`, '删除笔记', { type: 'warning' })
  } catch { return }
  const res = await noteApi.adminDelete(data.path)
  if (res.success) {
    ElMessage.success('已删除')
    if (currentPath.value === data.path) {
      currentPath.value = ''
      editorContent.value = ''
      isDirty.value = false
    }
    await loadTree()
  } else {
    ElMessage.error('删除失败')
  }
}

function showNewNoteDialog() { newNotePath.value = ''; newNoteVisible.value = true }
function showNewDirDialog() { newDirPath.value = ''; newDirVisible.value = true }

async function createNote() {
  let p = newNotePath.value.trim()
  if (!p) return
  if (!p.endsWith('.md')) p += '.md'
  const res = await noteApi.adminSave(p, '')
  if (res.success) {
    ElMessage.success('已创建')
    newNoteVisible.value = false
    await loadTree()
  } else {
    ElMessage.error('创建失败')
  }
}

async function createDir() {
  const p = newDirPath.value.trim()
  if (!p) return
  const res = await noteApi.adminCreateDir(p)
  if (res.success) {
    ElMessage.success('目录已创建')
    newDirVisible.value = false
    await loadTree()
  } else {
    ElMessage.error('创建失败')
  }
}

const beforeUnload = (e: BeforeUnloadEvent) => {
  if (isDirty.value) { e.preventDefault(); e.returnValue = '' }
}
window.addEventListener('beforeunload', beforeUnload)
onBeforeUnmount(() => window.removeEventListener('beforeunload', beforeUnload))

onBeforeRouteLeave(async (_to, _from, next) => {
  if (!isDirty.value) return next()
  try {
    await ElMessageBox.confirm('有未保存的更改，是否放弃？', '提示', { type: 'warning' })
    next()
  } catch { next(false) }
})
</script>

<style scoped>
.notes-manage {
  display: grid;
  grid-template-columns: 260px minmax(0, 1fr);
  height: calc(100vh - 140px);
  border: 1px solid var(--border-primary);
  border-radius: var(--radius-lg);
  overflow: hidden;
  background: var(--bg-panel);
}

/* ── 左侧树面板 ── */
.notes-tree-panel {
  border-right: 1px solid var(--border-primary);
  display: flex;
  flex-direction: column;
  overflow: hidden;
  background: var(--bg-page);
}

.notes-tree-toolbar {
  padding: 0.75rem 1rem;
  border-bottom: 1px solid var(--border-primary);
  display: flex;
  gap: 0.5rem;
  flex-shrink: 0;
}

.notes-tree {
  flex: 1;
  overflow-y: auto;
  padding: 0.5rem 0;
  --el-tree-node-hover-bg-color: var(--bg-panel-muted);
  --el-tree-node-content-height: 34px;
  --el-font-size-base: 0.875rem;
}

/* el-tree 节点字体对齐 Admin 风格 */
.notes-tree :deep(.el-tree-node__content) {
  font-size: 0.875rem;
  color: var(--text-primary);
  border-radius: var(--radius-md);
  margin: 0 0.5rem;
  padding-right: 0.25rem;
}

.notes-tree :deep(.el-tree-node__content:hover) {
  background: var(--bg-panel-muted);
}

.notes-tree :deep(.el-tree-node.is-current > .el-tree-node__content) {
  background: var(--bg-panel-muted);
  color: var(--text-primary);
  font-weight: 600;
}

.notes-tree :deep(.el-tree-node__expand-icon) {
  color: var(--text-muted);
}

.tree-node {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
  min-width: 0;
  gap: 0.25rem;
}

.tree-node__label {
  display: flex;
  align-items: center;
  gap: 0.3rem;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  min-width: 0;
}

.tree-node__star {
  color: #f59e0b;
  font-size: 0.75rem;
  flex-shrink: 0;
}

.tree-node__actions {
  display: none;
  align-items: center;
  gap: 0.35rem;
  flex-shrink: 0;
}

.tree-node:hover .tree-node__actions {
  display: flex;
}

/* ── 右侧编辑区 ── */
.notes-editor-panel {
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.notes-editor-empty {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  color: var(--text-muted);
  font-size: 0.9rem;
}

.notes-editor-toolbar {
  padding: 0.65rem 1.1rem;
  border-bottom: 1px solid var(--border-primary);
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
  flex-shrink: 0;
  background: var(--bg-page);
}

.notes-editor-path {
  font-family: var(--font-mono);
  font-size: 0.78rem;
  color: var(--text-muted);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.notes-editor-toolbar__actions {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  flex-shrink: 0;
}

.notes-editor-body {
  flex: 1;
  display: grid;
  grid-template-columns: 1fr 1fr;
  overflow: hidden;
}

.notes-editor-textarea {
  resize: none;
  border: none;
  border-right: 1px solid var(--border-primary);
  outline: none;
  padding: 1.1rem 1.25rem;
  font-family: var(--font-mono);
  font-size: 0.82rem;
  line-height: 1.65;
  background: var(--bg-page);
  color: var(--text-primary);
  overflow-y: auto;
}

.notes-editor-textarea::placeholder {
  color: var(--text-muted);
}

/* ── 预览 prose 样式 ── */
.notes-editor-preview {
  padding: 1.25rem 1.5rem;
  overflow-y: auto;
  font-size: 0.9rem;
  line-height: 1.75;
  color: var(--text-primary);
}

.notes-editor-preview :deep(h1),
.notes-editor-preview :deep(h2),
.notes-editor-preview :deep(h3),
.notes-editor-preview :deep(h4) {
  font-family: var(--font-display);
  font-weight: 700;
  letter-spacing: -0.02em;
  margin: 1.4em 0 0.5em;
  line-height: 1.3;
  color: var(--text-primary);
}

.notes-editor-preview :deep(h1) { font-size: 1.5rem; }
.notes-editor-preview :deep(h2) { font-size: 1.2rem; border-bottom: 1px solid var(--border-primary); padding-bottom: 0.3em; }
.notes-editor-preview :deep(h3) { font-size: 1rem; }

.notes-editor-preview :deep(p),
.notes-editor-preview :deep(ul),
.notes-editor-preview :deep(ol) {
  margin: 0.75em 0;
}

.notes-editor-preview :deep(ul),
.notes-editor-preview :deep(ol) {
  padding-left: 1.5em;
}

.notes-editor-preview :deep(li) {
  margin: 0.25em 0;
}

.notes-editor-preview :deep(code) {
  font-family: var(--font-mono);
  font-size: 0.82em;
  background: var(--bg-panel-muted);
  border: 1px solid var(--border-primary);
  border-radius: var(--radius-sm);
  padding: 0.1em 0.4em;
}

.notes-editor-preview :deep(pre) {
  margin: 1em 0;
  border-radius: var(--radius-md);
  overflow: hidden;
  border: 1px solid var(--border-primary);
}

.notes-editor-preview :deep(pre code) {
  display: block;
  padding: 1rem 1.1rem;
  font-size: 0.82rem;
  line-height: 1.6;
  background: var(--bg-panel-muted);
  border: none;
  border-radius: 0;
  overflow-x: auto;
}

.notes-editor-preview :deep(pre.hljs) {
  margin: 1em 0;
  border-radius: var(--radius-md);
  border: 1px solid var(--border-primary);
  overflow: hidden;
}

.notes-editor-preview :deep(pre.hljs code) {
  display: block;
  padding: 1rem 1.1rem;
  font-size: 0.82rem;
  line-height: 1.6;
  overflow-x: auto;
  background: transparent;
  border: none;
  border-radius: 0;
}

.notes-editor-preview :deep(blockquote) {
  margin: 1em 0;
  padding: 0.5em 1em;
  border-left: 3px solid var(--border-strong);
  color: var(--text-secondary);
  background: var(--bg-panel-muted);
  border-radius: 0 var(--radius-sm) var(--radius-sm) 0;
}

.notes-editor-preview :deep(a) {
  color: var(--accent-primary);
  text-decoration: underline;
  text-underline-offset: 2px;
}

.notes-editor-preview :deep(hr) {
  border: none;
  border-top: 1px solid var(--border-primary);
  margin: 1.5em 0;
}

.notes-editor-preview :deep(table) {
  width: 100%;
  border-collapse: collapse;
  font-size: 0.85rem;
  margin: 1em 0;
}

.notes-editor-preview :deep(th),
.notes-editor-preview :deep(td) {
  padding: 0.5em 0.75em;
  border: 1px solid var(--border-primary);
  text-align: left;
}

.notes-editor-preview :deep(th) {
  background: var(--bg-panel-muted);
  font-weight: 600;
}
</style>
