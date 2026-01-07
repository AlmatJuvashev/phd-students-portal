import React, { useState } from 'react';
import { useMutation } from '@tanstack/react-query';
import { Loader2, Wand2, Check, AlertCircle, Copy, FileJson } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { Textarea } from '@/components/ui/textarea';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { generateCourseStructure, generateQuiz, generateSurvey, generateAssessmentItems } from '../api';
import { Module } from '../types';
import { useToast } from '@/components/ui/use-toast';

interface AIAssistantModalProps {
  isOpen: boolean;
  onClose: () => void;
  onApply: (modules: Module[]) => void;
}

export const AIAssistantModal: React.FC<AIAssistantModalProps> = ({ isOpen, onClose, onApply }) => {
  const { toast } = useToast();
  const [activeTab, setActiveTab] = useState('structure');
  
  // Structure State
  const [syllabus, setSyllabus] = useState('');
  const [generatedModules, setGeneratedModules] = useState<Module[] | null>(null);

  // Quiz State
  const [quizTopic, setQuizTopic] = useState('');
  const [quizDifficulty, setQuizDifficulty] = useState('Medium');
  const [quizCount, setQuizCount] = useState(5);
  const [generatedQuiz, setGeneratedQuiz] = useState<any>(null);

  // Survey State
  const [surveyTopic, setSurveyTopic] = useState('');
  const [surveyCount, setSurveyCount] = useState(5);
  const [generatedSurvey, setGeneratedSurvey] = useState<any>(null);

  // Items State
  const [itemsTopic, setItemsTopic] = useState('');
  const [itemsType, setItemsType] = useState('multiple_choice');
  const [itemsCount, setItemsCount] = useState(5);
  const [generatedItems, setGeneratedItems] = useState<any[] | null>(null);

  // --- Mutations ---

  const structureMutation = useMutation({
    mutationFn: generateCourseStructure,
    onSuccess: (response: { modules: any[] }) => {
      // API returns { modules: [...] } directly
      const mappedModules = response.modules.map((m: any, mIdx: number) => ({
        ...m,
        id: `m_ai_${Date.now()}_${mIdx}`,
        lessons: m.lessons?.map((l: any, lIdx: number) => ({
          ...l,
          id: `l_ai_${Date.now()}_${mIdx}_${lIdx}`,
          activities: l.activities?.map((a: any, aIdx: number) => ({
            ...a,
            id: `a_ai_${Date.now()}_${mIdx}_${lIdx}_${aIdx}`,
            type: validateActivityType(a.type),
            points: a.points || 0,
            is_optional: a.is_optional || false,
            content: '',
            attachments: [],
            citations: []
          })) || []
        })) || []
      }));
      setGeneratedModules(mappedModules);
    }
  });

  const quizMutation = useMutation({
    mutationFn: () => generateQuiz(quizTopic, quizDifficulty, quizCount),
    onSuccess: (res) => setGeneratedQuiz(res)
  });

  const surveyMutation = useMutation({
    mutationFn: () => generateSurvey(surveyTopic, surveyCount),
    onSuccess: (res) => setGeneratedSurvey(res)
  });

  const itemsMutation = useMutation({
    mutationFn: () => generateAssessmentItems(itemsTopic, itemsType, itemsCount),
    // Response is { items: [...] }
    onSuccess: (res) => setGeneratedItems(res.items)
  });

  // --- Helpers ---

  const validateActivityType = (type: string): string => {
    const validTypes = ['text', 'video', 'quiz', 'survey', 'assignment', 'resource', 'live'];
    return validTypes.includes(type) ? type : 'text';
  };

  const activeMutation = () => {
    switch(activeTab) {
      case 'structure': return structureMutation;
      case 'quiz': return quizMutation;
      case 'survey': return surveyMutation;
      case 'items': return itemsMutation;
      default: return structureMutation;
    }
  };

  const handleGenerate = () => {
    switch(activeTab) {
      case 'structure': structureMutation.mutate(syllabus); break;
      case 'quiz': quizMutation.mutate(); break;
      case 'survey': surveyMutation.mutate(); break;
      case 'items': itemsMutation.mutate(); break;
    }
  };

  const handleApplyStructure = () => {
    if (generatedModules) {
      onApply(generatedModules);
      onClose();
    }
  };

  const copyToClipboard = (data: any) => {
    navigator.clipboard.writeText(JSON.stringify(data, null, 2));
    toast({
      title: "Copied",
      description: "JSON copied to clipboard"
    });
  };

  const reset = () => {
    setSyllabus('');
    setGeneratedModules(null);
    setGeneratedQuiz(null);
    setGeneratedSurvey(null);
    setGeneratedItems(null);
    structureMutation.reset();
    quizMutation.reset();
    surveyMutation.reset();
    itemsMutation.reset();
  };

  return (
    <Dialog open={isOpen} onOpenChange={(open) => !open && onClose()}>
      <DialogContent className="sm:max-w-[750px] max-h-[85vh] flex flex-col">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2">
            <Wand2 className="w-5 h-5 text-indigo-500" />
            AI Content Assistant
          </DialogTitle>
          <DialogDescription>
            Generate course structures, quizzes, surveys, and assessment items using AI.
          </DialogDescription>
        </DialogHeader>

        <Tabs value={activeTab} onValueChange={setActiveTab} className="flex-1 flex flex-col min-h-0">
          <TabsList className="grid w-full grid-cols-4">
            <TabsTrigger value="structure">Structure</TabsTrigger>
            <TabsTrigger value="quiz">Quiz</TabsTrigger>
            <TabsTrigger value="survey">Survey</TabsTrigger>
            <TabsTrigger value="items">Item Bank</TabsTrigger>
          </TabsList>
          
          <div className="flex-1 py-4 overflow-y-auto min-h-[350px]">
            {/* --- COURSE STRUCTURE --- */}
            <TabsContent value="structure" className="mt-0 h-full flex flex-col">
              {generatedModules ? (
                <div className="space-y-4">
                   <div className="bg-emerald-50 border border-emerald-200 rounded-md p-4 flex gap-2">
                      <Check className="w-5 h-5 text-emerald-600 mt-0.5" />
                      <div className="text-sm text-emerald-800">
                        <p className="font-bold">Structure Generated!</p>
                        <p>Review the modules below before applying.</p>
                      </div>
                   </div>
                   <div className="space-y-4">
                     {generatedModules.map((m) => (
                       <div key={m.id} className="border p-3 rounded bg-slate-50">
                         <div className="font-bold">{m.title}</div>
                         <div className="pl-4 mt-2 space-y-1">
                           {m.lessons.map(l => (
                             <div key={l.id} className="text-sm text-slate-600 flex items-center gap-2">
                               • {l.title} 
                               <span className="text-xs text-slate-400">({l.activities.length} activities)</span>
                             </div>
                           ))}
                         </div>
                       </div>
                     ))}
                   </div>
                </div>
              ) : (
                <div className="space-y-3">
                  <Label>Syllabus / Outline Text</Label>
                  <Textarea 
                    value={syllabus} 
                    onChange={e => setSyllabus(e.target.value)}
                    placeholder="Paste your syllabus here..." 
                    className="flex-1 min-h-[250px] font-mono"
                  />
                  <p className="text-xs text-slate-500">
                    The AI will parse this text and propose a module structure.
                  </p>
                </div>
              )}
            </TabsContent>

            {/* --- QUIZ --- */}
            <TabsContent value="quiz" className="mt-0 space-y-4">
              {generatedQuiz ? (
                 <div className="space-y-4">
                    <div className="bg-slate-900 text-slate-50 p-4 rounded-md font-mono text-xs overflow-auto max-h-[400px]">
                       <pre>{JSON.stringify(generatedQuiz, null, 2)}</pre>
                    </div>
                    <Button variant="outline" className="w-full" onClick={() => copyToClipboard(generatedQuiz)}>
                       <Copy className="w-4 h-4 mr-2" /> Copy JSON
                    </Button>
                 </div>
              ) : (
                <div className="space-y-4">
                  <div className="grid grid-cols-2 gap-4">
                    <div className="space-y-2">
                      <Label>Topic</Label>
                      <Input value={quizTopic} onChange={e => setQuizTopic(e.target.value)} placeholder="e.g. Data Ethics" />
                    </div>
                    <div className="space-y-2">
                       <Label>Difficulty</Label>
                       <Select value={quizDifficulty} onValueChange={setQuizDifficulty}>
                         <SelectTrigger><SelectValue /></SelectTrigger>
                         <SelectContent>
                           <SelectItem value="Easy">Easy</SelectItem>
                           <SelectItem value="Medium">Medium</SelectItem>
                           <SelectItem value="Hard">Hard</SelectItem>
                         </SelectContent>
                       </Select>
                    </div>
                    <div className="space-y-2">
                      <Label>Question Count</Label>
                      <Input type="number" value={quizCount} onChange={e => setQuizCount(Number(e.target.value))} min={1} max={20} />
                    </div>
                  </div>
                </div>
              )}
            </TabsContent>

            {/* --- SURVEY --- */}
            <TabsContent value="survey" className="mt-0 space-y-4">
              {generatedSurvey ? (
                 <div className="space-y-4">
                    <div className="bg-slate-900 text-slate-50 p-4 rounded-md font-mono text-xs overflow-auto max-h-[400px]">
                       <pre>{JSON.stringify(generatedSurvey, null, 2)}</pre>
                    </div>
                    <Button variant="outline" className="w-full" onClick={() => copyToClipboard(generatedSurvey)}>
                       <Copy className="w-4 h-4 mr-2" /> Copy JSON
                    </Button>
                 </div>
              ) : (
                <div className="space-y-4">
                   <div className="space-y-2">
                      <Label>Survey Topic</Label>
                      <Input value={surveyTopic} onChange={e => setSurveyTopic(e.target.value)} placeholder="e.g. Course Feedback" />
                    </div>
                    <div className="space-y-2">
                      <Label>Question Count</Label>
                      <Input type="number" value={surveyCount} onChange={e => setSurveyCount(Number(e.target.value))} min={1} max={20} />
                    </div>
                </div>
              )}
            </TabsContent>

            {/* --- ITEMS --- */}
            <TabsContent value="items" className="mt-0 space-y-4">
               {generatedItems ? (
                 <div className="space-y-4">
                    <div className="bg-slate-50 border rounded-md p-4 space-y-4 max-h-[400px] overflow-auto">
                       {generatedItems.map((item, idx) => (
                         <div key={idx} className="bg-white border p-3 rounded shadow-sm text-sm">
                            <div className="flex items-center justify-between mb-2">
                               <span className="font-bold uppercase text-xs text-slate-500">{item.type}</span>
                               <span className="text-xs bg-slate-100 px-1.5 py-0.5 rounded">Diff: {item.difficulty}</span>
                            </div>
                            <pre className="whitespace-pre-wrap font-sans text-slate-800">
                               {JSON.stringify(item.content, null, 2)}
                            </pre>
                         </div>
                       ))}
                    </div>
                    <Button variant="outline" className="w-full" onClick={() => copyToClipboard(generatedItems)}>
                       <Copy className="w-4 h-4 mr-2" /> Copy All Items JSON
                    </Button>
                 </div>
               ) : (
                 <div className="space-y-4">
                   <div className="grid grid-cols-2 gap-4">
                    <div className="space-y-2">
                      <Label>Topic</Label>
                      <Input value={itemsTopic} onChange={e => setItemsTopic(e.target.value)} placeholder="e.g. Cell Biology" />
                    </div>
                    <div className="space-y-2">
                       <Label>Item Type</Label>
                       <Select value={itemsType} onValueChange={setItemsType}>
                         <SelectTrigger><SelectValue /></SelectTrigger>
                         <SelectContent>
                           <SelectItem value="multiple_choice">Multiple Choice</SelectItem>
                           <SelectItem value="true_false">True / False</SelectItem>
                           <SelectItem value="essay">Essay</SelectItem>
                         </SelectContent>
                       </Select>
                    </div>
                    <div className="space-y-2">
                      <Label>Count</Label>
                      <Input type="number" value={itemsCount} onChange={e => setItemsCount(Number(e.target.value))} min={1} max={10} />
                    </div>
                  </div>
                 </div>
               )}
            </TabsContent>

          </div>

          <div className="mt-4 pt-4 border-t flex items-center justify-between">
             {activeMutation().isError && (
                 <div className="text-sm text-red-600 flex items-center gap-2">
                    <AlertCircle className="w-4 h-4" />
                    Error generating content.
                 </div>
             )}
             
             <div className="flex items-center gap-2 ml-auto">
               {(generatedModules || generatedQuiz || generatedSurvey || generatedItems) ? (
                 <>
                   <Button variant="ghost" onClick={reset}>Reset</Button>
                   {activeTab === 'structure' && generatedModules && (
                     <Button onClick={handleApplyStructure} className="bg-emerald-600 hover:bg-emerald-700 text-white gap-2">
                       <Check className="w-4 h-4" /> Apply to Course
                     </Button>
                   )}
                 </>
               ) : (
                 <>
                   <Button variant="ghost" onClick={onClose}>Cancel</Button>
                   <Button onClick={handleGenerate} disabled={activeMutation().isPending} className="gap-2">
                      {activeMutation().isPending ? <Loader2 className="w-4 h-4 animate-spin" /> : <Wand2 className="w-4 h-4" />}
                      Generate
                   </Button>
                 </>
               )}
             </div>
          </div>
        </Tabs>
      </DialogContent>
    </Dialog>
  );
};
