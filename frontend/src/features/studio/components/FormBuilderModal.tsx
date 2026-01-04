import React, { useState } from 'react';
import { motion, Reorder } from 'framer-motion';
import { 
  X, Plus, Type, Hash, List, Calendar, UploadCloud, 
  Trash2, GripVertical, Settings, CheckSquare, Layers, 
  Info, Globe 
} from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Switch } from '@/components/ui/switch';
import { Badge } from '@/components/ui/badge';
import { Label } from '@/components/ui/label';
import { cn } from '@/lib/utils';
import { FormField, FormFieldType } from '../types';

interface FormBuilderModalProps {
  initialFields?: FormField[];
  onSave: (fields: FormField[]) => void;
  onClose: () => void;
}

// Mock Global Dictionaries
const GLOBAL_DICTIONARIES = [
  { id: 'specialties', label: 'PhD Specialties' },
  { id: 'departments', label: 'University Departments' },
  { id: 'programs', label: 'Educational Programs' },
  { id: 'countries', label: 'Countries' }
];

const FIELD_TYPES: { type: FormFieldType; label: string; icon: any }[] = [
  { type: 'text', label: 'Text Input', icon: Type },
  { type: 'number', label: 'Number', icon: Hash },
  { type: 'select', label: 'Dropdown', icon: List },
  { type: 'date', label: 'Date Picker', icon: Calendar },
  { type: 'boolean', label: 'Checkbox', icon: CheckSquare },
  { type: 'upload', label: 'File Upload', icon: UploadCloud },
  { type: 'note', label: 'Instruction Note', icon: Info },
  { type: 'collection', label: 'Repeating Group', icon: Layers },
];

export const FormBuilderModal: React.FC<FormBuilderModalProps> = ({ initialFields = [], onSave, onClose }) => {
  const [fields, setFields] = useState<FormField[]>(
    initialFields.map(f => ({ ...f, id: f.id || `f_${Math.random().toString(36).substr(2, 9)}` }))
  );
  const [selectedFieldId, setSelectedFieldId] = useState<string | null>(null);

  const selectedField = fields.find(f => f.id === selectedFieldId);

  const addField = (type: FormFieldType) => {
    const newField: FormField = {
      id: `f_${Date.now()}`,
      key: `field_${Date.now()}`,
      type: type,
      label: 'New Field',
      required: false
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
    <div className="fixed inset-0 z-[100] flex items-center justify-center p-4 bg-slate-900/60 backdrop-blur-sm">
      <motion.div 
        initial={{ opacity: 0, scale: 0.95 }}
        animate={{ opacity: 1, scale: 1 }}
        exit={{ opacity: 0, scale: 0.95 }}
        className="bg-slate-50 w-full max-w-6xl h-[90vh] rounded-2xl shadow-2xl overflow-hidden flex flex-col"
      >
        {/* Header */}
        <div className="h-16 bg-white border-b border-slate-200 px-6 flex items-center justify-between flex-shrink-0">
          <div className="flex items-center gap-3">
             <div className="bg-indigo-100 text-indigo-600 p-2 rounded-lg"><Settings size={20} /></div>
             <div>
               <h2 className="font-bold text-slate-800 text-lg leading-none">Form Builder</h2>
               <p className="text-xs text-slate-500 mt-1">Design the data structure for this step</p>
             </div>
          </div>
          <div className="flex items-center gap-2">
            <Button variant="ghost" onClick={onClose}>Cancel</Button>
            <Button onClick={() => onSave(fields)}>Save Form</Button>
          </div>
        </div>

        <div className="flex-1 flex overflow-hidden">
          
          {/* Left: Toolbox */}
          <div className="w-64 bg-white border-r border-slate-200 flex-shrink-0 flex flex-col overflow-y-auto">
             <div className="p-4 border-b border-slate-100">
               <h3 className="text-xs font-bold text-slate-400 uppercase tracking-wider">Field Types</h3>
             </div>
             <div className="p-3 grid grid-cols-1 gap-2">
                {FIELD_TYPES.map(t => (
                  <button
                    key={t.type}
                    onClick={() => addField(t.type)}
                    className="flex items-center gap-3 p-3 rounded-xl border border-slate-200 hover:border-indigo-300 hover:bg-indigo-50 hover:text-indigo-700 transition-all text-slate-600 bg-slate-50 text-left group"
                  >
                    <t.icon size={18} className="text-slate-400 group-hover:text-indigo-500" />
                    <span className="text-sm font-medium">{t.label}</span>
                    <Plus size={14} className="ml-auto opacity-0 group-hover:opacity-100 transition-opacity" />
                  </button>
                ))}
             </div>
          </div>

          {/* Center: Canvas */}
          <div className="flex-1 bg-slate-100/50 overflow-y-auto p-8 flex flex-col items-center">
             <div className="w-full max-w-2xl">
               {fields.length === 0 ? (
                 <div className="border-2 border-dashed border-slate-300 rounded-2xl p-12 text-center text-slate-400 flex flex-col items-center">
                    <List size={48} className="opacity-20 mb-4" />
                    <p>Your form is empty.</p>
                    <p className="text-sm">Click items on the left to add fields.</p>
                 </div>
               ) : (
                 <Reorder.Group axis="y" values={fields} onReorder={setFields} className="space-y-3">
                    {fields.map(field => (
                      <Reorder.Item key={field.id} value={field}>
                        <div 
                          onClick={() => setSelectedFieldId(field.id)}
                          className={cn(
                            "bg-white p-4 rounded-xl border-2 transition-all cursor-pointer relative group shadow-sm flex items-start gap-4",
                            selectedFieldId === field.id ? "border-indigo-500 ring-4 ring-indigo-500/10 z-10" : "border-slate-200 hover:border-indigo-300"
                          )}
                        >
                           <div className="mt-1 text-slate-300 cursor-move hover:text-slate-500"><GripVertical size={20} /></div>
                           
                           <div className="flex-1 space-y-1 pointer-events-none">
                              <div className="flex items-center justify-between">
                                 <Label className="text-sm font-bold text-slate-700">
                                   {field.label} {field.required && <span className="text-red-500">*</span>}
                                 </Label>
                                 <Badge variant="secondary" className="text-[10px] uppercase">{field.type}</Badge>
                              </div>
                              
                              {/* Preview Representation */}
                              {field.type === 'text' && <div className="h-9 bg-slate-50 border border-slate-200 rounded-lg w-full" />}
                              {field.type === 'select' && (
                                <div className="h-9 bg-slate-50 border border-slate-200 rounded-lg w-full flex items-center justify-between px-3 text-slate-400 text-xs">
                                  <span>Select option...</span> <List size={14} />
                                </div>
                              )}
                              {field.type === 'collection' && (
                                <div className="p-3 bg-slate-50 border border-slate-200 border-dashed rounded-lg text-xs text-slate-500 text-center">
                                  Repeating Group Content
                                </div>
                              )}
                              {field.type === 'note' && (
                                <div className="text-xs text-slate-500 italic bg-blue-50 p-2 rounded border border-blue-100">
                                  Instructional text will appear here.
                                </div>
                              )}
                              
                              {field.help_text && <p className="text-xs text-slate-400">{field.help_text}</p>}
                           </div>

                           <button 
                             onClick={(e) => { e.stopPropagation(); removeField(field.id); }}
                             className="opacity-0 group-hover:opacity-100 p-2 text-slate-400 hover:text-red-500 hover:bg-red-50 rounded-lg transition-all"
                           >
                             <Trash2 size={18} />
                           </button>
                        </div>
                      </Reorder.Item>
                    ))}
                 </Reorder.Group>
               )}
             </div>
          </div>

          {/* Right: Inspector */}
          <div className="w-80 bg-white border-l border-slate-200 flex-shrink-0 flex flex-col overflow-y-auto">
             {selectedField ? (
               <div className="p-6 space-y-6">
                  <div>
                    <h3 className="text-sm font-bold text-slate-900 mb-1">Field Properties</h3>
                    <p className="text-xs text-slate-500">Editing {selectedField.type} field</p>
                  </div>

                  <div className="space-y-4">
                    <div className="space-y-1.5">
                       <Label className="text-xs font-bold text-slate-500 uppercase">Label</Label>
                       <Input value={selectedField.label} onChange={(e) => updateField(selectedField.id, { label: e.target.value })} />
                    </div>
                    
                    <div className="space-y-1.5">
                       <Label className="text-xs font-bold text-slate-500 uppercase">System Key</Label>
                       <Input value={selectedField.key} onChange={(e) => updateField(selectedField.id, { key: e.target.value })} className="font-mono text-xs" />
                    </div>

                    <div className="space-y-1.5">
                       <Label className="text-xs font-bold text-slate-500 uppercase">Help Text</Label>
                       <Input value={selectedField.help_text || ''} onChange={(e) => updateField(selectedField.id, { help_text: e.target.value })} placeholder="Optional instructions" />
                    </div>

                    <div className="flex items-center justify-between pt-2">
                       <span className="text-sm font-medium text-slate-700">Required Field</span>
                       <Switch checked={!!selectedField.required} onCheckedChange={(c) => updateField(selectedField.id, { required: c })} />
                    </div>

                    {/* Type Specific Options */}
                    {selectedField.type === 'select' && (
                      <div className="pt-4 border-t border-slate-100 space-y-4">
                         <div className="flex items-center justify-between">
                            <Label className="text-xs font-bold text-slate-500 uppercase flex items-center gap-1"><Globe size={12} /> Use Dictionary</Label>
                            <Switch checked={!!selectedField.dictionaryKey} onCheckedChange={(c) => updateField(selectedField.id, { dictionaryKey: c ? 'specialties' : undefined })} />
                         </div>

                         {selectedField.dictionaryKey ? (
                           <div className="space-y-1.5">
                              <Label className="text-xs font-bold text-slate-500 uppercase">Dictionary Source</Label>
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
                              <Label className="text-xs font-bold text-slate-500 uppercase">Manual Options</Label>
                              <div className="space-y-2">
                                {(selectedField.options || []).map((opt, idx) => (
                                  <div key={idx} className="flex gap-2">
                                    <Input 
                                      value={opt.label} 
                                      onChange={(e) => {
                                        const newOpts = [...(selectedField.options || [])];
                                        newOpts[idx] = { ...opt, label: e.target.value, value: e.target.value };
                                        updateField(selectedField.id, { options: newOpts });
                                      }}
                                      className="flex-1 text-xs h-8"
                                    />
                                    <button onClick={() => {
                                       const newOpts = selectedField.options?.filter((_, i) => i !== idx);
                                       updateField(selectedField.id, { options: newOpts });
                                    }} className="text-slate-400 hover:text-red-500"><X size={14} /></button>
                                  </div>
                                ))}
                                <Button size="sm" variant="outline" className="w-full border-dashed" onClick={() => updateField(selectedField.id, { options: [...(selectedField.options || []), { label: 'New Option', value: 'new_option' }] })}>
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
                 <p className="text-sm">Select a field on the canvas to edit properties.</p>
               </div>
             )}
          </div>

        </div>
      </motion.div>
    </div>
  );
};
