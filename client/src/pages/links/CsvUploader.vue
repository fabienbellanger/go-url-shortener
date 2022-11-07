<template>
    <div>
        <q-uploader
            @rejected="onRejected"
            @uploaded ="onUploaded"
            @failed ="onFailed"
            :factory="uploadFile"
            ref="uploader"
            label="Upload links from CSV file"
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
    emits: ['uploaded'],
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
            const errors = response.errors ?? [];
            const insertedLinks = response.inserted_links ?? 0;

            let notifyType = 'positive';
            if (errors.length > 0) {
                if (insertedLinks == 0) {
                    notifyType = 'negative';
                } else {
                    notifyType = 'warning';
                }
            } else if (insertedLinks == 0) {
                notifyType = 'warning';
            }

            uploader.value.reset();
            ctx.emit('uploaded', insertedLinks, errors);

            $q.notify({
                type: notifyType,
                html: true,
                message: `
                    Inserted links: <b>${insertedLinks}</b><br>
                    Errors: <b>${errors.length}</b>
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
