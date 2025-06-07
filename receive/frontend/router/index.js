import { createRouter, createWebHistory } from 'vue-router'

// 导入你的页面组件
import App from '../src/App.vue'
import Receive from "../src/Receive.vue";

const routes = [
    {
        path: '/',
        name: 'Receive',
        component: Receive
    }
]

const router = createRouter({
    history: createWebHistory(),
    routes
})

export default router