import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { useTranslation } from 'react-i18next';
import { useState, useMemo, useCallback } from 'react';
import {
  Settings,
  Save,
  Plus,
  Trash2,
  Info,
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
} from '@/components/ui/dialog';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { Switch } from '@/components/ui/switch';
import { settingsApi, Setting } from '../api';
import { SuperadminTableToolbar, SuperadminPagination, ConfirmDialog } from '../components';

const CATEGORY_COLORS: Record<string, string> = {
  system: 'bg-red-500/10 text-red-700 border-red-500/30',
  security: 'bg-orange-500/10 text-orange-700 border-orange-500/30',
  limits: 'bg-blue-500/10 text-blue-700 border-blue-500/30',
  defaults: 'bg-green-500/10 text-green-700 border-green-500/30',
  general: 'bg-gray-500/10 text-gray-700 border-gray-500/30',
};

function SettingEditor({
  setting,
  onClose,
}: {
  setting?: Setting | { key: string };
  onClose: () => void;
}) {
  const { t } = useTranslation('common');
  const queryClient = useQueryClient();
  const isNew = !('value' in (setting || {}));
  
  const [key, setKey] = useState(setting?.key || '');
  const [valueStr, setValueStr] = useState(
    'value' in (setting || {}) ? JSON.stringify((setting as Setting).value) : ''
  );
  const [description, setDescription] = useState(
    'description' in (setting || {}) ? (setting as Setting).description || '' : ''
  );
  const [category, setCategory] = useState(
    'category' in (setting || {}) ? (setting as Setting).category : 'general'
  );

  const mutation = useMutation({
    mutationFn: async () => {
      let parsedValue;
      try {
        parsedValue = JSON.parse(valueStr);
      } catch {
        parsedValue = valueStr;
      }
      return settingsApi.update(key, {
        value: parsedValue,
        description: description || undefined,
        category,
      });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['superadmin', 'settings'] });
      onClose();
    },
  });

  return (
    <div className="space-y-4">
      <div className="space-y-2">
        <Label htmlFor="key">{t('superadmin.settings.key', 'Key')}</Label>
        <Input
          id="key"
          value={key}
          onChange={(e) => setKey(e.target.value)}
          disabled={!isNew}
          placeholder="setting_key"
        />
      </div>

      <div className="space-y-2">
        <Label htmlFor="value">{t('superadmin.settings.value', 'Value')} (JSON)</Label>
        <Textarea
          id="value"
          value={valueStr}
          onChange={(e) => setValueStr(e.target.value)}
          placeholder='true, "string", 123, or {"key": "value"}'
          rows={3}
        />
        <p className="text-xs text-muted-foreground">
          {t('superadmin.settings.value_hint', 'Enter a valid JSON value: string, number, boolean, object, or array')}
        </p>
      </div>

      <div className="space-y-2">
        <Label htmlFor="description">{t('superadmin.settings.description', 'Description')}</Label>
        <Input
          id="description"
          value={description}
          onChange={(e) => setDescription(e.target.value)}
          placeholder="What this setting does"
        />
      </div>

      <div className="space-y-2">
        <Label htmlFor="category">{t('superadmin.settings.category', 'Category')}</Label>
        <Select value={category} onValueChange={setCategory}>
          <SelectTrigger>
            <SelectValue />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="general">General</SelectItem>
            <SelectItem value="system">System</SelectItem>
            <SelectItem value="security">Security</SelectItem>
            <SelectItem value="limits">Limits</SelectItem>
            <SelectItem value="defaults">Defaults</SelectItem>
          </SelectContent>
        </Select>
      </div>

      <div className="flex justify-end gap-2 pt-4">
        <Button variant="outline" onClick={onClose}>
          {t('common.cancel', 'Cancel')}
        </Button>
        <Button onClick={() => mutation.mutate()} disabled={!key || mutation.isPending}>
          {mutation.isPending ? '...' : t('common.save', 'Save')}
        </Button>
      </div>
    </div>
  );
}

export function SettingsPage() {
  const { t } = useTranslation('common');
  const queryClient = useQueryClient();
  const [editingSetting, setEditingSetting] = useState<Setting | { key: string } | null>(null);
  const [selectedCategory, setSelectedCategory] = useState<string>('');
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const [settingToDelete, setSettingToDelete] = useState<Setting | null>(null);

  // Table state
  const [searchQuery, setSearchQuery] = useState('');
  const [sortColumn, setSortColumn] = useState<string>('key');
  const [sortDirection, setSortDirection] = useState<'asc' | 'desc'>('asc');
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(25);

  const { data: settings, isLoading } = useQuery({
    queryKey: ['superadmin', 'settings', selectedCategory],
    queryFn: () => settingsApi.list(selectedCategory || undefined),
  });

  const { data: categories } = useQuery({
    queryKey: ['superadmin', 'settings', 'categories'],
    queryFn: settingsApi.getCategories,
  });

  // Filter, sort, and paginate settings
  const processedData = useMemo(() => {
    if (!settings) return { items: [], total: 0, totalPages: 0 };

    let filtered = [...settings];

    // Search filter
    if (searchQuery) {
      const query = searchQuery.toLowerCase();
      filtered = filtered.filter(
        (s) =>
          s.key.toLowerCase().includes(query) ||
          s.description?.toLowerCase().includes(query) ||
          s.category.toLowerCase().includes(query)
      );
    }

    // Sort
    filtered.sort((a, b) => {
      let aVal: any, bVal: any;
      switch (sortColumn) {
        case 'key':
          aVal = a.key.toLowerCase();
          bVal = b.key.toLowerCase();
          break;
        case 'category':
          aVal = a.category.toLowerCase();
          bVal = b.category.toLowerCase();
          break;
        case 'updated':
          aVal = new Date(a.updated_at).getTime();
          bVal = new Date(b.updated_at).getTime();
          break;
        default:
          aVal = a.key;
          bVal = b.key;
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
  }, [settings, searchQuery, sortColumn, sortDirection, currentPage, pageSize]);

  const handleSort = useCallback((column: string) => {
    if (sortColumn === column) {
      setSortDirection((d) => (d === 'asc' ? 'desc' : 'asc'));
    } else {
      setSortColumn(column);
      setSortDirection('asc');
    }
  }, [sortColumn]);

  const handleSearchChange = useCallback((value: string) => {
    setSearchQuery(value);
    setCurrentPage(1);
  }, []);

  const handleClearFilters = useCallback(() => {
    setSearchQuery('');
    setSelectedCategory('');
    setCurrentPage(1);
  }, []);

  const handlePageSizeChange = useCallback((size: number) => {
    setPageSize(size);
    setCurrentPage(1);
  }, []);

  const deleteMutation = useMutation({
    mutationFn: settingsApi.delete,
    onSuccess: () =>
      queryClient.invalidateQueries({ queryKey: ['superadmin', 'settings'] }),
  });

  const getCategoryColor = (cat: string) =>
    CATEGORY_COLORS[cat] || CATEGORY_COLORS.general;

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

  const formatValue = (value: unknown): string => {
    if (typeof value === 'boolean') return value ? 'true' : 'false';
    if (typeof value === 'string') return value;
    if (typeof value === 'number') return String(value);
    return JSON.stringify(value);
  };

  const renderValuePreview = (setting: Setting) => {
    const value = setting.value;
    
    if (typeof value === 'boolean') {
      return (
        <Switch checked={value} disabled className="pointer-events-none" />
      );
    }
    
    if (typeof value === 'number') {
      return <span className="font-mono">{value}</span>;
    }
    
    if (typeof value === 'string') {
      return (
        <span className="font-mono text-sm max-w-[200px] truncate inline-block">
          "{value}"
        </span>
      );
    }
    
    return (
      <code className="text-xs bg-muted px-2 py-1 rounded max-w-[200px] truncate inline-block">
        {JSON.stringify(value)}
      </code>
    );
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold flex items-center gap-2">
            <Settings className="h-6 w-6 text-violet-500" />
            {t('superadmin.settings.title', 'Global Settings')}
          </h1>
          <p className="text-muted-foreground">
            {t('superadmin.settings.description', 'Configure platform-wide settings')}
          </p>
        </div>
        <Button onClick={() => setEditingSetting({ key: '' })}>
          <Plus className="h-4 w-4 mr-2" />
          {t('superadmin.settings.add', 'Add Setting')}
        </Button>
      </div>

      {/* Edit Dialog */}
      <Dialog open={!!editingSetting} onOpenChange={() => setEditingSetting(null)}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>
              {'value' in (editingSetting || {})
                ? t('superadmin.settings.edit', 'Edit Setting')
                : t('superadmin.settings.add', 'Add Setting')}
            </DialogTitle>
          </DialogHeader>
          {editingSetting && (
            <SettingEditor
              setting={editingSetting}
              onClose={() => setEditingSetting(null)}
            />
          )}
        </DialogContent>
      </Dialog>

      {/* Search + Category Filter */}
      <SuperadminTableToolbar
        searchPlaceholder={t('superadmin.settings.search', 'Search settings...')}
        searchValue={searchQuery}
        onSearchChange={handleSearchChange}
        filters={[{
          key: 'category',
          label: 'Category',
          options: [
            { value: 'all', label: 'All Categories' },
            ...(categories?.map(c => ({ value: c, label: c.charAt(0).toUpperCase() + c.slice(1) })) || []),
          ],
        }]}
        filterValues={{ category: selectedCategory || 'all' }}
        onFilterChange={(key, value) => setSelectedCategory(value === 'all' ? '' : value)}
        onClearFilters={handleClearFilters}
      />

      {/* Settings Table */}
      {isLoading ? (
        <div className="text-center py-8 text-muted-foreground">Loading...</div>
      ) : (
        <>
          <div className="border rounded-lg">
            <Table>
              <TableHeader>
                <TableRow>
                  <SortableHeader column="key">{t('superadmin.settings.key', 'Key')}</SortableHeader>
                  <TableHead>{t('superadmin.settings.value', 'Value')}</TableHead>
                  <TableHead>{t('superadmin.settings.description', 'Description')}</TableHead>
                  <SortableHeader column="category">{t('superadmin.settings.category', 'Category')}</SortableHeader>
                  <TableHead className="w-24"></TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {processedData.items.map((setting) => (
                  <TableRow key={setting.key}>
                    <TableCell className="font-mono text-sm font-medium">
                      {setting.key}
                    </TableCell>
                    <TableCell>{renderValuePreview(setting)}</TableCell>
                    <TableCell className="max-w-[250px]">
                      {setting.description ? (
                        <div className="flex items-start gap-1">
                          <Info className="h-4 w-4 text-muted-foreground shrink-0 mt-0.5" />
                          <span className="text-sm text-muted-foreground">
                            {setting.description}
                          </span>
                        </div>
                      ) : (
                        <span className="text-muted-foreground">—</span>
                      )}
                    </TableCell>
                    <TableCell>
                      <Badge variant="outline" className={getCategoryColor(setting.category)}>
                        {setting.category}
                      </Badge>
                    </TableCell>
                    <TableCell>
                      <div className="flex items-center gap-1">
                        <Button
                          variant="ghost"
                          size="icon"
                          onClick={() => setEditingSetting(setting)}
                          title="Edit"
                        >
                          <Save className="h-4 w-4" />
                        </Button>
                        <Button
                          variant="ghost"
                          size="icon"
                          onClick={() => {
                            setSettingToDelete(setting);
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
                    <TableCell colSpan={5} className="text-center py-8 text-muted-foreground">
                      {searchQuery || selectedCategory
                        ? t('superadmin.settings.no_results', 'No matching settings found')
                        : t('superadmin.settings.empty', 'No settings found')}
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
        title={t('superadmin.settings.delete_title', 'Delete Setting')}
        description={t('superadmin.settings.delete_description', `Are you sure you want to delete the setting "${settingToDelete?.key}"? This action cannot be undone.`)}
        confirmLabel={t('common.delete', 'Delete')}
        onConfirm={() => {
          if (settingToDelete) {
            deleteMutation.mutate(settingToDelete.key);
          }
        }}
        loading={deleteMutation.isPending}
      />
    </div>
  );
}

export default SettingsPage;
