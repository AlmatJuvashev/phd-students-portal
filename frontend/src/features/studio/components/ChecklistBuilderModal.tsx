import React, { useState } from 'react';
import { motion, Reorder } from 'framer-motion';
import { 
  Plus, Save, Settings, Trash2, GripVertical, CheckSquare, 
  FileText, Shield, Paperclip
} from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Badge } from '@/components/ui/badge';
import { Switch } from '@/components/ui/switch';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { ScrollArea } from '@/components/ui/scroll-area';
import { Separator } from '@/components/ui/separator';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { cn } from '@/lib/utils';
import { ChecklistItem, ChecklistConfig } from '../types';

interface ChecklistBuilderModalProps {
  initialConfig?: ChecklistConfig;
  onSave: (config: ChecklistConfig) => void;
  onClose: () => void;
}

export const ChecklistBuilderModal: React.FC<ChecklistBuilderModalProps> = ({ initialConfig, onSave, onClose }) => {
  const [items, setItems] = useState<ChecklistItem[]>(initialConfig?.items || [
    { id: '1', text: 'New Item', required: true }
  ]);
  const [activeItemId, setActiveItemId] = useState<string | null>(items.length > 0 ? items[0].id : null);
  const [config, setConfig] = useState<Omit<ChecklistConfig, 'items'>>({
    intro: initialConfig?.intro || '',
    reviewer_role: initialConfig?.reviewer_role || 'secretary',
    templates: initialConfig?.templates || []
  });

  const addItem = () => {
    const newItem = { id: `item_${Date.now()}`, text: 'New Item', required: true };
    setItems([...items, newItem]);
    setActiveItemId(newItem.id);
  };

  const updateItem = (id: string, updates: Partial<ChecklistItem>) => {
    setItems(items.map(i => i.id === id ? { ...i, ...updates } : i));
  };

  const deleteItem = (id: string) => {
    setItems(items.filter(i => i.id !== id));
    if (activeItemId === id) setActiveItemId(null);
  };

  const addTemplate = () => {
    // Placeholder for actual template logic
    setConfig({ ...config, templates: [...(config.templates || []), { name: 'Guideline.pdf', url: '#', size: '1.2 MB' }] });
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
             <div className="bg-orange-100 text-orange-600 p-2 rounded-lg"><CheckSquare size={20} /></div>
             <div>
               <h2 className="font-bold text-slate-800 text-lg leading-none">Checklist Builder</h2>
               <p className="text-xs text-slate-500 mt-1">Design the tasks required for this step</p>
             </div>
          </div>
          <div className="flex items-center gap-2">
            <Button variant="ghost" onClick={onClose}>Cancel</Button>
            <Button onClick={() => onSave({ ...config, items })}>Save Checklist</Button>
          </div>
        </div>

        <div className="flex-1 flex overflow-hidden">
          {/* Left: Items List */}
          <div className="w-80 bg-white border-r border-slate-200 flex flex-col flex-shrink-0">
             <div className="p-4 border-b border-slate-100 flex justify-between items-center">
                <h3 className="text-xs font-black text-slate-400 uppercase tracking-widest">Items</h3>
                <Button size="sm" variant="ghost" onClick={addItem}><Plus size={16} /></Button>
             </div>
             <ScrollArea className="flex-1 p-3">
                <Reorder.Group axis="y" values={items} onReorder={setItems} className="space-y-2">
                   {items.map(item => (
                     <Reorder.Item key={item.id} value={item}>
                        <div 
                          onClick={() => setActiveItemId(item.id)}
                          className={cn(
                            "p-3 rounded-xl border flex items-center gap-3 cursor-pointer transition-all group",
                            activeItemId === item.id ? "bg-orange-50 border-orange-200 shadow-sm" : "bg-white border-slate-200 hover:border-orange-200"
                          )}
                        >
                           <GripVertical size={16} className="text-slate-300" />
                           <div className={cn("flex-1 text-sm font-bold truncate", activeItemId === item.id ? "text-orange-700" : "text-slate-700")}>{item.text}</div>
                           <button className="opacity-0 group-hover:opacity-100 p-1 hover:bg-slate-100 rounded" onClick={(e) => { e.stopPropagation(); deleteItem(item.id); }}>
                             <Trash2 size={14} className="text-slate-400 hover:text-red-500" />
                           </button>
                        </div>
                     </Reorder.Item>
                   ))}
                </Reorder.Group>
             </ScrollArea>
          </div>

          {/* Center: Canvas / Editor */}
          <div className="flex-1 bg-slate-100/50 overflow-y-auto p-8 flex flex-col items-center">
             <div className="w-full max-w-2xl space-y-6">
                <div className="bg-white p-8 rounded-[2rem] shadow-sm border border-slate-200 space-y-6">
                   <div>
                      <h3 className="text-lg font-black text-slate-900 mb-2">Introduction</h3>
                      <Textarea 
                        value={config.intro}
                        onChange={(e) => setConfig({ ...config, intro: e.target.value })}
                        className="min-h-[100px] resize-none"
                        placeholder="Instructions for the student..."
                      />
                   </div>

                   <Separator />

                   {/* Item Editor */ }
                   {activeItemId ? (
                      <div className="p-6 bg-slate-50 rounded-2xl border border-slate-200 animate-in fade-in slide-in-from-bottom-2">
                         <h4 className="text-xs font-black text-slate-400 uppercase tracking-widest mb-4">Editing Item</h4>
                         <div className="space-y-4">
                            <div className="space-y-1">
                               <Label className="text-xs font-bold text-slate-500">Item Text</Label>
                               <Input 
                                 value={items.find(i => i.id === activeItemId)?.text} 
                                 onChange={(e) => updateItem(activeItemId, { text: e.target.value })} 
                                 className="bg-white"
                               />
                            </div>
                            <div className="space-y-1">
                               <Label className="text-xs font-bold text-slate-500">Help Text / Tooltip</Label>
                               <Input 
                                 value={items.find(i => i.id === activeItemId)?.helpText || ''} 
                                 onChange={(e) => updateItem(activeItemId, { helpText: e.target.value })} 
                                 className="bg-white"
                                 placeholder="Optional"
                               />
                            </div>
                            <div className="flex items-center justify-between pt-2">
                               <span className="text-sm font-bold text-slate-700">Mandatory</span>
                               <Switch 
                                 checked={!!items.find(i => i.id === activeItemId)?.required} 
                                 onCheckedChange={(c) => updateItem(activeItemId, { required: c })} 
                               />
                            </div>
                         </div>
                      </div>
                   ) : (
                      <div className="text-center py-10 text-slate-400 text-sm">Select an item from the left to edit.</div>
                   )}
                </div>
             </div>
          </div>

          {/* Right: Settings */}
          <div className="w-80 bg-white border-l border-slate-200 flex flex-col flex-shrink-0 p-6 space-y-8 overflow-y-auto">
             <div className="space-y-4">
                <Label className="text-xs font-black text-slate-400 uppercase tracking-widest flex items-center gap-2"><Shield size={14} /> Reviewer</Label>
                <Select 
                   value={config.reviewer_role} 
                   onValueChange={(v: any) => setConfig({ ...config, reviewer_role: v })}
                >
                  <SelectTrigger>
                    <SelectValue placeholder="Select Reviewer" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="none">Auto-Approve</SelectItem>
                    <SelectItem value="advisor">Scientific Advisor</SelectItem>
                    <SelectItem value="secretary">Academic Secretary</SelectItem>
                    <SelectItem value="admin">Admin</SelectItem>
                  </SelectContent>
                </Select>
             </div>

             <Separator />

             <div className="space-y-4">
                <div className="flex items-center justify-between">
                   <Label className="text-xs font-black text-slate-400 uppercase tracking-widest flex items-center gap-2"><Paperclip size={14} /> Templates</Label>
                   <Button size="sm" variant="ghost" className="h-6 w-6 p-0" onClick={addTemplate}><Plus size={14} /></Button>
                </div>
                <div className="space-y-2">
                   {(config.templates || []).map((t, i) => (
                      <div key={i} className="flex items-center gap-3 p-3 bg-slate-50 border border-slate-200 rounded-xl">
                         <div className="p-1.5 bg-white rounded border border-slate-200"><FileText size={14} className="text-slate-400" /></div>
                         <div className="flex-1 min-w-0">
                            <div className="text-xs font-bold text-slate-700 truncate">{t.name}</div>
                            <div className="text-[10px] text-slate-400">{t.size}</div>
                         </div>
                         <button onClick={() => setConfig({ ...config, templates: config.templates?.filter((_, idx) => idx !== i) })} className="text-slate-400 hover:text-red-500">
                            <Trash2 size={14} />
                         </button>
                      </div>
                   ))}
                   {(!config.templates || config.templates.length === 0) && <div className="text-xs text-slate-400 italic text-center py-4 border border-dashed border-slate-200 rounded-xl">No templates added.</div>}
                </div>
             </div>
          </div>
        </div>
      </motion.div>
    </div>
  );
};
