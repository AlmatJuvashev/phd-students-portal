
import React from 'react';
import { Link } from 'react-router-dom';
import { 
  Activity, 
  ArrowRight, 
} from 'lucide-react';

const Landing: React.FC = () => {
  return (
    <div className="min-h-screen bg-white font-sans text-slate-900 selection:bg-teal-100 selection:text-teal-900 flex flex-col">
      
      {/* Navbar */}
      <nav className="w-full bg-white/80 backdrop-blur-md border-b border-slate-100 z-50 fixed top-0">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center h-16">
            <div className="flex items-center gap-2">
              <div className="bg-teal-600 p-1.5 rounded-lg">
                <Activity className="h-5 w-5 text-white" />
              </div>
              <span className="font-bold text-xl tracking-tight text-slate-900">ClinAssessor</span>
            </div>
            <div className="flex items-center gap-4">
              <Link 
                to="/dashboard" 
                className="px-5 py-2 rounded-full bg-slate-900 text-white text-sm font-semibold hover:bg-slate-800 transition-all flex items-center group"
              >
                Launch App
                <ArrowRight className="h-4 w-4 ml-2 group-hover:translate-x-0.5 transition-transform" />
              </Link>
            </div>
          </div>
        </div>
      </nav>

      {/* Hero Section */}
      <div className="relative flex-1 flex flex-col justify-center items-center pt-24 overflow-hidden">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 relative z-10 text-center">
          <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-teal-50 border border-teal-100 text-teal-700 text-sm font-medium mb-8 animate-in fade-in slide-in-from-bottom-4 duration-700">
            <span className="relative flex h-2 w-2">
              <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-teal-400 opacity-75"></span>
              <span className="relative inline-flex rounded-full h-2 w-2 bg-teal-500"></span>
            </span>
            Now powered by Qwen 2 72B
          </div>
          
          <h1 className="text-5xl md:text-7xl font-bold tracking-tight text-slate-900 mb-6 max-w-4xl mx-auto leading-[1.1] animate-in fade-in slide-in-from-bottom-6 duration-700 delay-100">
            Automated Clinical <br />
            <span className="text-transparent bg-clip-text bg-gradient-to-r from-teal-600 to-blue-600">Competency Assessment</span>
          </h1>
          
          <p className="text-lg md:text-xl text-slate-500 mb-10 max-w-2xl mx-auto leading-relaxed animate-in fade-in slide-in-from-bottom-6 duration-700 delay-200">
            Enhance medical education with AI-driven grading rooted in clinical protocols. 
            Provide instant, structured feedback to students and deep analytics for methodologists.
          </p>
          
          <div className="flex flex-col sm:flex-row gap-4 justify-center items-center animate-in fade-in slide-in-from-bottom-6 duration-700 delay-300">
            <Link 
              to="/dashboard" 
              className="px-8 py-4 rounded-xl bg-teal-600 text-white font-bold text-lg shadow-lg shadow-teal-600/20 hover:bg-teal-700 hover:shadow-xl hover:shadow-teal-600/30 transition-all transform hover:-translate-y-0.5 w-full sm:w-auto"
            >
              Enter Dashboard
            </Link>
          </div>
        </div>

        {/* Abstract Background Decoration */}
        <div className="absolute top-0 left-1/2 -translate-x-1/2 w-full h-full -z-10 pointer-events-none">
          <div className="absolute top-20 left-[10%] w-72 h-72 bg-teal-200/20 rounded-full blur-3xl mix-blend-multiply filter opacity-70 animate-blob"></div>
          <div className="absolute top-20 right-[10%] w-72 h-72 bg-blue-200/20 rounded-full blur-3xl mix-blend-multiply filter opacity-70 animate-blob animation-delay-2000"></div>
          <div className="absolute -bottom-8 left-[20%] w-72 h-72 bg-indigo-200/20 rounded-full blur-3xl mix-blend-multiply filter opacity-70 animate-blob animation-delay-4000"></div>
        </div>
      </div>
      
      {/* Simple Footer */}
      <div className="py-6 text-center text-slate-400 text-sm">
        © 2024 Clinical Assessment Systems. All rights reserved.
      </div>
    </div>
  );
};

export default Landing;
