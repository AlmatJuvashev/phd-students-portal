import { useCallback, useMemo, useState } from "react";
import { useTranslation } from "react-i18next";

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
};

export function CollectionField({
  field,
  value,
  onChange,
  canEdit = true,
  disabled = false,
  renderField,
}: CollectionFieldProps) {
  const items = useMemo(() => (Array.isArray(value) ? value : EMPTY_ARRAY), [value]);
  const itemFields: FieldDef[] = useMemo(() => field.item_fields ?? [], [field.item_fields]);
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
    [itemFields, items],
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
          nextErrors[f.key] = `${labelText}: ${T("forms.required", "Обязательное поле")}`;
        } else if (list.length > 0) {
          normalized[f.key] = list;
        }
      } else {
        let valueToStore = rawValue;
        if (typeof valueToStore === "string") {
          valueToStore = valueToStore.trim();
        }
        if (f.required && (valueToStore === undefined || valueToStore === "")) {
          nextErrors[f.key] = `${labelText}: ${T("forms.required", "Обязательное поле")}`;
        }
        if (valueToStore !== undefined && valueToStore !== "") {
          normalized[f.key] = valueToStore;
        }
      }
      const otherKey = `${f.key}_other`;
      if (draft[otherKey] !== undefined) {
        const otherValue = typeof draft[otherKey] === "string" ? draft[otherKey].trim() : draft[otherKey];
        if (otherValue) {
          normalized[otherKey] = otherValue;
        }
      }
    });

    if (Object.keys(nextErrors).length > 0) {
      setErrors(nextErrors);
      return null;
    }

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

    return normalized;
  }, [T, draft, itemFields]);

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
    [items, onChange],
  );

  const total = items.length;
  const minItems = field.min_items ?? 0;
  const maxItems = field.max_items;

  return (
    <div className="space-y-2">
      <div className="flex items-start justify-between gap-2">
        <div className="font-medium text-sm">
          {label}
          {field.required ? <span className="text-destructive">*</span> : null}
          {itemLabel ? <div className="text-xs text-muted-foreground mt-1">{itemLabel}</div> : null}
        </div>
        {canEdit && !disabled ? (
          <Button size="sm" onClick={() => openEditor(null)} disabled={disabled}>
            {T("forms.add_item", "Добавить запись")}
          </Button>
        ) : null}
      </div>

      <Card className="p-3">
        {total === 0 ? (
          <div className="text-sm text-muted-foreground">
            {T("forms.collection_empty", "Записей пока нет.")}
          </div>
        ) : (
          <div className="space-y-2">
            <table className="w-full text-sm">
              <thead className="text-muted-foreground">
                <tr>
                  <th className="text-left font-medium w-12">#</th>
                  {previewFields.map((f) => (
                    <th key={f.key} className="text-left font-medium">
                      {pbT(f.label, f.key) || f.key}
                    </th>
                  ))}
                  {canEdit ? <th className="w-20" /> : null}
                </tr>
              </thead>
              <tbody className="divide-y">
                {items.map((item, index) => (
                  <tr key={index} className="align-top">
                    <td className="py-2 pr-2 text-muted-foreground">{index + 1}</td>
                    {previewFields.map((f) => {
                      const cell = item[f.key];
                      const otherKey = `${f.key}_other`;
                      const otherValue = item[otherKey];
                      const rendered = Array.isArray(cell)
                        ? cell.join(", ")
                        : cell ?? "";
                      const display = rendered || otherValue || "—";
                      return (
                        <td key={f.key} className="py-2 pr-2">
                          {display}
                        </td>
                      );
                    })}
                    {canEdit ? (
                      <td className="py-2 flex gap-2 justify-end">
                        <Button
                          size="sm"
                          variant="secondary"
                          onClick={() => openEditor(index)}
                          disabled={disabled}
                        >
                          {T("forms.edit", "Редактировать")}
                        </Button>
                        <Button
                          size="sm"
                          variant="destructive"
                          onClick={() => handleDelete(index)}
                          disabled={disabled}
                        >
                          {T("forms.delete", "Удалить")}
                        </Button>
                      </td>
                    ) : null}
                  </tr>
                ))}
              </tbody>
            </table>
            <div className="text-xs text-muted-foreground">
              {T("forms.collection_total", "Всего записей")}: {total}
            </div>
          </div>
        )}
      </Card>

      {(minItems || maxItems) && (
        <div className="text-xs text-muted-foreground">
          {minItems ? `${T("forms.collection_min", "Минимум")}: ${minItems}. ` : ""}
          {maxItems ? `${T("forms.collection_max", "Максимум")}: ${maxItems}.` : ""}
        </div>
      )}

      <Modal open={modalOpen} onClose={closeEditor}>
        <div className="space-y-4">
          <div>
            <div className="text-lg font-semibold">
              {editingIndex === null
                ? T("forms.collection_add_title", "Новая запись")
                : T("forms.collection_edit_title", "Редактирование записи")}
            </div>
            {itemLabel ? (
              <div className="text-sm text-muted-foreground">{itemLabel}</div>
            ) : null}
          </div>
          <div className="space-y-3 max-h-[60vh] overflow-y-auto pr-1">
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
                  <div className="text-xs text-destructive">{errors[itemField.key]}</div>
                ) : null}
              </div>
            ))}
          </div>
          <Separator />
          <div className="flex justify-end gap-2">
            <Button variant="secondary" onClick={closeEditor}>
              {T("forms.cancel", "Отмена")}
            </Button>
            <Button onClick={handleSave}>
              {T("forms.save", "Сохранить")}
            </Button>
          </div>
        </div>
      </Modal>
    </div>
  );
}
