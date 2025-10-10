import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { useTranslation } from "react-i18next";
import { Copy, Mail, Phone, Check } from "lucide-react";
import { Link } from "react-router-dom";
import contacts from "@/playbooks/supervisors_contacts.json";
import React from "react";

type Localized = { ru?: string; kz?: string; en?: string } | string;
function textOf(x: Localized, lang: "ru" | "kz" | "en"): string {
  if (!x) return "";
  if (typeof x === "string") return x;
  return x[lang] || x.ru || x.kz || x.en || "";
}

export function ContactsPage() {
  const { t: T, i18n } = useTranslation("common");
  const lang = (i18n.language as "ru" | "kz" | "en") || "ru";
  const list = Array.isArray(contacts) ? (contacts as any[]) : [];
  const [copiedKey, setCopiedKey] = React.useState<string | null>(null);
  const clearTimer = React.useRef<number | null>(null);

  const copy = (val: string) => {
    try {
      navigator.clipboard.writeText(val);
    } catch (e) {
      console.error("copy failed", e);
    }
  };

  const noteCopied = (key: string) => {
    setCopiedKey(key);
    if (clearTimer.current) window.clearTimeout(clearTimer.current);
    clearTimer.current = window.setTimeout(() => setCopiedKey(null), 1500);
  };

  return (
    <div className="max-w-3xl mx-auto px-4 py-6 space-y-4">
      <div className="flex items-center justify-between">
        <h2 className="text-xl font-semibold">{T("home.supervisors", { defaultValue: "Supervisors" })}</h2>
        <Link to="/">
          <Button variant="secondary">{T("common.back", { defaultValue: "Back" })}</Button>
        </Link>
      </div>
      {list.length === 0 ? (
        <div className="text-sm text-muted-foreground">
          {T("home.no_contacts", { defaultValue: "No contacts available." })}
        </div>
      ) : (
        <div className="grid gap-3">
          {list.map((c, idx) => (
            <Card key={idx} className="p-4">
              <CardContent className="px-0 pb-0 space-y-2">
                <div className="text-base font-medium">
                  {textOf(c.name as Localized, lang)}
                </div>
                {((c as any).role || (c as any).title) && (
                  <div className="text-xs text-muted-foreground">
                    {textOf(((c as any).role || (c as any).title) as Localized, lang)}
                  </div>
                )}
                <div className="flex flex-wrap gap-3 items-center text-sm">
                  {c.email && (
                    <div className="inline-flex items-center gap-1">
                      <Mail className="h-4 w-4" />
                      <a className="underline" href={`mailto:${c.email}`}>{c.email}</a>
                      <button
                        className="ml-1 text-muted-foreground hover:underline"
                        onClick={() => { copy(String(c.email)); noteCopied(`${idx}-email`); }}
                        title={T("common.copy", { defaultValue: "Copy" })}
                      >
                        <Copy className="h-4 w-4" />
                      </button>
                      {copiedKey === `${idx}-email` && (
                        <span className="ml-1 text-xs text-emerald-600 inline-flex items-center gap-1" aria-live="polite">
                          <Check className="h-3 w-3" /> {T("common.copied", { defaultValue: "Copied" })}
                        </span>
                      )}
                    </div>
                  )}
                  {c.phone && (
                    <div className="inline-flex items-center gap-1">
                      <Phone className="h-4 w-4" />
                      <a className="underline" href={`tel:${c.phone}`}>{c.phone}</a>
                      <button
                        className="ml-1 text-muted-foreground hover:underline"
                        onClick={() => { copy(String(c.phone)); noteCopied(`${idx}-phone`); }}
                        title={T("common.copy", { defaultValue: "Copy" })}
                      >
                        <Copy className="h-4 w-4" />
                      </button>
                      {copiedKey === `${idx}-phone` && (
                        <span className="ml-1 text-xs text-emerald-600 inline-flex items-center gap-1" aria-live="polite">
                          <Check className="h-3 w-3" /> {T("common.copied", { defaultValue: "Copied" })}
                        </span>
                      )}
                    </div>
                  )}
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      )}
    </div>
  );
}

export default ContactsPage;
