import User from 'src/models/User';
import { MutationTree } from 'vuex';
import { UserStateInterface } from './state';
import * as userType from './type';

const mutation: MutationTree<UserStateInterface> = {
    [userType.M_INIT](state, user: User) {
        state.user = user;

        if (user.isAuthenticated()) {
            sessionStorage.setItem('user', JSON.stringify(user));
        } else {
            sessionStorage.removeItem('user');
        }
    },

    [userType.M_LOGOUT](state) {
        state.user = User.initEmpty();

        sessionStorage.removeItem('user');
    },
};

export default mutation;
