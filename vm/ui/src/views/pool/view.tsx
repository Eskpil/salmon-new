import {
    AlertDialog,
    Badge,
    Box,
    Button,
    Flex,
    Link,
    Separator,
    Table,
    Text,
} from "@radix-ui/themes";
import { useMutation, useQuery } from "@tanstack/react-query";
import { useParams } from "react-router";
import { getVolumes } from "../../data/queries/volumes";
import { Field, Form, Formik, FormikHelpers } from "formik";
import { convert, Units } from "../../utils/conversion";
import { createVolume } from "../../data/mutations/volumes";
import {
    CreateResourceInput,
    Resource,
    ResourceKind,
} from "../../types/resource";
import { getPool } from "../../data/queries/pools";
import { Pool } from "../../types/pool";
import { getNode } from "../../data/queries/nodes";
import { useNavigate } from "react-router";

interface CreateVolumeValues {
    name: string;
    capacity: number;
}

// TODO: Do not use link in breadcrum navigation. This causes fullpage reload.

const Title: React.FC<{ pool: Resource<Pool> }> = ({ pool }) => {
    const navigate = useNavigate();

    const node = useQuery({
        queryKey: ["nodes", pool.owner?.id],
        queryFn: () => getNode(pool.owner!.id!),
    });

    if (node.isError) {
        console.log(node.error);
        return <p>error</p>;
    }

    if (node.isLoading) {
        return <p>loading</p>;
    }

    return (
        <Box>
            <Link
                href=""
                color="purple"
                onClick={() => navigate(`/nodes/${node.data!.id}`)}
            >
                <Text size="6">{node.data?.spec?.hostname}</Text>
            </Link>
            <Text size="5" mr="1" ml="1">
                /
            </Text>
            <Text size="6">{pool.spec?.name}</Text>
        </Box>
    );
};

export const PoolView: React.FC<unknown> = () => {
    const { id } = useParams<{ id: string }>();

    const pool = useQuery({
        queryKey: ["pools", id],
        queryFn: () => getPool(id!),
    });

    const volumes = useQuery({
        queryKey: [id, `volumes`],
        queryFn: () => getVolumes(id!),
    });

    const { mutate } = useMutation({ mutationFn: createVolume });

    if (volumes.isError || pool.isError) {
        console.log(volumes.error, pool.error);
        return <p>error</p>;
    }

    if (pool.isLoading || volumes.isLoading) {
        return <p>loading</p>;
    }

    return (
        <Box p="9" width="100%">
            <Title pool={pool.data!} />
            <Box width="100%" pt="2">
                <Separator size="4" />
            </Box>
            <Box pt="3">
                <AlertDialog.Root>
                    <AlertDialog.Trigger>
                        <Button variant="solid" color="purple">
                            Create
                        </Button>
                    </AlertDialog.Trigger>
                    <AlertDialog.Content maxWidth="450px">
                        <Formik
                            initialValues={{ name: "volume123", capacity: 12 }}
                            onSubmit={(
                                values,
                                {
                                    setSubmitting,
                                }: FormikHelpers<CreateVolumeValues>,
                            ) => {
                                const capacity = convert(
                                    values.capacity,
                                    Units.Gigabyte,
                                    Units.Bytes,
                                );

                                values.capacity = capacity;

                                const input: CreateResourceInput = {
                                    owner_ref: {
                                        id: id!,
                                        kind: ResourceKind.StoragePool,
                                    },
                                    annotations: {},
                                    kind: ResourceKind.StorageVolume,
                                    spec: {
                                        name: values.name,
                                        capacity: values.capacity,
                                        allocation: values.capacity,
                                    },
                                };

                                mutate(input, {
                                    onSuccess: () => setSubmitting(false),
                                });
                            }}
                        >
                            <Form>
                                <AlertDialog.Title>
                                    Create Volume
                                </AlertDialog.Title>
                                <AlertDialog.Description size="2">
                                    <Box>
                                        <Box pb="1">
                                            <label htmlFor="name">Name</label>
                                        </Box>
                                        <Field
                                            id="name"
                                            name="name"
                                            placeholder="volume123"
                                        ></Field>
                                    </Box>
                                    <Box pt="3">
                                        <Box pb="1">
                                            <label htmlFor="capacity">
                                                Capacity
                                            </label>
                                        </Box>
                                        <Field
                                            id="capacity"
                                            name="capacity"
                                            placeholder="12"
                                        ></Field>
                                    </Box>
                                </AlertDialog.Description>

                                <Flex gap="3" mt="4" justify="end">
                                    <AlertDialog.Cancel>
                                        <Button variant="soft" color="red">
                                            Cancel
                                        </Button>
                                    </AlertDialog.Cancel>
                                    <AlertDialog.Action>
                                        <Button
                                            variant="solid"
                                            color="purple"
                                            type="submit"
                                        >
                                            Create
                                        </Button>
                                    </AlertDialog.Action>
                                </Flex>
                            </Form>
                        </Formik>
                    </AlertDialog.Content>
                </AlertDialog.Root>
            </Box>
            <Table.Root layout="auto">
                <Table.Header>
                    <Table.Row>
                        <Table.ColumnHeaderCell>Name</Table.ColumnHeaderCell>
                        <Table.ColumnHeaderCell>Key</Table.ColumnHeaderCell>
                        <Table.ColumnHeaderCell>
                            Virtual Machine
                        </Table.ColumnHeaderCell>
                        <Table.ColumnHeaderCell>Usage</Table.ColumnHeaderCell>
                        <Table.ColumnHeaderCell>Phase</Table.ColumnHeaderCell>
                    </Table.Row>
                </Table.Header>

                <Table.Body>
                    {volumes.data?.list?.map((resource) => {
                        const volume = resource.spec!;

                        const vm_name = resource.annotations!["vm.name"];

                        const capacity_gb = Math.round(
                            convert(
                                volume.capacity,
                                Units.Bytes,
                                Units.Gigabyte,
                            ),
                        );

                        const allocated_gb = Math.round(
                            convert(
                                volume.allocation,
                                Units.Bytes,
                                Units.Gigabyte,
                            ),
                        );

                        return (
                            <Table.Row key={resource.id}>
                                <Table.RowHeaderCell>
                                    {volume.name}
                                </Table.RowHeaderCell>
                                <Table.Cell>{volume.key}</Table.Cell>
                                <Table.Cell>
                                    {vm_name ? (
                                        <Badge color="purple">{vm_name}</Badge>
                                    ) : (
                                        <Badge color="red">unassigned</Badge>
                                    )}
                                </Table.Cell>
                                <Table.Cell>
                                    <Badge color="green">
                                        {allocated_gb} Gb
                                    </Badge>
                                    /
                                    <Badge color="purple">
                                        {capacity_gb} Gb
                                    </Badge>
                                </Table.Cell>
                                <Table.Cell>
                                    <Badge color="amber">
                                        {resource.status.phase}
                                    </Badge>
                                </Table.Cell>
                            </Table.Row>
                        );
                    })}
                </Table.Body>
            </Table.Root>
        </Box>
    );
};
