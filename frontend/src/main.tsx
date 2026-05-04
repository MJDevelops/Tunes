import { createRoot } from "react-dom/client";
import { StrictMode } from "react";
import { routeTree } from "./routeTree.gen";
import { createRouter, RouterProvider, createHashHistory } from "@tanstack/react-router";
import { useGSAP } from "@gsap/react";
import gsap from "gsap";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import "@wailsio/runtime";

gsap.registerPlugin(useGSAP);

const hashHistory = createHashHistory();
const router = createRouter({ routeTree, history: hashHistory });

declare module "@tanstack/react-router" {
  interface Register {
    router: typeof router;
  }
}

const queryClient = new QueryClient();

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <QueryClientProvider client={queryClient}>
      <RouterProvider router={router} />
    </QueryClientProvider>
  </StrictMode>,
);
