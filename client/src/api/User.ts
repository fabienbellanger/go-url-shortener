import User from 'src/models/User';
import Http from '../services/Http';

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
            Http.request('POST', 'login', false, { username: authUser.email, password: authUser.password })
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
            Http.request('GET', 'users', true)
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
            if (id !== '')
            {
                Http.request('DELETE', `/users/${id}`)
                    .then((user: User) => {
                        resolve(user);
                    })
                    .catch((error) => {
                        reject(error);
                    });
            }
            else
            {
                reject(new Error('invalid id'));
            }
        });
    }
}

export { UserAPI };
export type { AuthUser };
