
import React, { useState } from 'react';
import { Reorder, motion } from 'framer-motion';
import { 
  ArrowLeft, Plus, Save, Settings, Trash2, 
  GripVertical, CheckSquare, FileText, X,
  Type, Hash, List, Calendar, UploadCloud, Info, Layers, Globe
} from 'lucide-react';
import { Button, Input, Switch, Badge, IconButton, AvatarGroup } from '@/features/admin/components/AdminUI';
import { cn } from '@/lib/utils';
import { useNavigate, useParams } from 'react-router-dom';

const ACTIVE_DESIGNERS = [
  { initials: 'AR', color: 'bg-indigo-600' },
  { initials: 'JS', color: 'bg-emerald-500' },
];

// Mock Global Dictionaries
const GLOBAL_DICTIONARIES = [
  { id: 'specialties', label: 'PhD Specialties' },
  { id: 'departments', label: 'University Departments' },
  { id: 'programs', label: 'Educational Programs' },
  { id: 'countries', label: 'Countries' }
];

interface FormField {
  id: string;
  key: string;
  type: 'text' | 'number' | 'select' | 'date' | 'boolean' | 'upload' | 'collection' | 'note';
  label: string;
  required?: boolean;
  placeholder?: string;
  help_text?: string;
  options?: { value: string; label: string }[]; 
  dictionaryKey?: string; 
}

const FIELD_TYPES = [
  { type: 'text', label: 'Text Input', icon: Type },
  { type: 'number', label: 'Number', icon: Hash },
  { type: 'select', label: 'Dropdown', icon: List },
  { type: 'date', label: 'Date Picker', icon: Calendar },
  { type: 'boolean', label: 'Checkbox', icon: CheckSquare },
  { type: 'upload', label: 'File Upload', icon: UploadCloud },
  { type: 'note', label: 'Instruction', icon: Info },
  { type: 'collection', label: 'Repeater', icon: Layers },
];

export const FormBuilder: React.FC = () => {
  const navigate = useNavigate();
  const { programId, nodeId } = useParams();
  
  const [fields, setFields] = useState<FormField[]>([
    { id: 'f1', key: 'full_name', type: 'text', label: 'Full Name', required: true, placeholder: 'Enter name' }
  ]);
  const [selectedFieldId, setSelectedFieldId] = useState<string | null>(null);
  const [title, setTitle] = useState('Student Profile Form');

  const selectedField = fields.find(f => f.id === selectedFieldId);

  const onNavigate = (path: string) => navigate(path);

  const addField = (type: string) => {
    const newField: FormField = {
      id: `f_${Date.now()}`,
      key: `field_${Date.now()}`,
      type: type as any,
      label: 'New Field',
      required: false,
      options: type === 'select' ? [{ value: 'opt1', label: 'Option 1' }] : undefined
    };
    setFields([...fields, newField]);
    setSelectedFieldId(newField.id);
  };

  const updateField = (id: string, updates: Partial<FormField>) => {
    setFields(fields.map(f => f.id === id ? { ...f, ...updates } : f));
  };

  const removeField = (id: string) => {
    setFields(fields.filter(f => f.id !== id));
    if (selectedFieldId === id) setSelectedFieldId(null);
  };

  return (
    <div className="flex flex-col h-[calc(100vh-4rem)] bg-slate-50 font-sans overflow-hidden">
      {/* Header */}
      <div className="h-20 bg-white border-b border-slate-200 px-8 flex items-center justify-between flex-shrink-0 z-30 shadow-sm">
        <div className="flex items-center gap-6">
          <IconButton icon={ArrowLeft} onClick={() => onNavigate(`/admin/studio/programs/${programId}/builder`)} />
          <div>
             <div className="flex items-center gap-2 mb-1">
                <span className="text-[9px] font-black uppercase text-blue-600 bg-blue-50 px-2 py-0.5 rounded-full tracking-widest border border-blue-100">Form Studio</span>
             </div>
             <Input 
               value={title} 
               onChange={(e: any) => setTitle(e.target.value)}
               className="font-black text-slate-900 text-xl border-none p-0 h-auto focus:ring-0 w-96 bg-transparent"
             />
          </div>
        </div>
        <div className="flex items-center gap-6">
          <AvatarGroup users={ACTIVE_DESIGNERS} />
          <Button variant="primary" icon={Save} onClick={() => alert('Saved!')}>Publish Form</Button>
        </div>
      </div>

      <div className="flex-1 flex overflow-hidden">
        {/* Toolbox */}
        <div className="w-64 bg-white border-r border-slate-200 flex-col flex flex-shrink-0 overflow-y-auto">
           <div className="p-5 border-b border-slate-100"><h3 className="text-xs font-black text-slate-400 uppercase tracking-widest">Inputs</h3></div>
           <div className="p-3 grid grid-cols-1 gap-2">
              {FIELD_TYPES.map(t => (
                <button key={t.type} onClick={() => addField(t.type)} className="flex items-center gap-3 p-3 rounded-xl border border-slate-100 hover:border-blue-200 hover:bg-blue-50 hover:text-blue-700 transition-all text-slate-600 bg-white text-left group shadow-sm">
                  <t.icon size={16} className="text-slate-400 group-hover:text-blue-500" />
                  <span className="text-sm font-bold">{t.label}</span>
                  <Plus size={14} className="ml-auto opacity-0 group-hover:opacity-100" />
                </button>
              ))}
           </div>
        </div>

        {/* Canvas */}
        <div className="flex-1 bg-slate-100/50 overflow-y-auto p-8 flex flex-col items-center">
           <div className="w-full max-w-2xl space-y-4">
             {fields.length === 0 ? (
               <div className="border-2 border-dashed border-slate-300 rounded-2xl p-12 text-center text-slate-400">
                  <p>Your form is empty. Add fields from the left.</p>
               </div>
             ) : (
               <Reorder.Group axis="y" values={fields} onReorder={setFields} className="space-y-3">
                  {fields.map(field => (
                    <Reorder.Item key={field.id} value={field}>
                      <motion.div 
                        layout
                        onClick={() => setSelectedFieldId(field.id)}
                        className={cn(
                          "bg-white p-5 rounded-2xl border-2 transition-all cursor-pointer relative group flex items-start gap-4 shadow-sm",
                          selectedFieldId === field.id ? "border-blue-500 ring-4 ring-blue-500/10 z-10" : "border-slate-200 hover:border-blue-200"
                        )}
                      >
                         <div className="mt-1 text-slate-300 cursor-move"><GripVertical size={20} /></div>
                         <div className="flex-1 space-y-1 pointer-events-none">
                            <div className="flex items-center justify-between">
                               <label className="text-sm font-bold text-slate-800">{field.label} {field.required && <span className="text-red-500">*</span>}</label>
                               <Badge variant="secondary" className="text-[9px] uppercase">{field.type}</Badge>
                            </div>
                            
                            {/* Visual Preview based on type */}
                            {field.type === 'text' && <div className="h-10 bg-slate-50 border border-slate-200 rounded-xl w-full" />}
                            
                            {field.type === 'select' && (
                                <div className="h-10 bg-slate-50 border border-slate-200 rounded-xl w-full flex items-center justify-between px-3 text-slate-400">
                                    <span className="text-xs">{field.dictionaryKey ? `[Dictionary: ${field.dictionaryKey}]` : 'Select option...'}</span>
                                    <List size={14}/>
                                </div>
                            )}
                            
                            {field.type === 'upload' && (
                                <div className="h-12 bg-slate-50 border border-slate-200 border-dashed rounded-xl w-full flex items-center justify-center text-slate-400">
                                    <UploadCloud size={16} className="mr-2" /> Upload File
                                </div>
                            )}

                            {field.help_text && <p className="text-xs text-slate-400">{field.help_text}</p>}
                         </div>
                         <button onClick={(e) => { e.stopPropagation(); removeField(field.id); }} className="opacity-0 group-hover:opacity-100 p-2 text-slate-400 hover:text-red-500"><Trash2 size={18} /></button>
                      </motion.div>
                    </Reorder.Item>
                  ))}
               </Reorder.Group>
             )}
           </div>
        </div>

        {/* Inspector */}
        <div className="w-80 bg-white border-l border-slate-200 flex-shrink-0 flex flex-col overflow-y-auto">
           {selectedField ? (
             <div className="p-6 space-y-6">
                <div className="border-b border-slate-100 pb-4">
                  <h3 className="text-sm font-black text-slate-900">Field Properties</h3>
                  <p className="text-xs text-slate-500">Configuring {selectedField.type}</p>
                </div>
                <div className="space-y-4">
                  <div className="space-y-1">
                     <label className="text-xs font-bold text-slate-500">Label</label>
                     <Input value={selectedField.label} onChange={(e: any) => updateField(selectedField.id, { label: e.target.value })} />
                  </div>
                  <div className="space-y-1">
                     <label className="text-xs font-bold text-slate-500">Key Name</label>
                     <Input value={selectedField.key} onChange={(e: any) => updateField(selectedField.id, { key: e.target.value })} className="font-mono text-xs" />
                  </div>
                  <div className="space-y-1">
                     <label className="text-xs font-bold text-slate-500">Helper Text</label>
                     <Input value={selectedField.help_text || ''} onChange={(e: any) => updateField(selectedField.id, { help_text: e.target.value })} />
                  </div>
                  <div className="flex items-center justify-between pt-2">
                     <span className="text-sm font-medium text-slate-700">Required</span>
                     <Switch checked={!!selectedField.required} onCheckedChange={(c) => updateField(selectedField.id, { required: c })} />
                  </div>

                  {/* SELECT OPTIONS CONFIGURATION */}
                  {selectedField.type === 'select' && (
                      <div className="pt-4 border-t border-slate-100 space-y-4 animate-in fade-in">
                         <div className="flex items-center justify-between">
                            <label className="text-xs font-bold text-slate-500 uppercase flex items-center gap-1"><Globe size={12} /> Use Dictionary</label>
                            <Switch checked={!!selectedField.dictionaryKey} onCheckedChange={(c) => updateField(selectedField.id, { dictionaryKey: c ? 'specialties' : undefined })} />
                         </div>

                         {selectedField.dictionaryKey ? (
                           <div className="space-y-1.5">
                              <label className="text-xs font-bold text-slate-500 uppercase">Dictionary Source</label>
                              <select 
                                className="w-full p-2 bg-white border border-slate-200 rounded-lg text-sm"
                                value={selectedField.dictionaryKey}
                                onChange={(e) => updateField(selectedField.id, { dictionaryKey: e.target.value })}
                              >
                                {GLOBAL_DICTIONARIES.map(d => (
                                  <option key={d.id} value={d.id}>{d.label}</option>
                                ))}
                              </select>
                           </div>
                         ) : (
                           <div className="space-y-2">
                              <label className="text-xs font-bold text-slate-500 uppercase">Manual Options</label>
                              <div className="space-y-2">
                                {(selectedField.options || []).map((opt, idx) => (
                                  <div key={idx} className="flex gap-2">
                                    <Input 
                                      value={opt.label} 
                                      onChange={(e: any) => {
                                        const newOpts = [...(selectedField.options || [])];
                                        newOpts[idx] = { ...opt, label: e.target.value, value: e.target.value };
                                        updateField(selectedField.id, { options: newOpts });
                                      }}
                                      className="flex-1 text-xs h-8"
                                      placeholder="Option Label"
                                    />
                                    <button onClick={() => {
                                       const newOpts = selectedField.options?.filter((_, i) => i !== idx);
                                       updateField(selectedField.id, { options: newOpts });
                                    }} className="text-slate-300 hover:text-red-500"><X size={14} /></button>
                                  </div>
                                ))}
                                <Button size="sm" variant="outline" className="w-full border-dashed" onClick={() => updateField(selectedField.id, { options: [...(selectedField.options || []), { label: '', value: '' }] })}>
                                  + Add Option
                                </Button>
                              </div>
                           </div>
                         )}
                      </div>
                  )}
                </div>
             </div>
           ) : (
             <div className="flex-1 flex flex-col items-center justify-center text-slate-400 p-8 text-center">
               <Settings size={32} className="opacity-20 mb-2" />
               <p className="text-sm">Select a field to configure.</p>
             </div>
           )}
        </div>
      </div>
    </div>
  );
};
