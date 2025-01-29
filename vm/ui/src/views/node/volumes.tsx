import { useQuery } from "@tanstack/react-query";
import { getVolumes } from "../../data/queries/volumes";
import { Badge, Box, Table } from "@radix-ui/themes";
import { useNavigate } from "react-router";

interface Props {
    id: string;
}

export const VolumesView: React.FC<Props> = ({ id }) => {};
