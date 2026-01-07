
import React, { useState } from 'react';
import { 
  X, Save, Sparkles, Settings, Scale, 
  Users, UploadCloud, FileText, Type, Link as LinkIcon, 
  Mic, Trash2, Plus, Wand2, Paperclip, ArrowLeft
} from 'lucide-react';
import { Dialog, DialogContent } from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Switch } from "@/components/ui/switch";
import { Badge } from "@/components/ui/badge";
import { cn } from "@/lib/utils";
import { motion, AnimatePresence } from 'framer-motion';

interface AssignmentBuilderModalProps {
  isOpen: boolean;
  onClose: () => void;
  initialConfig?: AssignmentConfig;
  onSave: (config: AssignmentConfig) => void;
}

export interface RubricLevel {
  id: string;
  points: number;
  label: string;
  description: string;
}

export interface RubricCriterion {
  id: string;
  title: string;
  description: string;
  levels: RubricLevel[];
}

export interface AssignmentConfig {
  title: string;
  description: string;
  submission_types: ('file' | 'text' | 'url' | 'media')[];
  allowed_extensions?: string;
  group_assignment: boolean;
  peer_review: boolean;
  peer_review_count: number;
  rubric: RubricCriterion[];
  points: number;
  instruction_files?: { name: string; size: string; type: string }[];
}

export const AssignmentBuilderModal: React.FC<AssignmentBuilderModalProps> = ({ isOpen, onClose, initialConfig, onSave }) => {
  const [activeTab, setActiveTab] = useState('Configuration');
  const [isAiOpen, setIsAiOpen] = useState(false);
  const [aiPrompt, setAiPrompt] = useState('');
  const [isGenerating, setIsGenerating] = useState(false);

  const [config, setConfig] = useState<AssignmentConfig>(initialConfig || {
    title: 'New Assignment',
    description: '',
    submission_types: ['file'],
    group_assignment: false,
    peer_review: false,
    peer_review_count: 2,
    points: 100,
    rubric: [
      {
        id: 'rc1',
        title: 'Criterion 1',
        description: '',
        levels: [
          { id: 'rl1', points: 10, label: 'Excellent', description: '' },
          { id: 'rl2', points: 5, label: 'Good', description: '' },
          { id: 'rl3', points: 0, label: 'Poor', description: '' }
        ]
      }
    ],
    instruction_files: []
  });

  const updateConfig = (updates: Partial<AssignmentConfig>) => {
    setConfig(prev => ({ ...prev, ...updates }));
  };

  const toggleSubmissionType = (type: 'file' | 'text' | 'url' | 'media') => {
    const types = new Set(config.submission_types);
    if (types.has(type)) types.delete(type);
    else types.add(type);
    updateConfig({ submission_types: Array.from(types) as ('file' | 'text' | 'url' | 'media')[] });
  };

  const handleInstructionUpload = () => {
    // Mock upload
    const newFile = { name: "Instructions.pdf", size: "1.2 MB", type: "application/pdf" };
    updateConfig({ instruction_files: [...(config.instruction_files || []), newFile] });
  };

  const removeInstructionFile = (index: number) => {
    updateConfig({ instruction_files: config.instruction_files?.filter((_, i) => i !== index) });
  };

  // --- Rubric Handlers ---
  const addCriterion = () => {
    const newCrit: RubricCriterion = {
      id: `rc_${Date.now()}`,
      title: 'New Criterion',
      description: '',
      levels: [
        { id: 'rl1', points: 5, label: 'Mastery', description: '' },
        { id: 'rl2', points: 3, label: 'Competent', description: '' },
        { id: 'rl3', points: 1, label: 'Developing', description: '' }
      ]
    };
    updateConfig({ rubric: [...config.rubric, newCrit] });
  };

  const updateCriterion = (id: string, updates: Partial<RubricCriterion>) => {
    updateConfig({ rubric: config.rubric.map(r => r.id === id ? { ...r, ...updates } : r) });
  };

  const deleteCriterion = (id: string) => {
    updateConfig({ rubric: config.rubric.filter(r => r.id !== id) });
  };

  // --- AI Logic (Mock) ---
  const handleAiGenerate = () => {
    if (!aiPrompt) return;
    setIsGenerating(true);
    
    setTimeout(() => {
      let newRubric = [...config.rubric];
      const isEssay = aiPrompt.toLowerCase().includes('essay');
      
      if (isEssay) {
         newRubric = [
            ...newRubric,
            {
               id: `gen_${Date.now()}_1`,
               title: 'Thesis Statement',
               description: 'Clarity and strength of the main argument.',
               levels: [
                  { id: 'l1', points: 10, label: 'Strong', description: 'Clear, arguable, and specific.' },
                  { id: 'l2', points: 5, label: 'Weak', description: 'Vague or obvious.' }
               ]
            }
         ];
      }

      updateConfig({
         title: isEssay ? `Essay: ${aiPrompt}` : config.title,
         description: `AI Generated Draft for: ${aiPrompt}`,
         rubric: newRubric
      });
      
      setIsGenerating(false);
      setIsAiOpen(false);
    }, 1500);
  };

  return (
    <Dialog open={isOpen} onOpenChange={onClose}>
        <DialogContent className="max-w-[95vw] w-[1400px] h-[90vh] flex flex-col p-0 gap-0 bg-slate-50 overflow-hidden">
            {/* Header */}
            <div className="h-16 border-b border-slate-200 bg-white px-6 flex items-center justify-between shrink-0 z-30">
                <div className="flex items-center gap-4">
                    <div className="p-2 bg-indigo-100 text-indigo-600 rounded-lg"><FileText size={20} /></div>
                    <Input 
                        value={config.title} 
                        onChange={(e) => updateConfig({ title: e.target.value })}
                        className="font-black text-lg border-none focus-visible:ring-0 bg-transparent w-96 placeholder:text-slate-300"
                        placeholder="Assignment Title..."
                    />
                </div>
                <div className="flex items-center gap-2">
                    <Button 
                        variant="ghost" 
                        onClick={() => setIsAiOpen(true)}
                        className="text-purple-600 hover:bg-purple-50 hover:text-purple-700"
                    >
                        <Sparkles size={16} className="mr-2"/> AI Assist
                    </Button>
                    <div className="h-4 w-px bg-slate-200 mx-2" />
                    <Button variant="outline" onClick={onClose}>Cancel</Button>
                    <Button onClick={() => { onSave(config); onClose(); }}>
                        <Save size={16} className="mr-2"/> Save Assignment
                    </Button>
                </div>
            </div>

            <div className="flex-1 flex overflow-hidden">
                {/* Sidebar */}
                <div className="w-64 bg-white border-r border-slate-200 flex flex-col shrink-0 z-20">
                    <div className="p-6 space-y-2">
                        <button 
                            onClick={() => setActiveTab('Configuration')}
                            className={cn(
                                "w-full flex items-center gap-3 px-4 py-3 rounded-xl text-sm font-bold transition-all text-left",
                                activeTab === 'Configuration' ? "bg-indigo-50 text-indigo-700" : "text-slate-500 hover:bg-slate-50"
                            )}
                        >
                            <Settings size={18} /> Configuration
                        </button>
                        <button 
                            onClick={() => setActiveTab('Rubric')}
                            className={cn(
                                "w-full flex items-center gap-3 px-4 py-3 rounded-xl text-sm font-bold transition-all text-left",
                                activeTab === 'Rubric' ? "bg-indigo-50 text-indigo-700" : "text-slate-500 hover:bg-slate-50"
                            )}
                        >
                            <Scale size={18} /> Grading Rubric
                        </button>
                    </div>
                    
                     <div className="mt-auto p-6 border-t border-slate-100">
                        <div className="bg-slate-50 rounded-xl p-4 border border-slate-200 text-center">
                            <div className="text-xs font-bold text-slate-400 uppercase tracking-wider mb-2">Total Score</div>
                            <div className="text-3xl font-black text-slate-900">{config.rubric.reduce((acc, c) => acc + Math.max(...c.levels.map(l => l.points)), 0)}</div>
                            <div className="text-[10px] text-slate-400 font-bold mt-1">POINTS</div>
                        </div>
                    </div>
                </div>

                {/* Main Workspace */}
                <div className="flex-1 bg-slate-50/50 overflow-y-auto p-8 relative">
                    <AnimatePresence mode="wait">
                        {activeTab === 'Configuration' && (
                            <motion.div 
                                key="config"
                                initial={{ opacity: 0, y: 10 }}
                                animate={{ opacity: 1, y: 0 }}
                                exit={{ opacity: 0, scale: 0.98 }}
                                className="space-y-8 max-w-4xl mx-auto"
                            >
                                <div className="bg-white p-8 rounded-[2rem] shadow-sm border border-slate-200">
                                    <div className="flex justify-between items-start mb-4">
                                        <label className="text-xs font-black text-slate-400 uppercase tracking-widest flex items-center gap-2">
                                            <FileText size={16} /> Instructions & Prompt
                                        </label>
                                    </div>
                                    <Textarea 
                                        value={config.description}
                                        onChange={(e) => updateConfig({ description: e.target.value })}
                                        className="w-full h-48 p-0 border-none resize-none focus-visible:ring-0 text-lg leading-relaxed text-slate-700 placeholder:text-slate-300 mb-6 shadow-none"
                                        placeholder="# Enter assignment details, objectives, and deliverables..."
                                    />
                                    
                                    <div className="pt-6 border-t border-slate-100">
                                        <div className="flex items-center justify-between mb-3">
                                            <label className="text-xs font-black text-slate-400 uppercase tracking-widest flex items-center gap-2">
                                                <Paperclip size={16} /> Attached Materials
                                            </label>
                                            <Button size="sm" variant="secondary" onClick={handleInstructionUpload}><UploadCloud size={14} className="mr-2"/> Upload File</Button>
                                        </div>

                                        {(config.instruction_files?.length || 0) > 0 && (
                                            <div className="grid grid-cols-1 gap-2">
                                                {config.instruction_files?.map((file, i) => (
                                                    <div key={i} className="flex items-center gap-3 p-3 bg-slate-50 border border-slate-200 rounded-xl shadow-sm">
                                                        <div className="p-2 bg-red-50 text-red-600 rounded-lg border border-red-100"><FileText size={16} /></div>
                                                        <div className="flex-1 min-w-0">
                                                            <div className="text-sm font-bold text-slate-700 truncate">{file.name}</div>
                                                            <div className="text-[10px] text-slate-400 font-bold">{file.size} • {file.type.split('/')[1].toUpperCase()}</div>
                                                        </div>
                                                        <Button variant="ghost" size="icon" className="text-slate-400 hover:text-red-500 hover:bg-red-50" onClick={() => removeInstructionFile(i)}><Trash2 size={16}/></Button>
                                                    </div>
                                                ))}
                                            </div>
                                        )}
                                    </div>
                                </div>

                                <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
                                    <div className="bg-white p-8 rounded-[2rem] shadow-sm border border-slate-200">
                                        <label className="text-xs font-black text-slate-400 uppercase tracking-widest mb-6 flex items-center gap-2">
                                            <UploadCloud size={16} /> Submission Format
                                        </label>
                                        <div className="grid grid-cols-2 gap-4">
                                            {[
                                                { id: 'file', label: 'File Upload', icon: FileText },
                                                { id: 'text', label: 'Text Entry', icon: Type },
                                                { id: 'url', label: 'Website URL', icon: LinkIcon },
                                                { id: 'media', label: 'Media Rec', icon: Mic },
                                            ].map(type => (
                                                <button
                                                    key={type.id}
                                                    onClick={() => toggleSubmissionType(type.id as any)}
                                                    className={cn(
                                                        "flex flex-col items-center justify-center p-4 rounded-2xl border-2 transition-all gap-2 h-28",
                                                        config.submission_types.includes(type.id as any)
                                                            ? "border-indigo-600 bg-indigo-50 text-indigo-700"
                                                            : "border-slate-100 bg-slate-50 text-slate-400 hover:border-slate-200"
                                                    )}
                                                >
                                                    <type.icon size={24} />
                                                    <span className="text-xs font-bold">{type.label}</span>
                                                </button>
                                            ))}
                                        </div>
                                    </div>

                                    <div className="bg-white p-8 rounded-[2rem] shadow-sm border border-slate-200 space-y-6">
                                        <label className="text-xs font-black text-slate-400 uppercase tracking-widest mb-6 flex items-center gap-2">
                                            <Users size={16} /> Logistics
                                        </label>
                                        
                                        <div className="flex items-center justify-between p-4 bg-slate-50 rounded-2xl border border-slate-100">
                                            <div className="flex items-center gap-3">
                                                <div className="p-2 bg-indigo-100 text-indigo-600 rounded-lg"><Users size={18} /></div>
                                                <div>
                                                    <div className="text-sm font-bold text-slate-900">Group Assignment</div>
                                                    <div className="text-[10px] text-slate-500 font-bold uppercase">One submission per team</div>
                                                </div>
                                            </div>
                                            <Switch checked={config.group_assignment} onCheckedChange={(c) => updateConfig({ group_assignment: c })} />
                                        </div>

                                        <div className="flex items-center justify-between p-4 bg-slate-50 rounded-2xl border border-slate-100">
                                            <div className="flex items-center gap-3">
                                                <div className="p-2 bg-purple-100 text-purple-600 rounded-lg"><Scale size={18} /></div>
                                                <div>
                                                    <div className="text-sm font-bold text-slate-900">Peer Review</div>
                                                    <div className="text-[10px] text-slate-500 font-bold uppercase">Students grade each other</div>
                                                </div>
                                            </div>
                                            <Switch checked={config.peer_review} onCheckedChange={(c) => updateConfig({ peer_review: c })} />
                                        </div>
                                        
                                        {config.peer_review && (
                                            <div className="pl-4 border-l-2 border-purple-100 ml-4 animate-in fade-in slide-in-from-left-4">
                                                <div className="text-xs font-bold text-slate-500 mb-2">Required Reviews per Student</div>
                                                <Input 
                                                    type="number" 
                                                    className="w-20 bg-slate-50 border-slate-200 font-bold"
                                                    value={config.peer_review_count}
                                                    onChange={(e) => updateConfig({ peer_review_count: parseInt(e.target.value) })}
                                                />
                                            </div>
                                        )}
                                    </div>
                                </div>
                            </motion.div>
                        )}

                        {activeTab === 'Rubric' && (
                            <motion.div 
                                key="rubric"
                                initial={{ opacity: 0, y: 10 }}
                                animate={{ opacity: 1, y: 0 }}
                                exit={{ opacity: 0, scale: 0.98 }}
                                className="space-y-6 max-w-4xl mx-auto pb-10"
                            >
                                {config.rubric.map((crit, idx) => (
                                    <div key={crit.id} className="bg-white rounded-[2rem] border border-slate-200 overflow-hidden shadow-sm">
                                        <div className="p-6 bg-slate-50 border-b border-slate-200 flex items-start gap-4">
                                            <div className="w-8 h-8 bg-white rounded-lg flex items-center justify-center text-slate-400 font-black border border-slate-200 shadow-sm flex-shrink-0">
                                                {idx + 1}
                                            </div>
                                            <div className="flex-1 space-y-2">
                                                <Input 
                                                    value={crit.title} 
                                                    onChange={(e) => updateCriterion(crit.id, { title: e.target.value })}
                                                    className="font-bold text-lg bg-transparent border-none p-0 focus-visible:ring-0 text-slate-900 shadow-none h-auto"
                                                    placeholder="Criterion Title"
                                                />
                                                <Input 
                                                    value={crit.description}
                                                    onChange={(e) => updateCriterion(crit.id, { description: e.target.value })}
                                                    className="text-sm bg-transparent border-none p-0 focus-visible:ring-0 text-slate-500 shadow-none h-auto"
                                                    placeholder="Describe what is being evaluated..."
                                                />
                                            </div>
                                            <Button variant="ghost" size="icon" onClick={() => deleteCriterion(crit.id)} className="text-slate-400 hover:text-red-500 hover:bg-red-50"><Trash2 size={16}/></Button>
                                        </div>
                                        
                                        <div className="p-6 grid grid-cols-1 md:grid-cols-4 gap-4">
                                            {crit.levels.map((level, lIdx) => (
                                                <div key={level.id} className="border-2 border-slate-100 rounded-2xl p-4 hover:border-indigo-100 hover:bg-indigo-50/30 transition-all cursor-text group relative">
                                                    <div className="flex justify-between items-center mb-2">
                                                        <input 
                                                            value={level.label}
                                                            onChange={(e) => {
                                                                const newLevels = [...crit.levels];
                                                                newLevels[lIdx].label = e.target.value;
                                                                updateCriterion(crit.id, { levels: newLevels });
                                                            }}
                                                            className="font-bold text-xs bg-transparent border-none p-0 focus:ring-0 w-2/3 text-slate-700 outline-none"
                                                        />
                                                        <div className="flex items-center bg-white px-2 py-0.5 rounded-lg border border-slate-200 shadow-sm">
                                                            <input 
                                                                type="number"
                                                                value={level.points}
                                                                onChange={(e) => {
                                                                    const newLevels = [...crit.levels];
                                                                    newLevels[lIdx].points = parseInt(e.target.value) || 0;
                                                                    updateCriterion(crit.id, { levels: newLevels });
                                                                }}
                                                                className="w-8 text-center font-black text-xs bg-transparent border-none p-0 focus:ring-0 text-indigo-600 outline-none"
                                                            />
                                                            <span className="text-[8px] font-bold text-slate-400">PTS</span>
                                                        </div>
                                                    </div>
                                                    <textarea 
                                                        value={level.description}
                                                        onChange={(e) => {
                                                            const newLevels = [...crit.levels];
                                                            newLevels[lIdx].description = e.target.value;
                                                            updateCriterion(crit.id, { levels: newLevels });
                                                        }}
                                                        className="w-full h-20 bg-transparent border-none resize-none focus:ring-0 text-xs text-slate-500 leading-relaxed p-0 outline-none"
                                                        placeholder="Level description..."
                                                    />
                                                </div>
                                            ))}
                                        </div>
                                    </div>
                                ))}

                                <button 
                                    onClick={addCriterion}
                                    className="w-full py-6 border-2 border-dashed border-slate-300 rounded-[2rem] text-slate-400 font-bold hover:text-indigo-600 hover:border-indigo-300 hover:bg-indigo-50/50 transition-all flex items-center justify-center gap-2"
                                >
                                    <Plus size={20} /> Add Criterion
                                </button>
                            </motion.div>
                        )}
                    </AnimatePresence>
                </div>
            </div>

            {/* AI Modal Overlay */}
            <AnimatePresence>
                {isAiOpen && (
                    <div className="absolute inset-0 z-[100] flex items-center justify-center p-4 bg-slate-900/60 backdrop-blur-sm" onClick={() => setIsAiOpen(false)}>
                        <motion.div 
                            initial={{ opacity: 0, scale: 0.95 }}
                            animate={{ opacity: 1, scale: 1 }}
                            exit={{ opacity: 0, scale: 0.95 }}
                            className="bg-white w-full max-w-lg rounded-[2rem] shadow-2xl p-8 border border-white/20"
                            onClick={(e) => e.stopPropagation()}
                        >
                            <div className="flex flex-col items-center text-center mb-8">
                                <div className="w-16 h-16 bg-purple-50 text-purple-600 rounded-2xl flex items-center justify-center mb-4 shadow-inner">
                                    <Sparkles size={32} />
                                </div>
                                <h3 className="text-2xl font-black text-slate-900">AI Assignment Architect</h3>
                                <p className="text-slate-500 text-sm mt-2 max-w-xs">Describe your learning objective, and I'll draft the assignment details and rubric.</p>
                            </div>

                            <div className="space-y-4">
                                <Textarea 
                                    value={aiPrompt}
                                    onChange={(e) => setAiPrompt(e.target.value)}
                                    className="w-full h-32 p-4 bg-slate-50 border border-slate-200 rounded-2xl text-sm focus-visible:ring-purple-100 outline-none resize-none"
                                    placeholder="e.g. Create a persuasive essay assignment about climate change ethics..."
                                />
                                
                                <Button 
                                    onClick={handleAiGenerate} 
                                    disabled={isGenerating}
                                    className="w-full py-6 text-lg bg-purple-600 hover:bg-purple-700 shadow-xl shadow-purple-200 rounded-xl"
                                >
                                    {isGenerating ? "Generating..." : <><Wand2 className="mr-2" /> Generate Draft</>}
                                </Button>
                            </div>
                        </motion.div>
                    </div>
                )}
            </AnimatePresence>
        </DialogContent>
    </Dialog>
  );
};
