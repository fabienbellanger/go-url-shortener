/* eslint-disable @typescript-eslint/no-explicit-any */
/* eslint-disable @typescript-eslint/no-unsafe-assignment */
import { boot } from 'quasar/wrappers';
import axios, { AxiosError, AxiosInstance, AxiosRequestConfig, Method, AxiosResponse } from 'axios';
import User from 'src/models/User';
import { Router } from 'vue-router';

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
    private router: Router | null = null;

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

    /**
     * Ajout du router
     *
     * @author Fabien Bellanger
     * @param {Router} router
     */
    public addRouter(router: Router): void
    {
        this.router = router;
    }

    /**
     * Gestion de l'erreur 401
     *
     * @author Fabien Bellanger
     * @param {AxiosResponse<any>} response
     */
    public manage401Error(response: AxiosResponse<any>): void
    {
        if (response.status === 401)
        {
            this.router?.push({name: 'logout'});
        }
        else if (response.status === 200)
        {
            const data = response.data;

            if (typeof data.code === 'number' && data.code === 401)
            {
                this.router?.push({name: 'logout'});
            }
        }
    }

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
                                requestConfig.headers.Authorization = `Bearer ${user.token}`;
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
                        this.manage401Error(response);

                        return response;
                    },
                    (error: AxiosError) =>
                    {
                        const responseError = error.response;
                        if (responseError !== undefined) {
                            this.manage401Error(responseError);
                        }
                        
                        return Promise.reject(error);
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

export default boot(({ router }) => {
    http.addRouter(router);
});

const http: Http = new Http();
export { http };
