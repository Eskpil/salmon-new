import { useQuery } from "@tanstack/react-query";
import { getNodes } from "../../data/queries/nodes";
import { Badge, Box, Table, Text } from "@radix-ui/themes";
import { useNavigate } from "react-router";

export const NodesView: React.FC<unknown> = () => {
    const navigate = useNavigate();
    const nodes = useQuery({ queryKey: ["nodes"], queryFn: getNodes });

    if (nodes.error) {
        return <div>Error fetching nodes</div>;
    }

    return (
        <Box p="9">
            <Text size="8">Nodes</Text>
            <Box pt="3">
                <Table.Root>
                    <Table.Header>
                        <Table.Row>
                            <Table.ColumnHeaderCell></Table.ColumnHeaderCell>
                            <Table.ColumnHeaderCell>
                                Name
                            </Table.ColumnHeaderCell>
                            <Table.ColumnHeaderCell>Url</Table.ColumnHeaderCell>
                            <Table.ColumnHeaderCell>
                                Machines
                            </Table.ColumnHeaderCell>
                        </Table.Row>
                    </Table.Header>

                    <Table.Body>
                        {nodes.data?.list.map((node) => {
                            const color = "green";

                            const annotations = new Map();

                            Object.entries(node.annotations!).forEach(
                                (annotation) => {
                                    annotations.set(
                                        annotation[0],
                                        annotation[1],
                                    );
                                },
                            );

                            const url = annotations.get("node.url")!;

                            console.log(annotations);

                            return (
                                <Table.Row
                                    key={node.id}
                                    onClick={() => {
                                        navigate(`/nodes/${node.id}`);
                                    }}
                                >
                                    <Table.RowHeaderCell>
                                        <Badge color={color as any}>
                                            Connected
                                        </Badge>
                                    </Table.RowHeaderCell>
                                    <Table.RowHeaderCell>
                                        {node.spec!.hostname}
                                    </Table.RowHeaderCell>
                                    <Table.Cell>{url}</Table.Cell>
                                    <Table.Cell>
                                        <Badge color="green">
                                            {node.spec!.active_machines}
                                        </Badge>
                                        /
                                        <Badge color="purple">
                                            {node.spec!.total_machines}
                                        </Badge>
                                    </Table.Cell>
                                </Table.Row>
                            );
                        })}
                    </Table.Body>
                </Table.Root>
            </Box>
        </Box>
    );
};
