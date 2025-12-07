import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { useTranslation } from 'react-i18next';
import {
  Building2,
  Plus,
  Pencil,
  Trash2,
  Users,
  CheckCircle,
  XCircle,
  MessageCircle,
  Calendar,
} from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { Switch } from '@/components/ui/switch';
import { tenantsApi, Tenant, CreateTenantRequest } from '../api';
import { useState } from 'react';

const TENANT_TYPES = [
  { value: 'university', label: 'University', icon: '🎓' },
  { value: 'college', label: 'College', icon: '📚' },
  { value: 'vocational', label: 'Vocational School', icon: '🔧' },
  { value: 'school', label: 'School', icon: '🏫' },
];

function TenantForm({
  tenant,
  onSuccess,
  onCancel,
}: {
  tenant?: Tenant;
  onSuccess: () => void;
  onCancel: () => void;
}) {
  const { t } = useTranslation('common');
  const queryClient = useQueryClient();
  const [formData, setFormData] = useState<CreateTenantRequest>({
    slug: tenant?.slug || '',
    name: tenant?.name || '',
    tenant_type: tenant?.tenant_type || 'university',
    domain: tenant?.domain || '',
    app_name: tenant?.app_name || '',
    primary_color: tenant?.primary_color || '#3b82f6',
    secondary_color: tenant?.secondary_color || '#1e40af',
  });
  const [isActive, setIsActive] = useState(tenant?.is_active ?? true);

  const createMutation = useMutation({
    mutationFn: (data: CreateTenantRequest) => tenantsApi.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['superadmin', 'tenants'] });
      onSuccess();
    },
  });

  const updateMutation = useMutation({
    mutationFn: (data: { id: string; updates: any }) =>
      tenantsApi.update(data.id, data.updates),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['superadmin', 'tenants'] });
      onSuccess();
    },
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (tenant) {
      updateMutation.mutate({
        id: tenant.id,
        updates: { ...formData, is_active: isActive },
      });
    } else {
      createMutation.mutate(formData);
    }
  };

  const isLoading = createMutation.isPending || updateMutation.isPending;

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <div className="grid grid-cols-2 gap-4">
        <div className="space-y-2">
          <Label htmlFor="name">{t('superadmin.tenants.name', 'Name')} *</Label>
          <Input
            id="name"
            value={formData.name}
            onChange={(e) => setFormData({ ...formData, name: e.target.value })}
            placeholder="Kazakh National Medical University"
            required
          />
        </div>
        <div className="space-y-2">
          <Label htmlFor="slug">{t('superadmin.tenants.slug', 'Slug')} *</Label>
          <Input
            id="slug"
            value={formData.slug}
            onChange={(e) =>
              setFormData({ ...formData, slug: e.target.value.toLowerCase() })
            }
            placeholder="kaznmu"
            required
            pattern="[a-z0-9-]+"
            title="Lowercase letters, numbers, and hyphens only"
          />
        </div>
      </div>

      <div className="grid grid-cols-2 gap-4">
        <div className="space-y-2">
          <Label htmlFor="type">{t('superadmin.tenants.type', 'Type')}</Label>
          <Select
            value={formData.tenant_type}
            onValueChange={(value) =>
              setFormData({ ...formData, tenant_type: value })
            }
          >
            <SelectTrigger>
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              {TENANT_TYPES.map((type) => (
                <SelectItem key={type.value} value={type.value}>
                  {type.icon} {type.label}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>
        <div className="space-y-2">
          <Label htmlFor="domain">{t('superadmin.tenants.domain', 'Custom Domain')}</Label>
          <Input
            id="domain"
            value={formData.domain}
            onChange={(e) => setFormData({ ...formData, domain: e.target.value })}
            placeholder="portal.kaznmu.kz"
          />
        </div>
      </div>

      <div className="space-y-2">
        <Label htmlFor="app_name">{t('superadmin.tenants.app_name', 'Custom App Name')}</Label>
        <Input
          id="app_name"
          value={formData.app_name}
          onChange={(e) => setFormData({ ...formData, app_name: e.target.value })}
          placeholder="KazNMU PhD Portal"
        />
      </div>

      <div className="grid grid-cols-2 gap-4">
        <div className="space-y-2">
          <Label htmlFor="primary_color">{t('superadmin.tenants.primary_color', 'Primary Color')}</Label>
          <div className="flex gap-2">
            <Input
              type="color"
              id="primary_color"
              value={formData.primary_color}
              onChange={(e) =>
                setFormData({ ...formData, primary_color: e.target.value })
              }
              className="w-12 h-10 p-1"
            />
            <Input
              value={formData.primary_color}
              onChange={(e) =>
                setFormData({ ...formData, primary_color: e.target.value })
              }
              className="flex-1"
            />
          </div>
        </div>
        <div className="space-y-2">
          <Label htmlFor="secondary_color">{t('superadmin.tenants.secondary_color', 'Secondary Color')}</Label>
          <div className="flex gap-2">
            <Input
              type="color"
              id="secondary_color"
              value={formData.secondary_color}
              onChange={(e) =>
                setFormData({ ...formData, secondary_color: e.target.value })
              }
              className="w-12 h-10 p-1"
            />
            <Input
              value={formData.secondary_color}
              onChange={(e) =>
                setFormData({ ...formData, secondary_color: e.target.value })
              }
              className="flex-1"
            />
          </div>
        </div>
      </div>

      {tenant && (
        <div className="flex items-center gap-2">
          <Switch
            id="is_active"
            checked={isActive}
            onCheckedChange={setIsActive}
          />
          <Label htmlFor="is_active">{t('superadmin.tenants.active', 'Active')}</Label>
        </div>
      )}

      <div className="flex justify-end gap-2 pt-4">
        <Button type="button" variant="outline" onClick={onCancel}>
          {t('common.cancel', 'Cancel')}
        </Button>
        <Button type="submit" disabled={isLoading}>
          {isLoading ? '...' : tenant ? t('common.save', 'Save') : t('common.create', 'Create')}
        </Button>
      </div>
    </form>
  );
}

export function TenantsPage() {
  const { t } = useTranslation('common');
  const queryClient = useQueryClient();
  const [dialogOpen, setDialogOpen] = useState(false);
  const [editingTenant, setEditingTenant] = useState<Tenant | undefined>();

  const { data: tenants, isLoading } = useQuery({
    queryKey: ['superadmin', 'tenants'],
    queryFn: tenantsApi.list,
  });

  const deleteMutation = useMutation({
    mutationFn: tenantsApi.delete,
    onSuccess: () =>
      queryClient.invalidateQueries({ queryKey: ['superadmin', 'tenants'] }),
  });

  const updateServicesMutation = useMutation({
    mutationFn: ({ id, services }: { id: string; services: string[] }) =>
      tenantsApi.updateServices(id, services),
    onSuccess: () =>
      queryClient.invalidateQueries({ queryKey: ['superadmin', 'tenants'] }),
  });

  const toggleService = (tenant: Tenant, service: 'chat' | 'calendar') => {
    const current = tenant.enabled_services || [];
    const newServices = current.includes(service)
      ? current.filter((s) => s !== service)
      : [...current, service];
    updateServicesMutation.mutate({ id: tenant.id, services: newServices });
  };

  const openCreate = () => {
    setEditingTenant(undefined);
    setDialogOpen(true);
  };

  const openEdit = (tenant: Tenant) => {
    setEditingTenant(tenant);
    setDialogOpen(true);
  };

  const closeDialog = () => {
    setDialogOpen(false);
    setEditingTenant(undefined);
  };

  const getTenantTypeInfo = (type: string) =>
    TENANT_TYPES.find((t) => t.value === type) || TENANT_TYPES[0];

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold flex items-center gap-2">
            <Building2 className="h-6 w-6 text-violet-500" />
            {t('superadmin.tenants.title', 'Institutions')}
          </h1>
          <p className="text-muted-foreground">
            {t('superadmin.tenants.description', 'Manage universities, colleges, and schools')}
          </p>
        </div>
        <Dialog open={dialogOpen} onOpenChange={setDialogOpen}>
          <DialogTrigger asChild>
            <Button onClick={openCreate}>
              <Plus className="h-4 w-4 mr-2" />
              {t('superadmin.tenants.add', 'Add Institution')}
            </Button>
          </DialogTrigger>
          <DialogContent className="max-w-2xl">
            <DialogHeader>
              <DialogTitle>
                {editingTenant
                  ? t('superadmin.tenants.edit', 'Edit Institution')
                  : t('superadmin.tenants.add', 'Add Institution')}
              </DialogTitle>
            </DialogHeader>
            <TenantForm
              tenant={editingTenant}
              onSuccess={closeDialog}
              onCancel={closeDialog}
            />
          </DialogContent>
        </Dialog>
      </div>

      {isLoading ? (
        <div className="text-center py-8 text-muted-foreground">Loading...</div>
      ) : (
        <div className="border rounded-lg">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>{t('superadmin.tenants.name', 'Name')}</TableHead>
                <TableHead>{t('superadmin.tenants.slug', 'Slug')}</TableHead>
                <TableHead>{t('superadmin.tenants.type', 'Type')}</TableHead>
                <TableHead>{t('superadmin.tenants.services', 'Services')}</TableHead>
                <TableHead>{t('superadmin.tenants.users', 'Users')}</TableHead>
                <TableHead>{t('superadmin.tenants.status', 'Status')}</TableHead>
                <TableHead className="w-24"></TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {tenants?.map((tenant) => {
                const typeInfo = getTenantTypeInfo(tenant.tenant_type);
                return (
                  <TableRow key={tenant.id}>
                    <TableCell className="font-medium">
                      <div className="flex items-center gap-2">
                        {tenant.logo_url && (
                          <img
                            src={tenant.logo_url}
                            alt=""
                            className="h-6 w-6 rounded"
                          />
                        )}
                        {tenant.name}
                      </div>
                    </TableCell>
                    <TableCell>
                      <code className="text-xs bg-muted px-2 py-1 rounded">
                        {tenant.slug}
                      </code>
                    </TableCell>
                    <TableCell>
                      <span className="flex items-center gap-1">
                        {typeInfo.icon} {typeInfo.label}
                      </span>
                    </TableCell>
                    <TableCell>
                      <div className="flex items-center gap-2">
                        <button
                          className={`p-1.5 rounded-md transition-colors ${
                            tenant.enabled_services?.includes('chat')
                              ? 'bg-blue-100 text-blue-600 dark:bg-blue-900 dark:text-blue-300'
                              : 'bg-muted text-muted-foreground'
                          }`}
                          onClick={() => toggleService(tenant, 'chat')}
                          title={`Chat: ${tenant.enabled_services?.includes('chat') ? 'Enabled' : 'Disabled'}`}
                        >
                          <MessageCircle className="h-4 w-4" />
                        </button>
                        <button
                          className={`p-1.5 rounded-md transition-colors ${
                            tenant.enabled_services?.includes('calendar')
                              ? 'bg-green-100 text-green-600 dark:bg-green-900 dark:text-green-300'
                              : 'bg-muted text-muted-foreground'
                          }`}
                          onClick={() => toggleService(tenant, 'calendar')}
                          title={`Calendar: ${tenant.enabled_services?.includes('calendar') ? 'Enabled' : 'Disabled'}`}
                        >
                          <Calendar className="h-4 w-4" />
                        </button>
                      </div>
                    </TableCell>
                    <TableCell>
                      <div className="flex items-center gap-2 text-sm">
                        <Users className="h-4 w-4 text-muted-foreground" />
                        {tenant.user_count}
                        <span className="text-muted-foreground">
                          ({tenant.admin_count} admins)
                        </span>
                      </div>
                    </TableCell>
                    <TableCell>
                      {tenant.is_active ? (
                        <Badge variant="outline" className="text-green-600 border-green-600">
                          <CheckCircle className="h-3 w-3 mr-1" />
                          Active
                        </Badge>
                      ) : (
                        <Badge variant="outline" className="text-red-600 border-red-600">
                          <XCircle className="h-3 w-3 mr-1" />
                          Inactive
                        </Badge>
                      )}
                    </TableCell>
                    <TableCell>
                      <div className="flex items-center gap-1">
                        <Button
                          variant="ghost"
                          size="icon"
                          onClick={() => openEdit(tenant)}
                        >
                          <Pencil className="h-4 w-4" />
                        </Button>
                        <Button
                          variant="ghost"
                          size="icon"
                          onClick={() => {
                            if (confirm(t('superadmin.tenants.confirm_delete', 'Deactivate this institution?'))) {
                              deleteMutation.mutate(tenant.id);
                            }
                          }}
                        >
                          <Trash2 className="h-4 w-4 text-red-500" />
                        </Button>
                      </div>
                    </TableCell>
                  </TableRow>
                );
              })}
              {(!tenants || tenants.length === 0) && (
                <TableRow>
                  <TableCell colSpan={7} className="text-center py-8 text-muted-foreground">
                    {t('superadmin.tenants.empty', 'No institutions yet')}
                  </TableCell>
                </TableRow>
              )}
            </TableBody>
          </Table>
        </div>
      )}
    </div>
  );
}

export default TenantsPage;
