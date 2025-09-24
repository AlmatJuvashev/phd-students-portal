import React, { useRef, useState } from 'react'
import { api } from '../api/client'

export function PresignedUploader({ docId }:{docId:string}) {
  const [status, setStatus] = useState('')
  const input = useRef<HTMLInputElement|null>(null)
  async function upload() {
    const f = input.current?.files?.[0]
    if (!f) return
    const presign = await api(`/documents/${docId}/presign`, {
      method:'POST',
      body: JSON.stringify({ filename: f.name, content_type: f.type || 'application/octet-stream' })
    })
    const res = await fetch(presign.url, { method:'PUT', body: f, headers: { 'Content-Type': f.type } })
    if (!res.ok) throw new Error('upload failed')
    setStatus('Uploaded via S3')
  }
  return (
    <div className="text-sm space-y-2">
      <input type="file" ref={input} accept=".pdf,.docx" />
      <button className="px-3 py-1 rounded border" onClick={upload}>Upload</button>
      <span className="ml-2 text-gray-600">{status}</span>
    </div>
  )
}
