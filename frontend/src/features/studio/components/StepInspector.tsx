import React, { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { 
  X, FileText, CheckSquare, CheckCircle2, ClipboardList, 
  Stamp, Calendar, BookOpen, CreditCard, Sparkles, Info, 
  Flag, Layout, Award, Settings2, Trash2, Zap, GitMerge, AlertCircle
} from 'lucide-react';
import { Button, Input, IconButton, Badge, Tabs, Textarea } from '@/features/admin/components/AdminUI';
import { cn } from '@/lib/utils';
import { ProgramVersionNode, ProgramNodeType, ProgramPhase } from '../types';

interface StepInspectorProps {
  node: ProgramVersionNode;
  phases: ProgramPhase[];
  onUpdate: (updates: Partial<ProgramVersionNode>) => void;
  onDelete: (id: string) => void;
  onClose: () => void;
  onNavigate: (path: string) => void;
}

const NODE_TYPE_GROUPS = [
  {
    label: "Input & Actions",
    types: ['form', 'checklist', 'survey', 'confirmTask'] as ProgramNodeType[]
  },
  {
    label: "Content & Events",
    types: ['info', 'course', 'meeting'] as ProgramNodeType[]
  },
  {
    label: "Operations",
    types: ['approval', 'payment', 'sync_ops', 'milestone'] as ProgramNodeType[]
  }
];

const NODE_VISUALS: Record<string, { icon: any, color: string, bg: string, border: string, label: string }> = {
  'course': { icon: BookOpen, color: 'text-purple-600', bg: 'bg-purple-50', border: 'border-purple-200', label: 'Learning Event' },
  'payment': { icon: CreditCard, color: 'text-amber-600', bg: 'bg-amber-50', border: 'border-amber-200', label: 'Financial Gate' },
  'sync_ops': { icon: Sparkles, color: 'text-emerald-600', bg: 'bg-emerald-50', border: 'border-emerald-200', label: 'Ops Automation' },
  'approval': { icon: Stamp, color: 'text-slate-600', bg: 'bg-slate-50', border: 'border-slate-300', label: 'Admin Gate' },
  'form': { icon: FileText, color: 'text-blue-600', bg: 'bg-blue-50', border: 'border-blue-200', label: 'Data Collection' },
  'meeting': { icon: Calendar, color: 'text-pink-600', bg: 'bg-pink-50', border: 'border-pink-200', label: 'Sync Event' },
  'checklist': { icon: CheckSquare, color: 'text-orange-600', bg: 'bg-orange-50', border: 'border-orange-200', label: 'Requirement' },
  'milestone': { icon: Flag, color: 'text-indigo-600', bg: 'bg-indigo-50', border: 'border-indigo-200', label: 'Milestone' },
  'info': { icon: Info, color: 'text-cyan-600', bg: 'bg-cyan-50', border: 'border-cyan-200', label: 'Information' },
  'survey': { icon: ClipboardList, color: 'text-teal-600', bg: 'bg-teal-50', border: 'border-teal-200', label: 'Feedback' },
  'confirmTask': { icon: CheckCircle2, color: 'text-green-600', bg: 'bg-green-50', border: 'border-green-200', label: 'Confirmation' },
  'default': { icon: Layout, color: 'text-slate-500', bg: 'bg-white', border: 'border-slate-200', label: 'Process Step' },
};

export const StepInspector: React.FC<StepInspectorProps> = ({ 
  node, 
  phases, 
  onUpdate, 
  onDelete, 
  onClose, 
  onNavigate 
}) => {
  const { t } = useTranslation();
  const [activeTab, setActiveTab] = useState('GENERAL');
  
  const visuals = NODE_VISUALS[node.type] || NODE_VISUALS['default'];
  const Icon = visuals.icon;

  return (
    <div className="w-85 bg-white border-l border-slate-200 flex flex-col z-30 shadow-2xl animate-in slide-in-from-right duration-300">
      {/* Header */}
      <div className="p-6 border-b border-slate-100 flex justify-between items-start bg-slate-50/50">
        <div className="flex-1 min-w-0">
          <div className="flex items-center gap-2 mb-2">
            <div className={cn("p-1.5 rounded-lg", visuals.bg, visuals.color)}>
              <Icon size={16} />
            </div>
            <span className="text-[10px] font-black text-slate-400 uppercase tracking-widest">{t(`builder.nodeTypes.${node.type}`, visuals.label)}</span>
          </div>
          <Input 
            value={node.title} 
            onChange={(e: any) => onUpdate({ title: e.target.value })}
            className="font-black text-xl border-transparent px-0 focus:bg-white focus:px-3 focus:border-slate-200 transition-all w-full h-auto py-1 shadow-none" 
            placeholder="Step Title"
          />
        </div>
        <IconButton icon={X} onClick={onClose} className="hover:bg-slate-200/50" />
      </div>

      {/* Tabs */}
      <div className="px-6 py-4 border-b border-slate-100">
        <Tabs 
          tabs={['GENERAL', 'SETUP', 'LOGIC']} 
          activeTab={activeTab} 
          onChange={setActiveTab} 
        />
      </div>

      {/* Content */}
      <div className="flex-1 overflow-y-auto px-6 py-6 space-y-6 custom-scrollbar">
        
        {activeTab === 'GENERAL' && (
          <div className="space-y-6 animate-in fade-in slide-in-from-right-2">
            {/* Step Function */}
            <div className="space-y-4">
              <label className="text-[10px] font-black text-slate-400 uppercase tracking-wider">{t('builder.inspector.step_function', 'Step Function')}</label>
              {NODE_TYPE_GROUPS.map((group, idx) => (
                <div key={idx} className="space-y-2">
                  <div className="text-[10px] font-bold text-slate-300 uppercase tracking-widest px-1">{group.label}</div>
                  <div className="grid grid-cols-2 gap-2">
                    {group.types.map(type => {
                      const v = NODE_VISUALS[type] || NODE_VISUALS['default'];
                      const VIcon = v.icon;
                      return (
                        <button 
                          key={type}
                          onClick={() => onUpdate({ type })}
                          className={cn(
                            "p-2.5 rounded-xl text-[11px] font-bold text-left border transition-all flex items-center gap-2",
                            node.type === type 
                              ? "bg-indigo-50 border-indigo-200 text-indigo-700 shadow-sm ring-1 ring-indigo-200" 
                              : "bg-white border-slate-100 text-slate-500 hover:border-slate-300 hover:bg-slate-50"
                          )}
                        >
                          <VIcon size={14} className={cn(node.type === type ? visuals.color : "text-slate-400")} />
                          <span className="truncate">{t(`builder.nodeTypes.${type}`, v.label)}</span>
                        </button>
                      );
                    })}
                  </div>
                </div>
              ))}
            </div>

            <div className="h-px bg-slate-100" />

            {/* Description */}
            <div className="space-y-2">
              <label className="text-[10px] font-black text-slate-400 uppercase tracking-wider">{t('builder.inspector.description', 'Internal Notes / Description')}</label>
              <Textarea 
                value={node.description || ''} 
                onChange={(e: any) => onUpdate({ description: e.target.value })}
                className="h-24 py-3 bg-slate-50/50 border-slate-200 focus:bg-white transition-colors text-xs leading-relaxed"
                placeholder="Describe the purpose of this step..."
              />
            </div>

            <div className="h-px bg-slate-100" />

            {/* Phase Selection */}
            <div className="space-y-2">
              <label className="text-[10px] font-black text-slate-400 uppercase tracking-wider">{t('builder.inspector.assigned_phase', 'Assigned Phase')}</label>
              <select 
                value={node.module_key} 
                onChange={(e) => onUpdate({ module_key: e.target.value })}
                className="w-full h-10 bg-slate-50/50 border border-slate-200 rounded-lg text-xs font-bold px-3 outline-none focus:ring-2 focus:ring-indigo-500/20"
              >
                {phases.map(p => (
                  <option key={p.id} value={p.id}>{p.title}</option>
                ))}
              </select>
            </div>

            <div className="h-px bg-slate-100" />

            {/* XP Reward */}
            <div className="space-y-2">
              <label className="text-[10px] font-black text-slate-400 uppercase tracking-wider flex items-center gap-2">
                <Award size={14} className="text-amber-500" /> {t('builder.inspector.experience_points', 'Gamification Reward')}
              </label>
              <div className="relative">
                <Input 
                  type="number" 
                  value={node.points || 0} 
                  onChange={(e: any) => onUpdate({ points: parseInt(e.target.value) || 0 })}
                  className="pl-3 font-mono font-bold text-slate-900 bg-slate-50/50 border-slate-200 focus:bg-white"
                />
                <span className="absolute right-3 top-1/2 -translate-y-1/2 text-[10px] font-black text-slate-400 uppercase">XP</span>
              </div>
            </div>
          </div>
        )}

        {activeTab === 'SETUP' && (
          <div className="space-y-6 animate-in fade-in slide-in-from-right-2">
            {/* Dynamic Config Based on Type */}
            {node.type === 'form' && (
              <div className="space-y-4">
                <div className="p-5 bg-blue-50 border border-blue-100 rounded-2xl space-y-4">
                  <div className="flex items-center gap-3">
                    <div className="w-10 h-10 rounded-xl bg-white shadow-sm flex items-center justify-center text-blue-600">
                      <FileText size={20} />
                    </div>
                    <div>
                      <h4 className="font-bold text-blue-900 text-sm">Form Composition</h4>
                      <p className="text-[10px] text-blue-700 font-medium">Define fields for data collection.</p>
                    </div>
                  </div>
                  <Badge variant="outline" className="bg-white/50 border-blue-200 text-blue-700">
                    {node.config?.fields?.length || 0} Fields defined
                  </Badge>
                  <Button 
                    className="w-full bg-blue-600 hover:bg-blue-700 text-white shadow-blue-200"
                    onClick={() => onNavigate(`/admin/studio/programs/form/${node.id}/builder`)}
                  >
                    Launch Form Designer
                  </Button>
                </div>
              </div>
            )}

            {node.type === 'course' && (
              <div className="space-y-4">
                <div className="p-5 bg-purple-50 border border-purple-100 rounded-2xl space-y-4">
                  <div className="flex items-center gap-3">
                    <div className="w-10 h-10 rounded-xl bg-white shadow-sm flex items-center justify-center text-purple-600">
                      <BookOpen size={20} />
                    </div>
                    <div>
                      <h4 className="font-bold text-purple-900 text-sm">Linked Learning Unit</h4>
                      <p className="text-[10px] text-purple-700 font-medium">Connect this step to a course.</p>
                    </div>
                  </div>
                  
                  {node.config?.course_id ? (
                    <div className="space-y-3">
                      <div className="p-3 bg-white border border-purple-200 rounded-xl">
                        <div className="text-xs font-bold text-slate-800">{node.config.course_title || 'Selected Unit'}</div>
                        <div className="text-[10px] text-slate-500 font-mono mt-0.5">{node.config.course_code || 'CODE-123'}</div>
                      </div>
                      <div className="grid grid-cols-2 gap-2">
                        <Button variant="outline" size="sm" className="bg-white" onClick={() => onNavigate(`/admin/studio/courses/${node.config.course_id}/builder`)}>
                          Curriculum
                        </Button>
                        <Button variant="secondary" size="sm" className="bg-white" onClick={() => onNavigate(`/admin/studio/courses`)}>
                          Change
                        </Button>
                      </div>
                    </div>
                  ) : (
                    <Button variant="outline" className="w-full bg-white border-purple-200 text-purple-700 hover:bg-purple-100" onClick={() => onNavigate('/admin/studio/courses')}>
                      Select Course Content
                    </Button>
                  )}
                </div>
              </div>
            )}

            {node.type === 'confirmTask' && (
              <div className="space-y-4">
                <div className="p-5 bg-emerald-50 border border-emerald-100 rounded-2xl space-y-4">
                  <div className="flex items-center gap-3">
                    <div className="w-10 h-10 rounded-xl bg-white shadow-sm flex items-center justify-center text-emerald-600">
                      <CheckCircle2 size={20} />
                    </div>
                    <div>
                      <h4 className="font-bold text-emerald-900 text-sm">Task Verification</h4>
                      <p className="text-[10px] text-emerald-700 font-medium">Verify file uploads and evidence.</p>
                    </div>
                  </div>
                  <Button 
                    variant="primary" 
                    className="w-full bg-emerald-600 hover:bg-emerald-700 shadow-emerald-200"
                    onClick={() => onNavigate(`/admin/studio/programs/confirm-task/${node.id}/builder`)}
                  >
                    Open Task Studio
                  </Button>
                </div>
              </div>
            )}

            {node.type === 'payment' && (
              <div className="space-y-4">
                <div className="p-5 bg-amber-50 border border-amber-100 rounded-2xl space-y-4">
                  <div className="flex items-center gap-3">
                    <div className="w-10 h-10 rounded-xl bg-white shadow-sm flex items-center justify-center text-amber-600">
                      <CreditCard size={20} />
                    </div>
                    <div>
                      <h4 className="font-bold text-amber-900 text-sm">Financial Gate</h4>
                      <p className="text-[10px] text-amber-700 font-medium">Configure fee and gateway.</p>
                    </div>
                  </div>
                  <div className="grid grid-cols-2 gap-3">
                    <div className="space-y-1">
                      <label className="text-[10px] font-black text-amber-700 uppercase">Amount</label>
                      <Input 
                        type="number" 
                        value={node.config?.amount || 0} 
                        onChange={(e: any) => onUpdate({ config: { ...node.config, amount: parseInt(e.target.value) || 0 } })}
                        className="bg-white border-amber-200 h-9"
                      />
                    </div>
                    <div className="space-y-1">
                      <label className="text-[10px] font-black text-amber-700 uppercase">Currency</label>
                      <select className="w-full h-9 bg-white border border-amber-200 rounded-lg text-xs font-bold px-2">
                        <option>KZT</option>
                        <option>USD</option>
                      </select>
                    </div>
                  </div>
                </div>
              </div>
            )}

            {node.type === 'approval' && (
              <div className="space-y-4">
                <div className="p-5 bg-slate-50 border border-slate-200 rounded-2xl space-y-4">
                  <div className="flex items-center gap-3">
                    <div className="w-10 h-10 rounded-xl bg-white shadow-sm flex items-center justify-center text-slate-600">
                      <Stamp size={20} />
                    </div>
                    <div>
                      <h4 className="font-bold text-slate-900 text-sm">Administrative Gate</h4>
                      <p className="text-[10px] text-slate-500 font-medium">Configure approval workflow.</p>
                    </div>
                  </div>
                  <div className="space-y-2">
                    <label className="text-[10px] font-black text-slate-400 uppercase tracking-wider">Required Participant</label>
                    <select 
                      className="w-full h-10 bg-white border border-slate-200 rounded-xl text-xs font-bold px-3 outline-none shadow-sm transition-all focus:ring-2 focus:ring-indigo-100"
                      value={node.config?.role || ''}
                      onChange={(e) => onUpdate({ config: { ...node.config, role: e.target.value } })}
                    >
                      <option value="">Select Role...</option>
                      <option value="advisor">Scientific Advisor</option>
                      <option value="secretary">Academic Secretary</option>
                      <option value="dean">Head of Faculty</option>
                      <option value="registrar">Registrar Office</option>
                    </select>
                  </div>
                </div>
              </div>
            )}

            {!['form', 'course', 'confirmTask', 'payment', 'approval'].includes(node.type) && (
              <div className="p-10 text-center space-y-3">
                <div className="w-16 h-16 rounded-full bg-slate-50 flex items-center justify-center mx-auto text-slate-300">
                  <Settings2 size={32} />
                </div>
                <p className="text-sm font-bold text-slate-400 italic">No specific configuration required for this step type.</p>
              </div>
            )}
          </div>
        )}

        {activeTab === 'LOGIC' && (
          <div className="space-y-6 animate-in fade-in slide-in-from-right-2">
            {/* Sync Ops / Automation */}
            {node.type === 'sync_ops' && (
              <div className="p-5 bg-emerald-50 border border-emerald-100 rounded-2xl space-y-4">
                <div className="flex items-center gap-3">
                  <div className="w-10 h-10 rounded-xl bg-white shadow-sm flex items-center justify-center text-emerald-600">
                    <Zap size={20} fill="currentColor" />
                  </div>
                  <div>
                    <h4 className="font-bold text-emerald-900 text-sm">Auto-Workflow</h4>
                    <p className="text-[10px] text-emerald-700 font-medium">Background system trigger.</p>
                  </div>
                </div>
                <div className="space-y-2">
                  <label className="text-[10px] font-black text-emerald-600 uppercase tracking-wider">Trigger Action</label>
                  <select 
                    className="w-full h-10 bg-white border border-emerald-200 rounded-xl text-xs font-bold px-3 text-slate-700 outline-none shadow-sm focus:ring-2 focus:ring-emerald-500/20"
                    value={node.config?.action || ''}
                    onChange={(e) => onUpdate({ config: { ...node.config, action: e.target.value } })}
                  >
                    <option value="">Select Logic...</option>
                    <option value="assign_advisor">Assign Department Advisor</option>
                    <option value="unlock_cohort">Provision Next Cohort</option>
                    <option value="generate_credentials">Issue Digital Credentials</option>
                    <option value="archive_portfolio">Archive Evidence Portfolio</option>
                  </select>
                </div>
              </div>
            )}

            {/* Prerequisites */}
            <div className="space-y-4">
              <label className="text-[10px] font-black text-slate-400 uppercase tracking-wider flex items-center gap-2">
                 <GitMerge size={14} className="text-slate-400" /> Progression Logic
              </label>
              <div className="p-4 border border-dashed border-slate-200 rounded-2xl text-center space-y-2">
                 <p className="text-[11px] text-slate-400 font-medium italic">Advanced dependency logic can be managed in the Map view via edges.</p>
                 <Badge variant="outline" className="text-[9px]">v2.2 Pipeline</Badge>
              </div>
            </div>

            <div className="h-px bg-slate-100" />

            {/* Step ID / Slug */}
            <div className="space-y-2">
              <label className="text-[10px] font-black text-slate-400 uppercase tracking-wider">Unique Identifier (Slug)</label>
              <Input 
                value={node.slug} 
                onChange={(e: any) => onUpdate({ slug: e.target.value })}
                className="bg-slate-50 font-mono text-xs border-slate-200"
              />
              <p className="text-[10px] text-slate-400 font-medium leading-relaxed">System-wide key for API calls and conditional logic across the PhD portal.</p>
            </div>
          </div>
        )}

      </div>

      {/* Footer */}
      <div className="p-6 border-t border-slate-100 bg-slate-50/80 backdrop-blur-sm">
        <Button 
          variant="secondary" 
          className="w-full bg-white text-red-500 hover:text-red-700 hover:bg-red-50 border-red-100 hover:border-red-200"
          icon={Trash2}
          onClick={() => {
            if (confirm(t('builder.inspector.confirm_delete', 'Are you sure you want to remove this step? This cannot be undone.'))) {
              onDelete(node.id);
            }
          }}
        >
          {t('builder.inspector.delete_step', 'Delete Step')}
        </Button>
      </div>
    </div>
  );
};
