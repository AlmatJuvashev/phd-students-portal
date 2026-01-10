
import React, { useState } from 'react';
import { Reorder, motion, AnimatePresence } from 'framer-motion';
import { 
  ArrowLeft, Plus, Save, Settings, Trash2, GripVertical, CheckSquare, 
  HelpCircle, AlertCircle, Clock, Trophy, Shuffle, Target, Eye, 
  Play, MoreVertical, Check, X, Image as ImageIcon, Layout
} from 'lucide-react';
import { Button, Input, Switch, Badge, IconButton, AvatarGroup } from '@/features/admin/components/AdminUI';
import { cn } from '@/lib/utils';
import { useNavigate, useParams } from 'react-router-dom';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { getCourseModules, updateCourseActivity, createCourseActivity } from '@/features/curriculum/api';
import { api } from '@/api/client';
import { toast } from 'sonner';
import { useEffect } from 'react';

const ACTIVE_DESIGNERS = [
  { initials: 'MK', color: 'bg-purple-500' },
  { initials: 'JS', color: 'bg-emerald-500' },
];

type QuestionType = 'multiple_choice' | 'single_choice' | 'true_false';

interface QuizQuestion {
  id: string;
  type: QuestionType;
  text: string;
  points: number;
  options: { id: string; text: string; isCorrect: boolean }[];
  explanation?: string;
}

const INITIAL_QUESTIONS: QuizQuestion[] = [
  { 
    id: 'q1', 
    type: 'single_choice', 
    text: 'What is the primary function of the mitochondrion?', 
    points: 10,
    options: [
      { id: 'o1', text: 'Energy production', isCorrect: true },
      { id: 'o2', text: 'Protein synthesis', isCorrect: false },
      { id: 'o3', text: 'Cell division', isCorrect: false }
    ],
    explanation: 'Mitochondria are known as the powerhouses of the cell.'
  }
];

export const QuizBuilder: React.FC = () => {
    const navigate = useNavigate();
    const { courseId, nodeId } = useParams(); 
    const queryClient = useQueryClient();
    
    const [questions, setQuestions] = useState<QuizQuestion[]>([]);
    const [activeQuestionId, setActiveQuestionId] = useState<string | null>(null);
    const [quizConfig, setQuizConfig] = useState({
      title: 'Quiz Designer',
      timeLimit: 15,
      passingScore: 70,
      shuffle: true
    });

    const [parentIds, setParentIds] = useState<{moduleId: string, lessonId: string} | null>(null);

    // --- API Connectivity ---
    const { data: modulesData, isLoading } = useQuery({
      queryKey: ['courseModules', courseId],
      queryFn: () => getCourseModules(courseId!),
      enabled: !!courseId
    });

    const updateActivityMutation = useMutation({
      mutationFn: (content: any) => {
        if (!parentIds || !nodeId) throw new Error("Context missing");
        return updateCourseActivity(nodeId, parentIds.lessonId, parentIds.moduleId, courseId!, { content: JSON.stringify(content) });
      },
      onSuccess: () => {
        queryClient.invalidateQueries({ queryKey: ['courseModules', courseId] });
        toast.success('Quiz published successfully');
      },
      onError: () => toast.error('Failed to publish quiz')
    });

    const createActivityMutation = useMutation({
      mutationFn: (data: any) => {
        if (!parentIds) throw new Error("Module/Lesson context missing for creation");
        return createCourseActivity(parentIds.lessonId, parentIds.moduleId, courseId!, data);
      },
      onSuccess: (response: any) => {
        queryClient.invalidateQueries({ queryKey: ['courseModules', courseId] });
        toast.success('New quiz created');
        const newNodeId = response.data.id;
        navigate(`/admin/studio/courses/${courseId}/quiz/${newNodeId}/builder`, { replace: true });
      },
      onError: () => toast.error('Failed to create quiz')
    });

    // Hydrate from API
    useEffect(() => {
      if (modulesData && Array.isArray(modulesData)) {
        let foundActivity: any = null;
        
        if (nodeId) {
          for (const m of modulesData) {
            for (const l of (m.lessons || [])) {
              const a = (l.activities || []).find((act: any) => act.id === nodeId);
              if (a) {
                foundActivity = a;
                setParentIds({ moduleId: m.id, lessonId: l.id });
                break;
              }
            }
            if (foundActivity) break;
          }
        } else {
           // Creation mode: Default to first module/lesson
           if (modulesData.length > 0 && modulesData[0].lessons?.length > 0) {
             setParentIds({ moduleId: modulesData[0].id, lessonId: modulesData[0].lessons[0].id });
           }
           setQuestions([]);
           setQuizConfig({
             title: 'New Quiz Step',
             timeLimit: 15,
             passingScore: 70,
             shuffle: true
           });
           return;
        }

        if (foundActivity) {
          const content = typeof foundActivity.content === 'string' ? JSON.parse(foundActivity.content || '{}') : (foundActivity.content || {});
          if (content.questions) {
            setQuestions(content.questions.map((q: any) => ({
              ...q,
              id: q.id || `q_${Math.random()}`,
              options: q.options || []
            })));
          }
          setQuizConfig({
            title: foundActivity.title || 'Quiz Designer',
            timeLimit: content.timeLimit || 15,
            passingScore: content.passingScore || 70,
            shuffle: !!content.shuffleQuestions
          });
          if (content.questions?.length > 0 && !activeQuestionId) {
            setActiveQuestionId(content.questions[0].id);
          }
        }
      }
    }, [modulesData, nodeId]);

    const handleSave = () => {
      const content = {
        ...quizConfig,
        shuffleQuestions: quizConfig.shuffle,
        questions: questions
      };
      
      if (nodeId) {
        updateActivityMutation.mutate(content);
      } else {
        createActivityMutation.mutate({
          title: quizConfig.title,
          type: 'quiz',
          content: JSON.stringify(content),
          order: 99 // Default to end
        });
      }
    };

  const activeQuestion = questions.find(q => q.id === activeQuestionId);

  const addQuestion = () => {
    const newQ: QuizQuestion = {
      id: `q${Date.now()}`,
      type: 'single_choice',
      text: 'New Question',
      points: 10,
      options: [
        { id: `o${Date.now()}_1`, text: 'Option 1', isCorrect: true },
        { id: `o${Date.now()}_2`, text: 'Option 2', isCorrect: false }
      ]
    };
    setQuestions([...questions, newQ]);
    setActiveQuestionId(newQ.id);
  };

  const updateActiveQuestion = (updates: Partial<QuizQuestion>) => {
    if (!activeQuestionId) return;
    setQuestions(prev => prev.map(q => q.id === activeQuestionId ? { ...q, ...updates } : q));
  };

  const updateOption = (qId: string, oId: string, updates: any) => {
    setQuestions(prev => prev.map(q => {
        if (q.id !== qId) return q;
        return {
            ...q,
            options: q.options.map(o => o.id === oId ? { ...o, ...updates } : o)
        };
    }));
  };

  const deleteActiveQuestion = () => {
      const idx = questions.findIndex(q => q.id === activeQuestionId);
      const newQuestions = questions.filter(q => q.id !== activeQuestionId);
      setQuestions(newQuestions);
      if (newQuestions.length > 0) {
          setActiveQuestionId(newQuestions[Math.max(0, idx - 1)].id);
      } else {
          setActiveQuestionId(null);
      }
  };

  return (
    <div className="flex flex-col h-[calc(100vh-4rem)] bg-slate-50 font-sans overflow-hidden">
      {/* Header */}
      <div className="h-20 bg-white border-b border-slate-200 px-8 flex items-center justify-between flex-shrink-0 z-30 shadow-sm">
        <div className="flex items-center gap-6">
          <IconButton icon={ArrowLeft} onClick={() => navigate(`/admin/studio/courses/${courseId}/builder`)} />
          <div>
             <div className="flex items-center gap-2 mb-1">
                <span className="text-[9px] font-black uppercase text-indigo-600 bg-indigo-50 px-2 py-0.5 rounded-full tracking-widest border border-indigo-100">Quiz Studio</span>
             </div>
             <Input 
               value={quizConfig.title}
               onChange={(e: any) => setQuizConfig({...quizConfig, title: e.target.value})}
               className="font-black text-slate-900 text-xl border-none p-0 h-auto focus:ring-0 w-96 bg-transparent"
             />
          </div>
        </div>
        <div className="flex items-center gap-6">
          <AvatarGroup users={ACTIVE_DESIGNERS} />
          <Button variant="primary" icon={Save} onClick={handleSave} disabled={updateActivityMutation.isPending}>
            {updateActivityMutation.isPending ? 'Publishing...' : 'Publish Quiz'}
          </Button>
        </div>
      </div>

      <div className="flex-1 flex overflow-hidden">
        {/* Sidebar: Questions List */}
        <div className="w-80 bg-white border-r border-slate-200 flex flex-col flex-shrink-0 z-20">
           <div className="p-4 border-b border-slate-200 bg-slate-50 flex justify-between items-center">
              <h3 className="font-black text-xs uppercase tracking-wide text-slate-500">Questions ({questions.length})</h3>
              <Button size="sm" variant="ghost" icon={Plus} onClick={addQuestion}>Add</Button>
           </div>
           
           <div className="flex-1 overflow-y-auto p-3">
             <Reorder.Group axis="y" values={questions} onReorder={setQuestions} className="space-y-2">
               {questions.map((q, idx) => (
                 <Reorder.Item key={q.id} value={q}>
                   <div 
                     onClick={() => setActiveQuestionId(q.id)}
                     className={cn(
                       "p-3 rounded-xl border cursor-pointer relative group transition-all",
                       activeQuestionId === q.id 
                         ? "bg-indigo-50 border-indigo-200 shadow-sm" 
                         : "bg-white border-slate-200 hover:border-slate-300"
                     )}
                   >
                     <div className="flex items-start gap-3">
                        <span className={cn("flex-shrink-0 w-5 h-5 flex items-center justify-center rounded text-[10px] font-bold mt-0.5", activeQuestionId === q.id ? "bg-indigo-100 text-indigo-700" : "bg-slate-100 text-slate-500")}>
                            {idx + 1}
                        </span>
                        <div className="flex-1 min-w-0">
                           <div className={cn("text-xs font-bold truncate mb-1", activeQuestionId === q.id ? "text-indigo-900" : "text-slate-700")}>{q.text}</div>
                           <Badge variant="secondary" className="text-[9px] uppercase px-1 py-0">{q.type.replace('_', ' ')}</Badge>
                        </div>
                        <GripVertical size={14} className="text-slate-300 opacity-0 group-hover:opacity-100 cursor-grab" />
                     </div>
                   </div>
                 </Reorder.Item>
               ))}
             </Reorder.Group>
           </div>
           
           {/* Global Settings Summary */}
           <div className="p-4 border-t border-slate-200 bg-slate-50 space-y-3">
               <div className="flex items-center justify-between text-xs font-bold text-slate-600">
                   <div className="flex items-center gap-2"><Clock size={14} /> Time Limit</div>
                   <span>{quizConfig.timeLimit} mins</span>
               </div>
               <div className="flex items-center justify-between text-xs font-bold text-slate-600">
                   <div className="flex items-center gap-2"><Trophy size={14} /> Pass Score</div>
                   <span>{quizConfig.passingScore}%</span>
               </div>
               <div className="flex items-center justify-between text-xs font-bold text-slate-600">
                   <div className="flex items-center gap-2"><Shuffle size={14} /> Shuffle</div>
                   <Switch checked={quizConfig.shuffle} onCheckedChange={(v) => setQuizConfig({...quizConfig, shuffle: v})} />
               </div>
           </div>
        </div>

        {/* Main Editor */}
        <div className="flex-1 bg-slate-100/50 flex flex-col overflow-y-auto p-8 items-center">
           {activeQuestion ? (
              <div className="w-full max-w-3xl space-y-6">
                 {/* Question Card */}
                 <div className="bg-white rounded-2xl shadow-sm border border-slate-200 p-8 space-y-6 relative overflow-hidden">
                    <div className="flex justify-between items-start">
                        <div className="flex items-center gap-2">
                             <div className="p-2 bg-indigo-50 text-indigo-600 rounded-lg">
                                 <HelpCircle size={20} />
                             </div>
                             <div>
                                 <label className="text-[10px] font-black text-slate-400 uppercase tracking-widest block">Question Text</label>
                             </div>
                        </div>
                        <div className="flex items-center gap-2">
                            <select 
                                value={activeQuestion.type}
                                onChange={(e) => updateActiveQuestion({ type: e.target.value as any })}
                                className="bg-slate-50 border border-slate-200 rounded-lg text-xs font-bold p-2 outline-none"
                            >
                                <option value="single_choice">Single Choice</option>
                                <option value="multiple_choice">Multiple Choice</option>
                                <option value="true_false">True / False</option>
                            </select>
                            <IconButton icon={Trash2} onClick={deleteActiveQuestion} className="text-red-400 hover:bg-red-50" />
                        </div>
                    </div>

                    <textarea
                        value={activeQuestion.text}
                        onChange={(e) => updateActiveQuestion({ text: e.target.value })}
                        className="w-full text-lg font-bold text-slate-900 border-none p-0 focus:ring-0 resize-none h-24 placeholder:text-slate-300"
                        placeholder="Type your question here..."
                    />

                     {/* Options Editor */}
                    <div className="space-y-3">
                        <label className="text-[10px] font-black text-slate-400 uppercase tracking-widest block">Answer Options</label>
                        {activeQuestion.options.map((opt, idx) => (
                            <div key={opt.id} className={cn("flex items-center gap-3 p-3 rounded-xl border transition-all", opt.isCorrect ? "bg-emerald-50 border-emerald-200" : "bg-white border-slate-200")}>
                                <div className="flex-shrink-0">
                                    <button 
                                        onClick={() => updateOption(activeQuestion.id, opt.id, { isCorrect: !opt.isCorrect })}
                                        className={cn("w-6 h-6 rounded-full border-2 flex items-center justify-center transition-all", opt.isCorrect ? "border-emerald-500 bg-emerald-500 text-white" : "border-slate-300 text-transparent hover:border-slate-400")}
                                    >
                                        <Check size={14} strokeWidth={3} />
                                    </button>
                                </div>
                                <Input 
                                    value={opt.text}
                                    onChange={(e: any) => updateOption(activeQuestion.id, opt.id, { text: e.target.value })}
                                    className="border-none bg-transparent h-auto p-0 font-medium text-slate-700 focus:ring-0"
                                    placeholder={`Option ${idx + 1}`}
                                />
                                <IconButton 
                                    icon={X} 
                                    size="sm" 
                                    className="text-slate-300 hover:text-red-500" 
                                    onClick={() => updateActiveQuestion({ options: activeQuestion.options.filter(o => o.id !== opt.id) })}
                                />
                            </div>
                        ))}
                        <Button 
                            variant="outline" 
                            className="w-full border-dashed" 
                            icon={Plus} 
                            onClick={() => updateActiveQuestion({ 
                                options: [...activeQuestion.options, { id: `o${Date.now()}`, text: '', isCorrect: false }] 
                            })}
                        >
                            Add Option
                        </Button>
                    </div>

                    {/* LIVE PREVIEW AREA */}
                    <div className="pt-6 border-t border-slate-100 mt-6 space-y-4">
                        <div className="flex items-center justify-between">
                            <label className="text-[10px] font-black text-indigo-500 uppercase tracking-widest flex items-center gap-2">
                                <Eye size={12} /> Student View Preview
                            </label>
                            <Badge variant="outline" className="text-[9px] uppercase border-indigo-100 text-indigo-400">Live</Badge>
                        </div>
                        <div className="bg-slate-50 border border-slate-200 rounded-2xl p-6 space-y-4">
                            <div className="text-md font-bold text-slate-900 mb-2">{activeQuestion.text || 'Question text will appear here...'}</div>
                            <div className="space-y-2">
                                {activeQuestion.options.map((opt) => (
                                    <div key={opt.id} className="flex items-center gap-3 p-3 bg-white border border-slate-200 rounded-xl hover:border-indigo-300 transition-colors">
                                        <div className={cn(
                                            "w-5 h-5 rounded-full border-2 flex-shrink-0 border-slate-300",
                                            activeQuestion.type === 'multiple_choice' ? "rounded-md" : "rounded-full"
                                        )} />
                                        <span className="text-sm font-medium text-slate-700">{opt.text || <span className="text-slate-300 italic">Option text...</span>}</span>
                                    </div>
                                ))}
                                {activeQuestion.options.length === 0 && (
                                    <div className="text-xs text-slate-400 italic text-center py-4">Add options to see preview</div>
                                )}
                            </div>
                        </div>
                    </div>
                 </div>

                 {/* Settings & Explanation */}
                 <div className="bg-white rounded-2xl shadow-sm border border-slate-200 p-6 space-y-4">
                     <div>
                        <label className="text-[10px] font-black text-slate-400 uppercase tracking-widest flex items-center gap-2 mb-2">
                            <Target size={12} /> Correct Answer Explanation
                        </label>
                        <textarea 
                            value={activeQuestion.explanation || ''}
                            onChange={(e) => updateActiveQuestion({ explanation: e.target.value })}
                            className="w-full bg-slate-50 border border-slate-200 rounded-xl p-3 text-sm focus:ring-2 focus:ring-indigo-100 outline-none"
                            placeholder="Explain why the correct answer is correct..."
                            rows={3}
                        />
                     </div>
                     <div className="flex items-center justify-between border-t border-slate-100 pt-4">
                         <div className="flex items-center gap-2">
                             <span className="text-xs font-bold text-slate-600">Points Value</span>
                             <Input 
                                type="number" 
                                className="w-20 h-8 font-bold" 
                                value={activeQuestion.points}
                                onChange={(e: any) => updateActiveQuestion({ points: parseInt(e.target.value) || 0 })}
                             />
                         </div>
                     </div>
                 </div>
              </div>
           ) : (
                <div className="flex flex-col items-center justify-center h-full text-slate-400">
                    <Layout size={48} className="opacity-20 mb-4" />
                    <p>Select or add a question to begin.</p>
                </div>
           )}
        </div>
      </div>
    </div>
  );
};
