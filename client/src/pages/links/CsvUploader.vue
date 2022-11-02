<template>
    <div>
        <q-uploader
            @rejected="onRejected"
            @uploaded ="onUploaded"
            @failed ="onFailed"
            :factory="uploadFile"
            ref="uploader"
            label="Upload links from CSV file (100 lines max)"
            color="primary"
            max-file-size="1048576"
            accept=".csv, text/csv"
            style="max-width: 360px"
        />
    </div>
</template>

<script lang="ts">
import { QUploader, useQuasar } from 'quasar';
import { defineComponent, ref } from 'vue';
import { LinkAPI } from '../../api/Link';

export default defineComponent({
    name: 'CsvUploader',
    emits: ['finished'],
    methods: {
        uploadFile() {
            return LinkAPI.upload();
        }
    },
    setup(_props, ctx) {
        const $q = useQuasar();

        const uploader = ref<QUploader | null>(null);

        const onRejected = () => {
            $q.notify({
                type: 'negative',
                message: 'Error: Invalid file size (< 1MB) or type (.csv)',
            });
        }

        const onFailed = (error) => {
            let response = JSON.parse(error.xhr.response);

            $q.notify({
                type: 'negative',
                message: `Error: ${response?.details}`,
            });
        }

        const onUploaded = (info) => {
            const response = JSON.parse(info.xhr.response);
            const errors = response.errors ?? {};
            const insertedLinks = response.inserted_links ?? 0;

            uploader.value.reset();
            ctx.emit('finished', response);

            $q.notify({
                type: 'positive',
                html: true,
                message: `
                    Inserted links: <b>${insertedLinks}</b><br>
                    Errors: ${errors}
                `,
            });
        }

        return {
            uploader,
            onRejected,
            onUploaded,
            onFailed,
        }
    }
})
</script>
