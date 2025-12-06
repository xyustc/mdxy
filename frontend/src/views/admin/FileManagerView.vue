<template>
  <div class="file-manager">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>文件管理</span>
          <div>
            <el-button type="primary" @click="showUploadDialog">
              <el-icon><Upload /></el-icon>
              上传文件
            </el-button>
          </div>
        </div>
      </template>
      
      <el-table
        :data="noteTree"
        row-key="path"
        default-expand-all
        :tree-props="{ children: 'children', hasChildren: 'hasChildren' }"
        style="width: 100%"
      >
        <el-table-column prop="name" label="名称" min-width="300">
          <template #default="{ row }">
            <el-icon v-if="row.type === 'directory'"><Folder /></el-icon>
            <el-icon v-else><Document /></el-icon>
            <span style="margin-left: 10px">{{ row.name }}</span>
          </template>
        </el-table-column>
        
        <el-table-column prop="type" label="类型" width="120">
          <template #default="{ row }">
            <el-tag v-if="row.type === 'directory'" type="warning">文件夹</el-tag>
            <el-tag v-else type="success">文件</el-tag>
          </template>
        </el-table-column>
        
        <el-table-column prop="path" label="路径" min-width="200" />
        
        <el-table-column label="操作" width="150">
          <template #default="{ row }">
            <el-button
              type="danger"
              size="small"
              @click="handleDelete(row)"
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
    
    <!-- 上传对话框 -->
    <el-dialog
      v-model="uploadDialogVisible"
      title="上传文件"
      width="500px"
    >
      <el-form :model="uploadForm" label-width="80px">
        <el-form-item label="文件">
          <el-upload
            ref="uploadRef"
            :auto-upload="false"
            :limit="1"
            :on-change="handleFileChange"
            accept=".md,.markdown"
            drag
          >
            <el-icon class="el-icon--upload"><UploadFilled /></el-icon>
            <div class="el-upload__text">
              拖拽文件到此处或<em>点击上传</em>
            </div>
            <template #tip>
              <div class="el-upload__tip">
                仅支持 .md 和 .markdown 文件，最大 10MB
              </div>
            </template>
          </el-upload>
        </el-form-item>
        
        <el-form-item label="上传路径">
          <el-input
            v-model="uploadForm.path"
            placeholder="留空上传到根目录，如: folder/subfolder"
          />
        </el-form-item>
      </el-form>
      
      <template #footer>
        <el-button @click="uploadDialogVisible = false">取消</el-button>
        <el-button
          type="primary"
          :loading="uploading"
          @click="handleUpload"
        >
          上传
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getNoteTree } from '../../api/notes'
import { uploadNote, deleteNote } from '../../api/admin'
import {
  Upload,
  UploadFilled,
  Folder,
  Document
} from '@element-plus/icons-vue'

const noteTree = ref([])
const uploadDialogVisible = ref(false)
const uploading = ref(false)
const uploadRef = ref(null)
const uploadForm = ref({
  file: null,
  path: ''
})

const loadNoteTree = async () => {
  try {
    const response = await getNoteTree()
    if (response.success) {
      noteTree.value = response.data
    }
  } catch (error) {
    ElMessage.error('加载文件列表失败')
  }
}

const showUploadDialog = () => {
  uploadForm.value = { file: null, path: '' }
  uploadDialogVisible.value = true
}

const handleFileChange = (file) => {
  uploadForm.value.file = file.raw
}

const handleUpload = async () => {
  if (!uploadForm.value.file) {
    ElMessage.warning('请选择文件')
    return
  }
  
  uploading.value = true
  
  try {
    const response = await uploadNote(
      uploadForm.value.file,
      uploadForm.value.path
    )
    
    if (response.success) {
      ElMessage.success('上传成功')
      uploadDialogVisible.value = false
      await loadNoteTree()
    } else {
      ElMessage.error(response.message || '上传失败')
    }
  } catch (error) {
    ElMessage.error(error.response?.data?.detail || '上传失败')
  } finally {
    uploading.value = false
  }
}

const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除 ${row.name} 吗？此操作不可恢复！`,
      '警告',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    const response = await deleteNote(row.path)
    
    if (response.success) {
      ElMessage.success('删除成功')
      await loadNoteTree()
    } else {
      ElMessage.error(response.message || '删除失败')
    }
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(error.response?.data?.detail || '删除失败')
    }
  }
}

onMounted(() => {
  loadNoteTree()
})
</script>

<style scoped>
.file-manager {
  width: 100%;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.el-icon--upload {
  font-size: 67px;
  color: #8c939d;
  margin: 40px 0 16px;
}
</style>
