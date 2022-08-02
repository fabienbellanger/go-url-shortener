import SaleErrorHistory from './SaleErrorHistory';

/**
 * SaleError class
 *
 * @author Fabien Bellanger
 */
 export default class SaleError {
    id = '';
    project_id = '';
    sale_id = '';
    sale_ticket_number = '';
    sale_datetime = '';
    platform = '';
    gap = 0;
    history = [] as SaleErrorHistory[];
    created_at = '';
    updated_at = '';

    constructor(
        id: string,
        project_id: string,
        sale_id: string,
        sale_ticket_number: string,
        sale_datetime: string,
        platform: string,
        gap: number,
        history: SaleErrorHistory[],
        created_at: string,
        updated_at: string
    ) {
        this.id = id;
        this.project_id = project_id;
        this.sale_id = sale_id;
        this.sale_ticket_number = sale_ticket_number;
        this.sale_datetime = sale_datetime;
        this.platform = platform;
        this.gap = gap;
        this.history = history;
        this.created_at = created_at;
        this.updated_at = updated_at;
    }
}
