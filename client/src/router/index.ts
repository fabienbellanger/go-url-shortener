import { route } from 'quasar/wrappers';
import {
    createMemoryHistory,
    createRouter,
    createWebHashHistory,
    createWebHistory,
} from 'vue-router';
import { useUserStore } from 'src/stores/user';
import routes from './routes';

/*
 * If not building with SSR mode, you can
 * directly export the Router instantiation;
 *
 * The function below can be async too; either use
 * async/await or return a Promise which resolves
 * with the Router instance.
 */

export default route(function ({ /*store, ssrContext*/ }) {
    const createHistory = process.env.SERVER
        ? createMemoryHistory
        : process.env.VUE_ROUTER_MODE === 'history'
            ? createWebHistory
            : createWebHashHistory;

    const router = createRouter({
        scrollBehavior: () => ({ left: 0, top: 0 }),
        routes,

        // Leave this as is and make changes in quasar.conf.js instead!
        // quasar.conf.js -> build -> vueRouterMode
        // quasar.conf.js -> build -> publicPath
        history: createHistory(
            process.env.MODE === 'ssr' ? void 0 : process.env.VUE_ROUTER_BASE
        ),
    });

    router.beforeEach((to /* , from*/) => {
        const userStore = useUserStore();
        const isAuthenticated = userStore.isAuthenticated;

        if (to.name !== 'login' && to.name !== 'forgotten-password' && !isAuthenticated) {
            return 'login';
        } else if (to.name === 'logout') {
            userStore.logout();
            return 'login';
        } else if (to.name === 'login' && isAuthenticated) {
            return '';
        }
        return true;
    });

    return router;
});
