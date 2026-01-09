
import React, { useState } from 'react';
import { 
  ArrowLeft, Save, UploadCloud, FileText, Trash2, 
  CheckCircle2, Shield, Download, Paperclip
} from 'lucide-react';
import { Button, Input, Switch, Badge, IconButton, AvatarGroup } from '@/features/admin/components/AdminUI';
import { cn } from '@/lib/utils';
import { useNavigate, useParams } from 'react-router-dom';

const ACTIVE_DESIGNERS = [
  { initials: 'JS', color: 'bg-emerald-500' },
];

interface UploadRequirement {
  id: string;
  label: string;
  required: boolean;
}

interface ConfirmConfig {
  title: string;
  intro: string;
  buttonText: string;
  reviewer: string;
  templates: { name: string; size: string }[];
  uploads: UploadRequirement[];
}

export const ConfirmTaskBuilder: React.FC = () => {
  const navigate = useNavigate();
  const { programId, nodeId } = useParams();
  const onNavigate = (path: string) => navigate(path);

  const [config, setConfig] = useState<ConfirmConfig>({
    title: 'Thesis Final Submission',
    intro: 'Please upload your final thesis document and the signed approval form.',
    buttonText: 'Confirm Submission',
    reviewer: 'advisor',
    templates: [],
    uploads: [
      { id: 'u1', label: 'Final Thesis PDF', required: true }
    ]
  });

  const addUpload = () => {
    setConfig({ ...config, uploads: [...config.uploads, { id: Date.now().toString(), label: 'New Document', required: true }] });
  };

  const updateUpload = (id: string, updates: Partial<UploadRequirement>) => {
    setConfig({ ...config, uploads: config.uploads.map(u => u.id === id ? { ...u, ...updates } : u) });
  };

  const deleteUpload = (id: string) => {
    setConfig({ ...config, uploads: config.uploads.filter(u => u.id !== id) });
  };

  const addTemplate = () => {
    setConfig({ ...config, templates: [...config.templates, { name: 'Template.docx', size: '15 KB' }] });
  };

  return (
    <div className="flex flex-col h-[calc(100vh-4rem)] bg-slate-50 font-sans overflow-hidden">
      <div className="h-20 bg-white border-b border-slate-200 px-8 flex items-center justify-between flex-shrink-0 z-30 shadow-sm">
        <div className="flex items-center gap-6">
          <IconButton icon={ArrowLeft} onClick={() => onNavigate(`/admin/studio/programs/${programId}/builder`)} />
          <div>
             <div className="flex items-center gap-2 mb-1">
                <span className="text-[9px] font-black uppercase text-emerald-600 bg-emerald-50 px-2 py-0.5 rounded-full tracking-widest border border-emerald-100">Task Studio</span>
             </div>
             <Input 
               value={config.title} 
               onChange={(e: any) => setConfig({ ...config, title: e.target.value })}
               className="font-black text-slate-900 text-xl border-none p-0 h-auto focus:ring-0 w-96 bg-transparent"
             />
          </div>
        </div>
        <div className="flex items-center gap-6">
          <AvatarGroup users={ACTIVE_DESIGNERS} />
          <Button variant="primary" icon={Save} onClick={() => alert('Saved!')}>Publish Task</Button>
        </div>
      </div>

      <div className="flex-1 flex overflow-hidden">
        {/* Main Editor */}
        <div className="flex-1 bg-slate-100/50 overflow-y-auto p-8 flex flex-col items-center">
           <div className="w-full max-w-3xl space-y-8">
              
              {/* Intro Section */}
              <div className="bg-white p-8 rounded-[2rem] shadow-sm border border-slate-200">
                 <h3 className="text-lg font-black text-slate-900 mb-4">Instructions</h3>
                 <textarea 
                   value={config.intro}
                   onChange={(e) => setConfig({ ...config, intro: e.target.value })}
                   className="w-full p-4 bg-slate-50 border-none rounded-xl text-sm min-h-[120px] resize-none focus:ring-2 focus:ring-emerald-100 text-slate-600 leading-relaxed"
                   placeholder="Instructions for the student..."
                 />
              </div>

              {/* Uploads Section */}
              <div className="bg-white p-8 rounded-[2rem] shadow-sm border border-slate-200">
                 <div className="flex justify-between items-center mb-6">
                    <h3 className="text-lg font-black text-slate-900 flex items-center gap-2">
                       <UploadCloud className="text-emerald-500" /> Required Submissions
                    </h3>
                    <Button size="sm" variant="secondary" icon={UploadCloud} onClick={addUpload}>Add Upload Slot</Button>
                 </div>
                 
                 <div className="space-y-3">
                    {config.uploads.map((upload, i) => (
                       <div key={upload.id} className="flex items-center gap-4 p-4 bg-slate-50 border border-slate-200 rounded-2xl group transition-all hover:border-emerald-200">
                          <div className="w-8 h-8 rounded-full bg-white flex items-center justify-center font-bold text-xs text-slate-400 border border-slate-200 shadow-sm">{i + 1}</div>
                          <div className="flex-1">
                             <Input 
                               value={upload.label} 
                               onChange={(e: any) => updateUpload(upload.id, { label: e.target.value })} 
                               className="bg-transparent border-none p-0 h-auto font-bold text-slate-700 focus:ring-0 text-sm"
                               placeholder="Document Name"
                             />
                          </div>
                          <div className="flex items-center gap-3">
                             <div className="flex items-center gap-2 bg-white px-3 py-1.5 rounded-lg border border-slate-100">
                                <span className="text-[10px] font-bold text-slate-400 uppercase">Required</span>
                                <Switch checked={upload.required} onCheckedChange={(c) => updateUpload(upload.id, { required: c })} />
                             </div>
                             <IconButton icon={Trash2} className="text-slate-300 hover:text-red-500" onClick={() => deleteUpload(upload.id)} />
                          </div>
                       </div>
                    ))}
                    {config.uploads.length === 0 && (
                       <div className="text-center py-10 border-2 border-dashed border-slate-200 rounded-2xl text-slate-400 italic">
                          No uploads required for this task.
                       </div>
                    )}
                 </div>
              </div>

              {/* Confirmation Button Config */}
              <div className="bg-white p-8 rounded-[2rem] shadow-sm border border-slate-200 flex items-center gap-6">
                 <div className="p-4 bg-emerald-50 text-emerald-600 rounded-2xl">
                    <CheckCircle2 size={32} />
                 </div>
                 <div className="flex-1">
                    <label className="text-xs font-bold text-slate-400 uppercase tracking-wider block mb-2">Completion Action Button</label>
                    <Input 
                      value={config.buttonText} 
                      onChange={(e: any) => setConfig({ ...config, buttonText: e.target.value })}
                      className="bg-slate-50 border border-slate-200 font-bold text-slate-800"
                    />
                 </div>
              </div>
           </div>
        </div>

        {/* Right: Settings */}
        <div className="w-80 bg-white border-l border-slate-200 flex flex-col flex-shrink-0 p-6 space-y-8 overflow-y-auto">
           <div className="space-y-4">
              <label className="text-xs font-black text-slate-400 uppercase tracking-widest flex items-center gap-2"><Shield size={14} /> Reviewer</label>
              <div className="p-4 bg-slate-50 rounded-2xl border border-slate-200">
                 <select 
                   value={config.reviewer}
                   onChange={(e) => setConfig({ ...config, reviewer: e.target.value })}
                   className="w-full bg-transparent text-sm font-bold text-slate-700 outline-none"
                 >
                    <option value="none">Auto-Verify</option>
                    <option value="advisor">Scientific Advisor</option>
                    <option value="secretary">Academic Secretary</option>
                    <option value="admin">System Admin</option>
                 </select>
                 <p className="text-[10px] text-slate-400 mt-2 leading-tight">
                    {config.reviewer === 'none' ? 'Task completes immediately upon student submission.' : `Task status becomes "Under Review" until approved by ${config.reviewer}.`}
                 </p>
              </div>
           </div>

           <div className="space-y-4">
              <div className="flex items-center justify-between">
                 <label className="text-xs font-black text-slate-400 uppercase tracking-widest flex items-center gap-2"><Paperclip size={14} /> Downloads</label>
                 <button onClick={addTemplate} className="text-[10px] font-bold text-emerald-600 bg-emerald-50 px-2 py-1 rounded hover:bg-emerald-100">+ Add</button>
              </div>
              <div className="space-y-2">
                 {config.templates.map((t, i) => (
                    <div key={i} className="flex items-center gap-3 p-3 bg-slate-50 border border-slate-200 rounded-xl">
                       <div className="p-1.5 bg-white rounded border border-slate-200"><FileText size={14} className="text-slate-400" /></div>
                       <div className="flex-1 min-w-0">
                          <div className="text-xs font-bold text-slate-700 truncate">{t.name}</div>
                          <div className="text-[10px] text-slate-400">{t.size}</div>
                       </div>
                       <IconButton icon={Trash2} size="sm" onClick={() => setConfig({ ...config, templates: config.templates.filter((_, idx) => idx !== i) })} />
                    </div>
                 ))}
                 {config.templates.length === 0 && <div className="text-xs text-slate-400 italic text-center py-4 border border-dashed border-slate-200 rounded-xl">No templates provided.</div>}
              </div>
           </div>
        </div>
      </div>
    </div>
  );
};
