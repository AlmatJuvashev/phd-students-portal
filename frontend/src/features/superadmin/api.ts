import { API_URL } from '@/api/client';
import { getTenantHeaders } from '@/lib/tenant';

// Types for Superadmin API
export interface Tenant {
  id: string;
  slug: string;
  name: string;
  tenant_type: 'university' | 'college' | 'vocational' | 'school';
  domain?: string;
  logo_url?: string;
  app_name?: string;
  primary_color: string;
  secondary_color: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
  user_count: number;
  admin_count: number;
}

export interface CreateTenantRequest {
  slug: string;
  name: string;
  tenant_type?: string;
  domain?: string;
  app_name?: string;
  primary_color?: string;
  secondary_color?: string;
}

export interface UpdateTenantRequest {
  slug?: string;
  name?: string;
  tenant_type?: string;
  domain?: string;
  app_name?: string;
  primary_color?: string;
  secondary_color?: string;
  is_active?: boolean;
}

export interface Admin {
  id: string;
  username: string;
  email: string;
  first_name: string;
  last_name: string;
  role: string;
  is_active: boolean;
  is_superadmin: boolean;
  tenant_id?: string;
  tenant_name?: string;
  tenant_slug?: string;
  created_at: string;
  updated_at: string;
}

export interface CreateAdminRequest {
  username: string;
  email: string;
  password: string;
  first_name: string;
  last_name: string;
  role?: string;
  is_superadmin?: boolean;
  tenant_ids?: string[];
}

export interface ActivityLog {
  id: string;
  tenant_id?: string;
  tenant_name?: string;
  user_id?: string;
  username?: string;
  user_email?: string;
  action: string;
  entity_type?: string;
  entity_id?: string;
  description?: string;
  ip_address?: string;
  user_agent?: string;
  created_at: string;
}

export interface Setting {
  key: string;
  value: unknown;
  description?: string;
  category: string;
  updated_at: string;
  updated_by?: string;
}

// API helper
async function superadminFetch<T>(
  endpoint: string,
  options: RequestInit = {}
): Promise<T> {
  const token = localStorage.getItem('token');
  const response = await fetch(`${API_URL}${endpoint}`, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`,
      ...getTenantHeaders(),
      ...options.headers,
    },
  });

  if (!response.ok) {
    const error = await response.json().catch(() => ({ error: 'Request failed' }));
    throw new Error(error.error || 'Request failed');
  }

  return response.json();
}

// Tenants API
export const tenantsApi = {
  list: () => superadminFetch<Tenant[]>('/superadmin/tenants'),
  get: (id: string) => superadminFetch<Tenant>(`/superadmin/tenants/${id}`),
  create: (data: CreateTenantRequest) =>
    superadminFetch<Tenant>('/superadmin/tenants', {
      method: 'POST',
      body: JSON.stringify(data),
    }),
  update: (id: string, data: UpdateTenantRequest) =>
    superadminFetch<Tenant>(`/superadmin/tenants/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    }),
  delete: (id: string) =>
    superadminFetch<{ message: string }>(`/superadmin/tenants/${id}`, {
      method: 'DELETE',
    }),
};

// Admins API
export const adminsApi = {
  list: (tenantId?: string) =>
    superadminFetch<Admin[]>(
      `/superadmin/admins${tenantId ? `?tenant_id=${tenantId}` : ''}`
    ),
  get: (id: string) =>
    superadminFetch<{ admin: Admin; memberships: unknown[] }>(
      `/superadmin/admins/${id}`
    ),
  create: (data: CreateAdminRequest) =>
    superadminFetch<{ id: string; message: string }>('/superadmin/admins', {
      method: 'POST',
      body: JSON.stringify(data),
    }),
  update: (id: string, data: Partial<Admin> & { tenant_ids?: string[] }) =>
    superadminFetch<{ message: string }>(`/superadmin/admins/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    }),
  delete: (id: string) =>
    superadminFetch<{ message: string }>(`/superadmin/admins/${id}`, {
      method: 'DELETE',
    }),
  resetPassword: (id: string, password: string) =>
    superadminFetch<{ message: string }>(`/superadmin/admins/${id}/reset-password`, {
      method: 'POST',
      body: JSON.stringify({ password }),
    }),
};

// Logs API
export const logsApi = {
  list: (params: {
    page?: number;
    limit?: number;
    tenant_id?: string;
    user_id?: string;
    action?: string;
    entity_type?: string;
    start_date?: string;
    end_date?: string;
  }) => {
    const searchParams = new URLSearchParams();
    Object.entries(params).forEach(([key, value]) => {
      if (value !== undefined) searchParams.set(key, String(value));
    });
    return superadminFetch<{
      data: ActivityLog[];
      pagination: { page: number; limit: number; total: number; total_pages: number };
    }>(`/superadmin/logs?${searchParams.toString()}`);
  },
  getStats: () =>
    superadminFetch<{
      total_logs: number;
      logs_by_action: Record<string, number>;
      logs_by_tenant: { tenant_id: string; tenant_name: string; count: number }[];
      recent_activity: { date: string; count: number }[];
    }>('/superadmin/logs/stats'),
  getActions: () => superadminFetch<string[]>('/superadmin/logs/actions'),
  getEntityTypes: () => superadminFetch<string[]>('/superadmin/logs/entity-types'),
};

// Settings API
export const settingsApi = {
  list: (category?: string) =>
    superadminFetch<Setting[]>(
      `/superadmin/settings${category ? `?category=${category}` : ''}`
    ),
  get: (key: string) => superadminFetch<Setting>(`/superadmin/settings/${key}`),
  update: (key: string, data: { value: unknown; description?: string; category?: string }) =>
    superadminFetch<Setting>(`/superadmin/settings/${key}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    }),
  delete: (key: string) =>
    superadminFetch<{ message: string }>(`/superadmin/settings/${key}`, {
      method: 'DELETE',
    }),
  bulkUpdate: (settings: Record<string, unknown>) =>
    superadminFetch<{ updated: number }>('/superadmin/settings/bulk', {
      method: 'POST',
      body: JSON.stringify({ settings }),
    }),
  getCategories: () => superadminFetch<string[]>('/superadmin/settings/categories'),
};
