import React from 'react'

export function Dashboard() {
  return (
    <div className="mt-8 space-y-2">
      <h2 className="text-xl font-semibold">Dashboard</h2>
      <p className="text-sm text-gray-600">This is a placeholder. Add journey progress, recent uploads, and advisor feedback here.</p>
      <div className='mt-4 text-sm'>
        <a className='underline' href='/journey'>Open Journey Map</a>
      </div>
    </div>
  )
}
