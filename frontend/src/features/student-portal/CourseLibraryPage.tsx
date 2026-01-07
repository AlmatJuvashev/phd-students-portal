import React, { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useQuery } from '@tanstack/react-query';
import { Search, Loader2, BookOpen, GraduationCap } from 'lucide-react';
import { Input } from '@/components/ui/input';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { getStudentCatalog } from './api';
import { StudentCourse } from './types';
import { toast } from 'sonner';

export const CourseLibraryPage: React.FC = () => {
  const { t } = useTranslation('common');
  const [search, setSearch] = useState('');

  const catalogQuery = useQuery({
    queryKey: ['student', 'catalog'],
    queryFn: getStudentCatalog,
  });

  const courses = catalogQuery.data || [];
  const filtered = courses.filter(c => 
    c.title.toLowerCase().includes(search.toLowerCase()) || 
    c.code.toLowerCase().includes(search.toLowerCase())
  );

  const handleEnroll = (course: StudentCourse) => {
    // Phase 17: Implement enrollment request logic
    toast.info(`Enrollment request for ${course.code} sent! (Simulated)`);
  };

  if (catalogQuery.isLoading) {
    return (
      <div className="h-full flex items-center justify-center">
        <Loader2 className="animate-spin text-indigo-600" />
      </div>
    );
  }

  return (
    <div className="max-w-6xl mx-auto space-y-6 animate-in fade-in duration-500">
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-black text-slate-900 tracking-tight flex items-center gap-2">
            <BookOpen className="h-6 w-6 text-indigo-600" />
            {t('student.library.title', { defaultValue: 'Course Library' })}
          </h1>
          <p className="text-slate-500 text-sm mt-1">
            {t('student.library.subtitle', { defaultValue: 'Browse and enroll in available courses.' })}
          </p>
        </div>
        
        <div className="relative w-full md:w-72">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-slate-400" />
          <Input 
            placeholder={t('student.library.search', { defaultValue: 'Search courses...' })}
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            className="pl-9 bg-white border-slate-200"
          />
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {filtered.map((c) => (
          <div 
            key={c.code} // Using Code as key if ID is missing or shared in mock
            className="bg-white rounded-3xl border border-slate-200 shadow-sm p-6 flex flex-col gap-4 hover:shadow-md transition-shadow"
          >
            <div className="flex items-start justify-between gap-2">
              <div className="flex items-center gap-3 min-w-0">
                <div className="w-10 h-10 rounded-2xl bg-indigo-50 flex items-center justify-center text-indigo-600 shrink-0">
                   <GraduationCap size={20} />
                </div>
                <div className="min-w-0">
                  <div className="font-bold text-slate-900 truncate" title={c.title}>{c.title}</div>
                  <div className="text-xs text-slate-500 font-mono">{c.code}</div>
                </div>
              </div>
              <Badge variant="outline" className="shrink-0">{c.delivery_format || 'Standard'}</Badge>
            </div>

            <div className="text-sm text-slate-600 line-clamp-3 flex-1">
               {/* Description would go here if available in model */}
               There is currently no description available for this course.
            </div>

            <div className="pt-4 border-t border-slate-100 flex items-center justify-between mt-auto">
               <div className="text-xs text-slate-400">
                  {c.term_id || 'Fall 2025'}
               </div>
               <Button size="sm" onClick={() => handleEnroll(c)}>
                 {t('student.library.enroll', { defaultValue: 'Enroll' })}
               </Button>
            </div>
          </div>
        ))}

        {filtered.length === 0 && (
          <div className="col-span-full text-center py-12 text-slate-400 border-2 border-dashed border-slate-200 rounded-3xl">
            {t('student.library.empty', { defaultValue: 'No courses found.' })}
          </div>
        )}
      </div>
    </div>
  );
};
