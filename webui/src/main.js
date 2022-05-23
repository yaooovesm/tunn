import {createApp} from 'vue'
import App from './App.vue'

//router
import router from "./router"

//utils
import utils from "@/utils";

//axios
import axios from 'axios'
import VueAxios from 'vue-axios'

//element
import element from 'element-plus';
import 'element-plus/dist/index.css'

const app = createApp(App);
//注册工具类
app.config.globalProperties.$utils = utils
app
    //router
    .use(router)
    //element
    .use(element)
    //axios
    .use(VueAxios, axios)
//挂载
app.mount('#app');
