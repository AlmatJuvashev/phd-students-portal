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
  
  // Calculate next node based on outcomes conditions
  const getNextOnComplete = () => {
    if (node.outcomes && node.outcomes.length > 0) {
      // Find the first outcome that matches current form state
      for (const outcome of node.outcomes) {
        const outcomeAny = outcome as any;
        if (outcomeAny.when) {
          // Simple evaluation of when condition
          // For NK_package: "form.chk_thesis_unbound && form.chk_advisor_reviews && form.chk_pubs_app7 && form.chk_sc_extract && form.chk_lcb_defense"
          const requiredFields = outcomeAny.when.match(/form\.(\w+)/g)?.map((match: string) => match.replace('form.', '')) || [];
          const allRequired = requiredFields.every((field: string) => !!values[field]);
          
          if (allRequired) {
            return outcome.next?.[0];
          }
        }
      }
      // If no condition matches, return first outcome
      return node.outcomes[0]?.next?.[0];
    }
    
    // Fallback to simple next logic
    return (Array.isArray(node.next) ? node.next[0] : undefined);
  };
  
  const nextOnComplete = getNextOnComplete();

  return (
    <div className="h-full">
      <div className="lg:grid lg:grid-cols-[minmax(0,3fr)_minmax(0,2fr)] lg:gap-6 space-y-6 lg:space-y-0">
        <div className="space-y-4">
          {Boolean((node as any)?.description) && (
            <div className="text-sm text-muted-foreground mb-4 p-4 rounded-lg bg-muted/30 border-l-4 border-primary/50">
              {t((node as any).description, "")}
            </div>
          )}

          {/* Checklist items using ChecklistItem component */}
          <div className="space-y-3">
            {bools.map((f, index) => {
              const isChecked = !!values[f.key];

              return (
                <ChecklistItem
                  key={f.key}
                  checked={isChecked}
                  onChange={(checked) => {
                    if (!readOnly) {
                      setValues((s) => ({ ...s, [f.key]: checked }));
                    }
                  }}
                  label={`${t(f.label, f.key)}${f.required ? '*' : ''}`}
                  readOnly={readOnly}
                  disabled={disabled}
                />
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
          console.log('[ChecklistDetails] Submitting with next:', nextOnComplete);
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
