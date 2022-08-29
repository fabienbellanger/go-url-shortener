<template>
    <div class="text-center q-pa-md flex flex-center">
        <div>
            <div style="font-size: 10vh">404</div>
            <div class="text-h5" style="opacity: 0.4">Page non trouvée</div>

            <q-btn class="q-mt-xl" color="primary" text-color="white" to="/" label="Back" />
        </div>
    </div>
</template>

<script lang="ts">
import { useQuasar } from 'quasar';
import { defineComponent } from 'vue';

export default defineComponent({
    name: 'Error404Page',

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
