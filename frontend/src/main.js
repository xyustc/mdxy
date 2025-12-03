import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import './styles/main.css'
import { initCopyProtection, addWatermark } from './utils/copyProtection'

const app = createApp(App)
app.use(router)
app.mount('#app')

// 初始化防复制保护
initCopyProtection()

// 添加水印（可选，取消注释启用）
addWatermark('仅供个人学习(XingYu)')

