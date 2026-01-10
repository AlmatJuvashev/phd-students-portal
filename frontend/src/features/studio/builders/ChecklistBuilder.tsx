
import React, { useState } from 'react';
import { Reorder, motion, AnimatePresence } from 'framer-motion';
import { 
  ArrowLeft, Plus, Save, Settings, Trash2, GripVertical, CheckSquare, 
  FileText, Download, User, Info, Paperclip, Shield
} from 'lucide-react';
import { Button, Input, Switch, Badge, IconButton, AvatarGroup } from '@/features/admin/components/AdminUI';
import { cn } from '@/lib/utils';
import { useNavigate, useParams } from 'react-router-dom';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { getProgramVersionNodes, updateProgramVersionNode } from '@/features/curriculum/api';
import { api } from '@/api/client';
import { toast } from 'sonner';
import { useEffect } from 'react';

const ACTIVE_DESIGNERS = [
  { initials: 'AR', color: 'bg-indigo-600' },
];

interface ChecklistItem {
  id: string;
  text: string;
  required: boolean;
  helpText?: string;
}

interface ChecklistConfig {
  title: string;
  intro: string;
  reviewer: string;
  templates: { name: string; size: string }[];
}

export const ChecklistBuilder: React.FC = () => {
  const navigate = useNavigate();
  const { programId, nodeId } = useParams();
  const queryClient = useQueryClient();

  const [items, setItems] = useState<ChecklistItem[]>([]);
  const [activeItemId, setActiveItemId] = useState<string | null>(null);
  const [config, setConfig] = useState<ChecklistConfig>({
    title: 'Checklist Designer',
    intro: '',
    reviewer: 'advisor',
    templates: []
  });

  // --- API Connectivity ---
  const { data: nodesData, isLoading } = useQuery({
    queryKey: ['programNodes', programId],
    queryFn: () => getProgramVersionNodes(programId!),
    enabled: !!programId
  });

  const updateNodeMutation = useMutation({
    mutationFn: (config: any) => updateProgramVersionNode(programId!, nodeId!, { config: JSON.stringify(config) }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['programNodes', programId] });
      toast.success('Checklist published successfully');
    },
    onError: () => toast.error('Failed to publish checklist')
  });

  const createNodeMutation = useMutation({
    mutationFn: (data: any) => api.post(`/curriculum/programs/${programId}/builder/nodes`, data),
    onSuccess: (response: any) => {
      queryClient.invalidateQueries({ queryKey: ['programNodes', programId] });
      toast.success('New checklist step created');
      const newNodeId = response.data.id;
      navigate(`/admin/studio/programs/${programId}/checklist/${newNodeId}/builder`, { replace: true });
    },
    onError: () => toast.error('Failed to create checklist step')
  });

  // Hydrate from API
  useEffect(() => {
    if (nodeId && nodesData && Array.isArray(nodesData)) {
      const node = (nodesData as any[]).find(n => n.id === nodeId);
      if (node) {
        const nodeConfig = typeof node.config === 'string' ? JSON.parse(node.config || '{}') : (node.config || {});
        setItems(nodeConfig.items || []);
        setConfig({
          title: node.title || 'Checklist Designer',
          intro: nodeConfig.intro || '',
          reviewer: nodeConfig.reviewer || 'advisor',
          templates: nodeConfig.templates || []
        });
        if (nodeConfig.items?.length > 0 && !activeItemId) {
          setActiveItemId(nodeConfig.items[0].id);
        }
      }
    } else if (!nodeId && programId) {
      // Creation mode
      setItems([]);
      setConfig({
        title: 'New Checklist Step',
        intro: '',
        reviewer: 'advisor',
        templates: []
      });
    }
  }, [nodesData, nodeId]);

  const handleSave = () => {
    const payload = {
      ...config,
      items
    };
    
    if (nodeId) {
      updateNodeMutation.mutate(payload);
    } else {
      createNodeMutation.mutate({
        title: config.title,
        type: 'checklist',
        config: JSON.stringify(payload),
        module_key: 'I',
        coordinates: { x: 400, y: 300 }
      });
    }
  };

  const addItem = () => {
    const newItem = { id: Date.now().toString(), text: 'New Item', required: true };
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
    setConfig({ ...config, templates: [...config.templates, { name: 'Guideline.pdf', size: '1.2 MB' }] });
  };

  return (
    <div className="flex flex-col h-[calc(100vh-4rem)] bg-slate-50 font-sans overflow-hidden">
      <div className="h-20 bg-white border-b border-slate-200 px-8 flex items-center justify-between flex-shrink-0 z-30 shadow-sm">
        <div className="flex items-center gap-6">
          <IconButton icon={ArrowLeft} onClick={() => navigate(`/admin/studio/programs/${programId}/builder`)} />
          <div>
             <div className="flex items-center gap-2 mb-1">
                <span className="text-[9px] font-black uppercase text-orange-600 bg-orange-50 px-2 py-0.5 rounded-full tracking-widest border border-orange-100">Checklist Studio</span>
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
          <Button variant="primary" icon={Save} onClick={handleSave} disabled={updateNodeMutation.isPending}>
            {updateNodeMutation.isPending ? 'Publishing...' : 'Publish Checklist'}
          </Button>
        </div>
      </div>

      <div className="flex-1 flex overflow-hidden">
        {/* Left: Items List */}
        <div className="w-80 bg-white border-r border-slate-200 flex-col flex flex-shrink-0">
           <div className="p-5 border-b border-slate-100 flex justify-between items-center">
              <h3 className="text-xs font-black text-slate-400 uppercase tracking-widest">Checklist Items</h3>
              <Button size="sm" variant="ghost" icon={Plus} onClick={addItem}>Add</Button>
           </div>
           <div className="flex-1 overflow-y-auto p-3">
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
                         <IconButton icon={Trash2} size="sm" className="opacity-0 group-hover:opacity-100" onClick={(e: any) => { e.stopPropagation(); deleteItem(item.id); }} />
                      </div>
                   </Reorder.Item>
                 ))}
              </Reorder.Group>
           </div>
        </div>

        {/* Center: Canvas / Editor */}
        <div className="flex-1 bg-slate-100/50 overflow-y-auto p-8 flex flex-col items-center">
           <div className="w-full max-w-2xl space-y-6">
              <div className="bg-white p-8 rounded-[2rem] shadow-sm border border-slate-200 space-y-6">
                 <div>
                    <h3 className="text-lg font-black text-slate-900 mb-2">Introduction</h3>
                    <textarea 
                      value={config.intro}
                      onChange={(e) => setConfig({ ...config, intro: e.target.value })}
                      className="w-full p-4 bg-slate-50 border-none rounded-xl text-sm min-h-[100px] resize-none focus:ring-2 focus:ring-orange-100"
                      placeholder="Instructions for the student..."
                    />
                 </div>

                 {/* Item Editor */}
                 {activeItemId ? (
                    <div className="p-6 bg-slate-50 rounded-2xl border border-slate-200 animate-in fade-in">
                       <h4 className="text-xs font-black text-slate-400 uppercase tracking-widest mb-4">Editing Item</h4>
                       <div className="space-y-4">
                          <div className="space-y-1">
                             <label className="text-xs font-bold text-slate-500">Item Text</label>
                             <Input 
                               value={items.find(i => i.id === activeItemId)?.text} 
                               onChange={(e: any) => updateItem(activeItemId, { text: e.target.value })} 
                               className="bg-white"
                             />
                          </div>
                          <div className="flex items-center justify-between">
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
              <label className="text-xs font-black text-slate-400 uppercase tracking-widest flex items-center gap-2"><Shield size={14} /> Reviewer</label>
              <select 
                value={config.reviewer}
                onChange={(e) => setConfig({ ...config, reviewer: e.target.value })}
                className="w-full p-3 bg-slate-50 border border-slate-200 rounded-xl text-sm font-bold outline-none"
              >
                 <option value="none">Auto-Approve</option>
                 <option value="advisor">Scientific Advisor</option>
                 <option value="secretary">Academic Secretary</option>
                 <option value="admin">Admin</option>
              </select>
           </div>

           <div className="space-y-4">
              <div className="flex items-center justify-between">
                 <label className="text-xs font-black text-slate-400 uppercase tracking-widest flex items-center gap-2"><Paperclip size={14} /> Templates</label>
                 <button onClick={addTemplate} className="text-[10px] font-bold text-orange-600 bg-orange-50 px-2 py-1 rounded hover:bg-orange-100">+ Add</button>
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
                 {config.templates.length === 0 && <div className="text-xs text-slate-400 italic text-center py-4 border border-dashed border-slate-200 rounded-xl">No templates added.</div>}
              </div>
           </div>
        </div>
      </div>
    </div>
  );
};
