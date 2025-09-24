import React from 'react'

type Step = { code:string; title:string; status?: 'todo'|'in_progress'|'submitted'|'needs_changes'|'done' }
export function VerticalProgress({ steps }: { steps: Step[] }) {
  return (
    <ol className="relative border-s pl-6">
      {steps.map((s, i) => (
        <li key={s.code} className="mb-6 ms-4">
          <span className={`absolute -start-3 flex h-6 w-6 items-center justify-center rounded-full border text-xs ${s.status==='done' ? 'bg-green-500 text-white' : s.status==='submitted' ? 'bg-blue-500 text-white' : s.status==='needs_changes' ? 'bg-rose-500 text-white' : 'bg-white'}`}>{i+1}</span>
          <h3 className="font-medium">{s.code} — {s.title}</h3>
          <p className="text-xs text-gray-500">{s.status || 'todo'}</p>
        </li>
      ))}
    </ol>
  )
}
