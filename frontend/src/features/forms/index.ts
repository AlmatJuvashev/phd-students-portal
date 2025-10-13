/**
 * Unified Form System
 *
 * Phase 1.4 - Form Unification
 * Provides reusable components and hooks for form management
 */

export {
  FormProvider,
  useFormContext,
  type FormContextValue,
  type FormProviderProps,
} from "./FormProvider";
export { FormFields, type FormFieldsProps } from "./FormFields";
export { FormActions, type FormActionsProps } from "./FormActions";
export { FormStatus, type FormStatusProps } from "./FormStatus";

// Re-export existing components for convenience
export { FieldRenderer, type FieldRendererProps } from "./FieldRenderer";
export { ActionsBar } from "./ActionsBar";
export { ConfirmModal } from "./ConfirmModal";
export { TemplatesPanel } from "./TemplatesPanel";
