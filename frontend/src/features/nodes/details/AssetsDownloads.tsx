// components/node-details/AssetsDownloads.tsx
import { PublicAsset } from "@/lib/assets";
import { NodeVM, t } from "@/lib/playbook";
import { Separator } from "@/components/ui/separator";
import { Button } from "@/components/ui/button";
import { useMemo } from "react";
import { useTranslation } from "react-i18next";
import { useTemplatesForNode } from "@/features/nodes/useTemplatesForNode";

export function AssetsDownloads({ node }: { node: NodeVM }) {
  const { i18n, t: T } = useTranslation("common");
  const locale = (i18n.language as "ru" | "kz" | "en") || "ru";
  const assets: PublicAsset[] = useTemplatesForNode(node);
  const { order, groups } = useMemo(() => {
    const grouped: Record<string, PublicAsset[]> = {};
    for (const asset of assets) {
      const id = asset.id.toLowerCase();
      const match = id.match(/(app\d+)/i);
      let key = match ? match[1].toLowerCase() : id;
      if (!match && id.includes("omid")) key = "omid";
      if (!match && !id.includes("omid")) {
        key = id.replace(/_(ru|kz|en)(_.+)?$/, "");
      }
      grouped[key] = grouped[key] || [];
      grouped[key].push(asset);
    }
    return {
      order: Object.keys(grouped).sort(),
      groups: grouped,
    };
  }, [assets]);

  if (!order.length) return null;
  return (
    <div className="space-y-4 sticky top-4">
      <Separator className="my-4" />
      <div className="space-y-3">
        <h3
          id="templates-section"
          className="font-bold text-base sm:text-lg flex items-center gap-2 text-foreground"
        >
          <svg
            className="w-5 h-5 text-primary"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M7 21h10a2 2 0 002-2V9.414a1 1 0 00-.293-.707l-5.414-5.414A1 1 0 0012.586 3H7a2 2 0 00-2 2v14a2 2 0 002 2z"
            />
          </svg>
          {T("templates.heading")}
        </h3>
        <div className="space-y-2">
          {order.map((g, idx) => {
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
              <div
                key={g}
                className="animate-in fade-in slide-in-from-right-2 duration-300"
                style={{ animationDelay: `${idx * 50}ms` }}
              >
                <a
                  href={`/${preferred.storage.key}`}
                  target="_blank"
                  rel="noopener noreferrer"
                  download
                  className="block"
                >
                  <Button
                    variant="outline"
                    size="sm"
                    className="w-full justify-start gap-2 hover:bg-primary/5 hover:border-primary/40 hover:shadow-sm transition-all duration-200 group"
                  >
                    <svg
                      className="w-4 h-4 text-muted-foreground group-hover:text-primary transition-colors"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                    >
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth={2}
                        d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
                      />
                    </svg>
                    <span className="text-xs sm:text-sm truncate">{label}</span>
                  </Button>
                </a>
              </div>
            );
          })}
        </div>
      </div>
    </div>
  );
}
