import {
  useCallback,
  useEffect,
  useMemo,
  useReducer,
  type ReactNode,
} from "react";
import { useTranslation } from "react-i18next";
import { Card } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import { ActionsBar } from "@/features/forms/ActionsBar";
import { FieldRenderer } from "@/features/forms/FieldRenderer";
import { evalVisible as evalVisibleExpr } from "@/features/forms/Visibility";
import { AssetsDownloads } from "@/features/nodes/details/AssetsDownloads";
import { FieldDef, NodeVM, t as pbT } from "@/lib/playbook";

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
  validationErrors?: Record<string, Record<number, string>>;
  persistKey?: string;
  onValuesChange?: (values: Record<string, any>) => void;
};

export function FormTaskDetails({
  node,
  initial = {},
  onSubmit,
  canEdit = true,
  disabled = false,
  renderActions,
  validationErrors,
  persistKey,
  onValuesChange,
}: FormTaskDetailsProps) {
  const { t: T } = useTranslation("common");
  const fields: FieldDef[] = node.requirements?.fields ?? [];

  type FormState = Record<string, any>;
  type Action =
    | { type: "set"; key: string; value: any }
    | { type: "replace"; values: FormState };

  const formReducer = useCallback(
    (state: FormState, action: Action): FormState => {
      switch (action.type) {
        case "replace":
          return { ...action.values };
        case "set": {
          const next = { ...state, [action.key]: action.value } as FormState;
          if (action.key === "hearing_happened") {
            delete next.remarks_exist;
            delete next.plan_prepared;
            delete next.remarks_resolved;
          }
          if (action.key === "remarks_exist") {
            delete next.plan_prepared;
            delete next.remarks_resolved;
          }
          if (action.key === "plan_prepared") {
            delete next.remarks_resolved;
          }
          return next;
        }
        default:
          return state;
      }
    },
    []
  );

  const [values, dispatch] = useReducer(
    formReducer,
    initial ?? {},
    (start) => ({ ...(start ?? {}) })
  );

  useEffect(() => {
    let merged = initial ?? {};
    if (persistKey) {
      try {
        const saved = localStorage.getItem(persistKey);
        if (saved) {
          const parsed = JSON.parse(saved);
          merged = { ...merged, ...parsed };
        }
      } catch (e) {
        console.error("Failed to load draft from localStorage", e);
      }
    }
    dispatch({ type: "replace", values: merged });
  }, [initial, persistKey]);

  useEffect(() => {
    if (persistKey && canEdit && !disabled) {
      const handler = setTimeout(() => {
        try {
          localStorage.setItem(persistKey, JSON.stringify(values));
        } catch (e) {
          console.error("Failed to save draft to localStorage", e);
        }
      }, 500);
      return () => clearTimeout(handler);
    }
  }, [values, persistKey, canEdit, disabled]);

  useEffect(() => {
    onValuesChange?.(values);
  }, [values, onValuesChange]);

  const setField = useCallback(
    (key: string, value: any) => dispatch({ type: "set", key, value }),
    []
  );

  const evalVisible = useCallback(
    (expr?: string) => evalVisibleExpr(values, expr),
    [values]
  );

  const submit = useCallback(
    (extra: Record<string, any> = {}) => {
      if (!onSubmit) return;
      if (persistKey) {
        localStorage.removeItem(persistKey);
      }
      onSubmit({ ...values, ...extra });
    },
    [onSubmit, values, persistKey]
  );

  const saveDraft = useCallback(
    (extra: Record<string, any> = {}) => {
      if (!onSubmit) return;
      if (persistKey) {
        localStorage.removeItem(persistKey);
      }
      onSubmit({ ...values, ...extra, __draft: true });
    },
    [onSubmit, values, persistKey]
  );

  const ctx = useMemo<FormTaskContext>(
    () => ({
      node,
      values,
      setField,
      canEdit,
      disabled,
      submit,
      saveDraft,
      evalVisible,
    }),
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

  // NK_package: unified button style and simple completion check
  const nkActions = useMemo(() => {
    if (!canEdit) return null;
    const totalFields = (node.requirements?.fields ?? []).length;
    const required = (node.requirements?.fields ?? []).filter((f) => (f as any).required);
    const completedCount = required.reduce(
      (acc, f) => acc + (values[f.key] ? 1 : 0),
      0
    );
    const allRequiredFilled = required.every((f) => !!values[f.key]);

    const outcomeLabel = pbT(node.outcomes?.[0]?.label, "") || T("forms.proceed_next");

    return (
      <>
        <div className="flex items-center gap-2 text-sm text-muted-foreground">
          <span>
            {completedCount} / {required.length || totalFields}
          </span>
        </div>
        <div className="flex flex-col sm:flex-row gap-2">
          <Button
            variant="secondary"
            onClick={() => saveDraft()}
            disabled={disabled}
            className="w-full sm:w-auto"
          >
            {T("forms.save_draft")}
          </Button>
          <Button
            onClick={() => submit()}
            disabled={disabled || !allRequiredFilled}
            className="w-full sm:w-auto"
          >
            {outcomeLabel}
          </Button>
        </div>
      </>
    );
  }, [T, canEdit, disabled, node, saveDraft, submit, values]);

  return (
    <div className="flex flex-col h-full min-h-0" data-node-id={node.id}>
      {/* Mobile: show templates above the form */}
      <div className="lg:hidden mb-2">
        <AssetsDownloads node={node} />
      </div>
      <div className="grid grid-cols-1 lg:grid-cols-5 gap-4 flex-1 min-h-0">
        <div className="lg:col-span-3 flex flex-col min-h-0 min-w-0">
          <Card className="p-4 flex flex-col flex-1 min-h-0">
            {node.requirements?.notes && (
              <p className="text-sm text-muted-foreground mb-4">
                {node.requirements.notes}
              </p>
            )}

            <div className="space-y-3 pb-4 min-w-0">
              {fields.map((f) => {
                const visible = evalVisible((f as any).visible_when);
                if (!visible) return null;
                return (
                  <FieldRenderer
                    key={f.key}
                    field={f as any}
                    value={values[f.key]}
                    onChange={(v) => setField(f.key, v)}
                    setField={(k, v) => setField(k, v)}
                    otherValue={values[`${f.key}_other`]}
                    canEdit={canEdit}
                    disabled={disabled}
                    itemErrors={validationErrors?.[f.key]}
                  />
                );
              })}

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
              <div className="space-y-2 mt-4 form-actions">
                {renderActions
                  ? renderActions(ctx)
                  : node.id === "NK_package"
                  ? nkActions
                  : defaultActions}
              </div>
            )}
          </Card>
        </div>

        <div className="hidden lg:block lg:col-span-2 border-l pl-4 min-w-0">
          <AssetsDownloads node={node} />
        </div>
      </div>
    </div>
  );
}

export default FormTaskDetails;
