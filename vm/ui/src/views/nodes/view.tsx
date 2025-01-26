import { useQuery } from "@tanstack/react-query";
import { getNodes } from "../../data/queries/nodes";
import { Badge, Box, Container, Table, Text } from "@radix-ui/themes";
import { redirect } from "react-router";
import { useNavigate } from "react-router";

export const NodesView: React.FC<{}> = () => {
    const navigate = useNavigate();
    const nodes = useQuery({ queryKey: ["nodes"], queryFn: getNodes });

    if (nodes.error) {
        return <div>Error fetching nodes</div>;
    }

    return (
        <Container p="30" pt="9">
            <Table.Root>
                <Table.Header>
                    <Table.Row>
                        <Table.ColumnHeaderCell></Table.ColumnHeaderCell>
                        <Table.ColumnHeaderCell>Name</Table.ColumnHeaderCell>
                        <Table.ColumnHeaderCell>Url</Table.ColumnHeaderCell>
                        <Table.ColumnHeaderCell>
                            Machines
                        </Table.ColumnHeaderCell>
                    </Table.Row>
                </Table.Header>

                <Table.Body>
                    {nodes.data?.list.map((node) => {
                        return (
                            <Table.Row
                                key={node.name}
                                onClick={() => {
                                    console.log("redirecting");

                                    navigate(`/nodes/${node.name}`);
                                }}
                            >
                                <Table.RowHeaderCell>
                                    <Badge color="green">x</Badge>
                                </Table.RowHeaderCell>
                                <Table.RowHeaderCell>
                                    {node.name}
                                </Table.RowHeaderCell>
                                <Table.Cell>{node.url}</Table.Cell>
                                <Table.Cell>
                                    <Badge color="green">
                                        {node.active_machines}
                                    </Badge>
                                    /{node.total_machines}
                                </Table.Cell>
                            </Table.Row>
                        );
                    })}
                </Table.Body>
            </Table.Root>
        </Container>
    );
};
