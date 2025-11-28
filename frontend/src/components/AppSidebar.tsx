import {
  Home,
  ListMusic,
  Disc2,
  Settings,
  ArrowDownToLine,
} from "lucide-react";
import {
  Sidebar,
  SidebarContent,
  SidebarGroup,
  SidebarGroupContent,
  SidebarMenu,
  SidebarMenuButton,
} from "@/components/ui/sidebar";
import SidebarItem from "@/components/SidebarItem";
import SidebarButton from "@/components/SidebarButton";
import { Link, useLocation } from "@tanstack/react-router";
import { cn } from "@/lib/utils";

type SidebarItem = {
  title: string;
  icon?: React.ReactNode;
  url: string;
};

const items: SidebarItem[] = [
  {
    title: "Home",
    icon: <Home />,
    url: "/",
  },
  {
    title: "Playlists",
    icon: <ListMusic />,
    url: "/playlists",
  },
  {
    title: "Albums",
    icon: <Disc2 />,
    url: "/albums",
  },
  {
    title: "All Downloads",
    icon: <ArrowDownToLine />,
    url: "/downloads",
  },
];

export function AppSidebar() {
  const pathname = useLocation({
    select: (location) => location.pathname,
  });

  return (
    <Sidebar className="flex flex-col justify-between p-1 rounded-md">
      <SidebarContent>
        <SidebarGroup>
          <SidebarGroupContent>
            <SidebarMenu>
              {items.map((item) => (
                <SidebarItem key={item.title}>
                  <SidebarMenuButton
                    className={cn(
                      pathname === item.url && "bg-emerald-600 rounded-md",
                      "hover:bg-emerald-800 active:bg-emerald-800"
                    )}
                    asChild
                  >
                    <Link to={item.url}>
                      {item.icon}
                      <span>{item.title}</span>
                    </Link>
                  </SidebarMenuButton>
                </SidebarItem>
              ))}
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
      </SidebarContent>
      <SidebarButton
        className="m-2 p-2 hover:cursor-pointer flex justify-baseline gap-2 align-center"
        variant="outline"
      >
        <Settings />
        Settings
      </SidebarButton>
    </Sidebar>
  );
}
