
import React, { useState } from 'react';
import { Card } from '../components/ui/Card';
import { MOCK_PENDING_ANSWERS } from '../constants';
import { PendingAnswer } from '../types';
import { CheckCircle, Clock, User, MessageSquare, Send, ClipboardList } from 'lucide-react';

const GradingQueue: React.FC = () => {
    const [queue, setQueue] = useState<PendingAnswer[]>(MOCK_PENDING_ANSWERS);
    const [selectedAnswerId, setSelectedAnswerId] = useState<string | null>(null);
    
    // Grading Form State
    const [score, setScore] = useState<string>('');
    const [comment, setComment] = useState<string>('');

    const selectedAnswer = queue.find(a => a.id === selectedAnswerId);

    const handleSubmit = () => {
        if (!selectedAnswer) return;
        
        // Remove from queue mock
        setQueue(prev => prev.filter(a => a.id !== selectedAnswerId));
        setSelectedAnswerId(null);
        setScore('');
        setComment('');
        
        // In real app, API call here
    };

    return (
        <div className="h-[calc(100vh-8rem)] flex flex-col md:flex-row gap-6">
            {/* Left: Queue List */}
            <div className="w-full md:w-1/3 flex flex-col gap-4">
                <div className="flex justify-between items-center">
                    <h2 className="text-xl font-bold text-slate-900">Grading Queue</h2>
                    <span className="bg-teal-100 text-teal-800 text-xs font-bold px-2 py-1 rounded-full">
                        {queue.length} Pending
                    </span>
                </div>
                
                <div className="flex-1 overflow-y-auto custom-scrollbar space-y-3 pr-2">
                    {queue.length === 0 ? (
                        <div className="text-center py-10 bg-slate-50 rounded-lg border border-dashed border-slate-200">
                            <CheckCircle className="h-10 w-10 text-green-400 mx-auto mb-3" />
                            <p className="text-slate-500 font-medium">All caught up!</p>
                            <p className="text-xs text-slate-400">No pending answers to grade.</p>
                        </div>
                    ) : (
                        queue.map(item => (
                            <div 
                                key={item.id}
                                onClick={() => setSelectedAnswerId(item.id)}
                                className={`p-4 rounded-xl border cursor-pointer transition-all hover:shadow-md
                                    ${selectedAnswerId === item.id 
                                        ? 'bg-teal-50 border-teal-500 ring-1 ring-teal-500' 
                                        : 'bg-white border-slate-200 hover:border-teal-300'}
                                `}
                            >
                                <div className="flex justify-between items-start mb-2">
                                    <span className="text-xs font-bold text-slate-500 uppercase tracking-wider">{item.category}</span>
                                    <span className="text-xs text-slate-400 flex items-center">
                                        <Clock className="h-3 w-3 mr-1" /> 2h ago
                                    </span>
                                </div>
                                <h4 className="font-semibold text-slate-800 line-clamp-2 mb-2 text-sm">{item.question}</h4>
                                <div className="flex items-center gap-2">
                                    <div className="h-6 w-6 rounded-full bg-slate-100 flex items-center justify-center">
                                        <User className="h-3 w-3 text-slate-400" />
                                    </div>
                                    <span className="text-xs text-slate-500 font-mono">Student #{item.student_id}</span>
                                </div>
                            </div>
                        ))
                    )}
                </div>
            </div>

            {/* Right: Grading Area */}
            <div className="w-full md:w-2/3 flex flex-col">
                {selectedAnswer ? (
                    <Card className="flex-1 flex flex-col overflow-hidden">
                        <div className="border-b border-slate-100 p-6 bg-slate-50">
                            <div className="flex items-center gap-2 mb-2">
                                <span className="bg-white border border-slate-200 px-2 py-1 rounded text-xs font-bold text-slate-500 uppercase">
                                    {selectedAnswer.category}
                                </span>
                            </div>
                            <h3 className="text-lg font-bold text-slate-900">{selectedAnswer.question}</h3>
                        </div>
                        
                        <div className="flex-1 overflow-y-auto p-6">
                            <div className="mb-2 text-xs font-bold text-slate-400 uppercase tracking-wide">Student Answer</div>
                            <div className="bg-blue-50/50 p-6 rounded-xl border border-blue-100 text-slate-800 leading-relaxed text-lg font-serif">
                                {selectedAnswer.answer}
                            </div>
                            
                            <div className="mt-8 border-t border-slate-100 pt-6">
                                <h4 className="font-bold text-slate-900 mb-4 flex items-center">
                                    <MessageSquare className="h-4 w-4 mr-2 text-teal-600" />
                                    Assessment
                                </h4>
                                
                                <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
                                    <div className="md:col-span-1">
                                        <label className="block text-sm font-medium text-slate-700 mb-1">Score (0-100)</label>
                                        <input 
                                            type="number" 
                                            min="0" max="100"
                                            value={score}
                                            onChange={(e) => setScore(e.target.value)}
                                            className="w-full p-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-teal-500 focus:border-teal-500 text-center font-bold text-lg"
                                            placeholder="--"
                                        />
                                    </div>
                                    <div className="md:col-span-3">
                                        <label className="block text-sm font-medium text-slate-700 mb-1">Feedback Comment</label>
                                        <textarea 
                                            rows={3}
                                            value={comment}
                                            onChange={(e) => setComment(e.target.value)}
                                            className="w-full p-3 border border-slate-300 rounded-lg focus:ring-2 focus:ring-teal-500 focus:border-teal-500 text-sm"
                                            placeholder="Provide constructive feedback..."
                                        />
                                    </div>
                                </div>
                            </div>
                        </div>

                        <div className="p-4 bg-slate-50 border-t border-slate-100 flex justify-end">
                            <button 
                                onClick={handleSubmit}
                                disabled={!score || !comment}
                                className={`flex items-center px-6 py-2 rounded-lg font-bold text-white transition-all
                                    ${!score || !comment 
                                        ? 'bg-slate-300 cursor-not-allowed' 
                                        : 'bg-teal-600 hover:bg-teal-700 shadow-md'}
                                `}
                            >
                                <Send className="h-4 w-4 mr-2" />
                                Submit Grade
                            </button>
                        </div>
                    </Card>
                ) : (
                    <div className="flex-1 flex flex-col items-center justify-center bg-slate-50 rounded-xl border border-dashed border-slate-200">
                        <ClipboardList className="h-16 w-16 text-slate-300 mb-4" />
                        <h3 className="text-lg font-medium text-slate-500">Select an answer to grade</h3>
                        <p className="text-sm text-slate-400">Choose from the queue on the left</p>
                    </div>
                )}
            </div>
        </div>
    );
};

export default GradingQueue;
