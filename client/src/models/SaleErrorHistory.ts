/**
 * SaleErrorHistory class
 *
 * @author Fabien Bellanger
 */
 export default class SaleErrorHistory {
    id = '';
    sale_error_id = '';
    gap = 0;
    created_at = '';
    updated_at = '';

    constructor(
        id: string,
        sale_error_id: string,
        gap: number,
        created_at: string,
        updated_at: string
    ) {
        this.id = id;
        this.sale_error_id = sale_error_id;
        this.gap = gap;
        this.created_at = created_at;
        this.updated_at = updated_at;
    }
}
