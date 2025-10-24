import {
  Home,
  ListMusic,
  Disc2,
  Settings,
  ArrowDownToLine,
} from "lucide-react";
import { SiYoutube, SiSoundcloud } from "@icons-pack/react-simple-icons";
import {
  Sidebar,
  SidebarContent,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarMenu,
  SidebarMenuButton,
} from "@/components/ui/sidebar";
import { useAnimate } from "motion/react";
import SidebarItem from "@/components/SidebarItem";
import { MotionButton } from "@/components/Button";
import { Link } from "@tanstack/react-router";

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
    ],
  },
  {
    groupName: "Download",
    items: [
      {
        title: "All Downloads",
        icon: <ArrowDownToLine />,
        url: "/downloads",
      },
      {
        title: "YouTube",
        icon: <SiYoutube />,
        url: "/download/youtube",
      },
      {
        title: "Soundcloud",
        icon: <SiSoundcloud />,
        url: "/download/soundcloud",
      },
    ],
  },
];

export function AppSidebar() {
  const [scope, animate] = useAnimate();

  return (
    <Sidebar className="flex flex-col justify-between p-1 rounded-md">
      <SidebarContent>
        {items.map((group) => (
          <SidebarGroup key={group.groupName}>
            <SidebarGroupLabel>{group.groupName}</SidebarGroupLabel>
            <SidebarGroupContent>
              <SidebarMenu>
                {group.items.map((item) => (
                  <SidebarItem key={item.title}>
                    <SidebarMenuButton asChild>
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
        ))}
      </SidebarContent>
      <MotionButton
        className="m-2 p-2 hover:cursor-pointer flex justify-baseline gap-2 align-center"
        onMouseEnter={() => animate(scope.current, { rotate: -90 })}
        onMouseLeave={() => animate(scope.current, { rotate: 90 })}
        variant="outline"
      >
        <Settings ref={scope} />
        Settings
      </MotionButton>
    </Sidebar>
  );
}
