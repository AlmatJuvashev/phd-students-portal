import React from 'react'
import { Outlet } from 'react-router-dom'
import { RoleSwitcher } from '@/components/common/RoleSwitcher'

export const StudentLayout: React.FC = () => {
  return (
    <div className="flex h-screen bg-gray-50">
      {/* Sidebar Placeholder */}
      <aside className="w-64 bg-white border-r p-4 hidden md:block">
        <h2 className="text-xl font-bold mb-6 text-green-600">Student Portal</h2>
        <nav className="space-y-2">
            <a href="#" className="block p-2 rounded hover:bg-gray-100">My Learning</a>
            <a href="#" className="block p-2 rounded hover:bg-gray-100">Schedule</a>
            <a href="#" className="block p-2 rounded hover:bg-gray-100">Grades</a>
        </nav>
      </aside>

      <main className="flex-1 flex flex-col overflow-hidden">
        <header className="h-16 bg-white border-b flex items-center justify-between px-6">
            <h1 className="text-lg font-semibold">My Dashboard</h1>
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
