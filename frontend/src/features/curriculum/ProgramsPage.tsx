import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { useQuery } from '@tanstack/react-query';
import { Search, GraduationCap, Users, Activity, Loader2, Sparkles } from 'lucide-react';
import { Badge } from '@/components/ui/badge';
import { Input } from '@/components/ui/input';
import { cn } from '@/lib/utils';
import { getPrograms } from './api';
import { Program } from './types';

import { EduStudioHub } from '../studio/components/EduStudioHub';

export const ProgramsPage: React.FC = () => {
  const { t } = useTranslation('common');
  const navigate = useNavigate();
  const [filterStatus, setFilterStatus] = useState('all');
  const [search, setSearch] = useState('');
  const [showStudio, setShowStudio] = useState(false);

  const { data: programs = [], isLoading, error } = useQuery({
    queryKey: ['curriculum', 'programs'],
    queryFn: getPrograms,
  });

  const filteredPrograms = programs.filter((p: Program) => {
    const matchesSearch = p.title.toLowerCase().includes(search.toLowerCase()) || 
                         p.code.toLowerCase().includes(search.toLowerCase());
    const matchesStatus = filterStatus === 'all' || p.status === filterStatus;
    return matchesSearch && matchesStatus;
  });

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
        {t('common.error')}: {error instanceof Error ? error.message : 'Failed to load programs'}
      </div>
    );
  }

  return (
    <div className="space-y-8 animate-in fade-in duration-500 relative">
       <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
          <div>
             <h1 className="text-2xl font-black text-slate-900 tracking-tight">
               {t('curriculum.programs.title')}
             </h1>
             <p className="text-slate-500 text-sm mt-1">
               {t('curriculum.programs.subtitle')}
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

       {/* Filters */}
       <div className="bg-white p-2 rounded-xl border border-slate-200 shadow-sm flex flex-col sm:flex-row gap-2">
          <div className="relative flex-1">
             <Search size={16} className="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" />
             <Input 
               value={search}
               onChange={(e) => setSearch(e.target.value)}
               placeholder={t('curriculum.programs.search_placeholder')}
               className="w-full h-10 pl-9 pr-4 bg-transparent border-none focus-visible:ring-0 text-sm shadow-none"
             />
          </div>
          <div className="h-6 w-px bg-slate-200 hidden sm:block self-center" />
          <div className="flex gap-1 p-1 bg-slate-100 rounded-lg">
             {['all', 'active', 'draft', 'archived'].map(s => (
               <button
                 key={s}
                 onClick={() => setFilterStatus(s)}
                 className={cn(
                   "px-3 py-1.5 text-xs font-bold rounded-md capitalize transition-all",
                   filterStatus === s ? "bg-white text-slate-900 shadow-sm" : "text-slate-500 hover:text-slate-700"
                 )}
               >
                 {t(`curriculum.programs.status.${s}`)}
               </button>
             ))}
          </div>
       </div>

       {/* Grid */}
       <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {filteredPrograms.map((program: Program) => (
             <div 
               key={program.id}
               onClick={() => navigate(`/admin/studio/programs/${program.id}/builder`)}
               className="group bg-white p-6 rounded-2xl border border-slate-200 shadow-sm hover:shadow-lg hover:border-emerald-300 transition-all cursor-pointer flex flex-col relative overflow-hidden"
             >
                {/* Status Stripe */}
                <div className={cn(
                  "absolute top-0 left-0 w-full h-1",
                  program.status === 'active' ? "bg-emerald-500" : program.status === 'draft' ? "bg-amber-500" : "bg-slate-300"
                )} />

                <div className="flex justify-between items-start mb-4">
                   <div className="w-12 h-12 rounded-xl bg-slate-50 flex items-center justify-center text-slate-400 group-hover:bg-emerald-50 group-hover:text-emerald-600 transition-colors">
                      <GraduationCap size={24} />
                   </div>
                   <Badge variant={program.status === 'active' ? 'default' : program.status === 'draft' ? 'secondary' : 'outline'} className={cn(
                      program.status === 'active' && "bg-emerald-100 text-emerald-700 hover:bg-emerald-100",
                      program.status === 'draft' && "bg-amber-100 text-amber-700 hover:bg-amber-100"
                   )}>
                      {t(`curriculum.programs.status.${program.status}`)}
                   </Badge>
                </div>

                <h3 className="text-lg font-bold text-slate-900 mb-2 group-hover:text-emerald-700 transition-colors">{program.title}</h3>
                <p className="text-xs text-slate-500 mb-6 line-clamp-2 flex-1">{program.description}</p>

                <div className="grid grid-cols-2 gap-4 pt-4 border-t border-slate-100">
                   <div>
                      <div className="text-[10px] font-bold text-slate-400 uppercase tracking-wide mb-0.5">
                        {t('curriculum.programs.stats.students')}
                      </div>
                      <div className="text-sm font-bold text-slate-900 flex items-center gap-1">
                         <Users size={14} className="text-indigo-500" /> {0} {/* Backend needs to provide this count */}
                      </div>
                   </div>
                   <div>
                      <div className="text-[10px] font-bold text-slate-400 uppercase tracking-wide mb-0.5">
                        {t('curriculum.programs.stats.completion')}
                      </div>
                      <div className="text-sm font-bold text-slate-900 flex items-center gap-1">
                         <Activity size={14} className="text-emerald-500" /> {0}% {/* Backend calculation needed */}
                      </div>
                   </div>
                </div>
             </div>
          ))}
          
          {filteredPrograms.length === 0 && (
             <div className="col-span-full py-12 text-center text-slate-400 border-2 border-dashed border-slate-200 rounded-2xl">
               {t('curriculum.programs.no_programs')}
             </div>
          )}
       </div>
    </div>
  );
};
