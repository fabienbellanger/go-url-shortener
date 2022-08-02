/**
 * Link class
 *
 * @author Fabien Bellanger
 */
export default class Link {
    id = '';
    link = '';
    expired_at = '';

    constructor(
        id: string,
        link: string,
        expired_at: string
    ) {
        this.id = id;
        this.link = link;
        this.expired_at = expired_at;
    }
}
