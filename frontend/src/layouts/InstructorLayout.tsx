import { Link, Outlet } from 'react-router-dom'
import { RoleSwitcher } from '@/components/common/RoleSwitcher'
import { useTranslation } from 'react-i18next'
import { BookOpen, CheckSquare, LayoutDashboard, Users } from 'lucide-react'

export const InstructorLayout: React.FC = () => {
  const { t } = useTranslation('common')
  return (
    <div className="flex h-screen bg-gray-100">
      {/* Sidebar Placeholder */}
      <aside className="w-64 bg-white border-r p-4 hidden md:block flex-col">
        <h2 className="text-xl font-bold mb-6 text-blue-600 px-2">{t('layouts.instructor.title')}</h2>
        <nav className="space-y-1">
            <Link to="/teach/dashboard" className="flex items-center gap-3 p-3 rounded-xl hover:bg-blue-50 text-slate-600 hover:text-blue-700 font-medium transition-colors">
              <LayoutDashboard size={20} />
              {t('layouts.instructor.dashboard')}
            </Link>
            <Link to="/teach/courses" className="flex items-center gap-3 p-3 rounded-xl hover:bg-blue-50 text-slate-600 hover:text-blue-700 font-medium transition-colors">
              <BookOpen size={20} />
              {t('layouts.instructor.courses')}
            </Link>
            <Link to="/teach/grading" className="flex items-center gap-3 p-3 rounded-xl hover:bg-blue-50 text-slate-600 hover:text-blue-700 font-medium transition-colors">
              <CheckSquare size={20} />
              {t('layouts.instructor.grading')}
            </Link>
             {/* Attendance is typically course-specific, linking to courses for now as entry point */}
            <Link to="/teach/courses" className="flex items-center gap-3 p-3 rounded-xl hover:bg-blue-50 text-slate-600 hover:text-blue-700 font-medium transition-colors">
              <Users size={20} />
              {t('layouts.instructor.attendance')}
            </Link>
        </nav>
      </aside>

      <main className="flex-1 flex flex-col overflow-hidden">
        <header className="h-16 bg-white border-b flex items-center justify-between px-6">
            <h1 className="text-lg font-semibold">{t('layouts.instructor.dashboard')}</h1>
            <div className="flex items-center">
                <RoleSwitcher />
            </div>
        </header>
        <div className="flex-1 overflow-auto p-6">
            <Outlet />
        </div>
      </main>
    </div>
  )
}
