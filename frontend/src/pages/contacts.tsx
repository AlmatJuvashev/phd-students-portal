import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { BackButton } from "@/components/ui/back-button";
import { useTranslation } from "react-i18next";
import { Copy, Mail, Phone, Check, Users } from "lucide-react";
import { Link } from "react-router-dom";
import fallbackContacts from "@/playbooks/supervisors_contacts.json";
import React from "react";
import { useQuery } from "@tanstack/react-query";
import { fetchContacts, Contact } from "@/api/contacts";

type Localized = { ru?: string; kz?: string; en?: string } | string;
function textOf(x: Localized, lang: "ru" | "kz" | "en"): string {
  if (!x) return "";
  if (typeof x === "string") return x;
  return x[lang] || x.ru || x.kz || x.en || "";
}

export function ContactsPage() {
  const { t: T, i18n } = useTranslation("common");
  const lang = (i18n.language as "ru" | "kz" | "en") || "ru";
  const { data, isLoading, error } = useQuery<Contact[]>({
    queryKey: ["contacts"],
    queryFn: fetchContacts,
  });
  const list = (data && Array.isArray(data) ? data : []) as any[];
  const fallbackList = Array.isArray(fallbackContacts)
    ? (fallbackContacts as any[])
    : [];
  const resolvedList =
    list.length > 0 ? list : error ? fallbackList : ([] as any[]);
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
    <div className="max-w-4xl mx-auto px-4 py-6 space-y-6">
      <div className="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4 bg-gradient-to-r from-primary/10 via-primary/5 to-transparent rounded-2xl p-6 border-l-4 border-primary">
        <div className="flex-1">
          <h2 className="text-2xl sm:text-3xl font-bold bg-gradient-to-r from-primary to-primary/60 bg-clip-text text-transparent">
            {T("home.supervisors", { defaultValue: "Supervisors" })}
          </h2>
          <p className="text-sm text-muted-foreground mt-1">
            {T("home.contacts_description", {
              defaultValue: "Contact information for your academic supervisors",
            })}
          </p>
        </div>
        <BackButton to="/" showLabelOnMobile className="w-full sm:w-auto" />
      </div>
      {isLoading ? (
        <Card className="p-8">
          <div className="text-center text-muted-foreground">
            {T("common.loading", { defaultValue: "Loading..." })}
          </div>
        </Card>
      ) : resolvedList.length === 0 ? (
        <Card className="p-8">
          <div className="text-center space-y-2">
            <Users className="w-12 h-12 text-muted-foreground mx-auto" />
            <div className="text-sm text-muted-foreground">
              {T("home.no_contacts", {
                defaultValue: "No contacts available.",
              })}
              {error && (
                <div className="mt-2 text-xs text-red-500">
                  {T("common.error", { defaultValue: "Error" })}:{" "}
                  {String((error as any)?.message || "Unable to load contacts")}
                </div>
              )}
            </div>
          </div>
        </Card>
      ) : (
        <div className="grid gap-4">
          {resolvedList.map((c, idx) => (
            <Card
              key={idx}
              className="bg-gradient-to-br from-card to-card/50 border-2 hover:border-primary/30 transition-all duration-300 hover:shadow-lg"
            >
              <CardContent className="p-6 space-y-3">
                <div className="flex items-start gap-3">
                  <div className="rounded-full bg-primary/10 p-3 text-primary">
                    <Users className="h-6 w-6" />
                  </div>
                  <div className="flex-1">
                    <div className="text-lg font-bold">
                      {textOf(c.name as Localized, lang)}
                    </div>
                    {((c as any).role || (c as any).title) && (
                      <div className="text-sm text-muted-foreground font-medium">
                        {textOf(
                          ((c as any).role || (c as any).title) as Localized,
                          lang
                        )}
                      </div>
                    )}
                  </div>
                </div>
                <div className="flex flex-col sm:flex-row gap-3 items-start text-sm pt-2">
                  {c.email && (
                    <div className="flex items-center gap-2 bg-muted/30 rounded-lg p-3 flex-1">
                      <Mail className="h-5 w-5 text-primary flex-shrink-0" />
                      <a
                        className="underline text-sm flex-1 hover:text-primary transition-colors"
                        href={`mailto:${c.email}`}
                      >
                        {c.email}
                      </a>
                      <button
                        className="p-2 hover:bg-primary/10 rounded-md transition-colors text-muted-foreground hover:text-primary"
                        onClick={() => {
                          copy(String(c.email));
                          noteCopied(`${idx}-email`);
                        }}
                        title={T("common.copy", { defaultValue: "Copy" })}
                      >
                        <Copy className="h-4 w-4" />
                      </button>
                      {copiedKey === `${idx}-email` && (
                        <span
                          className="text-xs text-emerald-600 font-semibold inline-flex items-center gap-1 px-2 py-1 bg-emerald-50 dark:bg-emerald-900/20 rounded-md animate-in fade-in duration-200"
                          aria-live="polite"
                        >
                          <Check className="h-3 w-3" />{" "}
                          {T("common.copied", { defaultValue: "Copied" })}
                        </span>
                      )}
                    </div>
                  )}
                  {c.phone && (
                    <div className="flex items-center gap-2 bg-muted/30 rounded-lg p-3 flex-1">
                      <Phone className="h-5 w-5 text-primary flex-shrink-0" />
                      <a
                        className="underline text-sm flex-1 hover:text-primary transition-colors"
                        href={`tel:${c.phone}`}
                      >
                        {c.phone}
                      </a>
                      <button
                        className="p-2 hover:bg-primary/10 rounded-md transition-colors text-muted-foreground hover:text-primary"
                        onClick={() => {
                          copy(String(c.phone));
                          noteCopied(`${idx}-phone`);
                        }}
                        title={T("common.copy", { defaultValue: "Copy" })}
                      >
                        <Copy className="h-4 w-4" />
                      </button>
                      {copiedKey === `${idx}-phone` && (
                        <span
                          className="text-xs text-emerald-600 font-semibold inline-flex items-center gap-1 px-2 py-1 bg-emerald-50 dark:bg-emerald-900/20 rounded-md animate-in fade-in duration-200"
                          aria-live="polite"
                        >
                          <Check className="h-3 w-3" />{" "}
                          {T("common.copied", { defaultValue: "Copied" })}
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
