import { useCallback, useMemo, useState } from "react";
import { useTranslation } from "react-i18next";
import * as Dialog from "@radix-ui/react-dialog";

import { Button } from "@/components/ui/button";
import { assetsForNode } from "@/lib/assets";
import {
  FormTaskDetails,
  FormTaskDetailsProps,
  FormTaskContext,
} from "@/features/nodes/details/variants/FormTaskDetails";
import { t as pbT } from "@/lib/playbook";

const scrollToTemplates = () =>
  document
    .getElementById("templates-section")
    ?.scrollIntoView({ behavior: "smooth", block: "start" });

export function E1ApplyOmidDetails({
  node,
  ...rest
}: Omit<FormTaskDetailsProps, "renderActions">) {
  const { t: T, i18n } = useTranslation("common");
  const [open, setOpen] = useState(false);
  const assets = useMemo(() => assetsForNode(node), [node]);
  const lang = (i18n.language as "ru" | "kz" | "en") || "ru";

  // Extract instructions from playbook (cast node to any for accessing screen property)
  const primaryBtn = Array.isArray((node as any)?.screen?.buttons)
    ? (node as any).screen.buttons[0]
    : undefined;
  const instructionsRaw = primaryBtn?.instructions?.text as
    | string[]
    | Record<string, string[]>
    | undefined;
  const instructions: string[] = Array.isArray(instructionsRaw)
    ? instructionsRaw
    : Array.isArray((instructionsRaw as any)?.[lang])
    ? (instructionsRaw as any)[lang]
    : [];
  const accordionLabel = pbT(
    primaryBtn?.label as any,
    pbT(
      {
        ru: "Инструкция по прохождению",
        kz: "Өту бойынша нұсқаулық",
        en: "How to complete",
      },
      "Инструкция по прохождению"
    )
  );

  const openTemplate = useCallback(() => {
    const preferred =
      assets.find(
        (asset) =>
          asset.id.toLowerCase().includes("omid") &&
          asset.id.toLowerCase().includes(`_${lang}`)
      ) ||
      assets.find((asset) => asset.id.toLowerCase().includes("omid")) ||
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
                openTemplate();
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
    },
    [T, open, openTemplate]
  );

  return (
    <div className="space-y-4">
      {/* Inline guidance (always visible) */}
      {instructions.length > 0 && (
        <div className="space-y-3 p-5 sm:p-6 rounded-xl bg-muted/30 border-l-4 border-primary/40">
          <div className="text-sm font-semibold text-primary flex items-center gap-2">
            <svg
              className="w-4 h-4"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
              />
            </svg>
            {accordionLabel}
          </div>
          <ul className="list-disc pl-5 space-y-2 text-sm text-muted-foreground">
            {instructions.map((line: string, idx: number) => (
              <li key={idx} className="leading-relaxed">
                {line}
              </li>
            ))}
          </ul>
        </div>
      )}
      <FormTaskDetails node={node} {...rest} renderActions={renderActions} />
    </div>
  );
}

export default E1ApplyOmidDetails;

