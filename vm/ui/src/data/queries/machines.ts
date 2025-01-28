import { Machine } from "../../types/machine";
import { Resource, ResourceKind } from "../../types/resource";
import { List } from "../index";

export const getMachines = async (
    nodeId: string,
): Promise<List<Resource<Machine>>> => {
    return fetch(
        `http://10.100.102:8080/v1/resources?owner_id=${nodeId}&owner_kind=${ResourceKind.Node}&kind=${ResourceKind.Machine}`,
    ).then((res) => res.json());
};
