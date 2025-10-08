// components/node-details/AssetsDownloads.tsx
import { assetsForNode, PublicAsset, allAssets } from "@/lib/assets";
import { NodeVM, t } from "@/lib/playbook";
import { Separator } from "@/components/ui/separator";
import { Button } from "@/components/ui/button";
import { useTranslation } from "react-i18next";

export function AssetsDownloads({ node }: { node: NodeVM }) {
  const { i18n, t: T } = useTranslation("common");
  const locale = (i18n.language as "ru" | "kz" | "en") || "ru";
  // Prefer explicit templates listed on the node, fallback to heuristic mapping
  let assets: PublicAsset[] = [];
  const explicit = (node.requirements as any)?.templates as
    | string[]
    | undefined;
  if (explicit?.length) {
    const pool = allAssets();
    assets = explicit
      .map((id) => pool.find((a) => a.id === id))
      .filter(Boolean) as PublicAsset[];
  } else {
    assets = assetsForNode(node);
  }
  if (!assets.length) return null;

  // group by logical base template (e.g., app7, omid, etc.) so we show only one button per locale
  const groups: Record<string, PublicAsset[]> = {};
  for (const a of assets) {
    const id = a.id.toLowerCase();
    // appN (Appendix templates)
    const m = id.match(/(app\d+)/i);
    let key = m ? m[1].toLowerCase() : id;
    // OMiD application and similar localized assets
    if (!m && id.includes("omid")) key = "omid";
    // Fallback: strip locale suffix like _ru/_kz/_en and trailing segments
    if (!m && !id.includes("omid")) {
      key = id.replace(/_(ru|kz|en)(_.+)?$/, "");
    }
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
