import tailwindcss from "@tailwindcss/vite";
import { defineConfig } from "vite";
import tsconfigPaths from "vite-tsconfig-paths";
import react from "@vitejs/plugin-react-swc";
import { tanstackRouter } from "@tanstack/router-plugin/vite";
import wails from "@wailsio/runtime/plugins/vite";

export default defineConfig({
  plugins: [
    tanstackRouter({
      target: "react",
      autoCodeSplitting: true,
      routesDirectory: "./src/routes",
      generatedRouteTree: "./src/routeTree.gen.ts",
      routeFileIgnorePrefix: "-",
      quoteStyle: "double",
    }),
    react(),
    tailwindcss(),
    tsconfigPaths(),
    wails("./bindings"),
  ],
  build: {
    rollupOptions: {
      output: {
        manualChunks(id) {
          if (id.includes("node_modules")) {
            const modulePath = id.split("node_modules/")[1];
            const topLevelFolder = modulePath?.split("/")[0];
            if (topLevelFolder !== ".pnpm") {
              return topLevelFolder;
            }

            const scopedPackageName = modulePath?.split("/")[1];
            const chunkName =
              scopedPackageName?.split("@")[
                scopedPackageName.startsWith("@") ? 1 : 0
              ];

            return chunkName;
          }
        },
      },
    },
  },
});
