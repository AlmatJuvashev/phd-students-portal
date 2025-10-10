import { useState, useEffect } from "react";
import type { NodeVM, FieldDef } from "@/lib/playbook";
import { t } from "@/lib/playbook";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { ChecklistItem } from "@/components/ui/checklist-item";
import { Check } from "lucide-react";
import { useTranslation } from "react-i18next";
import { TemplatesPanel } from "@/features/forms/TemplatesPanel";
import { ConfirmModal } from "@/features/forms/ConfirmModal";

export default function D2ApplyToDsScene({
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
    (Array.isArray(node.next) ? node.next[0] : undefined) ||
    node.outcomes?.[0]?.next?.[0];

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
                ru: `Если документы были сданы${
                  submittedAt
                    ? ` (дата: ${new Date(submittedAt).toLocaleDateString()})`
                    : ""
                }.`,
                kz: `Егер құжаттар тапсырылған болса${
                  submittedAt
                    ? ` (күні: ${new Date(submittedAt).toLocaleDateString()})`
                    : ""
                }.`,
                en: `If the documents were submitted${
                  submittedAt
                    ? ` (date: ${new Date(submittedAt).toLocaleDateString()})`
                    : ""
                }.`,
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
        message={t((node as any).description, "")}
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
