
import React, { useState } from 'react';
import { NavLink, Outlet, useLocation, Link } from 'react-router-dom';
import { 
  LayoutDashboard, 
  Settings, 
  Users, 
  GraduationCap, 
  Menu, 
  X, 
  BrainCircuit, 
  Activity,
  ClipboardList
} from 'lucide-react';
import { UserRole } from '../../types';

interface LayoutProps {
  role: UserRole;
  setRole: (role: UserRole) => void;
  currentModel: string;
}

const Layout: React.FC<LayoutProps> = ({ role, setRole, currentModel }) => {
  const [isSidebarOpen, setIsSidebarOpen] = useState(false);
  const location = useLocation();

  const toggleSidebar = () => setIsSidebarOpen(!isSidebarOpen);

  const navItems = [
    { name: 'Dashboard', path: '/dashboard', icon: LayoutDashboard, roles: ['admin', 'teacher'] },
    { name: 'Training Ground', path: '/training', icon: GraduationCap, roles: ['student'] }, // Removed for teacher
    { name: 'Grading Queue', path: '/grading', icon: ClipboardList, roles: ['teacher'] }, // Added for teacher
    { name: 'Score Audit', path: '/audit', icon: Users, roles: ['admin', 'teacher'] },
    { name: 'Model & RAG', path: '/config', icon: BrainCircuit, roles: ['admin'] },
  ];

  return (
    <div className="min-h-screen bg-slate-50 flex font-sans text-slate-800">
      {/* Mobile Sidebar Overlay */}
      {isSidebarOpen && (
        <div 
          className="fixed inset-0 bg-slate-900/50 z-20 lg:hidden"
          onClick={() => setIsSidebarOpen(false)}
        />
      )}

      {/* Sidebar */}
      <aside className={`
        fixed inset-y-0 left-0 z-30 w-64 bg-white border-r border-slate-200 transform transition-transform duration-200 ease-in-out
        lg:translate-x-0 lg:static lg:inset-0
        ${isSidebarOpen ? 'translate-x-0' : '-translate-x-full'}
      `}>
        <div className="flex flex-col h-full">
          {/* Logo */}
          <Link to="/dashboard" className="h-16 flex items-center px-6 border-b border-slate-100 hover:bg-slate-50 transition-colors">
            <Activity className="h-6 w-6 text-teal-600 mr-2" />
            <span className="font-bold text-lg text-slate-900">ClinAssessor</span>
          </Link>

          {/* Navigation */}
          <nav className="flex-1 px-4 py-6 space-y-1">
            {navItems.filter(item => item.roles.includes(role)).map((item) => {
              const isActive = location.pathname === item.path || (item.path !== '/' && location.pathname.startsWith(item.path));
              return (
                <NavLink
                  key={item.name}
                  to={item.path}
                  onClick={() => setIsSidebarOpen(false)}
                  className={`
                    flex items-center px-4 py-3 text-sm font-medium rounded-lg transition-colors
                    ${isActive 
                      ? 'bg-teal-50 text-teal-700' 
                      : 'text-slate-600 hover:bg-slate-50 hover:text-slate-900'}
                  `}
                >
                  <item.icon className={`mr-3 h-5 w-5 ${isActive ? 'text-teal-600' : 'text-slate-400'}`} />
                  {item.name}
                </NavLink>
              );
            })}
          </nav>

          {/* User Role Switcher (Mock Auth) */}
          <div className="p-4 border-t border-slate-100">
            <label className="block text-xs font-semibold text-slate-400 uppercase tracking-wider mb-2">
              Viewing As
            </label>
            <select
              value={role}
              onChange={(e) => setRole(e.target.value as UserRole)}
              className="block w-full rounded-md border-slate-300 py-2 pl-3 pr-8 text-sm focus:border-teal-500 focus:outline-none focus:ring-teal-500 bg-slate-50"
            >
              <option value="admin">Admin</option>
              <option value="teacher">Teacher</option>
              <option value="student">Student</option>
            </select>
          </div>
        </div>
      </aside>

      {/* Main Content */}
      <div className="flex-1 flex flex-col overflow-hidden">
        {/* Top Navbar */}
        <header className="h-16 bg-white border-b border-slate-200 flex items-center justify-between px-4 sm:px-6 lg:px-8">
          <button
            onClick={toggleSidebar}
            className="lg:hidden p-2 rounded-md text-slate-400 hover:text-slate-500 hover:bg-slate-100"
          >
            <Menu className="h-6 w-6" />
          </button>

          <div className="flex-1 flex justify-between items-center ml-4 lg:ml-0">
            <h1 className="text-xl font-semibold text-slate-800 hidden sm:block">
               {navItems.find(i => i.path === location.pathname)?.name || 'Clinical Answer Assessor'}
            </h1>
            
            <div className="flex items-center space-x-4">
               {/* Model Status Indicator - Hidden for Teachers */}
              {role !== 'teacher' && (
                <div className="hidden md:flex items-center px-3 py-1 bg-slate-100 rounded-full border border-slate-200">
                  <BrainCircuit className="h-3 w-3 text-teal-600 mr-2" />
                  <span className="text-xs font-medium text-slate-600">
                    Model: <span className="text-slate-900">{currentModel}</span>
                  </span>
                </div>
              )}
              
              <div className="h-8 w-8 rounded-full bg-teal-100 flex items-center justify-center text-teal-700 font-bold text-sm">
                {role.charAt(0).toUpperCase()}
              </div>
            </div>
          </div>
        </header>

        {/* Page Content */}
        <main className="flex-1 overflow-y-auto p-4 sm:p-6 lg:p-8">
          <Outlet />
        </main>
      </div>
    </div>
  );
};

export default Layout;
