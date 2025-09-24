import React, { useEffect, useState } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { VerticalProgress } from '../components/VerticalProgress'
import { api } from '../api/client'

type Module = { id:string; code:string; title:string; sort_order:number }
type Step = { id:string; code:string; title:string; requires_upload:boolean; sort_order:number }

export function StudentChecklist() {
  const [mods, setMods] = useState<Module[]>([])
  const [stepsByMod, setStepsByMod] = useState<Record<string, Step[]>>({})
  useEffect(()=>{
    api('/checklist/modules').then(setMods)
  },[])
  async function openModule(code:string) {
    const steps = await api(`/checklist/steps?module=${encodeURIComponent(code)}`)
    setStepsByMod(s=>({...s, [code]: steps}))
  }
  return (
    <div className="mt-6 space-y-4">
      <h2 className="text-xl font-semibold">Checklist</h2>
      {mods.map(m=> (
        <div key={m.id} className="border rounded">
          <button className="w-full text-left p-3 font-medium bg-gray-50" onClick={()=>openModule(m.code)}>
            {m.code}. {m.title}
          </button>
          <AnimatePresence>
            {stepsByMod[m.code] && (
              <motion.div initial={{height:0, opacity:0}} animate={{height:"auto", opacity:1}} exit={{height:0, opacity:0}} className="p-3 space-y-2">
            {(stepsByMod[m.code]||[]).map(st=> (
              <div key={st.id} className="flex items-center justify-between border rounded p-2 shadow-sm">
                <div>
                  <div className="font-medium">{st.code}: {st.title}</div>
                  {st.requires_upload && <div className="text-xs text-gray-500">Requires upload (.pdf/.docx)</div>}
                </div>
                
                <button
                  className="text-sm underline"
                  onClick={async ()=>{
                    // naive create-or-open: create doc per step code (idempotent enough for demo)
                    const res = await api(`/students/${localStorage.getItem('student_id')||'me'}/documents`, { method:'POST', body: JSON.stringify({ kind: 'other', title: st.title }) })
                    location.href = `/documents/${res.id}`
                  }}
                >Open</button>

              </div>
            ))}
              </motion.div>
            )}
          </AnimatePresence>
        </div>
      ))}
      <div className='mt-6'>
        {/* Demo: vertical progress for current module if loaded */}
        {mods[0] && stepsByMod[mods[0].code] && (
          <VerticalProgress steps={stepsByMod[mods[0].code].map((st:any)=>({ code: st.code, title: st.title, status: st.status||'todo' }))} />
        )}
      </div>
    </div>
  )
}
