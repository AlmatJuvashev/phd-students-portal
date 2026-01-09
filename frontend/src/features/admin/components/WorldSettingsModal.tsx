
import React, { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { motion } from 'framer-motion';
import { X, Save, Palette, GitMerge, Trash2 } from 'lucide-react';
import { Button, Input, Switch } from '@/features/admin/components/AdminUI';
import { cn } from '@/lib/utils';

interface WorldSettingsModalProps {
  world: any;
  allNodes: any[];
  onSave: (world: any) => void;
  onClose: () => void;
  onDelete?: () => void;
}

const COLORS = [
  '#3b82f6', // Blue
  '#10b981', // Emerald
  '#f59e0b', // Amber
  '#ef4444', // Red
  '#8b5cf6', // Violet
  '#ec4899', // Pink
  '#64748b', // Slate
];

export const WorldSettingsModal: React.FC<WorldSettingsModalProps> = ({ world, allNodes, onSave, onClose, onDelete }) => {
  const [data, setData] = useState({
    title: world?.title || 'New Phase',
    description: world?.description || '',
    color: world?.color || COLORS[0],
    condition: world?.condition || null // { nodeId, fieldKey, operator, value }
  });

  const [activeTab, setActiveTab] = useState<'general' | 'logic'>('general');

  // Filter nodes that can drive logic (Forms)
  const formNodes = allNodes.filter(n => n.type === 'form');

  const handleSave = () => {
    onSave({
      ...world,
      ...data
    });
    onClose();
  };

  const { t } = useTranslation();

  return (
    <div className="fixed inset-0 z-[150] flex items-center justify-center p-4 bg-slate-900/60 backdrop-blur-sm" onClick={onClose}>
      <motion.div 
        initial={{ opacity: 0, scale: 0.95 }}
        animate={{ opacity: 1, scale: 1 }}
        exit={{ opacity: 0, scale: 0.95 }}
        className="bg-white w-full max-w-lg rounded-2xl shadow-2xl overflow-hidden flex flex-col"
        onClick={(e) => e.stopPropagation()}
      >
        <div className="p-4 border-b border-slate-100 flex justify-between items-center bg-slate-50">
          <h3 className="font-bold text-lg text-slate-800">
            {world?.id ? t('builder.world.edit_title', 'Edit Phase') : t('builder.world.create_title', 'Create Phase')}
          </h3>
          <button onClick={onClose} className="p-2 hover:bg-slate-200 rounded-full text-slate-500 transition-colors">
            <X size={20} />
          </button>
        </div>

        <div className="p-2 flex gap-2 border-b border-slate-100 px-6 bg-white">
           <button 
             onClick={() => setActiveTab('general')}
             className={cn("px-4 py-2 text-sm font-bold border-b-2 transition-colors", activeTab === 'general' ? "border-indigo-600 text-indigo-600" : "border-transparent text-slate-500 hover:text-slate-700")}
           >
             {t('builder.world.tabs.general', 'General')}
           </button>
           <button 
             onClick={() => setActiveTab('logic')}
             className={cn("px-4 py-2 text-sm font-bold border-b-2 transition-colors", activeTab === 'logic' ? "border-indigo-600 text-indigo-600" : "border-transparent text-slate-500 hover:text-slate-700")}
           >
             {t('builder.world.tabs.logic', 'Flow & Logic')}
           </button>
        </div>

        <div className="p-6 space-y-6 flex-1 overflow-y-auto bg-white min-h-[300px]">
          {activeTab === 'general' && (
            <div className="space-y-6">
              <div className="space-y-2">
                <label className="text-xs font-bold text-slate-500 uppercase tracking-wider">{t('builder.world.fields.title', 'Phase Title')}</label>
                <Input 
                  value={data.title} 
                  onChange={(e: any) => setData({ ...data, title: e.target.value })} 
                  placeholder="e.g. Preparation Stage"
                  autoFocus
                />
              </div>

              <div className="space-y-2">
                <label className="text-xs font-bold text-slate-500 uppercase tracking-wider">{t('builder.world.fields.description', 'Description')}</label>
                <textarea 
                  className="w-full p-3 bg-white border border-slate-200 rounded-xl text-sm focus:ring-2 focus:ring-indigo-100 outline-none resize-none"
                  rows={3}
                  value={data.description} 
                  onChange={(e: any) => setData({ ...data, description: e.target.value })}
                  placeholder={t('builder.world.description_placeholder', 'Briefly describe the goal of this phase...')} 
                />
              </div>

              <div className="space-y-3">
                <label className="text-xs font-bold text-slate-500 uppercase tracking-wider flex items-center gap-2">
                  <Palette size={14} /> {t('builder.world.fields.color', 'Theme Color')}
                </label>
                <div className="flex gap-3 flex-wrap">
                  {COLORS.map(c => (
                    <button
                      key={c}
                      onClick={() => setData({ ...data, color: c })}
                      className={cn(
                        "w-8 h-8 rounded-full border-2 transition-all shadow-sm",
                        data.color === c ? "border-slate-900 scale-110 ring-2 ring-slate-200" : "border-transparent hover:scale-105"
                      )}
                      style={{ backgroundColor: c }}
                    />
                  ))}
                </div>
              </div>
            </div>
          )}

          {activeTab === 'logic' && (
            <div className="space-y-6">
               <div className="bg-amber-50 border border-amber-100 rounded-xl p-4 text-xs text-amber-900 flex gap-3">
                  <GitMerge size={20} className="flex-shrink-0 text-amber-600" />
                  <div className="space-y-1">
                    <p className="font-bold">{t('builder.world.logic.conditional_visibility', 'Conditional Visibility')}</p>
                    <p className="opacity-90">{t('builder.world.logic.hint', 'Hide this phase until specific criteria are met in previous steps. Useful for branching paths (e.g., PhD vs. Masters tracks).')}</p>
                  </div>
               </div>

               <div className="space-y-4">
                  <div className="flex items-center justify-between p-3 bg-slate-50 rounded-xl border border-slate-200">
                     <span className="text-sm font-bold text-slate-700">{t('builder.world.logic.enable', 'Enable Logic Rule')}</span>
                     <Switch checked={!!data.condition} onCheckedChange={(c) => setData({ ...data, condition: c ? { nodeId: '', fieldKey: '', operator: 'equals', value: '' } : null })} />
                  </div>

                  {data.condition && (
                    <div className="p-4 bg-slate-50 rounded-xl border border-slate-200 space-y-4 animate-in slide-in-from-top-2">
                       <div className="space-y-1.5">
                          <label className="text-xs font-bold text-slate-400 uppercase">{t('builder.world.logic.based_on', 'Based on Form Step')}</label>
                          <select 
                            className="w-full p-2.5 bg-white border border-slate-200 rounded-lg text-sm outline-none focus:border-indigo-500 transition-colors"
                            value={(data.condition as any).nodeId}
                            onChange={(e) => setData({ ...data, condition: { ...(data.condition as any), nodeId: e.target.value } })}
                          >
                            <option value="">{t('builder.world.logic.select_form', 'Select a source form...')}</option>
                            {formNodes.map(n => (
                              <option key={n.id} value={n.id}>{n.title}</option>
                            ))}
                          </select>
                       </div>

                       <div className="grid grid-cols-2 gap-3">
                          <div className="space-y-1.5">
                            <label className="text-xs font-bold text-slate-400 uppercase">{t('builder.world.logic.field_key', 'Field Key')}</label>
                            <Input 
                              placeholder="e.g. student_type" 
                              value={(data.condition as any).fieldKey}
                              onChange={(e: any) => setData({ ...data, condition: { ...(data.condition as any), fieldKey: e.target.value } })}
                              className="bg-white"
                            />
                          </div>
                          <div className="space-y-1.5">
                            <label className="text-xs font-bold text-slate-400 uppercase">{t('builder.world.logic.operator', 'Operator')}</label>
                            <select 
                              className="w-full p-2.5 h-10 bg-white border border-slate-200 rounded-lg text-sm outline-none focus:border-indigo-500 transition-colors"
                              value={(data.condition as any).operator}
                              onChange={(e) => setData({ ...data, condition: { ...(data.condition as any), operator: e.target.value } })}
                            >
                              <option value="equals">{t('builder.world.logic.operators.equals', 'Equals')}</option>
                              <option value="not_equals">{t('builder.world.logic.operators.not_equals', 'Does Not Equal')}</option>
                              <option value="contains">{t('builder.world.logic.operators.contains', 'Contains')}</option>
                            </select>
                          </div>
                       </div>

                       <div className="space-y-1.5">
                          <label className="text-xs font-bold text-slate-400 uppercase">{t('builder.world.logic.value', 'Matching Value')}</label>
                          <Input 
                            placeholder="e.g. international" 
                            value={(data.condition as any).value}
                            onChange={(e: any) => setData({ ...data, condition: { ...(data.condition as any), value: e.target.value } })}
                            className="bg-white"
                          />
                       </div>
                    </div>
                  )}
               </div>
            </div>
          )}
        </div>

        <div className="p-4 bg-slate-50 border-t border-slate-100 flex justify-between items-center">
          {onDelete && world.id ? (
             <button 
               onClick={onDelete}
               className="text-red-500 hover:text-red-700 hover:bg-red-50 px-3 py-2 rounded-lg text-xs font-bold flex items-center gap-2 transition-colors"
             >
               <Trash2 size={16} /> {t('builder.world.delete_phase', 'Delete Phase')}
             </button>
          ) : <div />}
          
          <div className="flex gap-3">
            <Button variant="ghost" onClick={onClose}>{t('common.cancel', 'Cancel')}</Button>
            <Button onClick={handleSave} icon={Save}>{t('builder.world.save_changes', 'Save Changes')}</Button>
          </div>
        </div>
      </motion.div>
    </div>
  );
};
