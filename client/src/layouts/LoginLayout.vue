<template>
    <div class="q-pa-md">
        <div class="row items-center justify-center">
            <div class="col-lg-3 col-3 col-md-4 col-sm-6 col-xs-12">
                <q-card class="q-pa-none q-mt-xl">
                    <q-card-section class="q-pa-none">
                        <q-layout view="lhh lpr lFf" container style="min-height: 400px" class="rounded-borders">
                            <q-header color="red">
                                <div class="row items-center justify-center text-h5" style="height: 80px">
                                    <q-img src="/logo_white.png" class="q-mr-sm" style="height: 64px; max-width: 115px" />
                                </div>
                                <div class="row items-center justify-center text-h5 text-weight-medium" style="height: 40px">
                                    URL Shortener
                                </div>
                            </q-header>

                            <q-page-container>
                                <router-view />
                            </q-page-container>
                        </q-layout>
                    </q-card-section>

                    <q-card-actions class="row justify-center">
                        <div class="text-caption text-grey-6">
                            &copy; {{ year }} <a class="text-grey-6" href="https://www.apitic.com" target="_blank">Apitic</a>
                            <span v-if="version"> - {{ version }}</span>
                        </div>
                    </q-card-actions>
                </q-card>
            </div>
        </div>
    </div>
</template>

<script lang="ts">
import { defineComponent, ref } from 'vue';
import { useQuasar } from 'quasar';

export default defineComponent({
    name: 'LoginLayout',

    setup() {
        const $q = useQuasar();
        const version = process.env.VERSION;
        const year = ref(new Date().getFullYear());
        
        // Theme from OS
        // -------------
        const darkThemeOS = window.matchMedia('(prefers-color-scheme: dark)');
        if (localStorage.getItem('dark-mode') === null) {
            localStorage.setItem('dark-mode', darkThemeOS.matches.toString());
        }

        // Enable Dark mode
        // ----------------
        $q.dark.set(localStorage.getItem('dark-mode') !== 'false');

        return {
            year,
            version
        };
    },
});
</script>
