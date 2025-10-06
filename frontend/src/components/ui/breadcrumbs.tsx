import * as React from "react";
import { Link, useMatches } from "react-router-dom";
import { useQuery } from "@tanstack/react-query";
import { api } from "../../api/client";
import { useTranslation } from "react-i18next";

const paths = ["/", "/login", "/journey", "/advisor/inbox", "/admin/users"] as const;

export function Breadcrumbs() {
  const { t: T } = useTranslation("common");
  const matches = useMatches() as Array<{ pathname: string }>;
  const { data: me } = useQuery({
    queryKey: ["me"],
    queryFn: () => api("/me"),
  });
  const role = me?.role;
  const items = matches
    .filter((m) => m.pathname !== "/")
    .map((m) => {
      const p = m.pathname;
      let label = p;
      if (p === "/login") label = T("nav.login");
      if (p === "/journey") label = T("breadcrumbs.journey");
      if (p === "/advisor/inbox") label = T("breadcrumbs.advisor_inbox");
      if (p === "/admin/users") label = T("breadcrumbs.admin_users");
      if (p.startsWith("/documents/")) label = T("breadcrumbs.document");
      if (p.startsWith("/admin") && role !== "admin" && role !== "superadmin")
        return null;
      if (
        p.startsWith("/advisor") &&
        !["advisor", "chair", "admin", "superadmin"].includes(role)
      )
        return null;
      return { path: p, label };
    })
    .filter(Boolean) as { path: string; label: string }[];
  if (items.length === 0) return null;
  return (
    <nav className="text-sm text-gray-600 py-1">
      {items.map((it, idx) => (
        <span key={it.path}>
          <Link to={it.path} className="underline">
            {it.label}
          </Link>
          {idx < items.length - 1 ? " / " : ""}
        </span>
      ))}
    </nav>
  );
}
