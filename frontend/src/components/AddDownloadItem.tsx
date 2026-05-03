import { useState } from "react";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";

const AddDownloadItem = ({
  onChange,
  className,
}: {
  onChange: (val: string) => void;
  className?: string;
}) => {
  const [isEditing, setIsEditing] = useState(true);
  const [source, setSource] = useState("");

  const handleConfirm = () => {
    setIsEditing(false);
    onChange(source);
  };

  return (
    <div className={className}>
      <>
        <Input
          onChange={(e) => setSource(e.target.value)}
          value={source}
          disabled={!isEditing}
          placeholder="Enter download source"
        />
        {isEditing ? (
          <Button onClick={handleConfirm}>Confirm</Button>
        ) : (
          <Button onClick={() => setIsEditing(true)}>Edit</Button>
        )}
      </>
    </div>
  );
};

export default AddDownloadItem;
