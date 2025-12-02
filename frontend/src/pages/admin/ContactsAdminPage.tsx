import React from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { Contact, ContactPayload, createContact, deleteContact, fetchAdminContacts, updateContact } from "@/api/contacts";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Badge } from "@/components/ui/badge";
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { Switch } from "@/components/ui/switch";
import { useTranslation } from "react-i18next";
import { Plus, Pencil, Trash2, Loader2, Phone, Mail, RefreshCw } from "lucide-react";
import { Checkbox } from "@/components/ui/checkbox";

type ContactForm = {
  nameRu: string;
  nameKz: string;
  nameEn: string;
  titleRu: string;
  titleKz: string;
  titleEn: string;
  email: string;
  phone: string;
  sortOrder: string;
  isActive: boolean;
};

const emptyForm: ContactForm = {
  nameRu: "",
  nameKz: "",
  nameEn: "",
  titleRu: "",
  titleKz: "",
  titleEn: "",
  email: "",
  phone: "",
  sortOrder: "",
  isActive: true,
};

function compactMap(values: Record<string, string>): Record<string, string> | undefined {
  const out: Record<string, string> = {};
  Object.entries(values).forEach(([k, v]) => {
    if (v && v.trim() !== "") out[k] = v.trim();
  });
  return Object.keys(out).length ? out : undefined;
}

function toForm(c?: Contact): ContactForm {
  if (!c) return { ...emptyForm };
  return {
    nameRu: c.name?.ru || "",
    nameKz: c.name?.kz || "",
    nameEn: c.name?.en || "",
    titleRu: c.title?.ru || "",
    titleKz: c.title?.kz || "",
    titleEn: c.title?.en || "",
    email: c.email || "",
    phone: c.phone || "",
    sortOrder: c.sort_order != null ? String(c.sort_order) : "",
    isActive: c.is_active !== false,
  };
}

function buildPayload(form: ContactForm): ContactPayload {
  const name = compactMap({ ru: form.nameRu, kz: form.nameKz, en: form.nameEn });
  const title = compactMap({ ru: form.titleRu, kz: form.titleKz, en: form.titleEn });
  const payload: ContactPayload = {
    name: name || { ru: "" },
    title,
    email: form.email || undefined,
    phone: form.phone || undefined,
    sort_order: form.sortOrder ? Number(form.sortOrder) : undefined,
    is_active: form.isActive,
  };
  return payload;
}

export function ContactsAdminPage() {
  const { t } = useTranslation("common");
  const qc = useQueryClient();
  const [includeInactive, setIncludeInactive] = React.useState(true);
  const [showDialog, setShowDialog] = React.useState(false);
  const [form, setForm] = React.useState<ContactForm>({ ...emptyForm });
  const [editing, setEditing] = React.useState<Contact | null>(null);
  const [validationError, setValidationError] = React.useState<string | null>(null);

  const contactsQuery = useQuery<Contact[]>({
    queryKey: ["admin", "contacts", includeInactive],
    queryFn: () => fetchAdminContacts(includeInactive),
  });

  const creating = useMutation({
    mutationFn: (payload: ContactPayload) => createContact(payload),
    onSuccess: () => {
      setShowDialog(false);
      setForm({ ...emptyForm });
      qc.invalidateQueries({ queryKey: ["admin", "contacts"] });
    },
  });

  const updating = useMutation({
    mutationFn: ({ id, payload }: { id: string; payload: ContactPayload }) =>
      updateContact(id, payload),
    onSuccess: () => {
      setShowDialog(false);
      setEditing(null);
      setForm({ ...emptyForm });
      qc.invalidateQueries({ queryKey: ["admin", "contacts"] });
    },
  });

  const removing = useMutation({
    mutationFn: (id: string) => deleteContact(id),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["admin", "contacts"] }),
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    const nameMap = compactMap({ ru: form.nameRu, kz: form.nameKz, en: form.nameEn });
    if (!nameMap) {
      setValidationError(t("admin.contacts.errors.name_required", { defaultValue: "Name is required" }));
      return;
    }
    setValidationError(null);
    const payload = buildPayload(form);
    if (editing) {
      updating.mutate({ id: editing.id, payload });
    } else {
      creating.mutate(payload);
    }
  };

  const onEdit = (c: Contact) => {
    setEditing(c);
    setForm(toForm(c));
    setShowDialog(true);
  };

  const onCreate = () => {
    setEditing(null);
    setForm({ ...emptyForm, isActive: true });
    setShowDialog(true);
    setValidationError(null);
  };

  const sorted = (contactsQuery.data || []).slice().sort((a, b) => {
    const orderA = a.sort_order ?? 0;
    const orderB = b.sort_order ?? 0;
    if (orderA === orderB) return (a.name?.ru || "").localeCompare(b.name?.ru || "");
    return orderA - orderB;
  });

  const statusBadge = (active: boolean) => (
    <Badge variant={active ? "secondary" : "outline"} className={active ? "bg-emerald-100 text-emerald-800" : ""}>
      {active
        ? t("admin.forms.status_active", { defaultValue: "Active" })
        : t("admin.forms.status_inactive", { defaultValue: "Inactive" })}
    </Badge>
  );

  return (
    <div className="space-y-4">
      <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h1 className="text-2xl font-bold">{t("admin.contacts.title", { defaultValue: "Contacts" })}</h1>
          <p className="text-muted-foreground text-sm">
            {t("admin.contacts.subtitle", { defaultValue: "Manage visible contacts for students" })}
          </p>
        </div>
        <div className="flex items-center gap-2">
          <label className="flex items-center gap-2 text-sm">
            <Checkbox
              checked={includeInactive}
              onCheckedChange={(v) => setIncludeInactive(Boolean(v))}
              aria-label={t("admin.contacts.show_inactive", { defaultValue: "Show inactive" })}
            />
            {t("admin.contacts.show_inactive", { defaultValue: "Show inactive" })}
          </label>
          <Button onClick={contactsQuery.refetch} variant="outline" size="sm">
            <RefreshCw className="h-4 w-4 mr-2" />
            {t("common.refresh", { defaultValue: "Refresh" })}
          </Button>
          <Button onClick={onCreate} className="gap-2">
            <Plus className="h-4 w-4" />
            {t("admin.contacts.add", { defaultValue: "Add contact" })}
          </Button>
        </div>
      </div>

      <Card>
        <CardHeader>
          <CardTitle className="text-base">
            {t("admin.contacts.table_title", { defaultValue: "Contacts" })} · {sorted.length}
          </CardTitle>
        </CardHeader>
        <CardContent className="space-y-3">
          {contactsQuery.isLoading ? (
            <div className="py-8 text-center text-muted-foreground">
              <Loader2 className="h-5 w-5 animate-spin inline-block mr-2" />
              {t("common.loading", { defaultValue: "Loading..." })}
            </div>
          ) : sorted.length === 0 ? (
            <div className="py-8 text-center text-muted-foreground">
              {t("admin.contacts.empty", { defaultValue: "No contacts yet." })}
            </div>
          ) : (
            <div className="overflow-x-auto">
              <table className="min-w-full text-sm">
                <thead className="bg-muted/30">
                  <tr>
                    <th className="px-3 py-2 text-left w-12">#</th>
                    <th className="px-3 py-2 text-left">{t("table.name", { defaultValue: "Name" })}</th>
                    <th className="px-3 py-2 text-left">{t("admin.contacts.title_col", { defaultValue: "Title" })}</th>
                    <th className="px-3 py-2 text-left">{t("admin.contacts.phone", { defaultValue: "Phone" })}</th>
                    <th className="px-3 py-2 text-left">{t("admin.contacts.email", { defaultValue: "Email" })}</th>
                    <th className="px-3 py-2 text-left">{t("admin.contacts.sort", { defaultValue: "Sort" })}</th>
                    <th className="px-3 py-2 text-left">{t("admin.contacts.status", { defaultValue: "Status" })}</th>
                    <th className="px-3 py-2 text-right">{t("table.actions", { defaultValue: "Actions" })}</th>
                  </tr>
                </thead>
                <tbody>
                  {sorted.map((c, idx) => (
                    <tr key={c.id} className="border-b">
                      <td className="px-3 py-2 text-muted-foreground">{idx + 1}</td>
                      <td className="px-3 py-2">
                        <div className="font-medium">{c.name?.ru || c.name?.en || c.name?.kz || "-"}</div>
                        {c.title && (c.title.ru || c.title.en || c.title.kz) && (
                          <div className="text-xs text-muted-foreground">
                            {c.title.ru || c.title.en || c.title.kz}
                          </div>
                        )}
                      </td>
                      <td className="px-3 py-2 text-muted-foreground">
                        {c.title?.ru || c.title?.en || c.title?.kz || "—"}
                      </td>
                      <td className="px-3 py-2">
                        {c.phone ? (
                          <span className="inline-flex items-center gap-2">
                            <Phone className="h-4 w-4 text-muted-foreground" />
                            {c.phone}
                          </span>
                        ) : (
                          "—"
                        )}
                      </td>
                      <td className="px-3 py-2">
                        {c.email ? (
                          <span className="inline-flex items-center gap-2">
                            <Mail className="h-4 w-4 text-muted-foreground" />
                            {c.email}
                          </span>
                        ) : (
                          "—"
                        )}
                      </td>
                      <td className="px-3 py-2">{c.sort_order ?? 0}</td>
                      <td className="px-3 py-2">{statusBadge(c.is_active !== false)}</td>
                      <td className="px-3 py-2 text-right">
                        <div className="flex justify-end gap-2">
                          <Button variant="ghost" size="sm" onClick={() => onEdit(c)}>
                            <Pencil className="h-4 w-4 mr-1" />
                            {t("common.edit", { defaultValue: "Edit" })}
                          </Button>
                          {c.is_active !== false ? (
                            <Button
                              variant="ghost"
                              size="sm"
                              onClick={() => removing.mutate(c.id)}
                              disabled={removing.isPending}
                              className="text-destructive"
                            >
                              <Trash2 className="h-4 w-4 mr-1" />
                              {t("admin.forms.status_inactive", { defaultValue: "Inactive" })}
                            </Button>
                          ) : (
                            <Button
                              variant="ghost"
                              size="sm"
                              onClick={() =>
                                updating.mutate({
                                  id: c.id,
                                  payload: { name: c.name, is_active: true },
                                })
                              }
                              disabled={updating.isPending}
                            >
                              {t("admin.forms.status_active", { defaultValue: "Active" })}
                            </Button>
                          )}
                        </div>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </CardContent>
      </Card>

      <Dialog open={showDialog} onOpenChange={setShowDialog}>
        <DialogContent className="max-w-3xl">
          <DialogHeader>
            <DialogTitle>
              {editing
                ? t("admin.contacts.edit_title", { defaultValue: "Edit contact" })
                : t("admin.contacts.add", { defaultValue: "Add contact" })}
            </DialogTitle>
          </DialogHeader>
          <form className="space-y-4" onSubmit={handleSubmit}>
            <div className="grid grid-cols-1 md:grid-cols-3 gap-3">
              <div>
                <label className="text-xs text-muted-foreground">{t("admin.contacts.name_ru", { defaultValue: "Name (RU)" })}</label>
                <Input value={form.nameRu} onChange={(e) => setForm((f) => ({ ...f, nameRu: e.target.value }))} />
              </div>
              <div>
                <label className="text-xs text-muted-foreground">{t("admin.contacts.name_kz", { defaultValue: "Name (KZ)" })}</label>
                <Input value={form.nameKz} onChange={(e) => setForm((f) => ({ ...f, nameKz: e.target.value }))} />
              </div>
              <div>
                <label className="text-xs text-muted-foreground">{t("admin.contacts.name_en", { defaultValue: "Name (EN)" })}</label>
                <Input value={form.nameEn} onChange={(e) => setForm((f) => ({ ...f, nameEn: e.target.value }))} />
              </div>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-3 gap-3">
              <div>
                <label className="text-xs text-muted-foreground">{t("admin.contacts.title_ru", { defaultValue: "Title (RU)" })}</label>
                <Input value={form.titleRu} onChange={(e) => setForm((f) => ({ ...f, titleRu: e.target.value }))} />
              </div>
              <div>
                <label className="text-xs text-muted-foreground">{t("admin.contacts.title_kz", { defaultValue: "Title (KZ)" })}</label>
                <Input value={form.titleKz} onChange={(e) => setForm((f) => ({ ...f, titleKz: e.target.value }))} />
              </div>
              <div>
                <label className="text-xs text-muted-foreground">{t("admin.contacts.title_en", { defaultValue: "Title (EN)" })}</label>
                <Input value={form.titleEn} onChange={(e) => setForm((f) => ({ ...f, titleEn: e.target.value }))} />
              </div>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-3 gap-3">
              <div>
                <label className="text-xs text-muted-foreground">{t("admin.contacts.phone", { defaultValue: "Phone" })}</label>
                <Input value={form.phone} onChange={(e) => setForm((f) => ({ ...f, phone: e.target.value }))} />
              </div>
              <div>
                <label className="text-xs text-muted-foreground">{t("admin.contacts.email", { defaultValue: "Email" })}</label>
                <Input value={form.email} onChange={(e) => setForm((f) => ({ ...f, email: e.target.value }))} />
              </div>
              <div>
                <label className="text-xs text-muted-foreground">{t("admin.contacts.sort", { defaultValue: "Sort order" })}</label>
                <Input
                  type="number"
                  value={form.sortOrder}
                  onChange={(e) => setForm((f) => ({ ...f, sortOrder: e.target.value }))}
                />
              </div>
            </div>

            <div className="flex items-center gap-3">
              <Switch
                id="contact-active"
                checked={form.isActive}
                onCheckedChange={(v) => setForm((f) => ({ ...f, isActive: Boolean(v) }))}
              />
              <label htmlFor="contact-active" className="text-sm">
                {t("admin.forms.status_active", { defaultValue: "Active" })}
              </label>
            </div>

            {validationError && <div className="text-xs text-red-500">{validationError}</div>}

            <div className="flex justify-end gap-2 pt-2">
              <Button type="button" variant="outline" onClick={() => setShowDialog(false)}>
                {t("common.cancel", { defaultValue: "Cancel" })}
              </Button>
              <Button type="submit" disabled={creating.isPending || updating.isPending}>
                {(creating.isPending || updating.isPending) && <Loader2 className="h-4 w-4 animate-spin mr-2" />}
                {editing
                  ? t("common.save", { defaultValue: "Save" })
                  : t("admin.contacts.add", { defaultValue: "Add contact" })}
              </Button>
            </div>
          </form>
        </DialogContent>
      </Dialog>
    </div>
  );
}

export default ContactsAdminPage;
