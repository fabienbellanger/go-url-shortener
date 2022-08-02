<template>
    <div class="text-white text-center q-pa-md flex flex-center">
        <div>
            <div style="font-size: 10vh">404</div>
            <div class="text-h5" style="opacity: 0.4">Page non trouvée</div>

            <q-btn class="q-mt-xl" color="primary" text-color="white" to="/" label="Back" />
        </div>
    </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue';

export default defineComponent({
    name: 'Error404Page',
});
</script>
