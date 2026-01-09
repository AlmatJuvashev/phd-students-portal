import React, { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { 
  Search, Users, UserPlus, MoreHorizontal, Shield, 
  GraduationCap, Briefcase, Phone, Mail, Building2,
  Loader2, Check, X, Key, RefreshCw
} from 'lucide-react';
import { toast } from 'sonner';
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from '@/components/ui/dialog';
import { DropdownMenu, DropdownItem } from '@/components/ui/dropdown-menu';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { Label } from '@/components/ui/label';
import { cn } from '@/lib/utils';
import { getStaffList, createStaff, updateStaff, setStaffActive, resetStaffPassword, Staff, CreateStaffRequest } from '@/features/hr/api';

const ROLES = [
  { value: 'student', label: 'Student', icon: GraduationCap, color: 'bg-emerald-100 text-emerald-700' },
  { value: 'instructor', label: 'Teacher / Instructor', icon: Briefcase, color: 'bg-indigo-100 text-indigo-700' },
  { value: 'advisor', label: 'Advisor', icon: GraduationCap, color: 'bg-blue-100 text-blue-700' },
  { value: 'chair', label: 'Department Head / Chair', icon: Briefcase, color: 'bg-purple-100 text-purple-700' },
  { value: 'dean', label: 'Dean', icon: Shield, color: 'bg-rose-100 text-rose-700' },
  { value: 'admin', label: 'Administrator', icon: Shield, color: 'bg-amber-100 text-amber-700' },
];

export const HRPage: React.FC = () => {
  const { t } = useTranslation('common');
  const queryClient = useQueryClient();
  
  const [search, setSearch] = useState('');
  const [roleFilter, setRoleFilter] = useState('all');
  const [activeFilter, setActiveFilter] = useState('true');
  const [page, setPage] = useState(1);
  
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingStaff, setEditingStaff] = useState<Staff | null>(null);
  const [credentialsModal, setCredentialsModal] = useState<{ username: string; temp_password: string } | null>(null);
  
  // Form state
  const [formData, setFormData] = useState<CreateStaffRequest>({
    first_name: '',
    last_name: '',
    email: '',
    role: 'advisor',
    phone: '',
    department: '',
  });

  // Queries
  const { data: staffData, isLoading } = useQuery({
    queryKey: ['staff', page, roleFilter, activeFilter, search],
    queryFn: () => getStaffList({
      page,
      limit: 20,
      role: roleFilter !== 'all' ? roleFilter : undefined,
      active: activeFilter,
      search: search || undefined,
    }),
  });

  // Mutations
  const createMutation = useMutation({
    mutationFn: createStaff,
    onSuccess: (data) => {
      queryClient.invalidateQueries({ queryKey: ['staff'] });
      setIsModalOpen(false);
      setCredentialsModal(data);
      toast.success(t('hr.staff_created', 'Staff member created'));
    },
    onError: () => toast.error(t('hr.create_error', 'Failed to create staff')),
  });

  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: CreateStaffRequest }) => updateStaff(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['staff'] });
      setIsModalOpen(false);
      setEditingStaff(null);
      toast.success(t('hr.staff_updated', 'Staff member updated'));
    },
    onError: () => toast.error(t('hr.update_error', 'Failed to update staff')),
  });

  const toggleActiveMutation = useMutation({
    mutationFn: ({ id, active }: { id: string; active: boolean }) => setStaffActive(id, active),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['staff'] });
      toast.success(t('hr.status_updated', 'Status updated'));
    },
  });

  const resetPasswordMutation = useMutation({
    mutationFn: resetStaffPassword,
    onSuccess: (data) => setCredentialsModal(data),
    onError: () => toast.error(t('hr.reset_error', 'Failed to reset password')),
  });

  const handleOpenCreate = () => {
    setEditingStaff(null);
    setFormData({ first_name: '', last_name: '', email: '', role: 'advisor', phone: '', department: '' });
    setIsModalOpen(true);
  };

  const handleOpenEdit = (staff: Staff) => {
    const [firstName, ...lastParts] = staff.name.split(' ');
    setEditingStaff(staff);
    setFormData({
      first_name: firstName,
      last_name: lastParts.join(' '),
      email: staff.email,
      role: staff.role,
      phone: '',
      department: staff.department,
    });
    setIsModalOpen(true);
  };

  const handleSubmit = () => {
    if (editingStaff) {
      updateMutation.mutate({ id: editingStaff.id, data: formData });
    } else {
      createMutation.mutate(formData);
    }
  };

  const getRoleBadge = (role: string) => {
    const r = ROLES.find(x => x.value === role);
    if (!r) return <Badge variant="outline">{role}</Badge>;
    return <Badge className={cn(r.color, 'gap-1')}><r.icon size={12} />{r.label}</Badge>;
  };

  const staff = staffData?.data || [];
  const total = staffData?.total || 0;
  const totalPages = staffData?.total_pages || 1;

  // Stats
  const activeCount = staff.filter(s => s.is_active).length;
  const advisorCount = staff.filter(s => s.role === 'advisor').length;

  return (
    <div className="space-y-6 p-6 animate-in fade-in duration-300">
      {/* Header */}
      <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
        <div>
          <h1 className="text-2xl font-black text-slate-900">{t('hr.title', 'HR Management')}</h1>
          <p className="text-slate-500 text-sm">{t('hr.subtitle', 'Manage faculty and staff directory')}</p>
        </div>
        <Button onClick={handleOpenCreate} className="gap-2">
          <UserPlus size={16} /> {t('hr.add_staff', 'Add Staff')}
        </Button>
      </div>

      {/* Stats Cards */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <Card className="bg-gradient-to-br from-indigo-50 to-white border-indigo-100">
          <CardContent className="pt-6">
            <div className="flex items-center gap-4">
              <div className="p-3 bg-indigo-100 rounded-xl"><Users className="text-indigo-600" size={24} /></div>
              <div>
                <p className="text-2xl font-black text-slate-900">{total}</p>
                <p className="text-xs text-slate-500 font-bold">{t('hr.total_staff', 'Total Staff')}</p>
              </div>
            </div>
          </CardContent>
        </Card>
        <Card className="bg-gradient-to-br from-emerald-50 to-white border-emerald-100">
          <CardContent className="pt-6">
            <div className="flex items-center gap-4">
              <div className="p-3 bg-emerald-100 rounded-xl"><Check className="text-emerald-600" size={24} /></div>
              <div>
                <p className="text-2xl font-black text-slate-900">{activeCount}</p>
                <p className="text-xs text-slate-500 font-bold">{t('hr.active_staff', 'Active')}</p>
              </div>
            </div>
          </CardContent>
        </Card>
        <Card className="bg-gradient-to-br from-blue-50 to-white border-blue-100">
          <CardContent className="pt-6">
            <div className="flex items-center gap-4">
              <div className="p-3 bg-blue-100 rounded-xl"><GraduationCap className="text-blue-600" size={24} /></div>
              <div>
                <p className="text-2xl font-black text-slate-900">{advisorCount}</p>
                <p className="text-xs text-slate-500 font-bold">{t('hr.advisors', 'Advisors')}</p>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Filters */}
      <div className="flex flex-col sm:flex-row gap-3 bg-white p-3 rounded-xl border border-slate-200">
        <div className="relative flex-1">
          <Search size={16} className="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" />
          <Input 
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            placeholder={t('hr.search_placeholder', 'Search by name or email...')}
            className="pl-9"
          />
        </div>
        <Select value={roleFilter} onValueChange={setRoleFilter}>
          <SelectTrigger className="w-40"><SelectValue placeholder="Role" /></SelectTrigger>
          <SelectContent>
            <SelectItem value="all">{t('common.all', 'All Roles')}</SelectItem>
            {ROLES.map(r => <SelectItem key={r.value} value={r.value}>{r.label}</SelectItem>)}
          </SelectContent>
        </Select>
        <Select value={activeFilter} onValueChange={setActiveFilter}>
          <SelectTrigger className="w-32"><SelectValue /></SelectTrigger>
          <SelectContent>
            <SelectItem value="true">{t('hr.active', 'Active')}</SelectItem>
            <SelectItem value="false">{t('hr.inactive', 'Inactive')}</SelectItem>
            <SelectItem value="all">{t('common.all', 'All')}</SelectItem>
          </SelectContent>
        </Select>
      </div>

      {/* Table */}
      {isLoading ? (
        <div className="flex items-center justify-center py-12">
          <Loader2 className="animate-spin text-indigo-600" size={32} />
        </div>
      ) : (
        <div className="bg-white rounded-xl border border-slate-200 overflow-hidden">
          <table className="w-full text-sm">
            <thead className="bg-slate-50 border-b border-slate-200">
              <tr>
                <th className="text-left px-4 py-3 font-bold text-slate-600">{t('hr.name', 'Name')}</th>
                <th className="text-left px-4 py-3 font-bold text-slate-600">{t('hr.role', 'Role')}</th>
                <th className="text-left px-4 py-3 font-bold text-slate-600 hidden md:table-cell">{t('hr.department', 'Department')}</th>
                <th className="text-left px-4 py-3 font-bold text-slate-600 hidden lg:table-cell">{t('hr.email', 'Email')}</th>
                <th className="text-left px-4 py-3 font-bold text-slate-600">{t('hr.status', 'Status')}</th>
                <th className="text-right px-4 py-3 font-bold text-slate-600">{t('hr.actions', 'Actions')}</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-100">
              {staff.map((s) => (
                <tr key={s.id} className="hover:bg-slate-50 transition-colors">
                  <td className="px-4 py-3">
                    <div className="flex items-center gap-3">
                      <div className="w-9 h-9 rounded-full bg-gradient-to-br from-indigo-500 to-purple-500 flex items-center justify-center text-white font-bold text-xs">
                        {s.name.split(' ').map(n => n[0]).join('').slice(0, 2)}
                      </div>
                      <div>
                        <p className="font-bold text-slate-900">{s.name}</p>
                        <p className="text-xs text-slate-400">@{s.username}</p>
                      </div>
                    </div>
                  </td>
                  <td className="px-4 py-3">{getRoleBadge(s.role)}</td>
                  <td className="px-4 py-3 hidden md:table-cell text-slate-600">{s.department || '-'}</td>
                  <td className="px-4 py-3 hidden lg:table-cell text-slate-600">{s.email || '-'}</td>
                  <td className="px-4 py-3">
                    <Badge variant={s.is_active ? 'default' : 'secondary'} className={s.is_active ? 'bg-emerald-100 text-emerald-700' : 'bg-slate-100 text-slate-500'}>
                      {s.is_active ? t('hr.active', 'Active') : t('hr.inactive', 'Inactive')}
                    </Badge>
                  </td>
                  <td className="px-4 py-3 text-right">
                    <DropdownMenu trigger={<Button variant="ghost" size="sm"><MoreHorizontal size={16} /></Button>}>
                      <DropdownItem onClick={() => handleOpenEdit(s)}>
                        {t('common.edit', 'Edit')}
                      </DropdownItem>
                      <DropdownItem onClick={() => resetPasswordMutation.mutate(s.id)}>
                        <Key size={14} className="mr-2" /> {t('hr.reset_password', 'Reset Password')}
                      </DropdownItem>
                      <DropdownItem onClick={() => toggleActiveMutation.mutate({ id: s.id, active: !s.is_active })}>
                        {s.is_active ? t('hr.deactivate', 'Deactivate') : t('hr.activate', 'Activate')}
                      </DropdownItem>
                    </DropdownMenu>
                  </td>
                </tr>
              ))}
              {staff.length === 0 && (
                <tr>
                  <td colSpan={6} className="px-4 py-12 text-center text-slate-400">
                    {t('hr.no_staff', 'No staff members found')}
                  </td>
                </tr>
              )}
            </tbody>
          </table>
          
          {/* Pagination */}
          {totalPages > 1 && (
            <div className="flex items-center justify-between px-4 py-3 border-t border-slate-200 bg-slate-50">
              <p className="text-sm text-slate-500">
                {t('hr.showing', 'Showing')} {(page - 1) * 20 + 1}-{Math.min(page * 20, total)} {t('hr.of', 'of')} {total}
              </p>
              <div className="flex gap-2">
                <Button variant="outline" size="sm" disabled={page === 1} onClick={() => setPage(p => p - 1)}>
                  {t('common.previous', 'Previous')}
                </Button>
                <Button variant="outline" size="sm" disabled={page >= totalPages} onClick={() => setPage(p => p + 1)}>
                  {t('common.next', 'Next')}
                </Button>
              </div>
            </div>
          )}
        </div>
      )}

      {/* Create/Edit Modal */}
      <Dialog open={isModalOpen} onOpenChange={setIsModalOpen}>
        <DialogContent className="max-w-md">
          <DialogHeader>
            <DialogTitle>{editingStaff ? t('hr.edit_staff', 'Edit Staff') : t('hr.add_staff', 'Add Staff')}</DialogTitle>
          </DialogHeader>
          <div className="space-y-4 py-4">
            <div className="grid grid-cols-2 gap-3">
              <div>
                <Label>{t('hr.first_name', 'First Name')}</Label>
                <Input value={formData.first_name} onChange={(e) => setFormData({ ...formData, first_name: e.target.value })} />
              </div>
              <div>
                <Label>{t('hr.last_name', 'Last Name')}</Label>
                <Input value={formData.last_name} onChange={(e) => setFormData({ ...formData, last_name: e.target.value })} />
              </div>
            </div>
            <div>
              <Label>{t('hr.email', 'Email')}</Label>
              <Input type="email" value={formData.email} onChange={(e) => setFormData({ ...formData, email: e.target.value })} />
            </div>
            <div>
              <Label>{t('hr.role', 'Role')}</Label>
              <Select value={formData.role} onValueChange={(v) => setFormData({ ...formData, role: v })}>
                <SelectTrigger><SelectValue /></SelectTrigger>
                <SelectContent>
                  {ROLES.map(r => <SelectItem key={r.value} value={r.value}>{r.label}</SelectItem>)}
                </SelectContent>
              </Select>
            </div>
            <div>
              <Label>{t('hr.department', 'Department')}</Label>
              <Input value={formData.department} onChange={(e) => setFormData({ ...formData, department: e.target.value })} />
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setIsModalOpen(false)}>{t('common.cancel', 'Cancel')}</Button>
            <Button onClick={handleSubmit} disabled={createMutation.isPending || updateMutation.isPending}>
              {(createMutation.isPending || updateMutation.isPending) && <Loader2 className="animate-spin mr-2" size={16} />}
              {editingStaff ? t('common.save', 'Save') : t('common.create', 'Create')}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Credentials Modal */}
      <Dialog open={!!credentialsModal} onOpenChange={() => setCredentialsModal(null)}>
        <DialogContent className="max-w-sm">
          <DialogHeader>
            <DialogTitle className="flex items-center gap-2">
              <Key size={20} className="text-emerald-600" />
              {t('hr.credentials', 'Login Credentials')}
            </DialogTitle>
          </DialogHeader>
          <div className="space-y-4 py-4">
            <p className="text-sm text-slate-500">{t('hr.credentials_note', 'Please save these credentials. The password cannot be retrieved later.')}</p>
            <div className="bg-slate-50 p-4 rounded-lg space-y-2 font-mono text-sm">
              <div className="flex justify-between">
                <span className="text-slate-500">{t('hr.username', 'Username')}:</span>
                <span className="font-bold">{credentialsModal?.username}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-slate-500">{t('hr.password', 'Password')}:</span>
                <span className="font-bold">{credentialsModal?.temp_password}</span>
              </div>
            </div>
          </div>
          <DialogFooter>
            <Button onClick={() => setCredentialsModal(null)}>{t('common.done', 'Done')}</Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
};

export default HRPage;
