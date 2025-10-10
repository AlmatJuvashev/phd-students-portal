import React from "react";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Modal } from "@/components/ui/modal";
import { getAssetUrl } from "@/lib/assets";
import { t, safeText } from "@/lib/playbook";
import { Accordion, AccordionItem } from "@/components/ui/accordion";
import type { NodeVM } from "@/lib/playbook";
import i18n from "i18next";

type ConfirmTaskDetailsProps = {
  node: NodeVM | any;
  onComplete?: () => void;
};

const ConfirmTaskDetails: React.FC<ConfirmTaskDetailsProps> = ({
  node,
  onComplete,
}) => {
  const [isCompleted, setCompleted] = React.useState(
    node?.state === "done" || node?.status === "completed"
  );
  const [confirmOpen, setConfirmOpen] = React.useState(false);

  // Localized question (string or i18n map)
  const question: string | undefined = safeText(
    node?.screen?.question as any,
    ""
  );

  // Primary button (index 0) contains instructions
  const primaryBtn = Array.isArray(node?.screen?.buttons)
    ? node.screen.buttons[0]
    : undefined;
  // Normalize instruction text: could be array of strings or a locale->string[] map
  const instructionsRaw = primaryBtn?.instructions?.text as
    | string[]
    | Record<string, string[]>
    | undefined;
  const currentLang = (i18n?.language as "ru" | "kz" | "en") || "ru";
  const instructions: string[] = Array.isArray(instructionsRaw)
    ? instructionsRaw
    : Array.isArray((instructionsRaw as any)?.[currentLang])
    ? (instructionsRaw as any)[currentLang]
    : [];
  const download = (primaryBtn?.instructions?.download || undefined) as
    | { label?: string; asset_id?: string; asset_path?: string }
    | undefined;
  const accordionLabel = t(
    primaryBtn?.label as any,
    t(
      {
        ru: "Инструкция по прохождению",
        kz: "Өту бойынша нұсқаулық",
        en: "How to complete",
      },
      "Инструкция по прохождению"
    )
  );

  // Confirmation button (index 1)
  const confirmBtn = Array.isArray(node?.screen?.buttons)
    ? node.screen.buttons[1]
    : undefined;
  const confirmLabel: string = t(
    (confirmBtn?.label as any) || {},
    t({ ru: "Подтвердить", kz: "Растау", en: "Confirm" }, "Подтвердить")
  );
  const confirmText: string = safeText(
    (confirmBtn?.confirmation_text as any) || {},
    safeText(
      {
        ru: "Подтвердить выполнение шага?",
        kz: "Қадамды орындауды растау?",
        en: "Confirm completing this step?",
      },
      "Подтвердить выполнение шага?"
    )
  );

  const completedMessage: string = safeText(
    (node?.states as any)?.completed?.message as any,
    safeText(
      {
        ru: "Шаг подтверждён.",
        kz: "Қадам расталды.",
        en: "Step confirmed.",
      },
      "Шаг подтверждён."
    )
  );

  const handleConfirm = () => {
    setCompleted(true);
    setConfirmOpen(false);
    // toast removed; keep silent confirmation
    if (onComplete) onComplete();
  };

  return (
    <Card className="p-4">
      <CardContent className="space-y-4">
        {question && <p className="text-lg font-medium">{question}</p>}

        {!isCompleted && (
          <>
            <Accordion>
              <AccordionItem header={accordionLabel}>
                {Array.isArray(instructions) && instructions.length > 0 && (
                  <ul className="list-disc pl-5 space-y-1 text-sm text-muted-foreground">
                    {instructions.map((line: string, idx: number) => (
                      <li key={idx}>{line}</li>
                    ))}
                  </ul>
                )}

                {(() => {
                  if (!download) return null;
                  const resolved = download.asset_id
                    ? getAssetUrl(download.asset_id)
                    : undefined;
                  const href =
                    resolved && resolved !== "#"
                      ? resolved
                      : download.asset_path;
                  if (!href) return null;
                  return (
                    <div className="mt-3">
                      <Button asChild variant="secondary">
                        <a
                          href={href}
                          download
                          target="_blank"
                          rel="noopener noreferrer"
                        >
                          {download.label ||
                            t(
                              {
                                ru: "Скачать шаблон",
                                kz: "Үлгіні жүктеу",
                                en: "Download template",
                              },
                              "Скачать шаблон"
                            )}
                        </a>
                      </Button>
                    </div>
                  );
                })()}
              </AccordionItem>
            </Accordion>

            <div className="pt-2">
              <Button
                variant="default"
                className="mt-2"
                onClick={() => setConfirmOpen(true)}
              >
                {confirmLabel}
              </Button>
            </div>

            <Modal open={confirmOpen} onClose={() => setConfirmOpen(false)}>
              <div className="space-y-4">
                <div className="text-base font-medium">{confirmText}</div>
                <div className="flex justify-end gap-2">
                  <Button variant="ghost" onClick={() => setConfirmOpen(false)}>
                    {t(
                      { ru: "Отмена", kz: "Болдырмау", en: "Cancel" },
                      "Отмена"
                    )}
                  </Button>
                  <Button onClick={handleConfirm}>
                    {t(
                      {
                        ru: "Да, подтвердить",
                        kz: "Иә, растау",
                        en: "Yes, confirm",
                      },
                      "Да, подтвердить"
                    )}
                  </Button>
                </div>
              </div>
            </Modal>
          </>
        )}

        {isCompleted && (
          <div className="rounded-2xl bg-emerald-50 p-4 text-emerald-700">
            ✅ {completedMessage}
          </div>
        )}
      </CardContent>
    </Card>
  );
};

export default ConfirmTaskDetails;
