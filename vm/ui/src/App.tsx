import { Routes, Route, BrowserRouter } from "react-router";
import { NodesView } from "./views/nodes/view";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import "@radix-ui/themes/styles.css";
import { Theme, Text } from "@radix-ui/themes";
import { NodeView } from "./views/node/view";
import { PoolView } from "./views/pool/view";
import { CreateVmView } from "./views/vm/create";

const queryClient = new QueryClient();

function App() {
    return (
        <Theme appearance="dark" accentColor="purple">
            <QueryClientProvider client={queryClient}>
                <BrowserRouter>
                    <Routes>
                        <Route index element={<NodesView />} />
                        <Route path="nodes/:id" element={<NodeView />} />
                        <Route
                            path="nodes/:id/create-vm"
                            element={<CreateVmView />}
                        />
                        <Route path="pools/:id" element={<PoolView />} />
                        <Route path="*" element={<Text>404 </Text>} />
                    </Routes>
                </BrowserRouter>
            </QueryClientProvider>
        </Theme>
    );
}

export default App;
