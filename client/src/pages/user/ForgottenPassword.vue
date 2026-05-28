<template>
    <q-form @submit="submit">
        <q-page class="q-pa-md">
            <div class="text-h5 q-mb-md">Forgotten password</div>
            <div class="text-italic text-caption q-mb-md">Enter the email address associated with your account</div>

            <div class="q-mx-none q-my-sm">
                <q-input ref="emailInput" type="email" outlined dense autofocus label="Email" v-model="email">
                    <template v-slot:before>
                        <q-icon name="alternate_email" />
                    </template>
                </q-input>
            </div>

            <div class="q-mx-none q-mt-lg">
                <q-btn color="primary" label="Send" type="submit" class="full-width" :disable="!valid" />
            </div>
            <div class="q-mx-none q-mt-md">
                <q-btn flat color="primary" label="Back to login" class="full-width" 
                    :to="{ name: 'login' }"/>
            </div>
        </q-page>
    </q-form>
</template>

<script lang="ts">
import * as EmailValidator from 'email-validator';
import { useQuasar } from 'quasar';
import { computed, defineComponent, ref } from 'vue';
import { useRouter } from 'vue-router';
import { UserAPI } from '../../api/User';

export default defineComponent({
    name: 'PageForgottenPassword',

    setup() {
        const $q = useQuasar();
        const router = useRouter();
        const emailInput = ref<HTMLInputElement | null>(null);
        const email = ref('');
        const valid = computed(() => email.value !== '' && EmailValidator.validate(email.value));

        /**
         * Display authentication error
         *
         * @author Fabien Bellanger
         * @param error Error
         */
        const displayError = (error: Error) => {
            email.value = '';
            emailInput.value?.focus();

            $q.notify({
                type: 'negative',
                message: `Error during password reset process: ${error}`,
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
                UserAPI.forgottenPassword(email.value)
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
                                    message: 'Successfull password reset. Check your email to change your password',
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
            emailInput,
            email,
            valid,
            submit,
        };
    },
});
</script>