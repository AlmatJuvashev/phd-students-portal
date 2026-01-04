import React, { useState, useEffect } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Search, BookOpen, Loader2 } from 'lucide-react';
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter } from '@/components/ui/dialog';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { ScrollArea } from '@/components/ui/scroll-area';
import { getCourses } from '@/features/curriculum/api'; // Adjust path if needed
import { Course } from '@/features/curriculum/types';

interface CoursePickerDialogProps {
  isOpen: boolean;
  onClose: () => void;
  onSelect: (course: Course) => void;
  programId?: string; // Optional: filter by program
}

export const CoursePickerDialog: React.FC<CoursePickerDialogProps> = ({ 
  isOpen, 
  onClose, 
  onSelect,
  programId 
}) => {
  const [searchQuery, setSearchQuery] = useState('');
  
  const { data: courses, isLoading } = useQuery({
    queryKey: ['courses', programId],
    queryFn: async () => {
      const res = await getCourses(programId);
      // Handle both cases: AxiosResponse or direct data
      return (res as any).data || res; 
    },
    enabled: isOpen,
  });

  const filteredCourses = (courses as Course[])?.filter(c => 
    c.title.toLowerCase().includes(searchQuery.toLowerCase()) || 
    c.code.toLowerCase().includes(searchQuery.toLowerCase())
  );

  return (
    <Dialog open={isOpen} onOpenChange={onClose}>
      <DialogContent className="max-w-2xl max-h-[80vh] flex flex-col">
        <DialogHeader>
          <DialogTitle>Select a Course</DialogTitle>
        </DialogHeader>
        
        <div className="relative mb-4">
          <Search className="absolute left-3 top-2.5 h-4 w-4 text-slate-400" />
          <Input 
            placeholder="Search courses..." 
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="pl-9"
          />
        </div>

        <ScrollArea className="flex-1 h-[400px]">
          {isLoading ? (
             <div className="flex items-center justify-center h-40">
                <Loader2 className="animate-spin text-indigo-600" />
             </div>
          ) : (
            <div className="grid grid-cols-1 gap-2">
              {filteredCourses?.map(course => (
                 <button
                   key={course.id}
                   onClick={() => onSelect(course)}
                   className="flex items-center p-3 rounded-lg border border-slate-100 hover:border-indigo-200 hover:bg-indigo-50 transition-all text-left group"
                 >
                    <div className="w-10 h-10 rounded-lg bg-indigo-100 text-indigo-600 flex items-center justify-center mr-3 group-hover:bg-white group-hover:shadow-sm">
                       <BookOpen size={20} />
                    </div>
                    <div>
                       <div className="font-bold text-slate-900 group-hover:text-indigo-700">{course.title}</div>
                       <div className="text-xs text-slate-500">{course.code} • {course.credits} Credits</div>
                    </div>
                 </button>
              ))}
              {filteredCourses?.length === 0 && (
                 <div className="text-center py-10 text-slate-500">
                    No courses found matching "{searchQuery}"
                 </div>
              )}
            </div>
          )}
        </ScrollArea>

        <DialogFooter>
          <Button variant="outline" onClick={onClose}>Cancel</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
};
