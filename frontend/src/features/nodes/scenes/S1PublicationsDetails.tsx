import { useCallback, useMemo } from "react";
import { useTranslation } from "react-i18next";

import { Button } from "@/components/ui/button";
import { assetsForNode } from "@/lib/assets";
import {
  FormTaskDetails,
  FormTaskDetailsProps,
  FormTaskContext,
} from "@/features/nodes/details/variants/FormTaskDetails";

const scrollToTemplates = () =>
  document
    .getElementById("templates-section")
    ?.scrollIntoView({ behavior: "smooth", block: "start" });

export function S1PublicationsDetails({ node, ...rest }: Omit<FormTaskDetailsProps, "renderActions">) {
  const { t: T, i18n } = useTranslation("common");
  const assets = useMemo(() => assetsForNode(node), [node]);
  const lang = (i18n.language as "ru" | "kz" | "en") || "ru";

  const openTemplate = useCallback(() => {
    const preferred =
      assets.find(
        (asset) =>
          asset.id.toLowerCase().includes("app7") &&
          asset.id.toLowerCase().includes(`_${lang}`)
      ) ||
      assets.find((asset) => asset.id.toLowerCase().includes("app7")) ||
      assets[0];
    if (preferred?.storage?.key) {
      window.open(`/${preferred.storage.key}`, "_blank", "noopener,noreferrer");
    }
  }, [assets, lang]);

  const renderActions = useCallback(
    (ctx: FormTaskContext) => {
      if (!ctx.canEdit) return null;
      return (
        <>
          <div className="text-sm font-medium">
            {T("forms.app7_prompt")}{" "}
            <span className="text-destructive">*</span>
          </div>
          <div className="flex gap-2">
            <Button onClick={() => ctx.submit()} disabled={ctx.disabled}>
              {T("forms.yes")}
            </Button>
            <Button
              variant="secondary"
              disabled={ctx.disabled}
              onClick={() => {
                openTemplate();
                scrollToTemplates();
                ctx.saveDraft();
              }}
            >
              {T("forms.no")}
            </Button>
          </div>
        </>
      );
    },
    [T, openTemplate]
  );

  return <FormTaskDetails node={node} {...rest} renderActions={renderActions} />;
}

export default S1PublicationsDetails;
