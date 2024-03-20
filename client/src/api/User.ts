import User from 'src/models/User';
import { http } from 'src/boot/http';
import * as EmailValidator from 'email-validator';

interface AuthUser {
    email: string,
    password: string,
}

/**
 * Classe gérant les appels API des utilisateurs
 *
 * @author Fabien Bellanger
 */
class UserAPI {
    /**
     * Authentification et récupération des informations de l'utilisateur
     *
     * @author Fabien Bellanger
     * @param {AuthUser} authUser
     * @return {Promise<User>}
     */
     public static login(authUser: AuthUser): Promise<User> {
        return new Promise((resolve, reject) => {
            http.request('POST', 'login', false, {}, { username: authUser.email, password: authUser.password })
                .then((user: User) => {
                    resolve(User.fromUser(user));
                })
                .catch((error) => {
                    reject(error);
                });
        });
    }

    /**
     * Liste des utilisateurs
     *
     * @author Fabien Bellanger
     * @return {Promise<User[]>}
     */
    public static list(): Promise<User[]> {
        return new Promise((resolve, reject) => {
            http.request('GET', 'users', true)
                .then((data: User[]) => {
                    const users = [];

                    for (const i in data) {
                        users.push(new User(
                            data[i].id,
                            data[i].lastname,
                            data[i].firstname,
                            data[i].username,
                            data[i].created_at,
                            data[i].updated_at,
                            ''
                        ));
                    }

                    resolve(users);
                })
                .catch((error: Error) => {
                    reject(error);
                })
        });
    }

    /**
     * Suppression d'un utilisateur
     *
     * @author Fabien Bellanger
     * @param id string ID du lien
     * @return {Promise<user[]>}
     */
     public static delete(id: string): Promise<User> {
        return new Promise((resolve, reject) => {
            if (id !== '') {
                http.request('DELETE', `/users/${id}`)
                    .then((user: User) => {
                        resolve(user);
                    })
                    .catch((error) => {
                        reject(error);
                    });
            } else {
                reject(new Error('invalid id'));
            }
        });
    }

    /**
     * Ajout d'un utilisateur
     *
     * @author Fabien Bellanger
     * @param user User Utilisateur
     * @return {Promise<user>}
     */
     public static add(user: User): Promise<User> {
        return new Promise((resolve, reject) => {
            if (user.lastname !== '' && user.firstname !== '' && user.username !== '') {
                http.request('POST', '/register', true, {}, {
                    lastname: user.lastname,
                    firstname: user.firstname,
                    username: user.username,
                    password: user.password,
                })
                    .then((user: User) => {
                        resolve(user);
                    })
                    .catch((error) => {
                        reject(error);
                    });
            } else {
                reject(new Error('invalid user properties'));
            }
        });
    }

    /**
     * Mot de passe oublié
     *
     * @author Fabien Bellanger
     * @param email string Mail de l'utilisateur
     * @return {Promise<void>}
     */
    public static forgottenPassword(email: string): Promise<void> {
        return new Promise((resolve, reject) => {
            if (EmailValidator.validate(email)) {
                http.request('POST', `/forgotten-password/${email}`)
                    .then(() => {
                        resolve();
                    })
                    .catch((error) => {
                        reject(error);
                    });
            } else {
                reject(new Error('invalid email'))
            }
        });
    }

    /**
     * Changement du mot de passe
     *
     * @author Fabien Bellanger
     * @param token string Token
     * @param password string Mot de passe de l'utilisateur
     * @return {Promise<ForgottenPassword>}
     */
    public static updatePassword(token: string, password: string): Promise<void> {
        return new Promise((resolve, reject) => {
            if (token.length === 36 && password.length >= 8) {
                http.request('PATCH', `/update-password/${token}`, false, {}, {
                    password,
                })
                    .then(() => {
                        resolve();
                    })
                    .catch((error) => {
                        reject(error);
                    });
            } else {
                reject(new Error('invalid token or password'))
            }
        });
    }
}

export { UserAPI };
export type { AuthUser };
