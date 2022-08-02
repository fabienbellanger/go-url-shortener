/**
 * Project class
 *
 * @author Fabien Bellanger
 */
export default class Project {
    id = '';
    name = '';
    created_at = '';
    updated_at = '';
    deleted_at = '';

    constructor(
        id: string,
        name: string,
        created_at: string,
        updated_at: string,
        deleted_at: string
    ) {
        this.id = id;
        this.name = name;
        this.created_at = created_at;
        this.updated_at = updated_at;
        this.deleted_at = deleted_at;
    }

    public static getNameFromId(projects: Project[] | undefined, id: string): string {
        if (projects !== undefined) {
            for (const project of projects) {
                if (project.id === id) {
                    return project.name;
                }
            }
        }
        return '';
    }
}
