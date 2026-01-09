
import React, { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { useQuery } from '@tanstack/react-query';
import { Globe, Save, BookOpen, Trash2 } from 'lucide-react';
import { Modal, Button, Input, Textarea, Badge } from '@/features/admin/components/AdminUI';
import { cn } from '@/lib/utils';
import { getPrograms } from '../api';
import { Program } from '../types';

interface CourseModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSave: (data: any) => void;
  onDelete?: () => void;
  initialData?: any;
  isLoading?: boolean;
}

export const CourseModal: React.FC<CourseModalProps> = ({ 
  isOpen, 
  onClose, 
  onSave,
  onDelete,
  initialData,
  isLoading 
}) => {
  const { t } = useTranslation('common');
  const [lang, setLang] = useState<'en' | 'ru' | 'kk'>('en');
  
  const { data: programs = [] } = useQuery({ 
      queryKey: ['curriculum', 'programs'], 
      queryFn: getPrograms,
      enabled: isOpen 
  });

  const [formData, setFormData] = useState({
    code: '',
    program_id: '',
    title: { en: '', ru: '', kk: '' },
    description: { en: '', ru: '', kk: '' },
    credits: 5,
    workload_hours: 150
  });

  useEffect(() => {
    if (initialData && isOpen) {
      setFormData({
        code: initialData.code || '',
        program_id: initialData.program_id || '',
        title: typeof initialData.title === 'string' ? { en: initialData.title, ru: '', kk: '' } : { ...initialData.title },
        description: typeof initialData.description === 'string' ? { en: initialData.description, ru: '', kk: '' } : { ...initialData.description },
        credits: initialData.credits || 5,
        workload_hours: initialData.workload_hours || 150
      });
    } else if (!initialData && isOpen) {
        setFormData({
            code: '',
            program_id: '',
            title: { en: '', ru: '', kk: '' },
            description: { en: '', ru: '', kk: '' },
            credits: 5,
            workload_hours: 150
        });
    }
  }, [initialData, isOpen]);

  const handleSave = () => {
    if (!formData.code || !formData.title.en) return;
    onSave(formData);
  };

  const updateTitle = (val: string) => {
    setFormData(prev => ({
      ...prev,
      title: { ...prev.title, [lang]: val }
    }));
  };

  const updateDescription = (val: string) => {
    setFormData(prev => ({
      ...prev,
      description: { ...prev.description, [lang]: val }
    }));
  };

  const footer = (
    <div className="flex justify-between w-full">
      {onDelete && initialData ? (
        <Button variant="ghost" className="text-red-500 hover:text-red-700 hover:bg-red-50" onClick={onDelete} disabled={isLoading}>
           <Trash2 size={16} className="mr-2" /> {t('common.delete', 'Delete')}
        </Button>
      ) : <div />}
      <div className="flex gap-3">
        <Button variant="secondary" onClick={onClose} disabled={isLoading}>
            {t('common.cancel', 'Cancel')}
        </Button>
        <Button onClick={handleSave} isLoading={isLoading} icon={Save}>
            {initialData ? t('common.save', 'Save Changes') : t('common.create', 'Create Course')}
        </Button>
      </div>
    </div>
  );

  return (
    <Modal 
      isOpen={isOpen} 
      onClose={onClose} 
      title={initialData ? "Edit Course Metadata" : "Register New Global Course"} 
      footer={footer}
    >
      <div className="space-y-6">
        <div className="grid grid-cols-2 gap-4">
          <div className="space-y-1.5">
            <label className="text-[10px] font-black text-slate-400 uppercase tracking-widest">Course Code</label>
            <Input 
              placeholder="e.g. BIO-501" 
              value={formData.code}
              onChange={(e: any) => setFormData({...formData, code: e.target.value.toUpperCase()})}
            />
          </div>
          <div className="space-y-1.5">
            <label className="text-[10px] font-black text-slate-400 uppercase tracking-widest">Parent Program (Optional)</label>
            <select 
              className="w-full h-10 px-3 rounded-lg border border-slate-200 text-sm font-medium outline-none focus:ring-2 focus:ring-indigo-600 focus:ring-offset-1 transition-all"
              value={formData.program_id}
              onChange={(e: any) => setFormData({...formData, program_id: e.target.value})}
            >
              <option value="">-- No Parent Program --</option>
              {programs.map((p: Program) => (
                <option key={p.id} value={p.id}>{typeof p.title === 'string' ? p.title : (p.title as any).en}</option>
              ))}
            </select>
          </div>
        </div>

        {/* Localization Switcher */}
        <div className="pt-4 border-t border-slate-100">
           <div className="flex items-center justify-between mb-4">
              <label className="text-[10px] font-black text-slate-400 uppercase tracking-widest flex items-center gap-2">
                <Globe size={12} className="text-indigo-500" /> Course Content
              </label>
              <div className="flex bg-slate-100 p-1 rounded-lg">
                {(['en', 'ru', 'kk'] as const).map(l => (
                  <button
                    key={l}
                    onClick={() => setLang(l)}
                    className={cn(
                      "px-3 py-1 text-[10px] font-black uppercase rounded-md transition-all",
                      lang === l ? "bg-white text-indigo-600 shadow-sm" : "text-slate-500 hover:text-slate-700"
                    )}
                  >
                    {l}
                  </button>
                ))}
              </div>
           </div>

           <div className="space-y-4 animate-in fade-in duration-300">
              <div className="space-y-1.5">
                <div className="flex justify-between items-center">
                   <label className="text-xs font-bold text-slate-700">Course Title ({lang.toUpperCase()})</label>
                   {formData.title[lang] ? <Badge variant="success" className="text-[8px]">Filled</Badge> : <Badge variant="outline" className="text-[8px]">Required</Badge>}
                </div>
                <Input 
                  placeholder={`Enter title in ${lang.toUpperCase()}...`}
                  value={formData.title[lang]}
                  onChange={(e: any) => updateTitle(e.target.value)}
                />
              </div>

              <div className="space-y-1.5">
                <label className="text-xs font-bold text-slate-700">Course Scope ({lang.toUpperCase()})</label>
                <Textarea 
                  rows={3}
                  placeholder={`Summary of what this course covers...`}
                  value={formData.description[lang]}
                  onChange={(e: any) => updateDescription(e.target.value)}
                />
              </div>
           </div>
        </div>

        <div className="grid grid-cols-2 gap-4 pt-4 border-t border-slate-100">
           <div className="space-y-1.5">
              <label className="text-[10px] font-black text-slate-400 uppercase tracking-widest">Academic Credits (ECTS)</label>
              <Input 
                type="number" 
                value={formData.credits}
                onChange={(e: any) => setFormData({...formData, credits: parseInt(e.target.value) || 0})}
              />
           </div>
           <div className="space-y-1.5">
              <label className="text-[10px] font-black text-slate-400 uppercase tracking-widest">Study Load (Hours)</label>
              <Input 
                type="number" 
                value={formData.workload_hours}
                onChange={(e: any) => setFormData({...formData, workload_hours: parseInt(e.target.value) || 0})}
              />
           </div>
        </div>
      </div>
    </Modal>
  );
};
