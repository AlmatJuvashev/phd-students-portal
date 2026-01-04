import React, { useMemo, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { ArrowLeft, Loader2, Plus, Search, Trash2 } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Badge } from '@/components/ui/badge';
import { cn } from '@/lib/utils';
import { createQuestion, listBanks, listQuestions } from './api';
import { CreateQuestionRequest, Question, QuestionOption, QuestionType } from './types';

const QUESTION_TYPES: QuestionType[] = ['MCQ', 'MRQ', 'TRUE_FALSE', 'TEXT', 'LIKERT'];

export const BankItemsPage: React.FC = () => {
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const { bankId } = useParams<{ bankId: string }>();

  const [search, setSearch] = useState('');
  const [isCreateOpen, setIsCreateOpen] = useState(false);

  const [draft, setDraft] = useState<CreateQuestionRequest>({
    type: 'MCQ',
    stem: '',
    points_default: 1,
    difficulty_level: 'MEDIUM',
    options: [
      { text: 'Option A', is_correct: true },
      { text: 'Option B', is_correct: false },
    ],
  });

  const banksQuery = useQuery({ queryKey: ['item-banks', 'banks'], queryFn: listBanks });
  const questionsQuery = useQuery({
    queryKey: ['item-banks', 'banks', bankId, 'items'],
    queryFn: () => listQuestions(bankId!),
    enabled: !!bankId,
  });

  const createMutation = useMutation({
    mutationFn: (data: CreateQuestionRequest) => createQuestion(bankId!, data),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['item-banks', 'banks', bankId, 'items'] });
      setIsCreateOpen(false);
      setDraft({
        type: 'MCQ',
        stem: '',
        points_default: 1,
        difficulty_level: 'MEDIUM',
        options: [
          { text: 'Option A', is_correct: true },
          { text: 'Option B', is_correct: false },
        ],
      });
    },
  });

  const bank = useMemo(() => (banksQuery.data || []).find((b) => b.id === bankId), [banksQuery.data, bankId]);
  const questions = questionsQuery.data || [];

  const filtered = useMemo(() => {
    const q = search.trim().toLowerCase();
    if (!q) return questions;
    return questions.filter((it) => `${it.stem} ${it.type}`.toLowerCase().includes(q));
  }, [questions, search]);

  const showOptions = draft.type === 'MCQ' || draft.type === 'MRQ' || draft.type === 'TRUE_FALSE';

  const setOptionCorrect = (idx: number, isCorrect: boolean) => {
    if (!draft.options) return;
    const next = draft.options.map((o, i) => ({ ...o, is_correct: draft.type === 'MRQ' ? (i === idx ? isCorrect : o.is_correct) : i === idx }));
    setDraft({ ...draft, options: next });
  };

  const updateOption = (idx: number, updates: Partial<QuestionOption>) => {
    const options = (draft.options || []).map((o, i) => (i === idx ? { ...o, ...updates } : o));
    setDraft({ ...draft, options });
  };

  const addOption = () => {
    const options = [...(draft.options || []), { text: `Option ${String.fromCharCode(65 + (draft.options?.length || 0))}`, is_correct: false }];
    setDraft({ ...draft, options });
  };

  const removeOption = (idx: number) => {
    const options = (draft.options || []).filter((_, i) => i !== idx);
    setDraft({ ...draft, options });
  };

  const onCreate = () => {
    if (!draft.stem.trim()) return;
    const payload: CreateQuestionRequest = {
      type: draft.type,
      stem: draft.stem.trim(),
      points_default: draft.points_default || 1,
      difficulty_level: draft.difficulty_level,
    };
    if (showOptions) {
      payload.options = (draft.options || [])
        .map((o, i) => ({
          text: o.text.trim(),
          is_correct: !!o.is_correct,
          sort_order: i,
        }))
        .filter((o) => o.text);
    }
    createMutation.mutate(payload);
  };

  if (banksQuery.isLoading || questionsQuery.isLoading) {
    return (
      <div className="flex items-center justify-center h-full">
        <Loader2 className="animate-spin" />
      </div>
    );
  }

  if (!bankId) return <div className="p-8 text-slate-500">Bank not specified.</div>;

  return (
    <div className="space-y-6 animate-in fade-in duration-500">
      <div className="flex items-center gap-3">
        <Button variant="ghost" size="icon" onClick={() => navigate('/admin/item-banks')}>
          <ArrowLeft size={16} />
        </Button>
        <div className="flex-1 min-w-0">
          <div className="text-xs text-slate-500">Item Bank</div>
          <h1 className="text-xl font-black text-slate-900 truncate">{bank?.title || bankId}</h1>
        </div>
        <Button onClick={() => setIsCreateOpen(true)}>
          <Plus className="mr-2 h-4 w-4" />
          New Question
        </Button>
      </div>

      <div className="bg-white p-2 rounded-xl border border-slate-200 shadow-sm flex gap-2 max-w-md">
        <div className="relative flex-1">
          <Search size={16} className="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" />
          <Input
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            placeholder="Search questions…"
            className="w-full h-10 pl-9 pr-4 bg-transparent border-none focus-visible:ring-0 text-sm shadow-none"
          />
        </div>
      </div>

      <div className="bg-white border border-slate-200 rounded-2xl overflow-hidden shadow-sm overflow-x-auto">
        <table className="w-full text-sm text-left whitespace-nowrap">
          <thead className="bg-slate-50 border-b border-slate-200 text-xs font-bold text-slate-500 uppercase">
            <tr>
              <th className="px-6 py-4">Type</th>
              <th className="px-6 py-4">Stem</th>
              <th className="px-6 py-4">Points</th>
              <th className="px-6 py-4">Difficulty</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-slate-100">
            {filtered.map((q: Question) => (
              <tr key={q.id} className="hover:bg-slate-50 transition-colors">
                <td className="px-6 py-4">
                  <Badge variant="secondary" className="font-mono">
                    {q.type}
                  </Badge>
                </td>
                <td className="px-6 py-4 max-w-[600px] truncate">{q.stem}</td>
                <td className="px-6 py-4">{q.points_default}</td>
                <td className="px-6 py-4">{q.difficulty_level || '—'}</td>
              </tr>
            ))}
            {filtered.length === 0 && (
              <tr>
                <td colSpan={4} className="px-6 py-12 text-center text-slate-400 italic">
                  No questions found.
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>

      {isCreateOpen && (
        <div className="fixed inset-0 z-[100] bg-slate-900/60 backdrop-blur-sm flex items-center justify-center p-4" onClick={() => setIsCreateOpen(false)}>
          <div className="bg-white w-full max-w-2xl rounded-2xl shadow-2xl p-6" onClick={(e) => e.stopPropagation()}>
            <div className="flex justify-between items-center mb-6">
              <h3 className="text-lg font-bold text-slate-900">Create Question</h3>
              <button className="text-slate-400 hover:text-slate-600" onClick={() => setIsCreateOpen(false)}>
                ✕
              </button>
            </div>

            <div className="space-y-4">
              <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
                <div className="space-y-1 sm:col-span-2">
                  <label className="text-xs font-bold text-slate-500 uppercase">Type</label>
                  <select
                    value={draft.type}
                    onChange={(e) => setDraft({ ...draft, type: e.target.value as QuestionType })}
                    className="w-full h-10 px-3 bg-white border border-slate-200 rounded-lg text-sm"
                  >
                    {QUESTION_TYPES.map((t) => (
                      <option key={t} value={t}>
                        {t}
                      </option>
                    ))}
                  </select>
                </div>
                <div className="space-y-1">
                  <label className="text-xs font-bold text-slate-500 uppercase">Points</label>
                  <Input
                    type="number"
                    value={draft.points_default ?? 1}
                    onChange={(e) => setDraft({ ...draft, points_default: parseFloat(e.target.value) || 1 })}
                  />
                </div>
              </div>

              <div className="space-y-1">
                <label className="text-xs font-bold text-slate-500 uppercase">Difficulty</label>
                <select
                  value={draft.difficulty_level || ''}
                  onChange={(e) => setDraft({ ...draft, difficulty_level: (e.target.value || undefined) as any })}
                  className="w-full h-10 px-3 bg-white border border-slate-200 rounded-lg text-sm"
                >
                  <option value="">(none)</option>
                  <option value="EASY">EASY</option>
                  <option value="MEDIUM">MEDIUM</option>
                  <option value="HARD">HARD</option>
                </select>
              </div>

              <div className="space-y-1">
                <label className="text-xs font-bold text-slate-500 uppercase">Stem</label>
                <textarea
                  value={draft.stem}
                  onChange={(e) => setDraft({ ...draft, stem: e.target.value })}
                  className="w-full p-3 bg-slate-50 border border-slate-200 rounded-xl text-sm h-28 focus:ring-2 focus:ring-indigo-100 outline-none resize-none"
                  placeholder="Write the question…"
                />
              </div>

              {showOptions && (
                <div className="space-y-2">
                  <div className="flex justify-between items-center">
                    <label className="text-xs font-bold text-slate-500 uppercase">Options</label>
                    <Button size="sm" variant="outline" onClick={addOption}>
                      <Plus className="mr-2 h-4 w-4" />
                      Add
                    </Button>
                  </div>

                  <div className="space-y-2">
                    {(draft.options || []).map((opt, idx) => (
                      <div key={idx} className="flex items-center gap-2">
                        <button
                          type="button"
                          onClick={() => setOptionCorrect(idx, draft.type === 'MRQ' ? !opt.is_correct : true)}
                          className={cn(
                            'w-8 h-8 rounded-lg border flex items-center justify-center text-xs font-bold',
                            opt.is_correct ? 'bg-emerald-50 border-emerald-200 text-emerald-700' : 'bg-white border-slate-200 text-slate-400'
                          )}
                          title={draft.type === 'MRQ' ? 'Toggle correct' : 'Mark correct'}
                        >
                          ✓
                        </button>
                        <Input value={opt.text} onChange={(e) => updateOption(idx, { text: e.target.value })} />
                        <Button variant="ghost" size="icon" onClick={() => removeOption(idx)} disabled={(draft.options || []).length <= 2}>
                          <Trash2 size={16} className="text-slate-400" />
                        </Button>
                      </div>
                    ))}
                  </div>
                  <div className="text-xs text-slate-500">
                    {draft.type === 'MRQ' ? 'Mark all correct options.' : 'Select exactly one correct option.'}
                  </div>
                </div>
              )}
            </div>

            <div className="flex justify-end gap-2 mt-8 pt-4 border-t border-slate-100">
              <Button variant="outline" onClick={() => setIsCreateOpen(false)}>
                Cancel
              </Button>
              <Button onClick={onCreate} disabled={createMutation.isPending || !draft.stem.trim()}>
                {createMutation.isPending ? <Loader2 className="animate-spin h-4 w-4" /> : 'Create'}
              </Button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

