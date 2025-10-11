import { useState } from "react";
import { useTranslation } from "react-i18next";
import * as Dialog from "@radix-ui/react-dialog";

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

export function E1ApplyOmidDetails(
  props: Omit<FormTaskDetailsProps, "renderActions">
) {
  const { t: T, i18n } = useTranslation("common");
  const [open, setOpen] = useState(false);

  const renderActions = (ctx: FormTaskContext) => {
    if (!ctx.canEdit) return null;
    const lang = (i18n.language as "ru" | "kz" | "en") || "ru";
    return (
      <>
        <div className="text-sm font-medium">
          {T("forms.omid_prompt")}{" "}
          <span className="text-destructive">*</span>
        </div>
        <div className="flex gap-2">
          <Button onClick={() => setOpen(true)} disabled={ctx.disabled}>
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
                    a.id.toLowerCase().includes("omid") &&
                    a.id.toLowerCase().includes(`_${lang}`)
                ) ||
                assets.find((a) => a.id.toLowerCase().includes("omid")) ||
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
        <Dialog.Root open={open} onOpenChange={setOpen}>
          <Dialog.Portal>
            <Dialog.Overlay className="fixed inset-0 bg-black/50 z-[70]" />
            <Dialog.Content className="fixed z-[70] left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2 bg-white rounded-lg p-6 w-full max-w-md shadow-lg outline-none">
              <div className="mb-4 text-sm text-muted-foreground whitespace-pre-line">
                {T("forms.omid_info_after_yes")}
              </div>
              <div className="flex gap-2 justify-end">
                <Button variant="outline" onClick={() => setOpen(false)}>
                  {T("common.cancel")}
                </Button>
                <Button
                  onClick={() => {
                    setOpen(false);
                    ctx.submit();
                  }}
                  disabled={ctx.disabled}
                >
                  {T("forms.proceed_next")}
                </Button>
              </div>
            </Dialog.Content>
          </Dialog.Portal>
        </Dialog.Root>
      </>
    );
  };

  return <FormTaskDetails {...props} renderActions={renderActions} />;
}

export default E1ApplyOmidDetails;
