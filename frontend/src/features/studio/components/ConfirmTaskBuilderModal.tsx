
import React, { useState } from 'react';
import { 
  X, Save, Trash2, UploadCloud, FileText, CheckCircle2, Shield, Paperclip, CheckSquare
} from 'lucide-react';
import { Dialog, DialogContent } from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Switch } from "@/components/ui/switch";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";

interface UploadRequirement {
  id: string;
  label: string;
  required: boolean;
}

interface ConfirmConfig {
  intro: string;
  buttonText: string;
  reviewer: string;
  templates: { name: string; size: string }[];
  uploads: UploadRequirement[];
}

interface ConfirmTaskBuilderModalProps {
  isOpen: boolean;
  onClose: () => void;
  initialConfig?: ConfirmConfig;
  onSave: (config: ConfirmConfig) => void;
}

const DEFAULT_CONFIG: ConfirmConfig = {
  intro: 'Please upload the required documents.',
  buttonText: 'Confirm Submission',
  reviewer: 'advisor',
  templates: [],
  uploads: [{ id: 'u1', label: 'Document.pdf', required: true }]
};

export const ConfirmTaskBuilderModal: React.FC<ConfirmTaskBuilderModalProps> = ({ isOpen, onClose, initialConfig, onSave }) => {
  const [config, setConfig] = useState<ConfirmConfig>(initialConfig || DEFAULT_CONFIG);

  const updateUpload = (id: string, updates: Partial<UploadRequirement>) => {
    setConfig({ ...config, uploads: config.uploads.map(u => u.id === id ? { ...u, ...updates } : u) });
  };

  return (
    <Dialog open={isOpen} onOpenChange={onClose}>
      <DialogContent className="max-w-[1000px] h-[80vh] flex flex-col p-0 gap-0 bg-slate-50">
         <div className="h-16 border-b border-slate-200 bg-white px-6 flex items-center justify-between shrink-0">
             <div className="flex items-center gap-4">
                 <div className="p-2 bg-emerald-100 text-emerald-600 rounded-lg"><CheckCircle2 size={20} /></div>
                 <h2 className="font-bold text-lg text-slate-900">Confirm Task Config</h2>
             </div>
             <Button onClick={() => { onSave(config); onClose(); }} className="bg-emerald-600 hover:bg-emerald-700">
                 <Save size={16} className="mr-2"/> Save
             </Button>
         </div>

         <div className="flex-1 overflow-y-auto p-8">
            <div className="max-w-3xl mx-auto space-y-8">
                
                {/* Intro */}
                <div className="bg-white p-6 rounded-2xl shadow-sm border border-slate-200">
                    <label className="text-xs font-bold text-slate-400 uppercase mb-2 block">Student Instructions</label>
                    <Textarea 
                        value={config.intro} 
                        onChange={(e) => setConfig({...config, intro: e.target.value})}
                        className="min-h-[100px]"
                    />
                </div>

                {/* Uploads */}
                <div className="bg-white p-6 rounded-2xl shadow-sm border border-slate-200">
                    <div className="flex justify-between items-center mb-4">
                        <label className="text-xs font-bold text-slate-400 uppercase block">Required Submissions</label>
                        <Button size="sm" variant="outline" onClick={() => setConfig({...config, uploads: [...config.uploads, { id: Date.now().toString(), label: '', required: true }]})}>Add Upload Slot</Button>
                    </div>
                    <div className="space-y-3">
                        {config.uploads.map((u, i) => (
                            <div key={u.id} className="flex items-center gap-3 p-3 bg-slate-50 border border-slate-200 rounded-xl">
                                <div className="w-6 h-6 rounded-full bg-white flex items-center justify-center font-bold text-xs border border-slate-200">{i+1}</div>
                                <Input value={u.label} onChange={(e) => updateUpload(u.id, { label: e.target.value })} placeholder="Document Name" className="flex-1" />
                                <div className="flex items-center gap-2">
                                    <span className="text-[10px] uppercase font-bold text-slate-400">Required</span>
                                    <Switch checked={u.required} onCheckedChange={(c) => updateUpload(u.id, { required: c })} />
                                </div>
                                <Button variant="ghost" size="icon" onClick={() => setConfig({...config, uploads: config.uploads.filter(up => up.id !== u.id)})}><Trash2 size={16} className="text-slate-400"/></Button>
                            </div>
                        ))}
                    </div>
                </div>

                {/* Settings */}
                <div className="grid grid-cols-2 gap-6">
                    <div className="bg-white p-6 rounded-2xl shadow-sm border border-slate-200">
                        <label className="text-xs font-bold text-slate-400 uppercase mb-2 block flex items-center gap-2"><Shield size={14}/> Reviewer Role</label>
                        <Select value={config.reviewer} onValueChange={(v) => setConfig({...config, reviewer: v})}>
                            <SelectTrigger>
                                <SelectValue />
                            </SelectTrigger>
                            <SelectContent>
                                <SelectItem value="none">Auto-Verify</SelectItem>
                                <SelectItem value="advisor">Advisor</SelectItem>
                                <SelectItem value="secretary">Secretary</SelectItem>
                                <SelectItem value="admin">Admin</SelectItem>
                            </SelectContent>
                        </Select>
                    </div>
                    <div className="bg-white p-6 rounded-2xl shadow-sm border border-slate-200">
                         <label className="text-xs font-bold text-slate-400 uppercase mb-2 block flex items-center gap-2"><CheckCircle2 size={14}/> Action Button</label>
                         <Input value={config.buttonText} onChange={(e) => setConfig({...config, buttonText: e.target.value})} />
                    </div>
                </div>

            </div>
         </div>
      </DialogContent>
    </Dialog>
  );
};
