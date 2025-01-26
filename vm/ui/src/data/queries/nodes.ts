import { Node } from "../../types/node";
import { List } from "../index";

export const getNodes = async (): Promise<List<Node>> => {
    return fetch("http://10.100.102:8080/v1/nodes").then((res) => res.json());
};
