import { useState } from "react";
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogFooter,
  DialogHeader,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import AddDownloadItem from "./AddDownloadItem";
import { ScrollArea } from "@/components/ui/scroll-area";
import { v4 as uuidv4 } from "uuid";

type AddDownloadsProps = {
  open: boolean;
  onClose: () => void;
  onConfirm: (downloads: string[]) => void;
};

type AddDownloadItemWithKeys = {
  source: string;
  downloadId: string;
  buttonId: string;
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
        if (download.downloadId === downloadId) {
          return { ...download, download: changed };
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
        downloadId: uuidv4(),
        buttonId: uuidv4(),
      },
    ]);
  };

  const handleRemove = (id: string) => {
    setDownloads(downloads.filter((download) => download.downloadId !== id));
  };

  const handleClose = () => {
    onClose();
    setDownloads([]);
  };

  return (
    <Dialog open={open}>
      <DialogContent>
        <DialogHeader>Add Download Sources</DialogHeader>
        <ScrollArea className="flex w-full h-52 flex-col gap-1">
          {downloads.length > 0 &&
            downloads.map((download) => (
              <div className="flex gap-1">
                <AddDownloadItem
                  key={download.downloadId}
                  className="flex gap-1"
                  onChange={(val) => changeDownload(val, download.downloadId)}
                />
                <Button
                  key={download.buttonId}
                  onClick={() => handleRemove(download.downloadId)}
                >
                  Remove
                </Button>
              </div>
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
          >
            Add
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
};

export default AddDownloads;
