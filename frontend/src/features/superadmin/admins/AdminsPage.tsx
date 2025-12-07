import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { useTranslation } from 'react-i18next';
import {
  Users,
  Plus,
  Pencil,
  Trash2,
  Key,
  Shield,
  CheckCircle,
  XCircle,
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
import { adminsApi, tenantsApi, Admin, CreateAdminRequest } from '../api';
import { useState } from 'react';

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

  const { data: admins, isLoading } = useQuery({
    queryKey: ['superadmin', 'admins'],
    queryFn: () => adminsApi.list(),
  });

  const deleteMutation = useMutation({
    mutationFn: adminsApi.delete,
    onSuccess: () =>
      queryClient.invalidateQueries({ queryKey: ['superadmin', 'admins'] }),
  });

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
        <div className="border rounded-lg">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>{t('superadmin.admins.name', 'Name')}</TableHead>
                <TableHead>{t('superadmin.admins.email', 'Email')}</TableHead>
                <TableHead>{t('superadmin.admins.institution', 'Institution')}</TableHead>
                <TableHead>{t('superadmin.admins.role', 'Role')}</TableHead>
                <TableHead>{t('superadmin.admins.status', 'Status')}</TableHead>
                <TableHead className="w-32"></TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {admins?.map((admin) => (
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
                          if (confirm(t('superadmin.admins.confirm_delete', 'Deactivate this admin?'))) {
                            deleteMutation.mutate(admin.id);
                          }
                        }}
                        title="Delete"
                      >
                        <Trash2 className="h-4 w-4 text-red-500" />
                      </Button>
                    </div>
                  </TableCell>
                </TableRow>
              ))}
              {(!admins || admins.length === 0) && (
                <TableRow>
                  <TableCell colSpan={6} className="text-center py-8 text-muted-foreground">
                    {t('superadmin.admins.empty', 'No administrators yet')}
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

export default AdminsPage;
