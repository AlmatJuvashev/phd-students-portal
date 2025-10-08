import React from "react";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Modal } from "@/components/ui/modal";
import { useToast } from "@/components/toast";
import { getAssetUrl } from "@/lib/assets";
import type { NodeVM } from "@/lib/playbook";
import { t } from "@/lib/playbook";

type ConfirmTaskDetailsProps = {
  node: NodeVM | any;
  onComplete?: () => void;
};

const ConfirmTaskDetails: React.FC<ConfirmTaskDetailsProps> = ({ node, onComplete }) => {
  const { push } = useToast();
  const [isCompleted, setCompleted] = React.useState(node?.state === "done" || node?.status === "completed");
  const [confirmOpen, setConfirmOpen] = React.useState(false);

  const question: string | undefined = node?.screen?.question;
  const instructions: string[] = node?.screen?.buttons?.[0]?.instructions?.text || [];
  const instructionsTitle = t(
    {
      ru: "Инструкция по прохождению",
      kz: "Өту бойынша нұсқаулық",
      en: "How to complete",
    },
    "Инструкция по прохождению"
  );

  const handleConfirm = () => {
    setCompleted(true);
    setConfirmOpen(false);
    push({
      title: "Нормоконтроль подтверждён",
      description: "Квитанция получена. Вы можете перейти к следующему шагу.",
    });
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
              <div className="p-3">
                <ul className="list-disc pl-5 space-y-1 text-sm text-muted-foreground">
                  {instructions.map((line: string, idx: number) => (
                    <li key={idx}>{line}</li>
                  ))}
                </ul>
                <div className="mt-3">
                  <Button asChild variant="secondary">
                    <a href={getAssetUrl("tpl_ncste_normocontrol_letter_ru_docx")} download>
                      Скачать пример письма
                    </a>
                  </Button>
                </div>
              </div>
            </div>

            <div className="pt-2">
              <Button variant="default" className="mt-2" onClick={() => setConfirmOpen(true)}>
                Подтвердить получение квитанции
              </Button>
            </div>

            <Modal open={confirmOpen} onClose={() => setConfirmOpen(false)}>
              <div className="space-y-4">
                <div className="text-base font-medium">
                  Вы уверены, что получили квитанцию о прохождении нормоконтроля в НЦГНТЭ?
                </div>
                <div className="flex justify-end gap-2">
                  <Button variant="ghost" onClick={() => setConfirmOpen(false)}>
                    Отмена
                  </Button>
                  <Button onClick={handleConfirm}>Да, подтвердить</Button>
                </div>
              </div>
            </Modal>
          </>
        )}

        {isCompleted && (
          <div className="rounded-2xl bg-emerald-50 p-4 text-emerald-700">
            ✅ Нормоконтроль пройден и квитанция подтверждена.
          </div>
        )}
      </CardContent>
    </Card>
  );
};

export default ConfirmTaskDetails;
