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
            style="max-width: 320px"
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

        const onUploaded = () => {
            console.log('Uploaded');
            uploader.value.reset();

            ctx.emit('finished');
        }

        const onFailed = () => {
            console.log('Failed');
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
