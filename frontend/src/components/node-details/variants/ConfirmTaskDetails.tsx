import React from "react";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Modal } from "@/components/ui/modal";
import { getAssetUrl } from "@/lib/assets";
import { t, safeText } from "@/lib/playbook";
import type { NodeVM } from "@/lib/playbook";
import i18n from "i18next";

type ConfirmTaskDetailsProps = {
  node: NodeVM | any;
  onComplete?: () => void;
  onReset?: () => void;
};

const ConfirmTaskDetails: React.FC<ConfirmTaskDetailsProps> = ({
  node,
  onComplete,
  onReset,
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

  const handleReset = () => {
    setCompleted(false);
    if (onReset) onReset();
  };

  return (
    <Card className="bg-gradient-to-br from-card to-card/50">
      <CardContent className="space-y-5">
        {question && (
          <p className="text-lg sm:text-xl font-semibold text-foreground leading-relaxed">
            {question}
          </p>
        )}

        {!isCompleted && (
          <>
            {/* Inline guidance (no collapsible) */}
            <div className="space-y-3 p-4 rounded-xl bg-muted/30 border-l-4 border-primary/40">
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
              {Array.isArray(instructions) && instructions.length > 0 && (
                <ul className="list-disc pl-5 space-y-2 text-sm text-muted-foreground">
                  {instructions.map((line: string, idx: number) => (
                    <li key={idx} className="leading-relaxed">
                      {line}
                    </li>
                  ))}
                </ul>
              )}
              {(() => {
                if (!download) return null;
                const resolved = download.asset_id
                  ? getAssetUrl(download.asset_id)
                  : undefined;
                const href =
                  resolved && resolved !== "#" ? resolved : download.asset_path;
                if (!href) return null;
                return (
                  <div className="mt-3">
                    <Button
                      asChild
                      variant="secondary"
                      size="sm"
                      className="gap-2"
                    >
                      <a
                        href={href}
                        download
                        target="_blank"
                        rel="noopener noreferrer"
                      >
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
                            d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
                          />
                        </svg>
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
            </div>

            <div className="pt-2">
              <Button
                variant="default"
                size="lg"
                className="w-full sm:w-auto"
                onClick={() => setConfirmOpen(true)}
              >
                {confirmLabel}
              </Button>
            </div>

            <Modal open={confirmOpen} onClose={() => setConfirmOpen(false)}>
              <div className="space-y-5 p-2">
                <div className="text-center">
                  <div className="mx-auto w-14 h-14 bg-primary/10 rounded-full flex items-center justify-center mb-3">
                    <svg
                      className="w-7 h-7 text-primary"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                    >
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth={2}
                        d="M8.228 9c.549-1.165 2.03-2 3.772-2 2.21 0 4 1.343 4 3 0 1.4-1.278 2.575-3.006 2.907-.542.104-.994.54-.994 1.093m0 3h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                      />
                    </svg>
                  </div>
                  <div className="text-base sm:text-lg font-semibold text-foreground">
                    {confirmText}
                  </div>
                </div>
                <div className="flex flex-col-reverse sm:flex-row justify-center gap-2 sm:gap-3">
                  <Button
                    variant="outline"
                    onClick={() => setConfirmOpen(false)}
                    className="sm:min-w-[120px]"
                  >
                    {t(
                      { ru: "Отмена", kz: "Болдырмау", en: "Cancel" },
                      "Отмена"
                    )}
                  </Button>
                  <Button onClick={handleConfirm} className="sm:min-w-[120px]">
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
          <div className="rounded-2xl bg-gradient-to-br from-green-50 to-green-100/50 dark:from-green-900/20 dark:to-green-800/10 p-5 border-2 border-green-500/30 animate-in fade-in slide-in-from-bottom-2 duration-500">
            <div className="flex items-center gap-3">
              <div className="flex-shrink-0 w-10 h-10 bg-green-500 rounded-full flex items-center justify-center">
                <svg
                  className="w-6 h-6 text-white"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2.5}
                    d="M5 13l4 4L19 7"
                  />
                </svg>
              </div>
              <span className="text-green-700 dark:text-green-300 font-semibold flex-1">
                {completedMessage}
              </span>
              <Button variant="outline" size="sm" onClick={handleReset}>
                {t(
                  { ru: "Сбросить подтверждение", kz: "Растауын жою", en: "Reset confirmation" },
                  "Сбросить подтверждение"
                )}
              </Button>
            </div>
          </div>
        )}
      </CardContent>
    </Card>
  );
};

export default ConfirmTaskDetails;
