import { Button as SButton } from "@/components/ui/button";
import { motion } from "motion/react";

const MotionButton = motion.create(SButton);

export default function Button({
  children,
  ...props
}: React.ComponentPropsWithoutRef<typeof MotionButton>) {
  return (
    <MotionButton whileHover={{ scale: 1.05 }} {...props}>
      {children}
    </MotionButton>
  );
}

export { Button as MotionButton };
