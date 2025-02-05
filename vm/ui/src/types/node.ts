import { Topology } from "./topology";

export interface Node {
    hostname: string;
    kernel: string;
    active_machines: number;
    total_machines: number;
    topology: Topology;
    up_since: string;
}
