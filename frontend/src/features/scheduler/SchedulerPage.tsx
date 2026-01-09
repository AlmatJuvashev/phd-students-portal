
import React, { useState, useEffect, useMemo } from 'react';
import { 
  Calendar as CalendarIcon, ChevronLeft, ChevronRight, 
  ZoomIn, ZoomOut, Filter, MoreHorizontal, GripVertical,
  AlertTriangle, CheckCircle2, RefreshCw, Zap, Settings,
  Users, Layers, ArrowRight, Gauge, Search, Play,
  Building, MapPin, GraduationCap, FileText, Ban,
  X, Trash2, Calendar, Clock, Construction, History,
  LayoutList, BarChart3, ChevronDown, CheckSquare, Briefcase, BookOpen
} from 'lucide-react';
import { useNavigate } from 'react-router-dom';
import { Button, Badge, IconButton } from '@/features/admin/components/AdminUI';
import { 
  getRooms, getScheduleEvents, updateScheduleEvent, getBuildings, getGlobalCourses, addScheduleEvent,
  Room, ScheduleEvent, getActiveCourses, ActiveCourse, GlobalCourse, getDepartments, Department
} from './data/opsData';
import { cn } from '@/lib/utils';
import { motion, AnimatePresence } from 'framer-motion';
import { addDays, addWeeks, addMonths, format, isSameDay, startOfWeek, setHours, setMinutes, isWithinInterval, startOfMonth, endOfMonth, eachDayOfInterval, differenceInMinutes } from 'date-fns';

// --- Types & Constants ---

const DAY_START_HOUR = 8; // 8 AM
const DAY_END_HOUR = 20; // 8 PM
const TOTAL_MINUTES = (DAY_END_HOUR - DAY_START_HOUR) * 60;

interface SchedulerPageProps {}

interface OptimizerConfig {
  utilizationWeight: number; // 0-100
  satisfactionWeight: number; // 0-100
  allowOvertime: boolean;
  prioritizeBuildings: boolean;
}

type GroupBy = 'room' | 'instructor' | 'course';
type ViewType = 'timeline' | 'utilization';
type TimeRange = 'day' | 'week' | 'month' | 'quarter';
type Term = 'Fall 2024' | 'Spring 2025' | 'Summer 2025' | 'Fall 2025' | 'Spring 2026';

// --- Components ---

const OptimizationPanel = ({ config, setConfig, onRun, isRunning, score }: { config: OptimizerConfig, setConfig: any, onRun: () => void, isRunning: boolean, score: number }) => (
  <div className="absolute bottom-6 right-6 w-80 bg-slate-900 text-white rounded-2xl shadow-2xl border border-slate-700 overflow-hidden z-50 animate-in slide-in-from-right-10 fade-in duration-500">
     <div className="p-4 border-b border-slate-700 bg-slate-800/50 flex justify-between items-center">
        <div className="flex items-center gap-2 text-sm font-bold">
           <Zap size={16} className="text-yellow-400 fill-yellow-400" /> Auto-Scheduler
        </div>
        <div className="text-[10px] font-mono text-slate-400">SOLVER_V2.1</div>
     </div>
     
     <div className="p-5 space-y-5">
        <div className="space-y-3">
           <div className="flex justify-between text-xs font-bold uppercase tracking-wider text-slate-400">
              <span>Optimization Goal</span>
           </div>
           
           <div className="space-y-4">
              <div>
                 <div className="flex justify-between text-[10px] font-bold mb-1.5">
                    <span className="text-indigo-400">Space Efficiency</span>
                    <span>{config.utilizationWeight}%</span>
                 </div>
                 <input 
                   type="range" 
                   value={config.utilizationWeight}
                   onChange={(e) => setConfig({...config, utilizationWeight: parseInt(e.target.value), satisfactionWeight: 100 - parseInt(e.target.value)})}
                   className="w-full h-1.5 bg-slate-700 rounded-full appearance-none cursor-pointer accent-indigo-500"
                 />
              </div>
              
              <div>
                 <div className="flex justify-between text-[10px] font-bold mb-1.5">
                    <span className="text-emerald-400">Professor Preferences</span>
                    <span>{config.satisfactionWeight}%</span>
                 </div>
                 <input 
                   type="range" 
                   value={config.satisfactionWeight}
                   onChange={(e) => setConfig({...config, satisfactionWeight: parseInt(e.target.value), utilizationWeight: 100 - parseInt(e.target.value)})}
                   className="w-full h-1.5 bg-slate-700 rounded-full appearance-none cursor-pointer accent-emerald-500"
                 />
              </div>

              <div className="flex items-center justify-between pt-2">
                 <span className="text-[10px] text-slate-300">Minimize Building Hops</span>
                 <input 
                   type="checkbox" 
                   checked={config.prioritizeBuildings}
                   onChange={(e) => setConfig({...config, prioritizeBuildings: e.target.checked})}
                   className="rounded border-slate-600 bg-slate-700 text-indigo-500 focus:ring-indigo-500/50"
                 />
              </div>
           </div>
        </div>

        <div className="p-3 bg-slate-800 rounded-xl border border-slate-700 flex items-center justify-between">
           <div className="flex items-center gap-3">
              <Gauge size={20} className="text-slate-400" />
              <div>
                 <div className="text-[10px] text-slate-400 uppercase font-bold">Fitness Score</div>
                 <div className="text-lg font-black">{score.toFixed(1)}</div>
              </div>
           </div>
           <div className={cn("h-2 w-2 rounded-full", score > 80 ? "bg-emerald-500 shadow-[0_0_8px_rgba(16,185,129,0.6)]" : "bg-amber-500")} />
        </div>

        <button 
          onClick={onRun}
          disabled={isRunning}
          className="w-full py-3 bg-indigo-600 hover:bg-indigo-500 text-white rounded-xl font-bold text-sm shadow-lg shadow-indigo-900/20 transition-all flex items-center justify-center gap-2 disabled:opacity-50 disabled:cursor-not-allowed"
        >
           {isRunning ? <RefreshCw size={16} className="animate-spin" /> : <Play size={16} fill="currentColor" />}
           {isRunning ? 'Solving...' : 'Run Optimization'}
        </button>
     </div>
  </div>
);

const FilterSidebar = ({ 
  isOpen, 
  onClose,
  departments,
  buildings,
  filters,
  setFilters 
}: { 
  isOpen: boolean, 
  onClose: () => void,
  departments: Department[],
  buildings: any[],
  filters: any,
  setFilters: (f: any) => void
}) => {
  if (!isOpen) return null;

  return (
    <div className="absolute top-0 left-0 bottom-0 w-64 bg-white border-r border-slate-200 z-40 shadow-xl overflow-y-auto">
       <div className="p-4 border-b border-slate-100 flex justify-between items-center">
          <h3 className="font-bold text-slate-800 text-sm flex items-center gap-2"><Filter size={16} /> Filters</h3>
          <IconButton icon={X} size="sm" onClick={onClose} />
       </div>
       
       <div className="p-4 space-y-6">
          {/* Departments */}
          <div>
             <h4 className="text-[10px] font-black text-slate-400 uppercase tracking-widest mb-3">Departments</h4>
             <div className="space-y-2">
                {departments.map(dept => (
                   <label key={dept.id} className="flex items-center gap-2 cursor-pointer hover:bg-slate-50 p-1 rounded transition-colors">
                      <input 
                        type="checkbox" 
                        checked={filters.departments.includes(dept.id)}
                        onChange={(e) => {
                           if(e.target.checked) setFilters({...filters, departments: [...filters.departments, dept.id]});
                           else setFilters({...filters, departments: filters.departments.filter((id: string) => id !== dept.id)});
                        }}
                        className="rounded border-slate-300 text-indigo-600 focus:ring-indigo-500"
                      />
                      <span className="text-sm font-medium text-slate-700">{dept.name}</span>
                      <span className="w-2 h-2 rounded-full ml-auto" style={{ backgroundColor: dept.color }} />
                   </label>
                ))}
             </div>
          </div>

          {/* Buildings */}
          <div>
             <h4 className="text-[10px] font-black text-slate-400 uppercase tracking-widest mb-3">Buildings & Floors</h4>
             <div className="space-y-4">
                {buildings.map(b => (
                   <div key={b.id}>
                      <div className="flex items-center gap-2 mb-2 font-bold text-xs text-slate-800">
                         <Building size={12} className="text-slate-400" /> {b.name}
                      </div>
                      <div className="pl-6 space-y-1">
                         {[1, 2, 3, 4, 5].map(floor => (
                            <label key={`${b.id}_f${floor}`} className="flex items-center gap-2 cursor-pointer">
                               <input 
                                 type="checkbox" 
                                 checked={filters.floors.includes(`${b.id}_${floor}`)}
                                 onChange={(e) => {
                                    const val = `${b.id}_${floor}`;
                                    if(e.target.checked) setFilters({...filters, floors: [...filters.floors, val]});
                                    else setFilters({...filters, floors: filters.floors.filter((f: string) => f !== val)});
                                 }}
                                 className="rounded border-slate-300 text-indigo-600 focus:ring-indigo-500"
                               />
                               <span className="text-xs text-slate-600">Floor {floor}</span>
                            </label>
                         ))}
                      </div>
                   </div>
                ))}
             </div>
          </div>
          
          <Button variant="outline" size="sm" className="w-full" onClick={() => setFilters({ departments: [], floors: [] })}>Clear All</Button>
       </div>
    </div>
  );
};

const UtilizationCell = ({ utilization }: { utilization: number }) => {
  let color = "bg-slate-100";
  let text = "text-slate-500";
  
  if (utilization > 0) {
     if (utilization < 30) { color = "bg-emerald-100"; text = "text-emerald-700"; }
     else if (utilization < 70) { color = "bg-indigo-100"; text = "text-indigo-700"; }
     else if (utilization < 90) { color = "bg-amber-100"; text = "text-amber-700"; }
     else { color = "bg-red-100"; text = "text-red-700"; }
  }

  return (
    <div className={cn("h-full w-full flex flex-col justify-center items-center text-xs font-bold transition-all hover:brightness-95", color, text)}>
       {utilization > 0 ? `${utilization}%` : '-'}
    </div>
  );
};

export const SchedulerPage: React.FC<SchedulerPageProps> = () => {
  const navigate = useNavigate();
  const onNavigate = (path: string) => navigate(path);
  const [currentDate, setCurrentDate] = useState(new Date());
  const [rooms, setRooms] = useState<Room[]>([]);
  const [events, setEvents] = useState<ScheduleEvent[]>([]);
  const [activeCourses, setActiveCourses] = useState<ActiveCourse[]>([]);
  const [departments, setDepartments] = useState<Department[]>([]);
  const [buildings, setBuildings] = useState<any[]>([]); // Using any for mock buildings for simplicity in this context
  
  // Controls
  const [viewType, setViewType] = useState<ViewType>('timeline');
  const [timeRange, setTimeRange] = useState<TimeRange>('day');
  const [groupBy, setGroupBy] = useState<GroupBy>('room');
  const [zoomLevel, setZoomLevel] = useState(2.5);
  const [selectedTerm, setSelectedTerm] = useState<Term>('Fall 2025');
  
  // Filters
  const [isFilterOpen, setIsFilterOpen] = useState(false);
  const [filters, setFilters] = useState({
     departments: [] as string[],
     floors: [] as string[] // Format: "bldId_floorNum"
  });

  const [isOptimizerRunning, setIsOptimizerRunning] = useState(false);
  const [optimizerConfig, setOptimizerConfig] = useState<OptimizerConfig>({ utilizationWeight: 60, satisfactionWeight: 40, allowOvertime: false, prioritizeBuildings: true });
  
  useEffect(() => {
    setRooms(getRooms());
    setEvents(getScheduleEvents());
    setActiveCourses(getActiveCourses());
    setDepartments(getDepartments());
    setBuildings(getBuildings());
  }, []);

  // --- Date Range Calculation ---
  const { startDate, endDate, days } = useMemo(() => {
     let start = currentDate;
     let end = currentDate;
     let dayList: Date[] = [];

     if (timeRange === 'day') {
        start = currentDate;
        end = currentDate;
        dayList = [currentDate];
     } else if (timeRange === 'week') {
        start = startOfWeek(currentDate, { weekStartsOn: 1 });
        end = addDays(start, 6);
        dayList = eachDayOfInterval({ start, end });
     } else if (timeRange === 'month') {
        start = startOfMonth(currentDate);
        end = endOfMonth(currentDate);
        dayList = eachDayOfInterval({ start, end }); // Can be huge
     } else if (timeRange === 'quarter') {
        start = startOfMonth(currentDate);
        end = addMonths(start, 3);
        dayList = []; // Quarters usually just aggregated
     }

     return { startDate: start, endDate: end, days: dayList };
  }, [currentDate, timeRange]);

  // --- Filtering Logic ---
  const filteredRooms = useMemo(() => {
     return rooms.filter(r => {
        // Department Filter
        if (filters.departments.length > 0 && (!r.departmentId || !filters.departments.includes(r.departmentId))) {
           return false;
        }
        // Floor/Building Filter
        if (filters.floors.length > 0) {
           const floorKey = `${r.buildingId}_${r.floor}`;
           if (!filters.floors.includes(floorKey)) return false;
        }
        return true;
     });
  }, [rooms, filters]);

  // --- Grouping Logic (Timeline) ---
  const groupedResources = useMemo(() => {
     const groups: Record<string, { id: string, title: string, items: any[] }> = {};

     if (groupBy === 'room') {
        filteredRooms.forEach(r => {
           const key = `${r.buildingName} - Floor ${r.floor}`;
           if (!groups[key]) groups[key] = { id: key, title: key, items: [] };
           groups[key].items.push(r);
        });
     } else if (groupBy === 'instructor') {
        // Group by Instructor Name
        const uniqueInstructors = Array.from(new Set(events.filter(e => e.instructorName).map(e => e.instructorName!)));
        groups['Faculty'] = { id: 'faculty', title: 'Faculty Members', items: uniqueInstructors.map(name => ({ id: name, name, type: 'instructor' })) };
     } else {
        // Group by Course
        const uniqueCourses = Array.from(new Set(events.map(e => e.targetName)));
        groups['Courses'] = { id: 'courses', title: 'Active Courses', items: uniqueCourses.map(name => ({ id: name, name, type: 'course' })) };
     }

     return groups;
  }, [filteredRooms, events, groupBy]);

  // --- Utilization Logic (Analytics) ---
  const utilizationData = useMemo(() => {
     // Only for Month/Quarter views generally
     const data: Record<string, Record<string, number>> = {}; // { roomId: { '2024-10-01': 50, 'total': 45 } }
     
     filteredRooms.forEach(room => {
        data[room.id] = { total: 0 };
        
        // Find events for this room in range
        const roomEvents = events.filter(e => e.resourceId === room.id && isWithinInterval(new Date(e.date), { start: startDate, end: endDate }));
        
        // Aggregate by day or week depending on view
        roomEvents.forEach(e => {
           const dayKey = format(new Date(e.date), 'yyyy-MM-dd');
           const durationHours = (e.durationMinutes || 60) / 60;
           
           data[room.id][dayKey] = (data[room.id][dayKey] || 0) + durationHours;
           data[room.id].total += durationHours;
        });
     });
     return data;
  }, [filteredRooms, events, startDate, endDate]);

  const getUtilizationPct = (bookedHours: number) => {
     const availableHours = (DAY_END_HOUR - DAY_START_HOUR); 
     return Math.min(100, Math.round((bookedHours / availableHours) * 100));
  };

  // --- Optimizer Logic ---
  const globalScore = useMemo(() => {
     let score = 85; // Base score
     // Mock calculation based on resource conflict
     const conflicts = events.filter(e => e.warnings && e.warnings.length > 0).length;
     score -= conflicts * 5;
     
     // Utilization impact
     const totalEvents = events.length;
     const scheduledEvents = events.filter(e => e.resourceId).length;
     const unscheduledPenalty = (totalEvents - scheduledEvents) * 2;
     score -= unscheduledPenalty;

     return Math.max(0, Math.min(100, score));
  }, [events]);

  const handleOptimize = () => {
     setIsOptimizerRunning(true);
     // Simulate a solver running
     setTimeout(() => {
        // Mock optimization: simply assign unscheduled events to the first available room
        const newEvents = [...events];
        const unscheduled = newEvents.filter(e => !e.resourceId);
        
        unscheduled.forEach((evt) => {
           // Simple heuristic: assign to random room for visual effect
           if (rooms.length > 0) {
              const randomRoom = rooms[Math.floor(Math.random() * rooms.length)];
              evt.resourceId = randomRoom.id;
              evt.location = randomRoom.name;
              // Ensure it has a time if missing
              if (!evt.startTime) evt.startTime = '09:00';
           }
        });
        
        setEvents(newEvents);
        setIsOptimizerRunning(false);
     }, 1500);
  };

  // --- Handlers ---
  const handleNav = (dir: 'prev' | 'next') => {
     if (timeRange === 'day') setCurrentDate(addDays(currentDate, dir === 'next' ? 1 : -1));
     if (timeRange === 'week') setCurrentDate(addWeeks(currentDate, dir === 'next' ? 1 : -1));
     if (timeRange === 'month') setCurrentDate(addMonths(currentDate, dir === 'next' ? 1 : -1));
     if (timeRange === 'quarter') setCurrentDate(addMonths(currentDate, dir === 'next' ? 3 : -3));
  };

  return (
    <div className="flex h-full bg-slate-100 relative">
       {/* Filter Sidebar (Collapsible) */}
       <FilterSidebar 
         isOpen={isFilterOpen} 
         onClose={() => setIsFilterOpen(false)}
         departments={departments}
         buildings={buildings}
         filters={filters}
         setFilters={setFilters}
       />

       {/* Main Content */}
       <div className="flex-1 flex flex-col min-w-0 transition-all duration-300" style={{ marginLeft: isFilterOpen ? '256px' : '0' }}>
          
          {/* Top Toolbar */}
          <div className="h-16 bg-white border-b border-slate-200 px-6 flex items-center justify-between sticky top-0 z-30 shadow-sm">
             <div className="flex items-center gap-4">
                <IconButton icon={Filter} onClick={() => setIsFilterOpen(!isFilterOpen)} className={isFilterOpen ? "bg-indigo-100 text-indigo-600" : ""} />
                
                {/* Date Nav */}
                <div className="flex items-center gap-2 bg-slate-100 rounded-lg p-1">
                   <button onClick={() => handleNav('prev')} className="p-1 hover:bg-white rounded shadow-sm"><ChevronLeft size={16} /></button>
                   <span className="text-xs font-bold px-3 min-w-[120px] text-center">
                      {timeRange === 'day' && format(currentDate, 'MMM d, yyyy')}
                      {timeRange === 'week' && `Week of ${format(startDate, 'MMM d')}`}
                      {timeRange === 'month' && format(currentDate, 'MMMM yyyy')}
                      {timeRange === 'quarter' && `Q${Math.floor(currentDate.getMonth()/3)+1} ${currentDate.getFullYear()}`}
                   </span>
                   <button onClick={() => handleNav('next')} className="p-1 hover:bg-white rounded shadow-sm"><ChevronRight size={16} /></button>
                </div>

                <div className="h-6 w-px bg-slate-200" />

                {/* View Switcher */}
                <div className="flex bg-slate-100 p-1 rounded-lg">
                   {['day', 'week', 'month', 'quarter'].map(r => (
                      <button
                        key={r}
                        onClick={() => { setTimeRange(r as TimeRange); if(r==='month' || r==='quarter') setViewType('utilization'); else setViewType('timeline'); }}
                        className={cn("px-3 py-1.5 text-xs font-bold rounded capitalize transition-all", timeRange === r ? "bg-white text-indigo-600 shadow-sm" : "text-slate-500 hover:text-slate-700")}
                      >
                         {r}
                      </button>
                   ))}
                </div>
             </div>

             <div className="flex items-center gap-4">
                <div className="flex items-center gap-2">
                   <span className="text-[10px] font-bold text-slate-400 uppercase">View:</span>
                   <select 
                     value={viewType} 
                     onChange={(e) => setViewType(e.target.value as ViewType)}
                     className="bg-transparent text-xs font-bold text-slate-700 outline-none cursor-pointer border border-slate-200 rounded px-2 py-1 hover:border-indigo-300"
                   >
                      <option value="timeline">Timeline (Operational)</option>
                      <option value="utilization">Utilization (Analytics)</option>
                   </select>
                </div>

                <div className="flex items-center gap-2">
                   <span className="text-[10px] font-bold text-slate-400 uppercase">Group By:</span>
                   <select 
                     value={groupBy} 
                     onChange={(e) => setGroupBy(e.target.value as GroupBy)}
                     className="bg-transparent text-xs font-bold text-slate-700 outline-none cursor-pointer border border-slate-200 rounded px-2 py-1 hover:border-indigo-300"
                   >
                      <option value="room">Room (Building/Floor)</option>
                      <option value="instructor">Instructor</option>
                      <option value="course">Course</option>
                   </select>
                </div>
                
                <Button size="sm" variant="secondary" icon={BarChart3} onClick={() => setViewType(viewType === 'timeline' ? 'utilization' : 'timeline')}>
                   {viewType === 'timeline' ? 'Analyze' : 'Schedule'}
                </Button>
             </div>
          </div>

          {/* MAIN CONTENT AREA */}
          <div className="flex-1 overflow-auto relative bg-slate-50">
             
             {/* TIMELINE VIEW */}
             {viewType === 'timeline' && (
                <div className="min-w-max pb-20">
                   {/* Time Header */}
                   <div className="sticky top-0 z-20 bg-white border-b border-slate-200 flex h-10 shadow-sm">
                      <div className="w-56 flex-shrink-0 border-r border-slate-200 bg-slate-50 p-2 text-xs font-bold text-slate-500 uppercase flex items-center sticky left-0 z-30">
                         {groupBy === 'room' ? 'Location' : groupBy === 'instructor' ? 'Faculty' : 'Course'}
                      </div>
                      <div className="flex relative" style={{ width: TOTAL_MINUTES * zoomLevel }}>
                         {Array.from({ length: DAY_END_HOUR - DAY_START_HOUR }).map((_, i) => (
                            <div 
                              key={i} 
                              className="absolute top-0 bottom-0 border-l border-slate-100 pl-2 text-[10px] font-bold text-slate-400"
                              style={{ left: i * 60 * zoomLevel }}
                            >
                               {DAY_START_HOUR + i}:00
                            </div>
                         ))}
                      </div>
                   </div>

                   {/* Grouped Rows */}
                   {Object.entries(groupedResources).map(([groupName, groupData]: [string, any]) => (
                      <React.Fragment key={groupData.id}>
                         <div className="sticky left-0 right-0 z-10 bg-slate-100/90 backdrop-blur-sm border-b border-slate-200 px-4 py-1 text-[10px] font-black text-slate-500 uppercase tracking-widest flex items-center gap-2">
                            {groupBy === 'room' ? <Building size={10} /> : groupBy === 'instructor' ? <Users size={10} /> : <BookOpen size={10} />}
                            {groupName}
                         </div>
                         
                         {groupData.items.map((resource: any) => {
                            // Filter events for this resource row
                            const rowEvents = events.filter(e => {
                               if (groupBy === 'room') return e.resourceId === resource.id;
                               if (groupBy === 'instructor') return e.instructorName === resource.name;
                               return e.targetName === resource.name;
                            }).filter(e => isSameDay(new Date(e.date), currentDate)); // Filter by current view day

                            return (
                               <div key={resource.id} className="flex border-b border-slate-200 bg-white hover:bg-slate-50/50 transition-colors h-20 relative group/row">
                                  <div className="w-56 flex-shrink-0 border-r border-slate-200 p-3 sticky left-0 z-10 bg-white group-hover/row:bg-slate-50 transition-colors flex flex-col justify-center">
                                     <div className="font-bold text-slate-800 text-sm truncate">{resource.name}</div>
                                     {resource.capacity && <div className="text-[10px] text-slate-400">Cap: {resource.capacity}</div>}
                                  </div>
                                  
                                  <div className="relative h-full" style={{ width: TOTAL_MINUTES * zoomLevel }}>
                                     {/* Grid Lines */}
                                     {Array.from({ length: DAY_END_HOUR - DAY_START_HOUR }).map((_, i) => (
                                        <div 
                                          key={i} 
                                          className="absolute top-0 bottom-0 border-l border-slate-100"
                                          style={{ left: i * 60 * zoomLevel }}
                                        />
                                     ))}

                                     {/* Events */}
                                     {rowEvents.map(event => {
                                        const [h, m] = (event.startTime || '09:00').split(':').map(Number);
                                        const startMins = (h - DAY_START_HOUR) * 60 + m;
                                        return (
                                           <div 
                                             key={event.id}
                                             className={cn(
                                                "absolute top-2 bottom-2 rounded-lg border shadow-sm p-2 text-[10px] overflow-hidden hover:z-20 hover:shadow-md cursor-pointer transition-all",
                                                event.type === 'lecture' ? "bg-indigo-50 border-indigo-200 text-indigo-900" : "bg-amber-50 border-amber-200 text-amber-900"
                                             )}
                                             style={{ left: startMins * zoomLevel, width: (event.durationMinutes || 60) * zoomLevel }}
                                          >
                                             <div className="font-bold truncate">{event.title}</div>
                                             <div className="opacity-70 truncate">{event.startTime}</div>
                                          </div>
                                        );
                                     })}
                                  </div>
                               </div>
                            );
                         })}
                      </React.Fragment>
                   ))}
                </div>
             )}

             {/* UTILIZATION VIEW (Heatmap Table) */}
             {viewType === 'utilization' && (
                <div className="p-8">
                   <div className="bg-white rounded-2xl shadow-sm border border-slate-200 overflow-hidden">
                      <div className="overflow-x-auto">
                         <table className="w-full text-left border-collapse">
                            <thead>
                               <tr>
                                  <th className="p-4 border-b border-r border-slate-200 bg-slate-50 text-xs font-bold text-slate-500 uppercase sticky left-0 z-20 w-64">
                                     Resource / Room
                                  </th>
                                  {days.map(day => (
                                     <th key={day.toISOString()} className="p-2 border-b border-slate-200 bg-slate-50 text-center min-w-[60px]">
                                        <div className="text-[10px] font-bold text-slate-400 uppercase">{format(day, 'EEE')}</div>
                                        <div className="text-xs font-bold text-slate-700">{format(day, 'd')}</div>
                                     </th>
                                  ))}
                                  <th className="p-4 border-b border-l border-slate-200 bg-slate-50 text-xs font-bold text-slate-500 uppercase text-center sticky right-0 z-20 w-24">
                                     Avg Load
                                  </th>
                               </tr>
                            </thead>
                            <tbody>
                               {filteredRooms.map(room => {
                                  const roomStats = utilizationData[room.id] || { total: 0 };
                                  const avgLoad = roomStats.total / days.length; // rough avg hours/day
                                  const avgPct = Math.min(100, Math.round((avgLoad / (DAY_END_HOUR - DAY_START_HOUR)) * 100));

                                  return (
                                     <tr key={room.id} className="border-b border-slate-100 hover:bg-slate-50">
                                        <td className="p-3 border-r border-slate-200 sticky left-0 bg-white group-hover:bg-slate-50 z-10">
                                           <div className="font-bold text-sm text-slate-800">{room.name}</div>
                                           <div className="text-[10px] text-slate-400">{room.buildingName}, Fl {room.floor}</div>
                                        </td>
                                        {days.map(day => {
                                           const dayKey = format(day, 'yyyy-MM-dd');
                                           const hours = roomStats[dayKey] || 0;
                                           const pct = getUtilizationPct(hours);
                                           return (
                                              <td key={dayKey} className="p-1 h-12 border-r border-slate-50">
                                                 <UtilizationCell utilization={pct} />
                                              </td>
                                           );
                                        })}
                                        <td className="p-3 border-l border-slate-200 sticky right-0 bg-white group-hover:bg-slate-50 z-10 text-center">
                                           <div className={cn("inline-block px-2 py-1 rounded font-bold text-xs", avgPct > 80 ? "bg-red-100 text-red-700" : "bg-emerald-100 text-emerald-700")}>
                                              {avgPct}%
                                           </div>
                                        </td>
                                     </tr>
                                  );
                               })}
                            </tbody>
                         </table>
                      </div>
                   </div>
                </div>
             )}
          </div>

          {/* Floating Optimizer Card */}
          {viewType === 'timeline' && (
             <OptimizationPanel 
               config={optimizerConfig} 
               setConfig={setOptimizerConfig} 
               onRun={handleOptimize} 
               isRunning={isOptimizerRunning} 
               score={globalScore}
             />
          )}
       </div>
    </div>
  );
};

export default SchedulerPage;
