// components/node-details/AssetsDownloads.tsx
import { assetsForNode, PublicAsset } from "@/lib/assets";
import { NodeVM, t } from "@/lib/playbook";
import { Separator } from "@/components/ui/separator";
import { Button } from "@/components/ui/button";
import { useTranslation } from "react-i18next";

export function AssetsDownloads({ node }: { node: NodeVM }) {
  const { i18n, t: T } = useTranslation("common");
  const locale = (i18n.language as "ru" | "kz" | "en") || "ru";
  const assets = assetsForNode(node);
  if (!assets.length) return null;

  // group by base template (appN)
  const groups: Record<string, PublicAsset[]> = {};
  for (const a of assets) {
    const m = a.id.match(/(app\d+)/i);
    const key = m ? m[1].toLowerCase() : a.id;
    groups[key] = groups[key] || [];
    groups[key].push(a);
  }

  const order = Object.keys(groups).sort();
  return (
    <div className="space-y-3">
      <Separator />
      <div id="templates-section" className="font-semibold">
        {T("templates.heading")}
      </div>
      <div className="space-y-2">
        {order.map((g) => {
          const items = groups[g];
          const preferred =
            items.find((a) => a.id.toLowerCase().includes(`_${locale}`)) ||
            items.find((a) => a.title?.[locale as any]) ||
            items[0];
          if (!preferred) return null;
          const label =
            preferred.title?.[locale as any] ||
            t(preferred.title, preferred.id);
          return (
            <div key={g} className="flex flex-wrap items-center gap-2">
              <a
                href={`/${preferred.storage.key}`}
                target="_blank"
                rel="noopener noreferrer"
                download
              >
                <Button variant="outline" size="sm">
                  {label}
                </Button>
              </a>
            </div>
          );
        })}
      </div>
    </div>
  );
}
