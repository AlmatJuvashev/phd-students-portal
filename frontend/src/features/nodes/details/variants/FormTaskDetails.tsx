import { useCallback, useEffect, useMemo, useReducer, type ReactNode } from "react";
import { useTranslation } from "react-i18next";
import { Card } from "@/components/ui/card";
import { Separator } from "@/components/ui/separator";
import { ActionsBar } from "@/features/forms/ActionsBar";
import { FieldRenderer } from "@/features/forms/FieldRenderer";
import { evalVisible as evalVisibleExpr } from "@/features/forms/Visibility";
import { AssetsDownloads } from "@/features/nodes/details/AssetsDownloads";
import { FieldDef, NodeVM } from "@/lib/playbook";

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
  const { t: T } = useTranslation("common");
  const fields: FieldDef[] = node.requirements?.fields ?? [];

  type FormState = Record<string, any>;
  type Action =
    | { type: "set"; key: string; value: any }
    | { type: "replace"; values: FormState };

  const formReducer = useCallback((state: FormState, action: Action): FormState => {
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
  }, []);

  const [values, dispatch] = useReducer(
    formReducer,
    initial ?? {},
    (start) => ({ ...(start ?? {}) })
  );

  useEffect(() => {
    dispatch({ type: "replace", values: initial ?? {} });
  }, [initial]);

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
                    />
                  );
                })}
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
