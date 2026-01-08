
import React, { useState } from 'react';
import { Outlet, NavLink, useLocation, useNavigate } from "react-router-dom";
import { useAuth } from "@/contexts/AuthContext";
import { cn } from "@/lib/utils";
import {
  LayoutDashboard,
  Users,
  Monitor,
  Bell,
  BookOpen,
  Settings,
  ShieldCheck,
  Building,
  GraduationCap,
  CalendarClock,
  MapPin,
  Library,
  MessagesSquare,
  LogOut,
  ChevronDown,
  Search,
  Command,
  Menu,
  Sparkles,
  Layout,
  Map
} from "lucide-react";
import { useTranslation } from "react-i18next";
import { EduStudioHub } from "@/features/admin/components/EduStudioHub";
import { Button } from "@/features/admin/components/AdminUI";

type StudioType = 'academic' | 'people' | 'campus' | 'system';

interface SidebarItemProps {
    icon: any;
    label: string;
    path: string;
    collapsed?: boolean;
    badge?: string;
    disabled?: boolean;
}

const SidebarItem = ({ icon: Icon, label, path, collapsed, badge, disabled }: SidebarItemProps) => {
    const location = useLocation();
    const isActive = location.pathname === path || location.pathname.startsWith(path + '/');
    
    return (
        <NavLink
            to={path}
            onClick={(e) => disabled && e.preventDefault()}
            className={({ isActive: linkActive }) => cn(
                "flex items-center gap-3 px-3 py-2.5 rounded-xl text-sm font-medium transition-all group mb-1",
                isActive 
                    ? "bg-slate-900 text-white shadow-md" 
                    : "text-slate-500 hover:bg-slate-100 hover:text-slate-900",
                disabled && "opacity-50 cursor-not-allowed hover:bg-transparent",
                collapsed && "justify-center px-2"
            )}
            title={label}
        >
            <Icon size={18} className={cn(isActive ? "text-white" : "text-slate-400 group-hover:text-slate-600")} />
            {!collapsed && <span className="flex-1 text-left truncate">{label}</span>}
            {!collapsed && badge && <span className="text-[10px] bg-indigo-100 text-indigo-700 px-1.5 py-0.5 rounded font-bold">{badge}</span>}
        </NavLink>
    );
};

export function AdminLayout() {
  const { t } = useTranslation("common");
  const { user, logout } = useAuth();
  const navigate = useNavigate();
  const [collapsed, setCollapsed] = useState(false);
  const [activeStudio, setActiveStudio] = useState<StudioType>('academic');
  const [isHubOpen, setIsHubOpen] = useState(false);
  const [isMenuOpen, setIsMenuOpen] = useState(false); // Mobile menu

  // Mock studio switching logic or simple state
  // In a real app, this might be tied to a context or persisted
  
  const studios: { id: StudioType; name: string; desc: string; icon: any; color: string }[] = [
      { id: 'academic', name: 'Academic Studio', desc: 'Curriculum & Content', icon: BookOpen, color: 'bg-indigo-600' },
      { id: 'people', name: 'People Operations', desc: 'HR & Student Affairs', icon: Users, color: 'bg-emerald-600' },
      { id: 'campus', name: 'Campus Manager', desc: 'Facilities & Resources', icon: Building, color: 'bg-orange-600' },
      { id: 'system', name: 'System Admin', desc: 'Logs & Settings', icon: ShieldCheck, color: 'bg-slate-700' }
  ];

  const currentStudio = studios.find(s => s.id === activeStudio) || studios[0];

  return (
    <div className="flex h-screen bg-slate-50 font-sans">
      <EduStudioHub isOpen={isHubOpen} onClose={() => setIsHubOpen(false)} onNavigate={(path) => navigate(path)} />

      {/* Sidebar */}
      <aside 
        className={cn(
            "bg-white border-r border-slate-200 flex flex-col flex-shrink-0 z-30 transition-all duration-300",
            collapsed ? "w-20" : "w-64",
            // Mobile handling: usually hidden or absolute
            "hidden md:flex" 
        )}
      >
        {/* Org/Studio Selector */}
        <div className="h-20 flex items-center px-4 border-b border-slate-100">
           <div className="relative w-full group/studio">
                <button 
                    className="flex items-center gap-3 w-full p-2 hover:bg-slate-50 rounded-xl transition-colors text-left"
                    onClick={() => { /* Could toggle a dropdown here */ }}
                >
                    <div className={cn("w-10 h-10 rounded-lg flex items-center justify-center text-white font-black shadow-sm flex-shrink-0", currentStudio.color)}>
                        <currentStudio.icon size={20} />
                    </div>
                    {!collapsed && (
                        <div className="flex-1 min-w-0">
                            <div className="text-sm font-bold text-slate-900 truncate">{currentStudio.name}</div>
                            <div className="text-[10px] text-slate-500 font-medium truncate">{currentStudio.desc}</div>
                        </div>
                    )}
                    {!collapsed && <ChevronDown size={14} className="text-slate-400" />}
                </button>
                
                {/* Simple Hover Dropdown for Demo purposes - in prod use a proper Popover */}
                <div className="absolute top-full left-0 w-64 bg-white border border-slate-200 rounded-xl shadow-xl p-2 mt-2 hidden group-hover/studio:block z-50">
                    <div className="text-[10px] font-bold text-slate-400 uppercase px-2 py-1">Switch Studio</div>
                    {studios.map(s => (
                        <button 
                            key={s.id}
                            onClick={() => setActiveStudio(s.id)}
                            className={cn(
                                "flex items-center gap-3 w-full p-2 rounded-lg hover:bg-slate-50 transition-colors text-left",
                                activeStudio === s.id && "bg-slate-50"
                            )}
                        >
                            <div className={cn("w-8 h-8 rounded-lg flex items-center justify-center text-white shadow-sm flex-shrink-0", s.color)}>
                                <s.icon size={16} />
                            </div>
                            <div>
                                <div className="text-sm font-bold text-slate-900">{s.name}</div>
                                <div className="text-[10px] text-slate-500">{s.desc}</div>
                            </div>
                        </button>
                    ))}
                </div>
           </div>
        </div>

        {/* Nav Content based on Active Studio */}
        <div className="flex-1 overflow-y-auto py-6 px-3 space-y-6 custom-scrollbar">
           
           {/* ACADEMIC STUDIO */}
           {activeStudio === 'academic' && (
               <>
                   <div>
                       {!collapsed && <div className="px-3 text-[10px] font-black text-slate-400 uppercase tracking-widest mb-2">Overview</div>}
                       <SidebarItem icon={LayoutDashboard} label="Dashboard" path="/admin" collapsed={collapsed} />
                       <SidebarItem icon={Bell} label="Notifications" path="/admin/notifications" collapsed={collapsed} badge="9+" />
                   </div>
                   <div>
                       {!collapsed && <div className="px-3 text-[10px] font-black text-slate-400 uppercase tracking-widest mb-2">Curriculum</div>}
                       <SidebarItem icon={Map} label="Programs" path="/admin/programs" collapsed={collapsed} />
                       <SidebarItem icon={BookOpen} label="Courses" path="/admin/courses" collapsed={collapsed} />
                       <SidebarItem icon={BookOpen} label="Dictionaries" path="/admin/dictionaries" collapsed={collapsed} />
                   </div>
                   <div>
                       {!collapsed && <div className="px-3 text-[10px] font-black text-slate-400 uppercase tracking-widest mb-2">Delivery</div>}
                       <SidebarItem icon={CalendarClock} label="Scheduler AI" path="/admin/scheduler" collapsed={collapsed} badge="New" />
                       <SidebarItem icon={Library} label="Item Bank" path="/admin/item-banks" collapsed={collapsed} />
                   </div>
                   <div className="mt-4 px-2">
                       <button 
                           onClick={() => setIsHubOpen(true)}
                           className={cn(
                               "w-full bg-gradient-to-r from-indigo-600 to-violet-600 text-white p-3 rounded-xl shadow-lg shadow-indigo-200 hover:shadow-indigo-300 transition-all active:scale-95 flex items-center justify-center gap-2",
                               collapsed ? "aspect-square p-0" : ""
                           )}
                       >
                           <Sparkles size={18} />
                           {!collapsed && <span className="font-bold text-sm">Create New</span>}
                       </button>
                   </div>
               </>
           )}

           {/* PEOPLE STUDIO */}
           {activeStudio === 'people' && (
               <>
                   <div>
                       {!collapsed && <div className="px-3 text-[10px] font-black text-slate-400 uppercase tracking-widest mb-2">Directory</div>}
                       <SidebarItem icon={Users} label="All Users" path="/admin/users" collapsed={collapsed} />
                       <SidebarItem icon={GraduationCap} label="Students" path="/admin/students-monitor" collapsed={collapsed} />
                       <SidebarItem icon={Users} label="Advisors" path="/admin/advisors" collapsed={collapsed} />
                       <SidebarItem icon={Users} label="Contacts" path="/admin/contacts" collapsed={collapsed} />
                   </div>
                   <div>
                       {!collapsed && <div className="px-3 text-[10px] font-black text-slate-400 uppercase tracking-widest mb-2">Operations</div>}
                       <SidebarItem icon={Users} label="HR Management" path="/admin/hr" collapsed={collapsed} badge="New" />
                       <SidebarItem icon={GraduationCap} label="Enrollments" path="/admin/enrollments" collapsed={collapsed} />
                   </div>
               </>
           )}

           {/* CAMPUS STUDIO */}
           {activeStudio === 'campus' && (
               <>
                   <div>
                       {!collapsed && <div className="px-3 text-[10px] font-black text-slate-400 uppercase tracking-widest mb-2">Facilities</div>}
                       <SidebarItem icon={Building} label="Facilities" path="/admin/facilities" collapsed={collapsed} badge="New" />
                       <SidebarItem icon={MapPin} label="Rooms" path="/admin/rooms" collapsed={collapsed} />
                   </div>
                   <div>
                       {!collapsed && <div className="px-3 text-[10px] font-black text-slate-400 uppercase tracking-widest mb-2">Communication</div>}
                       <SidebarItem icon={MessagesSquare} label="Chat Rooms" path="/admin/chat-rooms" collapsed={collapsed} />
                   </div>
               </>
           )}

           {/* SYSTEM STUDIO */}
           {activeStudio === 'system' && (
               <>
                   <div>
                       {!collapsed && <div className="px-3 text-[10px] font-black text-slate-400 uppercase tracking-widest mb-2">Platform</div>}
                       <SidebarItem icon={LayoutDashboard} label="Tenants" path="/superadmin/tenants" collapsed={collapsed} />
                       <SidebarItem icon={ShieldCheck} label="Admins" path="/superadmin/admins" collapsed={collapsed} />
                   </div>
                   <div>
                       {!collapsed && <div className="px-3 text-[10px] font-black text-slate-400 uppercase tracking-widest mb-2">Configuration</div>}
                       <SidebarItem icon={Settings} label="Global Settings" path="/superadmin/settings" collapsed={collapsed} />
                       <SidebarItem icon={Layout} label="Logs" path="/superadmin/logs" collapsed={collapsed} />
                   </div>
               </>
           )}

        </div>

        {/* User Footer */}
        <div className="p-4 border-t border-slate-100">
           <div className="flex items-center gap-3 p-2 rounded-xl hover:bg-slate-50 transition-colors cursor-pointer" onClick={() => logout()}>
              <div className="w-9 h-9 rounded-full bg-slate-200 flex items-center justify-center text-xs font-bold text-slate-600 border-2 border-white shadow-sm">
                 {user?.email?.substring(0,2).toUpperCase()}
              </div>
              {!collapsed && (
                  <div className="flex-1 min-w-0">
                     <div className="text-sm font-bold text-slate-900 truncate">Admin User</div>
                     <div className="text-[10px] text-slate-500 truncate">{user?.email}</div>
                  </div>
              )}
              {!collapsed && <LogOut size={16} className="text-slate-400" />}
           </div>
        </div>
      </aside>

      {/* Main Content */}
      <div className="flex-1 flex flex-col min-w-0">
         {/* Top Bar */}
         <header className="h-20 bg-white border-b border-slate-200 px-6 flex items-center justify-between sticky top-0 z-20">
            <div className="flex items-center gap-4">
               <button className="md:hidden p-2 -ml-2 text-slate-500" onClick={() => setIsMenuOpen(!isMenuOpen)}><Menu size={20} /></button>
               {/* Mobile Menu would go here */}
               
               <div className="flex flex-col">
                   <h1 className="text-xl font-bold text-slate-900 tracking-tight">
                       {activeStudio === 'academic' && 'Academic Overview'}
                       {activeStudio === 'people' && 'People & Roster'}
                       {activeStudio === 'campus' && 'Facilities & Ops'}
                       {activeStudio === 'system' && 'System Configuration'}
                   </h1>
                   <div className="flex items-center gap-2 text-xs text-slate-500">
                      <span className="font-bold text-indigo-600">KazNMU Portal</span>
                      <span>•</span>
                      <span>{new Date().toLocaleDateString()}</span>
                   </div>
               </div>
            </div>

            <div className="flex items-center gap-4">
               <div className="relative hidden md:block group">
                  <Search size={16} className="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400 group-focus-within:text-indigo-500 transition-colors" />
                  <input 
                    type="text" 
                    placeholder="Global search..." 
                    className="h-10 pl-10 pr-12 rounded-xl bg-slate-100 border-transparent text-sm focus:bg-white focus:border-indigo-500 focus:ring-4 focus:ring-indigo-500/10 transition-all w-64 outline-none font-medium"
                  />
                  <div className="absolute right-3 top-1/2 -translate-y-1/2 flex items-center gap-1 pointer-events-none">
                     <Command size={10} className="text-slate-400" />
                     <span className="text-[10px] font-bold text-slate-400">K</span>
                  </div>
               </div>
               
               <div className="h-8 w-px bg-slate-200" />
               
               <button className="relative p-2.5 text-slate-400 hover:bg-slate-100 hover:text-indigo-600 rounded-xl transition-all">
                  <Bell size={20} />
                  <span className="absolute top-2.5 right-2.5 w-2.5 h-2.5 bg-red-500 rounded-full border-2 border-white shadow-sm" />
               </button>
            </div>
         </header>

         {/* Scroll Area */}
         <main className="flex-1 overflow-y-auto p-6 md:p-8 bg-slate-50 custom-scrollbar">
            <div className="max-w-7xl mx-auto space-y-8 pb-20">
               <Outlet />
            </div>
         </main>
      </div>
    </div>
  );
}

export default AdminLayout;
