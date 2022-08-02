import Http from '../services/Http';
import SaleError from '../models/SaleError';

/**
 * Classe gérant les appels API des erreurs de notes
 *
 * @author Fabien Bellanger
 */
class SalesErrorAPI {
    /**
     * Liste des erreurs de note
     *
     * @author Fabien Bellanger
     * @return {Promise<SaleError[]>}
     */
    public static list(projectId?: string, dateBegin?: string, dateEnd?: string): Promise<SaleError[]> {
        return new Promise((resolve, reject) => {
            let url = 'sales-errors';

            // TODO: Properly construct URL with query parameters
            if (projectId !== undefined && projectId !== '') {
                url += `?project_id=${projectId}`;
            }
            if (dateBegin !== undefined && dateBegin !== '') {
                url += `?date_begin=${dateBegin}`;
            }
            if (dateEnd !== undefined && dateEnd !== '') {
                url += `?date_end=${dateEnd}`;
            }

            Http.request('GET', url)
                .then((salesErrors: SaleError[]) => {
                    resolve(salesErrors);
                })
                .catch((error) => {
                    reject(error);
                });
        });
    }
}

export { SalesErrorAPI };
