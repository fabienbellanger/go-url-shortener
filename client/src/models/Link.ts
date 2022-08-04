/**
 * Link class
 *
 * @author Fabien Bellanger
 */
export default class Link {
    id = '';
    url = '';
    expired_at = '';

    constructor(
        id: string,
        url: string,
        expired_at: string
    ) {
        this.id = id;
        this.url = url;
        this.expired_at = (expired_at.length > 10) ? expired_at.substring(0, 10) : expired_at;
    }
}
