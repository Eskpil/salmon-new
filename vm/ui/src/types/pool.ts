export interface Pool {
    name: string;
    id: string;
    allocated_volumes: number;
    capacity: number;
    allocated: number;
    kind: string;
}
