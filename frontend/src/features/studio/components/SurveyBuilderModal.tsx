
import React, { useState } from 'react';
import { 
  X, Plus, Save, Settings, Trash2, GripVertical, Star, Hash, 
  MessageSquare, CheckSquare, Layout, Bookmark, ArrowLeft
} from 'lucide-react';
import { Dialog, DialogContent } from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Switch } from "@/components/ui/switch";
import { Badge } from "@/components/ui/badge";
import { cn } from "@/lib/utils";
import { Reorder } from 'framer-motion';

// Types for Survey (local for now, should move to centralized types)
export type SurveyQuestionType = 'rating_stars' | 'scale_10' | 'likert_matrix' | 'open_feedback' | 'multiple_choice' | 'section_header';

export interface SurveyQuestion {
  id: string;
  type: SurveyQuestionType;
  text: string;
  required: boolean;
  options?: { id: string; text: string }[];
  matrixRows?: string[];
  matrixCols?: string[];
}

interface SurveyBuilderModalProps {
  isOpen: boolean;
  onClose: () => void;
  initialQuestions?: SurveyQuestion[];
  initialConfig?: any;
  onSave: (questions: SurveyQuestion[], config: any) => void;
}

const INITIAL_QUESTIONS: SurveyQuestion[] = [
  { id: 'sq1', type: 'rating_stars', text: 'How would you rate this course?', required: true },
];

export const SurveyBuilderModal: React.FC<SurveyBuilderModalProps> = ({ isOpen, onClose, initialQuestions, initialConfig, onSave }) => {
  const [questions, setQuestions] = useState<SurveyQuestion[]>(initialQuestions || INITIAL_QUESTIONS);
  const [activeQuestionId, setActiveQuestionId] = useState<string | null>(questions[0]?.id || null);
  const [activeTab, setActiveTab] = useState<'editor' | 'settings'>('editor');
  
  const [config, setConfig] = useState(initialConfig || {
    title: 'New Survey',
    anonymous: true,
    showProgressBar: true
  });

  const activeQuestion = questions.find(q => q.id === activeQuestionId);

  const addQuestion = (type: SurveyQuestionType) => {
    const newQ: SurveyQuestion = {
      id: `sq${Date.now()}`,
      type,
      text: 'New Question',
      required: true,
      matrixRows: type === 'likert_matrix' ? ['Row 1', 'Row 2'] : undefined,
      matrixCols: type === 'likert_matrix' ? ['Poor', 'Fair', 'Good', 'Excellent'] : undefined,
      options: type === 'multiple_choice' ? [{ id: 'o1', text: 'Option 1' }, { id: 'o2', text: 'Option 2' }] : undefined
    };
    setQuestions([...questions, newQ]);
    setActiveQuestionId(newQ.id);
  };

  const updateActiveQuestion = (updates: Partial<SurveyQuestion>) => {
    if (!activeQuestionId) return;
    setQuestions(questions.map(q => q.id === activeQuestionId ? { ...q, ...updates } : q));
  };

  const deleteActiveQuestion = () => {
    const newQs = questions.filter(q => q.id !== activeQuestionId);
    setQuestions(newQs);
    setActiveQuestionId(newQs[0]?.id || null);
  };

  return (
    <Dialog open={isOpen} onOpenChange={onClose}>
      <DialogContent className="max-w-[95vw] w-[1400px] h-[90vh] flex flex-col p-0 gap-0 bg-slate-50 overflow-hidden">
        
        {/* Header */}
        <div className="h-16 border-b border-slate-200 bg-white px-6 flex items-center justify-between shrink-0">
            <div className="flex items-center gap-4">
                <div className="p-2 bg-rose-100 text-rose-600 rounded-lg"><Layout size={20} /></div>
                <h2 className="font-bold text-lg text-slate-900">Survey Builder</h2>
            </div>
            <div className="flex items-center gap-2">
                <Button variant="outline" onClick={() => setActiveTab('settings')} className={cn(activeTab==='settings' && "bg-slate-100")}>
                    <Settings size={16} className="mr-2"/> Settings
                </Button>
                <Button onClick={() => { onSave(questions, config); onClose(); }}>
                    <Save size={16} className="mr-2"/> Save Survey
                </Button>
            </div>
        </div>

        <div className="flex-1 flex overflow-hidden">
            {/* Sidebar: Toolset */}
            <div className="w-72 bg-white border-r border-slate-200 flex flex-col shrink-0">
                <div className="p-4 grid grid-cols-2 gap-2 border-b border-slate-100">
                     {[
                        { id: 'rating_stars', icon: Star, label: 'Stars' },
                        { id: 'scale_10', icon: Hash, label: '0-10' },
                        { id: 'open_feedback', icon: MessageSquare, label: 'Open' },
                        { id: 'likert_matrix', icon: Layout, label: 'Matrix' },
                        { id: 'multiple_choice', icon: CheckSquare, label: 'Choice' },
                        { id: 'section_header', icon: Bookmark, label: 'Header' }
                     ].map(type => (
                         <button 
                            key={type.id}
                            onClick={() => addQuestion(type.id as SurveyQuestionType)}
                            className="flex flex-col items-center justify-center p-3 rounded-lg border border-slate-200 hover:border-rose-200 hover:bg-rose-50 transition-all"
                         >
                            <type.icon size={16} className="mb-1 text-slate-500" />
                            <span className="text-[10px] uppercase font-bold text-slate-600">{type.label}</span>
                         </button>
                     ))}
                </div>
                <div className="flex-1 overflow-y-auto p-2 bg-slate-50/50">
                     <Reorder.Group axis="y" values={questions} onReorder={setQuestions}>
                        {questions.map((q, i) => (
                            <Reorder.Item key={q.id} value={q}>
                                <div 
                                    onClick={() => setActiveQuestionId(q.id)}
                                    className={cn(
                                        "p-3 mb-2 rounded-lg border cursor-pointer flex items-center gap-3 transition-all",
                                        activeQuestionId === q.id 
                                            ? "bg-white border-rose-500 shadow-md ring-1 ring-rose-500/20" 
                                            : "bg-white border-slate-200 hover:border-rose-200"
                                    )}
                                >
                                    <span className="text-xs font-bold text-slate-400 w-4">{i+1}</span>
                                    <div className="flex-1 min-w-0">
                                        <div className="text-xs font-bold text-slate-700 truncate">{q.text}</div>
                                        <div className="text-[10px] text-slate-400 uppercase">{q.type.replace('_',' ')}</div>
                                    </div>
                                    <GripVertical size={14} className="text-slate-300" />
                                </div>
                            </Reorder.Item>
                        ))}
                     </Reorder.Group>
                </div>
            </div>

            {/* Main Editor */}
            <div className="flex-1 flex flex-col bg-slate-100/50 overflow-hidden relative">
                {activeTab === 'editor' && activeQuestion ? (
                    <div className="flex-1 overflow-y-auto p-8">
                        <div className="max-w-3xl mx-auto space-y-6">
                            
                            <div className="bg-white rounded-2xl shadow-sm border border-slate-200 p-8">
                                <div className="flex justify-between items-start mb-6">
                                    <div className="flex items-center gap-2">
                                        <Badge variant="outline" className="uppercase text-[10px] tracking-widest text-rose-600 border-rose-200 bg-rose-50">{activeQuestion.type.replace('_',' ')}</Badge>
                                    </div>
                                    <Button variant="ghost" size="icon" onClick={deleteActiveQuestion} className="text-red-400 hover:text-red-500 hover:bg-red-50"><Trash2 size={16}/></Button>
                                </div>
                                
                                <div className="space-y-4 mb-8">
                                    <label className="text-xs font-black text-slate-400 uppercase">Question Text</label>
                                    <Textarea 
                                        value={activeQuestion.text}
                                        onChange={(e) => updateActiveQuestion({ text: e.target.value })}
                                        className="text-lg font-bold min-h-[100px] resize-none border-slate-200 focus:border-rose-500"
                                        placeholder="Enter question..."
                                    />
                                </div>

                                {/* Dynamic Editors */}
                                {activeQuestion.type === 'likert_matrix' && (
                                    <div className="grid grid-cols-2 gap-6">
                                        <div className="space-y-2">
                                            <label className="text-xs font-bold text-slate-500">Rows</label>
                                            <div className="space-y-2">
                                                {activeQuestion.matrixRows?.map((r, i) => (
                                                    <div key={i} className="flex gap-2">
                                                        <Input value={r} onChange={(e) => {
                                                            const newRows = [...(activeQuestion.matrixRows || [])];
                                                            newRows[i] = e.target.value;
                                                            updateActiveQuestion({ matrixRows: newRows });
                                                        }} />
                                                        <Button variant="ghost" size="icon" onClick={() => updateActiveQuestion({ matrixRows: activeQuestion.matrixRows?.filter((_, idx) => idx !== i) })}><X size={14}/></Button>
                                                    </div>
                                                ))}
                                                <Button size="sm" variant="outline" className="w-full" onClick={() => updateActiveQuestion({ matrixRows: [...(activeQuestion.matrixRows || []), 'New Row'] })}>Add Row</Button>
                                            </div>
                                        </div>
                                        <div className="space-y-2">
                                            <label className="text-xs font-bold text-slate-500">Columns</label>
                                            <div className="space-y-2">
                                                {activeQuestion.matrixCols?.map((c, i) => (
                                                    <div key={i} className="flex gap-2">
                                                        <Input value={c} onChange={(e) => {
                                                            const newCols = [...(activeQuestion.matrixCols || [])];
                                                            newCols[i] = e.target.value;
                                                            updateActiveQuestion({ matrixCols: newCols });
                                                        }} />
                                                        <Button variant="ghost" size="icon" onClick={() => updateActiveQuestion({ matrixCols: activeQuestion.matrixCols?.filter((_, idx) => idx !== i) })}><X size={14}/></Button>
                                                    </div>
                                                ))}
                                                <Button size="sm" variant="outline" className="w-full" onClick={() => updateActiveQuestion({ matrixCols: [...(activeQuestion.matrixCols || []), 'New Col'] })}>Add Column</Button>
                                            </div>
                                        </div>
                                    </div>
                                )}

                                {activeQuestion.type === 'multiple_choice' && (
                                     <div className="space-y-2">
                                        {activeQuestion.options?.map((opt, idx) => (
                                            <div key={opt.id} className="flex items-center gap-2">
                                                <div className="w-4 h-4 rounded-full border-2 border-slate-300" />
                                                <Input value={opt.text} onChange={(e) => {
                                                    const newOpts = [...(activeQuestion.options || [])];
                                                    newOpts[idx].text = e.target.value;
                                                    updateActiveQuestion({ options: newOpts });
                                                }} />
                                                <Button variant="ghost" size="icon" onClick={() => updateActiveQuestion({ options: activeQuestion.options?.filter(o => o.id !== opt.id) })}><X size={14}/></Button>
                                            </div>
                                        ))}
                                        <Button size="sm" variant="outline" onClick={() => updateActiveQuestion({ options: [...(activeQuestion.options || []), { id: `o${Date.now()}`, text: '' }] })}>Add Option</Button>
                                     </div>
                                )}

                                <div className="mt-8 pt-6 border-t border-slate-100 flex items-center justify-between">
                                     <div className="flex items-center gap-2">
                                         <span className="text-xs font-bold text-slate-700">Required</span>
                                         <Switch checked={activeQuestion.required} onCheckedChange={(c) => updateActiveQuestion({ required: c })} />
                                     </div>
                                </div>
                            </div>
                        </div>
                    </div>
                ) : activeTab === 'settings' ? (
                    <div className="flex-1 p-12 bg-white">
                        <div className="max-w-xl mx-auto space-y-6">
                             <h3 className="text-2xl font-black text-slate-900">Config</h3>
                             <div className="space-y-4">
                                 <div className="space-y-2">
                                     <label className="text-xs font-bold text-slate-500">Title</label>
                                     <Input value={config.title} onChange={(e) => setConfig({...config, title: e.target.value})} />
                                 </div>
                                 <div className="flex items-center justify-between p-4 border border-slate-200 rounded-xl">
                                     <span className="font-bold text-slate-900">Anonymous Responses</span>
                                     <Switch checked={config.anonymous} onCheckedChange={(c) => setConfig({...config, anonymous: c})} />
                                 </div>
                             </div>
                        </div>
                    </div>
                ) : (
                    <div className="flex-1 flex items-center justify-center text-slate-400">Select a question</div>
                )}
            </div>
        </div>
      </DialogContent>
    </Dialog>
  );
};
