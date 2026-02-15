<template>
  <div class="profile-edit-page">
    <h1 class="page-title">编辑个人信息</h1>

    <el-card>
      <el-form
        ref="formRef"
        :model="form"
        label-width="120px"
        :rules="rules"
      >
        <el-form-item label="姓名" prop="name">
          <el-input v-model="form.name" placeholder="请输入姓名" />
        </el-form-item>

        <el-form-item label="职位" prop="title">
          <el-input v-model="form.title" placeholder="请输入职位" />
        </el-form-item>

        <el-form-item label="个人简介" prop="bio">
          <el-input
            v-model="form.bio"
            type="textarea"
            :rows="4"
            placeholder="请输入个人简介"
          />
        </el-form-item>

        <el-form-item label="头像URL" prop="avatar">
          <el-input v-model="form.avatar" placeholder="请输入头像URL" />
        </el-form-item>

        <el-form-item label="邮箱" prop="email">
          <el-input v-model="form.email" placeholder="请输入邮箱" />
        </el-form-item>

        <el-form-item label="GitHub" prop="github">
          <el-input v-model="form.github" placeholder="请输入GitHub链接" />
        </el-form-item>

        <el-form-item label="LinkedIn" prop="linkedin">
          <el-input v-model="form.linkedin" placeholder="请输入LinkedIn链接" />
        </el-form-item>

        <el-form-item label="Twitter" prop="twitter">
          <el-input v-model="form.twitter" placeholder="请输入Twitter链接" />
        </el-form-item>

        <el-form-item label="个人网站" prop="website">
          <el-input v-model="form.website" placeholder="请输入个人网站链接" />
        </el-form-item>

        <el-form-item label="技能" prop="skills">
          <el-input
            v-model="form.skills"
            type="textarea"
            :rows="3"
            placeholder="请输入技能（JSON格式数组）"
          />
        </el-form-item>

        <el-form-item>
          <el-button type="primary" :loading="loading" @click="handleSubmit">
            保存
          </el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElCard, ElForm, ElFormItem, ElInput, ElButton, ElMessage } from 'element-plus'
import { profileApi } from '@/api/profile'
import type { Profile } from '@/api/types'

const formRef = ref()
const loading = ref(false)

const form = reactive({
  name: '',
  title: '',
  bio: '',
  avatar: '',
  email: '',
  github: '',
  linkedin: '',
  twitter: '',
  website: '',
  skills: ''
})

const rules = {
  name: [{ required: true, message: '请输入姓名', trigger: 'blur' }]
}

const fetchProfile = async () => {
  try {
    const response = await profileApi.get()
    if (response.success && response.data) {
      Object.assign(form, response.data)
    }
  } catch (error) {
    ElMessage.error('获取个人信息失败')
  }
}

const handleSubmit = async () => {
  try {
    await formRef.value.validate()
    loading.value = true

    const response = await profileApi.update(form)
    if (response.success) {
      ElMessage.success('保存成功')
    }
  } catch (error: any) {
    ElMessage.error(error.response?.data?.error || '保存失败')
  } finally {
    loading.value = false
  }
}

const handleReset = () => {
  fetchProfile()
}

onMounted(() => {
  fetchProfile()
})
</script>

<style scoped>
.profile-edit-page {
  padding: 20px;
}

.page-title {
  font-size: 24px;
  font-weight: 600;
  margin-bottom: 24px;
}
</style>
