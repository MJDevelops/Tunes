import tailwindcss from "@tailwindcss/vite";
import { defineConfig } from "vite";
import react, { reactCompilerPreset } from "@vitejs/plugin-react";
import { tanstackRouter } from "@tanstack/router-plugin/vite";
import wails from "@wailsio/runtime/plugins/vite";
import babel from "@rolldown/plugin-babel";

export default defineConfig({
  build: {
    rolldownOptions: {
      output: {
        codeSplitting: {
          groups: [
            {
              name: "libs",
              test: /node_modules/,
              minSize: 100000,
              maxSize: 250000,
              priority: 10,
            },
          ],
        },
      },
    },
  },
  resolve: {
    tsconfigPaths: true,
  },
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
    babel({
      presets: [reactCompilerPreset()],
    }),
    tailwindcss(),
    wails("./bindings"),
  ],
  server: {
    host: "127.0.0.1",
  },
});
