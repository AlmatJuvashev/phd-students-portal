
import React from 'react';
import { Card } from '../components/ui/Card';
import { MOCK_ANSWERS, PROTOCOLS, MOCK_PENDING_ANSWERS } from '../constants';
import { BarChart, Bar, XAxis, YAxis, Tooltip, ResponsiveContainer, Cell } from 'recharts';
import { Book, Users, FileText, AlertTriangle, ClipboardList, CheckCircle } from 'lucide-react';
import { Link } from 'react-router-dom';
import { UserRole } from '../types';

interface DashboardProps {
  role?: UserRole;
}

const Dashboard: React.FC<DashboardProps> = ({ role = 'admin' }) => {
  // Mock aggregations
  const totalAnswers = MOCK_ANSWERS.length;
  const uniqueStudents = new Set(MOCK_ANSWERS.map(a => a.user_id)).size;
  const uniqueTeachers = new Set(MOCK_ANSWERS.map(a => a.examiner_id)).size;
  const totalProtocols = PROTOCOLS.length;
  
  // Histogram data
  const scoreDistribution = [
    { range: '0-20', count: 0 },
    { range: '21-40', count: 0 },
    { range: '41-60', count: 0 },
    { range: '61-80', count: 0 },
    { range: '81-100', count: 0 },
  ];

  MOCK_ANSWERS.forEach(a => {
    if (a.score <= 20) scoreDistribution[0].count++;
    else if (a.score <= 40) scoreDistribution[1].count++;
    else if (a.score <= 60) scoreDistribution[2].count++;
    else if (a.score <= 80) scoreDistribution[3].count++;
    else scoreDistribution[4].count++;
  });

  const StatCard = ({ icon: Icon, label, value, color }: any) => (
    <Card className="flex flex-col justify-center">
      <div className="flex items-center">
        <div className={`p-3 rounded-lg ${color}`}>
          <Icon className="h-6 w-6 text-white" />
        </div>
        <div className="ml-4">
          <p className="text-sm font-medium text-slate-500">{label}</p>
          <p className="text-2xl font-bold text-slate-900">{value}</p>
        </div>
      </div>
    </Card>
  );

  return (
    <div className="space-y-6">
      {role === 'teacher' ? (
        // TEACHER DASHBOARD
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
           <StatCard icon={ClipboardList} label="Pending Grading" value={MOCK_PENDING_ANSWERS.length} color="bg-teal-500" />
           <StatCard icon={CheckCircle} label="Graded this Week" value={12} color="bg-blue-500" />
           <StatCard icon={AlertTriangle} label="AI Discrepancies" value={2} color="bg-amber-500" />
        </div>
      ) : (
        // ADMIN DASHBOARD
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          <StatCard icon={FileText} label="Total Answers" value={totalAnswers} color="bg-blue-500" />
          <StatCard icon={Users} label="Active Students" value={uniqueStudents} color="bg-teal-500" />
          <StatCard icon={Book} label="RAG Protocols" value={totalProtocols} color="bg-indigo-500" />
          <StatCard icon={Users} label="Examiners" value={uniqueTeachers} color="bg-rose-500" />
        </div>
      )}

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Score Distribution Chart */}
        <div className="lg:col-span-2">
          <Card title={role === 'teacher' ? "Your Grading Distribution" : "Global Score Distribution"} subtitle="Overview of student performance">
            <div className="h-64 w-full">
              <ResponsiveContainer width="100%" height="100%">
                <BarChart data={scoreDistribution}>
                  <XAxis dataKey="range" stroke="#94a3b8" fontSize={12} tickLine={false} axisLine={false} />
                  <YAxis stroke="#94a3b8" fontSize={12} tickLine={false} axisLine={false} />
                  <Tooltip 
                    cursor={{ fill: '#f1f5f9' }}
                    contentStyle={{ borderRadius: '8px', border: 'none', boxShadow: '0 4px 6px -1px rgb(0 0 0 / 0.1)' }} 
                  />
                  <Bar dataKey="count" radius={[4, 4, 0, 0]}>
                    {scoreDistribution.map((entry, index) => (
                      <Cell key={`cell-${index}`} fill="#0d9488" fillOpacity={0.7 + (index * 0.05)} />
                    ))}
                  </Bar>
                </BarChart>
              </ResponsiveContainer>
            </div>
          </Card>
        </div>

        {/* Sidebar / Quick Actions */}
        <div className="space-y-6">
          {role !== 'teacher' && (
             <Card title="System Status">
               <div className="space-y-4">
                  <div className="flex justify-between items-center py-2 border-b border-slate-50">
                     <span className="text-sm text-slate-600">Default Model</span>
                     <span className="text-sm font-semibold text-slate-900">Qwen 2 72B</span>
                  </div>
                  <div className="flex justify-between items-center py-2 border-b border-slate-50">
                     <span className="text-sm text-slate-600">Embeddings</span>
                     <span className="text-sm font-semibold text-slate-900">text-embedding-3-large</span>
                  </div>
                  <div className="flex justify-between items-center py-2">
                     <span className="text-sm text-slate-600">RAG Status</span>
                     <span className="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-green-100 text-green-700">Active</span>
                  </div>
               </div>
            </Card>
          )}

          {role === 'teacher' ? (
              <Card title="Quick Actions">
                <div className="space-y-3">
                   <Link to="/grading" className="block w-full py-2 px-4 bg-teal-600 text-white text-center rounded-lg font-medium hover:bg-teal-700">
                      Go to Grading Queue
                   </Link>
                   <Link to="/audit" className="block w-full py-2 px-4 border border-slate-300 text-slate-700 text-center rounded-lg font-medium hover:bg-slate-50">
                      Review My Stats
                   </Link>
                </div>
              </Card>
          ) : (
            <Card className="bg-amber-50 border-amber-200">
              <div className="flex items-start">
                <AlertTriangle className="h-5 w-5 text-amber-600 mt-0.5" />
                <div className="ml-3">
                  <h3 className="text-sm font-medium text-amber-800">Audit Alert</h3>
                  <p className="mt-1 text-sm text-amber-700">
                    3 teachers show significant deviation in scoring compared to the cohort average.
                  </p>
                  <div className="mt-3">
                    <Link to="/audit" className="text-sm font-medium text-amber-800 hover:text-amber-900 underline">
                      View Audit Report
                    </Link>
                  </div>
                </div>
              </div>
            </Card>
          )}
        </div>
      </div>
    </div>
  );
};

export default Dashboard;
