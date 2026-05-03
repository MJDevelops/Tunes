import { useState } from "react";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Field, FieldError } from "./ui/field";
import z from "zod";

const urlSchema = z.url();

const AddDownloadItem = ({
  onChange,
  children,
}: {
  onChange: (val: string) => void;
  children?: React.ReactNode;
}) => {
  const [isEditing, setIsEditing] = useState(true);
  const [source, setSource] = useState("");
  const [valid, setValid] = useState(false);

  const handleConfirm = () => {
    setIsEditing(false);
    onChange(source);
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setSource(e.target.value);
    setValid(urlSchema.safeParse(e.target.value).success);
  };

  return (
      <Field orientation="horizontal" data-invalid={!valid}>
        {!valid && <FieldError>Enter a valid URL</FieldError>}
        <Input
          onChange={handleChange}
          value={source}
          disabled={!isEditing}
          placeholder="Enter download source"
          aria-invalid={!valid}
        />
        {isEditing ? (
          <Button disabled={!valid} onClick={handleConfirm}>
            Confirm
          </Button>
        ) : (
          <Button variant="outline" onClick={() => setIsEditing(true)}>
            Edit
          </Button>
        )}
        {children}
      </Field>
  );
};

export default AddDownloadItem;
