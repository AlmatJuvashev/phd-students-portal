
import React, { useState } from 'react';
import { 
  X, Plus, Save, Settings, Trash2, GripVertical, CheckSquare, 
  MessageCircle, Target, Sparkles, Bold, Italic, Sigma, Layers, 
  ArrowLeft, ArrowRight, Table as TableIcon, Bookmark, Calendar
} from 'lucide-react';
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Switch } from "@/components/ui/switch";
import { Badge } from "@/components/ui/badge";
import { cn } from "@/lib/utils";
import { Reorder, AnimatePresence, motion } from 'framer-motion';
import { QuizQuestion, QuestionType } from '../types';

interface QuizBuilderModalProps {
  isOpen: boolean;
  onClose: () => void;
  initialQuestions?: QuizQuestion[];
  initialConfig?: any;
  onSave: (questions: QuizQuestion[], config: any) => void;
}

const INITIAL_QUESTIONS: QuizQuestion[] = [
  { id: 'q1', type: 'multiple_choice', text: 'New Question', points: 10, options: [{ id: 'o1', text: 'Option 1', is_correct: true }, { id: 'o2', text: 'Option 2', is_correct: false }] },
];

export const QuizBuilderModal: React.FC<QuizBuilderModalProps> = ({ isOpen, onClose, initialQuestions, initialConfig, onSave }) => {
  const [questions, setQuestions] = useState<QuizQuestion[]>(initialQuestions || INITIAL_QUESTIONS);
  const [activeQuestionId, setActiveQuestionId] = useState<string | null>(questions[0]?.id || null);
  const [activeTab, setActiveTab] = useState<'editor' | 'settings'>('editor');
  
  const [config, setConfig] = useState(initialConfig || {
    time_limit_minutes: 60,
    passing_score: 80,
    shuffle_questions: true,
    show_results: true
  });

  const activeQuestion = questions.find(q => q.id === activeQuestionId);

  const addQuestion = (type: QuestionType) => {
    const newQ: QuizQuestion = {
      id: `q${Date.now()}`,
      type,
      text: type === 'section_header' ? 'New Section' : 'Question Text',
      points: ['section_header', 'page_break'].includes(type) ? 0 : 10,
      options: ['multiple_choice', 'multi_select'].includes(type) 
        ? [{ id: `o${Date.now()}_1`, text: 'Option 1', is_correct: true }, { id: `o${Date.now()}_2`, text: 'Option 2', is_correct: false }] 
        : undefined,
    };
    setQuestions([...questions, newQ]);
    setActiveQuestionId(newQ.id);
  };

  const updateActiveQuestion = (updates: Partial<QuizQuestion>) => {
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
      <DialogContent className="max-w-[95vw] w-[1400px] h-[90vh] flex flex-col p-0 gap-0 bg-slate-50 overlow-hidden">
        
        {/* Header */}
        <div className="h-16 border-b border-slate-200 bg-white px-6 flex items-center justify-between shrink-0">
            <div className="flex items-center gap-4">
                <div className="p-2 bg-indigo-100 text-indigo-600 rounded-lg"><Target size={20} /></div>
                <h2 className="font-bold text-lg text-slate-900">Quiz Builder</h2>
            </div>
            <div className="flex items-center gap-2">
                <Button variant="outline" onClick={() => setActiveTab('settings')} className={cn(activeTab==='settings' && "bg-slate-100")}>
                    <Settings size={16} className="mr-2"/> Settings
                </Button>
                <Button onClick={() => { onSave(questions, config); onClose(); }}>
                    <Save size={16} className="mr-2"/> Save Quiz
                </Button>
            </div>
        </div>

        <div className="flex-1 flex overflow-hidden">
            {/* Sidebar: Question List */}
            <div className="w-80 bg-white border-r border-slate-200 flex flex-col shrink-0">
                <div className="p-4 grid grid-cols-2 gap-2 border-b border-slate-100">
                     {[
                        { id: 'multiple_choice', icon: CheckSquare, label: 'Choice' },
                        { id: 'short_text', icon: MessageCircle, label: 'Text' },
                        { id: 'ordering', icon: Layers, label: 'Order' },
                        { id: 'section_header', icon: Bookmark, label: 'Section' }
                     ].map(type => (
                         <button 
                            key={type.id}
                            onClick={() => addQuestion(type.id as QuestionType)}
                            className="flex flex-col items-center justify-center p-3 rounded-lg border border-slate-200 hover:border-indigo-200 hover:bg-indigo-50 transition-all"
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
                                            ? "bg-white border-indigo-500 shadow-md ring-1 ring-indigo-500/20" 
                                            : "bg-white border-slate-200 hover:border-indigo-200"
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
                            
                            {/* Question Card */}
                            <div className="bg-white rounded-2xl shadow-sm border border-slate-200 p-8">
                                <div className="flex justify-between items-start mb-6">
                                    <div className="flex items-center gap-2">
                                        <Badge variant="outline" className="uppercase text-[10px] tracking-widest">{activeQuestion.type.replace('_',' ')}</Badge>
                                    </div>
                                    <Button variant="ghost" size="icon" onClick={deleteActiveQuestion} className="text-red-400 hover:text-red-500 hover:bg-red-50"><Trash2 size={16}/></Button>
                                </div>
                                
                                <div className="space-y-4 mb-8">
                                    <label className="text-xs font-black text-slate-400 uppercase">Question Text</label>
                                    <Textarea 
                                        value={activeQuestion.text}
                                        onChange={(e) => updateActiveQuestion({ text: e.target.value })}
                                        className="text-lg font-bold min-h-[100px] resize-none border-slate-200 focus:border-indigo-500"
                                        placeholder="Enter your question text..."
                                    />
                                </div>

                                {/* Options Editor */}
                                {['multiple_choice', 'multi_select'].includes(activeQuestion.type) && (
                                    <div className="space-y-4">
                                        <label className="text-xs font-black text-slate-400 uppercase">Answer Options</label>
                                        <div className="space-y-2">
                                            {activeQuestion.options?.map((opt, idx) => (
                                                <div key={opt.id} className="flex items-center gap-3 p-2 rounded-lg border border-transparent hover:border-slate-200 group">
                                                    <button 
                                                        onClick={() => {
                                                            const newOpts = activeQuestion.options?.map(o => 
                                                                activeQuestion.type === 'multiple_choice'
                                                                    ? { ...o, is_correct: o.id === opt.id }
                                                                    : (o.id === opt.id ? { ...o, is_correct: !o.is_correct } : o)
                                                            );
                                                            updateActiveQuestion({ options: newOpts });
                                                        }}
                                                        className={cn(
                                                            "w-8 h-8 rounded-full border-2 flex items-center justify-center transition-all",
                                                            opt.is_correct ? "bg-emerald-500 border-emerald-500 text-white" : "border-slate-300 text-slate-300"
                                                        )}
                                                    >
                                                        {opt.is_correct && <CheckSquare size={14} />}
                                                    </button>
                                                    <Input 
                                                        value={opt.text}
                                                        onChange={(e) => {
                                                            const newOpts = [...(activeQuestion.options || [])];
                                                            newOpts[idx].text = e.target.value;
                                                            updateActiveQuestion({ options: newOpts });
                                                        }}
                                                        className="flex-1"
                                                    />
                                                    <Button variant="ghost" size="icon" className="opacity-0 group-hover:opacity-100" onClick={() => {
                                                        const newOpts = activeQuestion.options?.filter(o => o.id !== opt.id);
                                                        updateActiveQuestion({ options: newOpts });
                                                    }}><X size={14}/></Button>
                                                </div>
                                            ))}
                                            <Button variant="outline" size="sm" onClick={() => {
                                                const newOpts = [...(activeQuestion.options || []), { id: `o${Date.now()}`, text: '', is_correct: false }];
                                                updateActiveQuestion({ options: newOpts });
                                            }} className="ml-11 border-dashed text-slate-500"><Plus size={14} className="mr-2"/> Add Option</Button>
                                        </div>
                                    </div>
                                )}

                                {/* Points */}
                                <div className="mt-8 pt-6 border-t border-slate-100 flex items-center justify-between">
                                     <div className="flex flex-col">
                                         <span className="text-xs font-bold text-slate-500">Points Value</span>
                                         <Input 
                                            type="number" 
                                            className="w-24 mt-1 font-bold" 
                                            value={activeQuestion.points} 
                                            onChange={(e) => updateActiveQuestion({ points: parseInt(e.target.value) || 0 })}
                                         />
                                     </div>
                                </div>

                                <div className="mt-6 pt-6 border-t border-slate-100 grid grid-cols-2 gap-6">
                                     <div>
                                         <label className="text-xs font-bold text-emerald-600 uppercase mb-2 block">Correct Feedback</label>
                                         <Textarea 
                                            value={activeQuestion.feedback_correct || ''}
                                            onChange={(e) => updateActiveQuestion({ feedback_correct: e.target.value })}
                                            className="bg-emerald-50/50 border-emerald-100 text-xs min-h-[80px]"
                                            placeholder="Great job!"
                                         />
                                     </div>
                                     <div>
                                         <label className="text-xs font-bold text-red-500 uppercase mb-2 block">Incorrect Feedback</label>
                                         <Textarea 
                                            value={activeQuestion.feedback_incorrect || ''}
                                            onChange={(e) => updateActiveQuestion({ feedback_incorrect: e.target.value })}
                                            className="bg-red-50/50 border-red-100 text-xs min-h-[80px]"
                                            placeholder="Try reviewing the material..."
                                         />
                                     </div>
                                </div>
                            </div>
                        </div>
                    </div>
                ) : activeTab === 'settings' ? (
                    <div className="flex-1 p-12 bg-white">
                        <div className="max-w-xl mx-auto space-y-8">
                             <h3 className="text-2xl font-black text-slate-900">Quiz Configuration</h3>
                             
                             <div className="space-y-6">
                                 <div className="flex items-center justify-between p-4 border border-slate-200 rounded-xl">
                                     <div>
                                         <div className="font-bold text-slate-900">Time Limit</div>
                                         <div className="text-xs text-slate-500">Duration in minutes (0 for unlimited)</div>
                                     </div>
                                     <Input type="number" className="w-24 font-bold" value={config.time_limit_minutes} onChange={(e) => setConfig({...config, time_limit_minutes: parseInt(e.target.value)})} />
                                 </div>

                                 <div className="flex items-center justify-between p-4 border border-slate-200 rounded-xl">
                                     <div>
                                         <div className="font-bold text-slate-900">Passing Score (%)</div>
                                         <div className="text-xs text-slate-500">Percentage required to pass</div>
                                     </div>
                                     <Input type="number" className="w-24 font-bold" value={config.passing_score} onChange={(e) => setConfig({...config, passing_score: parseInt(e.target.value)})} />
                                 </div>

                                 <div className="flex items-center justify-between p-4 border border-slate-200 rounded-xl">
                                     <div>
                                         <div className="font-bold text-slate-900">Shuffle Questions</div>
                                         <div className="text-xs text-slate-500">Randomize order for each attempt</div>
                                     </div>
                                     <Switch checked={config.shuffle_questions} onCheckedChange={(c) => setConfig({...config, shuffle_questions: c})} />
                                 </div>
                             </div>
                        </div>
                    </div>
                ) : (
                    <div className="flex-1 flex items-center justify-center text-slate-400">Select a question to edit</div>
                )}
            </div>
        </div>
      </DialogContent>
    </Dialog>
  );
};
