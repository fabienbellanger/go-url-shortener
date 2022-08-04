<template>
    <q-layout view="hHh Lpr lFf">
        <q-header elevated>
            <q-toolbar>
                <q-btn flat dense round icon="menu" aria-label="Menu" @click="toggleLeftDrawer" />
                <q-toolbar-title>Apitic - URL Shortener</q-toolbar-title>
                <div>
                    <q-btn flat label="Logout" :to="{ name: 'logout' }" margin="sm" />
                </div>
            </q-toolbar>
        </q-header>

        <q-drawer v-model="leftDrawerOpen" show-if-above bordered :width="220"
            :mini="miniState" @mouseover="miniState = false" @mouseout="miniState = true"
            mini-to-overlay>
            <Drawer />
        </q-drawer>

        <q-page-container>
            <router-view />
            <div class="text-caption text-grey-6 text-center">&copy; {{ year }} Apitic</div>
        </q-page-container>
    </q-layout>
</template>

<script lang="ts">
import { defineComponent, ref } from 'vue';
import { useQuasar } from 'quasar';
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

        // Enable Dark mode
        // ----------------
        $q.dark.set(true);

        return {
            year,
            leftDrawerOpen,
            miniState,
            toggleLeftDrawer() {
                leftDrawerOpen.value = !leftDrawerOpen.value;
            },
        };
    },
});
</script>
