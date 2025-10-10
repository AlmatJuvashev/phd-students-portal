import { useTranslation } from "react-i18next";
import GuideCard from "@/components/ui/guide-card";
import type { NodeVM } from "@/lib/playbook";

export function useGuideForNode(node: NodeVM) {
  const { t } = useTranslation("guides");

  const mapping: Record<string, { key: string }> = {
    S1_text_ready: { key: "text_ready" },
    S1_antiplag: { key: "antiplag" },
    IV3_publication_certificate_ncste: { key: "ncste_pubcert" },
  };

  const entry = mapping[node.id];
  if (!entry) return null;

  const key = entry.key;
  const title = t(`${key}.title`) as string;
  const items = (t(`${key}.items`, { returnObjects: true }) as string[]) || [];
  const href = t(`${key}.linkUrl`) as string;
  const ctaLabel = t(`${key}.linkLabel`) as string;

  return function RenderGuide() {
    return (
      <GuideCard title={title} items={items} href={href || undefined} ctaLabel={ctaLabel} />
    );
  };
}

export default useGuideForNode;

