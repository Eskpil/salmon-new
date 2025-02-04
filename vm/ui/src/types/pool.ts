export interface Pool {
    name: string;
    id: string;
    allocated_volumes: number;
    capacity: number;
    available: number;
    allocation: number;
    type: string;
}
