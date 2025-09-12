import { Toaster } from "@/components/ui/sonner";
import { SidebarProvider, SidebarTrigger } from "@/components/ui/sidebar";
import { AppSidebar } from "@/components/AppSidebar";
import { Outlet } from "react-router";
import { ThemeProvider } from "@/components/ThemeProvider";
import "./app.css";

export default function App() {
  return (
    <ThemeProvider>
      <SidebarProvider>
        <AppSidebar />
        <SidebarTrigger />
        <Outlet />
        <Toaster />
      </SidebarProvider>
    </ThemeProvider>
  );
}
