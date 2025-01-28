export enum ResourceKind {
    Node = "node",
    StoragePool = "storagepool",
    StorageVolume = "storagevolume",
    Network = "network",
    Machine = "machine",
}

export interface OwnerRef {
    kind: string;
    id: string;
}

export interface Resource<T> {
    id: string;
    kind: string;
    annotations: Map<string, string> | undefined;
    owner_ref: OwnerRef | undefined;
    spec: T | undefined;
}
