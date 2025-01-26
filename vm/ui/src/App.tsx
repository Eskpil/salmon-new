import { Routes, Route, BrowserRouter } from "react-router";
import { NodesView } from "./views/nodes/view";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import "@radix-ui/themes/styles.css";
import { Theme } from "@radix-ui/themes";
import { NodeView } from "./views/node/view";

const queryClient = new QueryClient();

function App() {
    return (
        <Theme appearance="dark">
            <QueryClientProvider client={queryClient}>
                <BrowserRouter>
                    <Routes>
                        <Route index element={<NodesView />} />
                        <Route path="nodes/:id" element={<NodeView />} />
                    </Routes>
                </BrowserRouter>
            </QueryClientProvider>
        </Theme>
    );
}

export default App;
