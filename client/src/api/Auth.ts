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
class Auth {
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
}

export { Auth };
export type { AuthUser };
