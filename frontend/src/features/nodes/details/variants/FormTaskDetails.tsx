import { useCallback, useEffect, useMemo, useState, type ReactNode } from "react";
import { useTranslation } from "react-i18next";
import { AnimatePresence, motion } from "framer-motion";

import { Card } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import { ActionsBar } from "@/features/forms/ActionsBar";
import { FieldRenderer } from "@/features/forms/FieldRenderer";
import { evalVisible as evalVisibleExpr } from "@/features/forms/Visibility";
import { AssetsDownloads } from "@/features/nodes/details/AssetsDownloads";
import { FieldDef, NodeVM, t as pickLocale } from "@/lib/playbook";

export type FormTaskContext = {
  node: NodeVM;
  values: Record<string, any>;
  setField: (key: string, value: any) => void;
  canEdit: boolean;
  disabled: boolean;
  submit: (extra?: Record<string, any>) => void;
  saveDraft: (extra?: Record<string, any>) => void;
  evalVisible: (expr?: string) => boolean;
};

export type FormTaskDetailsProps = {
  node: NodeVM;
  onSubmit?: (payload: any) => void;
  initial?: Record<string, any>;
  canEdit?: boolean;
  disabled?: boolean;
  renderActions?: (ctx: FormTaskContext) => ReactNode;
};

export function FormTaskDetails({
  node,
  initial = {},
  onSubmit,
  canEdit = true,
  disabled = false,
  renderActions,
}: FormTaskDetailsProps) {
  const [values, setValues] = useState<Record<string, any>>(initial);
  useEffect(() => {
    setValues(initial ?? {});
  }, [initial]);

  const { t: T } = useTranslation("common");
  const fields: FieldDef[] = node.requirements?.fields ?? [];

  const cardsLayout = (node.requirements as any)?.ui_hints?.cards_layout;
  const buttonsStyle = (node.requirements as any)?.ui_hints?.buttons_style;

  const setField = useCallback((key: string, value: any) => {
    setValues((prev) => {
      const next = { ...prev, [key]: value };
      if (key === "hearing_happened") {
        delete next.remarks_exist;
        delete next.plan_prepared;
        delete next.remarks_resolved;
      }
      if (key === "remarks_exist") {
        delete next.plan_prepared;
        delete next.remarks_resolved;
      }
      if (key === "plan_prepared") {
        delete next.remarks_resolved;
      }
      return next;
    });
  }, []);

  const evalVisible = useCallback(
    (expr?: string) => evalVisibleExpr(values, expr),
    [values]
  );

  const submit = useCallback(
    (extra: Record<string, any> = {}) => {
      if (!onSubmit) return;
      onSubmit({ ...values, ...extra });
    },
    [onSubmit, values]
  );

  const saveDraft = useCallback(
    (extra: Record<string, any> = {}) => {
      if (!onSubmit) return;
      onSubmit({ ...values, ...extra, __draft: true });
    },
    [onSubmit, values]
  );

  const ctx = useMemo<FormTaskContext>(
    () => ({ node, values, setField, canEdit, disabled, submit, saveDraft, evalVisible }),
    [node, values, setField, canEdit, disabled, submit, saveDraft, evalVisible]
  );

  const defaultActions = useMemo(() => {
    if (!canEdit) return null;
    return (
      <ActionsBar
        onSubmit={() => submit()}
        onDraft={() => saveDraft()}
        disabled={disabled}
      />
    );
  }, [canEdit, disabled, saveDraft, submit]);

  return (
    <div className="flex flex-col h-full">
      <div className="grid grid-cols-1 lg:grid-cols-5 gap-4 flex-1 min-h-0">
        <div className="lg:col-span-3 flex flex-col min-h-0">
          <Card className="p-4 flex flex-col flex-1 min-h-0">
            {node.requirements?.notes && (
              <p className="text-sm text-muted-foreground mb-4">
                {node.requirements.notes}
              </p>
            )}

            <div className="flex-1 min-h-0 overflow-y-auto">
              <div className="space-y-3 pb-4">
                {cardsLayout?.style === "stacked" &&
                buttonsStyle === "yes_no" ? (
                  <div className="space-y-3">
                    <AnimatePresence initial={false}>
                      {fields.map((f) => {
                        const visible = evalVisible((f as any).visible_when);
                        if (!visible) return null;

                        if (f.key === "remarks_exist") {
                          const label = pickLocale(f.label as any, f.key);
                          return (
                            <motion.div
                              key={f.key}
                              initial={{ opacity: 0, y: 20 }}
                              animate={{ opacity: 1, y: 0 }}
                              exit={{ opacity: 0, y: 10 }}
                              transition={{ duration: 0.2 }}
                            >
                              <Card className="p-4">
                                <div className="mb-2 font-medium">{label}</div>
                                <div className="flex gap-2">
                                  <Button
                                    onClick={() => {
                                      setField(f.key, true);
                                      saveDraft({ [f.key]: true });
                                    }}
                                    disabled={!canEdit || disabled}
                                  >
                                    {T("forms.yes")}
                                  </Button>
                                  <Button
                                    variant="secondary"
                                    onClick={() => {
                                      setField(f.key, false);
                                      saveDraft({ [f.key]: false });
                                    }}
                                    disabled={!canEdit || disabled}
                                  >
                                    {T("forms.no")}
                                  </Button>
                                </div>
                              </Card>
                            </motion.div>
                          );
                        }

                        if (
                          f.key === "plan_prepared" &&
                          values["remarks_exist"] === true
                        ) {
                          const label = pickLocale(f.label as any, f.key);
                          return (
                            <motion.div
                              key={f.key}
                              initial={{ opacity: 0, y: 20 }}
                              animate={{ opacity: 1, y: 0 }}
                              exit={{ opacity: 0, y: 10 }}
                              transition={{ duration: 0.2 }}
                            >
                              <Card className="p-4">
                                <div className="mb-2 font-medium">{label}</div>
                                <div className="flex gap-2">
                                  <Button
                                    onClick={() => {
                                      setField(f.key, true);
                                      saveDraft({ [f.key]: true });
                                    }}
                                    disabled={!canEdit || disabled}
                                  >
                                    {T("forms.yes")}
                                  </Button>
                                  <Button
                                    variant="secondary"
                                    onClick={() => {
                                      setField(f.key, false);
                                      saveDraft({ [f.key]: false });
                                    }}
                                    disabled={!canEdit || disabled}
                                  >
                                    {T("forms.no")}
                                  </Button>
                                </div>
                              </Card>
                            </motion.div>
                          );
                        }

                        if (
                          f.key === "remarks_resolved" &&
                          values["plan_prepared"] === true
                        ) {
                          const label = pickLocale(f.label as any, f.key);
                          return (
                            <motion.div
                              key={f.key}
                              initial={{ opacity: 0, y: 20 }}
                              animate={{ opacity: 1, y: 0 }}
                              exit={{ opacity: 0, y: 10 }}
                              transition={{ duration: 0.2 }}
                            >
                              <Card className="p-4">
                                <div className="mb-2 font-medium">{label}</div>
                                <div className="flex gap-2">
                                  <Button
                                    onClick={() => {
                                      setField(f.key, true);
                                      saveDraft({ [f.key]: true });
                                    }}
                                    disabled={!canEdit || disabled}
                                  >
                                    {T("forms.yes")}
                                  </Button>
                                  <Button
                                    variant="secondary"
                                    onClick={() => {
                                      setField(f.key, false);
                                      saveDraft({ [f.key]: false });
                                    }}
                                    disabled={!canEdit || disabled}
                                  >
                                    {T("forms.no")}
                                  </Button>
                                </div>
                              </Card>
                            </motion.div>
                          );
                        }

                        return null;
                      })}
                    </AnimatePresence>
                  </div>
                ) : (
                  <div className="space-y-3">
                    {fields.map((f) => (
                      <FieldRenderer
                        key={f.key}
                        field={f as any}
                        value={values[f.key]}
                        onChange={(v) => setField(f.key, v)}
                        setField={(k, v) => setField(k, v)}
                        otherValue={values[`${f.key}_other`]}
                        canEdit={canEdit}
                        disabled={disabled}
                      />
                    ))}
                  </div>
                )}
              </div>

              {!!node.requirements?.validations?.length && (
                <>
                  <Separator />
                  <div>
                    <div className="mb-2 font-medium">
                      {T("forms.validations_title")}
                    </div>
                    <ul className="list-inside list-disc text-sm">
                      {node.requirements.validations!.map((v, i) => (
                        <li key={i}>
                          {v.rule}
                          {v.source ? ` @ ${v.source}` : ""}
                        </li>
                      ))}
                    </ul>
                  </div>
                </>
              )}
            </div>

            {canEdit && (
              <div className="space-y-2">
                {renderActions ? renderActions(ctx) : defaultActions}
              </div>
            )}
          </Card>
        </div>

        <div className="lg:col-span-2 border-l pl-4">
          <AssetsDownloads node={node} />
        </div>
      </div>
    </div>
  );
}

export default FormTaskDetails;
