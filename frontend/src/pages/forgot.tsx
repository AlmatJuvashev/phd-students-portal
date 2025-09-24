import React, { useState } from 'react'
import { api } from '../api/client'

export function ForgotPassword() {
  const [email, setEmail] = useState('')
  const [sent, setSent] = useState(false)

  async function submit(e: React.FormEvent) {
    e.preventDefault()
    await api('/auth/forgot', { method: 'POST', body: JSON.stringify({ email }) })
    setSent(true)
  }

  return (
    <div className="max-w-sm mx-auto mt-10">
      <h2 className="text-xl font-semibold mb-4">Forgot password</h2>
      {sent ? <p>Check your email for reset link (Mailpit in dev).</p> : (
        <form className="space-y-3" onSubmit={submit}>
          <input className="w-full border p-2 rounded" placeholder="Email" value={email} onChange={(e)=>setEmail(e.target.value)} />
          <button className="w-full bg-black text-white p-2 rounded">Send reset link</button>
        </form>
      )}
    </div>
  )
}
