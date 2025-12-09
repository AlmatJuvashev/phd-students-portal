/**
 * Tenant Context for Multitenancy
 * 
 * Resolves tenant from subdomain or provides dev override via X-Tenant-Slug header.
 * In production: kaznmu.phd-portal.kz -> tenant = "kaznmu"
 * In development: localhost:3000 with header X-Tenant-Slug: kaznmu
 */

export interface Tenant {
  id: string;
  slug: string;
  name: string;
  domain?: string;
  logoUrl?: string;
  settings?: Record<string, unknown>;
  isActive: boolean;
}

/**
 * Get tenant slug from current hostname
 * - Production: extracts subdomain (e.g., "kaznmu" from "kaznmu.phd-portal.kz")
 * - Development: returns default tenant or reads from localStorage
 */
export function getTenantSlug(): string {
  const hostname = window.location.hostname;
  
  // Local development - use localStorage or default
  if (hostname === 'localhost' || hostname === '127.0.0.1') {
    return localStorage.getItem('tenant-slug') || 'kaznmu';
  }
  
  // Vercel deployments - map specific app names to tenants
  // Add new entries here for additional Vercel deployments
  if (hostname.endsWith('.vercel.app')) {
    const vercelTenantMap: Record<string, string> = {
      'demo-phd-students-portal': 'demo.university',
      'phd-students-portal': 'kaznmu',
      // Add more mappings as needed: 'other-app-name': 'other-tenant-slug'
    };
    
    // Extract app name from hostname (e.g., 'demo-phd-students-portal' from 'demo-phd-students-portal.vercel.app')
    const appName = hostname.replace('.vercel.app', '');
    return vercelTenantMap[appName] || 'kaznmu';
  }
  
  // Production/staging - extract subdomain
  const parts = hostname.split('.');
  if (parts.length >= 2) {
    const subdomain = parts[0];
    // Skip common non-tenant subdomains
    if (subdomain !== 'www' && subdomain !== 'api' && subdomain !== 'app') {
      return subdomain;
    }
  }
  
  // Fallback to default tenant
  return 'kaznmu';
}

/**
 * Set tenant slug for development (persists to localStorage)
 */
export function setDevTenantSlug(slug: string): void {
  localStorage.setItem('tenant-slug', slug);
}

/**
 * Get tenant header for API requests
 * Returns the X-Tenant-Slug header object for fetch calls
 */
export function getTenantHeaders(): Record<string, string> {
  return {
    'X-Tenant-Slug': getTenantSlug(),
  };
}

/**
 * Check if we're in dev mode (localhost)
 */
export function isDevMode(): boolean {
  const hostname = window.location.hostname;
  return hostname === 'localhost' || hostname === '127.0.0.1';
}
