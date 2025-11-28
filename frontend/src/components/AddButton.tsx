import { Button } from "@/components/ui/button";
import { Plus } from "lucide-react";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";

const AddButton = () => (
  <DropdownMenu>
    <DropdownMenuTrigger className="absolute right-0 bottom-0 m-4">
      <Button className="hover:cursor-pointer" size="icon-sm" variant="outline">
        <Plus />
      </Button>
    </DropdownMenuTrigger>
    <DropdownMenuContent>
      <DropdownMenuItem>Import Track</DropdownMenuItem>
      <DropdownMenuItem>New Download</DropdownMenuItem>
    </DropdownMenuContent>
  </DropdownMenu>
);

export default AddButton;
