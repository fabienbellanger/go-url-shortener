export type UrlParameters = {
    name: string,
    value: string,
}

/**
 * Service générant une URL
 *
 * @author Fabien Bellanger
 */
export class Url
{
    /**
     * Generate URL
     * 
     * @param base string Base URL
     * @param parameters UrlParameters[] Query parameters
     * @return string Encoded URL
     */
    public static generate(base: string, parameters?: UrlParameters[]): string {
        let url = base;

        if (parameters !== undefined && parameters.length > 0) {
            let first = true;
            
            for (const parameter of parameters) {
                if (parameter.name !== '') {
                    if (first) {
                        url += '?';
                        first = false;
                    } else {
                        url += '&';
                    }
    
                    url += parameter.name + '=' + parameter.value;
                }
            }
        }

        return encodeURI(url);
    }
}
