import React from 'react';
import { useTranslation } from 'react-i18next';
import { useNavigate } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import { Calendar, ExternalLink, Loader2 } from 'lucide-react';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { cn } from '@/lib/utils';
import { getStudentAssignments } from './api';

export const StudentAssignments: React.FC = () => {
  const { t } = useTranslation('common');
  const navigate = useNavigate();

  const assignmentsQuery = useQuery({
    queryKey: ['student', 'assignments'],
    queryFn: getStudentAssignments,
  });

  const assignments = assignmentsQuery.data || [];

  if (assignmentsQuery.isLoading) {
    return (
      <div className="h-full flex items-center justify-center">
        <Loader2 className="animate-spin" />
      </div>
    );
  }

  if (assignmentsQuery.isError) {
    return (
      <div className="max-w-6xl mx-auto p-6">
        <div className="bg-red-50 border border-red-100 rounded-2xl p-4 text-red-900 text-sm">
          {t('student.assignments.load_error', { defaultValue: 'Failed to load assignments.' })}
        </div>
      </div>
    );
  }

  return (
    <div className="max-w-6xl mx-auto space-y-6 animate-in fade-in duration-500">
      <div>
        <h1 className="text-2xl font-black text-slate-900 tracking-tight">
          {t('student.assignments.title', { defaultValue: 'Assignments' })}
        </h1>
        <p className="text-slate-500 text-sm mt-1">
          {t('student.assignments.subtitle', { defaultValue: 'Upcoming tasks and deadlines.' })}
        </p>
      </div>

      <div className="bg-white border border-slate-200 rounded-3xl overflow-hidden shadow-sm">
        <div className="divide-y divide-slate-100">
          {assignments.map((a) => (
            <div key={a.id} className="p-5 flex items-start justify-between gap-4">
              <div className="min-w-0">
                <div className="flex items-center gap-2">
                  <div
                    className={cn(
                      'w-2 h-2 rounded-full',
                      a.severity === 'urgent' ? 'bg-red-500' : 'bg-slate-300'
                    )}
                  />
                  <div className="font-bold text-slate-900 truncate">{a.title}</div>
                  <Badge
                    variant="secondary"
                    className={cn(
                      a.severity === 'urgent' ? 'bg-red-100 text-red-700 hover:bg-red-100' : ''
                    )}
                  >
                    {a.severity || 'normal'}
                  </Badge>
                </div>
                <div className="mt-1 text-xs text-slate-500 flex items-center gap-2 flex-wrap">
                  <span className="inline-flex items-center gap-1">
                    <Calendar size={12} />{' '}
                    {a.due_at ? new Date(a.due_at).toLocaleDateString() : t('student.assignments.due_soon', { defaultValue: 'Due soon' })}
                  </span>
                  <span>·</span>
                  <span className="font-mono">{a.source}</span>
                  {a.status && (
                    <>
                      <span>·</span>
                      <span>{a.status}</span>
                    </>
                  )}
                </div>
              </div>

              {a.link ? (
                <Button
                  size="sm"
                  variant="outline"
                  onClick={() => navigate(a.link!)}
                  className="shrink-0"
                >
                  <ExternalLink className="mr-2 h-4 w-4" />
                  {t('student.assignments.open', { defaultValue: 'Open' })}
                </Button>
              ) : (
                <div className="text-xs text-slate-400 shrink-0">
                  {t('student.assignments.no_link', { defaultValue: 'No link' })}
                </div>
              )}
            </div>
          ))}

          {assignments.length === 0 && (
            <div className="p-10 text-center text-slate-400 italic">
              {t('student.assignments.empty', { defaultValue: 'No assignments.' })}
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

