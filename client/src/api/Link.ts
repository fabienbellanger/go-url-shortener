import { exportFile } from 'quasar';
import Link from 'src/models/Link';
import Http from '../services/Http';

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
     * @return {Promise<Link[]>}
     */
    public static list(): Promise<Link[]> {
        return new Promise((resolve, reject) => {
            Http.request('GET', 'links', true)
                .then((links: Link[]) => {
                    for (const i in links) {
                        links[i] = new Link(links[i].id, links[i].url, links[i].expired_at);
                    }
                    resolve(links);
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
            if (link.url !== '' && link.expired_at !== '')
            {
                Http.request('POST', '/links', true, {
                    url: link.url,
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
