import React, { useMemo, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { useQuery } from '@tanstack/react-query';
import { ExternalLink, Loader2, Search } from 'lucide-react';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { getCourses } from '@/features/curriculum/api';
import { getTeacherCourses } from './api';

export const TeacherCoursesPage: React.FC = () => {
  const { t } = useTranslation('common');
  const navigate = useNavigate();
  const [search, setSearch] = useState('');

  const offeringsQuery = useQuery({ queryKey: ['teacher', 'courses'], queryFn: getTeacherCourses });
  const catalogQuery = useQuery({
    queryKey: ['curriculum', 'courses'],
    queryFn: () => getCourses(),
    staleTime: 5 * 60 * 1000,
    retry: false,
  });

  const offerings = offeringsQuery.data || [];
  const courseById = useMemo(() => new Map((catalogQuery.data || []).map((c) => [c.id, c])), [catalogQuery.data]);

  const filtered = offerings.filter((o) => {
    const course = courseById.get(o.course_id);
    const hay = `${course?.title || ''} ${course?.code || ''} ${o.section} ${o.delivery_format}`.toLowerCase();
    return hay.includes(search.toLowerCase());
  });

  if (offeringsQuery.isLoading) {
    return (
      <div className="h-full flex items-center justify-center">
        <Loader2 className="animate-spin" />
      </div>
    );
  }

  return (
    <div className="max-w-6xl mx-auto space-y-8 animate-in fade-in duration-500">
      <div className="flex flex-col md:flex-row justify-between items-start md:items-center gap-4">
        <div>
          <h1 className="text-3xl font-black text-slate-900 tracking-tight">{t('teacher.courses.title')}</h1>
          <p className="text-slate-500 font-medium mt-1">{t('teacher.courses.subtitle')}</p>
        </div>
        <div className="flex gap-2 w-full md:w-auto">
          <div className="relative flex-1 md:flex-none">
            <Search size={16} className="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" />
            <Input
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              placeholder={t('teacher.courses.search_placeholder')}
              className="w-full md:w-80 pl-9"
            />
          </div>
          <Button variant="secondary" onClick={() => navigate('/admin/scheduler')}>
            <ExternalLink className="mr-2 h-4 w-4" />
            {t('teacher.courses.open_scheduler')}
          </Button>
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {filtered.map((offering) => {
          const course = courseById.get(offering.course_id);
          return (
            <button
              key={offering.id}
              onClick={() => navigate(`/admin/teacher/courses/${offering.id}`)}
              className="group bg-white rounded-3xl border border-slate-200 shadow-sm hover:shadow-xl hover:border-indigo-200 transition-all overflow-hidden flex flex-col h-full relative cursor-pointer text-left"
            >
              <div className="h-24 bg-slate-100 relative overflow-hidden">
                <div className="absolute inset-0 bg-gradient-to-br from-indigo-500/10 via-emerald-500/10 to-purple-500/10" />
                <div className="absolute bottom-4 left-6 right-6 flex items-center justify-between">
                  <Badge variant="secondary" className="bg-white/80">
                    {course?.code || 'COURSE'}
                  </Badge>
                  <Badge variant="outline" className="bg-white/80">
                    {offering.section}
                  </Badge>
                </div>
              </div>

              <div className="p-6 space-y-3 flex-1 flex flex-col">
                <div>
                  <div className="font-black text-slate-900 leading-tight">{course?.title || offering.course_id}</div>
                  <div className="text-xs text-slate-500 mt-1">{offering.delivery_format}</div>
                </div>
                <div className="mt-auto flex items-center justify-between text-xs text-slate-500">
                  <span>{t('teacher.courses.enrolled', { count: offering.current_enrolled || 0 })}</span>
                  <span className="text-indigo-600 font-bold group-hover:translate-x-1 transition-transform">
                    {t('teacher.courses.open')}
                  </span>
                </div>
              </div>
            </button>
          );
        })}

        {filtered.length === 0 && (
          <div className="col-span-full text-center py-12 text-slate-400 border-2 border-dashed border-slate-200 rounded-3xl">
            {t('teacher.courses.empty')}
          </div>
        )}
      </div>
    </div>
  );
};

