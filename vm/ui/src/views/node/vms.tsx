import { useQuery } from "@tanstack/react-query";
import { Badge, Box, Button, IconButton, Table } from "@radix-ui/themes";
import { getMachines } from "../../data/queries/machines";
import { useNavigate } from "react-router";
import { convert, Units } from "../../utils/conversion";

interface Props {
    id: string;
}

export const VmsView: React.FC<Props> = ({ id }) => {
    const navigate = useNavigate();

    const data = useQuery({
        queryKey: [id, `machines`],
        queryFn: () => getMachines(id),
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
                            <IconButton size="1">+</IconButton>
                        </Table.ColumnHeaderCell>
                    </Table.Row>
                </Table.Header>

                {data.data?.list.map((resource) => {
                    const machine = resource.spec!;

                    const memory = convert(
                        machine.topology.memory,
                        Units.Bytes,
                        Units.Gigabyte,
                    );

                    return (
                        <Table.Row
                            key={machine.uuid}
                            onClick={() => {
                                console.log("redirecting");
                            }}
                        >
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
                            <Table.Cell />
                        </Table.Row>
                    );
                })}
            </Table.Root>
        </Box>
    );
};
