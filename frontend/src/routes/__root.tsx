import { Toaster } from "@/components/ui/sonner";
import { SidebarProvider, SidebarTrigger } from "@/components/ui/sidebar";
import { AppSidebar } from "@/components/AppSidebar";
import { createRootRoute, Outlet } from "@tanstack/react-router";
import { ThemeProvider } from "@/components/ThemeProvider";
import AddButton, { type MenuItem } from "@/components/AddButton";
import { useEffect, useState } from "react";
import AddDownloads from "@/components/AddDownloads";
import { Dialogs, Events } from "@wailsio/runtime";
import { useDownloadStore } from "@/stores";
import "@/app.css";

const RootLayout = () => {
  const [addDownloads, setAddDownloads] = useState(false);
  const { enqueueDownloads, removeDownload } = useDownloadStore();
  const menuItems: MenuItem[] = [
    {
      name: "Import Tracks",
      action: async () => {
        const tracks = await Dialogs.OpenFile({
          Title: "Select tracks",
          Filters: [
            { DisplayName: "Audio", Pattern: "*.mp3;*.ogg;*.wav;*.flac" },
          ],
          AllowsMultipleSelection: true,
        });
      },
    },
    {
      name: "Add Download Sources",
      action: () => setAddDownloads(true),
    },
  ];

  useEffect(() => {
    Events.On("tunes:dl:finished", (event) => {
      removeDownload(event.data);
    });

    return () => {
      Events.Off("tunes:dl:finished");
    };
  }, []);

  return (
    <ThemeProvider>
      <SidebarProvider>
        <AppSidebar />
        <SidebarTrigger />
        <Outlet />
        <AddDownloads
          onConfirm={(downloads) => enqueueDownloads(downloads)}
          onClose={() => setAddDownloads(false)}
          open={addDownloads}
        />
        <AddButton menuItems={menuItems} />
        <Toaster />
      </SidebarProvider>
    </ThemeProvider>
  );
};

export const Route = createRootRoute({ component: RootLayout });
