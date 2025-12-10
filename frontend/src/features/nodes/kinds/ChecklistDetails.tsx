import { useEffect, useState } from "react";
import type { NodeVM, FieldDef } from "@/lib/playbook";
import { t } from "@/lib/playbook";
import { useTranslation } from "react-i18next";
import { ConfirmModal } from "@/features/forms/ConfirmModal";
import { motion } from "framer-motion";
import { Check } from "lucide-react";
import { clsx, type ClassValue } from "clsx";
import { twMerge } from "tailwind-merge";
import { AssetsDownloads } from "@/features/nodes/details/AssetsDownloads";

function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export default function ChecklistDetails({
  node,
  initial = {},
  disabled,
  onSubmit,
  canEdit,
}: {
  node: NodeVM;
  initial?: Record<string, any>;
  disabled?: boolean;
  onSubmit?: (payload: any) => void;
  canEdit?: boolean;
}) {
  const { t: T } = useTranslation("common");
  const [values, setValues] = useState<Record<string, any>>(initial);
  const [confirmOpen, setConfirmOpen] = useState(false);
  useEffect(() => setValues(initial ?? {}), [initial]);

  const fields: FieldDef[] = node.requirements?.fields ?? [];
  const bools = fields.filter((f) => f.type === "boolean");
  const requiredBools = bools.filter((f) => f.required);
  const ready = requiredBools.every((f) => !!values[f.key]);
  // If canEdit is explicitly provided, use it; otherwise fall back to state-based readonly
  const readOnly =
    canEdit !== undefined
      ? !canEdit
      : node.state === "submitted" ||
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
          const requiredFields =
            outcomeAny.when
              .match(/form\.(\w+)/g)
              ?.map((match: string) => match.replace("form.", "")) || [];
          const allRequired = requiredFields.every(
            (field: string) => !!values[field]
          );

          if (allRequired) {
            return outcome.next?.[0];
          }
        }
      }
      // If no condition matches, return first outcome
      return node.outcomes[0]?.next?.[0];
    }

    // Fallback to simple next logic
    return Array.isArray(node.next) ? node.next[0] : undefined;
  };

  const nextOnComplete = getNextOnComplete();

  return (
    <form className="h-full flex flex-col" onSubmit={(e) => e.preventDefault()}>
      <div className="flex-1 overflow-y-auto min-h-0">
          <div className="space-y-6 p-1">
            {Boolean((node as any)?.description) && (
                <div className="prose prose-slate prose-sm max-w-none text-slate-600 bg-slate-50/50 p-4 rounded-2xl border border-slate-100">
                <p>{t((node as any).description, "")}</p>
                </div>
            )}

            {/* Quest Progress Bar */}
            {bools.length > 0 && (
              <div className="mb-6 px-2">
                <div className="flex justify-between text-xs font-bold uppercase tracking-wider text-slate-400 mb-2">
                  <span>{T('journey.progress_label')}</span>
                  <span className={cn(bools.every(f => !!values[f.key]) && "text-emerald-500")}>
                    {bools.filter(f => !!values[f.key]).length} / {bools.length}
                  </span>
                </div>
                <div className="h-2 w-full bg-slate-200 rounded-full overflow-hidden">
                  <motion.div 
                    initial={false}
                    animate={{ width: `${(bools.filter(f => !!values[f.key]).length / bools.length) * 100}%` }}
                    transition={{ type: "spring", stiffness: 300, damping: 30 }}
                    className={cn(
                      "h-full rounded-full transition-colors duration-500",
                      bools.every(f => !!values[f.key]) ? "bg-emerald-500" : "bg-primary-500"
                    )}
                  />
                </div>
              </div>
            )}

            {/* Checklist Items */}
            <div className="space-y-3">
                {bools.map((f) => {
                const isChecked = !!values[f.key];

                return (
                    <motion.button
                    key={f.key}
                    type="button"
                    onClick={() => {
                        if (!readOnly) {
                        setValues((s) => ({ ...s, [f.key]: !isChecked }));
                        }
                    }}
                    disabled={disabled || readOnly}
                    whileTap={!readOnly ? { scale: 0.98 } : undefined}
                    initial={false}
                    animate={{
                        backgroundColor: isChecked ? "rgb(236 253 245)" : "rgba(255, 255, 255, 1)",
                        borderColor: isChecked ? "rgb(167 243 208)" : "rgb(226 232 240)"
                    }}
                    className={cn(
                        "w-full text-left p-4 rounded-xl border-2 transition-all flex items-start gap-4 group relative overflow-hidden",
                        isChecked 
                        ? "shadow-sm" 
                        : "hover:border-primary-300 hover:shadow-md"
                    )}
                    >
                    <div className={cn(
                        "w-6 h-6 rounded-md border-2 flex items-center justify-center transition-colors duration-200 flex-shrink-0 mt-0.5",
                        isChecked
                        ? "bg-emerald-500 border-emerald-500 text-white"
                        : "bg-white border-slate-300 group-hover:border-primary-400"
                    )}>
                        {isChecked && <Check size={16} strokeWidth={4} />}
                    </div>
                    <span className={cn(
                        "text-sm font-medium transition-colors",
                        isChecked ? "text-emerald-800 line-through opacity-70" : "text-slate-700"
                    )}>
                        {t(f.label, f.key)}{f.required ? "*" : ""}
                    </span>
                    </motion.button>
                );
                })}
            </div>

            {/* Read-only status message */}
            {readOnly && (
                <div className="mt-3 text-sm text-slate-400 whitespace-pre-line px-2">
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
            
            <div className="pt-4 px-2">
                 <AssetsDownloads node={node} />
            </div>
          </div>
      </div>

      {/* Footer Actions - Fixed at bottom */}
      <div className="mt-6 pt-4 border-t border-slate-100 flex-shrink-0 bg-white z-10">
        {!readOnly ? (
            <button 
            onClick={() => setConfirmOpen(true)}
            disabled={!ready || disabled}
            className={cn(
                "w-full py-4 font-bold text-lg rounded-2xl transition-all shadow-lg active:scale-[0.98] flex items-center justify-center gap-2",
                ready 
                ? "bg-emerald-500 hover:bg-emerald-600 text-white shadow-emerald-500/20" 
                : "bg-slate-200 text-slate-400 cursor-not-allowed"
            )}
            >
            {ready ? (
                <>
                {T("forms.proceed_next")} <Check size={20} />
                </>
            ) : (
                <span className="flex items-center gap-2">
                  <span className="text-xs uppercase tracking-wider font-extrabold">{values?.__draft ? "Save Draft" : "Complete All Items"}</span>
                </span>
            )}
            </button>
        ) : (
             <div className="w-full py-3 bg-emerald-50 text-emerald-600 rounded-xl border border-emerald-100 font-medium text-center flex items-center justify-center gap-2">
                <Check size={18} /> Completed
             </div>
        )}
        
        {!readOnly && (
            <div className="mt-3 text-center">
                 <button
                    type="button"
                    onClick={() => onSubmit?.({ ...values, __draft: true })}
                    disabled={disabled}
                    className="text-xs font-semibold text-slate-400 hover:text-slate-600 uppercase tracking-widest hover:underline"
                >
                    {T("forms.save_draft")}
                </button>
            </div>
        )}
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
    </form>

  );
}
