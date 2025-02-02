import { AlertDialog, Badge, Box, Button, Flex, Table } from "@radix-ui/themes";
import { useMutation, useQuery } from "@tanstack/react-query";
import { useParams } from "react-router";
import { getVolumes } from "../../data/queries/volumes";
import { Field, Form, Formik, FormikHelpers } from "formik";
import { convert, Units } from "../../utils/conversion";
import { createVolume } from "../../data/mutations/volumes";
import { CreateResourceInput, ResourceKind } from "../../types/resource";

interface Props {}

interface CreateVolumeValues {
    name: string;
    capacity: number;
}

export const PoolView: React.FC<Props> = () => {
    const { id } = useParams<{ id: string }>();

    const data = useQuery({
        queryKey: [id, `volumes`],
        queryFn: () => getVolumes(id!),
    });

    const { mutate } = useMutation({ mutationFn: createVolume });

    if (data.isError) {
        console.log(data.error);
        return <p>error</p>;
    }

    return (
        <Box pt="3">
            <Box>
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
                        <Table.ColumnHeaderCell>Usage</Table.ColumnHeaderCell>
                        <Table.ColumnHeaderCell>Phase</Table.ColumnHeaderCell>
                    </Table.Row>
                </Table.Header>

                {data.data?.list.map((resource) => {
                    const volume = resource.spec!;

                    const capacity_gb = Math.round(
                        convert(volume.capacity, Units.Bytes, Units.Gigabyte),
                    );

                    const allocated_gb = Math.round(
                        convert(volume.allocation, Units.Bytes, Units.Gigabyte),
                    );

                    return (
                        <Table.Row key={resource.id}>
                            <Table.RowHeaderCell>
                                {volume.name}
                            </Table.RowHeaderCell>
                            <Table.Cell>{volume.key}</Table.Cell>
                            <Table.Cell>
                                <Badge color="green">{allocated_gb} Gb</Badge>/
                                <Badge color="purple">{capacity_gb} Gb</Badge>
                            </Table.Cell>
                            <Table.Cell>
                                <Badge color="amber">
                                    {resource.status.phase}
                                </Badge>
                            </Table.Cell>
                        </Table.Row>
                    );
                })}
            </Table.Root>
        </Box>
    );
};
