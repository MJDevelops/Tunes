import { create } from "zustand";
import { PendingDownloads } from "~/wailsjs/go/main/App";
import { ytdlp } from "~/wailsjs/go/models";

type DownloadState = {
  downloads: ytdlp.Download[];
  fetchPending: () => Promise<void>;
};

const useDownloadStore = create<DownloadState>()((set) => ({
  downloads: [],
  fetchPending: async () => set({ downloads: await PendingDownloads() }),
}));

export default useDownloadStore;
