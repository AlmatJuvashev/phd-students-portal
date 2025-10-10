import React from "react";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Modal } from "@/components/ui/modal";
import { getAssetUrl } from "@/lib/assets";
import type { NodeVM } from "@/lib/playbook";
import { t, safeText } from "@/lib/playbook";
import i18n from "i18next";

type ConfirmUploadTaskDetailsProps = {
  node: NodeVM | any;
  onComplete?: () => void;
};

const ConfirmUploadTaskDetails: React.FC<ConfirmUploadTaskDetailsProps> = ({
  node,
  onComplete,
}) => {
  const [isCompleted, setCompleted] = React.useState(
    node?.state === "done" || node?.status === "completed"
  );
  const [confirmOpen, setConfirmOpen] = React.useState(false);

  const question = safeText(node?.screen?.question as any, "");
  // Determine current language and single template by language
  const lang = (i18n?.language as "ru" | "kz" | "en") || "ru";
  // Normalize instructions: could be array of strings or locale->string[] map
  const rawText = node?.screen?.buttons?.[0]?.instructions?.text as
    | string[]
    | Record<string, string[]>
    | undefined;
  const instructions: string[] = Array.isArray(rawText)
    ? rawText
    : Array.isArray((rawText as any)?.[lang])
    ? (rawText as any)[lang]
    : [];
  const templateByLang: Record<typeof lang, string> = {
    ru: "tpl_letter_rector_defense_request_ru_docx",
    kz: "tpl_letter_rector_defense_request_kz_docx",
    en: "tpl_letter_rector_defense_request_en_docx",
  } as const;
  const templateLabelByLang: Record<typeof lang, string> = {
    ru: "Скачать шаблон заявления",
    kz: "Үлгіні жүктеу",
    en: "Download template letter",
  } as const;
  const singleAssetId = templateByLang[lang];
  const singleAssetLabel = templateLabelByLang[lang];

  const instructionsTitle = safeText(
    {
      ru: "Как оформить заявление",
      kz: "Өтінішті қалай ресімдеу",
      en: "How to prepare the letter",
    },
    "Как оформить заявление"
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
            <div className="divide-y rounded-xl border overflow-hidden">
              <div className="w-full text-left p-3 font-medium bg-muted">
                {instructionsTitle}
              </div>
              <div className="p-3 space-y-3">
                <ul className="list-disc pl-5 space-y-1 text-sm text-muted-foreground">
                  {instructions.map((line: string, idx: number) => (
                    <li key={idx}>{line}</li>
                  ))}
                </ul>
                <div className="mt-1 flex flex-col gap-2">
                  <Button asChild variant="secondary" className="w-fit">
                    <a
                      href={getAssetUrl(singleAssetId)}
                      download
                      target="_blank"
                      rel="noopener noreferrer"
                    >
                      {singleAssetLabel}
                    </a>
                  </Button>
                </div>
              </div>
            </div>

            <div className="pt-2">
              <Button
                variant="default"
                className="mt-2"
                onClick={() => setConfirmOpen(true)}
              >
                {safeText(
                  node?.screen?.buttons?.[1]?.label,
                  "Подтвердить подачу заявления"
                )}
              </Button>
            </div>

            <Modal open={confirmOpen} onClose={() => setConfirmOpen(false)}>
              <div className="space-y-4">
                <div className="text-base font-medium">
                  {safeText(
                    node?.screen?.buttons?.[1]?.confirmation_text,
                    "Вы уверены, что подали заявление ректору о приёме к защите?"
                  )}
                </div>
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
            ✅ Заявление ректору подтверждено.
          </div>
        )}
      </CardContent>
    </Card>
  );
};

export default ConfirmUploadTaskDetails;
