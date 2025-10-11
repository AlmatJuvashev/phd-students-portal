import { useTranslation } from "react-i18next";

import { Button } from "@/components/ui/button";
import { assetsForNode } from "@/lib/assets";
import { FormTaskDetails, FormTaskDetailsProps, FormTaskContext } from "@/features/nodes/details/variants/FormTaskDetails";

const scrollToTemplates = () =>
  document
    .getElementById("templates-section")
    ?.scrollIntoView({ behavior: "smooth", block: "start" });

export function S1PublicationsDetails(props: Omit<FormTaskDetailsProps, "renderActions">) {
  const { t: T, i18n } = useTranslation("common");

  const renderActions = (ctx: FormTaskContext) => {
    if (!ctx.canEdit) return null;
    const lang = (i18n.language as "ru" | "kz" | "en") || "ru";
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
              const assets = assetsForNode(ctx.node);
              const preferred =
                assets.find(
                  (a) =>
                    a.id.toLowerCase().includes("app7") &&
                    a.id.toLowerCase().includes(`_${lang}`)
                ) ||
                assets.find((a) => a.id.toLowerCase().includes("app7")) ||
                assets[0];
              if (preferred?.storage?.key) {
                window.open(`/${preferred.storage.key}`, "_blank", "noopener,noreferrer");
              }
              scrollToTemplates();
              ctx.saveDraft();
            }}
          >
            {T("forms.no")}
          </Button>
        </div>
      </>
    );
  };

  return <FormTaskDetails {...props} renderActions={renderActions} />;
}

export default S1PublicationsDetails;
