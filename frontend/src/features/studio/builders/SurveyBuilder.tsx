
import React, { useState } from 'react';
import { Reorder, motion, AnimatePresence } from 'framer-motion';
import { 
  ArrowLeft, Plus, Save, Settings, Trash2, Heading, 
  SeparatorHorizontal, GripVertical, CheckSquare, Eye, 
  Lightbulb, MessageCircle, AlertCircle, FileText, X,
  Bookmark, MoreVertical, Play, Clock, Shuffle, Award,
  Database, GitBranch, Target, Layers, Sparkles, Image as ImageIcon,
  GitMerge, Maximize2, Minimize2, Bold, Italic, List, Sigma, Table as TableIcon,
  Star, Hash, MessageSquare, ClipboardList, Layout
} from 'lucide-react';
import { Button, Input, Tabs, Switch, IconButton, AvatarGroup, Tooltip, Badge } from '@/features/admin/components/AdminUI';
import { cn } from '@/lib/utils';
import { useNavigate, useParams } from 'react-router-dom';

const ACTIVE_DESIGNERS = [
  { initials: 'AD', color: 'bg-indigo-600' },
  { initials: 'MK', color: 'bg-purple-500' },
];

type SurveyQuestionType = 'rating_stars' | 'scale_10' | 'likert_matrix' | 'open_feedback' | 'multiple_choice' | 'section_header';

interface SurveyQuestion {
  id: string;
  type: SurveyQuestionType;
  text: string;
  subtitle?: string;
  required: boolean;
  options?: { id: string; text: string }[];
  matrixRows?: string[];
  matrixCols?: string[];
}

const INITIAL_QUESTIONS: SurveyQuestion[] = [
  { id: 'sq1', type: 'section_header', text: 'Course Experience', required: false },
  { id: 'sq2', type: 'rating_stars', text: 'Overall, how would you rate this module?', required: true },
  { id: 'sq3', type: 'scale_10', text: 'How likely are you to recommend this program to a colleague?', required: true },
];

export const SurveyBuilder: React.FC = () => {
    const navigate = useNavigate();
    const { programId, courseId, nodeId } = useParams();
    
    const onNavigate = (path: string) => navigate(path);
    const handleBack = () => {
        if (programId) {
            navigate(`/admin/studio/programs/${programId}/builder`);
        } else if (courseId) {
            navigate(`/admin/studio/courses/${courseId}/builder`);
        } else {
            navigate(-1);
        }
    };

  const [questions, setQuestions] = useState<SurveyQuestion[]>(INITIAL_QUESTIONS);
  const [activeQuestionId, setActiveQuestionId] = useState<string | null>(null);
  const [isZenMode, setIsZenMode] = useState(false);
  const [isPreviewMarkdown, setIsPreviewMarkdown] = useState(false);
  
  const [surveyConfig, setSurveyConfig] = useState({
    title: 'Student Feedback: Semester 1',
    anonymous: true,
    showProgressBar: true
  });

  const activeQuestion = questions.find(q => q.id === activeQuestionId);

  const addQuestion = (type: SurveyQuestionType) => {
    const newQ: SurveyQuestion = {
      id: `sq${Date.now()}`,
      type,
      text: 'Enter your question text...',
      required: true,
      matrixRows: type === 'likert_matrix' ? ['Course Content', 'Instructor Performance'] : undefined,
      matrixCols: type === 'likert_matrix' ? ['Poor', 'Fair', 'Good', 'Excellent'] : undefined,
      options: type === 'multiple_choice' ? [{ id: 'o1', text: 'Option 1' }, { id: 'o2', text: 'Option 2' }] : undefined
    };
    setQuestions([...questions, newQ]);
    setActiveQuestionId(newQ.id);
  };

  const updateActiveQuestion = (updates: Partial<SurveyQuestion>) => {
    if (!activeQuestionId) return;
    setQuestions(prev => prev.map(q => q.id === activeQuestionId ? { ...q, ...updates } : q));
  };

  const deleteActiveQuestion = () => {
    setQuestions(prev => prev.filter(q => q.id !== activeQuestionId));
    setActiveQuestionId(null);
  };

  const renderStem = (text: string) => {
    if (!isPreviewMarkdown) return null;
    const parts = text.split(/(\$\$.*?\$\$)/g);
    return (
      <div className="prose prose-slate max-w-none text-slate-800 font-medium">
        {parts.map((part, i) => {
          if (part.startsWith('$$') && part.endsWith('$$')) {
            return (
              <span key={i} className="inline-block px-2 py-1 mx-1 bg-indigo-50 text-indigo-700 rounded font-serif italic border border-indigo-100 shadow-sm">
                {part.slice(2, -2)}
              </span>
            );
          }
          return <span key={i}>{part}</span>;
        })}
      </div>
    );
  };

  return (
    <div className="flex flex-col h-[calc(100vh-4rem)] bg-slate-50 font-sans overflow-hidden">
      {/* 1. Header */}
      <div className="h-20 bg-white border-b border-slate-200 px-8 flex items-center justify-between flex-shrink-0 z-30 shadow-sm">
        <div className="flex items-center gap-6">
          {!isZenMode && <IconButton icon={ArrowLeft} onClick={handleBack} />}
          <div>
             <div className="flex items-center gap-2 mb-1">
                <span className="text-[9px] font-black uppercase text-rose-600 bg-rose-50 px-2 py-0.5 rounded-full tracking-widest border border-rose-100">Survey Authoring</span>
                {isZenMode && <span className="text-[9px] font-black uppercase text-amber-600 bg-amber-50 px-2 py-0.5 rounded-full tracking-widest border border-amber-100">Focus Flow</span>}
             </div>
             <h2 className="font-black text-slate-900 leading-none text-xl tracking-tight">{surveyConfig.title}</h2>
          </div>
        </div>
        
        <div className="flex items-center gap-8">
          <AvatarGroup users={ACTIVE_DESIGNERS} />
          <div className="flex items-center gap-2">
            <IconButton 
              icon={isZenMode ? Minimize2 : Maximize2} 
              onClick={() => setIsZenMode(!isZenMode)}
              variant={isZenMode ? "primary" : "ghost"}
            />
            <Button variant="secondary" icon={Play} onClick={() => window.open('/admin/studio/courses/c1/survey/preview', '_blank')}>Preview</Button>
            <Button variant="primary" icon={Save}>Publish Survey</Button>
          </div>
        </div>
      </div>

      <div className="flex-1 flex overflow-hidden relative">
        {/* 2. Sidebar: Tools */}
        <AnimatePresence>
          {!isZenMode && (
            <motion.div 
              initial={{ width: 0, opacity: 0 }}
              animate={{ width: 320, opacity: 1 }}
              exit={{ width: 0, opacity: 0 }}
              className="w-80 bg-white border-r border-slate-200 flex flex-col flex-shrink-0 z-20 shadow-xl overflow-hidden"
            >
              <div className="p-5 border-b border-slate-200 bg-slate-50/50 space-y-4">
                 <div className="grid grid-cols-3 gap-2">
                    {[
                      { id: 'rating_stars', icon: Star, label: 'Stars' },
                      { id: 'scale_10', icon: Hash, label: 'Scale 0-10' },
                      { id: 'likert_matrix', icon: TableIcon, label: 'Likert' },
                      { id: 'open_feedback', icon: MessageSquare, label: 'Open' },
                      { id: 'multiple_choice', icon: CheckSquare, label: 'Choice' },
                      { id: 'section_header', icon: Bookmark, label: 'Header' }
                    ].map(type => (
                      <button 
                        key={type.id} 
                        onClick={() => addQuestion(type.id as any)}
                        className="flex flex-col items-center gap-1.5 p-2 hover:bg-white hover:shadow-lg rounded-xl border border-transparent hover:border-rose-100 text-slate-500 hover:text-rose-600 transition-all group"
                      >
                        <div className="p-2 bg-slate-100 rounded-lg group-hover:bg-rose-50 transition-colors">
                          <type.icon size={16} />
                        </div>
                        <span className="text-[8px] font-black uppercase">{type.label}</span>
                      </button>
                    ))}
                 </div>
              </div>
              
              <div className="flex-1 overflow-y-auto p-3 bg-slate-50/30">
                <Reorder.Group axis="y" values={questions} onReorder={setQuestions} className="space-y-1.5">
                  {questions.map((q) => {
                    const isSection = q.type === 'section_header';
                    return (
                      <Reorder.Item key={q.id} value={q}>
                        <motion.div 
                          layout
                          onClick={() => setActiveQuestionId(q.id)}
                          className={cn(
                            "p-3 rounded-xl border bg-white cursor-pointer relative group flex gap-3 items-center",
                            activeQuestionId === q.id 
                              ? "border-rose-500 shadow-lg ring-2 ring-rose-500/10" 
                              : "border-slate-200 hover:border-rose-200",
                            isSection ? "mt-4 font-black" : "ml-4"
                          )}
                        >
                          <GripVertical size={12} className="text-slate-300 cursor-grab" />
                          <div className="flex-1 min-w-0">
                             <h4 className={cn("text-[11px] font-bold truncate", isSection ? "text-rose-900" : "text-slate-700")}>
                               {q.text}
                             </h4>
                          </div>
                        </motion.div>
                      </Reorder.Item>
                    );
                  })}
                </Reorder.Group>
              </div>
            </motion.div>
          )}
        </AnimatePresence>

        {/* 3. Main Area: Editor */}
        <div className="flex-1 bg-slate-100/50 flex flex-col overflow-hidden">
          <div className="flex-1 overflow-y-auto p-8 custom-scrollbar">
            <div className="max-w-4xl mx-auto">
              <AnimatePresence mode="wait">
              {activeQuestion ? (
                <motion.div 
                    key={activeQuestion.id}
                    initial={{ opacity: 0, y: 15 }}
                    animate={{ opacity: 1, y: 0 }}
                    exit={{ opacity: 0, scale: 0.98 }}
                    className="w-full bg-white rounded-[2.5rem] shadow-2xl border border-slate-100 relative flex flex-col overflow-hidden mb-12"
                >
                  <div className="p-10 space-y-10">
                      <div className="flex justify-between items-start">
                         <div className="space-y-1">
                            <label className="text-[10px] font-black text-rose-400 uppercase tracking-widest flex items-center gap-2">
                               <Layout size={12} /> {activeQuestion.type.replace('_', ' ')} Block
                            </label>
                            <h3 className="text-xl font-black text-slate-900">Survey Logic</h3>
                         </div>
                         <IconButton icon={Trash2} onClick={deleteActiveQuestion} className="text-red-400 hover:bg-red-50" />
                      </div>

                      <div className="space-y-4">
                         <div className="flex justify-between items-center">
                            <label className="text-[10px] font-black text-slate-400 uppercase tracking-widest">Question Text</label>
                            <button 
                              onClick={() => setIsPreviewMarkdown(!isPreviewMarkdown)}
                              className={cn("text-[10px] font-bold px-2 py-0.5 rounded", isPreviewMarkdown ? "bg-rose-600 text-white" : "bg-slate-100 text-slate-500")}
                            >
                               {isPreviewMarkdown ? 'Edit' : 'Preview'}
                            </button>
                         </div>
                         
                         {isPreviewMarkdown ? (
                            <div className="w-full min-h-[112px] p-4 bg-slate-50 rounded-2xl border border-slate-200">
                               {renderStem(activeQuestion.text)}
                            </div>
                         ) : (
                            <textarea
                              value={activeQuestion.text}
                              onChange={(e) => updateActiveQuestion({ text: e.target.value })} 
                              className="w-full text-2xl font-black text-slate-900 border-none p-0 focus:ring-0 resize-none leading-tight placeholder:text-slate-200 h-28"
                              placeholder="e.g. How satisfied are you with the curriculum?"
                            />
                         )}
                      </div>

                      {/* Question Content Based on Type */}
                      <div className="p-8 bg-slate-50 rounded-3xl border border-slate-200 flex flex-col items-center justify-center min-h-[200px] border-dashed">
                          {activeQuestion.type === 'rating_stars' && (
                             <div className="flex gap-4">
                                {[1,2,3,4,5].map(i => <Star key={i} size={48} className="text-slate-200" strokeWidth={1} />)}
                             </div>
                          )}
                          {activeQuestion.type === 'scale_10' && (
                             <div className="flex gap-2 w-full max-w-lg">
                                {[0,1,2,3,4,5,6,7,8,9,10].map(i => (
                                  <div key={i} className="flex-1 h-12 flex items-center justify-center border-2 border-slate-200 rounded-xl font-black text-slate-300">
                                    {i}
                                  </div>
                                ))}
                             </div>
                          )}
                          {activeQuestion.type === 'open_feedback' && (
                             <div className="w-full border-b-2 border-slate-200 text-slate-300 font-bold italic">
                                Placeholder for long-form response...
                             </div>
                          )}
                          {activeQuestion.type === 'likert_matrix' && (
                             <div className="w-full space-y-4">
                                <div className="grid grid-cols-2 gap-4">
                                   <div className="space-y-2">
                                      <label className="text-[9px] font-black text-slate-400 uppercase">Rows</label>
                                      {activeQuestion.matrixRows?.map((r, i) => <Input key={i} value={r} className="h-8 text-xs" />)}
                                   </div>
                                   <div className="space-y-2">
                                      <label className="text-[9px] font-black text-slate-400 uppercase">Columns</label>
                                      {activeQuestion.matrixCols?.map((c, i) => <Input key={i} value={c} className="h-8 text-xs" />)}
                                   </div>
                                </div>
                             </div>
                          )}
                      </div>

                      <div className="pt-10 border-t border-slate-100 flex justify-between items-center">
                         <div className="flex items-center gap-3">
                            <span className="text-xs font-bold text-slate-700">Required Response</span>
                            <Switch checked={activeQuestion.required} onCheckedChange={(v) => updateActiveQuestion({ required: v })} />
                         </div>
                      </div>
                  </div>
                </motion.div>
              ) : (
                <div className="flex flex-col items-center justify-center py-40 text-center max-w-sm mx-auto">
                   <div className="w-24 h-24 bg-white rounded-[2rem] shadow-xl flex items-center justify-center text-rose-500 mb-8 border border-slate-100">
                      <ClipboardList size={48} strokeWidth={2.5} />
                   </div>
                   <h3 className="text-xl font-black text-slate-900 mb-2">Survey Workspace</h3>
                   <p className="text-sm text-slate-500">Add feedback blocks from the left to start collecting data.</p>
                </div>
              )}
              </AnimatePresence>
            </div>
          </div>
        </div>

        {/* 4. Sidebar: Global Settings */}
        {!isZenMode && (
          <div className="w-80 bg-white border-l border-slate-200 flex flex-col flex-shrink-0 z-20">
             <div className="p-4 border-b border-slate-100 font-black text-xs uppercase tracking-widest text-slate-400">General Configuration</div>
             <div className="p-6 space-y-6">
                <div className="p-4 bg-rose-50 rounded-2xl border border-rose-100 space-y-4">
                   <div className="flex items-center justify-between">
                      <span className="text-[10px] font-black text-rose-700 uppercase">Anonymous Mode</span>
                      <Switch checked={surveyConfig.anonymous} onCheckedChange={(v) => setSurveyConfig({...surveyConfig, anonymous: v})} />
                   </div>
                   <p className="text-[9px] text-rose-600 font-medium leading-relaxed">Identity of participants will be hidden in results export.</p>
                </div>
                
                <div className="space-y-2">
                   <label className="text-[10px] font-black text-slate-400 uppercase tracking-widest">Survey Title</label>
                   <Input value={surveyConfig.title} onChange={(e:any) => setSurveyConfig({...surveyConfig, title: e.target.value})} />
                </div>
             </div>
          </div>
        )}
      </div>
    </div>
  );
};
