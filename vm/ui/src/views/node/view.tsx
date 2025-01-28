import {
    AlertDialog,
    Box,
    Button,
    Card,
    Flex,
    Grid,
    Separator,
    Tabs,
    Text,
} from "@radix-ui/themes";
import { PoolsView } from "./pools";
import { useParams } from "react-router";
import { VmsView } from "./vms";
import { PieChart } from "@mui/x-charts";
import { getNode } from "../../data/queries/nodes";
import { useQuery } from "@tanstack/react-query";

export const NodeView: React.FC<{}> = () => {
    const { id } = useParams<{ id: string }>();
    const data = useQuery({
        queryKey: ["nodes", id],
        queryFn: () => getNode(id!),
    });

    if (data.isLoading && !data.isError) {
        return <Text>Loading...</Text>;
    }

    return (
        <Box p="9">
            <Text size="8">{data.data?.list[0].spec?.hostname}</Text>
            <Box pt="3">
                <Tabs.Root defaultValue="overview">
                    <Tabs.List>
                        <Tabs.Trigger value="overview">Overview</Tabs.Trigger>
                        <Tabs.Trigger value="vms">
                            Virtual Machines
                        </Tabs.Trigger>
                        <Tabs.Trigger value="pools">Storage Pools</Tabs.Trigger>
                        <Tabs.Trigger value="networks">Networks</Tabs.Trigger>
                    </Tabs.List>

                    <Box pt="3">
                        <Tabs.Content value="overview">
                            <Grid columns="3" gap="4">
                                <Box gridColumn="1/3">
                                    <Card size="2"></Card>
                                </Box>
                                <Box gridColumnStart="3">
                                    <Card size="2">
                                        <Box>
                                            <Box>
                                                <Text size="2" color="gray">
                                                    Uptime
                                                </Text>
                                            </Box>
                                            <Text>24 hours</Text>
                                        </Box>
                                        <Separator size="4" mt="2" mb="2" />
                                        <Box>
                                            <Box>
                                                <Text size="2" color="gray">
                                                    Network
                                                </Text>
                                            </Box>
                                            <Text>10.100.0.101</Text>
                                        </Box>
                                        <Separator size="4" mt="2" mb="2" />
                                        <Box>
                                            <Box>
                                                <Text size="2" color="gray">
                                                    Kernel
                                                </Text>
                                            </Box>
                                            <Text>
                                                {
                                                    data.data?.list[0].spec!
                                                        .kernel
                                                }
                                            </Text>
                                        </Box>
                                        <Separator size="4" mt="2" mb="2" />
                                        <Box>
                                            <AlertDialog.Root>
                                                <AlertDialog.Trigger>
                                                    <Button
                                                        variant="solid"
                                                        color="red"
                                                    >
                                                        Reboot
                                                    </Button>
                                                </AlertDialog.Trigger>
                                                <AlertDialog.Content maxWidth="450px">
                                                    <AlertDialog.Title>
                                                        Reboot
                                                    </AlertDialog.Title>
                                                    <AlertDialog.Description size="2">
                                                        Are you sure? This node
                                                        will be rebooted and all
                                                        workloads paused
                                                    </AlertDialog.Description>

                                                    <Flex
                                                        gap="3"
                                                        mt="4"
                                                        justify="end"
                                                    >
                                                        <AlertDialog.Cancel>
                                                            <Button
                                                                variant="soft"
                                                                color="gray"
                                                            >
                                                                Cancel
                                                            </Button>
                                                        </AlertDialog.Cancel>
                                                        <AlertDialog.Action>
                                                            <Button
                                                                variant="solid"
                                                                color="red"
                                                            >
                                                                Reboot
                                                            </Button>
                                                        </AlertDialog.Action>
                                                    </Flex>
                                                </AlertDialog.Content>
                                            </AlertDialog.Root>
                                        </Box>
                                    </Card>
                                </Box>
                            </Grid>
                        </Tabs.Content>

                        <Tabs.Content value="vms">
                            <VmsView id={id!} />
                        </Tabs.Content>

                        <Tabs.Content value="pools">
                            <PoolsView id={id!} />
                        </Tabs.Content>
                        <Tabs.Content value="networks">
                            <PieChart
                                skipAnimation
                                series={[
                                    {
                                        data: [
                                            {
                                                id: 0,
                                                value: 10,
                                            },
                                            {
                                                id: 1,
                                                value: 15,
                                            },
                                            {
                                                id: 2,
                                                value: 20,
                                            },
                                        ],
                                        innerRadius: 30,
                                        paddingAngle: 5,
                                        cornerRadius: 5,
                                    },
                                ]}
                                width={400}
                                height={200}
                            />
                        </Tabs.Content>
                    </Box>
                </Tabs.Root>
            </Box>
        </Box>
    );
};
