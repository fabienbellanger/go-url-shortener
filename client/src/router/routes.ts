import { RouteRecordRaw } from 'vue-router';

const routes: RouteRecordRaw[] = [
    {
        path: '/login',
        component: () => import('layouts/LoginLayout.vue'),
        children: [
            {
                path: '',
                name: 'login',
                component: () => import('pages/Login.vue'),
            }
        ],
    },
    {
        path: '/logout',
        name: 'logout',
        component: () => import('layouts/LoginLayout.vue'),
    },
    {
        path: '/',
        component: () => import('layouts/MainLayout.vue'),
        children: [
            {
                path: '/',
                name: 'home',
                redirect: { name: 'links-list' },
            },
            {
                path: '/links',
                name: 'links-list',
                component: () => import('pages/links/List.vue'),
            },
            // {
            //     path: '/users',
            //     name: 'users-list',
            //     component: () => import('pages/users/List.vue'),
            // },
        ],
    },
    // Always leave this as last one,
    // but you can also remove it
    {
        path: '/:catchAll(.*)*',
        name: '404',
        component: () => import('pages/Error404.vue'),
    },
];

export default routes;
