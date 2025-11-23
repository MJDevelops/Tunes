import { createFileRoute } from "@tanstack/react-router";

function Index() {
  return (
    <div>
      <h1>Home</h1>
    </div>
  );
}

export const Route = createFileRoute("/")({
  component: Index,
});
