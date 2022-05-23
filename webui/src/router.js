import {createRouter, createWebHashHistory} from 'vue-router';

import LoginPage from "@/components/login/LoginPage";

const routers = [
    {
        path: '/',
        name: 'login',
        alias: "/login",
        component: LoginPage,
    }
]

const router = createRouter({
    history: createWebHashHistory(),
    routes: routers,
})

export default router