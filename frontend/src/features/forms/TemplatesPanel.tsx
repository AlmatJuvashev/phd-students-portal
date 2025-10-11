import { AssetsDownloads } from "@/components/node-details/AssetsDownloads";
import type { NodeVM } from "@/lib/playbook";
import clsx from "clsx";

export function TemplatesPanel({
  node,
  className,
}: {
  node: NodeVM;
  className?: string;
}) {
  return (
    <div className={clsx("lg:col-span-2 border-l pl-4", className)}>
      <AssetsDownloads node={node} />
    </div>
  );
}
