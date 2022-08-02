<template>
    <q-list>
        <q-item-label header>Sales errors</q-item-label>

        <q-item clickable tag="a" :to="link.url" v-for="link in links" :key="link.title" v-bind="link">
            <q-item-section v-if="link.icon" avatar>
                <q-icon :name="link.icon" />
            </q-item-section>

            <q-item-section>
                <q-item-label>{{ link.title }}</q-item-label>
                <q-item-label v-if="link.caption" q-item-label caption>
                    {{ link.caption }}
                </q-item-label>
            </q-item-section>
        </q-item>
    </q-list>
</template>

<script lang="ts">
import { defineComponent } from 'vue';

const links = [
    {
        title: 'Projects',
        caption: '',
        icon: 'dns',
        url: '/projects',
    },
    {
        title: 'Sales errors',
        caption: '',
        icon: 'error_outline',
        url: '/sales-errors',
    },
];

export default defineComponent({
    name: 'DrawerList',
    setup() {
        return {
            links,
        };
    },
});
</script>
