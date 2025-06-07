import {createApp} from 'vue'
import App from './App.vue'
import Router from "../router/index.js";
import Vant from 'vant';
import 'vant/lib/index.css';
createApp(App).use(Router).use(Vant).mount('#app')
