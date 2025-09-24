import React, { useState, useEffect } from 'react'
import { api } from '../api/client'

export function ResetPassword() {
  const [token, setToken] = useState('')
  const [newPw, setNewPw] = useState('')
  const [done, setDone] = useState(false)
  useEffect(()=>{
    const url = new URL(location.href)
    setToken(url.searchParams.get('token') || '')
  },[])
  async function submit(e: React.FormEvent) {
    e.preventDefault()
    await api('/auth/reset', { method: 'POST', body: JSON.stringify({ token, new_password: newPw }) })
    setDone(true)
  }
  return (
    <div className="max-w-sm mx-auto mt-10">
      <h2 className="text-xl font-semibold mb-4">Reset password</h2>
      {done ? <p>Password updated. You can now <a className="underline" href="/login">login</a>.</p> : (
        <form className="space-y-3" onSubmit={submit}>
          <input className="w-full border p-2 rounded" placeholder="New password" type="password" value={newPw} onChange={(e)=>setNewPw(e.target.value)} />
          <button className="w-full bg-black text-white p-2 rounded">Update password</button>
        </form>
      )}
    </div>
  )
}
