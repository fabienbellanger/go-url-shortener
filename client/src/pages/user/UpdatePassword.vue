<template>
    <q-form @submit="submit">
        <q-page class="q-pa-md">
            <div class="text-h5 q-mb-md">New password</div>
            
            <div class="q-mx-none q-my-sm">
                <q-input ref="passwordInput" type="password" outlined dense autofocus label="New password" v-model="password">
                    <template v-slot:before>
                        <q-icon name="lock" />
                    </template>
                </q-input>
            </div>

            <div class="q-mx-none q-my-sm">
                <q-input ref="password2Input" type="password" outlined dense autofocus label="New password confirmation" v-model="password2">
                    <template v-slot:before>
                        <q-icon name="lock" />
                    </template>
                </q-input>
            </div>
            <div class="q-mx-none q-mt-lg">
                <q-btn color="primary" label="Save" type="submit" class="full-width" :disable="!valid" />
            </div>
            <div class="q-mx-none q-mt-md">
                <q-btn flat color="primary" label="Back to login" class="full-width" 
                    :to="{ name: 'login' }"/>
            </div>
        </q-page>
    </q-form>
</template>

<script lang="ts">
import { defineComponent, ref, computed } from 'vue';
import { useQuasar } from 'quasar';
import { useRouter, useRoute } from 'vue-router';
import { UserAPI } from '../../api/User';

export default defineComponent({
    name: 'PageUpdatePassword',

    setup() {
        const $q = useQuasar();
        const router = useRouter();
        const route = useRoute();
        const passwordInput = ref<HTMLInputElement | null>(null);
        const password2Input = ref<HTMLInputElement | null>(null);
        const password = ref('');
        const password2 = ref('');
        const valid = computed(() => password.value.length >= 8 && password2.value.length >= 8 && password.value === password2.value);
        const token = route.params.token;

        /**
         * Display authentication error
         *
         * @author Fabien Bellanger
         * @param error Error
         */
        const displayError = (error: Error) => {
            password.value = '';
            password2.value = '';
            passwordInput.value?.focus();

            $q.notify({
                type: 'negative',
                message: 'Error when enter a new password: ' + error,
            });
            console.error(error);
        };

        /**
         * Envoi la demande e réinitialisation du mot de passe
         *
         * @author Fabien Bellanger
         */
        const submit = () => {
            if (valid.value) {
                UserAPI.updatePassword(token.toString(), password.value)
                    .then(() => {
                        // Redirect to login route
                        // -----------------------
                        router
                            .push({
                                name: 'login',
                            })
                            .then(() => {
                                $q.notify({
                                    type: 'positive',
                                    message: 'Successfull password update',
                                });
                            })
                            .catch((error: Error) => displayError(error));
                    })
                    .catch((error: Error) => {
                        displayError(error);
                    });
            }
        };

        return {
            passwordInput,
            password2Input,
            password,
            password2,
            valid,
            submit,
        };
    }
});
</script>