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
}: {
  field: FieldDef & { placeholder?: any };
  value: any;
  onChange: (v: any) => void;
  disabled?: boolean;
  canEdit?: boolean;
}) {
  const { t: T } = useTranslation("common");
  if (field.type === "boolean") {
    return (
      <label className="inline-flex items-center gap-2">
        <input
          id={field.key}
          type="checkbox"
          disabled={!canEdit || disabled}
          checked={!!value}
          onChange={(e) => onChange(e.target.checked)}
        />
        <span>
          {t(field.label, field.key)} {field.required ? <span className="text-destructive">*</span> : null}
        </span>
      </label>
    );
  }

  return (
    <div className="grid gap-1">
      <Label htmlFor={field.key}>
        {t(field.label, field.key)} {field.required ? <span className="text-destructive">*</span> : null}
      </Label>
      {field.type === "textarea" || field.type === "array" ? (
        <Textarea
          id={field.key}
          disabled={!canEdit || disabled}
          placeholder={field.type === "array" ? T("forms.array_hint") : t(field.placeholder, "")}
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

