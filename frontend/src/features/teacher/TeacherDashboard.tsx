import React from 'react';
import { useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { useQuery } from '@tanstack/react-query';
import { ArrowRight, Calendar, CheckCircle2, ChevronRight, Loader2, Play, TrendingUp, AlertTriangle } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { cn } from '@/lib/utils';
import { motion } from 'framer-motion';
import { getCourses } from '@/features/curriculum/api';
import { getTeacherCourses, getTeacherDashboard, getTeacherSubmissions } from './api';

const StatWidget = ({ title, value, sub, icon: Icon, color, onClick, index }: any) => (
  <motion.div
    initial={{ opacity: 0, y: 20 }}
    animate={{ opacity: 1, y: 0 }}
    transition={{ delay: index * 0.1 }}
    onClick={onClick}
    className="bg-white p-5 rounded-2xl border border-slate-200 shadow-sm hover:shadow-lg transition-all cursor-pointer group hover:-translate-y-1"
  >
    <div className="flex justify-between items-start mb-2">
      <div className={cn('p-2.5 rounded-xl transition-transform group-hover:scale-110', color)}>
        <Icon size={20} />
      </div>
      <div className="opacity-0 group-hover:opacity-100 text-slate-300 transition-opacity">
        <ArrowRight size={16} />
      </div>
    </div>
    <div className="text-2xl font-black text-slate-900">{value}</div>
    <div className="text-xs font-bold text-slate-500 uppercase tracking-wide mt-1">{title}</div>
    {sub && <div className="text-[10px] text-slate-400 mt-1 font-medium">{sub}</div>}
  </motion.div>
);

const ScheduleItemRow = ({ time, title, type, location }: any) => (
  <div className="flex gap-4 p-3 rounded-xl hover:bg-slate-50 transition-colors group">
    <div className="w-16 flex-shrink-0 text-right">
      <div className="text-sm font-black text-slate-900">{time}</div>
      <div className="text-[10px] text-slate-400 font-bold uppercase">{type}</div>
    </div>
    <div className="w-1 bg-slate-200 rounded-full relative group-hover:bg-indigo-400 transition-colors">
      <div className="absolute top-2 -left-[3px] w-2.5 h-2.5 bg-white border-2 border-slate-300 rounded-full group-hover:border-indigo-500 transition-colors" />
    </div>
    <div className="flex-1 pb-4">
      <div className="font-bold text-slate-800 text-sm">{title}</div>
      <div className="text-xs text-slate-500 flex items-center gap-1 mt-0.5">
        <Calendar size={10} /> {location}
      </div>
    </div>
  </div>
);

export const TeacherDashboard: React.FC = () => {
  const { t } = useTranslation('common');
  const navigate = useNavigate();

  const statsQuery = useQuery({ queryKey: ['teacher', 'dashboard'], queryFn: getTeacherDashboard });
  const coursesQuery = useQuery({ queryKey: ['teacher', 'courses'], queryFn: getTeacherCourses });
  const submissionsQuery = useQuery({ queryKey: ['teacher', 'submissions'], queryFn: getTeacherSubmissions });

  const catalogQuery = useQuery({
    queryKey: ['curriculum', 'courses'],
    queryFn: () => getCourses(),
    staleTime: 5 * 60 * 1000,
    retry: false,
  });

  const stats = statsQuery.data;
  const courses = coursesQuery.data || [];
  const submissions = submissionsQuery.data || [];
  const courseById = new Map((catalogQuery.data || []).map((c) => [c.id, c]));

  if (statsQuery.isLoading || coursesQuery.isLoading) {
    return (
      <div className="h-full flex items-center justify-center">
        <Loader2 className="animate-spin" />
      </div>
    );
  }

  const nextClass = stats?.next_class;

  return (
    <div className="space-y-8 animate-in fade-in duration-500">
      <div className="flex justify-between items-end">
        <div>
          <h1 className="text-3xl font-black text-slate-900 tracking-tight">{t('teacher.dashboard.title')}</h1>
          <p className="text-slate-500 font-medium mt-1">{t('teacher.dashboard.subtitle')}</p>
        </div>
        <Button variant="secondary" onClick={() => navigate('/admin/scheduler')}>
          <Calendar className="mr-2 h-4 w-4" />
          {t('teacher.dashboard.view_full_schedule')}
        </Button>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-5 gap-4">
        <StatWidget
          index={0}
          title={t('teacher.dashboard.stats.to_grade')}
          value={stats?.pending_grading || 0}
          sub={t('teacher.dashboard.stats.pending_grading_subtitle')}
          icon={CheckCircle2}
          color="bg-indigo-50 text-indigo-600"
          onClick={() => navigate('/admin/teacher/grading')}
        />
        <StatWidget
          index={1}
          title={t('teacher.dashboard.stats.active_courses')}
          value={stats?.active_courses || 0}
          sub={t('teacher.dashboard.stats.active_courses_subtitle')}
          icon={Play}
          color="bg-emerald-50 text-emerald-600"
          onClick={() => navigate('/admin/teacher/courses')}
        />
        <StatWidget
          index={2}
          title={t('teacher.dashboard.stats.today_classes')}
          value={stats?.today_classes_count || 0}
          sub={t('teacher.dashboard.stats.today_classes_subtitle')}
          icon={Calendar}
          color="bg-amber-50 text-amber-600"
          onClick={() => navigate('/admin/scheduler')}
        />
        <StatWidget
          index={3}
          title={t('teacher.dashboard.stats.next_class')}
          value={nextClass ? nextClass.start_time : '—'}
          sub={nextClass ? nextClass.title : t('teacher.dashboard.stats.next_class_none')}
          icon={Play}
          color="bg-slate-50 text-slate-700"
          onClick={() => navigate('/admin/scheduler')}
        />
        <StatWidget
          index={4}
          title={t('teacher.dashboard.stats.students_at_risk') || 'At Risk'}
          value={stats?.at_risk_count || 0}
          sub={t('teacher.dashboard.stats.at_risk_subtitle') || 'Requires attention'}
          icon={AlertTriangle}
          color="bg-red-50 text-red-600"
          onClick={() => navigate('/admin/teacher/students')}
        />
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
        <div className="lg:col-span-2 space-y-8">
          <div className="bg-white rounded-3xl border border-slate-200 shadow-sm p-6 relative overflow-hidden">
            <div className="flex justify-between items-center mb-6 relative z-10">
              <h3 className="font-bold text-slate-900 text-lg">{t('teacher.dashboard.schedule.title')}</h3>
              <div className="text-xs font-bold text-slate-400 uppercase tracking-widest bg-slate-50 px-3 py-1 rounded-lg">
                {new Date().toLocaleDateString()}
              </div>
            </div>

            {nextClass ? (
              <div className="bg-slate-900 rounded-2xl p-5 text-white shadow-xl mb-6 relative overflow-hidden">
                <div className="absolute top-0 right-0 w-40 h-40 bg-indigo-500 rounded-full blur-2xl opacity-30 -mr-10 -mt-10" />
                <div className="relative z-10">
                  <div className="text-xs font-bold text-indigo-300 uppercase tracking-widest mb-2">Next Class</div>
                  <div className="font-black text-xl mb-1">{nextClass.title}</div>
                  <div className="text-sm text-slate-400">
                    {nextClass.start_time} - {nextClass.end_time}
                  </div>
                </div>
                {nextClass.meeting_url && (
                  <a
                    className="absolute bottom-4 right-4 bg-white/10 hover:bg-white/20 backdrop-blur-md px-3 py-2 rounded-xl text-xs font-bold transition-colors flex items-center gap-2"
                    href={nextClass.meeting_url}
                    target="_blank"
                    rel="noreferrer"
                  >
                    Join Live <Play size={12} fill="currentColor" />
                  </a>
                )}
              </div>
            ) : (
              <div className="bg-slate-50 rounded-2xl p-5 text-slate-500 mb-6 border border-slate-200">
                {t('teacher.dashboard.schedule.no_next_class')}
              </div>
            )}

            <div className="space-y-1">
              {nextClass ? (
                <ScheduleItemRow
                  key={nextClass.id}
                  time={nextClass.start_time}
                  type={nextClass.type}
                  title={nextClass.title}
                  location={nextClass.room_id || 'TBA'}
                />
              ) : (
                <div className="text-center py-8 text-slate-400 text-sm italic">{t('teacher.dashboard.schedule.empty')}</div>
              )}
            </div>
          </div>

          <div className="bg-white rounded-3xl border border-slate-200 shadow-sm overflow-hidden">
            <div className="p-6 border-b border-slate-100 flex justify-between items-center">
              <h3 className="font-bold text-slate-900 text-lg">{t('teacher.dashboard.submissions.title')}</h3>
              <Button size="sm" variant="ghost" onClick={() => navigate('/admin/teacher/grading')}>
                {t('teacher.dashboard.submissions.view_all')}
              </Button>
            </div>
            <div className="divide-y divide-slate-50">
              {submissions.slice(0, 5).map((sub) => (
                <div
                  key={sub.id}
                  className="p-4 hover:bg-slate-50 transition-colors flex items-center justify-between group cursor-pointer"
                  onClick={() => navigate('/admin/teacher/grading')}
                >
                  <div className="flex items-center gap-4">
                    <div className="w-10 h-10 rounded-full bg-slate-100 flex items-center justify-center text-xs font-bold text-slate-500 border border-slate-200">
                      {(sub.student_name || sub.student_id || '?').toString().slice(0, 1)}
                    </div>
                    <div>
                      <div className="font-bold text-slate-900 text-sm">{sub.student_name || sub.student_id}</div>
                      <div className="text-xs text-slate-500">{sub.activity_title || sub.activity_id}</div>
                    </div>
                  </div>
                  <div className="flex items-center gap-4">
                    <span className="text-xs font-medium text-slate-400">{new Date(sub.submitted_at).toLocaleString()}</span>
                    <Badge variant={sub.status === 'SUBMITTED' ? 'secondary' : 'outline'} className="text-[10px] uppercase">
                      {sub.status}
                    </Badge>
                    <ChevronRight size={16} className="text-slate-300 group-hover:text-indigo-600" />
                  </div>
                </div>
              ))}
              {submissions.length === 0 && (
                <div className="p-8 text-center text-slate-400 text-sm italic">No recent submissions.</div>
              )}
            </div>
          </div>
        </div>

        <div className="space-y-6">
          <div className="bg-white rounded-3xl border border-slate-200 shadow-sm p-6">
            <h3 className="font-bold text-slate-900 text-lg mb-4">{t('teacher.dashboard.active_classes.title')}</h3>
            <div className="space-y-3">
              {courses.slice(0, 5).map((offering) => {
                const course = courseById.get(offering.course_id);
                const title = course?.title || offering.course_id;
                const code = course?.code || offering.section;
                return (
                  <div
                    key={offering.id}
                    className="p-4 rounded-2xl border border-slate-100 hover:border-slate-300 transition-all cursor-pointer group"
                    onClick={() => navigate(`/admin/teacher/courses/${offering.id}`)}
                  >
                    <div className="flex justify-between items-start mb-2">
                      <div>
                        <div className="font-bold text-slate-800 text-sm">{title}</div>
                        <div className="text-[10px] font-bold text-slate-400 uppercase tracking-wider">{code}</div>
                      </div>
                      <div className="text-xs font-bold bg-slate-100 px-2 py-1 rounded text-slate-600">
                        {t('teacher.dashboard.active_classes.students_count', { count: offering.current_enrolled || 0 })}
                      </div>
                    </div>
                    <div className="text-xs text-slate-500">{offering.delivery_format}</div>
                  </div>
                );
              })}
              {courses.length === 0 && <div className="text-center py-4 text-slate-400 text-xs italic">No active classes found.</div>}
            </div>
          </div>

          <div className="bg-gradient-to-br from-slate-900 to-slate-800 rounded-3xl p-6 text-white shadow-xl relative overflow-hidden">
            <div className="absolute top-0 right-0 p-4 opacity-10">
              <TrendingUp size={80} />
            </div>
            <div className="relative z-10">
              <div className="text-xs font-bold text-indigo-300 uppercase tracking-widest mb-2">{t('teacher.dashboard.insight.title')}</div>
              <h4 className="font-bold text-lg leading-tight mb-4">{t('teacher.dashboard.insight.placeholder_title')}</h4>
              <p className="text-sm text-slate-400 mb-6">{t('teacher.dashboard.insight.placeholder_body')}</p>
              <Button
                size="sm"
                className="w-full bg-white text-slate-900 hover:bg-slate-100 border-none transition-colors"
                onClick={() => navigate('/admin/analytics')}
              >
                {t('teacher.dashboard.insight.review_analytics')}
              </Button>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

