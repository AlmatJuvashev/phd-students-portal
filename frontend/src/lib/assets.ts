// lib/assets.ts
import type { NodeVM } from "@/lib/playbook";
// Allow JSON import via Vite bundler
// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore
import assetsList from "@/playbooks/assets_list.json";

export type PublicAsset = {
  id: string;
  kind: "template" | string;
  title: Record<string, string>;
  mime: string;
  storage: { provider: "public" | string; key: string };
  version?: string;
};

export function allAssets(): PublicAsset[] {
  const arr = (assetsList?.assets as PublicAsset[]) || [];
  return arr;
}

// Heuristic mapping from node → related templates (Appendix 5..9)
const rules: Array<{ match: (s: string, id: string) => boolean; tag: string }> = [
  {
    match: (s, id) => /publication|публикац|жариялан/.test(s) || id.includes("pub"),
    tag: "app7",
  },
  {
    match: (s) => /advisor|консультант|пікір/.test(s),
    tag: "app6",
  },
  {
    match: (s) => /application|заявлен|өтініш/.test(s),
    tag: "app5",
  },
  {
    match: (s) => /abstract|аннотац/.test(s),
    tag: "app8",
  },
  {
    match: (s) => /conclusion|заключ|қорытынды/.test(s),
    tag: "app9",
  },
];

export function assetsForNode(node: NodeVM): PublicAsset[] {
  const titles = node.title ? Object.values(node.title).join(" ") : "";
  const s = `${titles}`.toLowerCase();
  const id = node.id.toLowerCase();
  const tag = rules.find((r) => r.match(s, id))?.tag;
  if (!tag) return [];
  return allAssets().filter((a) => a.id.toLowerCase().includes(tag));
}

