import React from 'react'
import { Button } from './ui/button'

export type Mention = { id: string, name: string }

export function MentionChips({ list, onRemove }:{ list: Mention[], onRemove: (id:string)=>void }) {
  return (
    <div className="flex flex-wrap gap-1">
      {list.map(m => (
        <span key={m.id} className="inline-flex items-center gap-1 rounded-full border px-2 py-0.5 text-xs bg-white">
          <span className="inline-flex h-4 w-4 items-center justify-center rounded-full bg-gray-800 text-white text-[10px]">
            {m.name.split(' ').map(x=>x[0]).join('').slice(0,2).toUpperCase()}
          </span>
          @{m.name}
          <button className="text-gray-500 hover:text-gray-800" onClick={()=>onRemove(m.id)}>×</button>
        </span>
      ))}
    </div>
  )
}

export function MentionResults({ list, onPick }:{ list: Mention[], onPick:(m:Mention)=>void }) {
  return (
    <div className="mt-1 max-h-40 overflow-auto rounded border bg-white shadow">
      {list.map(m => (
        <button key={m.id} className="flex w-full items-center gap-2 px-2 py-1 text-left hover:bg-gray-50" onClick={()=>onPick(m)}>
          <span className="inline-flex h-6 w-6 items-center justify-center rounded-full bg-gray-800 text-white text-xs">
            {m.name.split(' ').map(x=>x[0]).join('').slice(0,2).toUpperCase()}
          </span>
          <span className="text-sm">@{m.name}</span>
        </button>
      ))}
      {list.length===0 && <div className="px-2 py-1 text-xs text-gray-500">No matches</div>}
    </div>
  )
}
