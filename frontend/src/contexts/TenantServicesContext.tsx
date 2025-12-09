import React, { createContext, useContext, useEffect, useState } from 'react';
import { useAuth } from '@/contexts/AuthContext';
import { api } from '@/api/client';

// Optional services that can be enabled/disabled per tenant
export const OPTIONAL_SERVICES = ['chat', 'calendar', 'smtp', 'email'] as const;
export type OptionalService = typeof OPTIONAL_SERVICES[number];

// Core services that are always enabled
export const CORE_SERVICES = ['journey', 'contacts', 'notifications', 'uploads'] as const;
export type CoreService = typeof CORE_SERVICES[number];

interface TenantServicesContextValue {
  enabledServices: string[];
  isServiceEnabled: (service: OptionalService) => boolean;
  isLoading: boolean;
}

const TenantServicesContext = createContext<TenantServicesContextValue>({
  enabledServices: [],
  isServiceEnabled: () => false,
  isLoading: true,
});

export function TenantServicesProvider({ children }: { children: React.ReactNode }) {
  const { user } = useAuth();
  const [enabledServices, setEnabledServices] = useState<string[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    async function fetchTenantServices() {
      if (!user) {
        setIsLoading(false);
        return;
      }

      try {
        // Fetch current tenant's services
        const response = await api<{ enabled_services?: string[] }>('/me/tenant');
        setEnabledServices(response.enabled_services || ['chat', 'calendar']); // Default all enabled
      } catch (error) {
        console.error('Failed to fetch tenant services:', error);
        // Default to all enabled on error
        setEnabledServices(['chat', 'calendar']);
      } finally {
        setIsLoading(false);
      }
    }

    fetchTenantServices();
  }, [user]);

  const isServiceEnabled = (service: OptionalService): boolean => {
    return enabledServices.includes(service);
  };

  return (
    <TenantServicesContext.Provider value={{ enabledServices, isServiceEnabled, isLoading }}>
      {children}
    </TenantServicesContext.Provider>
  );
}

export function useTenantServices() {
  return useContext(TenantServicesContext);
}

/**
 * Hook to check if a specific service is enabled
 */
export function useServiceEnabled(service: OptionalService): boolean {
  const { isServiceEnabled, isLoading } = useTenantServices();
  // While loading, return true to prevent flash of hidden content
  if (isLoading) return true;
  return isServiceEnabled(service);
}

/**
 * Higher-order component to protect routes based on enabled services
 */
export function withServiceRequired<P extends object>(
  Component: React.ComponentType<P>,
  service: OptionalService
) {
  return function ServiceProtectedComponent(props: P) {
    const { isServiceEnabled, isLoading } = useTenantServices();
    
    if (isLoading) {
      return <div className="p-4 text-sm text-muted-foreground">Loading...</div>;
    }
    
    if (!isServiceEnabled(service)) {
      return (
        <div className="p-8 text-center">
          <h2 className="text-xl font-semibold mb-2">Service Not Available</h2>
          <p className="text-muted-foreground">
            This feature is not enabled for your institution.
          </p>
        </div>
      );
    }
    
    return <Component {...props} />;
  };
}
