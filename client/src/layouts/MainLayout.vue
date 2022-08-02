<template>
    <q-layout view="hHh Lpr lFf">
        <q-header elevated>
            <q-toolbar>
                <q-btn flat dense round icon="menu" aria-label="Menu" @click="toggleLeftDrawer" />
                <q-toolbar-title>POS Tools</q-toolbar-title>
                <div>
                    <q-btn color="secondary" label="Logout" :to="{ name: 'logout' }" margin="sm" />
                </div>
            </q-toolbar>
        </q-header>

        <q-drawer v-model="leftDrawerOpen" show-if-above bordered :width="220">
            <Drawer />
        </q-drawer>

        <q-page-container>
            <router-view />
        </q-page-container>
    </q-layout>
</template>

<script lang="ts">
import { defineComponent, ref } from 'vue';
import { useQuasar } from 'quasar';
import Drawer from 'components/Drawer.vue';

const linksList = [
    {
        title: 'Sales errors',
        caption: '',
        icon: 'error_outline',
        link: 'http://localhost:8081',
    },
];

export default defineComponent({
    name: 'MainLayout',

    components: {
        Drawer,
    },

    setup() {
        const leftDrawerOpen = ref(false);
        const $q = useQuasar();

        // Enable Dark mode
        // ----------------
        $q.dark.set(true);

        return {
            linksList: linksList,
            leftDrawerOpen,
            toggleLeftDrawer() {
                leftDrawerOpen.value = !leftDrawerOpen.value;
            },
        };
    },
});
</script>
