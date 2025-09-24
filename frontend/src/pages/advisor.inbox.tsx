import React, { useEffect, useState } from 'react'
import { api } from '../api/client'
import { motion } from "framer-motion"

type Row = { student_id:string; student_name:string; step_id:string; step_code:string; step_title:string }

export function AdvisorInbox() {
  const [rows, setRows] = useState<Row[]>([])
  useEffect(()=>{ api('/advisor/inbox').then(setRows) },[])
  return (
    <div className="mt-6">
      <h2 className="text-xl font-semibold mb-2">Advisor Inbox</h2>
      {rows.length===0 ? <p className="text-sm text-gray-600">Nothing pending.</p> : (
        <div className="space-y-2">
          {rows.map((r,i)=> (
            <motion.div key={i} className="border rounded p-3 shadow-sm" initial={{opacity:0, y:4}} animate={{opacity:1, y:0}}>
              <div className="font-medium">{r.student_name}</div>
              <div className="text-sm">{r.step_code}: {r.step_title}</div>
              <div className="mt-2 space-x-2">
                <button className="text-sm underline" onClick={async()=>{ await api(`/reviews/${r.student_id}/steps/${r.step_id}/return`, {method:"POST", body: JSON.stringify({ comment: "Please revise", mentions: [] })}); location.reload() }}>Return with comments</button>
                <button className="text-sm underline" onClick={async()=>{ await api(`/reviews/${r.student_id}/steps/${r.step_id}/approve`, {method:"POST", body: JSON.stringify({ comment: "Approved" })}); location.reload() }}>Approve</button>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}
