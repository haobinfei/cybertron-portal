import { createApp } from 'vue'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import 'element-plus/theme-chalk/dark/css-vars.css'
import '@/assets/styles/tech.css'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'

// 启用 Element Plus 暗色模式
document.documentElement.classList.add('dark')

const app = createApp(App)

app.use(ElementPlus)
app.use(createPinia())
app.use(router)

app.mount('#app')
