import React, { useMemo, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { CheckSquare, ExternalLink, FileText, Loader2, Search, X } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Input } from '@/components/ui/input';
import { cn } from '@/lib/utils';
import { AnimatePresence, motion } from 'framer-motion';
import { getTeacherSubmissions, submitGradeForSubmission } from './api';
import { ActivitySubmission } from './types';

export const TeacherGradingPage: React.FC = () => {
  const { t } = useTranslation('common');
  const queryClient = useQueryClient();

  const [filter, setFilter] = useState<'all' | 'submitted' | 'graded'>('submitted');
  const [search, setSearch] = useState('');
  const [activeSubmission, setActiveSubmission] = useState<ActivitySubmission | null>(null);
  const [currentScore, setCurrentScore] = useState<number>(0);
  const [feedback, setFeedback] = useState('');

  const submissionsQuery = useQuery({
    queryKey: ['teacher', 'submissions'],
    queryFn: getTeacherSubmissions,
  });

  const gradeMutation = useMutation({
    mutationFn: submitGradeForSubmission,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['teacher', 'submissions'] });
      queryClient.invalidateQueries({ queryKey: ['teacher', 'dashboard'] });
      setActiveSubmission(null);
    },
  });

  const submissions = submissionsQuery.data || [];

  const filtered = useMemo(() => {
    const byStatus = submissions.filter((s) => {
      if (filter === 'all') return true;
      if (filter === 'submitted') return s.status === 'SUBMITTED';
      if (filter === 'graded') return s.status === 'GRADED';
      return true;
    });
    const bySearch = byStatus.filter((s) => {
      const hay = `${s.student_name || ''} ${s.student_email || ''} ${s.student_id} ${s.activity_title || ''} ${s.activity_id}`.toLowerCase();
      return hay.includes(search.toLowerCase());
    });
    return bySearch.sort((a, b) => new Date(b.submitted_at).getTime() - new Date(a.submitted_at).getTime());
  }, [filter, search, submissions]);

  const openSubmission = (sub: ActivitySubmission) => {
    setActiveSubmission(sub);
    setCurrentScore(0);
    setFeedback('');
  };

  const handleSubmitGrade = () => {
    if (!activeSubmission) return;
    gradeMutation.mutate({
      course_offering_id: activeSubmission.course_offering_id,
      activity_id: activeSubmission.activity_id,
      student_id: activeSubmission.student_id,
      score: currentScore,
      max_score: 100,
      feedback,
    });
  };

  if (submissionsQuery.isLoading) {
    return (
      <div className="h-full flex items-center justify-center">
        <Loader2 className="animate-spin" />
      </div>
    );
  }

  return (
    <div className="h-full flex flex-col space-y-6 animate-in fade-in duration-500">
      <div className="flex justify-between items-end flex-shrink-0">
        <div>
          <h1 className="text-2xl font-black text-slate-900 tracking-tight">{t('teacher.grading.title')}</h1>
          <p className="text-slate-500 text-sm mt-1">{t('teacher.grading.subtitle')}</p>
        </div>
      </div>

      <div className="bg-white p-2 rounded-2xl border border-slate-200 shadow-sm flex flex-col sm:flex-row gap-4 flex-shrink-0">
        <div className="relative flex-1">
          <Search size={16} className="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" />
          <input
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            placeholder={t('teacher.grading.search_placeholder')}
            className="w-full h-10 pl-9 pr-4 bg-transparent border-none focus:ring-0 text-sm outline-none"
          />
        </div>
        <div className="h-10 w-px bg-slate-200 hidden sm:block" />
        <div className="flex gap-1 bg-slate-100 p-1 rounded-xl">
          {(['all', 'submitted', 'graded'] as const).map((f) => (
            <button
              key={f}
              onClick={() => setFilter(f)}
              className={cn(
                'px-4 py-1.5 rounded-lg text-xs font-bold capitalize transition-all',
                filter === f ? 'bg-white text-slate-900 shadow-sm' : 'text-slate-500 hover:text-slate-700'
              )}
            >
              {t(`teacher.grading.filters.${f}`)}
            </button>
          ))}
        </div>
      </div>

      <div className="flex-1 bg-white border border-slate-200 rounded-3xl overflow-hidden shadow-sm flex flex-col">
        <div className="overflow-y-auto flex-1">
          <table className="w-full text-sm text-left">
            <thead className="bg-slate-50 border-b border-slate-200 text-xs font-bold text-slate-500 uppercase sticky top-0 z-10">
              <tr>
                <th className="px-6 py-4">{t('teacher.grading.table.student')}</th>
                <th className="px-6 py-4">{t('teacher.grading.table.activity')}</th>
                <th className="px-6 py-4">{t('teacher.grading.table.submitted')}</th>
                <th className="px-6 py-4">{t('teacher.grading.table.status')}</th>
                <th className="px-6 py-4 text-right">{t('teacher.grading.table.action')}</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-100">
              {filtered.map((sub) => (
                <tr key={sub.id} className="hover:bg-slate-50 transition-colors group cursor-pointer" onClick={() => openSubmission(sub)}>
                  <td className="px-6 py-4 font-bold text-slate-900">{sub.student_name || sub.student_id}</td>
                  <td className="px-6 py-4 text-slate-700">{sub.activity_title || sub.activity_id}</td>
                  <td className="px-6 py-4 text-slate-500">{new Date(sub.submitted_at).toLocaleString()}</td>
                  <td className="px-6 py-4">
                    <Badge variant={sub.status === 'SUBMITTED' ? 'secondary' : 'outline'} className="uppercase text-[10px]">
                      {sub.status}
                    </Badge>
                  </td>
                  <td className="px-6 py-4 text-right">
                    <Button size="sm" variant="outline" className="opacity-0 group-hover:opacity-100 transition-opacity">
                      {t('teacher.grading.grade_now')}
                    </Button>
                  </td>
                </tr>
              ))}
              {filtered.length === 0 && (
                <tr>
                  <td colSpan={5} className="p-12 text-center text-slate-400 italic">
                    {t('teacher.grading.no_submissions')}
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </div>
      </div>

      <AnimatePresence>
        {activeSubmission && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="fixed inset-0 z-[100] bg-slate-900/60 backdrop-blur-sm flex items-center justify-center p-4"
            onClick={() => setActiveSubmission(null)}
          >
            <motion.div
              initial={{ opacity: 0, scale: 0.95 }}
              animate={{ opacity: 1, scale: 1 }}
              exit={{ opacity: 0, scale: 0.95 }}
              className="bg-white rounded-3xl shadow-2xl w-full max-w-5xl overflow-hidden flex flex-col md:flex-row"
              onClick={(e) => e.stopPropagation()}
            >
              <div className="flex-1 bg-slate-50 border-b md:border-b-0 md:border-r border-slate-200 p-8 flex flex-col">
                <div className="flex justify-between items-start">
                  <div>
                    <div className="text-xs font-bold text-slate-400 uppercase tracking-widest">{t('teacher.grading.modal.submission')}</div>
                    <div className="text-xl font-black text-slate-900 mt-1">{activeSubmission.activity_title || activeSubmission.activity_id}</div>
                    <div className="text-sm text-slate-500 mt-1">{activeSubmission.student_name || activeSubmission.student_id}</div>
                  </div>
                  <button onClick={() => setActiveSubmission(null)} className="p-2 hover:bg-slate-200 rounded-full text-slate-500">
                    <X size={18} />
                  </button>
                </div>

                <div className="flex-1 flex items-center justify-center">
                  <div className="bg-white p-12 rounded-2xl shadow-sm border border-slate-200 text-center max-w-md">
                    <FileText size={48} className="mx-auto text-slate-300 mb-4" />
                    <h3 className="font-bold text-slate-700 text-lg mb-2">{t('teacher.grading.modal.preview_title')}</h3>
                    <p className="text-sm text-slate-500 mb-6">{t('teacher.grading.modal.preview_body')}</p>
                    <Button variant="outline" disabled>
                      <ExternalLink className="mr-2 h-4 w-4" />
                      {t('teacher.grading.modal.download')}
                    </Button>
                  </div>
                </div>
              </div>

              <div className="w-full md:w-96 bg-white flex flex-col overflow-y-auto">
                <div className="p-6 space-y-6">
                  <div className="p-4 bg-slate-50 rounded-2xl border border-slate-200 space-y-3">
                    <label className="text-xs font-bold text-slate-500 uppercase">{t('teacher.grading.score')}</label>
                    <div className="flex gap-2">
                      <Input
                        type="number"
                        className="text-lg font-bold"
                        placeholder="0"
                        value={currentScore}
                        onChange={(e) => setCurrentScore(parseInt(e.target.value) || 0)}
                      />
                      <div className="flex items-center justify-center bg-white border border-slate-200 rounded-lg px-3 font-bold text-slate-400 text-sm">
                        / 100
                      </div>
                    </div>
                  </div>

                  <div className="space-y-4">
                    <label className="text-xs font-bold text-slate-500 uppercase flex items-center gap-2">
                      <CheckSquare size={14} /> {t('teacher.grading.rubric')}
                    </label>
                    <div className="text-xs text-slate-500">{t('teacher.grading.rubric_placeholder')}</div>
                  </div>

                  <div className="space-y-2">
                    <label className="text-xs font-bold text-slate-500 uppercase">{t('teacher.grading.feedback')}</label>
                    <textarea
                      className="w-full p-3 bg-slate-50 border border-slate-200 rounded-xl text-sm h-32 focus:ring-2 focus:ring-indigo-100 outline-none resize-none"
                      placeholder={t('teacher.grading.feedback_placeholder')}
                      value={feedback}
                      onChange={(e) => setFeedback(e.target.value)}
                    />
                  </div>
                </div>

                <div className="mt-auto p-6 border-t border-slate-100">
                  <Button
                    onClick={handleSubmitGrade}
                    disabled={gradeMutation.isPending}
                    className="w-full py-6 bg-indigo-600 hover:bg-indigo-700 text-white shadow-lg shadow-indigo-200"
                  >
                    {gradeMutation.isPending ? <Loader2 size={16} className="animate-spin" /> : t('teacher.grading.save_grade')}
                  </Button>
                </div>
              </div>
            </motion.div>
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  );
};

