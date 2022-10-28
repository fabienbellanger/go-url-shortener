<template>
    <div>
        <q-uploader
            @rejected="onRejected"
            @finished="onFinished"
            @uploaded ="onUploaded"
            @failed ="onFailed"
            :factory="uploadFile"
            label="Upload links from CSV file"
            color="primary"
            max-file-size="2097152"
            accept=".csv, text/csv"
            style="max-width: 320px"
      />
    </div>
</template>

<script lang="ts">
import { useQuasar } from 'quasar';
import { defineComponent } from 'vue';
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

        const onRejected = () => {
            $q.notify({
                type: 'negative',
                message: 'Error: Invalid file size (< 2MB) or type (.csv)',
            });
        }

        const onFinished = () => {
            ctx.emit('finished');
            $q.notify({
                type: 'negative',
                message: 'Error: Invalid file size (< 2MB) or type (.csv)',
            });
        }

        const onUploaded = () => {
            console.log('Uploaded');
        }

        const onFailed = () => {
            console.log('Failed');
        }

        return {
            onRejected,
            onFinished,
            onUploaded,
            onFailed,
        }
    }
})
</script>
