import { createContext, useContext, useCallback, useMemo, useReducer, useEffect, type ReactNode } from "react";
import { FieldDef, NodeVM } from "@/lib/playbook";
import { evalVisible as evalVisibleExpr } from "./Visibility";

/**
 * Unified form context for all form components
 */
export interface FormContextValue {
  node: NodeVM;
  values: Record<string, any>;
  initial: Record<string, any>;
  setField: (key: string, value: any) => void;
  setValues: (values: Record<string, any>) => void;
  canEdit: boolean;
  disabled: boolean;
  readOnly: boolean;
  isValid: boolean;
  fields: FieldDef[];
  submit: (extra?: Record<string, any>) => void;
  saveDraft: (extra?: Record<string, any>) => void;
  evalVisible: (expr?: string) => boolean;
  getNextOnComplete: () => string | undefined;
}

const FormContext = createContext<FormContextValue | null>(null);

/**
 * Hook to access form context
 */
export function useFormContext() {
  const ctx = useContext(FormContext);
  if (!ctx) {
    throw new Error("useFormContext must be used within FormProvider");
  }
  return ctx;
}

/**
 * Form state reducer with field-specific logic
 */
type FormState = Record<string, any>;
type FormAction =
  | { type: "set"; key: string; value: any }
  | { type: "replace"; values: FormState }
  | { type: "reset" };

function formReducer(state: FormState, action: FormAction): FormState {
  switch (action.type) {
    case "replace":
      return { ...action.values };
    case "reset":
      return {};
    case "set": {
      const next = { ...state, [action.key]: action.value };

      // Specific field logic - clear dependent fields
      // Example: hearing_happened affects remarks fields
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
}

export interface FormProviderProps {
  node: NodeVM;
  initial?: Record<string, any>;
  onSubmit?: (payload: any) => void;
  canEdit?: boolean;
  disabled?: boolean;
  children: ReactNode;
}

/**
 * Unified FormProvider component
 * Manages form state, validation, and submission logic
 */
export function FormProvider({
  node,
  initial = {},
  onSubmit,
  canEdit: canEditProp,
  disabled = false,
  children,
}: FormProviderProps) {
  const fields: FieldDef[] = node.requirements?.fields ?? [];

  const [values, dispatch] = useReducer(
    formReducer,
    initial ?? {},
    (start) => ({ ...(start ?? {}) })
  );

  // Sync with initial values when they change
  useEffect(() => {
    dispatch({ type: "replace", values: initial ?? {} });
  }, [initial]);

  // Determine readonly state
  const submittedAt: string | undefined =
    (initial as any)?.__submittedAt || values?.__submittedAt;
  
  const readOnly =
    canEditProp !== undefined
      ? !canEditProp
      : node.state === "submitted" ||
        node.state === "done" ||
        Boolean(submittedAt);

  const canEdit = canEditProp !== undefined ? canEditProp : !readOnly;

  // Field setters
  const setField = useCallback(
    (key: string, value: any) => {
      if (!readOnly) {
        dispatch({ type: "set", key, value });
      }
    },
    [readOnly]
  );

  const setValues = useCallback(
    (newValues: Record<string, any>) => {
      if (!readOnly) {
        dispatch({ type: "replace", values: newValues });
      }
    },
    [readOnly]
  );

  // Visibility evaluation
  const evalVisible = useCallback(
    (expr?: string) => evalVisibleExpr(values, expr),
    [values]
  );

  // Validation logic
  const isValid = useMemo(() => {
    const requiredFields = fields.filter((f) => f.required);
    return requiredFields.every((f) => {
      const value = values[f.key];
      
      // Check visibility first (if visible property exists)
      const fAny = f as any;
      if (fAny.visible !== undefined && !evalVisible(fAny.visible)) {
        return true; // Hidden fields don't need to be filled
      }

      // Boolean fields
      if (f.type === "boolean") {
        return !!value;
      }

      // Text/select/date fields
      if (f.type === "text" || f.type === "select" || f.type === "date") {
        return value !== undefined && value !== null && value !== "";
      }

      return true;
    });
  }, [fields, values, evalVisible]);

  // Calculate next node based on outcomes
  const getNextOnComplete = useCallback(() => {
    if (node.outcomes && node.outcomes.length > 0) {
      // Find the first outcome that matches current form state
      for (const outcome of node.outcomes) {
        const outcomeAny = outcome as any;
        
        if (outcomeAny.when) {
          // Extract required fields from when condition
          // e.g., "form.field1 && form.field2" -> ["field1", "field2"]
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
  }, [node, values]);

  // Submission handlers
  const submit = useCallback(
    (extra: Record<string, any> = {}) => {
      if (!onSubmit || !isValid) return;
      
      const nextOnComplete = getNextOnComplete();
      const payload = {
        ...values,
        ...extra,
        __submittedAt: new Date().toISOString(),
        ...(nextOnComplete ? { __nextOverride: nextOnComplete } : {}),
      };
      
      onSubmit(payload);
    },
    [onSubmit, values, isValid, getNextOnComplete]
  );

  const saveDraft = useCallback(
    (extra: Record<string, any> = {}) => {
      if (!onSubmit) return;
      onSubmit({ ...values, ...extra, __draft: true });
    },
    [onSubmit, values]
  );

  const contextValue = useMemo<FormContextValue>(
    () => ({
      node,
      values,
      initial,
      setField,
      setValues,
      canEdit,
      disabled,
      readOnly,
      isValid,
      fields,
      submit,
      saveDraft,
      evalVisible,
      getNextOnComplete,
    }),
    [
      node,
      values,
      initial,
      setField,
      setValues,
      canEdit,
      disabled,
      readOnly,
      isValid,
      fields,
      submit,
      saveDraft,
      evalVisible,
      getNextOnComplete,
    ]
  );

  return (
    <FormContext.Provider value={contextValue}>
      {children}
    </FormContext.Provider>
  );
}
