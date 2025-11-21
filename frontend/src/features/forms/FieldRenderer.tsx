import { memo } from "react";
import { useTranslation } from "react-i18next";

import { Label } from "@/components/ui/label";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { DatePicker } from "@/components/ui/date-picker";
import type { FieldDef } from "@/lib/playbook";
import { t } from "@/lib/playbook";
import { CollectionField } from "./CollectionField";

export type FieldRendererProps = {
  field: FieldDef & { placeholder?: any };
  value: any;
  onChange: (value: any) => void;
  disabled?: boolean;
  canEdit?: boolean;
  setField?: (key: string, value: any) => void;
  otherValue?: any;
};

const SelectField = memo(
  ({
    field,
    value,
    onChange,
    disabled,
    canEdit,
    setField,
    otherValue,
    T,
  }: FieldRendererProps & { T: ReturnType<typeof useTranslation>["t"] }) => {
    const opts = (field as any).options as Array<{
      value: string;
      label?: any;
    }>;
    const otherKey = (field as any).other_key || `${field.key}_other`;
    const isOther = value === "other";
    return (
      <div className="grid gap-1">
        <Label htmlFor={field.key}>
          {t(field.label, field.key)}{" "}
          {field.required ? <span className="text-destructive">*</span> : null}
        </Label>
        <select
          id={field.key}
          disabled={!canEdit || disabled}
          className="h-10 rounded-md border px-3 text-sm"
          value={value ?? ""}
          onChange={(e) => onChange(e.target.value)}
        >
          <option value="">—</option>
          {opts.map((o) => (
            <option key={o.value} value={o.value}>
              {t(o.label as any, o.value)}
            </option>
          ))}
        </select>
        {isOther && (
          <div className="mt-2 grid gap-1">
            <Label htmlFor={otherKey}>
              {T("fields.dissertation_form_other", "Другое (уточните)")}
            </Label>
            <Input
              id={otherKey}
              disabled={!canEdit || disabled}
              value={otherValue ?? ""}
              onChange={(e) => setField?.(otherKey, e.target.value)}
              placeholder={T(
                "fields.dissertation_form_other",
                "Другое (уточните)"
              )}
            />
          </div>
        )}
      </div>
    );
  }
);
SelectField.displayName = "SelectField";

const NoteField = memo(({ field }: FieldRendererProps) => {
  const labelText = t(field.label, field.key);
  const isInfoNote =
    labelText.startsWith("ℹ️") ||
    labelText.startsWith("⚠️") ||
    labelText.length > 100;

  if (isInfoNote) {
    return (
      <div className="text-sm text-muted-foreground bg-muted/50 p-3 rounded-md mt-4 mb-2">
        {labelText}
      </div>
    );
  }

  return (
    <div className="text-sm font-semibold text-foreground mt-4 mb-2 first:mt-0">
      {labelText}
    </div>
  );
});
NoteField.displayName = "NoteField";

const BooleanField = memo(
  ({ field, value, onChange, disabled, canEdit }: FieldRendererProps) => {
    const isReadOnly = !canEdit;
    const isChecked = !!value;

    if (isReadOnly) {
      return (
        <div
          className={`flex items-center justify-between gap-3 py-2 px-3 rounded-md min-w-0 ${
            isChecked
              ? "bg-green-50 dark:bg-green-900/10 border border-green-200 dark:border-green-800/30"
              : "bg-gray-50 dark:bg-gray-900/10 border border-gray-200 dark:border-gray-800/30 opacity-60"
          }`}
        >
          <span
            className={`flex-1 min-w-0 text-sm ${
              isChecked
                ? "text-green-900 dark:text-green-100"
                : "text-gray-700 dark:text-gray-400"
            }`}
          >
            {t(field.label, field.key)}{" "}
            {field.required ? (
              <span className="text-destructive">*</span>
            ) : null}
          </span>
          {isChecked && (
            <svg
              className="h-5 w-5 text-green-500 flex-shrink-0"
              fill="currentColor"
              viewBox="0 0 20 20"
            >
              <path
                fillRule="evenodd"
                d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z"
                clipRule="evenodd"
              />
            </svg>
          )}
        </div>
      );
    }

    return (
      <label className="flex items-center justify-between gap-3 cursor-pointer py-2 px-3 rounded-md hover:bg-muted/50 transition-colors min-w-0">
        <span className="flex-1 min-w-0">
          {t(field.label, field.key)}{" "}
          {field.required ? <span className="text-destructive">*</span> : null}
        </span>
        <input
          id={field.key}
          type="checkbox"
          className="h-5 w-5 accent-primary flex-shrink-0 cursor-pointer"
          disabled={disabled}
          checked={isChecked}
          onChange={(e) => onChange(e.target.checked)}
        />
      </label>
    );
  }
);
BooleanField.displayName = "BooleanField";

const DateField = memo(
  ({
    field,
    value,
    onChange,
    disabled,
    canEdit,
    T,
  }: FieldRendererProps & { T: ReturnType<typeof useTranslation>["t"] }) => (
    <div className="grid gap-1">
      <Label htmlFor={field.key}>
        {t(field.label, field.key)}{" "}
        {field.required ? <span className="text-destructive">*</span> : null}
      </Label>
      <DatePicker
        value={value ?? ""}
        onChange={onChange}
        disabled={!canEdit || disabled}
        placeholder={T("fields.select_date", "Select a date")}
      />
    </div>
  )
);
DateField.displayName = "DateField";

const TextField = memo(
  ({
    field,
    value,
    onChange,
    disabled,
    canEdit,
    T,
  }: FieldRendererProps & { T: ReturnType<typeof useTranslation>["t"] }) => (
    <div className="grid gap-1">
      <Label htmlFor={field.key}>
        {t(field.label, field.key)}{" "}
        {field.required ? <span className="text-destructive">*</span> : null}
      </Label>
      {field.type === "textarea" || field.type === "array" ? (
        <Textarea
          id={field.key}
          disabled={!canEdit || disabled}
          placeholder={
            field.type === "array"
              ? T("forms.array_hint")
              : t(field.placeholder, "")
          }
          value={value ?? ""}
          onChange={(e) => onChange(e.target.value)}
        />
      ) : (
        <Input
          id={field.key}
          disabled={!canEdit || disabled}
          type={field.type === "number" ? "number" : "text"}
          placeholder={t(field.placeholder, "")}
          value={value ?? ""}
          onChange={(e) => onChange(e.target.value)}
        />
      )}
    </div>
  )
);
TextField.displayName = "TextField";

function renderField(
  props: FieldRendererProps,
  T: ReturnType<typeof useTranslation>["t"]
) {
  const { field } = props;

  if (field.type === "collection") {
    return (
      <CollectionField
        field={field}
        value={props.value}
        onChange={props.onChange}
        canEdit={props.canEdit}
        disabled={props.disabled}
        renderField={(childProps) => renderField(childProps, T)}
      />
    );
  }

  if (field.type === "select" && Array.isArray((field as any).options)) {
    return <SelectField {...props} T={T} />;
  }

  if (field.type === "note") {
    return <NoteField {...props} />;
  }

  if (field.type === "boolean") {
    return <BooleanField {...props} />;
  }

  if (field.type === "date") {
    return <DateField {...props} T={T} />;
  }

  return <TextField {...props} T={T} />;
}

export const FieldRenderer = memo(function FieldRenderer(
  props: FieldRendererProps
) {
  const { t: T } = useTranslation("common");
  return renderField(props, T);
});

FieldRenderer.displayName = "FieldRenderer";
