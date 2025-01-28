import { Topology } from "./topology";

export interface MachineDisk {
    device: string;
    type: string;
    volume: string;
}

export interface MachineInterface {
    mac: string;
    model: string;
    network: string;
}

export interface Machine {
    name: string;
    uuid: string;
    schema: string;

    topology: Topology;

    disks: MachineDisk[];
    interfaces: MachineInterface[];
}
