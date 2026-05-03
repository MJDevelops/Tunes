import { create } from "zustand";
import {
  EnqueueDownload,
  PendingDownloads,
} from "@bindings/internal/pkg/services/downloadservice";
import { Download } from "@bindings/internal/pkg/exec/ytdlp";

type DownloadState = {
  downloads: Download[];
  enqueueDownloads: (downloads: string[]) => void;
  fetchPending: () => Promise<void>;
  removeDownload: (id: string) => void;
};

const useDownloadStore = create<DownloadState>()((set) => ({
  downloads: [],
  fetchPending: async () => set({ downloads: await PendingDownloads() }),
  enqueueDownloads: (downloads) => {
    const newDownloads: Download[] = [];
    downloads.forEach(async (download) => {
      newDownloads.push(
        new Download({ ID: await EnqueueDownload(download), Url: download }),
      );
    });
    set((s) => ({ downloads: [...s.downloads, ...newDownloads] }));
  },
  removeDownload: (id) => {
    set((s) => ({ downloads: s.downloads.filter((down) => down.ID !== id) }));
  },
}));

export default useDownloadStore;
