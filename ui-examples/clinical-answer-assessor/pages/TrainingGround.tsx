
import React, { useState, useMemo } from 'react';
import { Card } from '../components/ui/Card';
import { TOPICS, MOCK_ANSWERS } from '../constants';
import { Topic, AnswerRecord } from '../types';
import { BarChart, ChevronRight, Search, Filter, History, LayoutGrid, Calendar, Trophy, ArrowRight, X, CheckCircle, AlertTriangle, FileText, Sparkles } from 'lucide-react';
import { useNavigate } from 'react-router-dom';

const TrainingGround: React.FC = () => {
  const navigate = useNavigate();
  const [activeTab, setActiveTab] = useState<'browse' | 'history'>('browse');
  const [searchQuery, setSearchQuery] = useState('');
  const [difficultyFilter, setDifficultyFilter] = useState<'All' | 'Beginner' | 'Intermediate' | 'Advanced'>('All');
  
  // State for History Details Modal
  const [selectedHistoryRecord, setSelectedHistoryRecord] = useState<AnswerRecord | null>(null);

  // Filter Topics
  const filteredTopics = useMemo(() => {
    const query = searchQuery.toLowerCase();
    return TOPICS.filter(topic => {
      const matchesSearch = topic.title.toLowerCase().includes(query) ||
                            topic.description.toLowerCase().includes(query) ||
                            topic.questions.some(q => q.text.toLowerCase().includes(query));
      
      const matchesDifficulty = difficultyFilter === 'All' || 
                                topic.questions.some(q => q.difficulty === difficultyFilter);

      return matchesSearch && matchesDifficulty;
    });
  }, [searchQuery, difficultyFilter]);

  // Mock History (User 101)
  const history = useMemo(() => {
     return MOCK_ANSWERS
        .filter(a => a.user_id === 101)
        .sort((a,b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime());
  }, []);

  return (
    <div className="space-y-6 relative">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
           <h2 className="text-2xl font-bold text-slate-900">Training Ground</h2>
           <p className="text-slate-500 mt-1">Practice with AI-powered assessment and clinical protocol references.</p>
        </div>
        <div className="flex bg-white p-1 rounded-lg border border-slate-200 self-start shadow-sm">
           <button 
             onClick={() => setActiveTab('browse')}
             className={`px-4 py-2 rounded-md text-sm font-medium transition-all flex items-center ${activeTab === 'browse' ? 'bg-teal-50 text-teal-700 shadow-sm ring-1 ring-teal-200' : 'text-slate-500 hover:text-slate-700 hover:bg-slate-50'}`}
           >
             <LayoutGrid className="h-4 w-4 mr-2" /> Browse
           </button>
           <button 
             onClick={() => setActiveTab('history')}
             className={`px-4 py-2 rounded-md text-sm font-medium transition-all flex items-center ${activeTab === 'history' ? 'bg-teal-50 text-teal-700 shadow-sm ring-1 ring-teal-200' : 'text-slate-500 hover:text-slate-700 hover:bg-slate-50'}`}
           >
             <History className="h-4 w-4 mr-2" /> My Progress
           </button>
        </div>
      </div>

      {activeTab === 'browse' ? (
        <div className="space-y-6">
          {/* Filters */}
          <div className="bg-white p-4 rounded-xl border border-slate-200 shadow-sm flex flex-col md:flex-row gap-4 items-center">
             <div className="relative flex-1 w-full">
                <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-slate-400" />
                <input 
                  type="text"
                  placeholder="Search topics or questions by keywords..."
                  className="w-full pl-10 pr-4 py-2 rounded-lg border border-slate-200 focus:outline-none focus:ring-2 focus:ring-teal-500 focus:border-transparent text-sm"
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                />
             </div>
             <div className="flex items-center gap-2 w-full md:w-auto overflow-x-auto pb-1 md:pb-0">
                <Filter className="h-4 w-4 text-slate-400 shrink-0 ml-1" />
                <span className="text-sm font-medium text-slate-700 mr-2 whitespace-nowrap">Difficulty:</span>
                <div className="flex space-x-1">
                    {['All', 'Beginner', 'Intermediate', 'Advanced'].map(level => (
                    <button
                        key={level}
                        onClick={() => setDifficultyFilter(level as any)}
                        className={`px-3 py-1.5 rounded-full text-xs font-medium border transition-colors whitespace-nowrap
                        ${difficultyFilter === level 
                            ? 'bg-teal-600 border-teal-600 text-white' 
                            : 'bg-white border-slate-200 text-slate-600 hover:border-slate-300 hover:bg-slate-50'}
                        `}
                    >
                        {level}
                    </button>
                    ))}
                </div>
             </div>
          </div>

          {/* Grid */}
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {filteredTopics.length > 0 ? (
                filteredTopics.map((topic) => (
                <TopicCard key={topic.id} topic={topic} onStart={() => navigate(`/training/${topic.id}`)} />
                ))
            ) : (
                <div className="col-span-full py-16 text-center bg-slate-50 rounded-xl border border-dashed border-slate-200">
                    <Search className="h-10 w-10 text-slate-300 mx-auto mb-3" />
                    <p className="text-slate-500 font-medium">No topics found matching your criteria.</p>
                    <button 
                        onClick={() => {setSearchQuery(''); setDifficultyFilter('All');}}
                        className="mt-4 text-teal-600 hover:text-teal-700 text-sm font-medium underline"
                    >
                        Clear filters
                    </button>
                </div>
            )}
          </div>
        </div>
      ) : (
        <div>
            <HistoryView 
              history={history} 
              onViewDetails={(record) => setSelectedHistoryRecord(record)} 
            />
        </div>
      )}

      {/* History Details Modal */}
      {selectedHistoryRecord && (
        <HistoryDetailsModal 
          record={selectedHistoryRecord} 
          onClose={() => setSelectedHistoryRecord(null)} 
        />
      )}
    </div>
  );
};

const TopicCard: React.FC<{ topic: Topic; onStart: () => void }> = ({ topic, onStart }) => {
  const difficulties = Array.from(new Set(topic.questions.map(q => q.difficulty)));
  
  const getDifficultyColor = (d: string) => {
    switch(d) {
        case 'Beginner': return 'bg-green-100 text-green-700 border-green-200';
        case 'Intermediate': return 'bg-orange-100 text-orange-700 border-orange-200';
        case 'Advanced': return 'bg-red-100 text-red-700 border-red-200';
        default: return 'bg-slate-100 text-slate-700';
    }
  };

  return (
    <Card className="hover:border-teal-300 transition-all duration-200 cursor-pointer group h-full flex flex-col shadow-sm hover:shadow-md">
        <div onClick={onStart} className="flex-1 flex flex-col h-full">
            <div className="flex justify-between items-start mb-2">
                <h3 className="text-lg font-bold text-slate-900 group-hover:text-teal-700 transition-colors line-clamp-1">{topic.title}</h3>
            </div>
        
            <p className="text-sm text-slate-600 mb-4 flex-1 line-clamp-3">{topic.description}</p>
        
            <div className="space-y-3 mb-6">
                <div className="flex flex-wrap gap-1.5">
                    {difficulties.map(d => (
                        <span key={d as string} className={`text-[10px] uppercase tracking-wider font-bold px-2 py-0.5 rounded border ${getDifficultyColor(d as string)}`}>
                            {d as string}
                        </span>
                    ))}
                </div>

                <div className="h-px bg-slate-100 w-full my-3"></div>

                <h4 className="text-xs font-semibold text-slate-400 uppercase tracking-wider mb-1">Key Objectives</h4>
                <ul className="text-sm text-slate-500 space-y-1">
                    {topic.objectives.slice(0, 2).map((obj, i) => (
                        <li key={i} className="flex items-start line-clamp-1">
                        <span className="mr-2 text-teal-500 shrink-0">•</span> {obj}
                        </li>
                    ))}
                </ul>
            </div>

            <div className="pt-4 border-t border-slate-100 flex justify-between items-center text-xs font-medium text-slate-500 mt-auto">
                <span className="flex items-center bg-slate-100 px-2 py-1 rounded text-slate-600">
                    <BarChart className="h-3 w-3 mr-1.5" /> {topic.questions.length} Questions
                </span>
                <span className="flex items-center text-teal-600 group-hover:translate-x-1 transition-transform font-semibold">
                    Start Practice <ChevronRight className="h-4 w-4 ml-1" />
                </span>
            </div>
        </div>
    </Card>
  );
}

const HistoryView: React.FC<{ history: AnswerRecord[], onViewDetails: (record: AnswerRecord) => void }> = ({ history, onViewDetails }) => {
    return (
        <Card className="overflow-hidden">
            <div className="p-6 border-b border-slate-100 flex justify-between items-center bg-slate-50">
                <div>
                    <h3 className="text-lg font-bold text-slate-900">Training History</h3>
                    <p className="text-sm text-slate-500">Your recent practice sessions and scores.</p>
                </div>
                <div className="flex items-center space-x-2">
                   <div className="bg-white px-3 py-1 rounded-md border border-slate-200 shadow-sm">
                      <span className="text-xs text-slate-400 uppercase font-semibold mr-2">Avg Score</span>
                      <span className="text-lg font-bold text-teal-600">
                         {history.length > 0 ? Math.round(history.reduce((a,b) => a + b.score, 0) / history.length) : 0}
                      </span>
                   </div>
                </div>
            </div>
            
            {history.length === 0 ? (
                <div className="text-center py-12">
                   <Trophy className="h-12 w-12 text-slate-200 mx-auto mb-3" />
                   <p className="text-slate-500">No practice history yet. Start a session!</p>
                </div>
            ) : (
                <div className="overflow-x-auto">
                    <table className="min-w-full divide-y divide-slate-200">
                        <thead className="bg-white">
                            <tr>
                                <th className="px-6 py-3 text-left text-xs font-medium text-slate-500 uppercase tracking-wider">Topic</th>
                                <th className="px-6 py-3 text-left text-xs font-medium text-slate-500 uppercase tracking-wider">Question</th>
                                <th className="px-6 py-3 text-left text-xs font-medium text-slate-500 uppercase tracking-wider">Score</th>
                                <th className="px-6 py-3 text-left text-xs font-medium text-slate-500 uppercase tracking-wider">Date</th>
                                <th className="px-6 py-3 text-right text-xs font-medium text-slate-500 uppercase tracking-wider">Action</th>
                            </tr>
                        </thead>
                        <tbody className="bg-white divide-y divide-slate-200">
                            {history.map((record, idx) => (
                                <tr key={idx} className="hover:bg-slate-50 transition-colors">
                                    <td className="px-6 py-4 whitespace-nowrap">
                                        <div className="text-sm font-medium text-slate-900">{record.topic}</div>
                                        <div className="text-xs text-slate-500">{record.name_rus}</div>
                                    </td>
                                    <td className="px-6 py-4">
                                        <div className="text-sm text-slate-500 line-clamp-2 min-w-[200px]">{record.question}</div>
                                    </td>
                                    <td className="px-6 py-4 whitespace-nowrap">
                                        <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium
                                            ${record.score >= 80 ? 'bg-green-100 text-green-800' :
                                              record.score >= 50 ? 'bg-yellow-100 text-yellow-800' :
                                              'bg-red-100 text-red-800'}
                                        `}>
                                            {record.score}
                                        </span>
                                    </td>
                                    <td className="px-6 py-4 whitespace-nowrap text-sm text-slate-500">
                                        <div className="flex items-center">
                                            <Calendar className="h-4 w-4 mr-2 text-slate-400" />
                                            {record.created_at}
                                        </div>
                                    </td>
                                    <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                                        <button 
                                          onClick={() => onViewDetails(record)}
                                          className="text-teal-600 hover:text-teal-900 flex items-center justify-end ml-auto group"
                                        >
                                            Details <ArrowRight className="h-3 w-3 ml-1 group-hover:translate-x-0.5 transition-transform" />
                                        </button>
                                    </td>
                                </tr>
                            ))}
                        </tbody>
                    </table>
                </div>
            )}
        </Card>
    );
};

// ----------------------------------------------------------------------
// History Details Modal
// ----------------------------------------------------------------------

interface HistoryDetailsModalProps {
  record: AnswerRecord;
  onClose: () => void;
}

const HistoryDetailsModal: React.FC<HistoryDetailsModalProps> = ({ record, onClose }) => {
  // Find the original rubric items if available to map IDs to Names (optional, if needed more than what's in record)
  // For now, we assume record.section_scores has what we need or we look up topic.
  const rubricItems = useMemo(() => {
     const t = TOPICS.find(t => t.title === record.topic); // simple match by title
     const q = t?.questions.find(q => q.text === record.question);
     return q?.rubric || [];
  }, [record]);
  
  const getRubricName = (id: string) => {
    return rubricItems.find(r => r.id === id)?.criteria || "Criterion";
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4 sm:p-6" role="dialog">
       {/* Backdrop */}
       <div 
         className="fixed inset-0 bg-slate-900/60 backdrop-blur-sm transition-opacity" 
         onClick={onClose}
       ></div>

       {/* Modal Content */}
       <div className="bg-white rounded-2xl shadow-2xl w-full max-w-4xl max-h-[90vh] flex flex-col z-10 overflow-hidden animate-in zoom-in-95 duration-200">
          
          {/* Header */}
          <div className="px-6 py-4 border-b border-slate-100 flex items-start justify-between bg-slate-50/50 shrink-0">
             <div>
                <div className="flex items-center gap-2 mb-1">
                   <span className="text-xs font-bold text-teal-600 uppercase tracking-wider">{record.topic}</span>
                   <span className="text-slate-300">|</span>
                   <span className="text-xs text-slate-500">{new Date(record.created_at).toLocaleDateString()}</span>
                </div>
                <h3 className="text-lg font-bold text-slate-900 line-clamp-1">{record.question}</h3>
             </div>
             
             <div className="flex items-center gap-4">
                <div className={`flex flex-col items-end px-3 py-1 rounded-lg border 
                   ${record.score >= 80 ? 'bg-green-50 border-green-100 text-green-700' :
                     record.score >= 50 ? 'bg-yellow-50 border-yellow-100 text-yellow-700' :
                     'bg-red-50 border-red-100 text-red-700'}
                `}>
                   <span className="text-2xl font-bold leading-none">{record.score}</span>
                   <span className="text-[10px] uppercase font-bold opacity-80">Total Score</span>
                </div>
                <button 
                  onClick={onClose}
                  className="p-2 rounded-full hover:bg-slate-200 text-slate-400 hover:text-slate-600 transition-colors"
                >
                   <X className="h-6 w-6" />
                </button>
             </div>
          </div>

          {/* Body */}
          <div className="flex-1 overflow-y-auto p-6 custom-scrollbar">
             <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-8">
                {/* Feedback: Strengths */}
                <div className="bg-green-50/50 rounded-xl border border-green-100 p-4">
                   <h4 className="flex items-center text-sm font-bold text-green-800 mb-3 uppercase tracking-wide">
                      <CheckCircle className="h-4 w-4 mr-2" /> Strengths
                   </h4>
                   <ul className="space-y-2">
                      {record.strengths && record.strengths.length > 0 ? (
                         record.strengths.map((s, i) => (
                           <li key={i} className="flex items-start text-sm text-slate-700">
                              <span className="mr-2 text-green-500 mt-1">•</span>
                              {s}
                           </li>
                         ))
                      ) : (
                         <li className="text-sm text-slate-500 italic">No specific strengths listed.</li>
                      )}
                   </ul>
                </div>

                {/* Feedback: Weaknesses */}
                <div className="bg-amber-50/50 rounded-xl border border-amber-100 p-4">
                   <h4 className="flex items-center text-sm font-bold text-amber-800 mb-3 uppercase tracking-wide">
                      <AlertTriangle className="h-4 w-4 mr-2" /> Improvements
                   </h4>
                   <ul className="space-y-2">
                      {record.weaknesses && record.weaknesses.length > 0 ? (
                         record.weaknesses.map((w, i) => (
                           <li key={i} className="flex items-start text-sm text-slate-700">
                              <span className="mr-2 text-amber-500 mt-1">•</span>
                              {w}
                           </li>
                         ))
                      ) : (
                         <li className="text-sm text-slate-500 italic">No specific improvements needed.</li>
                      )}
                   </ul>
                </div>
             </div>

             {/* Rubric Breakdown (if available) */}
             {record.section_scores && record.section_scores.length > 0 && (
                <div className="mb-8">
                   <h4 className="flex items-center text-sm font-bold text-slate-900 mb-4 uppercase tracking-wide">
                      <Sparkles className="h-4 w-4 mr-2 text-teal-500" /> Scoring Breakdown
                   </h4>
                   <div className="space-y-3">
                      {record.section_scores.map((section) => (
                         <div key={section.id} className="bg-white border border-slate-200 rounded-lg p-4 shadow-sm hover:border-teal-200 transition-colors">
                            <div className="flex justify-between items-start mb-2">
                               <h5 className="font-semibold text-slate-800 text-sm">
                                  {getRubricName(section.id)}
                               </h5>
                               <span className={`text-xs font-bold px-2 py-1 rounded border
                                  ${section.score === section.max ? 'bg-green-50 text-green-700 border-green-200' : 'bg-slate-50 text-slate-600 border-slate-200'}
                               `}>
                                  {section.score} / {section.max}
                               </span>
                            </div>
                            <p className="text-sm text-slate-600 bg-slate-50 p-2 rounded">
                               <span className="font-semibold text-slate-500 text-xs uppercase mr-1">Feedback:</span>
                               {section.feedback}
                            </p>
                         </div>
                      ))}
                   </div>
                </div>
             )}

             {/* Student Answer */}
             <div>
                <h4 className="flex items-center text-sm font-bold text-slate-900 mb-3 uppercase tracking-wide">
                   <FileText className="h-4 w-4 mr-2 text-slate-400" /> Your Response
                </h4>
                <div className="bg-slate-50 rounded-xl border border-slate-200 p-5">
                   <p className="whitespace-pre-wrap text-slate-700 text-sm leading-relaxed font-serif">
                      {record.answer}
                   </p>
                </div>
             </div>
          </div>

          {/* Footer */}
          <div className="p-4 border-t border-slate-100 bg-slate-50 flex justify-end">
             <button 
                onClick={onClose}
                className="px-6 py-2 bg-white border border-slate-300 text-slate-700 font-medium rounded-lg hover:bg-slate-50 shadow-sm transition-colors"
             >
                Close Details
             </button>
          </div>

       </div>
    </div>
  );
};

export default TrainingGround;
