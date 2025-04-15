import { createApp } from 'vue'
import { createPinia } from 'pinia'

// Импорт Element Plus и его стилей
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'

import App from './App.vue'
import router from './router'
import './assets/main.css' // Убедитесь, что этот файл не конфликтует с Tailwind/Element Plus

const app = createApp(App)

app.use(createPinia())
app.use(router)

// Используем Element Plus глобально
app.use(ElementPlus)

app.mount('#app')
