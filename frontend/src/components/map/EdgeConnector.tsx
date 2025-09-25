// components/map/EdgeConnector.tsx
export function EdgeConnector({
  x1,
  y1,
  x2,
  y2,
  kind = "default",
}: {
  x1: number;
  y1: number;
  x2: number;
  y2: number;
  kind?: "default" | "conditional" | "outcome";
}) {
  const midX = (x1 + x2) / 2;
  const path = `M ${x1} ${y1} C ${midX} ${y1}, ${midX} ${y2}, ${x2} ${y2}`;
  const dash =
    kind === "conditional" ? "4,4" : kind === "outcome" ? "0.1,6" : "0";
  return (
    <path
      d={path}
      fill="none"
      stroke="currentColor"
      strokeWidth={2}
      strokeDasharray={dash}
      className="text-muted-foreground/70"
    />
  );
}
