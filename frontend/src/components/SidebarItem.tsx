import { motion } from "motion/react";
import { SidebarMenuItem } from "@/components/ui/sidebar";

const MotionSidebarItem = motion.create(SidebarMenuItem);

export default function SidebarItem({
  children,
  ...props
}: React.ComponentPropsWithoutRef<typeof MotionSidebarItem>) {
  return (
    <MotionSidebarItem whileHover={{ scale: 1.05 }} {...props}>
      {children}
    </MotionSidebarItem>
  );
}
