
import React, { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { Globe, Save, Trash2 } from 'lucide-react';
import { Modal, Button, Input, Textarea, Badge } from '@/features/admin/components/AdminUI';
import { cn } from '@/lib/utils';

interface ProgramModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSave: (data: any) => void;
  onDelete?: () => void;
  initialData?: any;
  isLoading?: boolean;
}

export const ProgramModal: React.FC<ProgramModalProps> = ({ 
  isOpen, 
  onClose, 
  onSave,
  onDelete,
  initialData,
  isLoading 
}) => {
  const { t } = useTranslation('common');
  const [lang, setLang] = useState<'en' | 'ru' | 'kk'>('en');
  
  const [formData, setFormData] = useState({
    code: '',
    title: { en: '', ru: '', kk: '' },
    description: { en: '', ru: '', kk: '' },
    type: 'doctoral' as const,
    total_credits: 180,
    duration_semesters: 6
  });

  useEffect(() => {
    if (initialData && isOpen) {
      setFormData({
        code: initialData.code || '',
        title: typeof initialData.title === 'string' ? { en: initialData.title, ru: '', kk: '' } : { ...initialData.title },
        description: typeof initialData.description === 'string' ? { en: initialData.description, ru: '', kk: '' } : { ...initialData.description },
        type: initialData.type || 'doctoral',
        total_credits: initialData.total_credits || 180,
        duration_semesters: initialData.duration_semesters || 6
      });
    } else if (!initialData && isOpen) {
        setFormData({
            code: '',
            title: { en: '', ru: '', kk: '' },
            description: { en: '', ru: '', kk: '' },
            type: 'doctoral',
            total_credits: 180,
            duration_semesters: 6
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
            {initialData ? t('common.save', 'Save Changes') : t('common.create', 'Create Program')}
        </Button>
      </div>
    </div>
  );

  return (
    <Modal 
      isOpen={isOpen} 
      onClose={onClose} 
      title={initialData ? "Edit Program Metadata" : "Create New PhD Program"} 
      footer={footer}
    >
      <div className="space-y-6">
        {/* Basic Metadata */}
        <div className="grid grid-cols-2 gap-4">
          <div className="space-y-1.5">
            <label className="text-[10px] font-black text-slate-400 uppercase tracking-widest">Program Code / ID</label>
            <Input 
              placeholder="e.g. PHD-MED-2026" 
              value={formData.code}
              onChange={(e: any) => setFormData({...formData, code: e.target.value.toUpperCase()})}
            />
          </div>
          <div className="space-y-1.5">
            <label className="text-[10px] font-black text-slate-400 uppercase tracking-widest">Program Type</label>
            <select 
              className="w-full h-10 px-3 rounded-lg border border-slate-200 text-sm font-medium outline-none focus:ring-2 focus:ring-indigo-600 focus:ring-offset-1 transition-all"
              value={formData.type}
              onChange={(e: any) => setFormData({...formData, type: e.target.value as any})}
            >
              <option value="doctoral">Doctoral (PhD)</option>
              <option value="master">Master's</option>
              <option value="bachelor">Bachelor's</option>
              <option value="certificate">Professional Certificate</option>
            </select>
          </div>
        </div>

        {/* Localization Switcher */}
        <div className="pt-4 border-t border-slate-100">
           <div className="flex items-center justify-between mb-4">
              <label className="text-[10px] font-black text-slate-400 uppercase tracking-widest flex items-center gap-2">
                <Globe size={12} className="text-indigo-500" /> Localized Content
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
                   <label className="text-xs font-bold text-slate-700">Program Title ({lang.toUpperCase()})</label>
                   {formData.title[lang] ? <Badge variant="success" className="text-[8px]">Filled</Badge> : <Badge variant="outline" className="text-[8px]">Required</Badge>}
                </div>
                <Input 
                  placeholder={`Enter title in ${lang.toUpperCase()}...`}
                  value={formData.title[lang]}
                  onChange={(e: any) => updateTitle(e.target.value)}
                />
              </div>

              <div className="space-y-1.5">
                <label className="text-xs font-bold text-slate-700">Description ({lang.toUpperCase()})</label>
                <Textarea 
                  rows={4}
                  placeholder={`Describe the program goals in ${lang.toUpperCase()}...`}
                  value={formData.description[lang]}
                  onChange={(e: any) => updateDescription(e.target.value)}
                />
              </div>
           </div>
        </div>

        {/* Stats */}
        <div className="grid grid-cols-2 gap-4 pt-4 border-t border-slate-100">
           <div className="space-y-1.5">
              <label className="text-[10px] font-black text-slate-400 uppercase tracking-widest">Total Credits (ECTS)</label>
              <Input 
                type="number" 
                value={formData.total_credits}
                onChange={(e: any) => setFormData({...formData, total_credits: parseInt(e.target.value) || 0})}
              />
           </div>
           <div className="space-y-1.5">
              <label className="text-[10px] font-black text-slate-400 uppercase tracking-widest">Duration (Semesters)</label>
              <Input 
                type="number" 
                value={formData.duration_semesters}
                onChange={(e: any) => setFormData({...formData, duration_semesters: parseInt(e.target.value) || 0})}
              />
           </div>
        </div>
      </div>
    </Modal>
  );
};
