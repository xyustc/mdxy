import { createApp } from 'vue'
import { createPinia } from 'pinia'
import ElementPlus from 'element-plus'
import router from './router'
import App from './App.vue'
import { overflowTooltipDirective } from './directives/overflowTooltip'
import { useAppStore } from './stores/app'
import { useSiteStore } from './stores/site'

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

const appStore = useAppStore(pinia)
const siteStore = useSiteStore(pinia)

appStore.initTheme()
siteStore.hydrateFromCache()
void siteStore.ensureProfile()
void siteStore.ensureHome()

app.mount('#app')
