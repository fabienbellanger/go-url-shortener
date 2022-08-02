import User from 'src/models/User';

export interface UserStateInterface {
    user: User;
}

function state(): UserStateInterface {
    return {
        user: User.fromSession(),
    }
}

export default state;
