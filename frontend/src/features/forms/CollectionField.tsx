import { useCallback, useMemo, useState } from "react";
import { useTranslation } from "react-i18next";
import { Globe, BookOpen, Users, Lightbulb } from "lucide-react";

import { Button } from "@/components/ui/button";
import { Modal } from "@/components/ui/modal";
import { Card } from "@/components/ui/card";
import { Separator } from "@/components/ui/separator";
import { t as pbT, type FieldDef } from "@/lib/playbook";
import type { FieldRendererProps } from "./FieldRenderer";

const splitCoauthors = (raw: unknown): string[] => {
  if (Array.isArray(raw)) {
    return raw
      .map((item) => (typeof item === "string" ? item.trim() : ""))
      .filter(Boolean);
  }
  if (typeof raw !== "string") {
    return [];
  }
  return raw
    .split(/[\n;,]+/)
    .map((item) => item.trim())
    .filter(Boolean);
};

const prepareDraftValue = (value: any, field: FieldDef) => {
  if (field.type === "array") {
    if (Array.isArray(value)) {
      return value.join("\n");
    }
    if (typeof value === "string") {
      return value;
    }
    return "";
  }
  if (typeof value === "string") {
    return value;
  }
  return value ?? "";
};

const EMPTY_ARRAY: any[] = [];

export type CollectionFieldProps = {
  field: FieldDef;
  value: any;
  onChange: (value: any) => void;
  canEdit?: boolean;
  disabled?: boolean;
  renderField: (props: FieldRendererProps) => JSX.Element;
  itemErrors?: Record<number, string>;
};

export function CollectionField({
  field,
  value,
  onChange,
  canEdit = true,
  disabled = false,
  renderField,
  itemErrors,
}: CollectionFieldProps) {
  const items = useMemo(
    () => (Array.isArray(value) ? value : EMPTY_ARRAY),
    [value]
  );
  const itemFields: FieldDef[] = useMemo(
    () => field.item_fields ?? [],
    [field.item_fields]
  );
  const { t: T } = useTranslation("common");

  const [modalOpen, setModalOpen] = useState(false);
  const [editingIndex, setEditingIndex] = useState<number | null>(null);
  const [draft, setDraft] = useState<Record<string, any>>({});
  const [errors, setErrors] = useState<Record<string, string>>({});

  const label = pbT(field.label, field.key);
  const itemLabel = pbT(field.item_label, field.key);

  const previewFields = useMemo(() => {
    return itemFields.filter((f) => f.type !== "note").slice(0, 4);
  }, [itemFields]);

  const openEditor = useCallback(
    (index: number | null) => {
      const base: Record<string, any> = {};
      if (index !== null && items[index]) {
        const item = items[index];
        itemFields.forEach((f) => {
          base[f.key] = prepareDraftValue(item[f.key], f);
          const otherKey = `${f.key}_other`;
          if (item[otherKey] !== undefined) {
            base[otherKey] = item[otherKey];
          }
        });
        Object.keys(item).forEach((key) => {
          if (base[key] === undefined) {
            base[key] = item[key];
          }
        });
      } else {
        itemFields.forEach((f) => {
          base[f.key] = f.type === "array" ? "" : "";
        });
      }
      setDraft(base);
      setErrors({});
      setEditingIndex(index);
      setModalOpen(true);
    },
    [itemFields, items]
  );

  const closeEditor = useCallback(() => {
    setModalOpen(false);
    setEditingIndex(null);
    setDraft({});
    setErrors({});
  }, []);

  const setDraftField = useCallback((key: string, val: any) => {
    setDraft((prev) => ({ ...prev, [key]: val }));
  }, []);

  const validateAndNormalize = useCallback(() => {
    const normalized: Record<string, any> = {};
    const nextErrors: Record<string, string> = {};
    const fieldKeys = new Set(itemFields.map((f) => f.key));

    itemFields.forEach((f) => {
      const rawValue = draft[f.key];
      const labelText = pbT(f.label, f.key) || f.key;
      if (f.type === "array" || f.key === "coauthors") {
        const list = splitCoauthors(rawValue);
        if (f.required && list.length === 0) {
          nextErrors[f.key] = `${labelText}: ${T(
            "forms.required",
            "Обязательное поле"
          )}`;
        } else if (list.length > 0) {
          normalized[f.key] = list;
        }
      } else {
        let valueToStore = rawValue;
        if (typeof valueToStore === "string") {
          valueToStore = valueToStore.trim();
        }
        if (f.required && (valueToStore === undefined || valueToStore === "")) {
          nextErrors[f.key] = `${labelText}: ${T(
            "forms.required",
            "Обязательное поле"
          )}`;
        }
        if (valueToStore !== undefined && valueToStore !== "") {
          normalized[f.key] = valueToStore;
        }
      }
      const otherKey = `${f.key}_other`;
      if (draft[otherKey] !== undefined) {
        const otherValue =
          typeof draft[otherKey] === "string"
            ? draft[otherKey].trim()
            : draft[otherKey];
        if (otherValue) {
          normalized[otherKey] = otherValue;
        }
      }
    });

    // Preserve any additional keys from draft (like IDs injected later)
    Object.keys(draft).forEach((key) => {
      if (
        normalized[key] === undefined &&
        !key.endsWith("_other") &&
        !fieldKeys.has(key)
      ) {
        normalized[key] = draft[key];
      }
    });

    if (Object.keys(nextErrors).length > 0) {
      setErrors(nextErrors);
      return null;
    }

    return normalized;
  }, [draft, itemFields, T]);

  const handleSave = useCallback(() => {
    const normalized = validateAndNormalize();
    if (!normalized) return;

    const next = Array.isArray(items) ? [...items] : [];
    if (editingIndex !== null) {
      next[editingIndex] = { ...next[editingIndex], ...normalized };
    } else {
      next.push(normalized);
    }
    setErrors({});
    onChange(next);
    closeEditor();
  }, [closeEditor, editingIndex, items, onChange, validateAndNormalize]);

  const handleDelete = useCallback(
    (idx: number) => {
      const next = items.filter((_, index) => index !== idx);
      onChange(next);
    },
    [items, onChange]
  );

  const total = items.length;
  const minItems = field.min_items ?? 0;
  const maxItems = field.max_items;

  const iconMap: Record<string, any> = {
    "lucide:globe": Globe,
    "lucide:book-open": BookOpen,
    "lucide:users": Users,
    "lucide:lightbulb": Lightbulb,
  };

  const SectionIcon = (field as any).icon ? iconMap[(field as any).icon] : null;

  return (
    <div className="space-y-6 md:space-y-8">
      <div className="flex items-center gap-3 md:static md:bg-transparent md:backdrop-blur-none">
        {SectionIcon && (
          <SectionIcon className="w-5 h-5 text-primary flex-shrink-0" />
        )}
        <h3 className="text-lg font-semibold">
          {label}
          {field.required ? <span className="text-destructive">*</span> : null}
        </h3>
      </div>
      {itemLabel ? (
        <div className="text-sm text-muted-foreground">{itemLabel}</div>
      ) : null}

      {total === 0 ? (
        <Card className="p-4">
          <div className="text-sm text-muted-foreground">
            {T("forms.collection_empty", "Записей пока нет.")}
          </div>
        </Card>
      ) : (
        <div className="space-y-4 pt-2 md:pt-0">
          {items.map((item, index) => {
            const error = itemErrors?.[index];
            return (
              <div
                key={index}
                className={`bg-card border rounded-lg p-4 md:p-6 space-y-4 ${
                  error ? "border-destructive ring-1 ring-destructive" : ""
                }`}
              >
                <div className="flex justify-between items-start">
                  <div className="text-xs text-muted-foreground">#{index + 1}</div>
                  {error && (
                    <div className="text-xs text-destructive font-medium">{error}</div>
                  )}
                </div>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                {previewFields.map((f) => {
                  const cell = item[f.key];
                  const otherKey = `${f.key}_other`;
                  const otherValue = (item as any)[otherKey];
                  const rendered = Array.isArray(cell)
                    ? cell.join(", ")
                    : cell ?? "";
                  const display = (rendered || otherValue || "—") as string;
                  const isTitle = f.key === "title";
                  return (
                    <div
                      key={f.key}
                      className={`space-y-1 ${isTitle ? "md:col-span-2" : ""}`}
                    >
                      <div className="text-xs text-muted-foreground">
                        {pbT(f.label, f.key) || f.key}
                      </div>
                      <div className="text-sm break-words">{display}</div>
                    </div>
                  );
                })}
              </div>

              {canEdit ? (
                <div className="flex flex-col sm:flex-row gap-2 pt-4 border-t">
                  <Button
                    size="sm"
                    variant="secondary"
                    className="w-full sm:w-auto min-h-[44px]"
                    onClick={() => openEditor(index)}
                    disabled={disabled}
                  >
                    {T("forms.edit", "Редактировать")}
                  </Button>
                  <Button
                    size="sm"
                    variant="destructive"
                    className="w-full sm:w-auto min-h-[44px]"
                    onClick={() => handleDelete(index)}
                    disabled={disabled}
                  >
                    {T("forms.delete", "Удалить")}
                  </Button>
                </div>
              ) : null}
            </div>
            );
          })}
          <div className="text-xs text-muted-foreground">
            {T("forms.collection_total", "Всего записей")}: {total}
          </div>
        </div>
      )}

      {(minItems || maxItems) && (
        <div className="text-xs text-muted-foreground">
          {minItems
            ? `${T("forms.collection_min", "Минимум")}: ${minItems}. `
            : ""}
          {maxItems
            ? `${T("forms.collection_max", "Максимум")}: ${maxItems}.`
            : ""}
        </div>
      )}

      {canEdit && !disabled ? (
        <Button
          variant="outline"
          className="w-full mt-2 border-dashed hover:border-solid min-h-[44px]"
          onClick={() => openEditor(null)}
          disabled={disabled}
        >
          {T("forms.add_item", "Добавить запись")}
        </Button>
      ) : null}

      <Modal open={modalOpen} onClose={closeEditor}>
        <div className="max-h-[85vh] flex flex-col">
          <div className="pb-4 border-b mb-4">
            <div className="text-lg font-semibold">
              {editingIndex === null
                ? T("forms.collection_add_title", "Новая запись")
                : T("forms.collection_edit_title", "Редактирование записи")}
            </div>
            {itemLabel ? (
              <div className="text-sm text-muted-foreground">{itemLabel}</div>
            ) : null}
          </div>
          <div className="flex-1 space-y-4 overflow-y-auto pr-2">
            {itemFields.map((itemField) => (
              <div key={itemField.key} className="space-y-1">
                {renderField({
                  field: itemField as FieldDef & { placeholder?: any },
                  value: draft[itemField.key],
                  onChange: (val: any) => setDraftField(itemField.key, val),
                  canEdit: true,
                  disabled: false,
                  setField: setDraftField,
                  otherValue: draft[`${itemField.key}_other`],
                })}
                {errors[itemField.key] ? (
                  <div className="text-xs text-destructive">
                    {errors[itemField.key]}
                  </div>
                ) : null}
              </div>
            ))}
          </div>
          <Separator />
          <div className="flex justify-end gap-2 pt-2">
            <Button
              variant="secondary"
              onClick={closeEditor}
              className="min-h-[44px]"
            >
              {T("forms.cancel", "Отмена")}
            </Button>
            <Button onClick={handleSave} className="min-h-[44px]">
              {T("forms.save", "Сохранить")}
            </Button>
          </div>
        </div>
      </Modal>
    </div>
  );
}
