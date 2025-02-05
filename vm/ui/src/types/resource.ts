export enum ResourceKind {
    Node = "node",
    StoragePool = "storagepool",
    StorageVolume = "storagevolume",
    Network = "network",
    Machine = "machine",
    MachineRequest = "machinerequest",
}

export enum Phase {
    Requsted = "requested",
    Creating = "creating",
    Errored = "errored",
    Created = "created",
}

export interface OwnerRef {
    kind: string;
    id: string;
}

export interface Resource<T> {
    id: string;
    kind: string;
    annotations: Record<string, string> | undefined;
    owner: OwnerRef | undefined;
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

export interface DeleteResourceInput {
    kind: ResourceKind;
    id: string;
}
