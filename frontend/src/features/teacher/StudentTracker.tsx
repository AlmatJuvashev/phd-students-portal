import React, { useMemo, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import { useTranslation } from 'react-i18next';
import { ArrowLeft, Loader2, Search, Users } from 'lucide-react';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Progress } from '@/components/ui/progress';
import { ScrollArea } from '@/components/ui/scroll-area';
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
} from '@/components/ui/sheet';
import { cn } from '@/lib/utils';
import { getCourses } from '@/features/curriculum/api';
import { getAtRiskStudents, getCourseStudents, getStudentActivityLog, getTeacherCourses } from './api';
import { RiskLevel, StudentRiskProfile } from './types';

const riskBadgeClass = (level: RiskLevel) => {
  switch (level) {
    case 'critical':
      return 'bg-red-100 text-red-700 border-red-200';
    case 'high':
      return 'bg-amber-100 text-amber-800 border-amber-200';
    case 'medium':
      return 'bg-yellow-100 text-yellow-800 border-yellow-200';
    default:
      return 'bg-emerald-100 text-emerald-800 border-emerald-200';
  }
};

export const StudentTracker: React.FC = () => {
  const { t } = useTranslation('common');
  const navigate = useNavigate();
  const { courseId: offeringId } = useParams<{ courseId: string }>();

  const [search, setSearch] = useState('');
  const [filter, setFilter] = useState<'all' | 'at-risk' | RiskLevel>('all');
  const [selected, setSelected] = useState<StudentRiskProfile | null>(null);
  const [drawerOpen, setDrawerOpen] = useState(false);

  const offeringsQuery = useQuery({ queryKey: ['teacher', 'courses'], queryFn: getTeacherCourses });
  const catalogQuery = useQuery({
    queryKey: ['curriculum', 'courses'],
    queryFn: () => getCourses(),
    staleTime: 5 * 60 * 1000,
    retry: false,
  });

  const studentsQuery = useQuery({
    queryKey: ['teacher', 'courses', offeringId, 'students'],
    queryFn: () => getCourseStudents(offeringId!),
    enabled: !!offeringId,
  });

  const atRiskQuery = useQuery({
    queryKey: ['teacher', 'courses', offeringId, 'at-risk'],
    queryFn: () => getAtRiskStudents(offeringId!),
    enabled: !!offeringId,
  });

  const activityQuery = useQuery({
    queryKey: ['teacher', 'students', selected?.student_id, 'activity', offeringId],
    queryFn: () => getStudentActivityLog(selected!.student_id, offeringId!, 50),
    enabled: drawerOpen && !!selected?.student_id && !!offeringId,
  });

  const offering = useMemo(
    () => (offeringsQuery.data || []).find((o) => o.id === offeringId),
    [offeringsQuery.data, offeringId]
  );

  const course = useMemo(() => {
    if (!offering) return undefined;
    return (catalogQuery.data || []).find((c) => c.id === offering.course_id);
  }, [catalogQuery.data, offering]);

  const students = studentsQuery.data || [];
  const counts = useMemo(() => {
    const acc = { low: 0, medium: 0, high: 0, critical: 0 };
    for (const s of students) acc[s.risk_level]++;
    return acc;
  }, [students]);

  const filteredStudents = useMemo(() => {
    const term = search.trim().toLowerCase();
    return students
      .filter((s) => {
        if (!term) return true;
        return `${s.student_name} ${s.student_id}`.toLowerCase().includes(term);
      })
      .filter((s) => {
        if (filter === 'all') return true;
        if (filter === 'at-risk') return s.risk_level === 'high' || s.risk_level === 'critical';
        return s.risk_level === filter;
      })
      .sort((a, b) => {
        const rank: Record<RiskLevel, number> = { critical: 4, high: 3, medium: 2, low: 1 };
        return rank[b.risk_level] - rank[a.risk_level];
      });
  }, [students, search, filter]);

  if (offeringsQuery.isLoading || studentsQuery.isLoading) {
    return (
      <div className="h-full flex items-center justify-center">
        <Loader2 className="animate-spin" />
      </div>
    );
  }

  if (!offeringId) {
    return <div className="p-8 text-center text-slate-500">{t('teacher.tracker.not_found')}</div>;
  }

  if (!offering) {
    return <div className="p-8 text-center text-slate-500">{t('teacher.tracker.not_found')}</div>;
  }

  const atRiskCount = counts.high + counts.critical;

  return (
    <div className="max-w-6xl mx-auto space-y-8 animate-in fade-in duration-500">
      <div className="space-y-4">
        <button
          onClick={() => navigate(`/admin/teacher/courses/${offeringId}`)}
          className="flex items-center gap-2 text-sm font-bold text-slate-400 hover:text-slate-700 transition-colors"
        >
          <ArrowLeft size={16} /> {t('teacher.tracker.back')}
        </button>

        <div className="flex flex-col md:flex-row justify-between items-start md:items-center gap-6">
          <div>
            <div className="flex items-center gap-3 mb-1">
              <Badge variant="outline" className="bg-white border-slate-200 text-slate-500">
                {course?.code || 'COURSE'}
              </Badge>
              <span className="text-xs font-bold text-slate-400 uppercase tracking-wider">{offering.section}</span>
              <Badge variant="secondary" className="bg-slate-100 text-slate-700">
                {offering.delivery_format}
              </Badge>
            </div>
            <h1 className="text-3xl font-black text-slate-900 tracking-tight">{t('teacher.tracker.title')}</h1>
            <p className="text-slate-500 font-medium mt-1">
              {t('teacher.tracker.subtitle', { course: course?.title || offering.course_id })}
            </p>
          </div>

          <Button variant="secondary" onClick={() => navigate('/admin/teacher/grading')}>
            {t('teacher.tracker.open_grading')}
          </Button>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div className="bg-white p-4 rounded-2xl border border-slate-200 shadow-sm flex items-center gap-3">
            <div className="p-2 bg-indigo-50 text-indigo-600 rounded-lg">
              <Users size={20} />
            </div>
            <div>
              <div className="text-[10px] font-bold text-slate-400 uppercase">{t('teacher.tracker.stats.total')}</div>
              <div className="text-lg font-black text-slate-900">{students.length}</div>
            </div>
          </div>
          <div className="bg-white p-4 rounded-2xl border border-slate-200 shadow-sm flex items-center gap-3">
            <div className="p-2 bg-amber-50 text-amber-700 rounded-lg">
              <Users size={20} />
            </div>
            <div>
              <div className="text-[10px] font-bold text-slate-400 uppercase">{t('teacher.tracker.stats.at_risk')}</div>
              <div className="text-lg font-black text-slate-900">{atRiskCount}</div>
            </div>
          </div>
          <div className="bg-white p-4 rounded-2xl border border-slate-200 shadow-sm flex items-center gap-3">
            <div className="p-2 bg-emerald-50 text-emerald-700 rounded-lg">
              <Users size={20} />
            </div>
            <div>
              <div className="text-[10px] font-bold text-slate-400 uppercase">{t('teacher.tracker.stats.on_track')}</div>
              <div className="text-lg font-black text-slate-900">{counts.low}</div>
            </div>
          </div>
        </div>
      </div>

      <div className="flex flex-col md:flex-row gap-3">
        <div className="relative flex-1">
          <Search size={16} className="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" />
          <Input
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            placeholder={t('teacher.tracker.search_placeholder')}
            className="pl-9"
          />
        </div>
        <select
          value={filter}
          onChange={(e) => setFilter(e.target.value as any)}
          className="h-10 px-3 rounded-md border border-slate-200 bg-white text-sm"
        >
          <option value="all">{t('teacher.tracker.filters.all')}</option>
          <option value="at-risk">{t('teacher.tracker.filters.at_risk')}</option>
          <option value="critical">{t('teacher.tracker.filters.critical')}</option>
          <option value="high">{t('teacher.tracker.filters.high')}</option>
          <option value="medium">{t('teacher.tracker.filters.medium')}</option>
          <option value="low">{t('teacher.tracker.filters.low')}</option>
        </select>
      </div>

      <div className="bg-white border border-slate-200 rounded-2xl overflow-hidden shadow-sm animate-in fade-in">
        <table className="w-full text-sm text-left">
          <thead className="bg-slate-50 border-b border-slate-200 text-xs font-bold text-slate-500 uppercase">
            <tr>
              <th className="px-6 py-4">{t('teacher.tracker.table.student')}</th>
              <th className="px-6 py-4">{t('teacher.tracker.table.progress')}</th>
              <th className="px-6 py-4">{t('teacher.tracker.table.grade')}</th>
              <th className="px-6 py-4">{t('teacher.tracker.table.last_active')}</th>
              <th className="px-6 py-4">{t('teacher.tracker.table.status')}</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-slate-100">
            {filteredStudents.map((s) => (
              <tr
                key={s.student_id}
                className="hover:bg-slate-50 transition-colors cursor-pointer"
                onClick={() => {
                  setSelected(s);
                  setDrawerOpen(true);
                }}
              >
                <td className="px-6 py-4 font-bold text-slate-900">{s.student_name || s.student_id}</td>
                <td className="px-6 py-4">
                  <div className="flex items-center gap-3">
                    <div className="w-36">
                      <Progress value={Math.round(s.overall_progress)} className="h-2" />
                    </div>
                    <div className="text-xs font-bold text-slate-700">{Math.round(s.overall_progress)}%</div>
                  </div>
                </td>
                <td className="px-6 py-4 font-bold text-slate-900">
                  {s.average_grade > 0 ? `${Math.round(s.average_grade)}%` : '—'}
                </td>
                <td className="px-6 py-4 text-slate-600">
                  {s.days_inactive >= 999 ? '—' : t('teacher.tracker.days_ago', { count: s.days_inactive })}
                </td>
                <td className="px-6 py-4">
                  <Badge variant="outline" className={cn('uppercase text-[10px]', riskBadgeClass(s.risk_level))}>
                    {t(`teacher.tracker.risk.${s.risk_level}`)}
                  </Badge>
                </td>
              </tr>
            ))}
            {filteredStudents.length === 0 && (
              <tr>
                <td colSpan={5} className="p-12 text-center text-slate-400 italic">
                  {t('teacher.tracker.empty')}
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>

      <div className="text-xs text-slate-500">
        {t('teacher.tracker.api_hint', { count: (atRiskQuery.data || []).length })}
      </div>

      <Sheet open={drawerOpen} onOpenChange={setDrawerOpen}>
        <SheetContent side="right" className="sm:max-w-xl w-full">
          <SheetHeader>
            <SheetTitle>{selected?.student_name || selected?.student_id}</SheetTitle>
            <SheetDescription>
              {selected ? (
                <span className="inline-flex items-center gap-2">
                  <Badge variant="outline" className={cn('uppercase text-[10px]', riskBadgeClass(selected.risk_level))}>
                    {t(`teacher.tracker.risk.${selected.risk_level}`)}
                  </Badge>
                  <span className="text-xs text-slate-500">
                    {t('teacher.tracker.drawer.progress', {
                      percent: Math.round(selected.overall_progress),
                      completed: selected.assignments_completed,
                      total: selected.assignments_total,
                    })}
                  </span>
                </span>
              ) : null}
            </SheetDescription>
          </SheetHeader>

          {!selected ? null : (
            <div className="mt-6 space-y-6">
              <div className="grid grid-cols-2 gap-3">
                <div className="rounded-xl border border-slate-200 bg-white p-3">
                  <div className="text-[10px] font-bold text-slate-400 uppercase">{t('teacher.tracker.drawer.avg')}</div>
                  <div className="text-lg font-black text-slate-900">
                    {selected.average_grade > 0 ? `${Math.round(selected.average_grade)}%` : '—'}
                  </div>
                </div>
                <div className="rounded-xl border border-slate-200 bg-white p-3">
                  <div className="text-[10px] font-bold text-slate-400 uppercase">{t('teacher.tracker.drawer.last')}</div>
                  <div className="text-lg font-black text-slate-900">
                    {selected.days_inactive >= 999 ? '—' : t('teacher.tracker.days_ago', { count: selected.days_inactive })}
                  </div>
                </div>
              </div>

              <div className="rounded-2xl border border-slate-200 bg-white p-4">
                <div className="font-black text-slate-900">{t('teacher.tracker.drawer.factors')}</div>
                <ul className="mt-2 space-y-1 text-sm text-slate-700">
                  {(selected.risk_factors || []).length > 0 ? (
                    selected.risk_factors.map((f, idx) => <li key={idx}>• {f}</li>)
                  ) : (
                    <li className="text-slate-500 italic">{t('teacher.tracker.drawer.no_factors')}</li>
                  )}
                </ul>
              </div>

              <div className="rounded-2xl border border-slate-200 bg-white p-4">
                <div className="font-black text-slate-900">{t('teacher.tracker.drawer.actions')}</div>
                <div className="mt-3 flex flex-wrap gap-2">
                  {(selected.suggested_actions || []).length > 0 ? (
                    selected.suggested_actions.map((a) => (
                      <Button key={a} variant="secondary" size="sm">
                        {a}
                      </Button>
                    ))
                  ) : (
                    <div className="text-sm text-slate-500 italic">{t('teacher.tracker.drawer.no_actions')}</div>
                  )}
                </div>
              </div>

              <div className="rounded-2xl border border-slate-200 bg-white p-4">
                <div className="font-black text-slate-900">{t('teacher.tracker.drawer.activity')}</div>
                {activityQuery.isLoading ? (
                  <div className="py-8 text-center text-slate-500">
                    <Loader2 className="animate-spin mx-auto" />
                  </div>
                ) : (
                  <ScrollArea className="h-72 mt-3">
                    <div className="space-y-3 pr-4">
                      {(activityQuery.data || []).map((e, idx) => (
                        <div key={idx} className="rounded-xl border border-slate-200 p-3">
                          <div className="flex items-center justify-between gap-3">
                            <div className="font-bold text-slate-900">{e.title}</div>
                            <Badge variant="secondary" className="uppercase text-[10px]">
                              {e.kind}
                            </Badge>
                          </div>
                          <div className="mt-1 text-xs text-slate-500">
                            {new Date(e.occurred_at).toLocaleString()}
                            {e.status ? ` • ${e.status}` : ''}
                            {e.score !== null && e.score !== undefined && e.max_score
                              ? ` • ${e.score}/${e.max_score}`
                              : ''}
                          </div>
                        </div>
                      ))}
                      {(activityQuery.data || []).length === 0 && (
                        <div className="py-8 text-center text-slate-500 italic">{t('teacher.tracker.drawer.no_activity')}</div>
                      )}
                    </div>
                  </ScrollArea>
                )}
              </div>
            </div>
          )}
        </SheetContent>
      </Sheet>
    </div>
  );
};
