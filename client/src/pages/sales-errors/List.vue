<template>
    <div class="q-px-md">
        <h4 class="q-mt-lg">Sales errors list</h4>

        <q-table
            :rows="salesErrors"
            :columns="headers"
            :filter="filter"
            :pagination="initialPagination"
            row-key="id"
            no-data-label="No error"
            color="primary">
            <template v-slot:body="props">
                <q-tr :props="props">
                    <q-td key="project" :props="props">{{ getProjectNameFromId(props.row.project_id) }}</q-td>
                    <q-td key="account" :props="props">{{ props.row.account_id }}</q-td>
                    <q-td key="sale_id" :props="props">{{ props.row.sale_id }}</q-td>
                    <q-td key="sale_ticket_number" :props="props">{{ props.row.sale_ticket_number }}</q-td>
                    <q-td key="sale_datetime" :props="props">{{ displayDatetime(props.row.sale_datetime) }}</q-td>
                    <q-td key="platform" :props="props">{{ props.row.platform }}</q-td>
                    <q-td key="gap" :props="props">{{ props.row.gap }}</q-td>
                    <q-td key="updated_at" :props="props">
                        {{ displayDatetime(props.row.updated_at) }}
                    </q-td>
                    <q-td key="actions" :props="props"></q-td>
                </q-tr>
            </template>
            <template v-slot:top-right>
                <q-input dense debounce="300" v-model="filter" placeholder="Search">
                    <template v-slot:append>
                        <q-icon name="search" />
                    </template>
                </q-input>
            </template>
        </q-table>
    </div>
</template>

<script lang="ts">
// import { useQuasar } from 'quasar';
import SaleError from '../../models/SaleError';
import Project from '../../models/Project';
import { defineComponent, ref } from 'vue';
import { SalesErrorAPI } from '../../api/SalesError';
import { ProjectAPI } from '../../api/Project';
import { Url } from '../../services/Url';

// TODO: Get accounts list

export default defineComponent({
    name: 'SalesErrorsList',

    setup() {
        // const $q = useQuasar();
        const salesErrors = ref<SaleError[]>();
        const projects = ref<Project[]>();

        const headers = [
            {
                name: 'project',
                label: 'Project',
                field: 'project',
                align: 'left',
                sortable: true,
            },
            {
                name: 'account',
                label: 'Account',
                field: 'account',
                align: 'left',
                sortable: true,
            },
            {
                name: 'sale_id',
                label: 'Sale ID',
                field: 'sale_id',
                align: 'left',
                sortable: true,
                style: 'width: 100px',
            },
            {
                name: 'sale_ticket_number',
                label: 'Sale number',
                field: 'sale_ticket_number',
                align: 'center',
                sortable: true,
                style: 'width: 100px',
            },
            {
                name: 'sale_datetime',
                label: 'Sale datetime',
                field: 'sale_datetime',
                align: 'center',
                sortable: true,
                style: 'width: 180',
            },
            {
                name: 'platform',
                label: 'Platform',
                field: 'platform',
                align: 'left',
                sortable: true,
            },
            {
                name: 'gap',
                label: 'Gap',
                field: 'gap',
                align: 'right',
                sortable: true,
                style: 'width: 80',
            },
            {
                name: 'updated_at',
                label: 'Updated at',
                field: 'updated_at',
                align: 'left',
                sortable: true,
                style: 'width: 180px',
            },
            {
                name: 'actions',
                label: 'Actions',
                align: 'center',
                style: 'width: 120px',
            },
        ];

        // TODO: Factoriser
        const displayDatetime = (datetime: string) => {
            if (datetime) {
                return datetime.substr(0, 10) + ' ' + datetime.substr(11, 5);
            }
            return '';
        };

        const getProjectNameFromId = (id: string): string => {
            return Project.getNameFromId(projects.value, id);
        };

        const getProjectsList = () => {
            ProjectAPI.listAll()
                .then((projectsList: Project[]) => {
                    projects.value = projectsList;
                })
                .catch((error) => {
                    console.error(error);
                });
        };

        const getList = () => {
            const parameters = [
                { name: 'date_begin', value: '2022-01-26' },
                { name: 'date_end', value: '2022-01-26' },
                { name: 'project_id', value: 'my-project-id' },
            ];
            console.log(Url.generate('sales-errors', parameters));

            SalesErrorAPI.list()
                .then((data: SaleError[]) => {
                    salesErrors.value = data;
                })
                .catch((error) => {
                    console.error(error);
                });
        };
        void getProjectsList();
        void getList();

        return {
            displayDatetime,
            getList,
            getProjectNameFromId,
            headers,
            salesErrors,
            projects,
            filter: ref(''),
            initialPagination: {
                sortBy: 'updated_at',
                descending: false,
            },
        };
    },
});
</script>
