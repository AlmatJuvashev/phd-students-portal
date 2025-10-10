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
    <Card className="bg-gradient-to-br from-card to-card/50">
      <CardContent className="space-y-5">
        {question && (
          <p className="text-lg sm:text-xl font-semibold text-foreground leading-relaxed">
            {question}
          </p>
        )}

        {!isCompleted && (
          <>
            <div className="divide-y rounded-2xl border-2 border-border/50 overflow-hidden shadow-sm">
              <div className="w-full text-left p-4 font-semibold bg-gradient-to-r from-muted to-muted/50 flex items-center gap-2">
                <svg
                  className="w-5 h-5 text-primary"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
                  />
                </svg>
                {instructionsTitle}
              </div>
              <div className="p-5 space-y-4 bg-gradient-to-b from-background to-muted/5">
                <ul className="list-disc pl-5 space-y-2 text-sm text-muted-foreground">
                  {instructions.map((line: string, idx: number) => (
                    <li key={idx} className="leading-relaxed">
                      {line}
                    </li>
                  ))}
                </ul>
                <div className="flex flex-col gap-2 pt-2">
                  <Button
                    asChild
                    variant="secondary"
                    size="default"
                    className="w-full sm:w-fit gap-2"
                  >
                    <a
                      href={getAssetUrl(singleAssetId)}
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
                      {singleAssetLabel}
                    </a>
                  </Button>
                </div>
              </div>
            </div>

            <div className="pt-2">
              <Button
                variant="default"
                size="lg"
                className="w-full sm:w-auto"
                onClick={() => setConfirmOpen(true)}
              >
                {safeText(
                  node?.screen?.buttons?.[1]?.label,
                  "Подтвердить подачу заявления"
                )}
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
                        d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
                      />
                    </svg>
                  </div>
                  <div className="text-base sm:text-lg font-semibold text-foreground leading-relaxed px-2">
                    {safeText(
                      node?.screen?.buttons?.[1]?.confirmation_text,
                      "Вы уверены, что подали заявление ректору о приёме к защите?"
                    )}
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
              <span className="text-green-700 dark:text-green-300 font-semibold">
                Заявление ректору подтверждено.
              </span>
            </div>
          </div>
        )}
      </CardContent>
    </Card>
  );
};

export default ConfirmUploadTaskDetails;
