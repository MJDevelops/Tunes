import { createHashRouter, RouterProvider } from "react-router";
import { createRoot } from "react-dom/client";
import { StrictMode } from "react";
import Home from "@/components/Home";
import App from "./App";
import DownloadOverview from "@/components/DownloadOverview";

const router = createHashRouter([
  {
    path: "/",
    Component: App,
    children: [
      { index: true, Component: Home },
      { path: "/downloads", Component: DownloadOverview },
    ],
  },
]);

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <RouterProvider router={router} />
  </StrictMode>
);
