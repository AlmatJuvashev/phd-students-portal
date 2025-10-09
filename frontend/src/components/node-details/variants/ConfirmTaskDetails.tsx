import React from "react";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Modal } from "@/components/ui/modal";
import { useToast } from "@/components/toast";
import { getAssetUrl } from "@/lib/assets";
import { t } from "@/lib/playbook";
import { Accordion, AccordionItem } from "@/components/ui/accordion";
import type { NodeVM } from "@/lib/playbook";

type ConfirmTaskDetailsProps = {
  node: NodeVM | any;
  onComplete?: () => void;
};

const ConfirmTaskDetails: React.FC<ConfirmTaskDetailsProps> = ({
  node,
  onComplete,
}) => {
  const { push } = useToast();
  const [isCompleted, setCompleted] = React.useState(
    node?.state === "done" || node?.status === "completed"
  );
  const [confirmOpen, setConfirmOpen] = React.useState(false);

  // Helpers: localize string or array-of-strings from either raw or {ru,kz,en}
  const localizeString = (val: any): string | undefined => {
    if (val == null) return undefined;
    if (typeof val === "string") return val;
    if (typeof val === "object") return t(val as any, "");
    return undefined;
  };
  const localizeList = (val: any): string[] => {
    if (Array.isArray(val)) return val as string[];
    if (val && typeof val === "object") {
      const s = t(val as any, "");
      if (Array.isArray(s)) return s as any;
      if (typeof s === "string") return s.split(/\n+/).filter(Boolean);
    }
    return [];
  };

  // Localized question (string or i18n map)
  const question: string | undefined = localizeString(node?.screen?.question);

  // Primary button (index 0) contains instructions
  const primaryBtn = Array.isArray(node?.screen?.buttons)
    ? node.screen.buttons[0]
    : undefined;
  const instructions: string[] = localizeList(primaryBtn?.instructions?.text);
  const download = (primaryBtn?.instructions?.download || undefined) as
    | { label?: string; asset_id?: string; asset_path?: string }
    | undefined;
  const downloads = (primaryBtn?.instructions?.downloads || undefined) as
    | Array<{ label?: any; asset_id?: string; asset_path?: string }>
    | undefined;
  const accordionLabel =
    localizeString(primaryBtn?.label) ||
    t(
      {
        ru: "Инструкция по прохождению",
        kz: "Өту бойынша нұсқаулық",
        en: "How to complete",
      },
      "Инструкция по прохождению"
    );

  // Confirmation button (index 1)
  const confirmBtn = Array.isArray(node?.screen?.buttons)
    ? node.screen.buttons[1]
    : undefined;
  const confirmLabel: string =
    localizeString(confirmBtn?.label) ||
    t({ ru: "Подтвердить", kz: "Растау", en: "Confirm" }, "Подтвердить");
  const confirmText: string =
    localizeString(confirmBtn?.confirmation_text) ||
    t(
      {
        ru: "Подтвердить выполнение шага?",
        kz: "Қадамды орындауды растау?",
        en: "Confirm completing this step?",
      },
      "Подтвердить выполнение шага?"
    );

  const completedMessage: string =
    localizeString(node?.states?.completed?.message) ||
    t(
      { ru: "Шаг подтверждён.", kz: "Қадам расталды.", en: "Step confirmed." },
      "Шаг подтверждён."
    );

  const handleConfirm = () => {
    setCompleted(true);
    setConfirmOpen(false);
    push({
      title: t(
        { ru: "Шаг подтверждён", kz: "Қадам расталды", en: "Step confirmed" },
        "Шаг подтверждён"
      ),
      description: t(
        {
          ru: "Действие успешно отмечено как выполненное.",
          kz: "Әрекет сәтті аяқталған ретінде белгіленді.",
          en: "Action marked as completed.",
        },
        "Действие успешно отмечено как выполненное."
      ),
    });
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
                  const renderItem = (item: any, key: string | number) => {
                    const resolved = item?.asset_id
                      ? getAssetUrl(item.asset_id)
                      : undefined;
                    const href = resolved && resolved !== "#" ? resolved : item?.asset_path;
                    if (!href) return null;
                    const label = localizeString(item?.label) ||
                      t(
                        { ru: "Скачать шаблон", kz: "Үлгіні жүктеу", en: "Download template" },
                        "Скачать шаблон"
                      );
                    return (
                      <div key={key} className="mt-3">
                        <Button asChild variant="secondary">
                          <a href={href} download target="_blank" rel="noopener noreferrer">
                            {label}
                          </a>
                        </Button>
                      </div>
                    );
                  };
                  if (Array.isArray(downloads) && downloads.length) {
                    return <>{downloads.map((d, idx) => renderItem(d, idx))}</>;
                  }
                  if (download) {
                    return renderItem(download, "single");
                  }
                  return null;
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
