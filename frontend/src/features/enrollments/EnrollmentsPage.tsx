import React, { useState, useMemo } from 'react';
import { useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { 
  Search, Filter, UserPlus, FileDown, MoreHorizontal, 
  AlertCircle, CheckCircle2, PauseCircle, LayoutGrid, 
  List as ListIcon, Mail, Trash2, ArrowRight, Calendar, Users,
  X, Loader2
} from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Input } from '@/components/ui/input';
import { cn } from '@/lib/utils';
import { motion, AnimatePresence } from 'framer-motion';
import { getEnrollments, updateEnrollmentStatus, bulkEnroll } from './api';
import { Enrollment, Student } from './types';
import { getPrograms } from '../curriculum/api';

// --- Helper Components ---

const EnrollmentCard: React.FC<{ enrollment: Enrollment; onSelect: () => void; isSelected: boolean }> = ({ enrollment, onSelect, isSelected }) => {
  const { t } = useTranslation('common');
  
  return (
    <motion.div 
      layoutId={enrollment.id}
      initial={{ opacity: 0, y: 10 }}
      animate={{ opacity: 1, y: 0 }}
      className={cn(
        "bg-white p-4 rounded-xl border shadow-sm cursor-pointer group relative transition-all",
        isSelected ? "border-indigo-500 ring-1 ring-indigo-500 shadow-md" : "border-slate-200 hover:border-indigo-300 hover:shadow-md"
      )}
      onClick={onSelect}
    >
      <div className="flex items-start justify-between mb-3">
         <div className="flex items-center gap-3">
            <div className="w-10 h-10 rounded-full bg-slate-100 flex items-center justify-center text-sm font-bold text-slate-600 border border-slate-200">
               {enrollment.student.avatar || (enrollment.student.name?.[0])}
            </div>
            <div>
               <div className="font-bold text-slate-900 text-sm leading-tight">{enrollment.student.name}</div>
               <div className="text-xs text-slate-500">{enrollment.program.title}</div>
            </div>
         </div>
         {enrollment.overdue_tasks > 0 && (
            <div className="w-2 h-2 rounded-full bg-red-500 animate-pulse" title={`${enrollment.overdue_tasks} overdue items`} />
         )}
      </div>
      
      <div className="space-y-3">
         <div>
            <div className="flex justify-between text-[10px] font-bold text-slate-400 uppercase tracking-wider mb-1">
               <span>{t('enrollments.table.progress')}</span>
               <span className={enrollment.progress > 80 ? "text-emerald-600" : "text-slate-600"}>{enrollment.progress}%</span>
            </div>
            <div className="w-full h-1.5 bg-slate-100 rounded-full overflow-hidden">
               <div 
                 className={cn("h-full rounded-full", enrollment.progress > 80 ? "bg-emerald-500" : enrollment.progress < 30 ? "bg-amber-500" : "bg-indigo-500")} 
                 style={{ width: `${enrollment.progress}%` }} 
               />
            </div>
         </div>
         
         <div className="flex items-center justify-between pt-3 border-t border-slate-50">
            <div className="flex items-center gap-1 text-[10px] font-medium text-slate-400 bg-slate-50 px-2 py-1 rounded-lg">
               <Users size={12} /> {enrollment.cohort_id}
            </div>
            <div className="text-[10px] text-slate-400">
               {enrollment.last_activity ? new Date(enrollment.last_activity).toLocaleDateString() : 'N/A'}
            </div>
         </div>
      </div>

      <div className={cn(
         "absolute top-2 right-2 w-5 h-5 rounded border bg-white flex items-center justify-center transition-all",
         isSelected ? "border-indigo-500 bg-indigo-500 text-white" : "border-slate-200 text-transparent opacity-0 group-hover:opacity-100"
      )}>
         <CheckCircle2 size={14} />
      </div>
    </motion.div>
  );
};

const WizardStep = ({ number, title, active, completed }: any) => (
  <div className={cn("flex items-center gap-2", active ? "text-indigo-600" : completed ? "text-emerald-600" : "text-slate-400")}>
     <div className={cn(
       "w-6 h-6 rounded-full flex items-center justify-center text-xs font-bold border-2 transition-all",
       active ? "border-indigo-600 bg-indigo-50" : completed ? "border-emerald-600 bg-emerald-600 text-white" : "border-slate-200"
     )}>
        {completed ? <CheckCircle2 size={14} /> : number}
     </div>
     <span className="text-xs font-bold uppercase tracking-wider hidden sm:inline">{title}</span>
     {number < 3 && <div className="w-8 h-px bg-slate-200 mx-2 hidden sm:block" />}
  </div>
);

export const EnrollmentsPage: React.FC = () => {
  const { t } = useTranslation('common');
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  
  const [viewMode, setViewMode] = useState<'list' | 'board'>('list');
  const [search, setSearch] = useState('');
  const [selectedIds, setSelectedIds] = useState<Set<string>>(new Set());
  
  // Wizard State
  const [isWizardOpen, setIsWizardOpen] = useState(false);
  const [wizardStep, setWizardStep] = useState(1);
  const [wizardData, setWizardData] = useState<{
    programId: string;
    cohortId: string;
    studentIds: string[];
    startDate: string;
  }>({
    programId: '',
    cohortId: '',
    studentIds: [],
    startDate: new Date().toISOString().split('T')[0]
  });

  // Queries
  const { data: enrollments = [], isLoading, error } = useQuery({
    queryKey: ['enrollments'],
    queryFn: () => getEnrollments(),
  });

  const { data: programs = [] } = useQuery({
    queryKey: ['curriculum', 'programs'],
    queryFn: getPrograms,
    enabled: isWizardOpen
  });

  // Mutations
  const statusMutation = useMutation({
    mutationFn: ({ id, status }: { id: string; status: string }) => updateEnrollmentStatus(id, status),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['enrollments'] });
    }
  });

  const enrollMutation = useMutation({
    mutationFn: bulkEnroll,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['enrollments'] });
      setIsWizardOpen(false);
      setWizardStep(1);
      setWizardData({ programId: '', cohortId: '', studentIds: [], startDate: new Date().toISOString().split('T')[0] });
    }
  });

  // Filter Logic
  const filteredEnrollments = useMemo(() => {
    return enrollments.filter(e => 
      e.student.name.toLowerCase().includes(search.toLowerCase()) || 
      e.program.title.toLowerCase().includes(search.toLowerCase()) ||
      e.cohort_id.toLowerCase().includes(search.toLowerCase())
    );
  }, [enrollments, search]);

  // Kanban Columns
  const columns = useMemo(() => {
     return {
        active: filteredEnrollments.filter(e => e.status === 'active'),
        paused: filteredEnrollments.filter(e => e.status === 'paused'),
        completed: filteredEnrollments.filter(e => e.status === 'completed'),
        dropped: filteredEnrollments.filter(e => e.status === 'dropped'),
     };
  }, [filteredEnrollments]);

  // Bulk Actions
  const toggleSelection = (id: string) => {
    const newSet = new Set(selectedIds);
    if (newSet.has(id)) newSet.delete(id);
    else newSet.add(id);
    setSelectedIds(newSet);
  };

  const selectAll = () => {
    if (selectedIds.size === filteredEnrollments.length && filteredEnrollments.length > 0) setSelectedIds(new Set());
    else setSelectedIds(new Set(filteredEnrollments.map(e => e.id)));
  };

  if (isLoading) {
    return <div className="p-8 flex items-center justify-center h-full"><Loader2 className="animate-spin" /></div>;
  }

  return (
    <div className="space-y-6 animate-in fade-in duration-500 h-[calc(100vh-8rem)] flex flex-col">
       {/* Top Toolbar */}
       <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 flex-shrink-0">
          <div>
             <h1 className="text-2xl font-black text-slate-900 tracking-tight flex items-center gap-2">
               {t('enrollments.title')}
               <span className="text-sm font-medium text-slate-400 bg-slate-100 px-2 py-0.5 rounded-full">{enrollments.length}</span>
             </h1>
             <p className="text-slate-500 text-sm mt-1">{t('enrollments.subtitle')}</p>
          </div>
          
          <div className="flex items-center gap-3">
             <div className="flex bg-slate-100 p-1 rounded-xl">
                <button 
                  onClick={() => setViewMode('list')}
                  className={cn("p-2 rounded-lg transition-all", viewMode === 'list' ? "bg-white shadow text-indigo-600" : "text-slate-400 hover:text-slate-600")}
                  title={t('enrollments.views.list')}
                >
                   <ListIcon size={18} />
                </button>
                <button 
                  onClick={() => setViewMode('board')}
                  className={cn("p-2 rounded-lg transition-all", viewMode === 'board' ? "bg-white shadow text-indigo-600" : "text-slate-400 hover:text-slate-600")}
                  title={t('enrollments.views.board')}
                >
                   <LayoutGrid size={18} />
                </button>
             </div>
             <div className="h-8 w-px bg-slate-200" />
             <Button variant="outline" className="hidden sm:flex items-center gap-2">
                <FileDown size={16} />
                {t('enrollments.export_button')}
             </Button>
             <Button className="flex items-center gap-2" onClick={() => setIsWizardOpen(true)}>
                <UserPlus size={16} />
                {t('enrollments.enroll_button')}
             </Button>
          </div>
       </div>

       {/* Controls & Filters */}
       <div className="bg-white p-2 rounded-xl border border-slate-200 shadow-sm flex flex-col sm:flex-row gap-2 flex-shrink-0">
          <div className="relative flex-1">
             <Search size={16} className="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" />
             <Input 
               value={search}
               onChange={(e) => setSearch(e.target.value)}
               placeholder={t('enrollments.search_placeholder')}
               className="w-full h-10 pl-9 pr-4 bg-transparent border-none focus-visible:ring-0 text-sm shadow-none"
             />
          </div>
          <div className="h-10 w-px bg-slate-200 mx-1 hidden sm:block" />
          <Button variant="ghost" className="flex items-center gap-2 text-slate-500">
             <Filter size={16} />
             {t('common.filter')}
          </Button>
       </div>

       {/* Main Content Area */}
       <div className="flex-1 min-h-0 overflow-hidden relative">
          
          {/* LIST VIEW */}
          {viewMode === 'list' && (
             <div className="bg-white border border-slate-200 rounded-2xl shadow-sm overflow-hidden h-full flex flex-col">
                <div className="overflow-y-auto flex-1 overflow-x-auto">
                   <table className="w-full text-sm text-left">
                      <thead className="bg-slate-50 border-b border-slate-200 text-xs font-bold text-slate-500 uppercase sticky top-0 z-10">
                         <tr>
                            <th className="px-6 py-4 w-10 bg-slate-50">
                               <input 
                                 type="checkbox" 
                                 checked={selectedIds.size === filteredEnrollments.length && filteredEnrollments.length > 0}
                                 onChange={selectAll}
                                 className="rounded border-slate-300 text-indigo-600 focus:ring-indigo-500" 
                               />
                            </th>
                            <th className="px-6 py-4 bg-slate-50">{t('enrollments.table.student')}</th>
                            <th className="px-6 py-4 bg-slate-50">{t('enrollments.table.program_cohort')}</th>
                            <th className="px-6 py-4 bg-slate-50">{t('enrollments.table.progress')}</th>
                            <th className="px-6 py-4 bg-slate-50">{t('enrollments.table.status')}</th>
                            <th className="px-6 py-4 bg-slate-50 text-right">{t('enrollments.table.actions')}</th>
                         </tr>
                      </thead>
                      <tbody className="divide-y divide-slate-100">
                         {filteredEnrollments.map(enr => (
                            <tr 
                              key={enr.id} 
                              className={cn("hover:bg-slate-50 transition-colors group cursor-pointer", selectedIds.has(enr.id) && "bg-indigo-50/30")}
                              onClick={() => toggleSelection(enr.id)}
                            >
                               <td className="px-6 py-4" onClick={(e) => e.stopPropagation()}>
                                  <input 
                                    type="checkbox" 
                                    checked={selectedIds.has(enr.id)}
                                    onChange={() => toggleSelection(enr.id)}
                                    className="rounded border-slate-300 text-indigo-600 focus:ring-indigo-500" 
                                  />
                               </td>
                               <td className="px-6 py-4">
                                  <div className="flex items-center gap-3">
                                     <div className="w-8 h-8 rounded-full bg-white border border-slate-200 flex items-center justify-center text-xs font-bold text-slate-600 uppercase">
                                        {enr.student.avatar || enr.student.name[0]}
                                     </div>
                                     <div>
                                        <div className="font-bold text-slate-900 whitespace-nowrap">{enr.student.name}</div>
                                        <div className="text-xs text-slate-500">{enr.student.email}</div>
                                     </div>
                                  </div>
                               </td>
                               <td className="px-6 py-4">
                                  <div className="font-medium text-slate-800 whitespace-nowrap">{enr.program.title}</div>
                                  <div className="text-xs text-slate-500 bg-slate-100 px-1.5 py-0.5 rounded w-fit mt-1 border border-slate-200">{enr.cohort_id}</div>
                               </td>
                               <td className="px-6 py-4">
                                  <div className="flex items-center gap-2">
                                     <div className="w-20 h-1.5 bg-slate-100 rounded-full overflow-hidden">
                                        <div className={cn("h-full", enr.progress > 80 ? "bg-emerald-500" : "bg-indigo-500")} style={{ width: `${enr.progress}%` }} />
                                      </div>
                                      <span className="text-xs font-bold text-slate-600">{enr.progress}%</span>
                                  </div>
                                  {enr.overdue_tasks > 0 && (
                                     <div className="flex items-center gap-1 text-[10px] text-red-500 font-bold mt-1">
                                        <AlertCircle size={10} /> {enr.overdue_tasks} overdue
                                     </div>
                                  )}
                               </td>
                               <td className="px-6 py-4">
                                  <Badge variant={enr.status === 'active' ? 'default' : enr.status === 'paused' ? 'secondary' : 'outline'}>
                                     {t(`enrollments.status.${enr.status}`)}
                                  </Badge>
                               </td>
                               <td className="px-6 py-4 text-right">
                                  <button className="p-2 text-slate-400 hover:bg-slate-200 hover:text-slate-600 rounded-lg transition-colors opacity-0 group-hover:opacity-100">
                                     <MoreHorizontal size={16} />
                                  </button>
                               </td>
                            </tr>
                         ))}
                      </tbody>
                   </table>
                </div>
             </div>
          )}

          {/* BOARD VIEW (Kanban) */}
          {viewMode === 'board' && (
             <div className="h-full overflow-x-auto overflow-y-hidden pb-4">
                <div className="flex h-full gap-6 min-w-max px-1">
                   {/* Column: Active */}
                   <div className="w-80 flex flex-col h-full bg-slate-100/50 rounded-2xl border border-slate-200/60">
                      <div className="p-4 flex items-center justify-between border-b border-slate-200/60 sticky top-0 bg-slate-100/50 backdrop-blur-sm rounded-t-2xl z-10">
                         <div className="flex items-center gap-2">
                            <div className="w-2 h-2 rounded-full bg-emerald-500" />
                            <h3 className="font-bold text-sm text-slate-700">{t('enrollments.status.active')}</h3>
                         </div>
                         <span className="bg-white px-2 py-0.5 rounded text-xs font-bold text-slate-500 shadow-sm border border-slate-200">{columns.active.length}</span>
                      </div>
                      <div className="p-3 space-y-3 overflow-y-auto flex-1">
                         {columns.active.map(enr => (
                            <EnrollmentCard key={enr.id} enrollment={enr} isSelected={selectedIds.has(enr.id)} onSelect={() => toggleSelection(enr.id)} />
                         ))}
                      </div>
                   </div>

                   {/* Column: Paused */}
                   <div className="w-80 flex flex-col h-full bg-slate-100/50 rounded-2xl border border-slate-200/60">
                      <div className="p-4 flex items-center justify-between border-b border-slate-200/60 sticky top-0 bg-slate-100/50 backdrop-blur-sm rounded-t-2xl z-10">
                         <div className="flex items-center gap-2">
                            <div className="w-2 h-2 rounded-full bg-amber-500" />
                            <h3 className="font-bold text-sm text-slate-700">{t('enrollments.status.paused')}</h3>
                         </div>
                         <span className="bg-white px-2 py-0.5 rounded text-xs font-bold text-slate-500 shadow-sm border border-slate-200">{columns.paused.length}</span>
                      </div>
                      <div className="p-3 space-y-3 overflow-y-auto flex-1">
                         {columns.paused.map(enr => (
                            <EnrollmentCard key={enr.id} enrollment={enr} isSelected={selectedIds.has(enr.id)} onSelect={() => toggleSelection(enr.id)} />
                         ))}
                      </div>
                   </div>

                   {/* Column: Completed */}
                   <div className="w-80 flex flex-col h-full bg-slate-100/50 rounded-2xl border border-slate-200/60">
                      <div className="p-4 flex items-center justify-between border-b border-slate-200/60 sticky top-0 bg-slate-100/50 backdrop-blur-sm rounded-t-2xl z-10">
                         <div className="flex items-center gap-2">
                            <div className="w-2 h-2 rounded-full bg-indigo-500" />
                            <h3 className="font-bold text-sm text-slate-700">{t('enrollments.status.completed')}</h3>
                         </div>
                         <span className="bg-white px-2 py-0.5 rounded text-xs font-bold text-slate-500 shadow-sm border border-slate-200">{columns.completed.length}</span>
                      </div>
                      <div className="p-3 space-y-3 overflow-y-auto flex-1">
                         {columns.completed.map(enr => (
                            <EnrollmentCard key={enr.id} enrollment={enr} isSelected={selectedIds.has(enr.id)} onSelect={() => toggleSelection(enr.id)} />
                         ))}
                      </div>
                   </div>
                </div>
             </div>
          )}

          {/* Floating Action Bar (Bulk Actions) */}
          <AnimatePresence>
             {selectedIds.size > 0 && (
                <motion.div 
                   initial={{ y: 100, opacity: 0 }}
                   animate={{ y: 0, opacity: 1 }}
                   exit={{ y: 100, opacity: 0 }}
                   className="absolute bottom-6 left-1/2 -translate-x-1/2 bg-slate-900 text-white p-2 rounded-2xl shadow-2xl flex items-center gap-4 pl-6 z-50 border border-slate-700"
                >
                   <span className="text-sm font-bold whitespace-nowrap">{selectedIds.size} selected</span>
                   <div className="h-6 w-px bg-slate-700" />
                   <div className="flex gap-1">
                      <button className="p-2 hover:bg-slate-800 rounded-lg transition-colors text-slate-300 hover:text-white" title="Message Students">
                         <Mail size={18} />
                      </button>
                      <button className="p-2 hover:bg-slate-800 rounded-lg transition-colors text-slate-300 hover:text-white" title="Change Status">
                         <PauseCircle size={18} />
                      </button>
                      <button className="p-2 hover:bg-red-900/50 rounded-lg transition-colors text-red-400 hover:text-red-300" title="Remove">
                         <Trash2 size={18} />
                      </button>
                   </div>
                   <button 
                     onClick={() => setSelectedIds(new Set())}
                     className="ml-2 bg-slate-800 hover:bg-slate-700 rounded-full p-1"
                   >
                      <X size={14} />
                   </button>
                </motion.div>
             )}
          </AnimatePresence>
       </div>

       {/* WIZARD MODAL */}
       <AnimatePresence>
         {isWizardOpen && (
           <div className="fixed inset-0 z-[100] flex items-center justify-center p-4 bg-slate-900/60 backdrop-blur-sm">
             <motion.div 
               initial={{ opacity: 0, scale: 0.95 }}
               animate={{ opacity: 1, scale: 1 }}
               exit={{ opacity: 0, scale: 0.95 }}
               className="bg-white w-full max-w-2xl rounded-3xl shadow-2xl overflow-hidden flex flex-col max-h-[90vh]"
             >
                {/* Wizard Header */}
                <div className="p-6 border-b border-slate-100 bg-slate-50 flex justify-between items-center">
                   <div>
                      <h2 className="text-xl font-black text-slate-900">{t('enrollments.wizard.title')}</h2>
                      <p className="text-sm text-slate-500">{t('enrollments.wizard.subtitle')}</p>
                   </div>
                   <button onClick={() => setIsWizardOpen(false)} className="p-2 hover:bg-slate-200 rounded-full text-slate-400 transition-colors">
                      <X size={20} />
                   </button>
                </div>

                {/* Wizard Steps */}
                <div className="p-4 border-b border-slate-100 flex justify-center gap-4 bg-white">
                   <WizardStep number={1} title={t('enrollments.wizard.steps.program')} active={wizardStep === 1} completed={wizardStep > 1} />
                   <WizardStep number={2} title={t('enrollments.wizard.steps.students')} active={wizardStep === 2} completed={wizardStep > 2} />
                   <WizardStep number={3} title={t('enrollments.wizard.steps.review')} active={wizardStep === 3} completed={wizardStep > 3} />
                </div>

                {/* Wizard Content */}
                <div className="flex-1 overflow-y-auto p-8">
                   {wizardStep === 1 && (
                      <div className="space-y-6">
                         <div className="space-y-2">
                            <label className="text-xs font-bold text-slate-500 uppercase">{t('enrollments.wizard.fields.program')}</label>
                            <div className="grid gap-3">
                               {programs.map(p => (
                                  <div 
                                    key={p.id} 
                                    onClick={() => setWizardData({ ...wizardData, programId: p.id })}
                                    className={cn(
                                       "p-4 rounded-xl border-2 cursor-pointer transition-all flex items-center justify-between",
                                       wizardData.programId === p.id ? "border-indigo-500 bg-indigo-50" : "border-slate-200 hover:border-indigo-200"
                                    )}
                                  >
                                     <div>
                                        <div className="font-bold text-slate-900">{p.title}</div>
                                        <div className="text-xs text-slate-500 lowercase">{p.code}</div>
                                     </div>
                                     {wizardData.programId === p.id && <CheckCircle2 className="text-indigo-600" />}
                                  </div>
                               ))}
                            </div>
                         </div>
                      </div>
                   )}

                   {wizardStep === 2 && (
                      <div className="space-y-6">
                         <div className="relative">
                            <Search size={16} className="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" />
                            <Input placeholder={t('enrollments.wizard.fields.search_students')} className="pl-9" />
                         </div>
                         <div className="text-center py-8 text-slate-400">
                            {/* Simple student selection mock for now - in reality this would be another useQuery */}
                            <p className="text-xs italic">Student search and selection implementation pending...</p>
                         </div>
                      </div>
                   )}

                   {wizardStep === 3 && (
                      <div className="space-y-6">
                         <div className="bg-slate-50 p-6 rounded-2xl border border-slate-200 space-y-4">
                            <h3 className="font-bold text-slate-900 flex items-center gap-2">
                               <CheckCircle2 className="text-emerald-500" /> {t('enrollments.wizard.fields.ready_title')}
                            </h3>
                            <div className="space-y-2 text-sm">
                               <div className="flex justify-between">
                                  <span className="text-slate-500">Program</span>
                                  <span className="font-bold text-slate-900">{programs.find(p => p.id === wizardData.programId)?.title || 'Selected Program'}</span>
                               </div>
                               <div className="flex justify-between">
                                  <span className="text-slate-500">Students</span>
                                  <span className="font-bold text-indigo-600">{wizardData.studentIds.length} selected</span>
                               </div>
                            </div>
                         </div>
                         
                         <label className="flex items-center gap-3 p-4 border border-indigo-100 bg-indigo-50/50 rounded-xl cursor-pointer">
                            <input type="checkbox" className="rounded text-indigo-600 focus:ring-indigo-500" defaultChecked />
                            <span className="text-sm font-medium text-indigo-900">{t('enrollments.wizard.fields.welcome_email')}</span>
                         </label>
                      </div>
                   )}
                </div>

                {/* Footer Buttons */}
                <div className="p-6 border-t border-slate-100 flex justify-between bg-white">
                   <Button variant="ghost" onClick={() => {
                      if (wizardStep > 1) setWizardStep(wizardStep - 1);
                      else setIsWizardOpen(false);
                   }}>
                      {wizardStep === 1 ? t('common.cancel') : t('common.back')}
                   </Button>
                   
                   {wizardStep < 3 ? (
                      <Button 
                        onClick={() => setWizardStep(wizardStep + 1)}
                        disabled={(wizardStep === 1 && !wizardData.programId) || (wizardStep === 2 && wizardData.studentIds.length === 0 && false /* bypass for dev */)}
                        className="flex items-center gap-2"
                      >
                         {t('enrollments.wizard.buttons.next')}
                         <ArrowRight size={16} />
                      </Button>
                   ) : (
                      <Button 
                        onClick={() => enrollMutation.mutate({
                           program_id: wizardData.programId,
                           cohort_id: wizardData.cohortId,
                           student_ids: wizardData.studentIds,
                           start_date: wizardData.startDate
                        })} 
                        className="bg-emerald-600 hover:bg-emerald-700 text-white shadow-lg shadow-emerald-200"
                        disabled={enrollMutation.isPending}
                      >
                         {t('enrollments.wizard.buttons.confirm')}
                      </Button>
                   )}
                </div>
             </motion.div>
           </div>
         )}
       </AnimatePresence>
    </div>
  );
};
