import { create } from "zustand";
import { PendingDownloads } from "@bindings/internal/pkg/services/downloadservice";
import { Download } from "@bindings/internal/pkg/exec/ytdlp";

type DownloadState = {
  downloads: Download[];
  fetchPending: () => Promise<void>;
};

const useDownloadStore = create<DownloadState>()((set) => ({
  downloads: [],
  fetchPending: async () => set({ downloads: await PendingDownloads() }),
}));

export default useDownloadStore;
