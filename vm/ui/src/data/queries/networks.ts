import { List } from "..";
import { Network } from "../../types/network";
import { Resource, ResourceKind } from "../../types/resource";

export const getNetworks = async (
    nodeId: string,
): Promise<List<Resource<Network>>> => {
    return fetch(
        `http://10.100.102:8080/v1/resources?owner_id=${nodeId}&owner_kind=${ResourceKind.Node}&kind=${ResourceKind.Network}`,
    ).then((res) => res.json());
};
