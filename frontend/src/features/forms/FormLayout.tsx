import type { NodeVM } from "@/lib/playbook";

export function FormLayout({
  left,
  right,
}: {
  left: React.ReactNode;
  right?: React.ReactNode;
}) {
  return (
    <div className="grid grid-cols-1 lg:grid-cols-5 gap-4 h-full">
      <div className="lg:col-span-3 min-h-0 overflow-auto pr-1">{left}</div>
      {right ? <div className="lg:col-span-2 border-l pl-4 overflow-auto">{right}</div> : null}
    </div>
  );
}

