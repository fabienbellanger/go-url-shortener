<template>
    <q-card class="q-pa-none q-mt-xl">
        <q-card-section class="q-pa-none">
            <q-layout view="lhh lpr lFf" container style="height: 340px" class="rounded-borders">
                <q-header color="red">
                    <div class="row items-center justify-center text-h5" style="height: 80px">
                        <q-img src="/logo_white.png" class="q-mr-sm" style="height: 64px; max-width: 115px" />
                    </div>
                    <div class="row items-center justify-center text-h5 text-weight-medium" style="height: 40px">
                        URL Shortener
                    </div>
                </q-header>

                <q-page-container>
                    <q-form @submit="signIn">
                        <q-page class="q-pa-md">
                            <div class="q-mx-none q-my-sm">
                                <q-input ref="loginInput" type="email" outlined dense autofocus label="Login" v-model="login">
                                    <template v-slot:before>
                                        <q-icon name="person" />
                                    </template>
                                </q-input>
                            </div>
                            <div class="q-mx-none q-my-sm">
                                <q-input outlined dense type="password" label="Password" v-model="password" >
                                    <template v-slot:before>
                                        <q-icon name="lock" />
                                    </template>
                                </q-input>
                            </div>
                            <div class="q-mx-none q-mt-lg">
                                <q-btn color="primary" label="Se connecter" type="submit" class="full-width" :disable="!valid" />
                            </div>
                        </q-page>
                    </q-form>
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
</template>

<script lang="ts">
import { defineComponent, ref, computed } from 'vue';
import { useRouter } from 'vue-router';
import { useQuasar } from 'quasar';
import { useUserStore } from '../stores/user';
import { AuthUser } from '../api/User';

export default defineComponent({
    name: 'PageLogin',

    setup() {
        const $q = useQuasar();
        const store = useUserStore();
        const router = useRouter();

        const version = process.env.VERSION;
        const loginInput = ref<HTMLInputElement | null>(null);
        const year = ref(new Date().getFullYear());
        const login = ref('');
        const password = ref('');
        const valid = computed(() => login.value !== '' && password.value !== '');

        /**
         * Display authentication error
         *
         * @author Fabien Bellanger
         * @param error Error
         */
        const displayError = (error: Error) => {
            loginInput.value?.focus();

            $q.notify({
                type: 'negative',
                message: 'Wrong login and/or password',
            });
            console.error(error);
        };

        /**
         * Connexion à l'application
         *
         * @author Fabien Bellanger
         */
        const signIn = () => {
            if (valid.value) {
                let user: AuthUser = { email: login.value, password: password.value };

                store.init(user)
                    .then(() => {
                        $q.notify({
                            type: 'positive',
                            message: 'Successfull login',
                        });

                        // Redirect to home route
                        // ----------------------
                        router
                            .push({
                                name: 'home',
                            })
                            .catch((error: Error) => displayError(error));
                    })
                    .catch((error: Error) => displayError(error));
            }
        };

        return {
            loginInput,
            login,
            password,
            year,
            valid,
            version,
            signIn,
        };
    },
});
</script>
