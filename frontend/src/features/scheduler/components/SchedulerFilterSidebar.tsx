import React, { useState } from 'react';
import { Filter, X, Building as BuildingIcon } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Department, Building } from '../types';

interface SchedulerFilterSidebarProps {
  isOpen: boolean;
  onClose: () => void;
  departments: Department[];
  buildings: Building[];
  filters: {
    departments: string[];
    floors: string[];
  };
  setFilters: (f: any) => void;
}

export const SchedulerFilterSidebar: React.FC<SchedulerFilterSidebarProps> = ({ 
  isOpen, 
  onClose,
  departments,
  buildings,
  filters,
  setFilters 
}) => {
  if (!isOpen) return null;

  return (
    <div className="absolute top-0 left-0 bottom-0 w-64 bg-white border-r border-slate-200 z-40 shadow-xl overflow-y-auto">
       <div className="p-4 border-b border-slate-100 flex justify-between items-center">
          <h3 className="font-bold text-slate-800 text-sm flex items-center gap-2"><Filter size={16} /> Filters</h3>
          <Button variant="ghost" size="icon" onClick={onClose} className="h-6 w-6"><X size={16} /></Button>
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
                      {dept.color && <span className="w-2 h-2 rounded-full ml-auto" style={{ backgroundColor: dept.color }} />}
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
                         <BuildingIcon size={12} className="text-slate-400" /> {b.name}
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
