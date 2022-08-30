/**
 * Link class
 *
 * @author Fabien Bellanger
 */
export default class Link {
    id = '';
    url = '';
    name = '';
    expired_at = '';

    constructor(
        id: string,
        url: string,
        name: string,
        expired_at: string
    ) {
        this.id = id;
        this.url = url;
        this.name = name;
        this.expired_at = (expired_at.length > 10) ? expired_at.substring(0, 10) : expired_at;
    }

    isExpired(): boolean {
        return false;
    }
}
