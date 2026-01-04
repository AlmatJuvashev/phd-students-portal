import React from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import { useTranslation } from 'react-i18next';
import { ArrowLeft, Loader2, Sparkles } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { getCourses, getProgram } from './api';

export const ProgramDetailPage: React.FC = () => {
  const { t } = useTranslation('common');
  const navigate = useNavigate();
  const { id } = useParams<{ id: string }>();

  const programQuery = useQuery({
    queryKey: ['curriculum', 'programs', id],
    queryFn: () => getProgram(id!),
    enabled: !!id,
  });

  const coursesQuery = useQuery({
    queryKey: ['curriculum', 'courses', { programId: id }],
    queryFn: () => getCourses(id),
    enabled: !!id,
  });

  if (programQuery.isLoading || coursesQuery.isLoading) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <Loader2 className="animate-spin" />
      </div>
    );
  }

  if (programQuery.isError || !programQuery.data) {
    return (
      <div className="p-8 text-center text-red-700 bg-red-50 rounded-xl border border-red-100">
        {t('curriculum.programs.load_error', { defaultValue: 'Failed to load program.' })}
      </div>
    );
  }

  const program = programQuery.data;
  const courses = coursesQuery.data || [];

  return (
    <div className="space-y-8 animate-in fade-in duration-500">
      <div className="flex items-start justify-between gap-4">
        <div className="flex items-start gap-3">
          <Button variant="ghost" size="icon" onClick={() => navigate('/admin/programs')}>
            <ArrowLeft size={16} />
          </Button>
          <div className="min-w-0">
            <div className="flex items-center gap-2 flex-wrap">
              <h1 className="text-2xl font-black text-slate-900 tracking-tight truncate">{program.title}</h1>
              <Badge variant="secondary" className="font-mono">
                {program.code}
              </Badge>
              <Badge variant="outline" className="capitalize">
                {program.status}
              </Badge>
            </div>
            <p className="text-slate-500 text-sm mt-1">{program.description || '—'}</p>
          </div>
        </div>

        <Button
          onClick={() => navigate(`/admin/studio/programs/${program.id}/builder`)}
          className="bg-indigo-600 hover:bg-indigo-700 text-white shadow-lg shadow-indigo-200"
        >
          <Sparkles className="mr-2 h-4 w-4" />
          {t('curriculum.programs.edit_in_builder', { defaultValue: 'Edit in Builder' })}
        </Button>
      </div>

      <div className="bg-white border border-slate-200 rounded-3xl shadow-sm overflow-hidden">
        <div className="p-6 border-b border-slate-100">
          <div className="text-xs font-bold text-slate-400 uppercase tracking-widest">
            {t('curriculum.programs.courses', { defaultValue: 'Courses in Program' })}
          </div>
          <div className="mt-1 text-sm text-slate-500">
            {t('curriculum.programs.courses_count', { defaultValue: '{{count}} courses', count: courses.length })}
          </div>
        </div>

        <div className="divide-y divide-slate-100">
          {courses.map((c) => (
            <div key={c.id} className="p-5 flex items-start justify-between gap-4">
              <div className="min-w-0">
                <div className="font-bold text-slate-900 truncate">{c.title}</div>
                <div className="text-xs text-slate-500 mt-1">
                  <span className="font-mono">{c.code}</span> · {c.credits} {t('curriculum.courses.credits', { defaultValue: 'credits' })}
                </div>
              </div>
              <Badge variant="secondary" className="shrink-0">
                {c.category}
              </Badge>
            </div>
          ))}

          {courses.length === 0 && (
            <div className="p-10 text-center text-slate-400 italic">
              {t('curriculum.programs.no_courses', { defaultValue: 'No courses found for this program.' })}
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

