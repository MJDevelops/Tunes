import { create } from "zustand";
import { GetPlaylistTracks } from "@bindings/internal/pkg/services/audioservice";
import { Track } from "@bindings/internal/pkg/db/models";

const initialState = { tracks: [] };

type QueueState = {
  tracks: Track[];
  enqueue: (track: Track) => Promise<void>;
  enqueuePlaylist: (playlistId: number) => Promise<void>;
  shuffle: () => void;
  next: () => Track | undefined;
};

const useQueueStore = create<QueueState>()((set) => ({
  ...initialState,
  enqueue: async (track) => {},
  enqueuePlaylist: async (playlistId) =>
    set({ tracks: await GetPlaylistTracks(playlistId) }),
  shuffle: () => {
    set((s) => {
      for (let i = s.tracks.length - 1; i >= 1; i--) {
        const j = Math.floor(Math.random() * (i + 1));
        [s.tracks[i], s.tracks[j]] = [s.tracks[j], s.tracks[i]];
      }
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
