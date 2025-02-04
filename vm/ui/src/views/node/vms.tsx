import { useMutation, useQuery } from "@tanstack/react-query";
import { Badge, Box, Button, IconButton, Table } from "@radix-ui/themes";
import { getMachines } from "../../data/queries/machines";
import { useNavigate } from "react-router";
import { convert, Units } from "../../utils/conversion";
import { PlusIcon, TrashIcon } from "@radix-ui/react-icons";
import { deleteMachine } from "../../data/mutations/machine";
import { ResourceKind } from "../../types/resource";

interface Props {
    id: string;
}

export const VmsView: React.FC<Props> = ({ id }) => {
    const navigate = useNavigate();

    const data = useQuery({
        queryKey: [id, `machines`],
        queryFn: () => getMachines(id),
    });

    const { mutate: deleteMutation } = useMutation({
        mutationFn: deleteMachine,
    });

    if (data.isError) {
        console.log(data.error);
        return <p>error</p>;
    }

    return (
        <Box pt="3">
            <Box>
                <Button onClick={() => navigate(`/nodes/${id}/create-vm`)}>
                    Create
                </Button>
            </Box>

            <Table.Root layout="auto">
                <Table.Header>
                    <Table.Row>
                        <Table.ColumnHeaderCell>Status</Table.ColumnHeaderCell>
                        <Table.ColumnHeaderCell>Name</Table.ColumnHeaderCell>
                        <Table.ColumnHeaderCell>Cores</Table.ColumnHeaderCell>
                        <Table.ColumnHeaderCell>Memory</Table.ColumnHeaderCell>
                        <Table.ColumnHeaderCell>
                            Mac (network)
                        </Table.ColumnHeaderCell>
                        <Table.ColumnHeaderCell>
                            Interfaces
                        </Table.ColumnHeaderCell>
                        <Table.ColumnHeaderCell>Drives</Table.ColumnHeaderCell>
                        <Table.ColumnHeaderCell>
                            <Badge color="pink">
                                <PlusIcon />
                            </Badge>
                        </Table.ColumnHeaderCell>
                    </Table.Row>
                </Table.Header>

                {data.data?.list?.map((resource) => {
                    const machine = resource.spec!;

                    const memory = convert(
                        machine.topology.memory,
                        Units.Bytes,
                        Units.Gigabyte,
                    );

                    return (
                        <Table.Row key={machine.uuid}>
                            <Table.RowHeaderCell>
                                <Badge color="green">Running</Badge>
                            </Table.RowHeaderCell>
                            <Table.RowHeaderCell>
                                {machine.name}
                            </Table.RowHeaderCell>
                            <Table.Cell>
                                <Badge color="purple">
                                    {machine.topology.cores *
                                        machine.topology.threads}
                                </Badge>
                            </Table.Cell>
                            <Table.Cell>
                                <Badge color="purple">{memory} Gb</Badge>
                            </Table.Cell>
                            <Table.Cell>
                                {machine.interfaces[0].mac} (
                                <Badge color="amber">
                                    {machine.interfaces[0].network}
                                </Badge>
                                )
                            </Table.Cell>
                            <Table.Cell>
                                <Badge color="amber">
                                    {machine.interfaces.length}
                                </Badge>
                            </Table.Cell>
                            <Table.Cell>
                                <Badge color="amber">
                                    {machine.disks.length}
                                </Badge>
                            </Table.Cell>
                            <Table.Cell>
                                <IconButton
                                    color="red"
                                    variant="soft"
                                    size="1"
                                    onClick={() => {
                                        // TODO: This whole thing should bring up an alert.

                                        console.log("deleting");

                                        deleteMutation({
                                            kind: ResourceKind.Machine,
                                            id: resource.id,
                                        });
                                    }}
                                >
                                    <TrashIcon width="15" height="15" />
                                </IconButton>
                            </Table.Cell>
                        </Table.Row>
                    );
                })}
            </Table.Root>
        </Box>
    );
};
