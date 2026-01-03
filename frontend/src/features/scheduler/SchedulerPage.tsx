import React, { useState, useEffect, useMemo } from 'react';
import { 
  ChevronLeft, ChevronRight, Filter, Play, RefreshCw, 
  Building, Users, BookOpen, BarChart3, Loader2 
} from 'lucide-react';
import { Button } from '@/components/ui/button';
import { cn } from '@/lib/utils';
import { addDays, addWeeks, addMonths, format, isSameDay, startOfWeek, isWithinInterval, startOfMonth, endOfMonth, eachDayOfInterval } from 'date-fns';
import { 
  fetchTerms, fetchBuildings, fetchRooms, fetchOfferings, fetchSessions, optimizeSchedule 
} from './api';
import { 
  AcademicTerm, Building as BuildingType, Room, CourseOffering, ClassSession, SolverConfig 
} from './types';
import { SchedulerFilterSidebar } from './components/SchedulerFilterSidebar';
import { OptimizationPanel } from './components/OptimizationPanel';

// Constants
const DAY_START_HOUR = 8;
const DAY_END_HOUR = 20;
const TOTAL_MINUTES = (DAY_END_HOUR - DAY_START_HOUR) * 60;

type ViewType = 'timeline' | 'utilization';
type TimeRange = 'day' | 'week' | 'month' | 'quarter';
type GroupBy = 'room' | 'instructor' | 'course';

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

export default function SchedulerPage() {
  // --- Data State ---
  const [terms, setTerms] = useState<AcademicTerm[]>([]);
  const [selectedTerm, setSelectedTerm] = useState<AcademicTerm | null>(null);
  
  const [buildings, setBuildings] = useState<BuildingType[]>([]);
  const [rooms, setRooms] = useState<Room[]>([]);
  const [offerings, setOfferings] = useState<CourseOffering[]>([]);
  const [sessions, setSessions] = useState<ClassSession[]>([]);
  
  const [loading, setLoading] = useState(true);

  // --- UI State ---
  const [currentDate, setCurrentDate] = useState(new Date());
  const [viewType, setViewType] = useState<ViewType>('timeline');
  const [timeRange, setTimeRange] = useState<TimeRange>('day');
  const [groupBy, setGroupBy] = useState<GroupBy>('room');
  const [zoomLevel, setZoomLevel] = useState(2.5);
  
  // --- Filter State ---
  const [isFilterOpen, setIsFilterOpen] = useState(false);
  const [filters, setFilters] = useState<{departments: string[], floors: string[]}>({
     departments: [],
     floors: []
  });

  // --- Optimizer State ---
  const [isOptimizerRunning, setIsOptimizerRunning] = useState(false);
  const [optimizerConfig, setOptimizerConfig] = useState<SolverConfig>({ 
    max_iterations: 1000,
    utilization_weight: 60, 
    satisfaction_weight: 40, 
    prioritize_buildings: true,
    enable_department_constraints: true
  });

  // --- Initial Load ---
  useEffect(() => {
    async function loadInitialData() {
      try {
        const [termsData, buildingsData, roomsData] = await Promise.all([
          fetchTerms(),
          fetchBuildings(),
          fetchRooms()
        ]);
        
        setTerms(termsData);
        setBuildings(buildingsData);
        setRooms(roomsData);
        
        // Select active term or first term
        const active = termsData.find(t => t.is_active) || termsData[0];
        if (active) setSelectedTerm(active);

      } catch (err) {
        console.error("Failed to load initial data", err);
      } finally {
        setLoading(false);
      }
    }
    loadInitialData();
  }, []);

  // --- Fetch Offers & Sessions when Term changes ---
  useEffect(() => {
    if (!selectedTerm) return;
    
    async function loadTermData() {
        setLoading(true);
        try {
            const offs = await fetchOfferings(selectedTerm!.id);
            setOfferings(offs);
            
            // For now fetch all sessions for term (optimize later with date range)
            const sess = await fetchSessions(); // TODO: Filter by term via backend or here
            // Client-side filter for now if API doesn't support term_id for sessions directly yet (it maps via offerings)
            // Ideally we pass date range to fetchSessions
            setSessions(sess); 
        } catch (e) {
            console.error(e);
        } finally {
            setLoading(false);
        }
    }
    loadTermData();
  }, [selectedTerm]);

  // --- Derived State: Date Range ---
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
        dayList = eachDayOfInterval({ start, end }); 
     } else if (timeRange === 'quarter') {
        start = startOfMonth(currentDate);
        end = addMonths(start, 3);
        dayList = []; 
     }

     return { startDate: start, endDate: end, days: dayList };
  }, [currentDate, timeRange]);

  // --- Filtering Logic ---
  const filteredRooms = useMemo(() => {
     return rooms.filter(r => {
        if (filters.floors.length > 0) {
           const floorKey = `${r.building_id}_${r.floor}`;
           if (!filters.floors.includes(floorKey)) return false;
        }
        // Department filter TODO: Need departments in Room model on frontend or joined
        return true;
     });
  }, [rooms, filters]);

  // --- Grouping Logic ---
  const groupedResources = useMemo(() => {
     const groups: Record<string, { id: string, title: string, items: any[] }> = {};

     if (groupBy === 'room') {
        filteredRooms.forEach(r => {
            // Find building name
           const bName = buildings.find(b => b.id === r.building_id)?.name || 'Unknown';
           const key = `${bName} - Floor ${r.floor}`;
           if (!groups[key]) groups[key] = { id: key, title: key, items: [] };
           groups[key].items.push({...r, buildingName: bName});
        });
     } else if (groupBy === 'instructor') {
         // TODO: Implement instructor grouping if we have instructor data
     } else {
        // Group by Course
        const uniqueCourses = Array.from(new Set(offerings.map(o => o.course_id))); // Should be Title ideally
        groups['Courses'] = { id: 'courses', title: 'Active Courses', items: uniqueCourses.map(id => ({ id, name: id, type: 'course' })) };
     }
     return groups;
  }, [filteredRooms, offerings, groupBy, buildings]);

  // --- Optimizer Handler ---
  const handleOptimize = async () => {
    if (!selectedTerm) return;
    setIsOptimizerRunning(true);
    try {
        await optimizeSchedule(selectedTerm.id, optimizerConfig);
        // Refresh sessions
        const sess = await fetchSessions();
        setSessions(sess);
    } catch (e) {
        console.error("Optimization failed", e);
        alert("Optimization failed. See console.");
    } finally {
        setIsOptimizerRunning(false);
    }
  };

  const handleNav = (dir: 'prev' | 'next') => {
     if (timeRange === 'day') setCurrentDate(addDays(currentDate, dir === 'next' ? 1 : -1));
     if (timeRange === 'week') setCurrentDate(addWeeks(currentDate, dir === 'next' ? 1 : -1));
     if (timeRange === 'month') setCurrentDate(addMonths(currentDate, dir === 'next' ? 1 : -1));
     if (timeRange === 'quarter') setCurrentDate(addMonths(currentDate, dir === 'next' ? 3 : -3));
  };
  
  // Calculate mock score for now
  const globalScore = 85; 

  if (loading && !terms.length) return <div className="flex items-center justify-center h-full"><Loader2 className="animate-spin" /></div>;

  return (
    <div className="flex h-full bg-slate-100 relative">
       {/* Filter Sidebar */}
       <SchedulerFilterSidebar 
         isOpen={isFilterOpen} 
         onClose={() => setIsFilterOpen(false)}
         departments={[]} // TODO
         buildings={buildings}
         filters={filters}
         setFilters={setFilters}
       />

       {/* Main Content */}
       <div className="flex-1 flex flex-col min-w-0 transition-all duration-300" style={{ marginLeft: isFilterOpen ? '256px' : '0' }}>
          
          {/* Top Toolbar */}
          <div className="h-16 bg-white border-b border-slate-200 px-6 flex items-center justify-between sticky top-0 z-30 shadow-sm">
             <div className="flex items-center gap-4">
                <Button variant="ghost" size="icon" onClick={() => setIsFilterOpen(!isFilterOpen)} className={isFilterOpen ? "bg-indigo-100 text-indigo-600" : ""}>
                    <Filter size={16} />
                </Button>
                
                {/* Term Selector */}
                <select 
                    value={selectedTerm?.id || ''} 
                    onChange={(e) => setSelectedTerm(terms.find(t => t.id === e.target.value) || null)}
                    className="text-sm font-bold border-none bg-slate-100 rounded px-2 py-1"
                >
                    {terms.map(t => <option key={t.id} value={t.id}>{t.name}</option>)}
                </select>

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
                
                 <div className="flex bg-slate-100 p-1 rounded-lg">
                   {['day', 'week', 'month'].map(r => (
                      <button
                        key={r}
                        onClick={() => { setTimeRange(r as TimeRange); if(r==='month') setViewType('utilization'); else setViewType('timeline'); }}
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
                      <option value="timeline">Timeline</option>
                      <option value="utilization">Utilization</option>
                   </select>
                </div>

                <div className="flex items-center gap-2">
                   <span className="text-[10px] font-bold text-slate-400 uppercase">Group By:</span>
                   <select 
                     value={groupBy} 
                     onChange={(e) => setGroupBy(e.target.value as GroupBy)}
                     className="bg-transparent text-xs font-bold text-slate-700 outline-none cursor-pointer border border-slate-200 rounded px-2 py-1 hover:border-indigo-300"
                   >
                      <option value="room">Room</option>
                      <option value="course">Course</option>
                   </select>
                </div>
                
                <Button size="sm" variant="secondary" onClick={() => setViewType(viewType === 'timeline' ? 'utilization' : 'timeline')}>
                   <BarChart3 size={16} className="mr-2"/>
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
                         {groupBy === 'room' ? 'Location' : 'Resource'}
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
                            {groupBy === 'room' ? <Building size={10} /> : <BookOpen size={10} />}
                            {groupName}
                         </div>
                         
                         {groupData.items.map((resource: any) => {
                            // Filter events for this resource row
                            const rowEvents = sessions.filter(e => {
                               if (groupBy === 'room') return e.room_id === resource.id;
                               // if (groupBy === 'course') return e.course_offering_id === resource.id;
                               return false;
                            }).filter(e => isSameDay(new Date(e.date), currentDate)); 

                            return (
                               <div key={resource.id} className="flex border-b border-slate-200 bg-white hover:bg-slate-50/50 transition-colors h-20 relative group/row">
                                  <div className="w-56 flex-shrink-0 border-r border-slate-200 p-3 sticky left-0 z-10 bg-white group-hover/row:bg-slate-50 transition-colors flex flex-col justify-center">
                                     <div className="font-bold text-slate-800 text-sm truncate">{resource.name}</div>
                                     <div className="text-[10px] text-slate-400">Cap: {resource.capacity}</div>
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
                                        const [h, m] = (event.start_time || '09:00').split(':').map(Number);
                                        const startMins = (h - DAY_START_HOUR) * 60 + m;
                                        // TODO: Calculate duration from start/end
                                        const [endH, endM] = (event.end_time || '09:00').split(':').map(Number);
                                        const durationMinutes = ((endH * 60 + endM) - (h * 60 + m));

                                        return (
                                           <div 
                                             key={event.id}
                                             className={cn(
                                                "absolute top-2 bottom-2 rounded-lg border shadow-sm p-2 text-[10px] overflow-hidden hover:z-20 hover:shadow-md cursor-pointer transition-all bg-indigo-50 border-indigo-200 text-indigo-900"
                                             )}
                                             style={{ left: startMins * zoomLevel, width: durationMinutes * zoomLevel }}
                                          >
                                             <div className="font-bold truncate">{event.title}</div>
                                             <div className="opacity-70 truncate">{event.start_time} - {event.end_time}</div>
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
             
            {/* Utilization View placeholder logic (can be expanded) */}
             {viewType === 'utilization' && (
                 <div className="p-8 text-slate-500">Utilization view not yet fully implemented.</div>
             )}

          </div>

          {/* Floating Optimizer Card */}
          {selectedTerm && (
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
}
