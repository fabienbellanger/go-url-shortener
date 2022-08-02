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
                path: '',
                name: 'home',
                component: () => import('pages/Index.vue'),
            },
            {
                path: '/projects',
                name: 'projects-list',
                component: () => import('pages/projects/List.vue'),
            },
            {
                path: '/sales-errors',
                component: () => import('pages/sales-errors/Layout.vue'),
                children: [
                    {
                        path: '',
                        name: 'sales-errors-list',
                        component: () => import('pages/sales-errors/List.vue'),
                    }
                ],
            },
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
