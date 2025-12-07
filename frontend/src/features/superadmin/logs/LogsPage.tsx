import { useQuery } from '@tanstack/react-query';
import { useTranslation } from 'react-i18next';
import { useState } from 'react';
import {
  ScrollText,
  Filter,
  Calendar,
  User,
  Activity,
  Building2,
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
import { Input } from '@/components/ui/input';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { logsApi, tenantsApi, ActivityLog } from '../api';

const ACTION_COLORS: Record<string, string> = {
  login: 'bg-green-500/10 text-green-700 border-green-500/30',
  logout: 'bg-gray-500/10 text-gray-700 border-gray-500/30',
  create: 'bg-blue-500/10 text-blue-700 border-blue-500/30',
  update: 'bg-yellow-500/10 text-yellow-700 border-yellow-500/30',
  delete: 'bg-red-500/10 text-red-700 border-red-500/30',
  reset_password: 'bg-orange-500/10 text-orange-700 border-orange-500/30',
};

export function LogsPage() {
  const { t } = useTranslation('common');
  const [page, setPage] = useState(1);
  const [filters, setFilters] = useState<{
    tenant_id?: string;
    action?: string;
    entity_type?: string;
    start_date?: string;
    end_date?: string;
  }>({});

  const { data: logsData, isLoading } = useQuery({
    queryKey: ['superadmin', 'logs', page, filters],
    queryFn: () => logsApi.list({ page, limit: 50, ...filters }),
  });

  const { data: stats } = useQuery({
    queryKey: ['superadmin', 'logs', 'stats'],
    queryFn: logsApi.getStats,
  });

  const { data: actions } = useQuery({
    queryKey: ['superadmin', 'logs', 'actions'],
    queryFn: logsApi.getActions,
  });

  const { data: entityTypes } = useQuery({
    queryKey: ['superadmin', 'logs', 'entity-types'],
    queryFn: logsApi.getEntityTypes,
  });

  const { data: tenants } = useQuery({
    queryKey: ['superadmin', 'tenants'],
    queryFn: tenantsApi.list,
  });

  const logs = logsData?.data || [];
  const pagination = logsData?.pagination;

  const getActionColor = (action: string) =>
    ACTION_COLORS[action] || 'bg-gray-500/10 text-gray-700 border-gray-500/30';

  const formatDate = (dateStr: string) => {
    const date = new Date(dateStr);
    return date.toLocaleString();
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold flex items-center gap-2">
            <ScrollText className="h-6 w-6 text-violet-500" />
            {t('superadmin.logs.title', 'Activity Logs')}
          </h1>
          <p className="text-muted-foreground">
            {t('superadmin.logs.description', 'Track all platform activity')}
          </p>
        </div>
        {stats && (
          <div className="text-sm text-muted-foreground">
            {t('superadmin.logs.total', 'Total')}: <strong>{stats.total_logs.toLocaleString()}</strong>
          </div>
        )}
      </div>

      {/* Stats Cards */}
      {stats && (
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
          <div className="p-4 border rounded-lg">
            <div className="text-sm text-muted-foreground">{t('superadmin.logs.total_logs', 'Total Logs')}</div>
            <div className="text-2xl font-bold">{stats.total_logs.toLocaleString()}</div>
          </div>
          <div className="p-4 border rounded-lg">
            <div className="text-sm text-muted-foreground">{t('superadmin.logs.action_types', 'Action Types')}</div>
            <div className="text-2xl font-bold">{Object.keys(stats.logs_by_action).length}</div>
          </div>
          <div className="p-4 border rounded-lg">
            <div className="text-sm text-muted-foreground">{t('superadmin.logs.active_tenants', 'Active Tenants')}</div>
            <div className="text-2xl font-bold">{stats.logs_by_tenant.length}</div>
          </div>
          <div className="p-4 border rounded-lg">
            <div className="text-sm text-muted-foreground">{t('superadmin.logs.today', 'Today')}</div>
            <div className="text-2xl font-bold">
              {stats.recent_activity[0]?.count?.toLocaleString() || 0}
            </div>
          </div>
        </div>
      )}

      {/* Filters */}
      <div className="flex flex-wrap gap-3 p-4 border rounded-lg bg-muted/30">
        <div className="flex items-center gap-2">
          <Filter className="h-4 w-4 text-muted-foreground" />
          <span className="text-sm font-medium">{t('superadmin.logs.filters', 'Filters')}:</span>
        </div>
        
        <Select
          value={filters.tenant_id || '_all'}
          onValueChange={(v) => setFilters({ ...filters, tenant_id: v === '_all' ? undefined : v })}
        >
          <SelectTrigger className="w-[180px]">
            <Building2 className="h-4 w-4 mr-2" />
            <SelectValue placeholder={t('superadmin.logs.all_tenants', 'All Institutions')} />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="_all">{t('superadmin.logs.all_tenants', 'All Institutions')}</SelectItem>
            {tenants?.map((tenant) => (
              <SelectItem key={tenant.id} value={tenant.id}>
                {tenant.name}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>

        <Select
          value={filters.action || '_all'}
          onValueChange={(v) => setFilters({ ...filters, action: v === '_all' ? undefined : v })}
        >
          <SelectTrigger className="w-[150px]">
            <Activity className="h-4 w-4 mr-2" />
            <SelectValue placeholder={t('superadmin.logs.all_actions', 'All Actions')} />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="_all">{t('superadmin.logs.all_actions', 'All Actions')}</SelectItem>
            {actions?.map((action) => (
              <SelectItem key={action} value={action}>
                {action}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>

        <Select
          value={filters.entity_type || '_all'}
          onValueChange={(v) => setFilters({ ...filters, entity_type: v === '_all' ? undefined : v })}
        >
          <SelectTrigger className="w-[150px]">
            <SelectValue placeholder={t('superadmin.logs.all_entities', 'All Entities')} />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="_all">{t('superadmin.logs.all_entities', 'All Entities')}</SelectItem>
            {entityTypes?.map((type) => (
              <SelectItem key={type} value={type}>
                {type}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>

        <Input
          type="date"
          className="w-[150px]"
          value={filters.start_date || ''}
          onChange={(e) => setFilters({ ...filters, start_date: e.target.value })}
          placeholder="Start Date"
        />
        <Input
          type="date"
          className="w-[150px]"
          value={filters.end_date || ''}
          onChange={(e) => setFilters({ ...filters, end_date: e.target.value })}
          placeholder="End Date"
        />

        <Button
          variant="outline"
          size="sm"
          onClick={() => {
            setFilters({});
            setPage(1);
          }}
        >
          {t('superadmin.logs.clear', 'Clear')}
        </Button>
      </div>

      {/* Logs Table */}
      {isLoading ? (
        <div className="text-center py-8 text-muted-foreground">Loading...</div>
      ) : (
        <div className="border rounded-lg">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead className="w-[180px]">{t('superadmin.logs.time', 'Time')}</TableHead>
                <TableHead>{t('superadmin.logs.user', 'User')}</TableHead>
                <TableHead>{t('superadmin.logs.action', 'Action')}</TableHead>
                <TableHead>{t('superadmin.logs.entity', 'Entity')}</TableHead>
                <TableHead>{t('superadmin.logs.description', 'Description')}</TableHead>
                <TableHead>{t('superadmin.logs.institution', 'Institution')}</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {logs.map((log) => (
                <TableRow key={log.id}>
                  <TableCell className="text-sm text-muted-foreground whitespace-nowrap">
                    <div className="flex items-center gap-1">
                      <Calendar className="h-3 w-3" />
                      {formatDate(log.created_at)}
                    </div>
                  </TableCell>
                  <TableCell>
                    <div className="flex items-center gap-1">
                      <User className="h-3 w-3 text-muted-foreground" />
                      {log.username || log.user_email || '—'}
                    </div>
                  </TableCell>
                  <TableCell>
                    <Badge variant="outline" className={getActionColor(log.action)}>
                      {log.action}
                    </Badge>
                  </TableCell>
                  <TableCell>
                    {log.entity_type ? (
                      <span className="text-sm">
                        {log.entity_type}
                        {log.entity_id && (
                          <span className="text-muted-foreground text-xs ml-1">
                            ({log.entity_id.substring(0, 8)}...)
                          </span>
                        )}
                      </span>
                    ) : (
                      '—'
                    )}
                  </TableCell>
                  <TableCell className="max-w-[200px] truncate" title={log.description || ''}>
                    {log.description || '—'}
                  </TableCell>
                  <TableCell>
                    {log.tenant_name || <span className="text-muted-foreground">System</span>}
                  </TableCell>
                </TableRow>
              ))}
              {logs.length === 0 && (
                <TableRow>
                  <TableCell colSpan={6} className="text-center py-8 text-muted-foreground">
                    {t('superadmin.logs.empty', 'No logs found')}
                  </TableCell>
                </TableRow>
              )}
            </TableBody>
          </Table>
        </div>
      )}

      {/* Pagination */}
      {pagination && pagination.total_pages > 1 && (
        <div className="flex items-center justify-between">
          <div className="text-sm text-muted-foreground">
            {t('superadmin.logs.showing', 'Showing')} {(page - 1) * pagination.limit + 1} -{' '}
            {Math.min(page * pagination.limit, pagination.total)} {t('superadmin.logs.of', 'of')}{' '}
            {pagination.total}
          </div>
          <div className="flex gap-2">
            <Button
              variant="outline"
              size="sm"
              disabled={page <= 1}
              onClick={() => setPage(page - 1)}
            >
              {t('common.previous', 'Previous')}
            </Button>
            <Button
              variant="outline"
              size="sm"
              disabled={page >= pagination.total_pages}
              onClick={() => setPage(page + 1)}
            >
              {t('common.next', 'Next')}
            </Button>
          </div>
        </div>
      )}
    </div>
  );
}

export default LogsPage;
