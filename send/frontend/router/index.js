import { createRouter, createWebHistory } from 'vue-router'

// 导入你的页面组件
import App from '../src/App.vue'
import Send from "../src/Send.vue";

const routes = [
    {
        path: '/',
        name: 'Send',
        component: Send
    }
]

const router = createRouter({
    history: createWebHistory(),
    routes
})

export default router