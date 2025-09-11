import { Home } from "lucide-react";
import { Link } from "react-router";
import {
  Sidebar,
  SidebarContent,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "@/components/ui/sidebar";

type SidebarItem = {
  title: string;
  icon?: React.ReactNode;
  url: string;
};

type SidebarItems = {
  groupName: string;
  items: SidebarItem[];
};

const items: SidebarItems[] = [
  {
    groupName: "Your Library",
    items: [
      {
        title: "Home",
        icon: <Home />,
        url: "/",
      },
    ],
  },
  {
    groupName: "Download",
    items: [
      {
        title: "yt-dlp",
        url: "/download/yt-dlp",
      },
    ],
  },
];

export function AppSidebar() {
  return (
    <Sidebar>
      <SidebarContent>
        {items.map((group) => (
          <SidebarGroup key={group.groupName}>
            <SidebarGroupLabel>{group.groupName}</SidebarGroupLabel>
            <SidebarGroupContent>
              <SidebarMenu>
                {group.items.map((item) => (
                  <SidebarMenuItem key={item.title}>
                    <SidebarMenuButton asChild>
                      <Link to={item.url}>
                        {item.icon}
                        <span>{item.title}</span>
                      </Link>
                    </SidebarMenuButton>
                  </SidebarMenuItem>
                ))}
              </SidebarMenu>
            </SidebarGroupContent>
          </SidebarGroup>
        ))}
      </SidebarContent>
    </Sidebar>
  );
}
