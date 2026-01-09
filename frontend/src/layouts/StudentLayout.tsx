import { Link, Outlet } from 'react-router-dom'
import { RoleSwitcher } from '@/components/common/RoleSwitcher'
import { useTranslation } from 'react-i18next'
import { BookOpen, Calendar, GraduationCap, LayoutDashboard } from 'lucide-react'

export const StudentLayout: React.FC = () => {
  const { t } = useTranslation('common')
  return (
    <div className="flex h-screen bg-gray-50">
      {/* Sidebar Placeholder */}
      <aside className="w-64 bg-white border-r p-4 hidden md:block flex-col">
        <h2 className="text-xl font-bold mb-6 text-green-600 px-2">{t('layouts.student.title')}</h2>
        <nav className="space-y-1">
            <Link to="/student/dashboard" className="flex items-center gap-3 p-3 rounded-xl hover:bg-green-50 text-slate-600 hover:text-green-700 font-medium transition-colors">
              <LayoutDashboard size={20} />
              {t('layouts.student.dashboard')}
            </Link>
            <Link to="/student/courses" className="flex items-center gap-3 p-3 rounded-xl hover:bg-green-50 text-slate-600 hover:text-green-700 font-medium transition-colors">
              <BookOpen size={20} />
              {t('layouts.student.learning')}
            </Link>
            <Link to="/calendar" className="flex items-center gap-3 p-3 rounded-xl hover:bg-green-50 text-slate-600 hover:text-green-700 font-medium transition-colors">
              <Calendar size={20} />
              {t('layouts.student.schedule')}
            </Link>
            <Link to="/student/grades" className="flex items-center gap-3 p-3 rounded-xl hover:bg-green-50 text-slate-600 hover:text-green-700 font-medium transition-colors">
              <GraduationCap size={20} />
              {t('layouts.student.grades')}
            </Link>
        </nav>
      </aside>

      <main className="flex-1 flex flex-col overflow-hidden">
        <header className="h-16 bg-white/80 backdrop-blur-md border-b flex items-center justify-between px-8 sticky top-0 z-10">
            <h1 className="text-lg font-semibold">{t('layouts.student.dashboard')}</h1>
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
