import React, { Suspense } from "react";
import ReactDOM from "react-dom/client";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { RouterProvider } from "react-router-dom";
import "./index.css";
import { router } from "./routes";
import "./i18n";
import { AuthProvider } from '@/contexts/AuthContext'
import { TenantServicesProvider } from '@/contexts/TenantServicesContext';
import { registerSW } from "virtual:pwa-register";
import { DevNavigation } from "@/components/DevNavigation";

registerSW({ immediate: true });

const qc = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 5 * 60 * 1000, // 5m – keep results fresh to avoid reloading
      gcTime: 30 * 60 * 1000, // cache for 30m
      refetchOnWindowFocus: false,
      refetchOnReconnect: true,
      refetchOnMount: false,
      retry: 1,
    },
  },
});

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <QueryClientProvider client={qc}>
      <Suspense fallback={<div className="p-4 text-sm">Loading…</div>}>
        <AuthProvider>
          <TenantServicesProvider>
            <RouterProvider router={router} />
            {import.meta.env.DEV && (
              <DevNavigation onNavigate={(path: string) => router.navigate(path)} />
            )}
          </TenantServicesProvider>
        </AuthProvider>
      </Suspense>
    </QueryClientProvider>
  </React.StrictMode>
);

// Warm up common lazy chunks shortly after boot to reduce first-open latency
setTimeout(() => {
  // These imports create their own chunks via React.lazy; preloading them improves UX
}, 0);
