import { useQuery } from "@tanstack/react-query";
import { getPools } from "../../data/queries/pools";
import { Badge, Box, Table } from "@radix-ui/themes";
import { useNavigate } from "react-router";

interface Props {
    id: string;
}

export const PoolsView: React.FC<Props> = ({ id }) => {
    const navigate = useNavigate();

    const data = useQuery({
        queryKey: [id, `pools`],
        queryFn: () => getPools(id),
    });

    if (data.isError) {
        console.log(data.error);
        return <p>error</p>;
    }

    return (
        <Box pt="3">
            <Table.Root layout="auto">
                <Table.Header>
                    <Table.Row>
                        <Table.ColumnHeaderCell>Name</Table.ColumnHeaderCell>
                        <Table.ColumnHeaderCell>Volumes</Table.ColumnHeaderCell>
                        <Table.ColumnHeaderCell>Usage</Table.ColumnHeaderCell>
                        <Table.ColumnHeaderCell>Backend</Table.ColumnHeaderCell>
                    </Table.Row>
                </Table.Header>

                {data.data?.list.map((resource) => {
                    const pool = resource.spec!;

                    const capacity_gb = Math.round(pool.capacity / 1000000000);
                    const allocated_gb = Math.round(
                        pool.allocated / 1000000000,
                    );

                    return (
                        <Table.Row
                            key={pool.id}
                            onClick={() => {
                                navigate(`/pools/${resource.id}`);
                            }}
                        >
                            <Table.RowHeaderCell>
                                {pool.name}
                            </Table.RowHeaderCell>
                            <Table.Cell>{pool.allocated_volumes}</Table.Cell>
                            <Table.Cell>
                                <Badge color="green">{allocated_gb} Gb</Badge>/
                                <Badge color="purple">{capacity_gb} Gb</Badge>
                            </Table.Cell>
                            <Table.Cell>
                                <Badge color="amber">{pool.kind}</Badge>
                            </Table.Cell>
                        </Table.Row>
                    );
                })}
            </Table.Root>
        </Box>
    );
};
