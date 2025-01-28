import { Box, Text } from "@radix-ui/themes";
import { useParams } from "react-router";

interface Props {}

export const PoolView: React.FC<Props> = () => {
    const { id } = useParams<{ id: string }>();

    return (
        <Box>
            <Text>Hello {id}</Text>
        </Box>
    );
};
