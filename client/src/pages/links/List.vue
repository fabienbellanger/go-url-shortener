<template>
    <div class="q-px-md">
        <h4 class="q-mt-lg">Links list</h4>
        <q-table
            :rows="projects"
            :columns="headers"
            :filter="filter"
            :pagination="initialPagination"
            row-key="id"
            no-data-label="No link"
            color="primary">
            <template v-slot:body="props">
                <q-tr :props="props">
                    <q-td key="name" :props="props">
                        {{ props.row.name }}
                        <q-popup-edit v-model="props.row.name">
                            <q-input v-model="props.row.name" @change="updateName(props.row.id, props.row.name)" dense autofocus counter />
                        </q-popup-edit>
                    </q-td>
                    <q-td key="expired_at" :props="props">
                        {{ displayDatetime(props.row.expired_at) }}
                    </q-td>
                    <q-td key="actions" :props="props">
                        <q-btn
                            v-if="!displayDeleted"
                            size="sm"
                            icon="delete"
                            color="red"
                            @click="
                                confirmDeleteDialog = true;
                                projectDeleteId = props.row.id;
                            "></q-btn>

                        <q-btn v-else size="sm" icon="refresh" color="warning" @click="recoverProject(props.row.id)"></q-btn>
                    </q-td>
                </q-tr>
            </template>
            <template v-slot:top-right>
                <q-input dense debounce="300" v-model="filter" placeholder="Search">
                    <template v-slot:append>
                        <q-icon name="search" />
                    </template>
                </q-input>
            </template>
            <template v-slot:top-left>
                <q-btn round color="primary" icon="add" @click="confirmCreationDialog = true" />
                &nbsp;
                <q-checkbox v-model="displayDeleted" label="Deleted?" color="primary" @click="getList" />
            </template>
        </q-table>

        <!-- Confirm project deletion dialog -->
        <q-dialog v-model="confirmDeleteDialog" persistent>
            <q-card>
                <q-card-section class="row items-center">
                    <q-icon name="warning" color="warning" class="text-h4" />
                    <span class="q-ml-sm text-h6">Do you really want to delete this link?</span>
                </q-card-section>

                <q-card-actions align="right">
                    <q-btn flat label="Cancel" color="primary" v-close-popup />
                    <q-btn flat label="Delete" color="primary" v-close-popup @click="deleteProject" />
                </q-card-actions>
            </q-card>
        </q-dialog>

        <!-- Project creation dialog -->
        <q-dialog v-model="confirmCreationDialog" medium @hide="clearProjectCreationName">
            <q-card>
                <q-card-section class="row items-center">
                    <span class="q-ml-sm text-h6">Link creation</span>
                </q-card-section>

                <q-card-section>
                    <q-input v-model="projectCreationName" label="Project name" style="width: 240px" autofocus />
                </q-card-section>

                <q-card-actions align="right">
                    <q-btn flat label="Cancel" color="primary" v-close-popup @click="clearProjectCreationName" />
                    <q-btn label="Add" color="primary" v-close-popup @click="addProject" />
                </q-card-actions>
            </q-card>
        </q-dialog>
    </div>
</template>

<script lang="ts">
import { useQuasar } from 'quasar';
import Project from '../../models/Project';
import { defineComponent, ref } from 'vue';
import { ProjectAPI } from '../../api/Project';

export default defineComponent({
    name: 'ProjectsList',

    setup() {
        const $q = useQuasar();
        const projects = ref<Project[]>();
        const confirmDeleteDialog = ref<boolean>(false);
        const confirmCreationDialog = ref<boolean>(false);
        const confirmDeleteProject = ref<boolean>(false);
        const displayDeleted = ref<boolean>(false);
        const projectDeleteId = ref<string>('');
        const projectCreationName = ref<string>('');

        const headers = [
            {
                name: 'id',
                label: 'ID',
                field: 'id',
                align: 'left',
                sortable: true,
            },
            {
                name: 'name',
                label: 'Name',
                field: 'name',
                align: 'left',
                sortable: true,
            },
            {
                name: 'expired_at',
                label: 'Expired at',
                field: 'expired_at',
                align: 'left',
                sortable: true,
                style: 'width: 200px',
            },
            {
                name: 'actions',
                label: 'Actions',
                align: 'left',
                style: 'width: 120px',
                required: true,
            },
        ];

        const displayDatetime = (datetime: string) => {
            if (datetime) {
                return datetime.substr(0, 10) + ' ' + datetime.substr(11, 5);
            }
            return '';
        };

        const clearProjectDeletionId = () => {
            projectDeleteId.value = '';
        };

        const clearProjectCreationName = () => {
            projectCreationName.value = '';
        };

        const deleteProject = () => {
            ProjectAPI.delete(projectDeleteId.value)
                .then(() => {
                    getList();
                    clearProjectDeletionId();

                    $q.notify({
                        type: 'positive',
                        message: 'Successfull project deletion',
                    });
                })
                .catch((error) => {
                    clearProjectDeletionId();

                    $q.notify({
                        type: 'negative',
                        message: 'Error during project deletion',
                    });
                    console.error(error);
                });
        };

        const recoverProject = (id: string) => {
            ProjectAPI.recover(id)
                .then(() => {
                    getList();

                    $q.notify({
                        type: 'positive',
                        message: 'Successfull project recovery',
                    });
                })
                .catch((error) => {
                    $q.notify({
                        type: 'negative',
                        message: 'Error during project recovery',
                    });
                    console.error(error);
                });
        };

        const addProject = () => {
            ProjectAPI.add(projectCreationName.value)
                .then(() => {
                    getList();

                    $q.notify({
                        type: 'positive',
                        message: 'Successfull project creation',
                    });
                })
                .catch((error) => {
                    $q.notify({
                        type: 'negative',
                        message: 'Error during project creation',
                    });
                    console.error(error);
                });
        };

        const updateName = (id: string, name: string) => {
            ProjectAPI.update(id, name)
                .then(() => {
                    $q.notify({
                        type: 'positive',
                        message: 'Successfull project update',
                    });
                })
                .catch((error) => {
                    $q.notify({
                        type: 'negative',
                        message: 'Error during project update',
                    });
                    console.error(error);
                });
        };

        const getList = () => {
            ProjectAPI.list(displayDeleted.value)
                .then((projectsList: Project[]) => {
                    projects.value = projectsList;
                })
                .catch((error) => {
                    console.error(error);
                });
        };
        void getList();

        return {
            projects,
            headers,
            displayDatetime,
            deleteProject,
            recoverProject,
            addProject,
            updateName,
            getList,
            clearProjectCreationName,
            displayDeleted,
            confirmDeleteDialog,
            confirmCreationDialog,
            confirmDeleteProject,
            projectDeleteId,
            projectCreationName,
            filter: ref(''),
            initialPagination: {
                sortBy: 'name',
                descending: false,
            },
        };
    },
});
</script>
