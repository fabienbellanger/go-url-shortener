import { GetterTree } from 'vuex';
import { StateInterface } from '../index';
import { UserStateInterface } from './state';

const getters: GetterTree<UserStateInterface, StateInterface> = {
    getState(state) {
        return state;
    },
    getUser(state) {
        return state.user;
    },
};

export default getters;
