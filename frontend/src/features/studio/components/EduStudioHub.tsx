import React from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { 
  X, CheckSquare, AlignLeft, Activity, Sparkles, 
  BookOpen, Layers, MousePointer2, ShieldCheck 
} from 'lucide-react';
import { cn } from '@/lib/utils';
import { Dialog, DialogContent, DialogClose } from '@/components/ui/dialog';
import { Button } from '@/components/ui/button';

interface EduStudioHubProps {
  isOpen: boolean;
  onClose: () => void;
  onNavigate: (path: string) => void;
}

const ChoiceBtn = ({ icon: Icon, label, description, onClick, color }: any) => (
  <button 
    onClick={onClick} 
    className="group flex flex-col items-center justify-center p-6 border-2 border-slate-50 rounded-[2.5rem] hover:border-indigo-500 hover:bg-indigo-50/30 transition-all bg-white shadow-sm text-center w-full"
  >
    <div className={cn("p-4 rounded-2xl mb-4 group-hover:scale-110 transition-transform shadow-sm", color)}>
      <Icon size={28} />
    </div>
    <span className="text-sm font-black text-slate-800 mb-1">{label}</span>
    <span className="text-[10px] text-slate-400 font-medium leading-tight px-2">{description}</span>
  </button>
);

export const EduStudioHub: React.FC<EduStudioHubProps> = ({ isOpen, onClose, onNavigate }) => {
  const handleCreateQuestion = (type: string) => {
    onClose();
    onNavigate(`/admin/item-banks?type=${encodeURIComponent(type)}`);
  };

  const handleCreateCurriculum = (path: string) => {
    onClose();
    onNavigate(path);
  };

  if (!isOpen) return null;

  return (
    <Dialog open={isOpen} onOpenChange={onClose}>
        <DialogContent className="max-w-4xl p-0 bg-transparent border-none shadow-none text-slate-900">
          <motion.div
            initial={{ opacity: 0, scale: 0.9, y: 20 }}
            animate={{ opacity: 1, scale: 1, y: 0 }}
            exit={{ opacity: 0, scale: 0.9, y: 20 }}
            className="bg-slate-50 w-full rounded-[3rem] shadow-2xl overflow-hidden flex flex-col border border-white/20"
            onClick={(e) => e.stopPropagation()}
          >
            {/* Header */}
            <div className="p-8 border-b border-slate-200 flex justify-between items-center bg-white">
              <div>
                <h2 className="text-3xl font-black text-slate-900 flex items-center gap-3 tracking-tight">
                  <Sparkles className="text-indigo-500 fill-indigo-500" size={32} /> EduStudio Hub
                </h2>
                <p className="text-slate-500 text-sm font-medium mt-1">Central command for all educational assets</p>
              </div>
              <DialogClose asChild>
                <button className="p-3 hover:bg-slate-100 rounded-full transition-all text-slate-400">
                    <X size={24} />
                </button>
              </DialogClose>
            </div>

            <div className="p-10 space-y-12 overflow-y-auto max-h-[70vh] custom-scrollbar">
              {/* Category 1: Assessment Design */}
              <section className="space-y-6">
                <div className="flex items-center gap-3 px-2">
                  <div className="h-px flex-1 bg-slate-200" />
                  <h3 className="text-[11px] font-black text-slate-400 uppercase tracking-[0.2em]">Assessment Design</h3>
                  <div className="h-px flex-1 bg-slate-200" />
                </div>
                <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
                  <ChoiceBtn 
                    icon={CheckSquare} 
                    label="Multi-select" 
                    description="Standard MCQ with single or multi keys"
                    color="bg-blue-50 text-blue-600"
                    onClick={() => handleCreateQuestion('multi_select')} 
                  />
                  <ChoiceBtn 
                    icon={AlignLeft} 
                    label="Short Answer" 
                    description="Open-ended text with regex matching"
                    color="bg-emerald-50 text-emerald-600"
                    onClick={() => handleCreateQuestion('short_answer')} 
                  />
                  <ChoiceBtn 
                    icon={Activity} 
                    label="OSCE Station" 
                    description="Complex scenario with examiner rubric"
                    color="bg-purple-50 text-purple-600"
                    onClick={() => handleCreateQuestion('osce')} 
                  />
                </div>
              </section>

              {/* Category 2: Curriculum Strategy */}
              <section className="space-y-6">
                <div className="flex items-center gap-3 px-2">
                  <div className="h-px flex-1 bg-slate-200" />
                  <h3 className="text-[11px] font-black text-slate-400 uppercase tracking-[0.2em]">Curriculum Strategy</h3>
                  <div className="h-px flex-1 bg-slate-200" />
                </div>
                <div className="grid grid-cols-1 sm:grid-cols-2 gap-4 max-w-2xl mx-auto">
                  <ChoiceBtn 
                    icon={BookOpen} 
                    label="Course" 
                    description="Structure modules, lessons and readings"
                    color="bg-orange-50 text-orange-600"
                    onClick={() => handleCreateCurriculum('/admin/courses')} 
                  />
                  <ChoiceBtn 
                    icon={Layers} 
                    label="Program" 
                    description="Certification tracks and residency maps"
                    color="bg-indigo-50 text-indigo-600"
                    onClick={() => handleCreateCurriculum('/admin/programs')} 
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
        </DialogContent>
    </Dialog>
  );
};
