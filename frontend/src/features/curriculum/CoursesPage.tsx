import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { useQuery } from '@tanstack/react-query';
import { Search, Clock, ExternalLink, Library, Loader2, Sparkles } from 'lucide-react';
import { Badge } from '@/components/ui/badge';
import { Input } from '@/components/ui/input';
import { getCourses } from './api';
import { Course } from './types';

import { EduStudioHub } from '../studio/components/EduStudioHub';

export const CoursesPage: React.FC = () => {
  const { t } = useTranslation('common');
  const navigate = useNavigate();
  const [search, setSearch] = useState('');
  const [showStudio, setShowStudio] = useState(false);

  const { data: courses = [], isLoading, error } = useQuery({
    queryKey: ['curriculum', 'courses'],
    queryFn: () => getCourses(),
  });

  const filtered = courses.filter((c: Course) => 
    c.title.toLowerCase().includes(search.toLowerCase()) || 
    c.code.toLowerCase().includes(search.toLowerCase())
  );

  if (isLoading) {
    return (
      <div className="flex flex-col items-center justify-center min-h-[400px] space-y-4">
        <Loader2 className="w-8 h-8 animate-spin text-emerald-600" />
        <p className="text-sm text-slate-500">{t('common.loading')}</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className="p-8 text-center text-red-500 bg-red-50 rounded-xl border border-red-100">
        {t('common.error')}: {error instanceof Error ? error.message : 'Failed to load courses'}
      </div>
    );
  }

  return (
    <div className="space-y-8 animate-in fade-in duration-500 relative">
       <div className="flex justify-between items-end">
          <div>
             <h1 className="text-2xl font-black text-slate-900 tracking-tight">
               {t('curriculum.inventory.title')}
             </h1>
             <p className="text-slate-500 text-sm mt-1">
               {t('curriculum.inventory.subtitle')}
             </p>
          </div>
          <button 
             onClick={() => setShowStudio(true)}
             className="flex items-center gap-2 px-4 py-2 bg-indigo-600 text-white rounded-xl shadow-lg hover:bg-indigo-700 hover:scale-105 transition-all font-bold text-sm"
          >
             <div className="p-1 bg-white/20 rounded-md"><Sparkles size={14} /></div>
             Open Studio Hub
          </button>
       </div>

       <EduStudioHub 
         isOpen={showStudio} 
         onClose={() => setShowStudio(false)} 
         onNavigate={navigate} 
       />

       {/* Search Bar */}
       <div className="bg-white p-2 rounded-xl border border-slate-200 shadow-sm flex gap-2 max-w-md">
          <div className="relative flex-1">
             <Search size={16} className="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" />
             <Input 
               value={search}
               onChange={(e) => setSearch(e.target.value)}
               placeholder={t('curriculum.inventory.search_placeholder')}
               className="w-full h-10 pl-9 pr-4 bg-transparent border-none focus-visible:ring-0 text-sm shadow-none"
             />
          </div>
       </div>

       {/* List */}
       <div className="bg-white border border-slate-200 rounded-2xl overflow-hidden shadow-sm overflow-x-auto">
          <table className="w-full text-sm text-left whitespace-nowrap">
             <thead className="bg-slate-50 border-b border-slate-200 text-xs font-bold text-slate-500 uppercase">
                <tr>
                   <th className="px-6 py-4">{t('curriculum.inventory.table.code')}</th>
                   <th className="px-6 py-4">{t('curriculum.inventory.table.title')}</th>
                   <th className="px-6 py-4">{t('curriculum.inventory.table.category')}</th>
                   <th className="px-6 py-4">{t('curriculum.inventory.table.credits')}</th>
                   <th className="px-6 py-4">{t('curriculum.inventory.table.duration')}</th>
                   <th className="px-6 py-4 text-right">{t('curriculum.inventory.table.source')}</th>
                </tr>
             </thead>
             <tbody className="divide-y divide-slate-100">
                {filtered.map((course: Course) => (
                   <tr key={course.id} className="hover:bg-slate-50 transition-colors group">
                      <td className="px-6 py-4 font-mono text-slate-500 font-bold">{course.code}</td>
                      <td className="px-6 py-4 font-bold text-slate-900 flex items-center gap-2">
                         <div className="w-6 h-6 rounded bg-indigo-50 text-indigo-600 flex items-center justify-center">
                            <Library size={12} />
                         </div>
                         {course.title}
                      </td>
                      <td className="px-6 py-4">
                         <Badge variant="secondary" className="uppercase text-[10px]">{course.category}</Badge>
                      </td>
                      <td className="px-6 py-4">{course.credits}</td>
                      <td className="px-6 py-4">
                         <div className="flex items-center gap-2 text-slate-500">
                            <Clock size={14} /> 15 weeks {/* Default duration for now */}
                         </div>
                      </td>
                      <td className="px-6 py-4 text-right">
                         <button 
                           onClick={() => navigate(`/admin/studio/courses/${course.id}/builder`)}
                           className="inline-flex items-center gap-1.5 px-3 py-1.5 bg-slate-100 hover:bg-indigo-50 text-slate-600 hover:text-indigo-600 rounded-lg text-xs font-bold transition-colors"
                         >
                            {t('curriculum.inventory.table.view_studio')} <ExternalLink size={10} />
                         </button>
                      </td>
                   </tr>
                ))}
                {filtered.length === 0 && (
                   <tr>
                      <td colSpan={6} className="px-6 py-12 text-center text-slate-400 italic">
                         {t('curriculum.inventory.no_courses')}
                      </td>
                   </tr>
                )}
             </tbody>
          </table>
       </div>
    </div>
  );
};
