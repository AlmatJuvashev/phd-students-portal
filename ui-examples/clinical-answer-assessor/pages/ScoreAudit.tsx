
import React, { useMemo, useState } from 'react';
import { Card } from '../components/ui/Card';
import { MOCK_ANSWERS } from '../constants';
import { TeacherStat, AnswerRecord } from '../types';
import { UploadCloud, FileJson, AlertCircle, ChevronDown, ChevronRight, CheckCircle, Sliders, Users, FileText } from 'lucide-react';

const ScoreAudit: React.FC = () => {
  const [activeTab, setActiveTab] = useState<'stats' | 'discrepancy'>('stats');

  return (
    <div className="space-y-6">
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
        <div>
          <h2 className="text-2xl font-bold text-slate-900">Score Audit</h2>
          <p className="text-slate-500 mt-1">Analyze inter-rater reliability and detect scoring anomalies.</p>
        </div>
        
        {/* Tab Switcher */}
        <div className="flex p-1 bg-white border border-slate-200 rounded-lg shadow-sm">
           <button 
             onClick={() => setActiveTab('stats')}
             className={`px-4 py-2 text-sm font-medium rounded-md transition-colors flex items-center ${activeTab === 'stats' ? 'bg-teal-50 text-teal-700' : 'text-slate-600 hover:bg-slate-50'}`}
           >
             <Users className="h-4 w-4 mr-2" /> Examiner Statistics
           </button>
           <button 
             onClick={() => setActiveTab('discrepancy')}
             className={`px-4 py-2 text-sm font-medium rounded-md transition-colors flex items-center ${activeTab === 'discrepancy' ? 'bg-teal-50 text-teal-700' : 'text-slate-600 hover:bg-slate-50'}`}
           >
             <Sliders className="h-4 w-4 mr-2" /> Discrepancy Check
           </button>
        </div>
      </div>

      {activeTab === 'stats' ? <ExaminerStatsView /> : <DiscrepancyCheckView />}
    </div>
  );
};

// --------------------------------------------------------------------------------
// View 1: Examiner Stats (Original)
// --------------------------------------------------------------------------------
const ExaminerStatsView: React.FC = () => {
  const [selectedTeacherId, setSelectedTeacherId] = useState<number | null>(null);

  const teacherStats = useMemo(() => {
    const stats: Record<number, TeacherStat> = {};
    const globalTotalScore = MOCK_ANSWERS.reduce((sum, a) => sum + a.score, 0);
    const globalAvg = globalTotalScore / MOCK_ANSWERS.length;

    MOCK_ANSWERS.forEach(a => {
      if (!stats[a.examiner_id]) {
        stats[a.examiner_id] = {
          examiner_id: a.examiner_id,
          examiner_name: a.examiner_name,
          answers_count: 0,
          average_score: 0,
          std_dev: 0,
          deviation_from_global: 0,
          status: 'neutral'
        };
      }
      stats[a.examiner_id].answers_count++;
      stats[a.examiner_id].average_score += a.score;
    });

    Object.values(stats).forEach(stat => {
      stat.average_score = stat.average_score / stat.answers_count;
      stat.deviation_from_global = stat.average_score - globalAvg;
      if (stat.deviation_from_global > 10) stat.status = 'lenient';
      else if (stat.deviation_from_global < -10) stat.status = 'strict';
    });

    return Object.values(stats);
  }, []);

  const selectedTeacherAnswers = useMemo(() => {
    if (!selectedTeacherId) return [];
    return MOCK_ANSWERS.filter(a => a.examiner_id === selectedTeacherId);
  }, [selectedTeacherId]);

  return (
    <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Left: Teachers List */}
        <div className="lg:col-span-1">
          <Card title="Examiners" className="h-full">
            <div className="space-y-2 max-h-[600px] overflow-y-auto pr-2 custom-scrollbar">
              {teacherStats.map(stat => (
                <div 
                  key={stat.examiner_id}
                  onClick={() => setSelectedTeacherId(stat.examiner_id)}
                  className={`p-4 rounded-lg border cursor-pointer transition-colors ${
                    selectedTeacherId === stat.examiner_id 
                      ? 'bg-teal-50 border-teal-200' 
                      : 'bg-white border-slate-100 hover:border-slate-300'
                  }`}
                >
                  <div className="flex justify-between items-start">
                    <div>
                      <h4 className="font-semibold text-slate-900">{stat.examiner_name}</h4>
                      <p className="text-xs text-slate-500 mt-1">{stat.answers_count} graded answers</p>
                    </div>
                    <div className="flex flex-col items-end">
                      <span className="text-lg font-bold text-slate-700">{Math.round(stat.average_score)}</span>
                      <span className="text-xs text-slate-400">avg score</span>
                    </div>
                  </div>
                  
                  <div className="mt-3 flex items-center justify-between">
                     <span className={`text-xs px-2 py-1 rounded-full font-medium ${
                        stat.status === 'lenient' ? 'bg-red-100 text-red-800' : 
                        stat.status === 'strict' ? 'bg-blue-100 text-blue-800' : 
                        'bg-slate-100 text-slate-600'
                     }`}>
                        {stat.status.charAt(0).toUpperCase() + stat.status.slice(1)}
                     </span>
                     <span className={`text-xs font-medium ${stat.deviation_from_global > 0 ? 'text-green-600' : 'text-red-600'}`}>
                        {stat.deviation_from_global > 0 ? '+' : ''}{Math.round(stat.deviation_from_global)} vs avg
                     </span>
                  </div>
                </div>
              ))}
            </div>
          </Card>
        </div>

        {/* Right: Detailed Breakdown */}
        <div className="lg:col-span-2">
          {selectedTeacherId ? (
            <div className="space-y-6">
              <Card title={`Reviewing: ${teacherStats.find(t => t.examiner_id === selectedTeacherId)?.examiner_name}`}>
                 <div className="grid grid-cols-3 gap-4 mb-6">
                    <div className="p-4 bg-slate-50 rounded-lg text-center">
                       <div className="text-2xl font-bold text-slate-800">{selectedTeacherAnswers.length}</div>
                       <div className="text-xs text-slate-500 uppercase tracking-wide">Graded</div>
                    </div>
                    <div className="p-4 bg-slate-50 rounded-lg text-center">
                       <div className="text-2xl font-bold text-slate-800">
                          {Math.round(teacherStats.find(t => t.examiner_id === selectedTeacherId)?.average_score || 0)}
                       </div>
                       <div className="text-xs text-slate-500 uppercase tracking-wide">Avg Score</div>
                    </div>
                    <div className="p-4 bg-slate-50 rounded-lg text-center">
                       <div className="text-2xl font-bold text-amber-600">2</div>
                       <div className="text-xs text-slate-500 uppercase tracking-wide">Potential Bias</div>
                    </div>
                 </div>

                 <h4 className="text-sm font-semibold text-slate-700 mb-3">Grading History</h4>
                 <div className="space-y-3">
                    {selectedTeacherAnswers.map((answer, idx) => (
                       <AnswerRow key={idx} answer={answer} />
                    ))}
                 </div>
              </Card>
            </div>
          ) : (
            <div className="h-full min-h-[400px] flex items-center justify-center bg-slate-50 border-2 border-dashed border-slate-200 rounded-xl">
               <div className="text-center">
                  <Users className="h-12 w-12 text-slate-300 mx-auto mb-3" />
                  <p className="text-slate-500 font-medium">Select an examiner to audit grading patterns</p>
               </div>
            </div>
          )}
        </div>
      </div>
  );
};

// --------------------------------------------------------------------------------
// View 2: Discrepancy Check (New)
// --------------------------------------------------------------------------------
const DiscrepancyCheckView: React.FC = () => {
    const [auditData, setAuditData] = useState<AnswerRecord[]>([]);
    const [scoreThreshold, setScoreThreshold] = useState(2);
    const [isProcessing, setIsProcessing] = useState(false);
    
    // Group discrepancies by user_id
    const discrepancies = useMemo(() => {
        if (auditData.length === 0) return [];
        
        const flagged = auditData.filter(a => {
            const ai = a.ai_score || 0;
            return Math.abs(a.score - ai) > scoreThreshold;
        });

        // Group by user
        const grouped: Record<number, AnswerRecord[]> = {};
        flagged.forEach(a => {
            if (!grouped[a.user_id]) grouped[a.user_id] = [];
            grouped[a.user_id].push(a);
        });
        
        return Object.entries(grouped).map(([uid, answers]) => ({
            userId: Number(uid),
            answers
        }));
    }, [auditData, scoreThreshold]);

    const handleFileUpload = (e: React.ChangeEvent<HTMLInputElement>) => {
        const file = e.target.files?.[0];
        if (!file) return;

        setIsProcessing(true);
        // Simulate parsing and AI scoring
        setTimeout(() => {
            const reader = new FileReader();
            reader.onload = (event) => {
                try {
                    const json = JSON.parse(event.target?.result as string);
                    const records: any[] = Array.isArray(json) ? json : [json];
                    
                    // Transform to AnswerRecord and SIMULATE AI SCORING
                    const processed: AnswerRecord[] = records.map((r, i) => ({
                        user_id: r.user_id,
                        name_rus: r.name_rus, // Category
                        topic: r.name_rus, // Fallback topic to category
                        question: r.question,
                        answer: r.answer,
                        comment: r.comment,
                        score: r.score,
                        examiner_id: r.teacher_id,
                        examiner_name: `Examiner ${r.teacher_id}`,
                        attempt_id: `upload-${i}`,
                        created_at: new Date().toISOString(),
                        // Simulate AI Score (random deviation for demo)
                        ai_score: Math.max(0, Math.min(100, r.score + (Math.random() > 0.6 ? (Math.random() * 10 - 5) : 0))) 
                    }));
                    
                    setAuditData(processed);
                } catch (err) {
                    alert("Invalid JSON format");
                } finally {
                    setIsProcessing(false);
                }
            };
            reader.readAsText(file);
        }, 1500);
    };

    return (
        <div className="space-y-6">
            <Card>
                <div className="flex flex-col md:flex-row items-center gap-6 p-2">
                    <div className="flex-1 w-full">
                        <label className="flex flex-col items-center justify-center w-full h-32 border-2 border-slate-300 border-dashed rounded-lg cursor-pointer bg-slate-50 hover:bg-slate-100 transition-colors">
                            <div className="flex flex-col items-center justify-center pt-5 pb-6">
                                <UploadCloud className="w-8 h-8 mb-3 text-slate-400" />
                                <p className="mb-2 text-sm text-slate-500"><span className="font-semibold">Click to upload</span> student JSON dataset</p>
                                <p className="text-xs text-slate-500">.json files supported</p>
                            </div>
                            <input type="file" className="hidden" accept=".json" onChange={handleFileUpload} />
                        </label>
                    </div>
                    
                    <div className="w-full md:w-64 space-y-4 border-l border-slate-100 md:pl-6">
                        <div>
                            <label className="block text-xs font-semibold text-slate-500 uppercase tracking-wide mb-2">Score Diff Threshold</label>
                            <div className="flex items-center gap-3">
                                <input 
                                    type="range" min="1" max="20" step="1" 
                                    value={scoreThreshold} 
                                    onChange={(e) => setScoreThreshold(parseInt(e.target.value))}
                                    className="w-full h-2 bg-slate-200 rounded-lg appearance-none cursor-pointer accent-teal-600"
                                />
                                <span className="font-bold text-slate-900 w-8">{scoreThreshold}</span>
                            </div>
                            <p className="text-xs text-slate-400 mt-1">Flag answers where AI score differs by more than {scoreThreshold} points.</p>
                        </div>
                        
                        {isProcessing && (
                            <div className="flex items-center text-sm text-teal-600 font-medium animate-pulse">
                                <div className="h-2 w-2 rounded-full bg-teal-600 mr-2"></div>
                                Processing & AI Scoring...
                            </div>
                        )}
                    </div>
                </div>
            </Card>

            {discrepancies.length > 0 ? (
                <div className="grid grid-cols-1 gap-6">
                   <h3 className="text-lg font-bold text-slate-900 flex items-center">
                      <AlertCircle className="h-5 w-5 text-amber-500 mr-2" />
                      Flagged Discrepancies ({discrepancies.length} Students)
                   </h3>
                   {discrepancies.map((group) => (
                      <DiscrepancyGroupCard key={group.userId} group={group} />
                   ))}
                </div>
            ) : (
                auditData.length > 0 && !isProcessing && (
                    <div className="p-12 text-center bg-green-50 rounded-xl border border-green-100">
                        <CheckCircle className="h-12 w-12 text-green-500 mx-auto mb-3" />
                        <h3 className="text-lg font-bold text-green-900">High Agreement Found</h3>
                        <p className="text-green-700">No discrepancies found with the current threshold of {scoreThreshold} points.</p>
                    </div>
                )
            )}
        </div>
    );
};

const DiscrepancyGroupCard: React.FC<{ group: { userId: number, answers: AnswerRecord[] } }> = ({ group }) => {
    const [isExpanded, setIsExpanded] = useState(false);
    
    return (
        <Card className="overflow-hidden border-l-4 border-l-amber-400">
             <div 
                className="px-6 py-4 flex items-center justify-between cursor-pointer hover:bg-slate-50 transition-colors"
                onClick={() => setIsExpanded(!isExpanded)}
             >
                <div className="flex items-center gap-4">
                    <div className="h-10 w-10 rounded-full bg-slate-200 flex items-center justify-center font-bold text-slate-600">
                        #{group.userId}
                    </div>
                    <div>
                        <h4 className="font-bold text-slate-900">Student ID: {group.userId}</h4>
                        <p className="text-sm text-slate-500">{group.answers.length} flagged answers</p>
                    </div>
                </div>
                {isExpanded ? <ChevronDown className="h-5 w-5 text-slate-400" /> : <ChevronRight className="h-5 w-5 text-slate-400" />}
             </div>
             
             {isExpanded && (
                 <div className="bg-slate-50/50 border-t border-slate-100 p-6 space-y-6">
                    {group.answers.map((ans, idx) => (
                        <div key={idx} className="bg-white rounded-lg border border-slate-200 shadow-sm p-5">
                            <div className="flex justify-between items-start mb-3">
                                <div className="flex items-center gap-2">
                                    <span className="px-2 py-1 rounded bg-slate-100 text-slate-600 text-xs font-bold uppercase tracking-wider">
                                        {ans.name_rus}
                                    </span>
                                    <span className="text-slate-300">|</span>
                                    <span className="text-sm font-medium text-slate-700">{ans.examiner_name}</span>
                                </div>
                                <div className="flex items-center gap-4">
                                     <div className="text-right">
                                        <div className="text-xs text-slate-400 uppercase font-bold">Teacher</div>
                                        <div className="text-xl font-bold text-slate-900">{ans.score}</div>
                                     </div>
                                     <div className="text-right">
                                        <div className="text-xs text-slate-400 uppercase font-bold">AI Model</div>
                                        <div className="text-xl font-bold text-teal-600">{Math.round(ans.ai_score || 0)}</div>
                                     </div>
                                     <div className="ml-2 px-2 py-1 bg-amber-100 text-amber-800 text-xs font-bold rounded">
                                        Diff: {Math.abs(ans.score - (ans.ai_score || 0)).toFixed(1)}
                                     </div>
                                </div>
                            </div>
                            
                            <h5 className="font-bold text-slate-900 mb-2">{ans.question}</h5>
                            <div className="bg-slate-50 p-3 rounded border border-slate-100 text-sm text-slate-700 leading-relaxed mb-3 whitespace-pre-wrap">
                                {ans.answer}
                            </div>
                            
                            <div className="flex items-start gap-2">
                                <FileText className="h-4 w-4 text-slate-400 mt-0.5" />
                                <div>
                                    <span className="text-xs font-bold text-slate-500 uppercase">Examiner Comment:</span>
                                    <p className="text-sm text-slate-600 italic">"{ans.comment}"</p>
                                </div>
                            </div>
                        </div>
                    ))}
                 </div>
             )}
        </Card>
    );
};

// --------------------------------------------------------------------------------
// Shared Component
// --------------------------------------------------------------------------------
const AnswerRow: React.FC<{ answer: AnswerRecord }> = ({ answer }) => {
   const [expanded, setExpanded] = useState(false);
   
   // Simple logic if ai_score exists
   const aiScore = answer.ai_score || answer.score; 
   const diff = Math.abs(answer.score - aiScore);
   const isFlagged = diff > 5; // Default threshold for visual indicator in list

   return (
      <div className="border border-slate-200 rounded-lg hover:shadow-sm transition-shadow bg-white">
         <div 
            className="p-4 cursor-pointer flex justify-between items-center"
            onClick={() => setExpanded(!expanded)}
         >
            <div className="flex-1 min-w-0 pr-4">
               <div className="flex items-center space-x-2">
                  <span className="text-xs font-semibold text-teal-600 bg-teal-50 px-2 py-0.5 rounded">{answer.name_rus}</span>
                  {isFlagged && <AlertCircle className="h-4 w-4 text-amber-500" />}
               </div>
               <p className="text-sm font-medium text-slate-900 truncate mt-1">{answer.question}</p>
            </div>
            
            <div className="flex items-center space-x-6 text-sm">
               <div className="text-right">
                  <div className="font-bold text-slate-800">{answer.score}</div>
                  <div className="text-xs text-slate-400">Teacher</div>
               </div>
               <div className="text-right border-l pl-6 border-slate-100">
                  <div className="font-bold text-slate-500">{Math.round(aiScore)}</div>
                  <div className="text-xs text-slate-400">AI Model</div>
               </div>
               {expanded ? <ChevronDown className="h-4 w-4 text-slate-400" /> : <ChevronRight className="h-4 w-4 text-slate-400" />}
            </div>
         </div>

         {expanded && (
            <div className="px-4 pb-4 border-t border-slate-100 bg-slate-50/50">
               <div className="mt-3 space-y-3">
                  <div>
                     <h5 className="text-xs font-semibold text-slate-500 uppercase">Student Answer</h5>
                     <p className="text-sm text-slate-700 mt-1 whitespace-pre-wrap">{answer.answer}</p>
                  </div>
                  <div>
                     <h5 className="text-xs font-semibold text-slate-500 uppercase">Examiner Comment</h5>
                     <p className="text-sm text-slate-600 mt-1 italic">"{answer.comment}"</p>
                  </div>
               </div>
            </div>
         )}
      </div>
   );
};

export default ScoreAudit;
