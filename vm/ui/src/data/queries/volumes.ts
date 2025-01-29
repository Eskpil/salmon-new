import { Volume } from "../../types/volume";
import { Resource, ResourceKind } from "../../types/resource";
import { List } from "../index";

export const getVolumes = async (
    poolId: string,
): Promise<List<Resource<Volume>>> => {
    return fetch(
        `http://10.100.102:8080/v1/resources?owner_id=${poolId}&owner_kind=${ResourceKind.StoragePool}&kind=${ResourceKind.StorageVolume}`,
    ).then((res) => res.json());
};
