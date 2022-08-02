import { Module } from 'vuex';
import { StateInterface } from '../index';
import state, { UserStateInterface } from './state';
import actions from './actions';
import getters from './getters';
import mutations from './mutations';

const userModule: Module<UserStateInterface, StateInterface> = {
    namespaced: false,
    actions,
    getters,
    mutations,
    state,
};

export default userModule;
