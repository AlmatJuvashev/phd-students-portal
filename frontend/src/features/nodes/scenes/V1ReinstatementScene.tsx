import { useState, useEffect } from "react";
import type { NodeVM, FieldDef } from "@/lib/playbook";
import { t } from "@/lib/playbook";
import { Button } from "@/components/ui/button";
import { ChecklistItem } from "@/components/ui/checklist-item";
import { Check } from "lucide-react";
import { useTranslation } from "react-i18next";
import { AssetsDownloads } from "@/features/nodes/details/AssetsDownloads";
import { ConfirmModal } from "@/features/forms/ConfirmModal";
import { motion } from "framer-motion";
import { clsx, type ClassValue } from "clsx";
import { twMerge } from "tailwind-merge";

function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export default function V1ReinstatementScene({
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

  return (
    <div className="flex flex-col h-full overflow-auto p-1">
      <div className="space-y-4">
        {Boolean((node as any)?.description) && (
          <div className="text-sm text-muted-foreground mb-4 p-4 rounded-lg bg-muted/30 border-l-4 border-primary/50">
            {t((node as any).description, "")}
          </div>
        )}
        {/* Quest Progress Bar */}
        {fields.filter(f => f.type === 'boolean').length > 0 && (
          <div className="mb-6 px-1">
            <div className="flex justify-between text-xs font-bold uppercase tracking-wider text-muted-foreground mb-2">
              <span>Quest Progress</span>
              <span className={cn(fields.filter(f => f.type === 'boolean').every(f => !!values[f.key]) && "text-emerald-500")}>
                {fields.filter(f => f.type === 'boolean' && !!values[f.key]).length} / {fields.filter(f => f.type === 'boolean').length}
              </span>
            </div>
            <div className="h-2 w-full bg-secondary rounded-full overflow-hidden">
              <motion.div 
                initial={false}
                animate={{ width: `${(fields.filter(f => f.type === 'boolean' && !!values[f.key]).length / fields.filter(f => f.type === 'boolean').length) * 100}%` }}
                transition={{ type: "spring", stiffness: 300, damping: 30 }}
                className={cn(
                  "h-full rounded-full transition-colors duration-500",
                  fields.filter(f => f.type === 'boolean').every(f => !!values[f.key]) ? "bg-emerald-500" : "bg-primary"
                )}
              />
            </div>
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
                ru: `Если пакет на восстановление был подан${
                  submittedAt
                    ? ` (дата: ${new Date(submittedAt).toLocaleDateString()})`
                    : ""
                }.`,
                kz: `Егер қалпына келтіру топтамасы тапсырылған болса${
                  submittedAt
                    ? ` (күні: ${new Date(submittedAt).toLocaleDateString()})`
                    : ""
                }.`,
                en: `If the reinstatement package was submitted${
                  submittedAt
                    ? ` (date: ${new Date(submittedAt).toLocaleDateString()})`
                    : ""
                }.`,
              },
              ""
            )}
          </div>
        )}

        <div className="pt-8">
             <AssetsDownloads node={node} />
        </div>
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
