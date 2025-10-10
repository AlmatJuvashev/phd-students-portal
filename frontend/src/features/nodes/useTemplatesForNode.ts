import { allAssets, assetsForNode, PublicAsset } from "@/lib/assets";
import type { NodeVM } from "@/lib/playbook";
import { useMemo } from "react";
import { useTranslation } from "react-i18next";

export function useTemplatesForNode(node: NodeVM) {
  const { i18n } = useTranslation();
  const locale = (i18n.language as "ru" | "kz" | "en") || "ru";

  return useMemo(() => {
    const explicit = (node.requirements as any)?.templates as string[] | undefined;
    let assets: PublicAsset[] = [];
    if (explicit?.length) {
      const pool = allAssets();
      assets = explicit
        .map((id) => pool.find((a) => a.id === id))
        .filter(Boolean) as PublicAsset[];
    } else {
      assets = assetsForNode(node);
    }
    if (!assets.length) return [] as PublicAsset[];

    // group and pick single per logical template by locale preference
    const groups: Record<string, PublicAsset[]> = {};
    for (const a of assets) {
      const id = a.id.toLowerCase();
      const m = id.match(/(app\d+)/i);
      let key = m ? m[1].toLowerCase() : id;
      if (!m && id.includes("omid")) key = "omid";
      if (!m && !id.includes("omid")) key = id.replace(/_(ru|kz|en)(_.+)?$/, "");
      groups[key] = groups[key] || [];
      groups[key].push(a);
    }
    const pickOne = (arr: PublicAsset[]) =>
      arr.find((a) => a.id.toLowerCase().includes(`_${locale}`)) ||
      arr.find((a) => a.title?.[locale as any]) ||
      arr[0];
    const result: PublicAsset[] = [];
    Object.keys(groups)
      .sort()
      .forEach((k) => {
        const chosen = pickOne(groups[k]);
        if (chosen) result.push(chosen);
      });
    return result;
  }, [node, i18n.language]);
}

