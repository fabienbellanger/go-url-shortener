import User from 'src/models/User';
import { Auth, AuthUser } from 'src/api/Auth';
import { ActionTree } from 'vuex';
import { StateInterface } from '../index';
import { UserStateInterface } from './state';
import * as userType from './type';

const actions: ActionTree<UserStateInterface, StateInterface> = {
    [userType.A_INIT](store, userAuth: AuthUser): Promise<User> {
        return new Promise((resolve, reject) => {
            Auth.login(userAuth)
                .then((user) => {
                    store.commit(userType.M_INIT, user);
                    resolve(User.fromUser(user));
                })
                .catch((error) => {
                    store.commit(userType.M_INIT, User.initEmpty());
                    reject(error);
                });
        });
    },

    [userType.A_LOGOUT](store) {
        store.commit(userType.M_LOGOUT);
    },
};

export default actions;
