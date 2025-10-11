import { Label } from "@/components/ui/label";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import type { FieldDef } from "@/lib/playbook";
import { t } from "@/lib/playbook";
import { useTranslation } from "react-i18next";

export function FieldRenderer({
  field,
  value,
  onChange,
  disabled,
  canEdit = true,
  setField,
  otherValue,
}: {
  field: FieldDef & { placeholder?: any };
  value: any;
  onChange: (v: any) => void;
  disabled?: boolean;
  canEdit?: boolean;
  setField?: (k: string, v: any) => void;
  otherValue?: any;
}) {
  const { t: T } = useTranslation("common");
  // Select support
  if (field.type === "select" && Array.isArray((field as any).options)) {
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
  // Note field - used for section headers or info text
  if (field.type === "note") {
    const labelText = t(field.label, field.key);
    // Check if it's an info note (starts with emoji or contains multiple sentences)
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

    // Regular section header
    return (
      <div className="text-sm font-semibold text-foreground mt-4 mb-2 first:mt-0">
        {labelText}
      </div>
    );
  }

  if (field.type === "boolean") {
    return (
      <label className="flex items-center justify-between gap-3 cursor-pointer py-2 px-3 rounded-md hover:bg-muted/50 transition-colors">
        <span className="flex-1">
          {t(field.label, field.key)}{" "}
          {field.required ? <span className="text-destructive">*</span> : null}
        </span>
        <input
          id={field.key}
          type="checkbox"
          className="h-5 w-5 accent-primary flex-shrink-0 cursor-pointer"
          disabled={!canEdit || disabled}
          checked={!!value}
          onChange={(e) => onChange(e.target.checked)}
        />
      </label>
    );
  }

  // Date input (native for now; shadcn-style can be plugged in later)
  if (field.type === "date") {
    return (
      <div className="grid gap-1">
        <Label htmlFor={field.key}>
          {t(field.label, field.key)}{" "}
          {field.required ? <span className="text-destructive">*</span> : null}
        </Label>
        <Input
          id={field.key}
          type="date"
          disabled={!canEdit || disabled}
          value={value ?? ""}
          onChange={(e) => onChange(e.target.value)}
        />
      </div>
    );
  }

  return (
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
  );
}
