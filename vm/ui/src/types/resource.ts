export enum ResourceKind {
    Node = "node",
    StoragePool = "storagepool",
    StorageVolume = "storagevolume",
    Network = "network",
    Machine = "machine",
}

export enum Phase {
    Requsted = "requested",
    Creating = "creating",
    Created = "created",
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
    status: Status;
}

export interface Status {
    phase: Phase;
    error: string;
}

// TODO: avoid using any
export interface CreateResourceInput {
    kind: string;
    annotations: any;
    owner_ref: OwnerRef | undefined;
    spec: any;
}
