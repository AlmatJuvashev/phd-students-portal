import React, { useMemo, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { Plus, Search, Loader2, Layers, Trash2 } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Badge } from '@/components/ui/badge';
import { cn } from '@/lib/utils';
import { createBank, deleteBank, listBanks } from './api';
import { CreateBankRequest, QuestionBank } from './types';

export const BanksPage: React.FC = () => {
  const navigate = useNavigate();
  const queryClient = useQueryClient();

  const [search, setSearch] = useState('');
  const [isCreateOpen, setIsCreateOpen] = useState(false);
  const [draft, setDraft] = useState<CreateBankRequest>({ title: '', description: '', subject: '', is_public: false });

  const banksQuery = useQuery({ queryKey: ['item-banks', 'banks'], queryFn: listBanks });

  const createMutation = useMutation({
    mutationFn: createBank,
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['item-banks', 'banks'] });
      setIsCreateOpen(false);
      setDraft({ title: '', description: '', subject: '', is_public: false });
    },
  });

  const deleteMutation = useMutation({
    mutationFn: deleteBank,
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['item-banks', 'banks'] });
    },
  });

  const banks = banksQuery.data || [];

  const filtered = useMemo(() => {
    const q = search.trim().toLowerCase();
    if (!q) return banks;
    return banks.filter((b) => `${b.title} ${b.subject || ''}`.toLowerCase().includes(q));
  }, [banks, search]);

  const onCreate = () => {
    if (!draft.title.trim()) return;
    createMutation.mutate({
      title: draft.title.trim(),
      description: draft.description?.trim() || undefined,
      subject: draft.subject?.trim() || undefined,
      is_public: !!draft.is_public,
    });
  };

  if (banksQuery.isLoading) {
    return (
      <div className="flex items-center justify-center h-full">
        <Loader2 className="animate-spin" />
      </div>
    );
  }

  return (
    <div className="space-y-6 animate-in fade-in duration-500">
      <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
        <div>
          <h1 className="text-2xl font-black text-slate-900 tracking-tight">Item Banks</h1>
          <p className="text-slate-500 text-sm mt-1">Create and manage reusable question banks.</p>
        </div>
        <Button onClick={() => setIsCreateOpen(true)}>
          <Plus className="mr-2 h-4 w-4" /> New Bank
        </Button>
      </div>

      <div className="bg-white p-2 rounded-xl border border-slate-200 shadow-sm flex gap-2 max-w-md">
        <div className="relative flex-1">
          <Search size={16} className="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" />
          <Input
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            placeholder="Search banks…"
            className="w-full h-10 pl-9 pr-4 bg-transparent border-none focus-visible:ring-0 text-sm shadow-none"
          />
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {filtered.map((b: QuestionBank) => (
          <div
            key={b.id}
            onClick={() => navigate(`/admin/item-banks/${b.id}`)}
            className="group bg-white p-6 rounded-2xl border border-slate-200 shadow-sm hover:shadow-lg hover:border-indigo-200 transition-all cursor-pointer text-left"
            role="button"
            tabIndex={0}
            onKeyDown={(e) => {
              if (e.key === 'Enter' || e.key === ' ') navigate(`/admin/item-banks/${b.id}`);
            }}
          >
            <div className="flex items-start justify-between mb-3">
              <div className="w-10 h-10 rounded-xl bg-slate-50 flex items-center justify-center text-slate-500 group-hover:bg-indigo-50 group-hover:text-indigo-600 transition-colors">
                <Layers size={18} />
              </div>
              <div className="flex items-center gap-2">
                <Badge
                  variant={b.is_public ? 'default' : 'secondary'}
                  className={cn(b.is_public ? 'bg-emerald-100 text-emerald-700 hover:bg-emerald-100' : '')}
                >
                  {b.is_public ? 'Public' : 'Private'}
                </Badge>
                <button
                  type="button"
                  className="p-2 rounded-lg hover:bg-slate-100 text-slate-400 hover:text-red-600"
                  onClick={(e) => {
                    e.stopPropagation();
                    if (!window.confirm(`Delete bank "${b.title}"? This will remove all questions in it.`)) return;
                    deleteMutation.mutate(b.id);
                  }}
                  disabled={deleteMutation.isPending}
                  aria-label="Delete bank"
                >
                  <Trash2 size={16} />
                </button>
              </div>
            </div>
            <div className="font-black text-slate-900 leading-tight">{b.title}</div>
            <div className="text-xs text-slate-500 mt-1 line-clamp-2">{b.description || '—'}</div>
            <div className="mt-4 text-[10px] font-bold text-slate-400 uppercase tracking-widest">
              {b.subject || 'General'}
            </div>
          </div>
        ))}

        {filtered.length === 0 && (
          <div className="col-span-full text-center py-12 text-slate-400 border-2 border-dashed border-slate-200 rounded-2xl">
            No banks found.
          </div>
        )}
      </div>

      {isCreateOpen && (
        <div className="fixed inset-0 z-[100] bg-slate-900/60 backdrop-blur-sm flex items-center justify-center p-4" onClick={() => setIsCreateOpen(false)}>
          <div className="bg-white w-full max-w-lg rounded-2xl shadow-2xl p-6" onClick={(e) => e.stopPropagation()}>
            <div className="flex justify-between items-center mb-6">
              <h3 className="text-lg font-bold text-slate-900">Create Bank</h3>
              <button className="text-slate-400 hover:text-slate-600" onClick={() => setIsCreateOpen(false)}>
                ✕
              </button>
            </div>

            <div className="space-y-4">
              <div className="space-y-1">
                <label className="text-xs font-bold text-slate-500 uppercase">Title</label>
                <Input value={draft.title} onChange={(e) => setDraft({ ...draft, title: e.target.value })} placeholder="e.g., Anatomy MCQ Bank" />
              </div>
              <div className="space-y-1">
                <label className="text-xs font-bold text-slate-500 uppercase">Subject</label>
                <Input value={draft.subject} onChange={(e) => setDraft({ ...draft, subject: e.target.value })} placeholder="e.g., Anatomy" />
              </div>
              <div className="space-y-1">
                <label className="text-xs font-bold text-slate-500 uppercase">Description</label>
                <textarea
                  value={draft.description}
                  onChange={(e) => setDraft({ ...draft, description: e.target.value })}
                  className="w-full p-3 bg-slate-50 border border-slate-200 rounded-xl text-sm h-24 focus:ring-2 focus:ring-indigo-100 outline-none resize-none"
                  placeholder="Optional description…"
                />
              </div>

              <label className="flex items-center gap-2 text-sm text-slate-700">
                <input
                  type="checkbox"
                  checked={!!draft.is_public}
                  onChange={(e) => setDraft({ ...draft, is_public: e.target.checked })}
                />
                Make public (tenant-wide)
              </label>
            </div>

            <div className="flex justify-end gap-2 mt-8 pt-4 border-t border-slate-100">
              <Button variant="outline" onClick={() => setIsCreateOpen(false)}>
                Cancel
              </Button>
              <Button onClick={onCreate} disabled={createMutation.isPending || !draft.title.trim()}>
                {createMutation.isPending ? <Loader2 className="animate-spin h-4 w-4" /> : 'Create'}
              </Button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};
