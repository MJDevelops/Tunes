import { Button } from "@/components/ui/button";
import { Plus } from "lucide-react";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";

export type MenuItem = {
  name: string;
  action: () => void;
};

const AddButton = ({ menuItems }: { menuItems: MenuItem[] }) => (
  <DropdownMenu>
    <DropdownMenuTrigger className="absolute right-0 bottom-0 m-4">
      <Button
        className="hover:cursor-pointer p-1"
        size="icon-sm"
        variant="outline"
        asChild
      >
        <Plus />
      </Button>
    </DropdownMenuTrigger>
    <DropdownMenuContent className="m-2">
      {menuItems.map((item, i) => (
        <DropdownMenuItem
          onClick={item.action}
          key={i}
          className="hover:cursor-pointer"
        >
          {item.name}
        </DropdownMenuItem>
      ))}
    </DropdownMenuContent>
  </DropdownMenu>
);

export default AddButton;
