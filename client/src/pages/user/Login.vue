<template>
    <q-form @submit="signIn">
        <q-page class="q-pa-md">
            <div class="text-h5 q-mb-md">Connection</div>
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
                <q-btn color="primary" label="Sign in" type="submit" class="full-width" :disable="!valid" />
            </div>
            <div class="q-mx-none q-mt-md">
                <q-btn flat no-caps color="primary" label="Forgotten password" class="full-width" 
                    :to="{ name: 'forgotten-password' }"/>
            </div>
        </q-page>
    </q-form>
</template>

<script lang="ts">
import { defineComponent, ref, computed } from 'vue';
import { useRouter } from 'vue-router';
import { useQuasar } from 'quasar';
import * as EmailValidator from 'email-validator';
import { useUserStore } from '../../stores/user';
import { AuthUser } from '../../api/User';

export default defineComponent({
    name: 'PageLogin',

    setup() {
        const $q = useQuasar();
        const store = useUserStore();
        const router = useRouter();

        const loginInput = ref<HTMLInputElement | null>(null);
        const login = ref('');
        const password = ref('');
        const valid = computed(() => login.value !== '' 
            && EmailValidator.validate(login.value) 
            && password.value.length >= 8);

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
            valid,
            signIn,
        };
    },
});
</script>
