import { create } from "zustand";
import {
  GetAlbumTracks,
  GetPlaylist,
} from "@bindings/internal/pkg/services/dbservice";
import { Track } from "@bindings/internal/pkg/db/models";
import { shuffleArray } from "@/lib/utils";

const initialState = { tracks: [] };

type QueueState = {
  tracks: Track[];
  enqueue: (track: Track) => void;
  enqueuePlaylist: (playlistId: number) => Promise<void>;
  enqueueAlbum: (albumId: number) => Promise<void>;
  shuffle: () => void;
  next: () => Track | undefined;
};

const useQueueStore = create<QueueState>()((set) => ({
  ...initialState,
  enqueue: (track) => set((s) => ({ tracks: [...s.tracks, track] })),
  enqueuePlaylist: async (playlistId) =>
    set({ tracks: (await GetPlaylist(playlistId)).Tracks }),
  enqueueAlbum: async (albumId) =>
    set({ tracks: await GetAlbumTracks(albumId) }),
  shuffle: () => {
    set((s) => {
      s.tracks = shuffleArray(s.tracks);
      return s;
    });
  },
  next: () => {
    let track: Track | undefined;
    set((s) => {
      track = s.tracks.shift();
      return { ...s, tracks: s.tracks };
    });
    return track;
  },
}));

export default useQueueStore;
