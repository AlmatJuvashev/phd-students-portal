import { useCallback, useMemo, useState } from "react";
import { useTranslation } from "react-i18next";
import * as Dialog from "@radix-ui/react-dialog";

import { Button } from "@/components/ui/button";
import { allAssets } from "@/lib/assets";
import { t as pickLocale } from "@/lib/playbook";
import {
  FormTaskDetails,
  FormTaskDetailsProps,
  FormTaskContext,
} from "@/features/nodes/details/variants/FormTaskDetails";

const scrollToTemplates = () =>
  document
    .getElementById("templates-section")
    ?.scrollIntoView({ behavior: "smooth", block: "start" });

export function NkPackageDetails({
  node,
  ...rest
}: Omit<FormTaskDetailsProps, "renderActions">) {
  const { t: T, i18n } = useTranslation("common");
  const [confirm, setConfirm] = useState(false);
  const templateCandidates = useMemo(() => {
    const explicit = (node.requirements as any)?.templates as
      | string[]
      | undefined;
    const pool = allAssets();
    const fallback = pool.filter((a) =>
      a.id.toLowerCase().includes("lcb_request")
    );
    if (!explicit?.length) return fallback;
    return explicit
      .map((id) => pool.find((asset) => asset.id === id))
      .filter(Boolean) as typeof pool;
  }, [node]);

  const lang = (i18n.language as "ru" | "kz" | "en") || "ru";

  const pickLetterTemplate = useCallback(() => {
    const preferred =
      templateCandidates.find((asset) =>
        asset.id.toLowerCase().includes(`_${lang}`)
      ) || templateCandidates[0];
    if (preferred?.storage?.key) {
      window.open(`/${preferred.storage.key}`, "_blank", "noopener,noreferrer");
    }
  }, [templateCandidates, lang]);

  const renderActions = useCallback(
    (ctx: FormTaskContext) => {
      if (!ctx.canEdit) return null;
      const ui = (ctx.node.requirements as any)?.ui_hints;
      const btnAll = ui?.buttons?.find((b: any) => b.key === "all_ready");
      const btnNeed = ui?.buttons?.find(
        (b: any) => b.key === "need_lcb_letter"
      );
      const allLabel = pickLocale(btnAll?.label, T("forms.save_submit"));
      const needLabel = pickLocale(btnNeed?.label, T("forms.save_draft"));
      const ready =
        !!ctx.values["chk_thesis_unbound"] &&
        !!ctx.values["chk_advisor_reviews"] &&
        !!ctx.values["chk_pubs_app7"] &&
        !!ctx.values["chk_sc_extract"] &&
        !!ctx.values["chk_lcb_defense"];

      return (
        <>
          <div className="font-semibold">
            {T("forms.nk_required_package_title")}
          </div>
          <div className="text-sm text-muted-foreground">
            {pickLocale((ctx.node as any).description, "")}
          </div>
          <div className="flex gap-2">
            <Button
              onClick={() => setConfirm(true)}
              disabled={ctx.disabled || !ready}
            >
              {allLabel}
            </Button>
            <Button
              variant="secondary"
              disabled={ctx.disabled}
              onClick={() => {
                pickLetterTemplate();
                scrollToTemplates();
                ctx.saveDraft();
              }}
            >
              {needLabel}
            </Button>
          </div>
          <Dialog.Root open={confirm} onOpenChange={setConfirm}>
            <Dialog.Portal>
              <Dialog.Overlay className="fixed inset-0 bg-black/50 z-[70]" />
              <Dialog.Content className="fixed z-[70] left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2 bg-white rounded-lg p-6 w-full max-w-md shadow-lg outline-none">
                <Dialog.Title className="text-lg font-semibold mb-4">
                  {T("forms.confirm_submission", "Подтверждение")}
                </Dialog.Title>
                <div className="mb-4 text-sm text-muted-foreground whitespace-pre-line">
                  {T("forms.nk_confirm_info")}
                </div>
                <div className="flex gap-2 justify-end">
                  <Button variant="outline" onClick={() => setConfirm(false)}>
                    {T("common.cancel")}
                  </Button>
                  <Button
                    onClick={() => {
                      setConfirm(false);
                      const nextOnComplete =
                        ctx.node.outcomes?.find((o) => o.value === "complete")
                          ?.next?.[0] ||
                        (Array.isArray(ctx.node.next)
                          ? ctx.node.next[0]
                          : undefined);
                      ctx.submit(
                        nextOnComplete
                          ? { __nextOverride: nextOnComplete }
                          : undefined
                      );
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
    [T, confirm, pickLetterTemplate]
  );

  return (
    <FormTaskDetails node={node} {...rest} renderActions={renderActions} />
  );
}

export default NkPackageDetails;
