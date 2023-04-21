/* eslint-disable @typescript-eslint/no-explicit-any */
/* eslint-disable @typescript-eslint/no-unsafe-assignment */
import axios, { AxiosError, AxiosInstance, AxiosRequestConfig, Method } from 'axios';
import User from 'src/models/User';

/**
 * Service gérant les requêtes HTTP
 *
 * @author Fabien Bellanger
 */
class Http
{
    private static instance: Http;
    public baseURL = '';
    public headers?: any;
    public loading = '';
    public method: Method = 'GET';
    public parameters?: any;
    public url = '';
    public withCredentials = false;

    /**
     * Constructeur
     *
     * @author Fabien Bellanger
     */
    constructor()
    {
        if (!Http.instance)
        {
            // Initialisation
            // --------------
            this.baseURL = process.env.API_BASE_URL as string;
        }

        return Http.instance;
    }

    // /**
    //  * Gestion des erreurs
    //  *
    //  * @static
    //  * @author Fabien Bellanger
    //  * @param {AxiosResponse<any>} response
    //  */
    // public static manageError(response: AxiosResponse<any>): void
    // {
    //     if (response.status === 401)
    //     {
    //         // Not autorized : On déconnecte l'utilisateur
    //         // -------------------------------------------
    //         // tslint:disable-next-line
    //         app.$router.push('/logout').catch(() => {});
    //     }
    //     else if (response.status === 200)
    //     {
    //         const data = response.data;

    //         if (typeof data.code === 'number' && data.code === 401)
    //         {
    //             // Not autorized : On déconnecte l'utilisateur
    //             // -------------------------------------------
    //             // tslint:disable-next-line
    //             app.$router.push('/logout').catch(() => {});
    //         }
    //     }
    // }

    /**
     * Requête
     *
     * @author Fabien Bellanger
     * @param {string}      method                  Méthode {'get', 'post', 'put', 'patch', 'delete'}
     * @param {string}      url                     URL
     * @param {boolean}     withCredentials=true    Paramètres à transmettre (pour les requêtes POST, PUT et PATCH)
     * @param {object}      parameters={}           Paramètres type query parameters
     * @param {object}      data={}                 Body
     * @param {object}      headers={}              Headers
     * @param {string}      baseUrl=''              Base URL autre que celle de la caisse
     * @return {Promise<any>}
     */
     public request(
        method: Method,
        url: string,
        withCredentials = true,
        parameters = {},
        data = {},
        headers = {},
        baseUrl = ''): Promise<any>
    {
        return new Promise((resolve, reject) =>
        {
            if (['GET', 'POST', 'PUT', 'PATCH', 'DELETE'].indexOf(method) === -1)
            {
                reject(new Error('Bad method'));
            }
            else if (typeof url !== 'string')
            {
                reject(new Error('Bad request'));
            }
            else
            {
                const config: AxiosRequestConfig = {
                    url,
                    method,
                    headers,
                    withCredentials,
                    baseURL:         baseUrl === '' ? this.baseURL : baseUrl,
                    params:          parameters,
                    responseType:    'json',
                    data:            data,
                };
                const axiosInstance: AxiosInstance = axios.create(config);

                // Intercepteur - Requête
                // ----------------------
                axiosInstance.interceptors.request.use((requestConfig) =>
                    {
                        // Gestion des token d'authentification
                        // ------------------------------------
                        // S'il faut un token et qu'un token est bien présent alors on l'ajoute aux headers
                        if (withCredentials)
                        {
                            const user = User.fromSession();

                            if (user.token !== null && user.token !== '')
                            {
                                requestConfig.headers = { Authorization: `Bearer ${user.token}` };
                            }
                        }

                        return requestConfig;
                    },
                    (error: AxiosError) =>
                    {
                        return Promise.reject(error);
                    },
                );

                // Intercepteur - Réponse
                // ----------------------
                axiosInstance.interceptors.response.use((response) =>
                    {
                        // Gestion des erreurs
                        // Http.manageError(response);

                        return response;
                    },
                    (error: AxiosError) =>
                    {
                        // Http.manageError(error.response);
                        const message = error.response.data.message ?? 'Unknown error';
                        
                        return Promise.reject(message);
                    },
                );

                // Lancement de la requête
                // -----------------------
                axiosInstance.request(config)
                    .then((response) =>
                    {
                        resolve(response.data);
                    })
                    .catch((error) =>
                    {
                        reject(error);
                    });
            }
        });
    }
}

const instance: Http = new Http();
Object.freeze(instance);
export default instance;
