import { stringify } from 'csv-stringify/browser/esm';
import { arrayBuffer } from 'stream/consumers';

/**
 * Service gérant l'export de fichier CSV
 *
 * @author Fabien Bellanger
 */
class CSV
{
    public separator: string;

    /**
     * Constructeur
     *
     * @author Fabien Bellanger
     * @param string separator Séparateur de colonnes
     */
    constructor(separator = ';')
    {
        this.separator = separator;
    }

    /** 
     * Contruit le CSV sous forme de chaîne de caractères
     *
     * @author Fabien Bellanger
     * @param string[] headers Entête
     * @param string[][] body Contenu
     * @return Promise<string>
     */
    stringify(headers: string[], body: string[][]): Promise<string> {
        return new Promise((resolve, reject) => {
            body.unshift(headers);
            
            stringify(
                body,
                {
                    delimiter: this.separator,
                    quoted_string: true,
                }, (err, output) => {
                    if (err) {
                        reject(err);
                    } else {
                        resolve(output);
                    }
                });
        });
    }
}

export default CSV;
