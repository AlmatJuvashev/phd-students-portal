import { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { Check, Building2, ChevronDown } from 'lucide-react';
import { Button } from '@/components/ui/button';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { api } from '@/api/client';
import { setDevTenantSlug, getTenantSlug } from '@/lib/tenant';

interface TenantMembership {
  tenant_id: string;
  tenant_name: string;
  tenant_slug: string;
  role: string;
  is_primary: boolean;
}

interface TenantSwitcherProps {
  className?: string;
}

/**
 * TenantSwitcher component - allows users with access to multiple tenants
 * to switch between them. Only renders if user has 2+ tenant memberships.
 */
export function TenantSwitcher({ className }: TenantSwitcherProps) {
  const { t } = useTranslation('common');
  const [memberships, setMemberships] = useState<TenantMembership[]>([]);
  const [currentSlug, setCurrentSlug] = useState(getTenantSlug());
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    async function fetchMemberships() {
      try {
        const response = await api<{ memberships: TenantMembership[] }>('/me/tenants');
        setMemberships(response.memberships || []);
      } catch (error) {
        console.error('Failed to fetch tenant memberships:', error);
        setMemberships([]);
      } finally {
        setIsLoading(false);
      }
    }

    fetchMemberships();
  }, []);

  // Don't render if user has less than 2 tenants
  if (isLoading || memberships.length < 2) {
    return null;
  }

  const currentTenant = memberships.find((m) => m.tenant_slug === currentSlug);

  const handleSwitch = (membership: TenantMembership) => {
    if (membership.tenant_slug === currentSlug) return;
    
    // Update tenant in localStorage for dev mode
    setDevTenantSlug(membership.tenant_slug);
    setCurrentSlug(membership.tenant_slug);
    
    // Reload to apply new tenant context
    window.location.reload();
  };

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button variant="ghost" className={className}>
          <Building2 className="h-4 w-4 mr-2" />
          <span className="truncate max-w-[120px]">
            {currentTenant?.tenant_name || currentSlug}
          </span>
          <ChevronDown className="h-4 w-4 ml-1" />
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end" className="w-[220px]">
        {memberships.map((membership) => (
          <DropdownMenuItem
            key={membership.tenant_id}
            onClick={() => handleSwitch(membership)}
            className="flex items-center justify-between"
          >
            <div className="flex flex-col">
              <span className="font-medium">{membership.tenant_name}</span>
              <span className="text-xs text-muted-foreground">
                {membership.role} • {membership.tenant_slug}
              </span>
            </div>
            {membership.tenant_slug === currentSlug && (
              <Check className="h-4 w-4 text-green-500" />
            )}
          </DropdownMenuItem>
        ))}
      </DropdownMenuContent>
    </DropdownMenu>
  );
}

export default TenantSwitcher;
