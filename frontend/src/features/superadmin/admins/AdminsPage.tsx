import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { useTranslation } from 'react-i18next';
import { useState, useMemo, useCallback } from 'react';
import {
  Users,
  Plus,
  Pencil,
  Trash2,
  Key,
  Shield,
  CheckCircle,
  XCircle,
  ArrowUpDown,
  ArrowUp,
  ArrowDown,
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
import { Switch } from '@/components/ui/switch';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { adminsApi, tenantsApi, Admin, CreateAdminRequest, Tenant } from '../api';
import { SuperadminTableToolbar, SuperadminPagination, ConfirmDialog, type FilterConfig } from '../components';

function AdminForm({
  admin,
  onSuccess,
  onCancel,
}: {
  admin?: Admin;
  onSuccess: () => void;
  onCancel: () => void;
}) {
  const { t } = useTranslation('common');
  const queryClient = useQueryClient();
  const [formData, setFormData] = useState<CreateAdminRequest>({
    username: '',
    email: admin?.email || '',
    password: '',
    first_name: admin?.first_name || '',
    last_name: admin?.last_name || '',
    role: admin?.role || 'admin',
    is_superadmin: admin?.is_superadmin || false,
    tenant_ids: [],
  });
  const [isActive, setIsActive] = useState(admin?.is_active ?? true);
  const [selectedTenants, setSelectedTenants] = useState<string[]>([]);

  const { data: tenants } = useQuery({
    queryKey: ['superadmin', 'tenants'],
    queryFn: tenantsApi.list,
  });

  const createMutation = useMutation({
    mutationFn: (data: CreateAdminRequest) => adminsApi.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['superadmin', 'admins'] });
      onSuccess();
    },
  });

  const updateMutation = useMutation({
    mutationFn: (data: { id: string; updates: any }) =>
      adminsApi.update(data.id, data.updates),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['superadmin', 'admins'] });
      onSuccess();
    },
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (admin) {
      updateMutation.mutate({
        id: admin.id,
        updates: {
          email: formData.email,
          first_name: formData.first_name,
          last_name: formData.last_name,
          role: formData.role,
          is_superadmin: formData.is_superadmin,
          is_active: isActive,
          tenant_ids: selectedTenants,
        },
      });
    } else {
      createMutation.mutate({
        ...formData,
        tenant_ids: selectedTenants,
      });
    }
  };

  const isLoading = createMutation.isPending || updateMutation.isPending;

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      {!admin && (
        <div className="grid grid-cols-2 gap-4">
          <div className="space-y-2">
            <Label htmlFor="username">{t('superadmin.admins.username', 'Username')} *</Label>
            <Input
              id="username"
              value={formData.username}
              onChange={(e) => setFormData({ ...formData, username: e.target.value })}
              required
            />
          </div>
          <div className="space-y-2">
            <Label htmlFor="password">{t('superadmin.admins.password', 'Password')} *</Label>
            <Input
              id="password"
              type="password"
              value={formData.password}
              onChange={(e) => setFormData({ ...formData, password: e.target.value })}
              required
              minLength={6}
            />
          </div>
        </div>
      )}

      <div className="grid grid-cols-2 gap-4">
        <div className="space-y-2">
          <Label htmlFor="first_name">{t('superadmin.admins.first_name', 'First Name')} *</Label>
          <Input
            id="first_name"
            value={formData.first_name}
            onChange={(e) => setFormData({ ...formData, first_name: e.target.value })}
            required
          />
        </div>
        <div className="space-y-2">
          <Label htmlFor="last_name">{t('superadmin.admins.last_name', 'Last Name')} *</Label>
          <Input
            id="last_name"
            value={formData.last_name}
            onChange={(e) => setFormData({ ...formData, last_name: e.target.value })}
            required
          />
        </div>
      </div>

      <div className="space-y-2">
        <Label htmlFor="email">{t('superadmin.admins.email', 'Email')} *</Label>
        <Input
          id="email"
          type="email"
          value={formData.email}
          onChange={(e) => setFormData({ ...formData, email: e.target.value })}
          required
        />
      </div>

      <div className="grid grid-cols-2 gap-4">
        <div className="space-y-2">
          <Label htmlFor="role">{t('superadmin.admins.role', 'Role')}</Label>
          <Select
            value={formData.role}
            onValueChange={(value) => setFormData({ ...formData, role: value })}
          >
            <SelectTrigger>
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="admin">Admin</SelectItem>
              <SelectItem value="superadmin">Superadmin</SelectItem>
            </SelectContent>
          </Select>
        </div>
        <div className="space-y-2 flex items-end">
          <div className="flex items-center gap-2 h-10">
            <Switch
              id="is_superadmin"
              checked={formData.is_superadmin}
              onCheckedChange={(checked) =>
                setFormData({ ...formData, is_superadmin: checked })
              }
            />
            <Label htmlFor="is_superadmin" className="flex items-center gap-1">
              <Shield className="h-4 w-4 text-violet-500" />
              {t('superadmin.admins.is_superadmin', 'Platform Superadmin')}
            </Label>
          </div>
        </div>
      </div>

      <div className="space-y-2">
        <Label>{t('superadmin.admins.tenants', 'Assign to Institutions')}</Label>
        <div className="border rounded-md p-3 max-h-40 overflow-y-auto space-y-2">
          {tenants?.map((tenant) => (
            <label key={tenant.id} className="flex items-center gap-2 cursor-pointer">
              <input
                type="checkbox"
                checked={selectedTenants.includes(tenant.id)}
                onChange={(e) => {
                  if (e.target.checked) {
                    setSelectedTenants([...selectedTenants, tenant.id]);
                  } else {
                    setSelectedTenants(selectedTenants.filter((id) => id !== tenant.id));
                  }
                }}
                className="rounded"
              />
              <span>{tenant.name}</span>
              <span className="text-xs text-muted-foreground">({tenant.slug})</span>
            </label>
          ))}
        </div>
      </div>

      {admin && (
        <div className="flex items-center gap-2">
          <Switch id="is_active" checked={isActive} onCheckedChange={setIsActive} />
          <Label htmlFor="is_active">{t('superadmin.admins.active', 'Active')}</Label>
        </div>
      )}

      <div className="flex justify-end gap-2 pt-4">
        <Button type="button" variant="outline" onClick={onCancel}>
          {t('common.cancel', 'Cancel')}
        </Button>
        <Button type="submit" disabled={isLoading}>
          {isLoading ? '...' : admin ? t('common.save', 'Save') : t('common.create', 'Create')}
        </Button>
      </div>
    </form>
  );
}

function ResetPasswordDialog({ adminId, onClose }: { adminId: string; onClose: () => void }) {
  const { t } = useTranslation('common');
  const [password, setPassword] = useState('');
  
  const mutation = useMutation({
    mutationFn: () => adminsApi.resetPassword(adminId, password),
    onSuccess: () => {
      alert(t('superadmin.admins.password_reset_success', 'Password reset successfully'));
      onClose();
    },
  });

  return (
    <div className="space-y-4">
      <div className="space-y-2">
        <Label htmlFor="new_password">{t('superadmin.admins.new_password', 'New Password')}</Label>
        <Input
          id="new_password"
          type="password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          minLength={6}
        />
      </div>
      <div className="flex justify-end gap-2">
        <Button variant="outline" onClick={onClose}>
          {t('common.cancel', 'Cancel')}
        </Button>
        <Button
          onClick={() => mutation.mutate()}
          disabled={password.length < 6 || mutation.isPending}
        >
          {mutation.isPending ? '...' : t('superadmin.admins.reset_password', 'Reset Password')}
        </Button>
      </div>
    </div>
  );
}

export function AdminsPage() {
  const { t } = useTranslation('common');
  const queryClient = useQueryClient();
  const [dialogOpen, setDialogOpen] = useState(false);
  const [editingAdmin, setEditingAdmin] = useState<Admin | undefined>();
  const [resetPasswordAdminId, setResetPasswordAdminId] = useState<string | null>(null);
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const [adminToDelete, setAdminToDelete] = useState<Admin | null>(null);

  // Table state
  const [searchQuery, setSearchQuery] = useState('');
  const [filters, setFilters] = useState<Record<string, string>>({
    status: 'all',
    tenant: 'all',
    role: 'all',
  });
  const [sortColumn, setSortColumn] = useState<string>('name');
  const [sortDirection, setSortDirection] = useState<'asc' | 'desc'>('asc');
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(25);

  const { data: admins, isLoading } = useQuery({
    queryKey: ['superadmin', 'admins'],
    queryFn: () => adminsApi.list(),
  });

  const { data: tenants } = useQuery({
    queryKey: ['superadmin', 'tenants'],
    queryFn: tenantsApi.list,
  });

  // Filter configuration
  const filterConfig: FilterConfig[] = useMemo(() => [
    {
      key: 'status',
      label: 'Status',
      options: [
        { value: 'all', label: 'All Status' },
        { value: 'active', label: 'Active' },
        { value: 'inactive', label: 'Inactive' },
      ],
    },
    {
      key: 'role',
      label: 'Role',
      options: [
        { value: 'all', label: 'All Roles' },
        { value: 'admin', label: 'Admin' },
        { value: 'superadmin', label: 'Superadmin' },
      ],
    },
    {
      key: 'tenant',
      label: 'Institution',
      options: [
        { value: 'all', label: 'All Institutions' },
        ...(tenants?.map((t: Tenant) => ({ value: t.id, label: t.name })) || []),
      ],
    },
  ], [tenants]);

  // Filter, sort, and paginate data
  const processedData = useMemo(() => {
    if (!admins) return { items: [], total: 0, totalPages: 0 };

    let filtered = [...admins];

    // Search filter
    if (searchQuery) {
      const query = searchQuery.toLowerCase();
      filtered = filtered.filter(
        (a) =>
          a.first_name.toLowerCase().includes(query) ||
          a.last_name.toLowerCase().includes(query) ||
          a.email.toLowerCase().includes(query) ||
          a.username.toLowerCase().includes(query)
      );
    }

    // Status filter
    if (filters.status !== 'all') {
      filtered = filtered.filter((a) =>
        filters.status === 'active' ? a.is_active : !a.is_active
      );
    }

    // Role filter
    if (filters.role !== 'all') {
      filtered = filtered.filter((a) => a.role === filters.role);
    }

    // Tenant filter
    if (filters.tenant !== 'all') {
      filtered = filtered.filter((a) => a.tenant_id === filters.tenant);
    }

    // Sort
    filtered.sort((a, b) => {
      let aVal: any, bVal: any;
      switch (sortColumn) {
        case 'name':
          aVal = `${a.first_name} ${a.last_name}`.toLowerCase();
          bVal = `${b.first_name} ${b.last_name}`.toLowerCase();
          break;
        case 'email':
          aVal = a.email.toLowerCase();
          bVal = b.email.toLowerCase();
          break;
        case 'tenant':
          aVal = a.tenant_name || '';
          bVal = b.tenant_name || '';
          break;
        case 'role':
          aVal = a.role;
          bVal = b.role;
          break;
        case 'status':
          aVal = a.is_active ? 1 : 0;
          bVal = b.is_active ? 1 : 0;
          break;
        default:
          aVal = a.email;
          bVal = b.email;
      }
      if (aVal < bVal) return sortDirection === 'asc' ? -1 : 1;
      if (aVal > bVal) return sortDirection === 'asc' ? 1 : -1;
      return 0;
    });

    const total = filtered.length;
    const totalPages = Math.ceil(total / pageSize);

    // Paginate
    const start = (currentPage - 1) * pageSize;
    const items = filtered.slice(start, start + pageSize);

    return { items, total, totalPages };
  }, [admins, searchQuery, filters, sortColumn, sortDirection, currentPage, pageSize]);

  const handleSort = useCallback((column: string) => {
    if (sortColumn === column) {
      setSortDirection((d) => (d === 'asc' ? 'desc' : 'asc'));
    } else {
      setSortColumn(column);
      setSortDirection('asc');
    }
  }, [sortColumn]);

  const handleFilterChange = useCallback((key: string, value: string) => {
    setFilters((prev) => ({ ...prev, [key]: value }));
    setCurrentPage(1);
  }, []);

  const handleClearFilters = useCallback(() => {
    setSearchQuery('');
    setFilters({ status: 'all', tenant: 'all', role: 'all' });
    setCurrentPage(1);
  }, []);

  const handleSearchChange = useCallback((value: string) => {
    setSearchQuery(value);
    setCurrentPage(1);
  }, []);

  const handlePageSizeChange = useCallback((size: number) => {
    setPageSize(size);
    setCurrentPage(1);
  }, []);

  const deleteMutation = useMutation({
    mutationFn: adminsApi.delete,
    onSuccess: () =>
      queryClient.invalidateQueries({ queryKey: ['superadmin', 'admins'] }),
  });

  const SortableHeader = ({ column, children }: { column: string; children: React.ReactNode }) => (
    <TableHead
      className="cursor-pointer select-none hover:bg-muted/50"
      onClick={() => handleSort(column)}
    >
      <div className="flex items-center gap-1">
        {children}
        {sortColumn === column ? (
          sortDirection === 'asc' ? (
            <ArrowUp className="h-3 w-3" />
          ) : (
            <ArrowDown className="h-3 w-3" />
          )
        ) : (
          <ArrowUpDown className="h-3 w-3 opacity-30" />
        )}
      </div>
    </TableHead>
  );

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold flex items-center gap-2">
            <Users className="h-6 w-6 text-violet-500" />
            {t('superadmin.admins.title', 'Administrators')}
          </h1>
          <p className="text-muted-foreground">
            {t('superadmin.admins.description', 'Manage administrators across all institutions')}
          </p>
        </div>
        <Dialog open={dialogOpen} onOpenChange={setDialogOpen}>
          <DialogTrigger asChild>
            <Button onClick={() => { setEditingAdmin(undefined); setDialogOpen(true); }}>
              <Plus className="h-4 w-4 mr-2" />
              {t('superadmin.admins.add', 'Add Admin')}
            </Button>
          </DialogTrigger>
          <DialogContent className="max-w-2xl">
            <DialogHeader>
              <DialogTitle>
                {editingAdmin
                  ? t('superadmin.admins.edit', 'Edit Admin')
                  : t('superadmin.admins.add', 'Add Admin')}
              </DialogTitle>
            </DialogHeader>
            <AdminForm
              admin={editingAdmin}
              onSuccess={() => { setDialogOpen(false); setEditingAdmin(undefined); }}
              onCancel={() => { setDialogOpen(false); setEditingAdmin(undefined); }}
            />
          </DialogContent>
        </Dialog>
      </div>

      {/* Reset Password Dialog */}
      <Dialog open={!!resetPasswordAdminId} onOpenChange={() => setResetPasswordAdminId(null)}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>{t('superadmin.admins.reset_password', 'Reset Password')}</DialogTitle>
          </DialogHeader>
          {resetPasswordAdminId && (
            <ResetPasswordDialog
              adminId={resetPasswordAdminId}
              onClose={() => setResetPasswordAdminId(null)}
            />
          )}
        </DialogContent>
      </Dialog>

      {isLoading ? (
        <div className="text-center py-8 text-muted-foreground">Loading...</div>
      ) : (
        <>
          <SuperadminTableToolbar
            searchPlaceholder={t('superadmin.admins.search', 'Search administrators...')}
            searchValue={searchQuery}
            onSearchChange={handleSearchChange}
            filters={filterConfig}
            filterValues={filters}
            onFilterChange={handleFilterChange}
            onClearFilters={handleClearFilters}
          />
          <div className="border rounded-lg">
            <Table>
              <TableHeader>
                <TableRow>
                  <SortableHeader column="name">{t('superadmin.admins.name', 'Name')}</SortableHeader>
                  <SortableHeader column="email">{t('superadmin.admins.email', 'Email')}</SortableHeader>
                  <SortableHeader column="tenant">{t('superadmin.admins.institution', 'Institution')}</SortableHeader>
                  <SortableHeader column="role">{t('superadmin.admins.role', 'Role')}</SortableHeader>
                  <SortableHeader column="status">{t('superadmin.admins.status', 'Status')}</SortableHeader>
                  <TableHead className="w-32"></TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {processedData.items.map((admin) => (
                  <TableRow key={admin.id + (admin.tenant_id || '')}>
                    <TableCell className="font-medium">
                      {admin.first_name} {admin.last_name}
                      {admin.is_superadmin && (
                        <Shield className="inline h-4 w-4 ml-1 text-violet-500" />
                      )}
                    </TableCell>
                    <TableCell>{admin.email}</TableCell>
                    <TableCell>
                      {admin.tenant_name ? (
                        <span>{admin.tenant_name}</span>
                      ) : (
                        <span className="text-muted-foreground">—</span>
                      )}
                    </TableCell>
                    <TableCell>
                      <Badge variant="outline">{admin.role}</Badge>
                    </TableCell>
                    <TableCell>
                      {admin.is_active ? (
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
                          onClick={() => { setEditingAdmin(admin); setDialogOpen(true); }}
                          title="Edit"
                        >
                          <Pencil className="h-4 w-4" />
                        </Button>
                        <Button
                          variant="ghost"
                          size="icon"
                          onClick={() => setResetPasswordAdminId(admin.id)}
                          title="Reset Password"
                        >
                          <Key className="h-4 w-4" />
                        </Button>
                        <Button
                          variant="ghost"
                          size="icon"
                          onClick={() => {
                            setAdminToDelete(admin);
                            setDeleteDialogOpen(true);
                          }}
                          title="Delete"
                        >
                          <Trash2 className="h-4 w-4 text-red-500" />
                        </Button>
                      </div>
                    </TableCell>
                  </TableRow>
                ))}
                {processedData.items.length === 0 && (
                  <TableRow>
                    <TableCell colSpan={6} className="text-center py-8 text-muted-foreground">
                      {searchQuery || Object.values(filters).some(v => v !== 'all')
                        ? t('superadmin.admins.no_results', 'No matching administrators found')
                        : t('superadmin.admins.empty', 'No administrators yet')}
                    </TableCell>
                  </TableRow>
                )}
              </TableBody>
            </Table>
          </div>
          {processedData.total > 0 && (
            <SuperadminPagination
              currentPage={currentPage}
              totalPages={processedData.totalPages}
              totalItems={processedData.total}
              pageSize={pageSize}
              onPageChange={setCurrentPage}
              onPageSizeChange={handlePageSizeChange}
            />
          )}
        </>
      )}

      {/* Delete Confirmation Dialog */}
      <ConfirmDialog
        open={deleteDialogOpen}
        onOpenChange={setDeleteDialogOpen}
        title={t('superadmin.admins.delete_title', 'Deactivate Administrator')}
        description={t('superadmin.admins.delete_description', `Are you sure you want to deactivate "${adminToDelete?.email}"? They will no longer be able to access the admin panel.`)}
        confirmLabel={t('superadmin.admins.deactivate', 'Deactivate')}
        onConfirm={() => {
          if (adminToDelete) {
            deleteMutation.mutate(adminToDelete.id);
          }
        }}
        loading={deleteMutation.isPending}
      />
    </div>
  );
}

export default AdminsPage;
