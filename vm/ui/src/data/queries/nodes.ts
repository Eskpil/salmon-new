import { Node } from "../../types/node";
import { Resource } from "../../types/resource";
import { ResourceKind } from "../../types/resource";
import { List } from "../index";

export const getNodes = async (): Promise<List<Resource<Node>>> => {
    return fetch(
        `http://10.100.102:8080/v1/resources?kind=${ResourceKind.Node}`,
    ).then((res) => res.json());
};

export const getNode = async (id: string): Promise<Resource<Node>> => {
    const nodes: List<Resource<Node>> = await fetch(
        `http://10.100.102:8080/v1/resources?kind=${ResourceKind.Node}&id=${id}`,
    ).then((res) => res.json());

    return nodes.list[0];
};
