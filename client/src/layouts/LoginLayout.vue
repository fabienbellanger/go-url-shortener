<template>
    <div class="q-pa-md">
        <div class="row items-center justify-center">
            <div class="col-lg-3 col-3 col-md-4 col-sm-6 col-xs-12">
                <router-view />
            </div>
        </div>
    </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue';
import { useQuasar } from 'quasar';

export default defineComponent({
    name: 'LoginLayout',

    setup() {
        const $q = useQuasar();
        
        // Theme from OS
        // -------------
        const darkThemeOS = window.matchMedia('(prefers-color-scheme: dark)');
        if (localStorage.getItem('dark-mode') === null) {
            localStorage.setItem('dark-mode', darkThemeOS.matches.toString());
        }

        // Enable Dark mode
        // ----------------
        $q.dark.set(localStorage.getItem('dark-mode') !== 'false');

        return {};
    },
});
</script>
