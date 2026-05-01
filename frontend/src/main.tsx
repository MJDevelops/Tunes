import { createRoot } from "react-dom/client";
import { StrictMode } from "react";
import { routeTree } from "./routeTree.gen";
import { createRouter, RouterProvider } from "@tanstack/react-router";
import { useGSAP } from "@gsap/react";
import gsap from "gsap";
import "@wailsio/runtime";

gsap.registerPlugin(useGSAP);

const router = createRouter({ routeTree });

declare module "@tanstack/react-router" {
  interface Register {
    router: typeof router;
  }
}

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <RouterProvider router={router} />
  </StrictMode>,
);
