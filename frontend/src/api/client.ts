const getBaseUrl = () => {
  const defaultUrl = import.meta.env.VITE_API_URL || "http://localhost:8280/api";
  
  if (typeof window === "undefined") return defaultUrl;

  const hostname = window.location.hostname;
  
  // If we're on a subdomain of localhost (e.g. demo.localhost), 
  // we must talk to the API on that same host to maintain "same-site" context for cookies.
  if (hostname.includes(".localhost")) {
    return `http://${hostname}:8280/api`;
  }
  
  return defaultUrl;
};

export const API_URL = getBaseUrl();

import { getTenantHeaders } from '@/lib/tenant';

export async function api<T = any>(
  path: string,
  opts: RequestInit = {}
): Promise<T> {
  // Token is managed via HttpOnly cookies

  // Debug logging for auth issues
  if (path === "/notifications" || path.includes("/notifications")) {
    console.log("[API Debug] Notifications request:", {
      path,
      // Token is hidden in HttpOnly cookie
    });
  }

  // Debug logging for FormData uploads
  if (opts.body instanceof FormData) {
    console.log("[API Debug] FormData request:", {
      path,
      method: opts.method || "GET",
    });
  }

  const headers: Record<string, string> = {
    ...((opts.headers as Record<string, string>) || {}),
    // Token is now in HttpOnly cookie, so we don't send Authorization header
    ...getTenantHeaders(), // Add X-Tenant-Slug header
  };

  if (!headers["Content-Type"] && !(opts.body instanceof FormData)) {
    headers["Content-Type"] = "application/json";
  }

  // Ensure cookies are sent
  opts.credentials = 'include';
  
  console.log("[API Debug] Request details:", {
    url: `${API_URL}${path}`,
    method: opts.method || "GET",
    headers,
    isFormData: opts.body instanceof FormData,
  });
  
  const res = await fetch(`${API_URL}${path}`, { ...opts, headers });
  
  console.log("[API Debug] Response:", {
    status: res.status,
    statusText: res.statusText,
    contentType: res.headers.get("content-type"),
  });
  
  if (!res.ok) {
    const t = await res.text();
    // Enhanced error logging
    if (res.status === 401) {
      console.error("[API 401 Error]", {
        path,
        response: t,
      });
    }
    console.error("[API Error]", {
      path,
      status: res.status,
      response: t,
    });
    throw new Error(t || res.statusText);
  }
  const type = res.headers.get("content-type") || "";
  return (
    type.includes("application/json") ? res.json() : res.text()
  ) as Promise<T>;
}

// HTTP method helpers
api.get = (path: string, opts?: RequestInit) =>
  api(path, { ...opts, method: "GET" });
api.post = (path: string, data?: any, opts?: RequestInit) =>
  api(path, {
    ...opts,
    method: "POST",
    body: data ? JSON.stringify(data) : undefined,
  });
api.put = (path: string, data?: any, opts?: RequestInit) =>
  api(path, {
    ...opts,
    method: "PUT",
    body: data ? JSON.stringify(data) : undefined,
  });
api.patch = (path: string, data?: any, opts?: RequestInit) =>
  api(path, {
    ...opts,
    method: "PATCH",
    body: data ? JSON.stringify(data) : undefined,
  });
api.delete = (path: string, opts?: RequestInit) =>
  api(path, { ...opts, method: "DELETE" });
