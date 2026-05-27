import { createFileRoute, Link } from "@tanstack/react-router";
import { GetPlaylists } from "@bindings/internal/pkg/services/audioservice";
import { useQuery } from "@tanstack/react-query";

const Playlists = () => {
  const { isPending, data } = useQuery({
    queryKey: ["playlists"],
    queryFn: async () => await GetPlaylists(),
  });
  return isPending ? (
    <div>Loading...</div>
  ) : (
    <div>
      {data?.map((p) => (
        <Link
          to="/playlists/$playlistId"
          params={{ playlistId: p.ID }}
          key={p.ID}
        >
          {p.Title}
        </Link>
      ))}
    </div>
  );
};

export const Route = createFileRoute("/playlists/")({
  component: Playlists,
});
