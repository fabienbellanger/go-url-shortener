<template>
    <div class="q-px-md">
        <h4 class="q-mt-lg">Users list</h4>

        <q-table
            :rows="users"
            :columns="headers"
            :filter="filter"
            :rows-per-page-options="[25, 50, 100]"
            row-key="id"
            no-data-label="No user"
            color="primary"
            binary-state-sort>
            <template v-slot:body="props">
                <q-tr :props="props">
                    <q-td key="id" :props="props">
                        {{ props.row.id }}
                    </q-td>
                    <q-td key="firstname" :props="props">
                        {{ props.row.firstname }}
                    </q-td>
                    <q-td key="lastname" :props="props">
                        {{ props.row.lastname }}
                    </q-td>
                    <q-td key="username" :props="props">
                        {{ props.row.username }}
                    </q-td>
                    <q-td key="created_at" :props="props">
                        {{ formatDatetime(props.row.created_at) }}
                    </q-td>
                    <q-td key="actions" :props="props">
                        <!-- <q-btn
                            size="sm"
                            icon="edit"
                            color="green"
                            @click="
                                currentUser = props.row;
                                confirmCreationDialog = true;
                            "></q-btn>
                        &nbsp; -->
                        <q-btn
                            size="sm"
                            icon="delete"
                            color="red"
                            @click="
                                currentUser = props.row;
                                confirmDeleteDialog = true;
                            "></q-btn>
                    </q-td>
                </q-tr>
            </template>
            <template v-slot:top-right>
                <q-input clearable dense debounce="300" v-model="filter" placeholder="Search">
                    <template v-slot:prepend>
                        <q-icon name="search" />
                    </template>
                </q-input>
            </template>
            <template v-slot:top-left>
                <q-btn round color="primary" icon="add" @click="
                    clearUserCreation();
                    confirmCreationDialog = true;
                " />
            </template>
        </q-table>

        <!-- Confirm User deletion dialog -->
        <q-dialog v-model="confirmDeleteDialog" persistent>
            <q-card>
                <q-card-section class="row items-center">
                    <q-icon name="warning" color="warning" class="text-h4" />
                    <span class="q-ml-sm text-h6">Do you really want to delete this user?</span>
                </q-card-section>

                <q-card-actions align="right">
                    <q-btn flat label="Cancel" color="primary" v-close-popup />
                    <q-btn flat label="Delete" color="primary" v-close-popup @click="deleteUser" />
                </q-card-actions>
            </q-card>
        </q-dialog>

        <!-- User creation/update dialog -->
        <q-dialog v-model="confirmCreationDialog" medium @hide="clearUserCreation">
            <q-card>
                <q-form @submit="editUser">
                    <q-card-section class="row items-center">
                        <span class="q-ml-sm text-h6">
                            <span v-if="currentUser.id">User update</span>
                            <span v-else>User creation</span>
                        </span>
                    </q-card-section>

                    <q-card-section>
                        <q-input v-model="currentUser.lastname" label="Lastname" style="width: 320px" autofocus
                            :rules="[val => !!val || 'Lastname is required']"/>
                    </q-card-section>

                    <q-card-section>
                        <q-input v-model="currentUser.firstname" label="Firstname" style="width: 320px" autofocus
                            :rules="[val => !!val || 'Firstname is required']"/>
                    </q-card-section>

                    <q-card-section>
                        <q-input v-model="currentUser.username" label="Username" style="width: 320px" autofocus type="email"
                            :rules="[val => !!val || 'Username is required', val => checkEmail(val) || 'Email is not valid']"/>
                    </q-card-section>

                    <q-card-section v-if="!currentUser.id">
                        <q-input v-model="currentUser.password" label="Password" style="width: 320px" autofocus type="password" autocomplete="new-password"
                            :rules="[val => !!val || 'Password is required', val => val.length >= 8 || 'Please use at least 8 characters']"/>
                    </q-card-section>

                    <q-card-actions align="right">
                        <q-btn flat label="Cancel" color="primary" v-close-popup @click="clearUserCreation" />
                        <q-btn label="Save" color="primary" type="submit"/>
                    </q-card-actions>
                </q-form>
            </q-card>
        </q-dialog>
    </div>
</template>

<script lang="ts">
import { useQuasar } from 'quasar';
import User from '../../models/User';
import { defineComponent, ref } from 'vue';
import { useUserStore } from '../../stores/user';
import { UserAPI } from '../../api/User';
import * as EmailValidator from 'email-validator';

export default defineComponent({
    name: 'UsersList',

    setup() {
        const $q = useQuasar();
        const userStore = useUserStore();
        const filter = ref('');
        const users = ref<User[]>();
        const currentUser = ref<User>();
        const confirmDeleteDialog = ref<boolean>(false);
        const confirmCreationDialog = ref<boolean>(false);
        const confirmDeleteUser = ref<boolean>(false);

        const headers = [
            {
                name: 'id',
                label: 'ID',
                field: 'id',
                align: 'left',
                style: 'width: 160px',
            },
            {
                name: 'firstname',
                label: 'Firstname',
                field: 'firstname',
                align: 'left',
                sortable: true,
            },
            {
                name: 'lastname',
                label: 'Lastname',
                field: 'lastname',
                align: 'left',
                sortable: true,
            },
            {
                name: 'username',
                label: 'Username',
                field: 'username',
                align: 'left',
                sortable: true,
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
                style: 'width: 80px',
                required: true,
            },
        ];

        const formatDatetime = (datetime: string) => {
            if (datetime) {
                return datetime.substr(0, 10) + ' ' + datetime.substr(11, 5);
            }
            return '';
        };

        const clearUserCreation = () => {
            currentUser.value = User.initEmpty();
        };

        const checkEmail = (email: string) => {
            return EmailValidator.validate(email);
        };

        const getList = () => {
            UserAPI.list()
                .then((usersList: User[]) => {
                    users.value = usersList;
                })
                .catch((error) => {
                    console.error(error);
                });
        };

        const deleteUser = () => {
            if (userStore.user.id === currentUser.value.id) {
                $q.notify({
                    type: 'negative',
                    message: 'Error: authenticated user cannot be deleted',
                });

                return;
            }

            UserAPI.delete(currentUser.value.id)
                .then(() => {
                    getList();
                    currentUser.value.id = '';

                    $q.notify({
                        type: 'positive',
                        message: 'Successfull user deletion',
                    });
                })
                .catch((error) => {
                    getList();
                    currentUser.value.id = '';

                    $q.notify({
                        type: 'negative',
                        message: `Error: ${error}`,
                    });
                    console.error(error);
                });
        };

        const editUser = () => {
            if (currentUser.value.id) {
                updateUser();
            } else {
                addUser();
            }

            confirmCreationDialog.value = false;
        };

        const addUser = () => {
            UserAPI.add(currentUser.value)
                .then(() => {
                    getList();

                    $q.notify({
                        type: 'positive',
                        message: 'Successfull user creation',
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

        const updateUser = () => {
            return;
        };

        getList();

        return {
            headers,
            filter,
            users,
            confirmCreationDialog,
            confirmDeleteDialog,
            confirmDeleteUser,
            currentUser,
            formatDatetime,
            getList,
            deleteUser,
            clearUserCreation,
            editUser,
            checkEmail,
        };
    },
});
</script>

<style scoped>
tr:nth-child(odd) {
    background-color: #93939314 !important;
}
</style>