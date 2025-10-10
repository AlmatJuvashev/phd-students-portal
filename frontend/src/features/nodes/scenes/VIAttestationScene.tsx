import { useEffect, useState } from "react";
import type { NodeVM, FieldDef } from "@/lib/playbook";
import { t } from "@/lib/playbook";
import { Button } from "@/components/ui/button";
import { ChecklistItem } from "@/components/ui/checklist-item";
import { Check } from "lucide-react";
import { useTranslation } from "react-i18next";
import { TemplatesPanel } from "@/features/forms/TemplatesPanel";
import { ConfirmModal } from "@/features/forms/ConfirmModal";

export default function VIAttestationScene({
  node,
  initial = {},
  disabled,
  onSubmit,
}: {
  node: NodeVM;
  initial?: Record<string, any>;
  disabled?: boolean;
  onSubmit?: (payload: any) => void;
}) {
  const { t: T } = useTranslation("common");
  const [values, setValues] = useState<Record<string, any>>(initial);
  const [confirmOpen, setConfirmOpen] = useState(false);
  useEffect(() => setValues(initial ?? {}), [initial]);

  const fields: FieldDef[] = node.requirements?.fields ?? [];
  const requiredBools = fields.filter(
    (f) => f.type === "boolean" && f.required
  );
  const ready = requiredBools.every((f) => !!values[f.key]);
  const readOnly =
    node.state === "submitted" ||
    node.state === "done" ||
    Boolean((initial as any)?.__submittedAt);
  const submittedAt: string | undefined =
    (initial as any)?.__submittedAt || values?.__submittedAt;
  const nextOnComplete =
    (Array.isArray(node.next) ? node.next[0] : undefined) || undefined;

  const guardMessage = t(
    {
      ru: "Проверьте, что все пункты аттестационного дела отмечены. После подтверждения данные зафиксируются, и маршрут будет завершён.",
      kz: "Аттестаттау ісіндегі барлық тармақтардың белгіленгенін тексеріңіз. Растағаннан кейін деректер бекітіледі және маршрут аяқталады.",
      en: "Ensure all attestation file items are checked. After confirming, the data will be saved and the journey will complete.",
    },
    ""
  );

  return (
    <div className="grid grid-cols-1 lg:grid-cols-5 gap-4 h-full">
      <div className="lg:col-span-3 min-h-0 overflow-auto space-y-4">
        {Boolean((node as any)?.description) && (
          <div className="text-sm text-muted-foreground mb-4 p-4 rounded-lg bg-muted/30 border-l-4 border-primary/50">
            {t((node as any).description, "")}
          </div>
        )}
        <div className="space-y-3">
          {fields.map((f) => (
            <div key={f.key}>
              {f.type === "boolean" ? (
                <ChecklistItem
                  checked={!!values[f.key]}
                  onChange={(checked) =>
                    setValues((s) => ({ ...s, [f.key]: checked }))
                  }
                  label={t(f.label, f.key)}
                  required={f.required}
                  disabled={disabled}
                  readOnly={readOnly}
                />
              ) : null}
            </div>
          ))}
        </div>

        {!readOnly && (
          <div className="flex gap-2 pt-2">
            <Button
              onClick={() => setConfirmOpen(true)}
              disabled={!ready}
              aria-busy={disabled}
            >
              {T("forms.proceed_next")}
            </Button>
            <Button
              variant="secondary"
              onClick={() => onSubmit?.({ ...values, __draft: true })}
              aria-busy={disabled}
            >
              {T("forms.save_draft")}
            </Button>
          </div>
        )}

        {readOnly && (
          <div className="mt-3 text-sm text-muted-foreground whitespace-pre-line">
            {t(
              {
                ru: `Если аттестационное дело сформировано${
                  submittedAt
                    ? ` (дата: ${new Date(submittedAt).toLocaleDateString()})`
                    : ""
                }. ${guardMessage}`,
                kz: `Егер аттестаттау ісі қалыптастырылған болса${
                  submittedAt
                    ? ` (күні: ${new Date(submittedAt).toLocaleDateString()})`
                    : ""
                }. ${guardMessage}`,
                en: `If the attestation file is complete${
                  submittedAt
                    ? ` (date: ${new Date(submittedAt).toLocaleDateString()})`
                    : ""
                }. ${guardMessage}`,
              },
              ""
            )}
          </div>
        )}
      </div>

      <TemplatesPanel node={node} />

      <ConfirmModal
        open={confirmOpen}
        onOpenChange={setConfirmOpen}
        message={guardMessage}
        confirmLabel={T("forms.proceed_next")}
        cancelLabel={T("common.cancel")}
        onConfirm={() => {
          setConfirmOpen(false);
          onSubmit?.({
            ...values,
            __submittedAt: new Date().toISOString(),
            __nextOverride: nextOnComplete,
          });
        }}
      />
    </div>
  );
}
