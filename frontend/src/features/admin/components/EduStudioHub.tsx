import React from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { 
  X, CheckSquare, AlignLeft, Activity, Sparkles, 
  ShieldCheck, BookOpen, Layers, MousePointer2, LayoutDashboard
} from 'lucide-react';
import { cn } from '@/lib/utils';
import { useTranslation } from 'react-i18next';

interface EduStudioHubProps {
  isOpen: boolean;
  onClose: () => void;
  onNavigate: (path: string) => void;
}

const ChoiceBtn = ({ icon: Icon, label, description, onClick, color }: any) => (
  <button 
    onClick={onClick} 
    className="group flex flex-col items-center justify-center p-6 border-2 border-slate-50 rounded-[2.5rem] hover:border-indigo-500 hover:bg-indigo-50/30 transition-all bg-white shadow-sm text-center"
  >
    <div className={cn("p-4 rounded-2xl mb-4 group-hover:scale-110 transition-transform shadow-sm", color)}>
      <Icon size={28} />
    </div>
    <span className="text-sm font-black text-slate-800 mb-1">{label}</span>
    <span className="text-[10px] text-slate-400 font-medium leading-tight px-2">{description}</span>
  </button>
);

export const EduStudioHub: React.FC<EduStudioHubProps> = ({ isOpen, onClose, onNavigate }) => {
  const { t } = useTranslation('common');

  const handleCreateQuestion = (type: string) => {
    onClose();
    onNavigate(`/admin/item-bank/questions/new?type=${type}`);
  };

  const handleCreateCurriculum = (path: string) => {
    onClose();
    onNavigate(path);
  };

  return (
    <AnimatePresence>
      {isOpen && (
        <div className="fixed inset-0 z-[200] flex items-center justify-center p-4 bg-slate-950/60 backdrop-blur-md">
          <motion.div
            initial={{ opacity: 0, scale: 0.9, y: 20 }}
            animate={{ opacity: 1, scale: 1, y: 0 }}
            exit={{ opacity: 0, scale: 0.9, y: 20 }}
            className="bg-slate-50 w-full max-w-4xl rounded-[3rem] shadow-2xl overflow-hidden flex flex-col border border-white/20"
            onClick={(e) => e.stopPropagation()}
          >
            {/* Header */}
            <div className="p-8 border-b border-slate-200 flex justify-between items-center bg-white">
              <div>
                <h2 className="text-3xl font-black text-slate-900 flex items-center gap-3 tracking-tight">
                  <Sparkles className="text-indigo-500 fill-indigo-500" size={32} /> {t('edustudio.title')}
                </h2>
                <p className="text-slate-500 text-sm font-medium mt-1">{t('edustudio.subtitle')}</p>
              </div>
              <button onClick={onClose} className="p-3 hover:bg-slate-100 rounded-full transition-all text-slate-400">
                <X size={24} />
              </button>
            </div>

            <div className="p-10 space-y-12 overflow-y-auto max-h-[70vh] custom-scrollbar">
              {/* Category 1: Assessment Design */}
              <section className="space-y-6">
                <div className="flex items-center gap-3 px-2">
                  <div className="h-px flex-1 bg-slate-200" />
                  <h3 className="text-[11px] font-black text-slate-400 uppercase tracking-[0.2em]">{t('edustudio.cat.assessment')}</h3>
                  <div className="h-px flex-1 bg-slate-200" />
                </div>
                <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
                  <ChoiceBtn 
                    icon={CheckSquare} 
                    label={t('edustudio.btn.multiselect.label')} 
                    description={t('edustudio.btn.multiselect.desc')}
                    color="bg-blue-50 text-blue-600"
                    onClick={() => handleCreateQuestion('multi_select')} 
                  />
                  <ChoiceBtn 
                    icon={AlignLeft} 
                    label={t('edustudio.btn.short.label')} 
                    description={t('edustudio.btn.short.desc')}
                    color="bg-emerald-50 text-emerald-600"
                    onClick={() => handleCreateQuestion('short_answer')} 
                  />
                  <ChoiceBtn 
                    icon={Activity} 
                    label={t('edustudio.btn.osce.label')} 
                    description={t('edustudio.btn.osce.desc')}
                    color="bg-purple-50 text-purple-600"
                    onClick={() => handleCreateQuestion('osce')} 
                  />
                </div>
              </section>

              {/* Category 2: Curriculum Strategy */}
              <section className="space-y-6">
                <div className="flex items-center gap-3 px-2">
                  <div className="h-px flex-1 bg-slate-200" />
                  <h3 className="text-[11px] font-black text-slate-400 uppercase tracking-[0.2em]">{t('edustudio.cat.curriculum')}</h3>
                  <div className="h-px flex-1 bg-slate-200" />
                </div>
                <div className="grid grid-cols-1 sm:grid-cols-2 gap-4 max-w-2xl mx-auto">
                  <ChoiceBtn 
                    icon={BookOpen} 
                    label={t('edustudio.btn.course.label')} 
                    description={t('edustudio.btn.course.desc')}
                    color="bg-orange-50 text-orange-600"
                    onClick={() => handleCreateCurriculum('/admin/studio/courses/new/builder')} 
                  />
                  <ChoiceBtn 
                    icon={Layers} 
                    label={t('edustudio.btn.program.label')} 
                    description={t('edustudio.btn.program.desc')}
                    color="bg-indigo-50 text-indigo-600"
                    onClick={() => handleCreateCurriculum('/admin/studio/programs/new/journey')} 
                  />
                </div>
              </section>


              {/* Category 3: Quick Access Dashboards */}
              <section className="space-y-6">
                 <div className="flex items-center gap-3 px-2">
                   <div className="h-px flex-1 bg-slate-200" />
                   <h3 className="text-[11px] font-black text-slate-400 uppercase tracking-[0.2em]">{t('edustudio.cat.quick')}</h3>
                   <div className="h-px flex-1 bg-slate-200" />
                 </div>
                 <div className="grid grid-cols-1 sm:grid-cols-2 gap-4 max-w-2xl mx-auto">
                   <ChoiceBtn 
                     icon={LayoutDashboard} 
                     label={t('edustudio.btn.student.label')} 
                     description={t('edustudio.btn.student.desc')}
                     color="bg-emerald-50 text-emerald-600"
                     onClick={() => handleCreateCurriculum('/student/dashboard')} 
                   />
                   <ChoiceBtn 
                     icon={Activity} 
                     label={t('edustudio.btn.instructor.label')} 
                     description={t('edustudio.btn.instructor.desc')}
                     color="bg-indigo-50 text-indigo-600"
                     onClick={() => handleCreateCurriculum('/teach/dashboard')} 
                   />
                 </div>
              </section>

              {/* Category 4: Operations & Campus */}
              <section className="space-y-6">
                 <div className="flex items-center gap-3 px-2">
                   <div className="h-px flex-1 bg-slate-200" />
                   <h3 className="text-[11px] font-black text-slate-400 uppercase tracking-[0.2em]">{t('edustudio.cat.operations')}</h3>
                   <div className="h-px flex-1 bg-slate-200" />
                 </div>
                 <div className="grid grid-cols-1 sm:grid-cols-2 gap-4 max-w-2xl mx-auto">
                   <ChoiceBtn 
                     icon={ShieldCheck} 
                     label={t('edustudio.btn.hr.label')} 
                     description={t('edustudio.btn.hr.desc')}
                     color="bg-slate-50 text-slate-600"
                     onClick={() => handleCreateCurriculum('/admin/hr')} 
                   />
                   <ChoiceBtn 
                     icon={ShieldCheck} 
                     label={t('edustudio.btn.facilities.label')} 
                     description={t('edustudio.btn.facilities.desc')}
                     color="bg-orange-50 text-orange-600"
                     onClick={() => handleCreateCurriculum('/admin/facilities')} 
                   />
                 </div>
              </section>
            </div>

            <div className="p-6 bg-white border-t border-slate-200 flex items-center justify-between">
               <div className="flex items-center gap-2 text-[10px] font-black text-slate-400 uppercase tracking-widest">
                  <ShieldCheck size={14} className="text-emerald-500" /> Standardized Design Schema v4.2
               </div>
               <div className="flex items-center gap-1 text-[10px] font-bold text-indigo-500">
                 <MousePointer2 size={12} /> SELECT AN AUTHORING ENGINE
               </div>
            </div>
          </motion.div>
        </div>
      )}
    </AnimatePresence>
  );
};
