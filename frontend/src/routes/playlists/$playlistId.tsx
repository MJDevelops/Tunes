import { GetPlaylist } from "@bindings/internal/pkg/services/audioservice";
import { useQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";

const PlaylistPage = () => {
  const { playlistId } = Route.useParams();
  const { isPending, data } = useQuery({
    queryKey: ["playlist", playlistId],
    queryFn: async () => {
      if (isNaN(playlistId)) return null;
      return await GetPlaylist(playlistId);
    },
  });

  return isPending ? (
    <div>Loading...</div>
  ) : (
    <div>
      <h1>{data?.Title}</h1>
      {data?.Tracks.map((t) => (
        <div key={t.ID}>
          <h1>{t.Title}</h1>
          <p>{t.Artists.map((a) => a?.Name).join(", ")}</p>
        </div>
      ))}
    </div>
  );
};

export const Route = createFileRoute("/playlists/$playlistId")({
  component: PlaylistPage,
  params: {
    parse: ({ playlistId }) => ({ playlistId: Number.parseInt(playlistId) }),
  },
});
