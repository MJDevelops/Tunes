import { Outlet, Scripts, ScrollRestoration } from "react-router";
import { Toaster } from "~/components/ui/sonner";
import { SidebarProvider, SidebarTrigger } from "~/components/ui/sidebar";
import "./app.css";

export function Layout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en">
      <head>
        <meta charSet="utf-8" />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
      </head>
      <body>
        <SidebarProvider>
          <SidebarTrigger />
          {children}
          <Toaster />
        </SidebarProvider>
        <ScrollRestoration />
        <Scripts />
      </body>
    </html>
  );
}

export default function App() {
  return <Outlet />;
}
