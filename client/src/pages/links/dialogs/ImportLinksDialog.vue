<template>
    <q-dialog persistent>
        <q-card>
            <q-card-section class="row items-center">
                <q-icon name="warning" color="warning" class="text-h4" />
                <span class="q-ml-sm text-h6">Links not imported</span>
            </q-card-section>

            <q-card-section>
                <q-list bordered separator>
                    <q-item v-for="(error, line) in errors" :key="line">
                        <q-item-section>
                            <q-item-label caption>Line {{ line }}</q-item-label>
                            <q-item-label>{{ error}}</q-item-label>
                        </q-item-section>
                    </q-item>
                </q-list>
            </q-card-section>

            <q-card-actions align="right">
                <q-btn flat label="Close" color="primary" v-close-popup />
            </q-card-actions>
        </q-card>
    </q-dialog>
</template>

<script lang="ts">
import { defineComponent } from 'vue';

export default defineComponent({
    name: 'ImportLinksDialog',
    props: {
        errors: Object,
    },
    setup(props) {

        console.log(props.errors);

        return {
        }
    }
})
</script>
