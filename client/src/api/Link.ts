import { http } from 'src/boot/http';
import Link from 'src/models/Link';
import User from 'src/models/User';


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

            http.request('GET', url.pathname + url.search, true)
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
                http.request('POST', '/links', true, {}, {
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
                http.request('PUT', `/links/${link.id}`, true, {}, {
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
     * @return {Promise<Link>}
     */
    public static delete(id: string): Promise<Link> {
        return new Promise((resolve, reject) => {
            if (id !== '')
            {
                http.request('DELETE', `/links/${id}`)
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

    /**
     * Suppression d'un ensemble de liens
     *
     * @author Fabien Bellanger
     * @param ids string[] ID des liens
     * @return {Promise<Link[]>}
     */
    public static deleteSelectedLinks(ids: string[]): Promise<void> {
        return new Promise((resolve, reject) => {
            if (ids.length !== 0)
            {
                http.request('DELETE', '/links/selected', true, {}, ids)
                    .then(() => {
                        resolve();
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

    /**
     * Upload d'un fichier CSV
     *
     * @author Fabien Bellanger
     * @return {Promise<any>}
     */
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    public static upload(): Promise<any> {
        return new Promise((resolve, reject) => {
            const user = User.fromSession();
            if (user.token !== null && user.token !== '') {
                resolve({
                    url: `${http.baseURL}/links/upload`,
                    method: 'POST',
                    headers: [
                        { name: 'Authorization', value: `Bearer ${user.token}` }
                    ]
                })
            } else {
                reject('Empty token');
            }
        });
    }

    /**
     * Export des liens
     *
     * @author Fabien Bellanger
     * @return {Promise<string>}
     */
    public static export(filter: string): Promise<string> {
        return new Promise((resolve, reject) => {
            // Construct URL
            const url = new URL('https://www.apitic.com/links/export/csv');
            if (filter) {
                url.searchParams.append('s', filter);
            }

            http.request('GET', url.pathname + url.search, true)
                .then((data) => resolve(data))
                .catch((error) => {
                    reject(error);
                });
        });
    }
}

export { LinkAPI };
export type { LinkAPIList }
