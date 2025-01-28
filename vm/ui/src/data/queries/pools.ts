import { Pool } from "../../types/pool";
import { Resource, ResourceKind } from "../../types/resource";
import { List } from "../index";

export const getPools = async (
    nodeId: string,
): Promise<List<Resource<Pool>>> => {
    return fetch(
        `http://10.100.102:8080/v1/resources?owner_id=${nodeId}&owner_kind=${ResourceKind.Node}&kind=${ResourceKind.StoragePool}`,
    ).then((res) => res.json());
};
