<template>
    <q-dialog v-model="open" persistent>
        <q-card>
            <q-card-section class="row items-center">
                <q-icon name="warning" color="warning" class="text-h4" />
                <span class="q-ml-sm text-h6">Links import results</span>
            </q-card-section>

            <q-card-actions align="right">
                <q-btn flat label="Close" color="primary" v-close-popup />
            </q-card-actions>
        </q-card>
    </q-dialog>
</template>

<script lang="ts">
import { defineComponent, ref } from 'vue';

export default defineComponent({
    name: 'ImportLinksDialog',
    props: {
        results: Boolean,
        visible: Object,
    },
    setup(props) {
        const open = ref(props.visible);

        console.log(props.results);
        console.log(props.visible);

        return {
            open
        }
    }
})
</script>
