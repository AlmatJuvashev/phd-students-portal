import React from 'react';
import { useTranslation } from 'react-i18next';
import { useQuery } from '@tanstack/react-query';
import { CalendarClock, GraduationCap, Loader2 } from 'lucide-react';
import { Badge } from '@/components/ui/badge';
import { getStudentCourses } from './api';

export const StudentCourses: React.FC = () => {
  const { t } = useTranslation('common');

  const coursesQuery = useQuery({
    queryKey: ['student', 'courses'],
    queryFn: getStudentCourses,
  });

  const courses = coursesQuery.data || [];

  if (coursesQuery.isLoading) {
    return (
      <div className="h-full flex items-center justify-center">
        <Loader2 className="animate-spin" />
      </div>
    );
  }

  if (coursesQuery.isError) {
    return (
      <div className="max-w-6xl mx-auto p-6">
        <div className="bg-red-50 border border-red-100 rounded-2xl p-4 text-red-900 text-sm">
          {t('student.courses.load_error', { defaultValue: 'Failed to load courses.' })}
        </div>
      </div>
    );
  }

  return (
    <div className="max-w-6xl mx-auto space-y-6 animate-in fade-in duration-500">
      <div>
        <h1 className="text-2xl font-black text-slate-900 tracking-tight">
          {t('student.courses.title', { defaultValue: 'My Courses' })}
        </h1>
        <p className="text-slate-500 text-sm mt-1">
          {t('student.courses.subtitle', { defaultValue: 'Current enrollments and upcoming sessions.' })}
        </p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {courses.map((c) => (
          <div
            key={c.course_offering_id}
            className="bg-white rounded-3xl border border-slate-200 shadow-sm p-6 flex flex-col gap-3"
          >
            <div className="flex items-start justify-between gap-2">
              <div className="flex items-center gap-2 min-w-0">
                <div className="w-9 h-9 rounded-xl bg-slate-50 flex items-center justify-center text-slate-500">
                  <GraduationCap size={18} />
                </div>
                <div className="min-w-0">
                  <div className="font-black text-slate-900 truncate">{c.title}</div>
                  <div className="text-xs text-slate-500 truncate">{c.instructor_name || '—'}</div>
                </div>
              </div>
              <Badge variant="secondary" className="font-mono">
                {c.code}
              </Badge>
            </div>

            <div className="flex items-center gap-2 text-xs text-slate-500">
              <Badge variant="outline" className="text-[10px]">
                {t('student.courses.section', { defaultValue: 'Section' })} {c.section}
              </Badge>
              <Badge variant="outline" className="text-[10px]">
                {c.delivery_format}
              </Badge>
            </div>

            {c.next_session ? (
              <div className="mt-2 flex items-start gap-2 text-sm text-slate-700">
                <CalendarClock size={16} className="mt-0.5 text-slate-400" />
                <div>
                  <div className="font-bold">
                    {t('student.courses.next_session', { defaultValue: 'Next session' })} — {c.next_session.date}
                  </div>
                  <div className="text-xs text-slate-500">
                    {c.next_session.start_time}–{c.next_session.end_time} · {c.next_session.type}
                  </div>
                </div>
              </div>
            ) : (
              <div className="mt-2 text-xs text-slate-500">
                {t('student.courses.no_upcoming', { defaultValue: 'No upcoming sessions scheduled.' })}
              </div>
            )}
          </div>
        ))}

        {courses.length === 0 && (
          <div className="col-span-full text-center py-12 text-slate-400 border-2 border-dashed border-slate-200 rounded-3xl">
            {t('student.courses.empty', { defaultValue: 'No active enrollments.' })}
          </div>
        )}
      </div>
    </div>
  );
};

