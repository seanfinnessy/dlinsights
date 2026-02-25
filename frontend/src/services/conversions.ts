export const formatDuration = (seconds: number) => {
  const h = Math.floor(seconds / 3600)
  const m = Math.floor((seconds % 3600) / 60)
  const s = seconds % 60
  const mm = String(m).padStart(2, '0')
  const ss = String(s).padStart(2, '0')
  return h > 0 ? `${h}:${mm}:${ss}` : `${mm}:${ss}`
}

export const formatStartTime = (epoch: number) => {
  const d = new Date(epoch * 1000)
  const date = d.toLocaleString('en-US', { month: 'short', day: 'numeric' })
  const time = d.toLocaleString('en-US', { hour: 'numeric', minute: '2-digit', hour12: true })
  return { date, time }
}

export const formatGameMode = (mode: number): { label: string; style: 'sb' | 'comp' } | null => {
  if (mode === 4) return { label: 'SB', style: 'sb' }
  if (mode === 1) return { label: 'COMP', style: 'comp' }
  return { label: 'N/A', style: 'sb' }
}

export const formatSouls = (souls: number) => {
  return parseFloat((souls / 1000).toPrecision(3)) + 'k'
}