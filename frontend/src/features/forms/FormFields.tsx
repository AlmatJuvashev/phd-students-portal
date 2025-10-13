import { useFormContext } from "./FormProvider";
import { FieldRenderer } from "./FieldRenderer";
import { ChecklistItem } from "@/components/ui/checklist-item";
import { t } from "@/lib/playbook";

export interface FormFieldsProps {
  /**
   * Render only specific field types
   * If not provided, renders all fields
   */
  types?: ("boolean" | "text" | "select" | "date")[];
  
  /**
   * Custom className for the container
   */
  className?: string;
  
  /**
   * Spacing between fields
   */
  spacing?: "tight" | "normal" | "loose";
}

/**
 * Automatic field renderer using FormProvider context
 * Supports boolean (checklist), text, select, and date fields
 */
export function FormFields({
  types,
  className = "",
  spacing = "normal",
}: FormFieldsProps) {
  const { fields, values, setField, readOnly, disabled, evalVisible } = useFormContext();

  const spacingClass = {
    tight: "space-y-2",
    normal: "space-y-3",
    loose: "space-y-4",
  }[spacing];

  const filteredFields = types
    ? fields.filter((f) => types.includes(f.type as any))
    : fields;

  return (
    <div className={`${spacingClass} ${className} min-w-0`}>
      {filteredFields.map((f) => {
        const isChecked = !!values[f.key];
        const fieldValue = values[f.key];

        // Check visibility
        const fAny = f as any;
        if (fAny.visible !== undefined && !evalVisible(fAny.visible)) {
          return null; // Skip hidden fields
        }

        // Boolean fields use ChecklistItem component
        if (f.type === "boolean") {
          return (
            <ChecklistItem
              key={f.key}
              checked={isChecked}
              onChange={(checked) => setField(f.key, checked)}
              label={`${t(f.label, f.key)}${f.required ? "*" : ""}`}
              readOnly={readOnly}
              disabled={disabled}
            />
          );
        }

        // Other field types use FieldRenderer
        return (
          <FieldRenderer
            key={f.key}
            field={f}
            value={fieldValue}
            onChange={(v) => setField(f.key, v)}
            setField={setField}
            canEdit={!readOnly}
            disabled={disabled}
          />
        );
      })}
    </div>
  );
}
