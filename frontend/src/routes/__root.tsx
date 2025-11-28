import { Toaster } from "@/components/ui/sonner";
import { SidebarProvider, SidebarTrigger } from "@/components/ui/sidebar";
import { AppSidebar } from "@/components/AppSidebar";
import { createRootRoute, Outlet } from "@tanstack/react-router";
import { ThemeProvider } from "@/components/ThemeProvider";
import "@/app.css";
import AddButton from "@/components/AddButton";

const RootLayout = () => (
  <ThemeProvider>
    <SidebarProvider>
      <AppSidebar />
      <SidebarTrigger />
      <Outlet />
      <AddButton />
      <Toaster />
    </SidebarProvider>
  </ThemeProvider>
);

export const Route = createRootRoute({ component: RootLayout });
