<template>
    <q-dialog medium persistent>
        <q-card style="width: 600px">
            <q-card-section class="row items-center">
                <span class="q-ml-sm text-h6">Links import</span>
            </q-card-section>

            <q-card-section v-if="showUploader">
                <div class="q-mb-sm">
                    <div class="row items-center">
                        <q-icon name="warning" color="warning" class="text-h5" />
                        <span class="text-subtitle1 q-ml-sm">Information about CSV file</span>
                    </div>
                    <ul>
                        <li>Format: "URL";"Name";"Expiration datetime (YYYY-MM-DD HH:MM:SS)"</li>
                        <li>File must have an header at line 1 ("URL";"Name";"Expiration datetime")</li>
                        <li>; as separator</li>
                        <li>1MB max for file size</li>
                        <li>100 lines max (with header)</li>
                    </ul>
                </div>

                <q-separator />

                <div class="row justify-center q-mt-md">
                    <csv-uploader @uploaded="uploadFinished"></csv-uploader>
                </div>
            </q-card-section>

            <q-card-section v-else-if="errors.length > 0">
                <div class="row items-center q-mb-sm">
                    <q-icon name="error" color="red" class="text-h5" />
                    <span class="text-subtitle1 q-ml-sm">Errors</span>
                </div>
                <q-list bordered separator>
                    <q-item v-for="error in errors" :key="error.line">
                        <q-item-section>
                            <q-item-label caption lines="1">Line {{ error.line }} | {{ error.err }}</q-item-label>
                            <q-item-label lines="3">{{ error.data}}</q-item-label>
                        </q-item-section>
                    </q-item>
                </q-list>
            </q-card-section>

            <q-card-actions align="right">
                <q-btn flat label="Close" color="primary" @click="close" />
            </q-card-actions>
        </q-card>
    </q-dialog>
</template>

<script lang="ts">
import { defineComponent, ref } from 'vue';
import CsvUploader from '../CsvUploader.vue';

export default defineComponent({
    name: 'ImportLinksDialog',
    emits: ['close'],
    components: { CsvUploader },
    setup(_props, ctx) {
        const showUploader = ref(true);
        const errors = ref([]);
        const reload = ref(false);

        const uploadFinished = (insertedLinks, errorsList) => {
            // Hide uploader
            showUploader.value = false;

            // Open dialog with errors
            if (errorsList.length > 0) {
                errors.value = errorsList;
                showUploader.value = false;

                if (insertedLinks > 0) {
                    reload.value = true;
                }
            } else {
                reload.value = (insertedLinks > 0);

                close();
            }
        }

        const close = () => {
            // Close dialog and reload links list if links inserted
            ctx.emit('close', reload.value);

            showUploader.value = true;
            errors.value = [];
            reload.value = false;
        }

        return {
            showUploader,
            errors,
            uploadFinished,
            close,
        }
    }
})
</script>
