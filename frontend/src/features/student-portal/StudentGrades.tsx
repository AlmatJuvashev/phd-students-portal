import React from 'react';
import { useTranslation } from 'react-i18next';
import { useQuery } from '@tanstack/react-query';
import { Loader2 } from 'lucide-react';
import { Badge } from '@/components/ui/badge';
import { getStudentGrades } from './api';

export const StudentGrades: React.FC = () => {
  const { t } = useTranslation('common');

  const gradesQuery = useQuery({
    queryKey: ['student', 'grades'],
    queryFn: getStudentGrades,
  });

  const grades = gradesQuery.data || [];

  if (gradesQuery.isLoading) {
    return (
      <div className="h-full flex items-center justify-center">
        <Loader2 className="animate-spin" />
      </div>
    );
  }

  if (gradesQuery.isError) {
    return (
      <div className="max-w-6xl mx-auto p-6">
        <div className="bg-red-50 border border-red-100 rounded-2xl p-4 text-red-900 text-sm">
          {t('student.grades.load_error', { defaultValue: 'Failed to load grades.' })}
        </div>
      </div>
    );
  }

  return (
    <div className="max-w-6xl mx-auto space-y-6 animate-in fade-in duration-500">
      <div>
        <h1 className="text-2xl font-black text-slate-900 tracking-tight">
          {t('student.grades.title', { defaultValue: 'Grades' })}
        </h1>
        <p className="text-slate-500 text-sm mt-1">
          {t('student.grades.subtitle', { defaultValue: 'Your recent graded activities.' })}
        </p>
      </div>

      <div className="bg-white border border-slate-200 rounded-3xl overflow-hidden shadow-sm overflow-x-auto">
        <table className="w-full text-sm text-left whitespace-nowrap">
          <thead className="bg-slate-50 border-b border-slate-200 text-xs font-bold text-slate-500 uppercase">
            <tr>
              <th className="px-6 py-4">Course</th>
              <th className="px-6 py-4">Activity</th>
              <th className="px-6 py-4">Score</th>
              <th className="px-6 py-4">Grade</th>
              <th className="px-6 py-4">Date</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-slate-100">
            {grades.map((g) => (
              <tr key={g.id} className="hover:bg-slate-50 transition-colors">
                <td className="px-6 py-4">
                  <div className="font-bold text-slate-900">{g.course_title || g.course_code || g.course_id || '—'}</div>
                  {g.course_code && (
                    <div className="text-xs text-slate-500 font-mono">{g.course_code}</div>
                  )}
                </td>
                <td className="px-6 py-4 font-mono text-xs text-slate-700">{g.activity_id}</td>
                <td className="px-6 py-4">
                  <span className="font-bold text-slate-900">
                    {g.score}/{g.max_score}
                  </span>
                </td>
                <td className="px-6 py-4">
                  <Badge variant="secondary">{g.grade || '—'}</Badge>
                </td>
                <td className="px-6 py-4 text-slate-500">
                  {g.graded_at ? new Date(g.graded_at).toLocaleDateString() : '—'}
                </td>
              </tr>
            ))}
            {grades.length === 0 && (
              <tr>
                <td colSpan={5} className="px-6 py-12 text-center text-slate-400 italic">
                  {t('student.grades.empty', { defaultValue: 'No grades yet.' })}
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
};

