import Link from 'src/models/Link';
import Http from '../services/Http';

type LinkAPIList = {
    total: number;
    links: Link[];
};

/**
 * Classe gérant les appels API des liens
 *
 * @author Fabien Bellanger
 */
class LinkAPI {
    /**
     * Liste des liens
     *
     * @author Fabien Bellanger
     * @return {Promise<LinkAPIList>}
     */
    public static list(filter: string, page: number, rowsPerPage: number, sortBy: string, descending: boolean): Promise<LinkAPIList> {
        return new Promise((resolve, reject) => {
            // Construct URL
            const url = new URL('https://www.apitic.com/links');
            url.searchParams.append('page', page.toString());
            url.searchParams.append('limit', rowsPerPage.toString());
            url.searchParams.append('sort', descending ? 'desc' : 'asc');
            if (filter) {
                url.searchParams.append('s', filter);
            }
            if (sortBy) {
                url.searchParams.append('sort-by', sortBy);
            }

            Http.request('GET', url.pathname + url.search, true)
                .then((data) => {
                    const links = data.links;
                    for (const i in links) {
                        links[i] = new Link(links[i].id, links[i].url, links[i].name, links[i].expired_at, links[i].created_at);
                    }
                    resolve({
                        total: data.total,
                        links: links,
                    });
                })
                .catch((error) => {
                    reject(error);
                });
        });
    }

    /**
     * Ajout d'un lien
     *
     * @author Fabien Bellanger
     * @param link Link Lien
     * @return {Promise<Link[]>}
     */
     public static add(link: Link): Promise<Link> {
        return new Promise((resolve, reject) => {
            if (link.url !== '' && link.expired_at !== '') {
                Http.request('POST', '/links', true, {
                    url: link.url,
                    name: link.name,
                    expired_at: (new Date(link.expired_at)).toISOString()
                })
                    .then((link: Link) => {
                        resolve(link);
                    })
                    .catch((error) => {
                        reject(error);
                    });
            } else {
                reject(new Error('invalid URL or expired date'));
            }
        });
    }

    /**
     * Modification d'un lien
     *
     * @author Fabien Bellanger
     * @param link Link Lien
     * @return {Promise<Link[]>}
     */
     public static update(link: Link): Promise<Link> {
        return new Promise((resolve, reject) => {
            if (link.url !== '' && link.expired_at !== '')
            {
                Http.request('PUT', `/links/${link.id}`, true, {
                    url: link.url,
                    name: link.name,
                    expired_at: (new Date(link.expired_at)).toISOString()
                })
                    .then((link: Link) => {
                        resolve(link);
                    })
                    .catch((error) => {
                        reject(error);
                    });
            }
            else
            {
                reject(new Error('invalid URL or expired date'));
            }
        });
    }

    /**
     * Suppression d'un lien
     *
     * @author Fabien Bellanger
     * @param id string ID du lien
     * @return {Promise<Link[]>}
     */
    public static delete(id: string): Promise<Link> {
        return new Promise((resolve, reject) => {
            if (id !== '')
            {
                Http.request('DELETE', `/links/${id}`)
                    .then((link: Link) => {
                        resolve(link);
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

export { LinkAPI };
export type { LinkAPIList }
