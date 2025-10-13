import { Suspense, use } from "react";
import { PendingDownloads } from "~/wailsjs/go/main/App";
import { ytdlp } from "~/wailsjs/go/models";
import { Skeleton } from "@/components/ui/skeleton";

function Downloads({ promise }: { promise: Promise<ytdlp.Download[] | null> }) {
  const downloads = use(promise);

  return (
    <div>
      {downloads ? (
        downloads.map((d) => (
          <span>
            {d.ID} {d.Url}
          </span>
        ))
      ) : (
        <span>No Downloads</span>
      )}
    </div>
  );
}

export default function DownloadOverview() {
  return (
    <Suspense fallback={<Skeleton className="h-full w-full m-2" />}>
      <Downloads promise={PendingDownloads()} />
    </Suspense>
  );
}
