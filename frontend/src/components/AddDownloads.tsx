import { useState } from "react";
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import AddDownloadItem from "./AddDownloadItem";
import { ScrollArea } from "@/components/ui/scroll-area";
import { v4 as uuidv4 } from "uuid";
import { VisuallyHidden } from "@radix-ui/react-visually-hidden";

type AddDownloadsProps = {
  open: boolean;
  onClose: () => void;
  onConfirm: (downloads: string[]) => void;
};

type AddDownloadItemWithKeys = {
  source: string;
  key: string;
};

const AddDownloads = ({
  open = false,
  onClose,
  onConfirm,
}: AddDownloadsProps) => {
  const [downloads, setDownloads] = useState<AddDownloadItemWithKeys[]>([]);

  const changeDownload = (changed: string, downloadId: string) => {
    setDownloads(
      downloads.map((download) => {
        if (download.key === downloadId) {
          return { ...download, source: changed };
        }
        return download;
      }),
    );
  };

  const addDownload = () => {
    setDownloads([
      ...downloads,
      {
        source: "",
        key: uuidv4(),
      },
    ]);
  };

  const handleRemove = (id: string) => {
    setDownloads(downloads.filter((download) => download.key !== id));
  };

  const handleClose = () => {
    onClose();
    setDownloads([]);
  };

  return (
    <Dialog open={open}>
      <DialogContent>
        <DialogHeader>Add Download Sources</DialogHeader>
        <VisuallyHidden>
          <DialogTitle>Download Sources Dialog</DialogTitle>
        </VisuallyHidden>
        <ScrollArea className="flex w-full h-52 flex-col gap-1">
          {downloads.length > 0 &&
            downloads.map((download) => (
                <AddDownloadItem
                  key={download.key}
                  onChange={(val) => changeDownload(val, download.key)}
                >
                <Button
                  onClick={() => handleRemove(download.key)}
                >
                  Remove
                </Button>
                </AddDownloadItem>
            ))}
          <Button onClick={addDownload}>Insert Download</Button>
        </ScrollArea>
        <DialogFooter>
          <DialogClose asChild>
            <Button onClick={handleClose} variant="outline">
              Cancel
            </Button>
          </DialogClose>
          <Button
            onClick={() =>
              onConfirm(downloads.map((download) => download.source))
            }
            disabled={downloads.some((download) => download.source === "") || downloads.length === 0}
          >
            Add
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
};

export default AddDownloads;
