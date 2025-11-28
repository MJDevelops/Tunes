import { Button } from "@/components/ui/button";
import { motion } from "motion/react";

const MotionButton = motion.create(Button);

function SidebarButton({
  children,
  ...props
}: React.ComponentPropsWithoutRef<typeof MotionButton>) {
  return (
    <MotionButton whileHover={{ scale: 1.05 }} {...props}>
      {children}
    </MotionButton>
  );
}

export default SidebarButton;
