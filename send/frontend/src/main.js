import {createApp} from 'vue'
import App from './App.vue'
import './style.css';
import Vant from 'vant';
import 'vant/lib/index.css';
import router from "../router/index.js";

createApp(App).use(router).use(Vant).mount('#app')
