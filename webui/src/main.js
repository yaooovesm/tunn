import { createApp } from 'vue'
import App from './App.vue'

//router
import router from "./router"

//axios
import axios from 'axios'
import VueAxios from 'vue-axios'

//element
import element from 'element-plus';
import 'element-plus/dist/index.css'

const app = createApp(App);
app
    //router
    .use(router)
    //element
    .use(element)
    //axios
    .use(VueAxios, axios)
//挂载
app.mount('#app');
