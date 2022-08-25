import { defineStore } from 'pinia';
import User from 'src/models/User';
import { Auth, AuthUser } from 'src/api/Auth';

interface State {
    user: User;
}

export const useUserStore = defineStore('user', {
    state: (): State => {
        return {
            user: User.fromSession(),
        }
    },

    getters: {
        getUser: (state: State) => {
            return state.user;
        },

        isAuthenticated: (state: State) => {
            return state.user.token !== ''
        }
    },

    actions: {
        init(userAuth: AuthUser): Promise<User> {
            return new Promise((resolve, reject) => {
                Auth.login(userAuth)
                    .then((user) => {
                        this.user = user;
                        if (user.isAuthenticated()) {
                            sessionStorage.setItem('user', JSON.stringify(user));
                        } else {
                            sessionStorage.removeItem('user');
                        }

                        resolve(User.fromUser(user));
                    })
                    .catch((error) => {
                        this.user = User.initEmpty();
                        sessionStorage.removeItem('user');

                        reject(error);
                    });
            });
        },

        logout() {
            this.user = User.initEmpty();

            sessionStorage.removeItem('user');
        },
    },
})
