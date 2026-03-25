import { createApp } from 'vue'
import { createPinia } from 'pinia'
import ElementPlus from 'element-plus'
import router from './router'
import App from './App.vue'
import { overflowTooltipDirective } from './directives/overflowTooltip'

// 样式
import 'element-plus/dist/index.css'
import './styles/fonts.css'
import './styles/variables.css'
import './styles/typography.css'
import './styles/global.css'

const app = createApp(App)
const pinia = createPinia()

app.use(pinia)
app.use(ElementPlus)
app.use(router)
app.directive('overflow-tooltip', overflowTooltipDirective)

app.mount('#app')
