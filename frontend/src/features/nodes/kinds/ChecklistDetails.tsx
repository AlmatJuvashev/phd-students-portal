import { useEffect, useState } from "react";
import type { NodeVM, FieldDef } from "@/lib/playbook";
import { t } from "@/lib/playbook";
import { Button } from "@/components/ui/button";
import { ChecklistItem } from "@/components/ui/checklist-item";
import { useTranslation } from "react-i18next";
import { TemplatesPanel } from "@/features/forms/TemplatesPanel";
import { ConfirmModal } from "@/features/forms/ConfirmModal";

export default function ChecklistDetails({
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
  const bools = fields.filter((f) => f.type === "boolean");
  const requiredBools = bools.filter((f) => f.required);
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
    <div className="h-full">
      <div className="lg:grid lg:grid-cols-[minmax(0,3fr)_minmax(0,2fr)] lg:gap-6 space-y-6 lg:space-y-0">
        <div className="space-y-4">
          {Boolean((node as any)?.description) && (
            <div className="text-sm text-muted-foreground mb-4 p-4 rounded-lg bg-muted/30 border-l-4 border-primary/50">
              {t((node as any).description, "")}
            </div>
          )}
          
          {/* Checklist items - simplified container */}
          <div className="space-y-3" style={{ minHeight: 'auto' }}>
            {bools.map((f, index) => {
              const isChecked = !!values[f.key];
              
              return (
                <div 
                  key={f.key}
                  className="border p-3 rounded bg-white"
                  style={{ display: 'block', visibility: 'visible' }}
                >
                  <label className="flex items-center gap-3 cursor-pointer">
                    <span className="flex-1 text-sm">
                      {t(f.label, f.key)}
                      {f.required && <span className="text-red-500 ml-1">*</span>}
                    </span>
                    <input
                      type="checkbox"
                      checked={isChecked}
                      onChange={(e) => {
                        setValues((s) => ({ ...s, [f.key]: e.target.checked }));
                      }}
                      disabled={disabled}
                      className="w-5 h-5"
                    />
                  </label>
                </div>
              );
            })}
          </div>

          {/* Action buttons */}
          {!readOnly && (
            <div className="flex gap-2 pt-4 border-t bg-background/80 backdrop-blur-sm sticky bottom-0">
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

        <TemplatesPanel node={node} className="lg:border-l lg:pl-6" />
      </div>

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
