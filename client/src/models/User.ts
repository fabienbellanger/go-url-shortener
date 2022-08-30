/**
 * User class
 *
 * @author Fabien Bellanger
 */
export default class User {
    id = '';
    lastname = '';
    firstname = '';
    username = '';
    password = '';
    created_at = '';
    updated_at = '';
    token = '';

    constructor(
        id: string,
        lastname: string,
        firstname: string,
        username: string,
        created_at: string,
        updated_at: string,
        token: string,
    ) {
        this.id = id;
        this.lastname = lastname;
        this.firstname = firstname;
        this.username = username;
        this.created_at = created_at;
        this.updated_at = updated_at;
        this.token = token;
    }

    static initEmpty(): User {
        return new User('', '', '', '', '', '', '');
    }

    static fromUser(user: User): User {
        if (user === null) {
            return User.initEmpty();
        }
        return new User(
            user.id,
            user.lastname,
            user.firstname,
            user.username,
            user.created_at,
            user.updated_at,
            user.token
        );
    }

    static toSession() {
        sessionStorage.setItem('user', JSON.stringify(this));
    }

    static fromSession(): User {
        const user = <User> JSON.parse(sessionStorage.getItem('user') as string);
        return User.fromUser(user);
    }

    /**
     * Check if user is authenticated
     *
     * @author Fabien Bellanger
     * @return boolean
     */
    isAuthenticated(): boolean {
        return this.token !== '' && this.token !== null;
    }
}
