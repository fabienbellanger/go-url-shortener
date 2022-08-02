import Project from 'src/models/Project';
import Http from '../services/Http';

/**
 * Classe gérant les appels API des projets
 *
 * @author Fabien Bellanger
 */
class ProjectAPI {
    /**
     * Liste des projets
     *
     * @author Fabien Bellanger
     * @param deleted boolean Affichage des projets supprimés
     * @return {Promise<Project[]>}
     */
    public static list(deleted: boolean): Promise<Project[]> {
        return new Promise((resolve, reject) => {
            const url = deleted ? 'projects/deleted' : 'projects';

            Http.request('GET', url)
                .then((projects: Project[]) => {
                    resolve(projects);
                })
                .catch((error) => {
                    reject(error);
                });
        });
    }
    /**
     * Liste de tous les projets
     *
     * @author Fabien Bellanger
     * @return {Promise<Project[]>}
     */
    public static listAll(): Promise<Project[]> {
        return new Promise((resolve, reject) => {
            Http.request('GET', 'projects/all')
                .then((projects: Project[]) => {
                    resolve(projects);
                })
                .catch((error) => {
                    reject(error);
                });
        });
    }

    /**
     * Mise à jour d'un projet
     *
     * @author Fabien Bellanger
     * @param id string ID du projet
     * @param name string Nouveau nom du projet
     * @return {Promise<Project[]>}
     */
    public static update(id: string, name: string): Promise<Project> {
        return new Promise((resolve, reject) => {
            if (id !== '' && name !== '')
            {
                Http.request('PUT', `/projects/${id}`, true, { name })
                    .then((project: Project) => {
                        resolve(project);
                    })
                    .catch((error) => {
                        reject(error);
                    });
            }
            else
            {
                reject(new Error('invalid id or name'));
            }
        });
    }

    /**
     * Suppression (logique) d'un projet
     *
     * @author Fabien Bellanger
     * @param id string ID du projet
     * @return {Promise<Project[]>}
     */
    public static delete(id: string): Promise<Project> {
        return new Promise((resolve, reject) => {
            if (id !== '')
            {
                Http.request('DELETE', `/projects/${id}`)
                    .then((project: Project) => {
                        resolve(project);
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
     * Réactive un projet
     *
     * @author Fabien Bellanger
     * @param id string ID du projet
     * @return {Promise<Project[]>}
     */
    public static recover(id: string): Promise<Project> {
        return new Promise((resolve, reject) => {
            if (id !== '')
            {
                Http.request('PATCH', `/projects/${id}/recover`)
                    .then((project: Project) => {
                        resolve(project);
                    })
                    .catch((error) => {
                        reject(error as Error);
                    });
            }
            else
            {
                reject(new Error('invalid id'));
            }
        });
    }

    /**
     * Ajout d'un projet
     *
     * @author Fabien Bellanger
     * @param name string Nom du projet
     * @return {Promise<Project[]>}
     */
    public static add(name: string): Promise<Project> {
        return new Promise((resolve, reject) => {
            if (name !== '')
            {
                Http.request('POST', '/projects', true, {name})
                    .then((project: Project) => {
                        resolve(project);
                    })
                    .catch((error) => {
                        reject(error);
                    });
            }
            else
            {
                reject(new Error('invalid name'));
            }
        });
    }
}

export { ProjectAPI };
