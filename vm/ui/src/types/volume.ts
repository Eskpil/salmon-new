export interface Volume {
    name: string;
    pool: string;
    format: string;
    schema: string;

    key: string;
    type: number;
    allocation: number;
    capacity: number;
}
