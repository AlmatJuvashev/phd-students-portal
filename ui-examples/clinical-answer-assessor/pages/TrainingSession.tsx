
import React, { useState, useMemo } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { TOPICS, PROTOCOLS, MOCK_ANSWERS } from '../constants';
import { Card } from '../components/ui/Card';
import { ArrowLeft, BookOpen, Send, Clock, CheckCircle, History, List, AlertTriangle, ChevronRight, Check } from 'lucide-react';

const TrainingSession: React.FC = () => {
  const { topicId } = useParams();
  const navigate = useNavigate();
  const topic = TOPICS.find(t => t.id === topicId);
  const question = topic?.questions[0]; // Just take first for demo

  // Mock User Identity for history lookup
  const CURRENT_USER_ID = 101;

  // State
  const [ragEnabled, setRagEnabled] = useState(true);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [result, setResult] = useState<any>(null);
  
  // Single Answer State (Legacy)
  const [singleAnswer, setSingleAnswer] = useState('');
  
  // Rubric Answer State (New)
  const [rubricAnswers, setRubricAnswers] = useState<Record<string, string>>({});

  // Helper for Rubric
  const isRubricMode = question?.rubric && question.rubric.length > 0;
  
  // Computed
  const completedRubricCount = isRubricMode 
    ? Object.keys(rubricAnswers).filter(k => rubricAnswers[k]?.trim().length > 0).length 
    : 0;
  
  const totalRubricCount = question?.rubric?.length || 0;
  const progressPercent = isRubricMode ? (completedRubricCount / totalRubricCount) * 100 : 0;

  // Retrieve past attempts for this specific question
  const pastAttempts = useMemo(() => {
    if (!question) return [];
    return MOCK_ANSWERS
      .filter(a => a.user_id === CURRENT_USER_ID && a.question === question.text)
      .sort((a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime());
  }, [question]);

  if (!topic || !question) return <div>Topic not found</div>;

  const handleRubricChange = (id: string, text: string) => {
    setRubricAnswers(prev => ({ ...prev, [id]: text }));
  };

  const scrollToSection = (id: string) => {
    const el = document.getElementById(`section-${id}`);
    if (el) {
      el.scrollIntoView({ behavior: 'smooth', block: 'start' });
    }
  };

  const handleSubmit = () => {
    setIsSubmitting(true);
    // Mock API delay
    setTimeout(() => {
      setIsSubmitting(false);
      
      const mockResult: any = {
        score: isRubricMode ? 82 : 78,
        strengths: ["Clinical reasoning aligns with protocols", "Correct identification of symptoms"],
        weaknesses: ["Missing detail on drug dosage frequencies"],
        rag_snippets: [
          { title: "NRCHD Protocol #12", text: "Infusion therapy involves crystalloids at a rate of 40-50 ml/kg/day..." }
        ]
      };

      if (isRubricMode && question.rubric) {
        // Generate mock per-section scores
        mockResult.sectionScores = question.rubric.map(r => {
           // Randomish score based on length of input
           const ansLen = rubricAnswers[r.id]?.length || 0;
           const scorePercent = ansLen > 50 ? 0.8 + (Math.random() * 0.2) : 0.4;
           return {
             id: r.id,
             score: Math.min(r.max_score, Math.round(r.max_score * scorePercent)),
             max: r.max_score,
             feedback: ansLen > 50 ? "Satisfactory coverage of criteria." : "Response is too brief or incomplete."
           };
        });
        mockResult.score = mockResult.sectionScores.reduce((acc: number, curr: any) => acc + curr.score, 0);
      }

      setResult(mockResult);
      window.scrollTo({ top: 0, behavior: 'smooth' });
    }, 2000);
  };

  return (
    <div className="max-w-7xl mx-auto h-[calc(100vh-8rem)] flex flex-col">
      {/* Header */}
      <div className="flex items-center mb-4 shrink-0 px-4 md:px-0">
        <button onClick={() => navigate('/training')} className="p-2 mr-4 hover:bg-slate-100 rounded-full text-slate-500">
          <ArrowLeft className="h-5 w-5" />
        </button>
        <div className="flex-1">
           <div className="flex items-center gap-2">
              <span className="text-xs font-semibold text-teal-600 uppercase tracking-wide">{topic.title}</span>
              <span className="text-slate-300">|</span>
              <span className="text-xs text-slate-500">{topic.questions.length} Questions</span>
           </div>
           <h2 className="text-xl font-bold text-slate-900 mt-1 line-clamp-1">{question.text}</h2>
        </div>
        <div className="hidden md:flex items-center space-x-4">
           <div className="flex items-center text-sm text-slate-500 bg-white px-3 py-1.5 rounded-full border shadow-sm">
              <Clock className="h-4 w-4 mr-2" />
              {question.estimated_time_mins} min
           </div>
           <div className="flex items-center text-sm text-slate-500 bg-white px-3 py-1.5 rounded-full border shadow-sm">
              <span className={`h-2 w-2 rounded-full mr-2 ${question.difficulty === 'Beginner' ? 'bg-green-500' : 'bg-orange-500'}`} />
              {question.difficulty}
           </div>
        </div>
      </div>

      {/* Main Content Area */}
      <div className="flex-1 flex gap-6 min-h-0 px-4 md:px-0 pb-4">
         
         {/* LEFT COLUMN: Sidebar Navigation (Only for Rubric Mode) */}
         {isRubricMode && (
            <div className="hidden lg:flex w-64 flex-col gap-4 overflow-hidden shrink-0">
              {/* Progress Card */}
              <div className="bg-white rounded-xl border border-slate-200 shadow-sm p-4">
                  <div className="flex justify-between items-center mb-2">
                     <span className="text-xs font-semibold text-slate-500 uppercase">Progress</span>
                     <span className="text-xs font-bold text-teal-600">{Math.round(progressPercent)}%</span>
                  </div>
                  <div className="w-full bg-slate-100 rounded-full h-2 mb-4">
                     <div 
                       className="bg-teal-500 h-2 rounded-full transition-all duration-500 ease-out" 
                       style={{ width: `${progressPercent}%` }}
                     ></div>
                  </div>
                  <button 
                     onClick={handleSubmit}
                     disabled={completedRubricCount < totalRubricCount || isSubmitting || !!result}
                     className={`w-full py-2 px-4 rounded-lg text-sm font-medium flex items-center justify-center transition-all
                        ${result ? 'bg-green-100 text-green-800' : 
                          completedRubricCount === totalRubricCount ? 'bg-teal-600 text-white hover:bg-teal-700 shadow-sm' : 
                          'bg-slate-100 text-slate-400 cursor-not-allowed'}
                     `}
                  >
                     {result ? 'Submitted' : 'Submit Exam'}
                  </button>
              </div>

              {/* Table of Contents */}
              <div className="flex-1 bg-white rounded-xl border border-slate-200 shadow-sm overflow-y-auto custom-scrollbar">
                  <div className="p-3 border-b border-slate-100 bg-slate-50">
                     <h3 className="text-xs font-bold text-slate-500 uppercase flex items-center">
                        <List className="h-3 w-3 mr-2" /> Exam Sections
                     </h3>
                  </div>
                  <div className="py-2">
                     {question.rubric!.map((item, idx) => {
                        const isFilled = (rubricAnswers[item.id] || '').length > 0;
                        const scoreData = result?.sectionScores?.find((s:any) => s.id === item.id);
                        
                        return (
                           <button
                              key={item.id}
                              onClick={() => scrollToSection(item.id)}
                              className="w-full text-left px-4 py-2 text-sm hover:bg-slate-50 transition-colors flex items-start group relative"
                           >
                              <div className={`mt-0.5 w-4 h-4 rounded border flex items-center justify-center mr-3 shrink-0 text-[10px] 
                                 ${result 
                                   ? (scoreData?.score === item.max_score ? 'bg-green-50 border-green-200 text-green-600' : 'bg-amber-50 border-amber-200 text-amber-600')
                                   : (isFilled ? 'bg-teal-50 border-teal-200 text-teal-600' : 'border-slate-300 text-slate-400')
                                 }
                              `}>
                                 {result ? (scoreData?.score === item.max_score ? <Check className="h-3 w-3" /> : '!') : (isFilled ? <Check className="h-3 w-3" /> : (idx + 1))}
                              </div>
                              <div className="flex-1">
                                 <div className={`line-clamp-1 font-medium ${isFilled ? 'text-slate-900' : 'text-slate-500'}`}>
                                    {item.criteria}
                                 </div>
                                 {result && (
                                    <div className="text-xs text-slate-400 mt-0.5">
                                       Score: <span className="font-semibold text-slate-700">{scoreData?.score}</span>/{item.max_score}
                                    </div>
                                 )}
                              </div>
                           </button>
                        );
                     })}
                  </div>
              </div>
            </div>
         )}

         {/* CENTER COLUMN: Content & Form */}
         <div className="flex-1 flex flex-col gap-6 overflow-y-auto custom-scrollbar pr-1">
            
            {/* Clinical Scenario Card */}
            {(!result) && (
              <Card className="bg-blue-50/50 border-blue-100 shrink-0">
                 <div className="flex items-start">
                    <BookOpen className="h-5 w-5 text-blue-600 mt-0.5 mr-3 shrink-0" />
                    <div>
                       <h3 className="text-sm font-bold text-blue-900 uppercase tracking-wide mb-1">Clinical Scenario</h3>
                       <p className="text-slate-800 leading-relaxed text-lg">{question.text}</p>
                    </div>
                 </div>
              </Card>
            )}

            {/* Score Summary Banner */}
            {result && (
               <div className="bg-white rounded-xl border border-slate-200 shadow-sm p-6 text-center animate-in fade-in slide-in-from-top-4 duration-500">
                  <div className="inline-flex items-center justify-center p-3 bg-teal-50 rounded-full mb-4">
                     <span className="text-4xl font-bold text-teal-600">{result.score}</span>
                     <span className="text-sm text-teal-400 ml-1 font-medium self-end mb-1">/ 100</span>
                  </div>
                  <h3 className="text-xl font-bold text-slate-900">Assessment Complete</h3>
                  <p className="text-slate-500 max-w-lg mx-auto mt-2">
                     Review the breakdown below. The AI has evaluated each section against the clinical protocol criteria.
                  </p>
                  <button 
                     onClick={() => { setResult(null); setRubricAnswers({}); setSingleAnswer(''); }}
                     className="mt-6 px-6 py-2 bg-slate-900 text-white rounded-lg font-medium hover:bg-slate-800 transition-colors"
                  >
                     Start New Practice Session
                  </button>
               </div>
            )}

            {/* Form Area */}
            {isRubricMode ? (
               <div className="space-y-6 pb-20">
                  {question.rubric!.map((item, index) => {
                     const val = rubricAnswers[item.id] || '';
                     const wordCount = val.trim() === '' ? 0 : val.trim().split(/\s+/).length;
                     const scoreData = result?.sectionScores?.find((s:any) => s.id === item.id);
                     
                     return (
                        <div key={item.id} id={`section-${item.id}`} className="scroll-mt-6">
                           <Card className={`transition-all duration-300 ${scoreData ? (scoreData.score === item.max_score ? 'border-green-200 ring-1 ring-green-100' : 'border-amber-200 ring-1 ring-amber-100') : 'border-slate-200 focus-within:ring-2 ring-teal-50 focus-within:border-teal-400'}`}>
                              <div className="border-b border-slate-100 px-5 py-4 bg-slate-50/50 flex flex-col md:flex-row md:items-center justify-between gap-2">
                                 <div className="flex items-start md:items-center gap-3">
                                    <span className="flex items-center justify-center h-6 w-6 rounded bg-white border border-slate-200 text-xs font-bold text-slate-500 shrink-0 shadow-sm">
                                       {index + 1}
                                    </span>
                                    <div>
                                       <h4 className="text-base font-bold text-slate-800">{item.criteria}</h4>
                                       <p className="text-xs text-slate-500 mt-0.5 max-w-xl">{item.description}</p>
                                    </div>
                                 </div>
                                 <div className="flex items-center gap-3 shrink-0 ml-9 md:ml-0">
                                    <span className="text-xs font-medium text-slate-400 bg-white px-2 py-1 rounded border border-slate-200">
                                       Max Score: {item.max_score}
                                    </span>
                                    {scoreData && (
                                       <span className={`text-sm font-bold px-3 py-1 rounded-full ${scoreData.score === item.max_score ? 'bg-green-100 text-green-700' : 'bg-amber-100 text-amber-700'}`}>
                                          {scoreData.score} / {item.max_score}
                                       </span>
                                    )}
                                 </div>
                              </div>
                              
                              <div className="p-0">
                                 {result ? (
                                    <div className="p-5 bg-slate-50/30">
                                       <p className="text-slate-800 whitespace-pre-wrap mb-4 font-serif text-sm leading-relaxed">{val}</p>
                                       <div className="mt-4 p-3 bg-white rounded-lg border border-slate-100 shadow-sm">
                                          <div className="flex items-center gap-2 mb-1">
                                             <div className={`h-2 w-2 rounded-full ${scoreData.score === item.max_score ? 'bg-green-500' : 'bg-amber-500'}`}></div>
                                             <span className="text-xs font-bold text-slate-700 uppercase">AI Feedback</span>
                                          </div>
                                          <p className="text-sm text-slate-600">{scoreData.feedback}</p>
                                       </div>
                                    </div>
                                 ) : (
                                    <>
                                       <textarea
                                          className="w-full p-5 focus:outline-none resize-y min-h-[140px] text-slate-700 leading-relaxed placeholder:text-slate-300"
                                          placeholder={`Enter your answer for ${item.criteria.toLowerCase()}...`}
                                          value={val}
                                          onChange={(e) => handleRubricChange(item.id, e.target.value)}
                                          disabled={isSubmitting || !!result}
                                       />
                                       <div className="px-5 py-2 bg-white flex justify-end items-center text-xs text-slate-300 border-t border-slate-50">
                                          <span>{wordCount} words</span>
                                       </div>
                                    </>
                                 )}
                              </div>
                           </Card>
                        </div>
                     );
                  })}
                  
                  {/* Bottom Submit Area for Mobile/Tablet or end of form */}
                  {!result && (
                     <div className="flex justify-end pt-4">
                        <button 
                           onClick={handleSubmit}
                           disabled={completedRubricCount < totalRubricCount || isSubmitting}
                           className={`flex items-center px-8 py-3 rounded-xl font-bold text-white shadow-md transition-all transform hover:-translate-y-0.5
                              ${completedRubricCount < totalRubricCount || isSubmitting ? 'bg-slate-300 cursor-not-allowed shadow-none translate-y-0' : 'bg-teal-600 hover:bg-teal-700 hover:shadow-lg'}
                           `}
                        >
                           {isSubmitting ? 'Analyzing...' : `Submit Assessment (${completedRubricCount}/${totalRubricCount})`}
                           {!isSubmitting && <Send className="h-4 w-4 ml-2" />}
                        </button>
                     </div>
                  )}
               </div>
            ) : (
               // LEGACY SINGLE TEXTAREA MODE
               <Card className="flex flex-col h-full min-h-[400px]">
                  {result ? (
                      <div className="p-8">
                         <h4 className="text-lg font-bold text-slate-900 mb-2">Your Answer</h4>
                         <p className="text-slate-700 whitespace-pre-wrap mb-6 p-4 bg-slate-50 rounded-lg border border-slate-100">{singleAnswer}</p>
                      </div>
                  ) : (
                     <>
                        <div className="border-b border-slate-100 px-4 py-3 bg-slate-50 flex justify-between items-center">
                           <span className="text-xs font-medium text-slate-500 uppercase tracking-wide">Free Text Response</span>
                           <div className="flex items-center gap-2">
                               <button 
                                 onClick={() => setRagEnabled(!ragEnabled)}
                                 className={`flex items-center text-xs px-2.5 py-1 rounded-md transition-colors border ${ragEnabled ? 'bg-teal-50 border-teal-200 text-teal-700' : 'bg-white border-slate-200 text-slate-500'}`}
                               >
                                 <BookOpen className="h-3 w-3 mr-1.5" />
                                 Protocol RAG {ragEnabled ? 'ON' : 'OFF'}
                               </button>
                           </div>
                        </div>
                        <textarea
                           className="flex-1 w-full p-6 focus:outline-none resize-none text-slate-700 text-lg leading-relaxed"
                           placeholder="Type your clinical reasoning here..."
                           value={singleAnswer}
                           onChange={(e) => setSingleAnswer(e.target.value)}
                           disabled={isSubmitting}
                        />
                        <div className="px-6 py-3 bg-white flex justify-between items-center border-t border-slate-100">
                           <div className="text-xs text-slate-400">
                              {singleAnswer.trim() === '' ? 0 : singleAnswer.trim().split(/\s+/).length} words
                           </div>
                           <button 
                              onClick={handleSubmit}
                              disabled={!singleAnswer.trim() || isSubmitting}
                              className={`flex items-center px-6 py-2 rounded-lg font-medium text-white shadow-sm transition-all
                                 ${!singleAnswer.trim() || isSubmitting ? 'bg-slate-300 cursor-not-allowed' : 'bg-teal-600 hover:bg-teal-700 hover:shadow-md'}
                              `}
                           >
                              {isSubmitting ? 'Analyzing...' : 'Submit Answer'}
                           </button>
                        </div>
                     </>
                  )}
               </Card>
            )}
         </div>

         {/* RIGHT COLUMN: Feedback & Context (Sticky) */}
         <div className="hidden lg:flex lg:col-span-1 w-80 shrink-0 flex-col gap-4 overflow-y-auto custom-scrollbar h-full pl-2">
            
            {result ? (
               <div className="space-y-4 animate-in slide-in-from-right duration-500">
                  <Card title="Global Assessment">
                     <div className="space-y-4">
                        <div>
                           <h4 className="text-xs font-bold text-green-700 mb-2 flex items-center uppercase tracking-wide"><CheckCircle className="h-3 w-3 mr-1" /> Strengths</h4>
                           <ul className="text-sm text-slate-600 space-y-2">
                              {result.strengths.map((s: string, i: number) => (
                                 <li key={i} className="flex items-start"><span className="mr-2">•</span> {s}</li>
                              ))}
                           </ul>
                        </div>
                        <div className="pt-3 border-t border-slate-100">
                           <h4 className="text-xs font-bold text-amber-700 mb-2 flex items-center uppercase tracking-wide"><AlertTriangle className="h-3 w-3 mr-1" /> Improvements</h4>
                           <ul className="text-sm text-slate-600 space-y-2">
                              {result.weaknesses.map((w: string, i: number) => (
                                 <li key={i} className="flex items-start"><span className="mr-2">•</span> {w}</li>
                              ))}
                           </ul>
                        </div>
                     </div>
                  </Card>
                  
                  <Card title="Protocol References" className="bg-slate-50 border-slate-200">
                     {result.rag_snippets.map((snip: any, i: number) => (
                        <div key={i} className="text-sm">
                           <p className="font-semibold text-slate-800 mb-1 flex items-center gap-1">
                              <BookOpen className="h-3 w-3 text-slate-400" />
                              {snip.title}
                           </p>
                           <p className="text-slate-600 italic border-l-2 border-teal-300 pl-3 py-2 bg-white rounded-r text-xs leading-relaxed">
                              "{snip.text}"
                           </p>
                        </div>
                     ))}
                  </Card>
               </div>
            ) : (
               <div className="space-y-4">
                   {!isRubricMode && (
                      <Card className="bg-blue-50 border-blue-100 p-4">
                        <h4 className="font-bold text-blue-900 text-sm mb-1">Instructions</h4>
                        <p className="text-blue-800 text-xs leading-relaxed">
                           Provide a comprehensive answer based on clinical protocols. You can enable RAG to simulate reference availability.
                        </p>
                      </Card>
                   )}

                   {/* Protocol Status */}
                   <Card title="Knowledge Base" className="bg-slate-50">
                      <div className="text-xs text-slate-500 mb-2">Active protocols for this session:</div>
                      <div className="space-y-2">
                         {PROTOCOLS.filter(p => p.active).slice(0, 3).map(p => (
                            <div key={p.id} className="flex items-center text-xs font-medium text-slate-700 bg-white px-3 py-2 rounded border border-slate-200 shadow-sm">
                               <div className="h-1.5 w-1.5 rounded-full bg-teal-500 mr-2 shadow-[0_0_4px_rgba(20,184,166,0.5)]" />
                               {p.name}
                            </div>
                         ))}
                      </div>
                   </Card>
               </div>
            )}

            {/* History Section */}
            {pastAttempts.length > 0 && (
               <Card title="Previous Attempts">
                  <div className="divide-y divide-slate-100 -mx-6 px-6">
                     {pastAttempts.map((attempt, i) => (
                        <div key={i} className="py-3 flex justify-between items-center group cursor-pointer hover:bg-slate-50 transition-colors">
                           <div>
                              <div className="flex items-center text-xs font-semibold text-slate-700">
                                 <History className="h-3 w-3 text-slate-400 mr-1.5" />
                                 {new Date(attempt.created_at).toLocaleDateString()}
                              </div>
                           </div>
                           <span className={`inline-flex items-center px-2 py-0.5 rounded text-[10px] font-bold border
                              ${attempt.score >= 80 ? 'bg-green-50 text-green-700 border-green-100' :
                                attempt.score >= 50 ? 'bg-yellow-50 text-yellow-700 border-yellow-100' :
                                'bg-red-50 text-red-700 border-red-100'}
                           `}>
                              {attempt.score}%
                           </span>
                        </div>
                     ))}
                  </div>
               </Card>
            )}
         </div>

      </div>
    </div>
  );
};

export default TrainingSession;
