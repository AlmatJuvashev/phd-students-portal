import React, { useMemo, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { useQuery } from '@tanstack/react-query';
import { ArrowLeft, ExternalLink, Loader2, Users } from 'lucide-react';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { cn } from '@/lib/utils';
import { getCourses } from '@/features/curriculum/api';
import { getCourseGradebook, getCourseRoster, getTeacherCourses } from './api';

export const TeacherCourseDetail: React.FC = () => {
  const { t } = useTranslation('common');
  const navigate = useNavigate();
  const { courseId: offeringId } = useParams<{ courseId: string }>();
  const [activeTab, setActiveTab] = useState<'overview' | 'roster' | 'gradebook'>('overview');

  const offeringsQuery = useQuery({ queryKey: ['teacher', 'courses'], queryFn: getTeacherCourses });
  const catalogQuery = useQuery({
    queryKey: ['curriculum', 'courses'],
    queryFn: () => getCourses(),
    staleTime: 5 * 60 * 1000,
    retry: false,
  });

  const rosterQuery = useQuery({
    queryKey: ['teacher', 'courses', offeringId, 'roster'],
    queryFn: () => getCourseRoster(offeringId!),
    enabled: !!offeringId && activeTab === 'roster',
  });

  const gradebookQuery = useQuery({
    queryKey: ['teacher', 'courses', offeringId, 'gradebook'],
    queryFn: () => getCourseGradebook(offeringId!),
    enabled: !!offeringId && activeTab === 'gradebook',
  });

  const offering = useMemo(
    () => (offeringsQuery.data || []).find((o) => o.id === offeringId),
    [offeringsQuery.data, offeringId]
  );

  const course = useMemo(() => {
    if (!offering) return undefined;
    return (catalogQuery.data || []).find((c) => c.id === offering.course_id);
  }, [catalogQuery.data, offering]);

  const gradebook = gradebookQuery.data || [];
  const avgPercent = useMemo(() => {
    if (gradebook.length === 0) return null;
    const percents = gradebook.map((e) => (e.max_score > 0 ? (e.score / e.max_score) * 100 : 0));
    return Math.round(percents.reduce((a, b) => a + b, 0) / percents.length);
  }, [gradebook]);

  if (offeringsQuery.isLoading) {
    return (
      <div className="h-full flex items-center justify-center">
        <Loader2 className="animate-spin" />
      </div>
    );
  }

  if (!offeringId || !offering) {
    return <div className="p-8 text-center text-slate-500">{t('teacher.detail.not_found')}</div>;
  }

  return (
    <div className="max-w-6xl mx-auto space-y-8 animate-in fade-in duration-500">
      <div className="space-y-4">
        <button
          onClick={() => navigate('/admin/teacher/courses')}
          className="flex items-center gap-2 text-sm font-bold text-slate-400 hover:text-slate-700 transition-colors"
        >
          <ArrowLeft size={16} /> {t('teacher.detail.back')}
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
            <h1 className="text-3xl font-black text-slate-900 tracking-tight">{course?.title || offering.course_id}</h1>
          </div>

          <div className="flex gap-2">
            <Button variant="secondary" onClick={() => navigate(`/admin/studio/courses/${offering.course_id}/builder`)}>
              <ExternalLink className="mr-2 h-4 w-4" />
              {t('teacher.detail.edit_content')}
            </Button>
            <Button onClick={() => navigate('/admin/teacher/grading')}>{t('teacher.detail.open_grading')}</Button>
          </div>
        </div>

        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
          <div className="bg-white p-4 rounded-2xl border border-slate-200 shadow-sm flex items-center gap-3">
            <div className="p-2 bg-indigo-50 text-indigo-600 rounded-lg">
              <Users size={20} />
            </div>
            <div>
              <div className="text-[10px] font-bold text-slate-400 uppercase">{t('teacher.detail.stats.enrolled')}</div>
              <div className="text-lg font-black text-slate-900">{offering.current_enrolled || 0}</div>
            </div>
          </div>
          <div className="bg-white p-4 rounded-2xl border border-slate-200 shadow-sm flex items-center gap-3">
            <div className="p-2 bg-emerald-50 text-emerald-600 rounded-lg">
              <Users size={20} />
            </div>
            <div>
              <div className="text-[10px] font-bold text-slate-400 uppercase">{t('teacher.detail.stats.roster')}</div>
              <div className="text-lg font-black text-slate-900">{rosterQuery.data?.length ?? '—'}</div>
            </div>
          </div>
          <div className="bg-white p-4 rounded-2xl border border-slate-200 shadow-sm flex items-center gap-3">
            <div className="p-2 bg-amber-50 text-amber-600 rounded-lg">
              <Users size={20} />
            </div>
            <div>
              <div className="text-[10px] font-bold text-slate-400 uppercase">{t('teacher.detail.stats.gradebook')}</div>
              <div className="text-lg font-black text-slate-900">{gradebook.length || '—'}</div>
            </div>
          </div>
          <div className="bg-white p-4 rounded-2xl border border-slate-200 shadow-sm flex items-center gap-3">
            <div className="p-2 bg-slate-50 text-slate-700 rounded-lg">
              <Users size={20} />
            </div>
            <div>
              <div className="text-[10px] font-bold text-slate-400 uppercase">{t('teacher.detail.stats.avg')}</div>
              <div className="text-lg font-black text-slate-900">{avgPercent !== null ? `${avgPercent}%` : '—'}</div>
            </div>
          </div>
        </div>
      </div>

      <div className="flex gap-2 bg-white p-2 rounded-2xl border border-slate-200 shadow-sm">
        {(['overview', 'roster', 'gradebook'] as const).map((tab) => (
          <button
            key={tab}
            onClick={() => setActiveTab(tab)}
            className={cn(
              'px-4 py-2 rounded-xl text-xs font-bold uppercase tracking-wider transition-all',
              activeTab === tab ? 'bg-indigo-600 text-white shadow' : 'text-slate-500 hover:text-slate-700 hover:bg-slate-50'
            )}
          >
            {t(`teacher.detail.tabs.${tab}`)}
          </button>
        ))}
      </div>

      {activeTab === 'overview' && (
        <div className="bg-white border border-slate-200 rounded-2xl p-8 text-slate-600">
          <div className="font-bold text-slate-900 mb-2">{t('teacher.detail.overview_title')}</div>
          <div className="text-sm">{t('teacher.detail.overview_body')}</div>
        </div>
      )}

      {activeTab === 'roster' && (
        <div className="bg-white border border-slate-200 rounded-2xl overflow-hidden shadow-sm animate-in fade-in">
          {rosterQuery.isLoading ? (
            <div className="p-12 text-center">
              <Loader2 className="animate-spin mx-auto" />
            </div>
          ) : (
            <table className="w-full text-sm text-left">
              <thead className="bg-slate-50 border-b border-slate-200 text-xs font-bold text-slate-500 uppercase">
                <tr>
                  <th className="px-6 py-4">{t('teacher.detail.table.student')}</th>
                  <th className="px-6 py-4">{t('teacher.detail.table.email')}</th>
                  <th className="px-6 py-4">{t('teacher.detail.table.status')}</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-slate-100">
                {(rosterQuery.data || []).map((e) => (
                  <tr key={e.id} className="hover:bg-slate-50 transition-colors">
                    <td className="px-6 py-4 font-bold text-slate-900">{e.student_name || e.student_id}</td>
                    <td className="px-6 py-4 text-slate-600">{e.student_email || '—'}</td>
                    <td className="px-6 py-4">
                      <Badge variant="secondary" className="uppercase text-[10px]">
                        {e.status}
                      </Badge>
                    </td>
                  </tr>
                ))}
                {(rosterQuery.data || []).length === 0 && (
                  <tr>
                    <td colSpan={3} className="p-12 text-center text-slate-400 italic">
                      {t('teacher.detail.no_students')}
                    </td>
                  </tr>
                )}
              </tbody>
            </table>
          )}
        </div>
      )}

      {activeTab === 'gradebook' && (
        <div className="bg-white border border-slate-200 rounded-2xl overflow-hidden shadow-sm animate-in fade-in">
          {gradebookQuery.isLoading ? (
            <div className="p-12 text-center">
              <Loader2 className="animate-spin mx-auto" />
            </div>
          ) : (
            <table className="w-full text-sm text-left">
              <thead className="bg-slate-50 border-b border-slate-200 text-xs font-bold text-slate-500 uppercase">
                <tr>
                  <th className="px-6 py-4">{t('teacher.detail.gradebook.activity')}</th>
                  <th className="px-6 py-4">{t('teacher.detail.gradebook.student')}</th>
                  <th className="px-6 py-4">{t('teacher.detail.gradebook.score')}</th>
                  <th className="px-6 py-4">{t('teacher.detail.gradebook.grade')}</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-slate-100">
                {gradebook.map((e) => (
                  <tr key={e.id} className="hover:bg-slate-50 transition-colors">
                    <td className="px-6 py-4 font-mono text-slate-600">{e.activity_id}</td>
                    <td className="px-6 py-4 font-mono text-slate-600">{e.student_id}</td>
                    <td className="px-6 py-4">
                      {e.score} / {e.max_score}
                    </td>
                    <td className="px-6 py-4 font-bold">{e.grade}</td>
                  </tr>
                ))}
                {gradebook.length === 0 && (
                  <tr>
                    <td colSpan={4} className="p-12 text-center text-slate-400 italic">
                      {t('teacher.detail.gradebook.empty')}
                    </td>
                  </tr>
                )}
              </tbody>
            </table>
          )}
        </div>
      )}
    </div>
  );
};

