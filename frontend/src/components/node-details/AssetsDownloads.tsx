// components/node-details/AssetsDownloads.tsx
import { assetsForNode, PublicAsset } from "@/lib/assets";
import { NodeVM, t } from "@/lib/playbook";
import { Separator } from "@/components/ui/separator";
import { Button } from "@/components/ui/button";

export function AssetsDownloads({
  node,
  locale = "ru",
  heading = {
    ru: "Шаблоны",
    kz: "Үлгілер",
    en: "Templates",
  },
}: {
  node: NodeVM;
  locale?: "ru" | "kz" | "en";
  heading?: Record<string, string>;
}) {
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
      <div className="font-semibold">{t(heading, "Templates")}</div>
      <div className="space-y-2">
        {order.map((g) => {
          const items = groups[g];
          // prefer current locale first
          const sorted = items.slice().sort((a, b) => {
            const la = a.id.includes(`_${locale}`) ? 0 : 1;
            const lb = b.id.includes(`_${locale}`) ? 0 : 1;
            if (la !== lb) return la - lb;
            return a.id.localeCompare(b.id);
          });
          return (
            <div key={g} className="flex flex-wrap items-center gap-2">
              <div className="min-w-40 text-sm font-medium capitalize">
                {g.replace("app", "Appendix ")}
              </div>
              {sorted.map((a) => (
                <a
                  key={a.id}
                  href={`/${a.storage.key}`}
                  target="_blank"
                  rel="noopener noreferrer"
                  download
                >
                  <Button variant="outline" size="sm">
                    {a.title[locale] || t(a.title, a.id)}
                  </Button>
                </a>
              ))}
            </div>
          );
        })}
      </div>
    </div>
  );
}

