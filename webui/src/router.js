import {createRouter, createWebHashHistory} from 'vue-router';

import LoginPage from "@/components/login/LoginPage";
import DashboardPage from "@/components/dashboard/DashboardPage";
//import axios from "axios";

const routers = [
    {
        path: '/',
        name: 'default',
        component: LoginPage,
    },
    {
        path: '/login',
        name: 'login',
        component: LoginPage,
    },
    {
        path: '/dashboard',
        name: 'dashboard',
        component: DashboardPage,
    }
]

const router = createRouter({
    history: createWebHashHistory(),
    routes: routers,
})

// router.beforeEach((to, from, next) => {
//     if (to.name !== "loading") {
//         axios({
//             method: "get",
//             url: "/api/application/",
//             data: {}
//         }).then((res) => {
//             let response = res.data
//             console.log(response)
//             console.log(to)
//             if (!response.data.running && to.name === "dashboard") {
//                 router.push({path: "/login"})
//             } else {
//                 next()
//             }
//         }).catch(() => {
//             router.push({path: "/"})
//         })
//     } else {
//         next()
//     }
// })

export default router