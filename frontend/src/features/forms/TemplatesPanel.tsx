import { AssetsDownloads } from "@/components/node-details/AssetsDownloads";
import type { NodeVM } from "@/lib/playbook";

export function TemplatesPanel({ node }: { node: NodeVM }) {
  return (
    <div className="lg:col-span-2 border-l pl-4 overflow-auto">
      <AssetsDownloads node={node} />
    </div>
  );
}

