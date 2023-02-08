<template>
    <div class="q-px-md">
        <h4 class="q-mt-lg">Links list</h4>

        <q-table
            @request="getList"
            :rows="links"
            :columns="headers"
            :filter="filter"
            :rows-per-page-options="[25, 50, 100, 200, 500]"
            :loading="loading"
            v-model:pagination="pagination"
            row-key="id"
            no-data-label="No link"
            color="primary"
            binary-state-sort>
            <template v-slot:body="props">
                <q-tr :props="props">
                    <q-td key="id" :props="props" @click="copyLink(props.row.id)" style="cursor: pointer">
                        <q-tooltip transition-show="scale" transition-hide="scale" >
                            Copy link to clipboard
                        </q-tooltip>
                        {{ props.row.id }}
                    </q-td>
                    <q-td key="name" :props="props">
                        <div class="ellipsis">
                            {{ props.row.name }}
                            <q-tooltip>{{ props.row.name }}</q-tooltip>
                        </div>
                    </q-td>
                    <q-td key="url" :props="props">
                        <div class="ellipsis">
                            {{ props.row.url }}
                            <q-tooltip>{{ props.row.url }}</q-tooltip>
                        </div>
                    </q-td>
                    <q-td key="expired_at" :props="props">
                        {{ formatDatetime(props.row.expired_at) }}
                    </q-td>
                    <q-td key="created_at" :props="props">
                        {{ formatDatetime(props.row.created_at) }}
                    </q-td>
                    <q-td key="actions" :props="props">
                        <q-btn
                            size="sm"
                            icon="content_copy"
                            color="orange"
                            @click="copyLink(props.row.id)">
                            <q-tooltip transition-show="scale" transition-hide="scale" >
                                Copy link to clipboard
                            </q-tooltip>
                        </q-btn>
                        &nbsp;
                        <q-btn
                            size="sm"
                            icon="link"
                            color="blue"
                            @click="
                                currentLink = props.row;
                                openLink();
                            ">
                            <q-tooltip transition-show="scale" transition-hide="scale" >
                                Open short link
                            </q-tooltip>
                        </q-btn>
                        &nbsp;
                        <q-btn
                            size="sm"
                            icon="edit"
                            color="green"
                            @click="
                                currentLink = props.row;
                                confirmCreationDialog = true;
                            ">
                            <q-tooltip transition-show="scale" transition-hide="scale" >
                                Update link
                            </q-tooltip>
                        </q-btn>
                        &nbsp;
                        <q-btn
                            size="sm"
                            icon="delete"
                            color="red"
                            @click="
                                currentLink = props.row;
                                confirmDeleteDialog = true;
                            ">
                            <q-tooltip transition-show="scale" transition-hide="scale" >
                                Delete link
                            </q-tooltip>
                        </q-btn>
                    </q-td>
                </q-tr>
            </template>
            <template v-slot:top-right>
                <q-input clearable dense debounce="300" v-model="filter" placeholder="Search" class="search_input">
                    <template v-slot:prepend>
                        <q-icon name="search" />
                    </template>
                </q-input>
            </template>
            <template v-slot:top-left>
                <q-btn round color="primary" icon="add" @click="newLink">
                    <q-tooltip transition-show="scale" transition-hide="scale" >
                        Add a new link
                    </q-tooltip>
                </q-btn>
                <q-btn
                    color="primary"
                    icon-right="file_upload"
                    label="Import csv"
                    class="q-mx-md"
                    @click="showUploaderDialog = true"
                    no-caps></q-btn>
                <q-btn
                    color="primary"
                    icon-right="file_download"
                    label="Export to csv"
                    class="q-mx-md"
                    @click="exportCSV"
                    no-caps></q-btn>
            </template>
        </q-table>

        <!-- Confirm link deletion dialog -->
        <q-dialog v-model="confirmDeleteDialog" persistent>
            <q-card>
                <q-card-section class="row items-center">
                    <q-icon name="warning" color="warning" class="text-h4" />
                    <span class="q-ml-sm text-h6">Do you really want to delete this link?</span>
                </q-card-section>

                <q-card-actions align="right">
                    <q-btn flat label="Cancel" color="primary" v-close-popup />
                    <q-btn flat label="Delete" color="primary" v-close-popup @click="deleteLink" />
                </q-card-actions>
            </q-card>
        </q-dialog>

        <!-- link creation dialog -->
        <q-dialog v-model="confirmCreationDialog" medium @hide="clearLinkCreation">
            <q-card>
                <q-form @submit="editLink">
                    <q-card-section class="row items-center">
                        <span class="q-ml-sm text-h6">
                            <span v-if="currentLink.id">Link update ({{ currentLink.id }})</span>
                            <span v-else>Link creation</span>
                        </span>
                    </q-card-section>

                    <q-card-section>
                        <q-input v-model="currentLink.url" label="URL*" style="width: 320px" autofocus type="url" 
                            :rules="[val => (val.startsWith('http://') || val.startsWith('https://')) || 'URL is required']"/>
                    </q-card-section>

                    <q-card-section>
                        <q-input v-model="currentLink.name" label="Name" style="width: 320px"/>
                    </q-card-section>

                    <q-card-section>
                        <q-input v-model="currentLink.expired_at" label="Expired At*" style="width: 320px"
                            :rules="[val => !!val || 'Expiration date is required']">
                            <template v-slot:prepend>
                                <q-icon name="event" class="cursor-pointer">
                                    <q-popup-proxy cover transition-show="scale" transition-hide="scale">
                                        <q-date v-model="currentLink.expired_at" first-day-of-week="1" mask="YYYY-MM-DD">
                                            <div class="row items-center justify-end">
                                                <q-btn v-close-popup label="Close" color="primary" flat />
                                            </div>
                                        </q-date>
                                    </q-popup-proxy>
                                </q-icon>
                            </template>
                        </q-input>
                    </q-card-section>

                    <q-card-actions align="right">
                        <q-btn flat label="Cancel" color="primary" v-close-popup @click="clearLinkCreation" />
                        <q-btn label="Save" color="primary" type="submit" />
                    </q-card-actions>
                </q-form>
            </q-card>
        </q-dialog>

        <import-links-dialog v-model="showUploaderDialog" @close="uploadFinished"/>
    </div>
</template>

<script lang="ts">
import { exportFile, useQuasar, date, copyToClipboard } from 'quasar';
import Link from '../../models/Link';
import { defineComponent, ref } from 'vue';
import { LinkAPI, LinkAPIList } from '../../api/Link';
import CSV from '../../services/CSV';
import ImportLinksDialog from './dialogs/ImportLinksDialog.vue';

export default defineComponent({
    components: { ImportLinksDialog },
    name: 'LinksList',

    setup() {
        const $q = useQuasar();
        const links = ref<Link[]>();
        const confirmDeleteDialog = ref<boolean>(false);
        const confirmCreationDialog = ref<boolean>(false);
        const confirmDeleteLink = ref<boolean>(false);
        const loading = ref<boolean>(false);
        const showUploaderDialog = ref<boolean>(false);
        const currentLink = ref<Link>();
        const valid = ref<boolean>();
        const filter = ref('');
        const pagination = ref({
            sortBy: 'created_at',
            descending: true,
            rowsPerPage: 25,
            page: 1,
            rowsNumber: 50,
        });

        const headers = [
            {
                name: 'id',
                label: 'ID',
                field: 'id',
                align: 'left',
                style: 'width: 120px',
            },
            {
                name: 'name',
                label: 'Name',
                field: 'name',
                align: 'left',
                sortable: true,
                style: 'max-width: 200px',
            },
            {
                name: 'url',
                label: 'URL',
                field: 'url',
                align: 'left',
                sortable: true,
                style: 'max-width: 500px',
            },
            {
                name: 'expired_at',
                label: 'Expired at',
                field: 'expired_at',
                align: 'left',
                sortable: true,
                style: 'width: 100px',
            },
            {
                name: 'created_at',
                label: 'Created at',
                field: 'created_at',
                align: 'left',
                sortable: true,
                style: 'width: 100px',
            },
            {
                name: 'actions',
                label: 'Actions',
                align: 'left',
                style: 'width: 240px',
                required: true,
            },
        ];

        const formatDatetime = (datetime: string) => {
            if (datetime) {
                return datetime.substr(0, 10) + ' ' + datetime.substr(11, 5);
            }
            return '';
        };

        const clearLinkCreation = () => {
            currentLink.value = new Link(
                '',
                '',
                '',
                (date.addToDate(new Date(), {years: 50}).toISOString()).substr(0, 10),
                '',
            );
        };

        const deleteLink = () => {
            LinkAPI.delete(currentLink.value.id)
                .then(() => {
                    getList();
                    currentLink.value.id = '';

                    $q.notify({
                        type: 'positive',
                        message: 'Successfull link deletion',
                    });
                })
                .catch((error) => {
                    getList();
                    currentLink.value.id = '';

                    $q.notify({
                        type: 'negative',
                        message: `Error: ${error}`,
                    });
                    console.error(error);
                });
        };

        const newLink = () => {
            clearLinkCreation();
            confirmCreationDialog.value = true;
        };

        const editLink = () => {
            if (currentLink.value.id) {
                updateLink();
            } else {
                addLink();
            }
            
            confirmCreationDialog.value = false;
        };

        const addLink = () => {
            LinkAPI.add(currentLink.value)
                .then(() => {
                    getList();

                    $q.notify({
                        type: 'positive',
                        message: 'Successfull link creation',
                    });
                })
                .catch((error) => {
                    $q.notify({
                        type: 'negative',
                        message: `Error: ${error}`,
                    });
                    console.error(error);
                });
        };

        const updateLink = () => {
            LinkAPI.update(currentLink.value)
                .then(() => {
                    getList();

                    $q.notify({
                        type: 'positive',
                        message: 'Successfull link update',
                    });
                })
                .catch((error) => {
                    getList();

                    $q.notify({
                        type: 'negative',
                        message: `Error: ${error}`,
                    });
                    console.error(error);
                });
        };

        const openLink = () => {
            window.open(`${process.env.SORT_URL_BASE}/${currentLink.value.id}`, '_blank');
        };

        const getList = (props?) => {
            loading.value = true;

            const search = props ? props.filter : filter.value;
            const { page, rowsPerPage, sortBy, descending } = props ? props.pagination : pagination.value;

            LinkAPI.list(search, page, rowsPerPage, sortBy, descending)
                .then((linksList: LinkAPIList) => {
                    pagination.value.page = page;
                    pagination.value.rowsPerPage = rowsPerPage;
                    pagination.value.sortBy = sortBy;
                    pagination.value.descending = descending;
                    pagination.value.rowsNumber = linksList.total;

                    links.value = linksList.links;

                    loading.value = false;
                })
                .catch((error) => {
                    console.error(error);
                    loading.value = false;
                });
        };

        const exportCSV = () => {
            // Headers
            const csvHeaders = headers
                .filter((line) => line.name != 'actions')
                .map((line) => line.label);

            // Body
            const csvBody = links.value
                .map((link) => [
                    link.id,
                    link.name,
                    link.url,
                    link.expired_at,
                    link.created_at,
                ]);

            const CSVService = new CSV(';');
            CSVService.stringify(csvHeaders, csvBody)
                .then((content) => {
                    const status = exportFile(
                        `${date.formatDate(Date.now(), 'YYYYMMDDHHmmss')}_url-shortener.csv`,
                        content,
                        'text/csv',
                    );
                    if (status !== true) {
                        $q.notify({
                            color: 'negative',
                            icon: 'warning',
                            message: 'Browser denied file download',
                        });
                    }
                })
                .catch(() => {
                    $q.notify({
                        color: 'negative',
                        icon: 'warning',
                        message: 'Error during CSV creation',
                    });
                });
        };

        const copyLink = (id) => {
            copyToClipboard(`${process.env.SORT_URL_BASE}/${id}`)
                .then(() => {
                    $q.notify({
                        color: 'positive',
                        message: 'Link successfully copied',
                    });
                })
                .catch(() => {
                    $q.notify({
                        color: 'negative',
                        icon: 'warning',
                        message: 'Error when copying link to clipboard',
                    });
                });
        };

        const importCSV = () => {
            console.log('Import CSV file');
        }

        const uploadFinished = (reload) => {
            // Hide uploader dialog
            showUploaderDialog.value = false;

            // Reload links list
            if (reload) {
                getList();
            }
        }

        void getList();

        return {
            links,
            currentLink,
            headers,
            confirmDeleteDialog,
            confirmCreationDialog,
            confirmDeleteLink,
            showUploaderDialog,
            valid,
            pagination,
            filter,
            loading,
            formatDatetime,
            deleteLink,
            newLink,
            editLink,
            addLink,
            updateLink,
            openLink,
            getList,
            clearLinkCreation,
            exportCSV,
            importCSV,
            copyLink,
            uploadFinished,
        };
    },
});
</script>

<style scoped>
tr:nth-child(odd) {
    background-color: #93939314 !important;
}

.search_input {
    width: 320px;
}
</style>