<template>
    <q-layout view="hHh Lpr lFf">
        <q-header elevated>
            <q-toolbar>
                <q-btn flat dense round icon="menu" aria-label="Menu" @click="toggleLeftDrawer" />
                <q-toolbar-title>
                    <q-img src="/logo_white_small.png" class="q-mr-sm" style="height: 32px; width: 57px" />
                    <span>URL Shortener</span>
                </q-toolbar-title>
                <div>
                    <span>Hi {{userStore.user.firstname}},</span>
                    <q-btn flat round :to="{ name: 'logout' }" class="q-mx-sm" icon="logout">
                        <q-tooltip transition-show="scale" transition-hide="scale" >
                            Logout
                        </q-tooltip>
                    </q-btn>
                    <q-btn flat round
                        :icon="$q.dark.isActive ? 'light_mode' : 'dark_mode'"
                        @click="toggleDarkMode">
                        <q-tooltip transition-show="scale" transition-hide="scale" >
                            Change theme
                        </q-tooltip>
                    </q-btn>
                </div>
            </q-toolbar>
        </q-header>

        <q-drawer v-model="leftDrawerOpen" show-if-above bordered :width="220"
            :mini="miniState" @mouseover="miniState = false" @mouseout="miniState = true"
            mini-to-overlay dark>
            <Drawer />
        </q-drawer>

        <q-page-container class="column">
            <div>
                <router-view />
            </div>
            <div class="text-caption text-grey-6 text-center"></div>
        </q-page-container>

        <q-footer class="bg-transparent">
            <q-toolbar class="bg-dark q-mt-md">
                <q-toolbar-title>
                    <div class="text-caption text-grey-6 text-center">
                        &copy; {{ year }} <a class="text-grey-6" href="https://www.apitic.com" target="_blank">Apitic</a>
                        <span v-if="version"> - {{ version }}</span>
                    </div>
                </q-toolbar-title>
            </q-toolbar>
        </q-footer>
    </q-layout>
</template>

<script lang="ts">
import { defineComponent, ref } from 'vue';
import { useQuasar } from 'quasar';
import { useUserStore } from '../stores/user';
import Drawer from 'components/Drawer.vue';

export default defineComponent({
    name: 'MainLayout',

    components: {
        Drawer,
    },

    setup() {
        const $q = useQuasar();
        const leftDrawerOpen = ref(false);
        const miniState = ref(true)
        const year = ref(new Date().getFullYear());
        const userStore = useUserStore();
        const version = process.env.VERSION;

        // Theme from OS
        // -------------
        const darkThemeOS = window.matchMedia('(prefers-color-scheme: dark)');
        if (localStorage.getItem('dark-mode') === null) {
            localStorage.setItem('dark-mode', darkThemeOS.matches.toString());
        }

        // Enable Dark mode
        // ----------------
        $q.dark.set(localStorage.getItem('dark-mode') !== 'false');

        const toggleDarkMode = () => {
            $q.dark.toggle();
            localStorage.setItem('dark-mode', $q.dark.isActive.toString());
        };

        return {
            year,
            leftDrawerOpen,
            miniState,
            userStore,
            version,
            toggleDarkMode,
            toggleLeftDrawer() {
                leftDrawerOpen.value = !leftDrawerOpen.value;
            },
        };
    },
});
</script>

<style scoped>
</style>
