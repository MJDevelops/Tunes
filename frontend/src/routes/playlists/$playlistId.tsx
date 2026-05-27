import { GetPlaylistTracks } from "@bindings/internal/pkg/services/audioservice";
import { useQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";

const PlaylistPage = () => {
  const { playlistId } = Route.useParams();
  const { isPending, data } = useQuery({
    queryKey: ["playlist", playlistId],
    queryFn: async () => {
      if (isNaN(playlistId)) return null;
      return await GetPlaylistTracks(playlistId);
    },
  });
  return <div>Playlist Page</div>;
};

export const Route = createFileRoute("/playlists/$playlistId")({
  component: PlaylistPage,
  params: {
    parse: ({ playlistId }) => ({ playlistId: Number.parseInt(playlistId) }),
  },
});
